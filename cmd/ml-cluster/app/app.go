package app

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/lgarithm/az/arm/tpl"
)

const (
	numWorkers = 3
	vmSize     = compute.StandardNV6

	imageGroup = `teavana-image`
	imageName  = `gpu-mpi-3-1`
)

var (
	vmImage = &compute.ImageReference{
		// must be in the same location as VM
		// https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-template-functions-resource#resourceid
		ID: to.StringPtr(fmt.Sprintf(`[resourceId('%s', '%s', '%s')]`, imageGroup, `Microsoft.Compute/images`, imageName)),
	}
)

func New(cloudInitScript string) *tpl.Builder {
	b := tpl.NewBuilder()
	vnet := b.AddVN("default-vnet")

	// relay VM
	{
		rules := []network.SecurityRule{
			tpl.NewAllowInboundRule("allow-ssh", "22", 1000),
		}
		nsg := b.AddNSG("relay-nsg", rules...)
		opts := tpl.DefaultVMOptions()
		name := `relay`
		ip := b.AddIP(name + "-ip")
		ni := b.AddNI(name+"-nic", vnet, &nsg, &ip)
		b.AddVM(name, ni, &opts)
	}

	// worker VMs
	{
		opts := tpl.DefaultVMOptions()
		opts.AllowPassword = true
		opts.CloudInitScript = cloudInitScript
		nsg := b.AddNSG("internal-nsg")
		for i := 0; i < numWorkers; i++ {
			name := fmt.Sprintf("node-%02d", i)
			ni := b.AddNI(name+"-nic", vnet, &nsg, nil)
			vm := b.AddVM(name, ni, &opts)
			vm.HardwareProfile.VMSize = vmSize
			vm.StorageProfile.ImageReference = vmImage
		}
	}
	return b
}
