package general

import (
    "github.com/PaloAltoNetworks/pango/util"
)


// FwGeneral is a namespace struct, included as part of pango.Client.
type FwGeneral struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwGeneral) Initialize(con util.XapiClient) {
    c.con = con
}

// Show performs SHOW to retrieve the device's general settings.
func (c *FwGeneral) Show() (Config, error) {
    c.con.LogQuery("(show) general settings")
    return c.details(c.con.Show)
}

// Get performs GET to retrieve the device's general settings.
func (c *FwGeneral) Get() (Config, error) {
    c.con.LogQuery("(get) general settings")
    return c.details(c.con.Get)
}

// Set performs SET to create / update the device's general settings.
func (c *FwGeneral) Set(e Config) error {
    var err error
    _, fn := c.versioning()
    c.con.LogAction("(set) general settings")

    path := c.xpath()
    path = path[:len(path) - 1]

    _, err = c.con.Set(path, fn(e), nil, nil)
    return err
}

// Edit performs EDIT to update the device's general settings.
func (c *FwGeneral) Edit(e Config) error {
    var err error
    _, fn := c.versioning()
    c.con.LogAction("(edit) general settings")

    path := c.xpath()

    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

/** Internal functions for the FwGeneral struct **/

func (c *FwGeneral) versioning() (normalizer, func(Config) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwGeneral) details(fn util.Retriever) (Config, error) {
    path := c.xpath()
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Config{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwGeneral) xpath() []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "deviceconfig",
        "system",
    }
}
