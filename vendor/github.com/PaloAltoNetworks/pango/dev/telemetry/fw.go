package telemetry

import (
	"github.com/PaloAltoNetworks/pango/util"
)

// FwTelemetry is a namespace struct, included as part of pango.Firewall.
type FwTelemetry struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwTelemetry) Initialize(con util.XapiClient) {
	c.con = con
}

// Show performs SHOW to retrieve telemetry sharing settings.
func (c *FwTelemetry) Show() (Settings, error) {
	c.con.LogQuery("(show) telemetry settings")
	return c.details(c.con.Show)
}

// Get performs GET to retrieve telemetry sharing settings.
func (c *FwTelemetry) Get() (Settings, error) {
	c.con.LogQuery("(get) telemetry settings")
	return c.details(c.con.Get)
}

// Set performs SET to update telemetry sharing settings.
func (c *FwTelemetry) Set(e Settings) error {
	var err error
	_, fn := c.versioning()
	c.con.LogAction("(set) telemetry settings")

	path := c.xpath()
	path = path[:len(path)-1]

	_, err = c.con.Set(path, fn(e), nil, nil)
	return err
}

// Edit performs EDIT to update telemetry sharing settings.
func (c *FwTelemetry) Edit(e Settings) error {
	var err error
	_, fn := c.versioning()
	c.con.LogAction("(edit) telemetry settings")

	path := c.xpath()

	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes all telemetry sharing from the firewall.
func (c *FwTelemetry) Delete() error {
	c.con.LogAction("(delete) telemetry settings")
	path := c.xpath()

	_, err := c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for the FwTelemetry struct **/

func (c *FwTelemetry) versioning() (normalizer, func(Settings) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *FwTelemetry) details(fn util.Retriever) (Settings, error) {
	path := c.xpath()
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Settings{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *FwTelemetry) xpath() []string {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"deviceconfig",
		"system",
		"update-schedule",
		"statistics-service",
	}
}
