package main

import (
	"encoding/json"
	"flag"
	"io"

	"github.com/golang/glog"
	"github.com/lgarithm/az/arm"
	"github.com/lgarithm/az/cmd/example-1/app"
	"github.com/lgarithm/az/dep"
)

var (
	group    = flag.String("group", "test-01", "azure resource group")
	location = flag.String("location", "southeastasia", "azure location")
)

func main() {
	flag.Parse()
	cf, err := arm.NewClientFactory()
	if err != nil {
		glog.Exit(err)
	}
	client := cf.NewDepClient()
	builder := app.New("")
	d, err := dep.New("auto-deployment", *group, *location, builder)
	if err != nil {
		glog.Exit(err)
	}
	d.Up()

	// ctx := context.Background()
	// res, err := client.CreateOrUpdate(ctx, "test-01", "auto-deployment", dep.ToDeployment())
	// if err != nil {
	// 	glog.Exit(err)
	// }
	// saveJSON(res, os.Stdout)
}

func saveJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(i)
}
