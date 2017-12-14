// Package poli is the client.Policies namespace.
package poli

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/poli/security"
    "github.com/PaloAltoNetworks/pango/poli/nat"
)


// Poli is the client.Policies namespace.
type Poli struct {
    Security *security.Security
    Nat *nat.Nat
}

// Initialize is invoked on client.Initialize().
func (c *Poli) Initialize(i util.XapiClient) {
    c.Security = &security.Security{}
    c.Security.Initialize(i)

    c.Nat = &nat.Nat{}
    c.Nat.Initialize(i)
}
