package cloudinit

// SSHKeys is http://cloudinit.readthedocs.io/en/latest/topics/modules.html#ssh
type SSHKeys struct {
	RSAPrivate Resource `yaml:"rsa_private,omitempty"`
	RSAPpublic Resource `yaml:"rsa_public,omitempty"`
}
