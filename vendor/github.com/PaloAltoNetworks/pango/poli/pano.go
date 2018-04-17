package poli

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/poli/security"
    "github.com/PaloAltoNetworks/pango/poli/nat"
)


// Poli is the client.Policies namespace.
type PanoPoli struct {
    Security *security.PanoSecurity
    Nat *nat.PanoNat
}

// Initialize is invoked on client.Initialize().
func (c *PanoPoli) Initialize(i util.XapiClient) {
    c.Security = &security.PanoSecurity{}
    c.Security.Initialize(i)

    c.Nat = &nat.PanoNat{}
    c.Nat.Initialize(i)
}
