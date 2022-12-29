package azcli

func ListDisk(prefix string, group string) Proc {
	return AZ(
		`disk`, `list`,
		`-g`, group,
		`--query`,
		`[].name | [?starts_with(@, '`+prefix+`')]`,
		// `--debug`,
	)
}

func DeleteDisk(name string, group string) Proc {
	return AZ(
		`disk`, `delete`,
		`-g`, group,
		`-n`, name,
		`--yes`,
		// `--debug`,
	)
}
