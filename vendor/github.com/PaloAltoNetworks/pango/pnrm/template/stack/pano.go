package stack

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Stack is the client.Panorama.TemplateStack namespace.
type Stack struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Stack) Initialize(con util.XapiClient) {
    c.con = con
}

/*
SetDevice performs a SET to add specific device to template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Stack) SetDevice(st interface{}, d string) error {
    var name string

    switch v := st.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to set device: %s", v)
    }

    c.con.LogAction("(set) device in template stack: %s", name)

    path := c.xpath([]string{name})
    path = append(path, "devices")

    _, err := c.con.Set(path, util.Entry{Value: d}, nil, nil)
    return err
}

/*
EditDevice performs an EDIT to add specific device to template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Stack) EditDevice(st interface{}, d string) error {
    var name string

    switch v := st.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to edit device: %s", v)
    }

    c.con.LogAction("(edit) device in template stack: %s", name)

    path := c.xpath([]string{name})
    path = append(path, "devices", util.AsEntryXpath([]string{d}))

    _, err := c.con.Edit(path, util.Entry{Value: d}, nil, nil)
    return err
}

/*
DeleteDevice performs a DELETE to remove specific device d from template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Stack) DeleteDevice(st interface{}, d string) error {
    var name string

    switch v := st.(type) {
    case string:
        name = v
    case Entry:
        name = v.Name
    default:
        return fmt.Errorf("Unknown type sent to delete device: %s", v)
    }

    c.con.LogAction("(delete) device from template stack: %s", name)

    path := c.xpath([]string{name})
    path = append(path, "devices", util.AsEntryXpath([]string{d}))

    _, err := c.con.Delete(path, nil, nil)
    return err
}

// ShowList performs SHOW to retrieve a list of template stacks.
func (c *Stack) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of template stacks")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of template stacks.
func (c *Stack) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of template stacks")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given template stack.
func (c *Stack) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) template stack %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given template stack.
func (c *Stack) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) template stack %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more template stacks.
func (c *Stack) Set(e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "template-stack"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) template stacks: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the template stacks.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update a template stack.
func (c *Stack) Edit(e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) template stack %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the template stack.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given template stacks from the firewall.
//
// Objects can be a string or an Entry object.
func (c *Stack) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) template stacks: %v", names)

    // Remove the template stacks.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for this namespace struct **/

func (c *Stack) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Stack) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Stack) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "template-stack",
        util.AsEntryXpath(vals),
    }
}
