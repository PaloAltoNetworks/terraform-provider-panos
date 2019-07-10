package poli

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/poli/nat"
    "github.com/PaloAltoNetworks/pango/poli/pbf"
    "github.com/PaloAltoNetworks/pango/poli/security"
)


// Poli is the client.Policies namespace.
type FwPoli struct {
    Nat *nat.FwNat
    PolicyBasedForwarding *pbf.FwPbf
    Security *security.FwSecurity
}

// Initialize is invoked on client.Initialize().
func (c *FwPoli) Initialize(i util.XapiClient) {
    c.Nat = &nat.FwNat{}
    c.Nat.Initialize(i)

    c.PolicyBasedForwarding = &pbf.FwPbf{}
    c.PolicyBasedForwarding.Initialize(i)

    c.Security = &security.FwSecurity{}
    c.Security.Initialize(i)
}
