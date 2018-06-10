package tpl

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	disk "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
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
	Header
	storage.Account
}

type ManagedDiskResource struct {
	Header
	disk.Disk
}

type VirtualNetworkResource struct {
	Header
	network.VirtualNetwork
}

type NetworkSecurityGroupResource struct {
	Header
	network.SecurityGroup
}

type PublicIPResource struct {
	Header
	network.PublicIPAddress
}

type NetworkInterfaceResource struct {
	Header
	network.Interface
}

type VirtualMachineResource struct {
	Header
	compute.VirtualMachine
}

type VirtualMachineExtensionResource struct {
	Header
	compute.VirtualMachineExtension
}

type Builder struct {
	resources []interface{}
	ni2ip     map[string]string
	vm2ni     map[string]string
}

func (h *Header) addDep(ref string) {
	h.DependsOn = append(h.DependsOn, ref)
}

func NewBuilder() *Builder {
	return &Builder{
		ni2ip: map[string]string{},
		vm2ni: map[string]string{},
	}
}

func (b *Builder) Build() DeploymentTemplate {
	return map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters":     &map[string]Parameter{},
		"resources":      &b.resources,
	}
}

func (b *Builder) add(r interface{}) {
	b.resources = append(b.resources, r)
}

func (b *Builder) header(apiVersion string) Header {
	return Header{APIVersion: apiVersion}
}

// AddVN adds a Virtual Network to the template builder.
func (b *Builder) AddVN() VirtualNetworkResource {
	r := VirtualNetworkResource{
		b.header(APIVersions.VirtualNetwork),
		newVN("default"),
	}
	b.add(r)
	return r
}

// AddNSG adds a Network Security Group to the template builder.
func (b *Builder) AddNSG(rules ...network.SecurityRule) NetworkSecurityGroupResource {
	var p *[]network.SecurityRule
	if len(rules) > 0 {
		p = &rules
	}
	r := NetworkSecurityGroupResource{
		b.header(APIVersions.NetworkSecurityGroup),
		newNSG("default", p),
	}
	b.add(r)
	return r
}

// AddIP adds a Public IP Address to the template builder.
func (b *Builder) AddIP(name string) PublicIPResource {
	r := PublicIPResource{
		b.header(APIVersions.PublicIPAddress),
		newIP(name),
	}
	b.add(r)
	return r
}

// AddNI adds a network Interface to the template builder.
func (b *Builder) AddNI(name string, vn VirtualNetworkResource, nsg *NetworkSecurityGroupResource, ip *PublicIPResource) NetworkInterfaceResource {
	h := b.header(APIVersions.NetworkInterface)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *vn.Type, *vn.Name))
	if nsg != nil {
		h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *nsg.Type, *nsg.Name))
	}
	if ip != nil {
		h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *ip.Type, *ip.Name))
	}
	r := NetworkInterfaceResource{
		h,
		newNI(name, vn, nsg, ip),
	}
	b.add(r)
	if ip != nil {
		b.ni2ip[name] = *ip.Name
	}
	return r
}

// AddVM adds a Virtual Machine to the template builder.
func (b *Builder) AddVM(name string, ni NetworkInterfaceResource, opts *VMOptions) VirtualMachineResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *ni.Type, *ni.Name))
	r := VirtualMachineResource{
		h,
		newVM(name, ni, opts),
	}
	b.add(r)
	b.vm2ni[name] = *ni.Name
	return r
}

// AddVMExt adds a Virtual Machine Extension for a VM to the template builder.
func (b *Builder) AddVMExt(name string, vm VirtualMachineResource) VirtualMachineExtensionResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *vm.Type, *vm.Name))
	r := VirtualMachineExtensionResource{
		h,
		newVMExt(name, vm),
	}
	b.add(r)
	return r
}

// AddDockerVMExt adds a Virtual Machine Extension for a VM to the template builder.
func (b *Builder) AddDockerVMExt(name string, vm VirtualMachineResource) VirtualMachineExtensionResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *vm.Type, *vm.Name))
	r := VirtualMachineExtensionResource{
		h,
		newDockerVMExt(name, vm),
	}
	b.add(r)
	return r
}

// AddWindowsVM adds a Virtual Machine with windows image to the template builder.
func (b *Builder) AddWindowsVM(name string, ni NetworkInterfaceResource) VirtualMachineResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *ni.Type, *ni.Name))
	r := VirtualMachineResource{
		h,
		newWindowsVM(name, ni),
	}
	b.add(r)
	return r
}

func (b Builder) GetVM2IP() map[string]string {
	m := map[string]string{}
	for vm, ni := range b.vm2ni {
		if ip := b.ni2ip[ni]; ip != "" {
			m[vm] = ip
		}
	}
	return m
}
