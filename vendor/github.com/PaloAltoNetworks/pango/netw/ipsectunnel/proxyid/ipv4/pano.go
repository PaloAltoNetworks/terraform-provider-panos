package ipv4

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// PanoIpv4 is a namespace struct, included as part of pango.Firewall.
type PanoIpv4 struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *PanoIpv4) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of IPSec tunnel proxy IDs.
func (c *PanoIpv4) GetList(tmpl, ts, tun string) ([]string, error) {
    c.con.LogQuery("(get) list of ipsec tunnel proxy ids")
    path := c.xpath(tmpl, ts, tun, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of IPSec tunnel proxy IDs.
func (c *PanoIpv4) ShowList(tmpl, ts, tun string) ([]string, error) {
    c.con.LogQuery("(show) list of ipsec tunnel proxy ids")
    path := c.xpath(tmpl, ts, tun, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given IPSec tunnel proxy ID.
func (c *PanoIpv4) Get(tmpl, ts, tun, name string) (Entry, error) {
    c.con.LogQuery("(get) ipsec tunnel proxy id %q", name)
    return c.details(c.con.Get, tmpl, ts, tun, name)
}

// Get performs SHOW to retrieve information for the given IPSec tunnel proxy ID.
func (c *PanoIpv4) Show(tmpl, ts, tun, name string) (Entry, error) {
    c.con.LogQuery("(show) ipsec tunnel proxy id %q", name)
    return c.details(c.con.Show, tmpl, ts, tun, name)
}

// Set performs SET to create / update one or more IPSec tunnel proxy IDs.
func (c *PanoIpv4) Set(tmpl, ts, tun string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if tun == "" {
        return fmt.Errorf("tun must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "proxy-id"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) ipsec tunnel proxy ids: %v", names)

    // Set xpath.
    path := c.xpath(tmpl, ts, tun, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an IPSec tunnel proxy ID.
func (c *PanoIpv4) Edit(tmpl, ts, tun string, e Entry) error {
    var err error

    if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if tun == "" {
        return fmt.Errorf("tun must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) ipsec tunnel proxy id %q", e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, tun, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given IPSec tunnel proxy IDs from the firewall.
//
// Items can be either a string or an Entry object.
func (c *PanoIpv4) Delete(tmpl, ts, tun string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if tmpl == "" && ts == "" {
        return fmt.Errorf("tmpl or ts must be specified")
    } else if tun == "" {
        return fmt.Errorf("tun must be specified")
    }

    names := make([]string, len(e))
    for i := range e {
        switch v := e[i].(type) {
        case string:
            names[i] = v
        case Entry:
            names[i] = v.Name
        default:
            return fmt.Errorf("Unsupported type to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) ipsec tunnel proxy ids: %v", names)

    path := c.xpath(tmpl, ts, tun, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoIpv4) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *PanoIpv4) details(fn util.Retriever, tmpl, ts, tun, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, tun, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoIpv4) xpath(tmpl, ts, tun string, vals []string) []string {
    ans := make([]string, 0, 15)
    ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
    ans = append(ans,
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "tunnel",
        "ipsec",
        util.AsEntryXpath([]string{tun}),
        "auto-key",
        "proxy-id",
        util.AsEntryXpath(vals),
    )

    return ans
}
