package azcli

import "strconv"

func AZ(args ...string) Proc {
	return Proc{
		Prog: `az`,
		Args: args,
	}
}

func Login() Proc { return AZ(`login`) }

func ListResource() Proc { return AZ(`resource`, `list`, `-o`, `table`) }

type Place struct {
	Admin    string
	Group    string
	Location string
}

type (
	Size  string
	Image string
)

const (
	Ubuntu20 Image = `Canonical:0001-com-ubuntu-server-focal:20_04-lts:latest`
	Ubuntu22 Image = `Canonical:0001-com-ubuntu-server-jammy:22_04-lts:latest`

	B1ls  Size = `Standard_B1ls`
	K80x1 Size = `Standard_NC6_Promo`
)

var str = strconv.Itoa
