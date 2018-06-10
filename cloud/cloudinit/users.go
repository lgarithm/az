package cloudinit

// User is http://cloudinit.readthedocs.io/en/latest/topics/modules.html#users-and-groups
type User struct {
	Name           string     `yaml:"name"`
	AuthorizedKeys []Resource `yaml:"ssh-authorized-keys"`
}
