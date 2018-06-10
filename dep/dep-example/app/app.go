package app

import (
	"flag"
	"fmt"

	"teavana.com/cloud/app"
	"teavana.com/cloud/azure/arm/dep"
)

type Flags struct {
	Location string
	Group    string
	Action   string
}

type App struct {
	Flags Flags
}

func (a App) Main() error {
	cloudInitScript, err := genCloudInitConfig().Encode()
	if err != nil {
		return err
	}
	builder := gen(cloudInitScript)
	d, err := dep.New("auto-deployment", a.Flags.Group, a.Flags.Location, builder)
	if err != nil {
		return err
	}
	app.Save(*d)
	var op app.Action
	switch a.Flags.Action {
	case "show":
		fmt.Printf("%s\n", cloudInitScript)
		return nil
	case "up":
		op = d.UpAndInit
	case "down":
		op = d.Down
	default:
		flag.Usage()
		return fmt.Errorf("invalued action %s", a.Flags.Action)
	}
	return app.Measure(a.Flags.Action, op)
}
