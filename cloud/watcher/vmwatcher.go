package watcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/golang/glog"
	"github.com/lgarithm/az/cloud/azure"
	"github.com/lgarithm/az/cloud/cloudinit"
	"github.com/lgarithm/go/control"
	"github.com/lgarithm/go/net/ssh"
	"github.com/lgarithm/go/xterm"
	"github.com/pkg/errors"
)

// VMWatcher watches a VM
type VMWatcher struct {
	Group    string
	Name     string
	vmClient *compute.VirtualMachinesClient
	niClient *network.InterfacesClient
	ipClient *network.PublicIPAddressesClient
}

func (w VMWatcher) GetVMPublicIPs() ([]string, error) {
	vm, err := w.vmClient.Get(context.TODO(), w.Group, w.Name, "")
	if err != nil {
		return nil, err
	}
	var ipAddresses []string
	for _, niRef := range *vm.NetworkProfile.NetworkInterfaces {
		ni, err := w.niClient.Get(context.TODO(), w.Group, azure.ID2Name(*niRef.ID), "")
		if err != nil {
			glog.Warning(err)
			continue
		}
		for _, ipCfg := range *ni.IPConfigurations {
			if ipCfg.PublicIPAddress == nil {
				continue
			}
			pubIP := *ipCfg.PublicIPAddress
			ip, err := w.ipClient.Get(context.TODO(), w.Group, azure.ID2Name(*pubIP.ID), "")
			if err != nil {
				glog.Warning(err)
				continue
			}
			if ip.IPAddress == nil {
				continue
			}
			ipAddress := *ip.IPAddress
			ipAddresses = append(ipAddresses, ipAddress)
		}
	}
	return ipAddresses, nil
}

func (w VMWatcher) GetIP() (string, error) {
	ips, err := w.GetVMPublicIPs()
	if err != nil {
		return "", err
	}
	if len(ips) <= 0 {
		return "", fmt.Errorf("VM %s has no IP", w.Name)
	}
	return ips[0], nil
}

func (w VMWatcher) watchSSH() error {
	ips, err := w.GetVMPublicIPs()
	if err != nil {
		return err
	}
	if len(ips) <= 0 {
		err := fmt.Errorf("%s has no ip", w.Name)
		log.Print(err)
		return err
	}
	err = sshTest(ips[0])
	if err == nil {
		log.Printf("%s can be ssh into", w.Name)
	} else {
		log.Print(err)
	}
	return err
}

// WaitCloudInit waits until cloud-init finish in a VM.
func (w VMWatcher) WaitCloudInit(ctx context.Context) error {
	ip, err := w.waitIP(ctx)
	if err != nil {
		return err
	}
	const period = 15 * time.Second
	var finalErr error
	begin := time.Now()
	control.Wait(ctx, period, func() error {
		err, finalErr = w.watchCloudInit(ctx, ip)
		if err == nil {
			if finalErr == nil {
				glog.Infof("cloud init finished in %s at %s", xterm.Green.S(w.Name), ip)
			} else {
				glog.Warningf("cloud init finished in %s at %s with error: %v", w.Name, ip, finalErr)
			}
		} else {
			duration := time.Now().Sub(begin)
			log.Printf("still wating cloud init for %s at %s after %s: %v, retry in %s", w.Name, ip, duration, err, period)
		}
		return err
	})
	if err != nil {
		return err
	}
	return finalErr
}

func (w VMWatcher) waitIP(ctx context.Context) (ip string, err error) {
	return ip, control.Wait(ctx, 6*time.Second, func() error {
		ip, err = w.GetIP()
		if err == nil {
			log.Printf("IP of VM %s is %s", w.Name, xterm.Green.S(ip))
		} else {
			log.Printf("still wating IP of VM %s", w.Name)
		}
		return err
	})
}

func (w VMWatcher) watchCloudInit(ctx context.Context, ip string) (error, error) {
	const timeout = 8 * time.Second
	ctx, cancal := context.WithTimeout(ctx, timeout)
	defer cancal()
	const cmd = `tail /var/log/cloud-init.log`
	stdout, _, err := ssh.RunWith(ctx, defaultUser, ip, cmd)
	if err != nil {
		return err, nil
	}
	lines := strings.Split(strings.TrimRight(string(stdout), "\n"), "\n")
	if len(lines) <= 0 {
		return errors.New("/var/log/cloud-init.log is empty"), nil
	}
	last := lines[len(lines)-1]
	if !cloudinit.IsLastLine(last) {
		return errors.New("cloud-init is still running"), nil
	}
	if err := cloudinit.FinalError(last); err != nil {
		return nil, errors.Wrap(err, "cloud-init failed")
	}
	return nil, nil
}

const defaultUser = "cup"

func sshTest(ip string) error {
	const timeout = 3 * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()
	log.Printf("test ssh %s@%s", defaultUser, ip)
	_, _, err := ssh.RunWith(ctx, defaultUser, ip, `pwd`)
	return err
}
