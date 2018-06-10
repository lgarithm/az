package tpl

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	disk "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
)

const (
	defaultSNet = "snet01"

	TypeSA    = `Microsoft.Storage/storageAccounts`
	TypeMD    = `Microsoft.Compute/disks`
	TypeVM    = `Microsoft.Compute/virtualMachines`
	TypeVMExt = `Microsoft.Compute/virtualMachines/extensions`
	TypeVN    = `Microsoft.Network/virtualNetworks`
	TypeSG    = `Microsoft.Network/networkSecurityGroups`
	TypeIP    = `Microsoft.Network/publicIPAddresses`
	TypeNI    = `Microsoft.Network/networkInterfaces`
)

type Parameter struct {
	DefaultValue string `json:"defaultValue"`
	Type         string `json:"type"`
	Value        string `json:"value"`
}

type Header struct {
	APIVersion string   `json:"apiVersion"`
	DependsOn  []string `json:"dependsOn"`
}

func (h *Header) addDep(ref string) {
	h.DependsOn = append(h.DependsOn, ref)
}

type GenericResource struct {
	Header
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location string `json:"location"`
}

func (g *GenericResource) RefExpr() string {
	return fmt.Sprintf("[resourceId('%s', '%s')]", g.Type, g.Name)
}

type StorageAccountResource struct {
	storage.Account
}

type ManagedDiskResource struct {
	disk.Disk
}

type VirtualNetworkResource struct {
	network.VirtualNetwork
}

type NetworkSecurityGroupResource struct {
	network.SecurityGroup
}

type PublicIPResource struct {
	network.PublicIPAddress
}

type NetworkInterfaceResource struct {
	network.Interface
}

type VirtualMachineResource struct {
	compute.VirtualMachine
}

type VirtualMachineExtensionResource struct {
	compute.VirtualMachineExtension
}

// DeploymentTemplate can be used by
// az group deployment create --template-file template.json -g $group
// when saved as json file
type DeploymentTemplate map[string]interface{}

func (d DeploymentTemplate) ToDeployment() resources.Deployment {
	value := (map[string]interface{})(d)
	return resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &value,
			Parameters: &map[string]interface{}{},
			Mode:       resources.Incremental,
		},
	}
}
