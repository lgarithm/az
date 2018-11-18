package tpl

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-06-01/network"
	"github.com/Azure/go-autorest/autorest/to"
)

func newVN(name string) network.VirtualNetwork {
	return network.VirtualNetwork{
		Name:     to.StringPtr(name),
		Type:     to.StringPtr(TypeVN),
		Location: to.StringPtr("[resourceGroup().location]"),
		VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
			AddressSpace: &network.AddressSpace{
				AddressPrefixes: &[]string{"10.1.0.0/24"},
			},
			Subnets: &[]network.Subnet{
				network.Subnet{
					Name: to.StringPtr(defaultSNet),
					SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
						AddressPrefix: to.StringPtr("10.1.0.0/24"),
					},
				},
			},
		},
	}
}
