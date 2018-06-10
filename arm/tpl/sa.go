package tpl

import (
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

func newSA(name string) storage.Account {
	return storage.Account{
		Name:     to.StringPtr(name),
		Type:     to.StringPtr(TypeSA),
		Location: to.StringPtr("[resourceGroup().location]"),
		Sku: &storage.Sku{
			Name: storage.StandardLRS,
			Tier: storage.Standard,
		},
		Kind: storage.Storage,
	}
}
