package main

import (
	"flag"

	"github.com/golang/glog"

	"teavana.com/cloud/app"

	deploy_cluster_app "teavana.com/cloud/azure/arm/dep/dep-example/app"
)

var (
	location = flag.String("location", "southeastasia", "azure location")
	group    = flag.String("group", "teavana-example", "azure resource group")

	action = flag.String("action", "up", "show | up | down")
)

func main() {
	app.Init()
	defer glog.Flush()
	flags := deploy_cluster_app.Flags{
		Location: *location,
		Group:    *group,
		Action:   *action,
	}
	app := deploy_cluster_app.App{Flags: flags}
	if err := app.Main(); err != nil {
		glog.Exit(err)
	}
}
