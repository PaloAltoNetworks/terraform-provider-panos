package layer2

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// FwLayer2 is the client.Network.Layer2Subinterface namespace.
type FwLayer2 struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwLayer2) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwLayer2) ShowList(iType, eth, mType string) ([]string, error) {
    c.con.LogQuery("(show) list of %s", plural)
    path := c.xpath(iType, eth, mType, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwLayer2) GetList(iType, eth, mType string) ([]string, error) {
    c.con.LogQuery("(get) list of %s", plural)
    path := c.xpath(iType, eth, mType, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwLayer2) Get(iType, eth, mType, name string) (Entry, error) {
    c.con.LogQuery("(get) %s %q", singular, name)
    return c.details(c.con.Get, iType, eth, mType, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwLayer2) Show(iType, eth, mType, name string) (Entry, error) {
    c.con.LogQuery("(show) %s %q", singular, name)
    return c.details(c.con.Show, iType, eth, mType, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwLayer2) Set(iType, eth, mType, vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if iType == "" {
        return fmt.Errorf("iType must be specified")
    } else if eth == "" {
        return fmt.Errorf("eth must be specified")
    } else if mType == "" {
        return fmt.Errorf("mType must be specified")
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct.
    d := util.BulkElement{XMLName: xml.Name{Local: "units"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) %s: %v", plural, names)

    // Set xpath.
    path := c.xpath(iType, eth, mType, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    if err != nil {
        return err
    }

    // Remove from any vsys it's currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", names); err != nil {
        return err
    }

    // Perform vsys import.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, names)
}

// Edit performs EDIT to create / update one object.
func (c *FwLayer2) Edit(iType, eth, mType, vsys string, e Entry) error {
    var err error

    if iType == "" {
        return fmt.Errorf("iType must be specified")
    } else if eth == "" {
        return fmt.Errorf("eth must be specified")
    } else if mType == "" {
        return fmt.Errorf("mType must be specified")
    }

    _, fn := c.versioning()

    c.con.LogAction("(edit) %s %q", singular, e.Name)

    // Set xpath.
    path := c.xpath(iType, eth, mType, []string{e.Name})

    // Edit the object.
    if _, err = c.con.Edit(path, fn(e), nil, nil); err != nil {
        return err
    }

    // Remove from any vsys it's currently in.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", []string{e.Name}); err != nil {
        return err
    }

    // Perform vsys import.
    return c.con.VsysImport(util.InterfaceImport, "", "", vsys, []string{e.Name})
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwLayer2) Delete(iType, eth, mType string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    } else if iType == "" {
        return fmt.Errorf("iType must be specified")
    } else if eth == "" {
        return fmt.Errorf("eth must be specified")
    } else if mType == "" {
        return fmt.Errorf("mType must be specified")
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

    // Unimport interfaces.
    if err = c.con.VsysUnimport(util.InterfaceImport, "", "", names); err != nil {
        return err
    }

    // Remove the objects.
    path := c.xpath(iType, eth, mType, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *FwLayer2) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *FwLayer2) details(fn util.Retriever, iType, eth, mType, name string) (Entry, error) {
    path := c.xpath(iType, eth, mType, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *FwLayer2) xpath(iType, eth, mType string, vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "interface",
        iType,
        util.AsEntryXpath([]string{eth}),
        mType,
        "units",
        util.AsEntryXpath(vals),
    }
}
