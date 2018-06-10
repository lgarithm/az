package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/lgarithm/az/cmd/example-1/app"
)

func main() {
	flag.Parse()
	// cf, err := arm.NewClientFactory()
	// if err != nil {
	// 	glog.Exit(err)
	// }
	// cf.Info()
	dep := app.New("").Build()
	saveJSON(dep, os.Stdout)
}

func saveJSON(i interface{}, w io.Writer) {
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	e.Encode(i)
}
