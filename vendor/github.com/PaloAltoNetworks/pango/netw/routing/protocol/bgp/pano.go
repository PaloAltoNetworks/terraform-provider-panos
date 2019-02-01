package bgp

import (
    "fmt"
    //"encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// PanoBgp is the client.Network.RedistributionProfile namespace.
type PanoBgp struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoBgp) Initialize(con util.XapiClient) {
    c.con = con
}

// Get performs GET to retrieve the BGP config.
func (c *PanoBgp) Get(tmpl, ts, vr string) (Config, error) {
    c.con.LogQuery("(get) bgp config for %q", vr)
    return c.details(c.con.Get, tmpl, ts, vr)
}

// Show performs SHOW to retrieve the BGP config.
func (c *PanoBgp) Show(tmpl, ts, vr string) (Config, error) {
    c.con.LogQuery("(show) bgp config for %q", vr)
    return c.details(c.con.Show, tmpl, ts, vr)
}

// Set performs SET to create / update the BGP config.
func (c *PanoBgp) Set(tmpl, ts, vr string, e Config) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()
    c.con.LogAction("(set) bgp config for %q", vr)
    path := c.xpath(tmpl, ts, vr)
    path = path[:len(path) - 1]

    _, err = c.con.Set(path, fn(e), nil, nil)
    return err
}

// Edit performs EDIT to create / update the BGP config.
func (c *PanoBgp) Edit(tmpl, ts, vr string, e Config) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    _, fn := c.versioning()
    c.con.LogAction("(edit) bgp config for %q", vr)
    path := c.xpath(tmpl, ts, vr)

    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the BGP config for the given virtual router.
func (c *PanoBgp) Delete(tmpl, ts, vr string) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if vr == "" {
        return fmt.Errorf("vr must be specified")
    }

    c.con.LogAction("(delete) bgp config for %q", vr)

    // Remove the objects.
    path := c.xpath(tmpl, ts, vr)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoBgp) versioning() (normalizer, func(Config) (interface{})) {
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

func (c *PanoBgp) details(fn util.Retriever, tmpl, ts, vr string) (Config, error) {
    path := c.xpath(tmpl, ts, vr)
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Config{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoBgp) xpath(tmpl, ts, vr string) []string {
    ans := make([]string, 0, 13)
    ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
    ans = append(ans,
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "virtual-router",
        util.AsEntryXpath([]string{vr}),
        "protocol",
        "bgp",
    )

    return ans
}
