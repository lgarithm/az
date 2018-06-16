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
	// defaultVMSize = compute.StandardA1
	defaultVMSize = compute.VirtualMachineSizeTypesStandardDS1V2
)

type vmOptions struct {
	AddLocalSSHKey  bool
	AllowPassword   bool
	AdminUser       string
	AdminPassword   string
	CloudInitScript string
}

// DefaultVMOptions creates a VMOptions with recommended values
func DefaultVMOptions() vmOptions {
	const defaultAdminUser = "cup"
	const defaultAdminPassword = "Test1234"
	return vmOptions{
		AddLocalSSHKey:  true,
		AdminUser:       defaultAdminUser,
		AdminPassword:   defaultAdminPassword,
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

func newVM(name string, ni NetworkInterfaceResource, o *vmOptions) compute.VirtualMachine {
	opts := DefaultVMOptions()
	if o != nil {
		opts = *o
	}
	osDiskName := fmt.Sprintf("%s-disk-%d", name, time.Now().Unix())
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

func newWindowsVM(name string, ni NetworkInterfaceResource, o *vmOptions) compute.VirtualMachine {
	opts := DefaultVMOptions()
	if o != nil {
		opts = *o
	}
	osDiskName := fmt.Sprintf("%s-disk-%d", name, time.Now().Unix())
	return compute.VirtualMachine{
		Type:     to.StringPtr(TypeVM),
		Name:     to.StringPtr(name),
		Location: to.StringPtr("[resourceGroup().location]"),
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			OsProfile:       newWindowsOsProfile(name, opts),
			HardwareProfile: &compute.HardwareProfile{VMSize: defaultVMSize},
			NetworkProfile:  newNetworkProfile(ni),
			StorageProfile:  newStorageProfileFromImage(osDiskName, defaultWindowsImage),
		},
	}
}

func newLinuxOsProfile(name string, opts vmOptions) *compute.OSProfile {
	var defaultMyKey = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	var linuxConfiguration *compute.LinuxConfiguration
	if opts.AddLocalSSHKey {
		if bs, err := ioutil.ReadFile(defaultMyKey); err == nil {
			key := compute.SSHPublicKey{
				Path:    to.StringPtr(fmt.Sprintf("/home/%s/.ssh/authorized_keys", opts.AdminUser)),
				KeyData: to.StringPtr(string(bs)),
			}
			linuxConfiguration = &compute.LinuxConfiguration{
				SSH: &compute.SSHConfiguration{PublicKeys: &[]compute.SSHPublicKey{key}},
			}
		} else {
			glog.Warningf("failed to add local ssh key: %v", err)
		}
	}
	if opts.AllowPassword {
		linuxConfiguration.DisablePasswordAuthentication = to.BoolPtr(false)
	}
	var customData *string
	if opts.CloudInitScript != "" {
		customData = to.StringPtr(base64.StdEncoding.EncodeToString([]byte(opts.CloudInitScript)))
	}
	return &compute.OSProfile{
		CustomData:         customData,
		ComputerName:       to.StringPtr(name),
		AdminUsername:      to.StringPtr(opts.AdminUser),
		AdminPassword:      to.StringPtr(opts.AdminPassword),
		LinuxConfiguration: linuxConfiguration,
	}
}

func newWindowsOsProfile(name string, opts vmOptions) *compute.OSProfile {
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
		AdminUsername:        to.StringPtr(opts.AdminUser),
		AdminPassword:        to.StringPtr(opts.AdminPassword),
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
			ManagedDisk: &compute.ManagedDiskParameters{
				StorageAccountType: compute.StorageAccountTypesStandardLRS,
			},
			CreateOption: compute.DiskCreateOptionTypesFromImage,
		},
		DataDisks: &[]compute.DataDisk{},
	}
}
