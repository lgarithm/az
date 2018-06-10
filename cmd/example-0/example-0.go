package main

import (
	"context"
	"flag"
	"fmt"

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
	ctx := context.Background()
	client := cf.NewGroupsClient()
	res, err := client.List(ctx, "", nil)
	if err != nil {
		glog.Exit(err)
	}
	for _, g := range res.Values() {
		fmt.Printf("%s\n", *g.Name)
	}
}
