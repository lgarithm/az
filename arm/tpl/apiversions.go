package tpl

type APIVersionConfig struct {
	Default              string
	NetworkInterface     string
	NetworkSecurityGroup string
	PublicIPAddress      string
	VirtualNetwork       string
}

var APIVersions = APIVersionConfig{
	Default:              "2016-04-30-preview",
	NetworkInterface:     "2017-03-01",
	NetworkSecurityGroup: "2017-03-01",
	PublicIPAddress:      "2017-03-01",
	VirtualNetwork:       "2017-03-01",
}
