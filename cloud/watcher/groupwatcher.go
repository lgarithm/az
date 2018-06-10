package watcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/golang/glog"
	"github.com/lgarithm/az/arm"
	"github.com/lgarithm/go/control"
	"github.com/lgarithm/go/profile"
	"golang.org/x/sync/errgroup"
)

// GroupWatcher watches a resource group
type GroupWatcher struct {
	Group    string
	vmClient *compute.VirtualMachinesClient
	niClient *network.InterfacesClient
	ipClient *network.PublicIPAddressesClient
	gsClient *resources.Client
}

func setDefaultPolling(client *autorest.Client) {
	client.PollingDelay = 30 * time.Second   // was 1min
	client.PollingDuration = 4 * time.Minute // was 15min
}

// NewGroupWatcher creates a new Watcher
func NewGroupWatcher(group string) GroupWatcher {
	cf, err := arm.NewClientFactory()
	if err != nil {
		glog.Exit(err)
	}

	vmClient := cf.NewVMClient()
	ipClient := cf.NewIPClient()
	niClient := cf.NewNIClient()
	gsClient := cf.NewResourceClient()
	setDefaultPolling(&vmClient.Client)
	setDefaultPolling(&ipClient.Client)
	setDefaultPolling(&niClient.Client)
	return GroupWatcher{
		Group:    group,
		vmClient: vmClient,
		niClient: niClient,
		ipClient: ipClient,
		gsClient: gsClient,
	}
}

// WaitCloudInit waits until cloud-init finish in all VMs in a resource group
func (w GroupWatcher) WaitCloudInit(ctx context.Context, vmNames []string) error {
	defer profile.Profile("WaitCloudInit").Done()
	g, ctx := errgroup.WithContext(ctx)
	for _, name := range vmNames {
		func(name string) {
			g.Go(func() error { return w.NewVMWatcher(name).WaitCloudInit(ctx) })
		}(name)
	}
	return g.Wait()
}

// WaitSSH waits until all VMs can be ssh into
func (w GroupWatcher) WaitSSH(ctx context.Context, vmNames []string) error {
	defer profile.Profile("WaitSSH").Done()
	w.WaitAllVMs(ctx, vmNames, func(name string) error {
		return w.NewVMWatcher(name).watchSSH()
	})
	return nil
}

// WaitDown waits until all resources are deleted in a resource group
func (w GroupWatcher) WaitDown(ctx context.Context) error {
	defer profile.Profile("WaitDown").Done()
	return control.Wait(ctx, 15*time.Second, func() error {
		res, err := w.gsClient.ListByResourceGroup(context.TODO(), w.Group, "", "", nil)
		if err != nil {
			// if httpRes := res.Response.Response; httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			// 	return nil
			// }
			return err
		}
		var waits []string
		for _, r := range res.Values() {
			waits = append(waits, *r.Name)
		}
		if len(waits) > 0 {
			log.Printf("still deleting %s", strings.Join(waits, ", "))
			return fmt.Errorf("Still deleting %d resources", len(waits))
		}
		return nil
	})
}

// NewVMWatcher creates a VMWatcher with given VM name
func (w GroupWatcher) NewVMWatcher(name string) *VMWatcher {
	return &VMWatcher{
		Group:      w.Group,
		Name:       name,
		SSHTimeout: 6 * time.Second,
		User:       `cup`,

		vmClient: w.vmClient,
		niClient: w.niClient,
		ipClient: w.ipClient,
	}
}

func (w GroupWatcher) WaitAllVMs(ctx context.Context, vmNames []string, poll func(string) error) {
	var wg sync.WaitGroup
	wg.Add(len(vmNames))
	for _, name := range vmNames {
		go func(name string) {
			control.Wait(ctx, 10*time.Second, func() error { return poll(name) })
			wg.Done()
		}(name)
	}
	wg.Wait()
}
