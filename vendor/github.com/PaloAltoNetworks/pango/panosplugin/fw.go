package panosplugin

import (
	"github.com/PaloAltoNetworks/pango/panosplugin/cloudwatch"
	"github.com/PaloAltoNetworks/pango/util"
)

type Firewall struct {
	AwsCloudWatch *cloudwatch.Firewall
}

func FirewallNamespace(x util.XapiClient) *Firewall {
	return &Firewall{
		AwsCloudWatch: cloudwatch.FirewallNamespace(x),
	}
}
