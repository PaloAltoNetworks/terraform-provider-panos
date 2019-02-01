package bgp

import (
    "fmt"
    //"encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// FwBgp is the client.Network.RedistributionProfile namespace.
type FwBgp struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwBgp) Initialize(con util.XapiClient) {
    c.con = con
}

// Get performs GET to retrieve the BGP config.
func (c *FwBgp) Get(vr string) (Config, error) {
    c.con.LogQuery("(get) bgp config for %q", vr)
    return c.details(c.con.Get, vr)
}

// Show performs SHOW to retrieve the BGP config.
func (c *FwBgp) Show(vr string) (Config, error) {
    c.con.LogQuery("(show) bgp config for %q", vr)
    return c.details(c.con.Show, vr)
}

// Set performs SET to create / update the BGP config.
func (c *FwBgp) Set(vr string, e Config) error {
    var err error

    if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()
    c.con.LogAction("(set) bgp config for %q", vr)
    path := c.xpath(vr)
    path = path[:len(path) - 1]

    _, err = c.con.Set(path, fn(e), nil, nil)
    return err
}

// Edit performs EDIT to create / update the BGP config.
func (c *FwBgp) Edit(vr string, e Config) error {
    var err error

    if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()
    c.con.LogAction("(edit) bgp config for %q", vr)
    path := c.xpath(vr)

    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the BGP config for the given virtual router.
func (c *FwBgp) Delete(vr string) error {
    var err error

    if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    c.con.LogAction("(delete) bgp config for %q", vr)

    // Remove the objects.
    path := c.xpath(vr)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwBgp) versioning() (normalizer, func(Config) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &container_v4{}, specify_v4
    } else if v.Gte(version.Number{7, 1, 0, ""}) {
        return &container_v3{}, specify_v3
    } else if v.Gte(version.Number{7, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *FwBgp) details(fn util.Retriever, vr string) (Config, error) {
    path := c.xpath(vr)
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Config{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwBgp) xpath(vr string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "virtual-router",
        util.AsEntryXpath([]string{vr}),
        "protocol",
        "bgp",
    }
}
