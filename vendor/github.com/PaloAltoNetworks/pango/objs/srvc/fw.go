package srvc

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// FwSrvc is a namespace struct, included as part of pango.Client.
type FwSrvc struct {
	con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwSrvc) Initialize(con util.XapiClient) {
	c.con = con
}

// GetList performs GET to retrieve a list of service objects.
func (c *FwSrvc) GetList(vsys string) ([]string, error) {
	c.con.LogQuery("(get) list of service objects")
	path := c.xpath(vsys, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// ShowList performs SHOW to retrieve a list of service objects.
func (c *FwSrvc) ShowList(vsys string) ([]string, error) {
	c.con.LogQuery("(show) list of service objects")
	path := c.xpath(vsys, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given service object.
func (c *FwSrvc) Get(vsys, name string) (Entry, error) {
	c.con.LogQuery("(get) service object %q", name)
	listing, err := c.details(c.con.Get, vsys, name)
	if err == nil && len(listing) > 0 {
		return listing[0], nil
	}
	return Entry{}, err
}

// GetAll performs a GET to retrieve services.
func (c *FwSrvc) GetAll(vsys string) ([]Entry, error) {
	c.con.LogQuery("(get) all services")
	return c.details(c.con.Get, vsys, "")
}

// Get performs SHOW to retrieve information for the given service object.
func (c *FwSrvc) Show(vsys, name string) (Entry, error) {
	c.con.LogQuery("(show) service object %q", name)
	listing, err := c.details(c.con.Show, vsys, name)
	if err == nil && len(listing) > 0 {
		return listing[0], nil
	}
	return Entry{}, err
}

// ShowAll performs a SHOW to retrieve services.
func (c *FwSrvc) ShowAll(vsys string) ([]Entry, error) {
	c.con.LogQuery("(show) all services")
	return c.details(c.con.Show, vsys, "")
}

// Set performs SET to create / update one or more service objects.
func (c *FwSrvc) Set(vsys string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	}

	_, fn := c.versioning()
	names := make([]string, len(e))

	// Build up the struct with the given configs.
	d := util.BulkElement{XMLName: xml.Name{Local: "service"}}
	for i := range e {
		d.Data = append(d.Data, fn(e[i]))
		names[i] = e[i].Name
	}
	c.con.LogAction("(set) service objects: %v", names)

	// Set xpath.
	path := c.xpath(vsys, names)
	if len(e) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the objects.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update a service object.
func (c *FwSrvc) Edit(vsys string, e Entry) error {
	var err error

	_, fn := c.versioning()

	c.con.LogAction("(edit) service object %q", e.Name)

	// Set xpath.
	path := c.xpath(vsys, []string{e.Name})

	// Create the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given service objects from the firewall.
//
// Service objects can be either a string or an Entry object.
func (c *FwSrvc) Delete(vsys string, e ...interface{}) error {
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
	c.con.LogAction("(delete) service objects: %v", names)

	path := c.xpath(vsys, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for the FwSrvc struct **/

func (c *FwSrvc) versioning() (normalizer, func(Entry) interface{}) {
	v := c.con.Versioning()

	if v.Gte(version.Number{8, 1, 0, ""}) {
		return &container_v2{}, specify_v2
	} else {
		return &container_v1{}, specify_v1
	}
}

func (c *FwSrvc) details(fn util.Retriever, vsys, name string) ([]Entry, error) {
	path := c.xpath(vsys, []string{name})
	obj, _ := c.versioning()
	_, err := fn(path, nil, obj)
	if err != nil {
		return nil, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *FwSrvc) xpath(vsys string, vals []string) []string {
	if vsys == "" {
		vsys = "vsys1"
	}

	if vsys == "shared" {
		return []string{
			"config",
			"shared",
			"service",
			util.AsEntryXpath(vals),
		}
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"vsys",
		util.AsEntryXpath([]string{vsys}),
		"service",
		util.AsEntryXpath(vals),
	}
}
