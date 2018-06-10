package app

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/lgarithm/az/arm/tpl"
)

const vmName = "master"

func New(cloudInitScript string) *tpl.Builder {
	b := tpl.NewBuilder()
	rules := []network.SecurityRule{
		tpl.NewAllowInboundRule("allow-ssh", "22", 1000),
	}
	nsg := b.AddNSG(rules...)
	vn := b.AddVN()
	ip := b.AddIP(vmName)
	ni := b.AddNI(vmName, vn, &nsg, &ip)
	opts := tpl.DefaultVMOptions()
	opts.CloudInitScript = cloudInitScript
	vm := b.AddVM(vmName, ni, &opts)
	// b.AddVMExt("provision", vm)
	// b.AddDockerVMExt("docker", vm)
	use(vm)
	return b
}

func use(i interface{}) {}
