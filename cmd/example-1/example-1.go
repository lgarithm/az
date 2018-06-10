package main

import (
	"encoding/json"
	"flag"
	"io"

	"github.com/golang/glog"
	"github.com/lgarithm/az/cmd/example-1/app"
	"github.com/lgarithm/az/dep"
)

var (
	group    = flag.String("group", "test-01", "azure resource group")
	location = flag.String("location", "southeastasia", "azure location")
)

func main() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	glog.CopyStandardLogTo("INFO")
	defer glog.Flush()
	builder := app.New("")
	d, err := dep.New("auto-deployment", *group, *location, builder)
	if err != nil {
		glog.Exit(err)
	}
	d.Up()
}

func saveJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(i)
}
