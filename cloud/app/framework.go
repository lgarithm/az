package app

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/golang/glog"

	"teavana.com/cloud/azure/arm/dep"
)

type Action func() error

func Init() {
	pwd, _ := os.Getwd()
	dir := path.Join(pwd, "logs")
	os.Mkdir(dir, os.ModePerm)
	flag.Set("log_dir", dir)
	flag.Set("alsologtostderr", "true")
	flag.Parse()
	glog.Infof("args: %q", os.Args)
	glog.CopyStandardLogTo("INFO")
}

func Measure(name string, f Action) error {
	var prog = os.Args[0]
	begin := time.Now()
	glog.Infof("running %s::%s", prog, name)
	err := f()
	if err != nil {
		glog.Warningf("%s::%s took %s, with error: %v", prog, name, time.Now().Sub(begin), err)
	} else {
		glog.Infof("%s::%s took %s", prog, name, time.Now().Sub(begin))
	}
	return err
}

func Save(d dep.Deployment) error {
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}
	filename := path.Join("logs", d.Name+".json")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return d.SaveTemplate(f)
}
