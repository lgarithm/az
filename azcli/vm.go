package azcli

type VM struct {
	Image Image
	Size  Size
}

func AZVM(args ...string) Proc {
	args = append([]string{`vm`}, args...)
	return AZ(args...)
}

func CreateVM(name string, relay, admin string, location, group string, size, image string) Proc {
	return AZVM(
		`create`,
		`--admin-username`, admin,
		`-l`, location,
		`-g`, group,
		`--nsg`, relay+`NSG`,
		`--vnet-name`, relay+`VNET`,
		`--subnet`, relay+`Subnet`,
		`--public-ip-address`, ``,
		`--size`, size,
		`--image`, image,
		`-n`, name,
		`--debug`,
		`-o`, `table`,
	)
}

func CreatePublicVM(name string, relay, admin string, location, group string, size, image string) Proc {
	return AZVM(
		`create`,
		`--admin-username`, admin,
		`-l`, location,
		`-g`, group,
		`--nsg`, relay+`NSG`,
		`--vnet-name`, relay+`VNET`,
		`--subnet`, relay+`Subnet`,
		`--size`, size,
		`--image`, image,
		`-n`, name,
		`--debug`,
		`-o`, `table`,
	)
}

func StartVM(name string, group string) Proc {
	return AZVM(
		`start`,
		`-g`, group,
		`-n`, name,
		`-o`, `table`,
	)
}

func StopVM(name string, group string) Proc {
	return AZVM(
		`deallocate`,
		`-g`, group,
		`-n`, name,
		`--debug`,
		`-o`, `table`,
	)
}

func DeleteVM(name string, group string) Proc {
	return AZVM(
		`delete`,
		`-g`, group,
		`-n`, name,
		`--yes`,
		`--debug`,
		`-o`, `table`,
	)
}

func GetIP(name string, group string) Proc {
	return AZVM(
		`list-ip-addresses`,
		`-g`, group,
		`-n`, name,
		`--query`, `[0].virtualMachine.network.privateIpAddresses[0]`,
	)
}

func GetPubIP(name string, group string) Proc {
	return AZVM(
		`list-ip-addresses`,
		`-g`, group,
		`-n`, name,
		`--query`, `[0].virtualMachine.network.publicIpAddresses[0].ipAddress`,
	)
}
