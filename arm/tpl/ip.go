package tpl

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/Azure/go-autorest/autorest/to"
)

func newIP(name string) network.PublicIPAddress {
	return network.PublicIPAddress{
		Name:     to.StringPtr(name),
		Type:     to.StringPtr(TypeIP),
		Location: to.StringPtr("[resourceGroup().location]"),
	}
}
