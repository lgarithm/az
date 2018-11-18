package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
	"github.com/lgarithm/az/cmd/ml-cluster/app"
	"github.com/lgarithm/az/dep"
	"github.com/lgarithm/go/control"
	"github.com/lgarithm/go/profile"
)

var (
	group    = flag.String("group", "test-ml-cluster", "azure resource group")
	location = flag.String("location", "southeastasia", "azure location")
	action   = flag.String("action", "up", "up | down | reload")
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
	withFile("template.json", func(f io.Writer) error { return d.SaveTemplate(f) })
	if err := run(d); err != nil {
		glog.Exit(err)
	}
}

func run(d *dep.Deployment) error {
	defer profile.Profile("main::run").Done()
	switch *action {
	case "up":
		return d.Up()
	case "down":
		return d.Down()
	case "reload":
		return control.Chain(d.Down, d.Up)()
	default:
		return fmt.Errorf("invalid action: %s", *action)
	}
}

func saveJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(i)
}

func withFile(filename string, f func(io.Writer) error) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	return f(w)
}
