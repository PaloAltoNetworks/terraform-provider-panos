package poli

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/poli/pbf"
	"github.com/PaloAltoNetworks/pango/poli/security"
)

// Poli is the client.Policies namespace.
type FwPoli struct {
	Nat                   *nat.Firewall
	PolicyBasedForwarding *pbf.Firewall
	Security              *security.Firewall
}

// Initialize is invoked on client.Initialize().
func (c *FwPoli) Initialize(i util.XapiClient) {
	c.Nat = nat.FirewallNamespace(i)
	c.PolicyBasedForwarding = pbf.FirewallNamespace(i)
	c.Security = security.FirewallNamespace(i)
}
