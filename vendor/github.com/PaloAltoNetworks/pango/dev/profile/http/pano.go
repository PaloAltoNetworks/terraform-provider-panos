package http

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// PanoHttp is the client.Network.HttpServerProfile namespace.
type PanoHttp struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoHttp) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *PanoHttp) ShowList(tmpl, ts, vsys, dg string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(tmpl, ts, vsys, dg, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *PanoHttp) GetList(tmpl, ts, vsys, dg string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(tmpl, ts, vsys, dg, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *PanoHttp) Get(tmpl, ts, vsys, dg, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, tmpl, ts, vsys, dg, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *PanoHttp) Show(tmpl, ts, vsys, dg, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, tmpl, ts, vsys, dg, name)
}

// Set performs SET to create / update one or more objects.
func (c *PanoHttp) Set(tmpl, ts, vsys, dg string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "temp"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(tmpl, ts, vsys, dg, names)
    d.XMLName = xml.Name{Local: path[len(path) - 2]}
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
func (c *PanoHttp) Edit(tmpl, ts, vsys, dg string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(tmpl, ts, vsys, dg, []string{e.Name})

    // Edit the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// SetWithoutSubconfig performs a DELETE to remove any subconfig
// before performing a SET to create an object.
func (c *PanoHttp) SetWithoutSubconfig(tmpl, ts, vsys, dg string, e Entry) error {
    c.con.LogAction("(delete) %s subconfig for %s", singular, e.Name)

    path := c.xpath(tmpl, ts, vsys, dg, []string{e.Name})

    path = append(path, "server")
    _, _ = c.con.Delete(path, nil, nil)

    path[len(path) - 1] = "format"
    _, _ = c.con.Delete(path, nil, nil)

    return c.Set(tmpl, ts, vsys, dg, e)
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *PanoHttp) Delete(tmpl, ts, vsys, dg string, e ...interface{}) error {
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
    c.con.LogAction("(delete) %s: %v", plural, names)

    // Remove the objects.
    path := c.xpath(tmpl, ts, vsys, dg, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *PanoHttp) versioning() (normalizer, func(Entry) (interface{})) {
    v := c.con.Versioning()

    if v.Gte(version.Number{9, 0, 0, ""}) {
        return &container_v3{}, specify_v3
    } else if v.Gte(version.Number{8, 1, 0, ""}) {
        return &container_v2{}, specify_v2
    } else {
        return &container_v1{}, specify_v1
    }
}

func (c *PanoHttp) details(fn util.Retriever, tmpl, ts, vsys, dg, name string) (Entry, error) {
    path := c.xpath(tmpl, ts, vsys, dg, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *PanoHttp) xpath(tmpl, ts, vsys, dg string, vals []string) []string {
    var ans []string

    if tmpl != "" || ts != "" {
        if vsys == "" {
            vsys = "shared"
        }

        ans = make([]string, 0, 13)
        ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
        ans = append(ans, util.VsysXpathPrefix(vsys)...)
    } else {
        if dg == "" {
            dg = "shared"
        }

        ans = make([]string, 0, 8)
        ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
    }

    ans = append(ans,
        "log-settings",
        "http",
        util.AsEntryXpath(vals),
    )

    return ans
}
