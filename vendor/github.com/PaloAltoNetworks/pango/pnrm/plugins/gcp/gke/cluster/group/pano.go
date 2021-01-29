package group

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

// Group is the client.Panorama.GkeClusterGroup namespace.
type Group struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Group) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *Group) ShowList() ([]string, error) {
	c.con.LogQuery("(show) list of %s", plural)
	path := c.xpath(nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of values.
func (c *Group) GetList() ([]string, error) {
	c.con.LogQuery("(get) list of %s", plural)
	path := c.xpath(nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given uid.
func (c *Group) Get(name string) (Entry, error) {
	c.con.LogQuery("(get) %s %q", singular, name)
	return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *Group) Show(name string) (Entry, error) {
	c.con.LogQuery("(show) %s %q", singular, name)
	return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more objects.
func (c *Group) Set(e ...Entry) error {
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
	path := c.xpath(names)
	d.XMLName = xml.Name{Local: path[len(path)-2]}
	if len(e) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the objects.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update one object.
func (c *Group) Edit(e Entry) error {
	var err error

	_, fn := c.versioning()

	c.con.LogAction("(edit) %s %q", singular, e.Name)

	// Set xpath.
	path := c.xpath([]string{e.Name})

	// Edit the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *Group) Delete(e ...interface{}) error {
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
	path := c.xpath(names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

// ShowPortMapping returns the service port mappings.
func (c *Group) ShowPortMapping(group interface{}) ([]map[string]string, error) {
	var name string

	switch v := group.(type) {
	case string:
		name = v
	case Entry:
		name = v.Name
	default:
		return nil, fmt.Errorf("Unknown type sent to show port mapping: %s", v)
	}
	c.con.LogOp("(op) showing gke port mappings for %q", name)

	req := pmReq_v1{Group: name}
	ans := pmContainer_v1{}

	if _, err := c.con.Op(req, "", nil, &ans); err != nil {
		return nil, err
	}

	return ans.Normalize(), nil
}

/** Internal functions for this namespace struct **/

func (c *Group) versioning() (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *Group) details(fn util.Retriever, name string) (Entry, error) {
	path := c.xpath([]string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *Group) xpath(vals []string) []string {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"plugins",
		"gcp",
		"gke",
		util.AsEntryXpath(vals),
	}
}
