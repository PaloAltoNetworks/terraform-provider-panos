// Package dev is the client.Device namespace.
package dev

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/dev/general"
)


// Dev is the client.Device namespace.
type Dev struct {
    GeneralSettings *general.General
}

// Initialize is invoked on client.Initialize().
func (c *Dev) Initialize(i util.XapiClient) {
    c.GeneralSettings = &general.General{}
    c.GeneralSettings.Initialize(i)
}
