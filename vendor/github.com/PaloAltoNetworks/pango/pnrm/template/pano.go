package template

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)

// Template is the client.Panorama.Template namespace.
type Template struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Template) Initialize(con util.XapiClient) {
    c.con = con
}

/*
SetDeviceVsys performs a SET to add specific vsys from a device to
template t.

If you want all vsys to be included, or the device is a virtual firewall, then
leave the vsys list empty.

The template can be either a string or an Entry object.
*/
func (c *Template) SetDeviceVsys(t interface{}, d string, vsys []string) error {
    var name string

    switch v := t.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to set device vsys: %s", v)
    }

    c.con.LogAction("(set) device vsys in template: %s", name)

    m := util.MapToVsysEnt(map[string] []string{d: vsys})
    path := c.xpath([]string{name})
    path = append(path, "devices")

    _, err := c.con.Set(path, m.Entries[0], nil, nil)
    return err
}

/*
EditDeviceVsys performs an EDIT to add specific vsys from a device to
template t.

If you want all vsys to be included, or the device is a virtual firewall, then
leave the vsys list empty.

The template can be either a string or an Entry object.
*/
func (c *Template) EditDeviceVsys(t interface{}, d string, vsys []string) error {
    var name string

    switch v := t.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to edit device vsys: %s", v)
    }

    c.con.LogAction("(edit) device vsys in template: %s", name)

    m := util.MapToVsysEnt(map[string] []string{d: vsys})
    path := c.xpath([]string{name})
    path = append(path, "devices", util.AsEntryXpath([]string{d}))

    _, err := c.con.Edit(path, m.Entries[0], nil, nil)
    return err
}

/*
DeleteDeviceVsys performs a DELETE to remove specific vsys from device d from
template t.

If you want all vsys to be removed, or the device is a virtual firewall, then
leave the vsys list empty.

The template can be either a string or an Entry object.
*/
func (c *Template) DeleteDeviceVsys(t interface{}, d string, vsys []string) error {
    var name string

    switch v := t.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to delete device vsys: %s", v)
    }

    c.con.LogAction("(delete) device vsys from template: %s", name)

    path := make([]string, 0, 9)
    path = append(path, c.xpath([]string{name})...)
    path = append(path, "devices", util.AsEntryXpath([]string{d}))
    if len(vsys) > 0 {
        path = append(path, "vsys", util.AsEntryXpath(vsys))
    }

    _, err := c.con.Delete(path, nil, nil)
    return err
}

// ShowList performs SHOW to retrieve a list of templates.
func (c *Template) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of templates")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of templates.
func (c *Template) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of templates")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given template.
func (c *Template) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) template %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given template.
func (c *Template) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) template %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more templates.
func (c *Template) Set(e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "template"}}
    for i := range e {
        e[i].SetConfTree()
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) templates: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the templates.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update a template.
func (c *Template) Edit(e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) template %q", e.Name)
    e.SetConfTree()

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the template.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given templates from the firewall.
//
// Templates can be a string or an Entry object.
func (c *Template) Delete(e ...interface{}) error {
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
            return fmt.Errorf("Unknown type sent to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) templates: %v", names)

    // Remove the templates.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *Template) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 1, 0, ""}) {
        return &container_v3{}, specify_v3
    } else if v.Gte(version.Number{7, 0, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *Template) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Template) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "template",
        util.AsEntryXpath(vals),
    }
}
