package env

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/satori/uuid"
)

var (
	clientID       = flag.String("client_id", "", "AZURE_CLIENT_ID")
	subscriptionID = flag.String("subscription_id", "", "AZURE_SUBSCRIPTION_ID")
	tenantID       = flag.String("tenant_id", "", "AZURE_TENANT_ID")
)

type Env struct {
	ClientID       string
	SubscriptionID string
	TenantID       string
}

func New() *Env {
	tk := stoleToken()

	// get clientID
	if *clientID == "" {
		*clientID = os.Getenv("AZURE_CLIENT_ID")
	}
	if *clientID == "" {
		*clientID = tk.ClientID
	}
	if err := checkUUID(*clientID); err != nil {
		glog.Exitf("invalid AZURE_CLIENT_ID: %s", *clientID)
	}

	profile := loadAzureProfile()
	sub := profile.Subscriptions[0] // FIXME: check len

	// get subscriptionID
	if *subscriptionID == "" {
		*subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
	if *subscriptionID == "" {
		*subscriptionID = sub.ID
	}
	if err := checkUUID(*subscriptionID); err != nil {
		glog.Exitf("invalid AZURE_SUBSCRIPTION_ID: %s", *subscriptionID)
	}

	// get tenantID
	if *tenantID == "" {
		*tenantID = os.Getenv("AZURE_TENANT_ID")
	}
	if *tenantID == "" {
		*tenantID = sub.TenantID
	}
	if err := checkUUID(*tenantID); err != nil {
		glog.Exitf("invalid AZURE_TENANT_ID: %s", *tenantID)
	}

	return &Env{
		ClientID:       *clientID,
		SubscriptionID: *subscriptionID,
		TenantID:       *tenantID,
	}
}

func checkUUID(id string) error {
	_, err := uuid.FromString(id)
	return err
}
