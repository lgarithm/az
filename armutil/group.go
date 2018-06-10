package armutil

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/lgarithm/az/arm"
	"github.com/lgarithm/az/cloud/watcher"
	"github.com/pkg/errors"
)

// EnsureGroup creates a resource group if not exists
func EnsureGroup(cf *arm.ClientFactory, name, location string) error {
	client := cf.NewGroupsClient()
	_, err := client.CreateOrUpdate(context.TODO(), name, resources.Group{Location: &location})
	return err
}

// TearDownGroup deletes a resource group
func TearDownGroup(cf *arm.ClientFactory, name string) error {
	client := cf.NewGroupsClient()
	res, err := client.CheckExistence(context.TODO(), name)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return nil
	}

	const defaultTimeout = 14 * time.Minute
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()
	done := make(chan error)
	go func() { done <- watcher.NewGroupWatcher(name).WaitDown(ctx) }()
	if err := teardownGroup(cf, name); err != nil {
		return err
	}
	return <-done
}

func teardownGroup(cf *arm.ClientFactory, name string) error {
	client := cf.NewGroupsClient()
	promise, err := client.Delete(context.TODO(), name)
	if err != nil {
		return errors.Wrap(err, "client.Delete")
	}
	return promise.WaitForCompletionRef(context.TODO(), client.Client)
}
