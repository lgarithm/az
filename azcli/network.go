package azcli

func DeleteNIC(name string, group string) Proc {
	return AZ(`network`, `nic`,
		`delete`,
		`-g`, group,
		`-n`, name,
		// `--debug`,
	)
}

func DeletePublicIP(name string, group string) Proc {
	return AZ(
		`network`, `public-ip`,
		`delete`,
		`-g`, group,
		`-n`, name,
		// `--debug`,
	)
}

type VNet struct {
	Name   string
	Subnet string
}
