package dep

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/golang/glog"
	"github.com/lgarithm/az/armutil"
	"github.com/lgarithm/az/cloud/watcher"
	"github.com/lgarithm/go/control"
	"github.com/lgarithm/go/profile"
	"github.com/lgarithm/go/xterm"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Reload reloads the deployment
func (d Deployment) Reload() error {
	if err := d.Down(); err != nil {
		return err
	}
	return d.UpAndInit()
}

var defaultTimeout = 14 * time.Minute

// Up setup the deployment and waits all VMs can be ssh into
func (d Deployment) Up() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()
	var g errgroup.Group
	g.Go(func() error {
		return watcher.NewGroupWatcher(d.Group).WaitSSH(ctx, d.VMNames())
	})
	t, err := profile.Duration(func() error { _, err := d.Deploy(); return err })
	if err != nil {
		glog.Warningf("deploy %s, took %s, error: %s", xterm.Red.S("failed"), t, err)
		return err
	}
	glog.Infof("deploy finished, took %s, wait until VMs can be ssh into", t)
	return g.Wait()
}

// UpAndInit setup the deployment and waits until cloud-init finish
func (d Deployment) UpAndInit() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()
	var g errgroup.Group
	g.Go(func() error {
		return watcher.NewGroupWatcher(d.Group).WaitCloudInit(ctx, d.VMNames())
	})
	t, err := profile.Duration(func() error { _, err := d.Deploy(); return err })
	if err != nil {
		glog.Warningf("deploy %s, took %s, error: %s", xterm.Red.S("failed"), t, err)
		return err
	}
	glog.Infof("deploy finished, took %s, waiting cloud-init", t)
	return g.Wait()
}

// Down deletes the resource group
func (d Deployment) Down() error {
	return armutil.TearDownGroup(d.clientFactory, d.Group)
}

// Validate validates the deployment config
func (d Deployment) Validate() error {
	defer profile.Profile(fmt.Sprintf("Validating %s", d.Name)).Done()
	client := d.clientFactory.NewDepClient()
	client.PollingDelay = 10 * time.Second
	client.PollingDuration = 1 * time.Minute
	res, err := client.Validate(context.TODO(), d.Group, d.Name, d.template.ToDeployment())
	if err != nil {
		return err
	}
	if res.Error != nil {
		msg, _ := toJSON(res.Error)
		return errors.New(msg)
	}
	return nil
}

// Deploy pushes the deployment config
func (d Deployment) Deploy() (*resources.DeploymentExtended, error) {
	if err := armutil.EnsureGroup(d.clientFactory, d.Group, d.Location); err != nil {
		return nil, err
	}
	if err := d.Validate(); err != nil {
		return nil, errors.Wrap(err, "Validate failed")
	}
	glog.Infof("Deploying %s to %s in %s", d.Name, d.Group, d.Location)
	client := d.clientFactory.NewDepClient()
	client.PollingDelay = 30 * time.Second
	promise, err := client.CreateOrUpdate(context.TODO(), d.Group, d.Name, d.template.ToDeployment())
	if err := promise.WaitForCompletionRef(context.TODO(), client.Client); err != nil {
		return nil, errors.Wrap(err, "CreateOrUpdate failed")
	}
	res, err := promise.Result(*client)
	if err != nil {
		return nil, errors.Wrap(err, "CreateOrUpdate failed")
	}
	return &res, nil
}

// Info contains metadata of a deployment
type Info struct {
	Group string
	IPMap map[string]string
}

// GetInfo returns Info of a Deployment
func (d Deployment) GetInfo() (*Info, error) {
	ipMap, _, err := d.GetIPMap()
	if err != nil {
		return nil, err
	}
	return &Info{
		Group: d.Group,
		IPMap: ipMap,
	}, nil
}

// With runs a job with the deployment.
func (d Deployment) With(job func(info Info) error) error {
	d.Down() // ignore
	defer d.Down()
	if err := d.UpAndInit(); err != nil {
		return control.LogError(errors.Wrap(err, "UpAndInit failed"))
	}
	const n = 3
	var info *Info
	if err := control.Try(n, func() error {
		var err error
		info, err = d.GetInfo()
		return err
	}); err != nil {
		return control.LogError(errors.Wrap(err, fmt.Sprintf("GetInfo Failed for %d times", n)))
	}
	return control.LogError(errors.Wrap(job(*info), "job failed"))
}

func toJSON(v interface{}) (string, error) {
	j, err := json.MarshalIndent(v, "", "  ")
	return string(j), err
}
