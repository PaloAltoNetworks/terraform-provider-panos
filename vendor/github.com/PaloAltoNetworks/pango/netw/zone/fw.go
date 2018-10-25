package zone

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// FwZone is a namespace struct, included as part of pango.Client.
type FwZone struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwZone) Initialize(con util.XapiClient) {
    c.con = con
}

/*
SetInterface performs a SET to add an interface to a zone.

The zone can be either a string or an Entry object.
*/
func (c *FwZone) SetInterface(vsys string, zone interface{}, mode, iface string) error {
    var name string

    switch v := zone.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s set interface: %s", singular, v)
    }

    c.con.LogAction("(set) %s interface: %s", singular, name)

    path := c.xpath(vsys, []string{name})
    path = append(path, "network", mode)

    _, err := c.con.Set(path, util.Member{Value: iface}, nil, nil)
    return err
}

/*
DeleteInterface performs a DELETE to remove the interface from the zone.

The zone can be either a string or an Entry object.
*/
func (c *FwZone) DeleteInterface(vsys string, zone interface{}, mode, iface string) error {
    var name string

    switch v := zone.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to %s delete interface: %s", singular, v)
    }

    c.con.LogAction("(delete) %s interface: %s", singular, name)

    path := c.xpath(vsys, []string{name})
    path = append(path, "network", mode, util.AsMemberXpath([]string{iface}))

    _, err := c.con.Delete(path, nil, nil)
    return err
}

// GetList performs GET to retrieve a list of values.
func (c *FwZone) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwZone) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwZone) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given uid.
func (c *FwZone) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwZone) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "zone"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

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

// Edit performs EDIT to create / update one object.
func (c *FwZone) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given objects.
//
// Objects can be either a string or an Entry object.
func (c *FwZone) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) %s: %v", plural, names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwZone) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwZone) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwZone) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "zone",
        util.AsEntryXpath(vals),
    }
}
