package azcli

type VMSS struct {
	Name string
	VM   VM
}

func createVMSS(name string, admin string, location, group string, size, image string, n int, vn *VNet) Proc {
	args := []string{`vmss`, `create`,
		`-l`, location,
		`-g`, group,
		`--disable-overprovision`,
		`--instance-count`, str(n),
		`--vm-sku`, size,
		`--image`, image,
		`--admin-username`, admin,
		// `--nsg`, relay+`NSG`,
		`--lb`, ``,
		`-n`, name,
		`--debug`,
		`-o`, `table`,
	}
	if vn != nil {
		args = append(args,
			`--vnet-name`, vn.Name, // relay+`VNET`,
			`--subnet`, vn.Subnet, // relay+`Subnet`,
		)
	}
	return AZ(args...)
}

func DeleteVMSS(name string, group string) Proc {
	return AZ(`vmss`, `delete`,
		`-g`, group,
		`-n`, name,
		`--debug`,
		`-o`, `table`,
	)
}

func ScaleVMSS(name string, group string, n int) Proc {
	return AZ(`vmss`, `scale`,
		`-g`, group,
		`-n`, name,
		`--new-capacity`, str(n),
		`--debug`,
		`-o`, `table`,
	)
}

func (p Place) CreateVMSS(name string, image Image, size Size, n int, vn *VNet) Proc {
	return createVMSS(name, p.Admin, p.Location, p.Group, string(size), string(image), n, vn)
}
