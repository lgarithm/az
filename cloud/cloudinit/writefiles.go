package cloudinit

import (
	"fmt"
	"path"
)

type WriteFile struct {
	Encoding    string   `yaml:"encoding,omitempty"`
	Content     Resource `yaml:"content"`
	Owner       string   `yaml:"owner"`
	Path        string   `yaml:"path"`
	Permissions string   `yaml:"permissions,omitempty"`
}

func NewKeyPairsForUser(user, keyFilePrefix string) (WriteFile, WriteFile) {
	group := user
	owner := fmt.Sprintf("%s:%s", user, group)
	home := fmt.Sprintf("/home/%s", user)
	private := WriteFile{
		Owner:       owner,
		Content:     NewFileResource(path.Join(keyFilePrefix, "id_rsa")),
		Path:        path.Join(home, ".ssh", "id_rsa"),
		Permissions: "600",
	}
	public := WriteFile{
		Owner:       owner,
		Content:     NewFileResource(path.Join(keyFilePrefix, "id_rsa.pub")),
		Path:        path.Join(home, ".ssh", "id_rsa.pub"),
		Permissions: "644",
	}
	return private, public
}
