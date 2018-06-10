package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/lgarithm/az/arm"
)

func main() {
	flag.Parse()
	cf, err := arm.NewClientFactory()
	if err != nil {
		glog.Exit(err)
	}
	cf.Info()
}
