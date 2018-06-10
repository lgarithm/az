package armutil

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/golang/glog"

	"github.com/lgarithm/az/arm"
	"github.com/lgarithm/az/cloud/watcher"
)

// EnsureGroup creates a resource group if not exists
func EnsureGroup(cf *arm.ClientFactory, name, location string) error {
	client := cf.NewGroupsClient()
	_, err := client.CreateOrUpdate(name, resources.Group{Location: &location})
	return err
}

// TearDownGroup deletes a resource group
func TearDownGroup(cf *arm.ClientFactory, name string) error {
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
	resch, errch := client.Delete(name, nil)
	err := <-errch
	res := <-resch
	if err != nil {
		// FIXME
		if res.Response == nil {
			glog.Error("res.Response is null")
			return err
		}
		if res.StatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	return nil
}
