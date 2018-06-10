package tpl

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
)

// VirtualMachineExtension is an enum
// az vm extension image list --location southeastasia -o table
// https://github.com/Azure/azure-linux-extensions
type VirtualMachineExtension string

const (
	CustomScript    VirtualMachineExtension = "CustomScript"
	DockerExtension VirtualMachineExtension = "DockerExtension"
)

func newVMExt(name string, vm VirtualMachineResource) compute.VirtualMachineExtension {
	cmd := "pwd"
	settings := map[string]interface{}{
		"commandToExecute": cmd,
	}
	protectedSettings := map[string]interface{}{
		// "commandToExecute": cmd,
	}
	// A nested resource type must have identical number of segments as its resource name.
	// Please see https://aka.ms/arm-template/#resources for usage details.
	fullname := fmt.Sprintf("%s/%s", *vm.Name, name)
	return newVMExtFor(fullname, CustomScript, settings, protectedSettings)
}

func newDockerVMExt(name string, vm VirtualMachineResource) compute.VirtualMachineExtension {
	settings := map[string]interface{}{}
	protectedSettings := map[string]interface{}{}
	// A nested resource type must have identical number of segments as its resource name.
	// Please see https://aka.ms/arm-template/#resources for usage details.
	fullname := fmt.Sprintf("%s/%s", *vm.Name, name)
	return newVMExtFor(fullname, DockerExtension, settings, protectedSettings)
}

var handlerVersions = map[VirtualMachineExtension]string{
	CustomScript:    "2.0",
	DockerExtension: "1.2",
}

func newVMExtFor(fullname string, extType VirtualMachineExtension, settings map[string]interface{}, protectedSettings map[string]interface{}) compute.VirtualMachineExtension {
	return compute.VirtualMachineExtension{
		Type:     to.StringPtr(TypeVMExt),
		Name:     to.StringPtr(fullname),
		Location: to.StringPtr("[resourceGroup().location]"),
		VirtualMachineExtensionProperties: &compute.VirtualMachineExtensionProperties{
			Publisher:               to.StringPtr("Microsoft.Azure.Extensions"),
			Type:                    to.StringPtr(string(extType)),
			TypeHandlerVersion:      to.StringPtr(handlerVersions[extType]),
			AutoUpgradeMinorVersion: to.BoolPtr(true),
			Settings:                &settings,
			ProtectedSettings:       &protectedSettings,
		},
	}
}

// https://docs.microsoft.com/en-us/azure/virtual-machines/linux/agent-user-guide
// /var/lib/waagent
