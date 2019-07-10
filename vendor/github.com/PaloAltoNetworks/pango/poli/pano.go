package poli

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/poli/nat"
    "github.com/PaloAltoNetworks/pango/poli/pbf"
    "github.com/PaloAltoNetworks/pango/poli/security"
)


// Poli is the client.Policies namespace.
type PanoPoli struct {
    Nat *nat.PanoNat
    PolicyBasedForwarding *pbf.PanoPbf
    Security *security.PanoSecurity
}

// Initialize is invoked on client.Initialize().
func (c *PanoPoli) Initialize(i util.XapiClient) {
    c.Nat = &nat.PanoNat{}
    c.Nat.Initialize(i)

    c.PolicyBasedForwarding = &pbf.PanoPbf{}
    c.PolicyBasedForwarding.Initialize(i)

    c.Security = &security.PanoSecurity{}
    c.Security.Initialize(i)
}
