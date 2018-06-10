package app

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/lgarithm/az/arm/tpl"
)

const (
	n      = 4
	vmSize = compute.VirtualMachineSizeTypesStandardNV6
)

func New(cloudInitScript string) *tpl.Builder {
	b := tpl.NewBuilder()
	rules := []network.SecurityRule{
		tpl.NewAllowInboundRule("allow-ssh", "22", 1000),
	}
	nsg := b.AddNSG(rules...)
	vn := b.AddVN()
	opts := tpl.DefaultVMOptions()
	opts.CloudInitScript = cloudInitScript
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("node-%02d", i)
		ip := b.AddIP(name + "-ip")
		ni := b.AddNI(name+"-nic", vn, &nsg, &ip)
		vm := b.AddVM(name, ni, &opts)
		vm.HardwareProfile.VMSize = vmSize
		// b.AddVMExt("provision", vm)
		// b.AddDockerVMExt("docker", vm)
	}
	return b
}
