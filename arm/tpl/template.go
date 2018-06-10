package tpl

import "github.com/Azure/azure-sdk-for-go/services/resources"

// DeploymentTemplate can be used by
// az group deployment create --template-file template.json -g $group
// when saved as json file
type DeploymentTemplate map[string]interface{}

func (d DeploymentTemplate) ToDeployment() resources.Deployment {
	value := (map[string]interface{})(d)
	return resources.Deployment{
		Properties: &resources.DeploymentProperties{
			Template:   &value,
			Parameters: &map[string]interface{}{},
			Mode:       resources.Incremental,
		},
	}
}
