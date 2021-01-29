package matchlist

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

// FwMatchList is the client.Objects.LogForwardingProfileMatchList namespace.
type FwMatchList struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *FwMatchList) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *FwMatchList) ShowList(vsys, logfwd string) ([]string, error) {
	c.con.LogQuery("(show) list of %s", plural)
	path := c.xpath(vsys, logfwd, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of values.
func (c *FwMatchList) GetList(vsys, logfwd string) ([]string, error) {
	c.con.LogQuery("(get) list of %s", plural)
	path := c.xpath(vsys, logfwd, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given uid.
func (c *FwMatchList) Get(vsys, logfwd, name string) (Entry, error) {
	c.con.LogQuery("(get) %s %q", singular, name)
	return c.details(c.con.Get, vsys, logfwd, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *FwMatchList) Show(vsys, logfwd, name string) (Entry, error) {
	c.con.LogQuery("(show) %s %q", singular, name)
	return c.details(c.con.Show, vsys, logfwd, name)
}

// Set performs SET to create / update one or more objects.
func (c *FwMatchList) Set(vsys, logfwd string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
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
	path := c.xpath(vsys, logfwd, names)
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
func (c *FwMatchList) Edit(vsys, logfwd string, e Entry) error {
	var err error

	if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
	}

	_, fn := c.versioning()

	c.con.LogAction("(edit) %s %q", singular, e.Name)

	// Set xpath.
	path := c.xpath(vsys, logfwd, []string{e.Name})

	// Edit the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *FwMatchList) Delete(vsys, logfwd string, e ...interface{}) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
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
	path := c.xpath(vsys, logfwd, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *FwMatchList) versioning() (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *FwMatchList) details(fn util.Retriever, vsys, logfwd, name string) (Entry, error) {
	path := c.xpath(vsys, logfwd, []string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *FwMatchList) xpath(vsys, logfwd string, vals []string) []string {
	if vsys == "" {
		vsys = "shared"
	}

	ans := make([]string, 0, 10)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"log-settings",
		"profiles",
		util.AsEntryXpath([]string{logfwd}),
		"match-list",
		util.AsEntryXpath(vals),
	)

	return ans
}
