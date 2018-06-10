package cloudinit

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Config represents a cloud-init config file.
type Config struct {
	Timezone   string          `yaml:"timezone,omitempty"`
	Apt        Apt             `yaml:"apt,omitempty"`
	Packages   []Package       `yaml:"packages,omitempty"`
	SSHKeys    SSHKeys         `yaml:"ssh_keys,omitempty"`
	Users      []User          `yaml:"users,omitempty"`
	WriteFiles []WriteFile     `yaml:"write_files,omitempty"`
	Runcmd     []ScriptSection `yaml:"runcmd,omitempty"`
}

func (c Config) Encode() (string, error) {
	const shebang = "#cloud-config"
	bs, err := yaml.Marshal(&c)
	return fmt.Sprintf("%s\n%s\n", shebang, string(bs)), err
}
