package addrgrp

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// FwAddrGrp is a namespace struct, included as part of pango.Client.
type FwAddrGrp struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwAddrGrp) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of address groups.
func (c *FwAddrGrp) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of address groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of address groups.
func (c *FwAddrGrp) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of address groups")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given address group.
func (c *FwAddrGrp) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) address group %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given address group.
func (c *FwAddrGrp) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) address group %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more address groups.
func (c *FwAddrGrp) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "address-group"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) address groups: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an address group.
func (c *FwAddrGrp) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) address group %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the objects.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given address groups from the firewall.
//
// Address groups can be either a string or an Entry object.
func (c *FwAddrGrp) Delete(vsys string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
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
    c.con.LogAction("(delete) address groups: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the FwAddrGrp struct **/

func (c *FwAddrGrp) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwAddrGrp) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwAddrGrp) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    if vsys == "shared" {
        return []string {
            "config",
            "shared",
            "address-group",
            util.AsEntryXpath(vals),
        }
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "address-group",
        util.AsEntryXpath(vals),
    }
}
