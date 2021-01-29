package server

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

// FwServer is the client.Network.SyslogServer namespace.
type FwServer struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwServer) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwServer) ShowList(vsys, profile string) ([]string, error) {
	c.con.LogQuery("(show) list of %s", plural)
	path := c.xpath(vsys, profile, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwServer) GetList(vsys, profile string) ([]string, error) {
	c.con.LogQuery("(get) list of %s", plural)
	path := c.xpath(vsys, profile, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwServer) Get(vsys, profile, name string) (Entry, error) {
	c.con.LogQuery("(get) %s %q", singular, name)
	return c.details(c.con.Get, vsys, profile, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwServer) Show(vsys, profile, name string) (Entry, error) {
	c.con.LogQuery("(show) %s %q", singular, name)
	return c.details(c.con.Show, vsys, profile, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwServer) Set(vsys, profile string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if profile == "" {
		return fmt.Errorf("profile must be specified")
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
	path := c.xpath(vsys, profile, names)
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
func (c *FwServer) Edit(vsys, profile string, e Entry) error {
	var err error

	if profile == "" {
		return fmt.Errorf("profile must be specified")
	}

	_, fn := c.versioning()

	c.con.LogAction("(edit) %s %q", singular, e.Name)

	// Set xpath.
	path := c.xpath(vsys, profile, []string{e.Name})

	// Edit the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwServer) Delete(vsys, profile string, e ...interface{}) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if profile == "" {
		return fmt.Errorf("profile must be specified")
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
	path := c.xpath(vsys, profile, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *FwServer) versioning() (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *FwServer) details(fn util.Retriever, vsys, profile, name string) (Entry, error) {
	path := c.xpath(vsys, profile, []string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *FwServer) xpath(vsys, profile string, vals []string) []string {
	if vsys == "" {
		vsys = "shared"
	}

	ans := make([]string, 0, 10)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"log-settings",
		"syslog",
		util.AsEntryXpath([]string{profile}),
		"server",
		util.AsEntryXpath(vals),
	)

	return ans
}
