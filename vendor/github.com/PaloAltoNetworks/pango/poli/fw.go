package poli

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/poli/security"
    "github.com/PaloAltoNetworks/pango/poli/nat"
)


// Poli is the client.Policies namespace.
type FwPoli struct {
    Security *security.FwSecurity
    Nat *nat.FwNat
}

// Initialize is invoked on client.Initialize().
func (c *FwPoli) Initialize(i util.XapiClient) {
    c.Security = &security.FwSecurity{}
    c.Security.Initialize(i)

    c.Nat = &nat.FwNat{}
    c.Nat.Initialize(i)
}
