package tpl

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-05-01/network"
	"github.com/Azure/go-autorest/autorest/to"
)

var defaultNSGRules = []network.SecurityRule{
	NewAllowInboundRule("default-allow-ssh", "22", 1000),
	NewAllowInboundRule("default-allow-rdp", "3389", 999),
	NewAllowInboundRule("default-allow-winRM-http", "5985", 1100),
	NewAllowInboundRule("default-allow-winRM-https", "5986", 1101),

	NewAllowOutboundRule("default-allow-http", "80", 500),
	NewAllowOutboundRule("default-allow-https", "443", 501),
}

func newNSG(name string, rules *[]network.SecurityRule) network.SecurityGroup {
	if rules == nil {
		rules = &defaultNSGRules
	}
	return network.SecurityGroup{
		Name:     to.StringPtr(name),
		Type:     to.StringPtr(TypeSG),
		Location: to.StringPtr("[resourceGroup().location]"),
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			SecurityRules: rules,
		},
	}
}

func (nsg *NetworkSecurityGroupResource) AddRule(rule network.SecurityRule) {
	rules := nsg.SecurityGroupPropertiesFormat.SecurityRules
	newRules := append(*rules, rule)
	rules = &newRules
}

func NewAllowInboundRule(name string, portRange string, priority int32) network.SecurityRule {
	return newDefaultAllowRule(name, network.SecurityRuleDirectionInbound, portRange, priority)
}

func NewAllowOutboundRule(name string, portRange string, priority int32) network.SecurityRule {
	return newDefaultAllowRule(name, network.SecurityRuleDirectionInbound, portRange, priority)
}

func newDefaultAllowRule(name string, dir network.SecurityRuleDirection, portRange string, priority int32) network.SecurityRule {
	return network.SecurityRule{
		Name: to.StringPtr(name),
		SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
			Access: network.SecurityRuleAccessAllow,
			DestinationAddressPrefix: to.StringPtr("*"),
			DestinationPortRange:     to.StringPtr(portRange),
			Direction:                dir,
			Priority:                 to.Int32Ptr(priority),
			Protocol:                 network.SecurityRuleProtocolTCP,
			SourceAddressPrefix:      to.StringPtr("*"),
			SourcePortRange:          to.StringPtr("*"),
		},
	}
}
