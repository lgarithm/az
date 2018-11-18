package tpl

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/go-autorest/autorest/to"
)

func newNI(name string, vn VirtualNetworkResource, nsg *NetworkSecurityGroupResource, ip *PublicIPResource) network.Interface {
	var ipAddress *network.PublicIPAddress
	if ip != nil {
		ipAddress = &network.PublicIPAddress{
			ID: to.StringPtr(fmt.Sprintf("[resourceId('%s', '%s')]", *ip.Type, *ip.Name)),
		}
	}
	var networkSecurityGroup *network.SecurityGroup
	if nsg != nil {
		networkSecurityGroup = &network.SecurityGroup{
			ID: to.StringPtr(fmt.Sprintf("[resourceId('%s', '%s')]", *nsg.Type, *nsg.Name)),
		}
	}
	return network.Interface{
		Name:     to.StringPtr(name),
		Type:     to.StringPtr(TypeNI),
		Location: to.StringPtr("[resourceGroup().location]"),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				network.InterfaceIPConfiguration{
					Name: to.StringPtr("ipconfig01"),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: network.Dynamic,
						Subnet: &network.Subnet{
							ID: to.StringPtr(fmt.Sprintf("[concat(resourceId('%s', '%s'), '/subnets/%s')]", *vn.Type, *vn.Name, defaultSNet)),
						},
						PublicIPAddress: ipAddress,
					},
				},
			},
			NetworkSecurityGroup: networkSecurityGroup,
		},
	}
}
