package tpl

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
)

type Item struct {
	Header Header
	Body   interface{}
}

func toJSON(i interface{}) map[string]interface{} {
	bs, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	objectMap := make(map[string]interface{})
	if err := json.Unmarshal(bs, &objectMap); err != nil {
		return nil
	}
	return objectMap
}

func (i Item) MarshalJSON() ([]byte, error) {
	head := toJSON(i.Header)
	body := toJSON(i.Body)
	objectMap := make(map[string]interface{})
	for k, v := range body {
		objectMap[k] = v
	}
	for k, v := range head {
		objectMap[k] = v
	}
	return json.Marshal(objectMap)
}

type Builder struct {
	items []Item
	ni2ip map[string]string
	vm2ni map[string]string
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
		"resources":      &b.items,
	}
}

func (b *Builder) add(h Header, r interface{}) {
	b.items = append(b.items, Item{h, r})
}

func (b *Builder) header(apiVersion string) Header {
	return Header{APIVersion: apiVersion}
}

// AddVN adds a Virtual Network to the template builder.
func (b *Builder) AddVN(name string) VirtualNetworkResource {
	r := VirtualNetworkResource{
		newVN(name),
	}
	b.add(b.header(APIVersions.VirtualNetwork), r)
	return r
}

// AddNSG adds a Network Security Group to the template builder.
func (b *Builder) AddNSG(name string, rules ...network.SecurityRule) NetworkSecurityGroupResource {
	var p *[]network.SecurityRule
	if len(rules) > 0 {
		p = &rules
	}
	r := NetworkSecurityGroupResource{
		newNSG(name, p),
	}
	b.add(b.header(APIVersions.NetworkSecurityGroup), r)
	return r
}

// AddIP adds a Public IP Address to the template builder.
func (b *Builder) AddIP(name string) PublicIPResource {
	r := PublicIPResource{
		newIP(name),
	}
	b.add(b.header(APIVersions.PublicIPAddress), r)
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
		newNI(name, vn, nsg, ip),
	}
	b.add(h, r)
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
		newVM(name, ni, opts),
	}
	b.add(h, r)
	b.vm2ni[name] = *ni.Name
	return r
}

// AddVMExt adds a Virtual Machine Extension for a VM to the template builder.
func (b *Builder) AddVMExt(name string, vm VirtualMachineResource) VirtualMachineExtensionResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *vm.Type, *vm.Name))
	r := VirtualMachineExtensionResource{
		newVMExt(name, vm),
	}
	b.add(h, r)
	return r
}

// AddDockerVMExt adds a Virtual Machine Extension for a VM to the template builder.
func (b *Builder) AddDockerVMExt(name string, vm VirtualMachineResource) VirtualMachineExtensionResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *vm.Type, *vm.Name))
	r := VirtualMachineExtensionResource{
		newDockerVMExt(name, vm),
	}
	b.add(h, r)
	return r
}

// AddWindowsVM adds a Virtual Machine with windows image to the template builder.
func (b *Builder) AddWindowsVM(name string, ni NetworkInterfaceResource) VirtualMachineResource {
	h := b.header(APIVersions.Default)
	h.addDep(fmt.Sprintf("[resourceId('%s', '%s')]", *ni.Type, *ni.Name))
	r := VirtualMachineResource{
		newWindowsVM(name, ni),
	}
	b.add(h, r)
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
