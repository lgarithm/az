package azcli

type StorageSKU string

const Premium_LRS StorageSKU = `Premium_LRS`

type VMSS struct {
	Name    string
	VM      VM
	Storage StorageSKU
}

func (v VMSS) CreateAt(p Place, n int, vn *VNet) Proc {
	args := []string{`vmss`, `create`,
		`-l`, p.Location,
		`-g`, p.Group,
		`--disable-overprovision`,
		`--instance-count`, str(n),
		`--vm-sku`, string(v.VM.Size),
		`--image`, string(v.VM.Image),
		`--admin-username`, p.Admin,
		// `--nsg`, relay+`NSG`,
		`--lb`, ``,
		`-n`, v.Name,
		`--debug`,
		`-o`, `table`,
	}
	if len(v.Storage) > 0 {
		args = append(args, `--storage-sku`, string(v.Storage))
	}
	if vn != nil {
		args = append(args,
			`--vnet-name`, vn.Name, // relay+`VNET`,
			`--subnet`, vn.Subnet, // relay+`Subnet`,
		)
	}
	return AZ(args...)
}

func createVMSS(name string, admin string, location, group string, size, image string, n int, vn *VNet) Proc {
	v := VMSS{
		Name: name,
		VM: VM{
			Image: Image(image),
			Size:  Size(size),
		},
	}
	p := Place{
		Admin:    admin,
		Location: location,
		Group:    group,
	}
	return v.CreateAt(p, n, vn)
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
