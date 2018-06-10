package tpl

// disk "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
// "github.com/Azure/go-autorest/autorest/to"

// func newMD(name string) disk.Disk {
// 	// id := "[concat('/Subscriptions/', subscription(), '/Providers/Microsoft.Compute/Locations/', uniqueString(resourceGroup().location), '/Publishers/Canonical/ArtifactTypes/VMImage/Offers/UbuntuServer/Skus/17.04/Versions/latest')]"
// 	id := "" // FIXME
// 	return disk.Disk{
// 		Type:     to.StringPtr(TypeMD),
// 		Name:     to.StringPtr(name),
// 		Location: to.StringPtr("[resourceGroup().location]"),
// 		DiskProperties: &disk.DiskProperties{
// 			AccountType: disk.StandardLRS,
// 			CreationData: &disk.CreationData{
// 				CreateOption: disk.FromImage,
// 				ImageReference: &disk.ImageDiskReference{
// 					ID: to.StringPtr(id),
// 				},
// 			},
// 		},
// 	}
// }
