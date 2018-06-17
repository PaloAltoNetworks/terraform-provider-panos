package dev

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/dev/general"
    "github.com/PaloAltoNetworks/pango/dev/telemetry"
)


// FwDev is the client.Device namespace.
type FwDev struct {
    GeneralSettings *general.General
    Telemetry *telemetry.FwTelemetry
}

// Initialize is invoked on client.Initialize().
func (c *FwDev) Initialize(i util.XapiClient) {
    c.GeneralSettings = &general.General{}
    c.GeneralSettings.Initialize(i)

    c.Telemetry = &telemetry.FwTelemetry{}
    c.Telemetry.Initialize(i)
}
