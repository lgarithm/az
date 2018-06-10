package tpl

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/glog"
)

const (
	defaultAdminUser     = "cup"
	defaultAdminPassword = "Test1234"

	// defaultVMSize = compute.StandardA1
	defaultVMSize = compute.VirtualMachineSizeTypesStandardDS1V2
)

type VMOptions struct {
	AddLocalSSHKey  bool
	CloudInitScript string
}

// DefaultVMOptions creates a VMOptions with recommended values
func DefaultVMOptions() VMOptions {
	return VMOptions{
		AddLocalSSHKey:  true,
		CloudInitScript: "",
	}
}

var (
	defaultUbuntuImage = &compute.ImageReference{
		Publisher: to.StringPtr("Canonical"),
		Offer:     to.StringPtr("UbuntuServer"),
		Sku:       to.StringPtr("18.04-LTS"),
		Version:   to.StringPtr("latest"),
	}
	defaultWindowsImage = &compute.ImageReference{
		Publisher: to.StringPtr("MicrosoftWindowsServer"),
		Offer:     to.StringPtr("WindowsServer"),
		// Sku:       to.StringPtr("2016-Nano-Server"),
		Sku:     to.StringPtr("2016-Datacenter"),
		Version: to.StringPtr("latest"),
	}
)

func newVM(name string, ni NetworkInterfaceResource, o *VMOptions) compute.VirtualMachine {
	opts := DefaultVMOptions()
	if o != nil {
		opts = *o
	}
	osDiskName := fmt.Sprintf("%s-%d", name, time.Now().Unix())
	return compute.VirtualMachine{
		Type:     to.StringPtr(TypeVM),
		Name:     to.StringPtr(name),
		Location: to.StringPtr("[resourceGroup().location]"),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile:       newLinuxOsProfile(name, opts),
			HardwareProfile: &compute.HardwareProfile{VMSize: defaultVMSize},
			NetworkProfile:  newNetworkProfile(ni),
			StorageProfile:  newStorageProfileFromImage(osDiskName, defaultUbuntuImage),
		},
	}
}

func newWindowsVM(name string, ni NetworkInterfaceResource) compute.VirtualMachine {
	osDiskName := fmt.Sprintf("%s-%d", name, time.Now().Unix())
	return compute.VirtualMachine{
		Type:     to.StringPtr(TypeVM),
		Name:     to.StringPtr(name),
		Location: to.StringPtr("[resourceGroup().location]"),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile:       newWindowsOsProfile(name),
			HardwareProfile: &compute.HardwareProfile{VMSize: defaultVMSize},
			NetworkProfile:  newNetworkProfile(ni),
			StorageProfile:  newStorageProfileFromImage(osDiskName, defaultWindowsImage),
		},
	}
}

var defaultMyKey = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")

func newLinuxOsProfile(name string, opts VMOptions) *compute.OSProfile {
	var linuxConfiguration *compute.LinuxConfiguration
	if opts.AddLocalSSHKey {
		if bs, err := ioutil.ReadFile(defaultMyKey); err == nil {
			key := compute.SSHPublicKey{
				Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", defaultAdminUser)),
				KeyData: to.StringPtr(string(bs)),
			}
			linuxConfiguration = &compute.LinuxConfiguration{
				// DisablePasswordAuthentication: to.BoolPtr(false),
				SSH: &compute.SSHConfiguration{PublicKeys: &[]compute.SSHPublicKey{key}},
			}
		} else {
			glog.Warningf("failed to add local ssh key: %v", err)
		}
	}
	var customData *string
	if opts.CloudInitScript != "" {
		customData = to.StringPtr(base64.StdEncoding.EncodeToString([]byte(opts.CloudInitScript)))
	}
	return &compute.OSProfile{
		CustomData:         customData,
		ComputerName:       to.StringPtr(name),
		AdminUsername:      to.StringPtr(defaultAdminUser),
		AdminPassword:      to.StringPtr(defaultAdminPassword),
		LinuxConfiguration: linuxConfiguration,
	}
}

func newWindowsOsProfile(name string) *compute.OSProfile {
	windowsConfiguration := &compute.WindowsConfiguration{
		WinRM: &compute.WinRMConfiguration{
			Listeners: &[]compute.WinRMListener{
				compute.WinRMListener{
					Protocol: compute.HTTP,
				},
				// compute.WinRMListener{
				// 	Protocol: compute.HTTPS,
				// 	CertificateURL: to.StringPtr(""),
				// },
			},
		},
	}
	return &compute.OSProfile{
		ComputerName:         to.StringPtr(name),
		AdminUsername:        to.StringPtr(defaultAdminUser),
		AdminPassword:        to.StringPtr(defaultAdminPassword),
		WindowsConfiguration: windowsConfiguration,
	}
}

func newNetworkProfile(ni NetworkInterfaceResource) *compute.NetworkProfile {
	return &compute.NetworkProfile{
		NetworkInterfaces: &[]compute.NetworkInterfaceReference{
			compute.NetworkInterfaceReference{
				ID: to.StringPtr(fmt.Sprintf("[resourceId('%s', '%s')]", *ni.Type, *ni.Name)),
			},
		},
	}
}

func newStorageProfileFromImage(name string, image *compute.ImageReference) *compute.StorageProfile {
	return &compute.StorageProfile{
		ImageReference: image,
		OsDisk: &compute.OSDisk{
			Name: to.StringPtr(name),
			// ManagedDisk: &compute.ManagedDiskParameters{
			// 	StorageAccountType: compute.StandardLRS,
			// },
			CreateOption: compute.DiskCreateOptionTypesFromImage,
		},
		DataDisks: &[]compute.DataDisk{},
	}
}
