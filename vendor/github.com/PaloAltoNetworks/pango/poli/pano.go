package poli

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/poli/pbf"
	"github.com/PaloAltoNetworks/pango/poli/security"
)

// Poli is the client.Policies namespace.
type PanoPoli struct {
	Nat                   *nat.Panorama
	PolicyBasedForwarding *pbf.Panorama
	Security              *security.Panorama
}

// Initialize is invoked on client.Initialize().
func (c *PanoPoli) Initialize(i util.XapiClient) {
	c.Nat = nat.PanoramaNamespace(i)
	c.PolicyBasedForwarding = pbf.PanoramaNamespace(i)
	c.Security = security.PanoramaNamespace(i)
}
