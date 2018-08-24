package mngtprof

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// FwMngtProf is a namespace struct, included as part of pango.Client.
type FwMngtProf struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwMngtProf) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of interface management profiles.
func (c *FwMngtProf) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of interface management profiles")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of interface management profiles.
func (c *FwMngtProf) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of interface management profiles")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given interface management
// profile.
func (c *FwMngtProf) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) interface management profile %q", name)
    return c.details(c.con.Get, name)
}

// Get performs SHOW to retrieve information for the given interface management
// profile.
func (c *FwMngtProf) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) interface management profile %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more interface management profiles.
func (c *FwMngtProf) Set(e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "interface-management-profile"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) interface management profiles: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the profiles.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an interface management profile.
func (c *FwMngtProf) Edit(e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) interface management profile %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the profile.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given interface management profile(s) from the firewall.
//
// Profiles can be either a string or an Entry object.
func (c *FwMngtProf) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) interface management profiles: %v", names)

    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwMngtProf) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwMngtProf) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwMngtProf) xpath(vals []string) []string {
    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "profiles",
        "interface-management-profile",
        util.AsEntryXpath(vals),
    }
}
