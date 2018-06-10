package arm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-02-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/lgarithm/az/env"
)

var (
	cloud                 = azure.PublicCloud
	baseURI               = cloud.ResourceManagerEndpoint
	storageServiceBaseURL = cloud.StorageEndpointSuffix
)

func newSPT(e *env.Env) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(cloud.ActiveDirectoryEndpoint, e.TenantID)
	if err != nil {
		return nil, err
	}
	tok := StoleToken()
	t := adal.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
	}
	return adal.NewServicePrincipalTokenFromManualToken(*oauthConfig, e.ClientID, cloud.ResourceManagerEndpoint, t)
}

// ClientFactory helps creating azure API clients
type ClientFactory struct {
	env        *env.Env
	spt        *adal.ServicePrincipalToken
	authorizer autorest.Authorizer
}

// NewClientFactory creates a new ClentFactory
func NewClientFactory() (*ClientFactory, error) {
	e := env.New()
	spt, err := newSPT(e)
	if err != nil {
		return nil, err
	}
	cf := ClientFactory{
		env:        e,
		spt:        spt,
		authorizer: autorest.NewBearerAuthorizer(spt),
	}
	return &cf, nil
}

// NewGroupsClient creates an azure resources.GroupsClient
func (cf *ClientFactory) NewGroupsClient() *resources.GroupsClient {
	client := resources.NewGroupsClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// // NewGroupClient creates an azure resources.GroupClient
// func (cf *ClientFactory) NewGroupClient() *resources.GroupClient {
// 	client := resources.NewGroupClientWithBaseURI(baseURI, cf.env.SubscriptionID)
// 	client.Authorizer = cf.authorizer
// 	return &client
// }

// NewDepClient creates an azure resources.DeploymentsClient
func (cf *ClientFactory) NewDepClient() *resources.DeploymentsClient {
	client := resources.NewDeploymentsClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewVMClient creates an azure compute.VirtualMachinesClient
func (cf *ClientFactory) NewVMClient() *compute.VirtualMachinesClient {
	client := compute.NewVirtualMachinesClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewVNClient creates an azure network.VirtualNetworksClient
func (cf *ClientFactory) NewVNClient() *network.VirtualNetworksClient {
	client := network.NewVirtualNetworksClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewNIClient creates an azure network.InterfacesClient
func (cf *ClientFactory) NewNIClient() *network.InterfacesClient {
	client := network.NewInterfacesClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewIPClient creates an azure network.PublicIPAddressesClient
func (cf *ClientFactory) NewIPClient() *network.PublicIPAddressesClient {
	client := network.NewPublicIPAddressesClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewSAClient creates an azure storage.AccountsClient
func (cf *ClientFactory) NewSAClient() *storage.AccountsClient {
	client := storage.NewAccountsClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

// NewCRClient creates an azure containerregistry.RegistriesClient
func (cf *ClientFactory) NewCRClient() *containerregistry.RegistriesClient {
	client := containerregistry.NewRegistriesClientWithBaseURI(baseURI, cf.env.SubscriptionID)
	client.Authorizer = cf.authorizer
	return &client
}

func (cf *ClientFactory) Info() {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "    ")
	e.Encode(cf.env)
}

// Info displays metadata about autorest.Client
func Info(client autorest.Client) {
	fmt.Printf("autorest.Client:\n")
	fmt.Printf("\tPollingDelay    %s\n", client.PollingDelay)
	fmt.Printf("\tPollingDuration %s\n", client.PollingDuration)
	fmt.Printf("\tRetryAttempts   %d\n", client.RetryAttempts)
	fmt.Printf("\tRetryDuration   %s\n", client.RetryDuration)
}
