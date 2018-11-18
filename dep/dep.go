package dep

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lgarithm/az/arm"
	"github.com/lgarithm/az/arm/tpl"
)

// Deployment represents a user defined deployment
type Deployment struct {
	Name           string
	Group          string
	Location       string
	clientFactory  *arm.ClientFactory
	template       tpl.DeploymentTemplate
	vmNameToIPName map[string]string
}

// New creates a Deployment
func New(name, group, location string, builder *tpl.Builder) (*Deployment, error) {
	cf, err := arm.NewClientFactory()
	if err != nil {
		return nil, err
	}
	d := Deployment{
		Name:           name,
		Group:          group,
		Location:       location,
		clientFactory:  cf,
		template:       builder.Build(),
		vmNameToIPName: builder.GetVM2IP(),
	}
	return &d, nil
}

// SaveTemplate saves Deployment template
func (d Deployment) SaveTemplate(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(d.template)
}

func (d Deployment) summary() ([]tpl.GenericResource, error) {
	bs, err := json.Marshal((d.template)["resources"])
	if err != nil {
		return nil, err
	}
	var gres []tpl.GenericResource
	if err := json.Unmarshal(bs, &gres); err != nil {
		return nil, err
	}
	return gres, nil
}

// Show prints a summary of the defined resources
func (d Deployment) Show() error {
	gres, err := d.summary()
	if err != nil {
		return err
	}
	bs, err := json.MarshalIndent(gres, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return nil
}

// RelayVMNames returns the list of VM names that has public IP addresses in the Deployment
func (d Deployment) RelayVMNames() []string {
	var names []string
	for name := range d.vmNameToIPName {
		names = append(names, name)
	}
	return names
}

// GetIPMap returns (VM Name -> IP, VM Names without IP, error)
func (d Deployment) GetIPMap() (map[string]string, []string, error) {
	res, err := d.clientFactory.NewIPClient().List(context.TODO(), d.Group)
	if err != nil {
		return nil, nil, err
	}
	ipNameToIPAddress := map[string]string{}
	for _, ip := range res.Values() {
		ipNameToIPAddress[*ip.Name] = *ip.IPAddress
	}
	vmNameToIPAddress := map[string]string{}
	var missing []string
	for vmName, ipName := range d.vmNameToIPName {
		if addr := ipNameToIPAddress[ipName]; addr == "" {
			missing = append(missing, vmName)
		} else {
			vmNameToIPAddress[vmName] = addr
		}
	}
	return vmNameToIPAddress, missing, nil
}
