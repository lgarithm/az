package app

import (
	"os"
	"path"

	"teavana.com/cloud/cloudinit"
)

const timezone = "Hongkong"

var pkgs = []string{
	"libffi-dev",
	"libssl-dev",
	"git",
	"golang-go",
	"influxdb",
	"nginx",
	"python-pip",
	"ruby",
	"sqlite3",
}

func genCloudInitConfig() cloudinit.Config {
	// dockerKey := cloudinit.NewURLResource("https://download.docker.com/linux/ubuntu/gpg")
	// dockerKey := cloudinit.NewFileResource(path.Join(os.Getenv("HOME"), ".teavana.d", "docker.gpg"))
	apt := cloudinit.Apt{
		Sources: map[string]cloudinit.AptSource{
		// "docker.list": cloudinit.AptSource{
		// 	Source: "deb https://download.docker.com/linux/ubuntu xenial stable",
		// 	Key:    dockerKey,
		// },
		},
	}
	var packages []cloudinit.Package
	for _, p := range pkgs {
		packages = append(packages, cloudinit.NewNamedPackage(p))
	}
	// packages = append(packages, cloudinit.NewVersionedPackage("docker-ce", "17.03.2~ce-0~ubuntu-xenial"))
	keyFilePrefix := path.Join(os.Getenv("HOME"), ".teavana.d/id")
	zipfile := cloudinit.ZipFile{
		Files: []cloudinit.File{
			{
				Name:   path.Join(keyFilePrefix, "id_rsa"),
				Target: path.Join(".ssh", "id_rsa"),
			},
			{
				Name:   path.Join(keyFilePrefix, "id_rsa.pub"),
				Target: path.Join(".ssh", "id_rsa.pub"),
			},
		},
	}
	config := cloudinit.Config{
		Timezone: timezone,
		Apt:      apt,
		Packages: packages,
		// Users:    users,
		// WriteFiles: writeFiles,
		Runcmd: []cloudinit.ScriptSection{
			cloudinit.NewUnzipFileScript(zipfile, "cup"),
			cloudinit.NewFileScript("extension/example1.sh"),
			cloudinit.NewFileScript("extension/phone_slack.sh"),
		},
	}
	return config
}
