package cloudinit

import yaml "gopkg.in/yaml.v2"

// Apt represents apt config for cloud-init
type Apt struct {
	Sources map[string]AptSource `yaml:"sources"`
}

// AptSource represents apt source config for cloud-init
type AptSource struct {
	Source string   `yaml:"source"`
	Key    Resource `yaml:"key"`
}

type Package interface {
	yaml.Marshaler
}

type NamedPackage struct {
	Package

	Name string
}

func (p NamedPackage) MarshalYAML() (interface{}, error) {
	return p.Name, nil
}

type VersionedPackage struct {
	Package

	Name    string
	Version string
}

func (p VersionedPackage) MarshalYAML() (interface{}, error) {
	return []string{p.Name, p.Version}, nil
}

func NewNamedPackage(name string) Package {
	return NamedPackage{Name: name}
}

func NewVersionedPackage(name, version string) Package {
	return VersionedPackage{Name: name, Version: version}
}
