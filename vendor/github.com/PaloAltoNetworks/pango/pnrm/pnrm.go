/*
Package pnrm is the client.Panorama namespace.
*/
package pnrm


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/pnrm/dg"
)


// Pnrm is the panorama.DeviceGroup namespace.
type Pnrm struct {
    DeviceGroup *dg.Dg
}

// Initialize is invoked on panorama.Initialize().
func (c *Pnrm) Initialize(i util.XapiClient) {
    c.DeviceGroup = &dg.Dg{}
    c.DeviceGroup.Initialize(i)
}
