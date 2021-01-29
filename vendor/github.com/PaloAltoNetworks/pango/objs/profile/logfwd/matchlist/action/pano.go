package action

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// PanoAction is the client.Objects.LogForwardingProfileMatchListAction namespace.
type PanoAction struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *PanoAction) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of values.
func (c *PanoAction) ShowList(dg, logfwd, matchlist string) ([]string, error) {
	c.con.LogQuery("(show) list of %s", plural)
	path := c.xpath(dg, logfwd, matchlist, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of values.
func (c *PanoAction) GetList(dg, logfwd, matchlist string) ([]string, error) {
	c.con.LogQuery("(get) list of %s", plural)
	path := c.xpath(dg, logfwd, matchlist, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given uid.
func (c *PanoAction) Get(dg, logfwd, matchlist, name string) (Entry, error) {
	c.con.LogQuery("(get) %s %q", singular, name)
	return c.details(c.con.Get, dg, logfwd, matchlist, name)
}

// Show performs SHOW to retrieve information for the given uid.
func (c *PanoAction) Show(dg, logfwd, matchlist, name string) (Entry, error) {
	c.con.LogQuery("(show) %s %q", singular, name)
	return c.details(c.con.Show, dg, logfwd, matchlist, name)
}

// Set performs SET to create / update one or more objects.
func (c *PanoAction) Set(dg, logfwd, matchlist string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
	} else if matchlist == "" {
		return fmt.Errorf("matchlist must be specified")
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
	path := c.xpath(dg, logfwd, matchlist, names)
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
func (c *PanoAction) Edit(dg, logfwd, matchlist string, e Entry) error {
	var err error

	if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
	} else if matchlist == "" {
		return fmt.Errorf("matchlist must be specified")
	}

	_, fn := c.versioning()

	c.con.LogAction("(edit) %s %q", singular, e.Name)

	// Set xpath.
	path := c.xpath(dg, logfwd, matchlist, []string{e.Name})

	// Edit the object.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *PanoAction) Delete(dg, logfwd, matchlist string, e ...interface{}) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if logfwd == "" {
		return fmt.Errorf("logfwd must be specified")
	} else if matchlist == "" {
		return fmt.Errorf("matchlist must be specified")
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
	path := c.xpath(dg, logfwd, matchlist, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *PanoAction) versioning() (normalizer, func(Entry) interface{}) {
	v := c.con.Versioning()

	if v.Gte(version.Number{9, 0, 0, ""}) {
		return &container_v3{}, specify_v3
	} else if v.Gte(version.Number{8, 1, 0, ""}) {
		return &container_v2{}, specify_v2
	} else {
		return &container_v1{}, specify_v1
	}
}

func (c *PanoAction) details(fn util.Retriever, dg, logfwd, matchlist, name string) (Entry, error) {
	path := c.xpath(dg, logfwd, matchlist, []string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *PanoAction) xpath(dg, logfwd, matchlist string, vals []string) []string {
	if dg == "" {
		dg = "shared"
	}

	ans := make([]string, 0, 12)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		"log-settings",
		"profiles",
		util.AsEntryXpath([]string{logfwd}),
		"match-list",
		util.AsEntryXpath([]string{matchlist}),
		"actions",
		util.AsEntryXpath(vals),
	)

	return ans
}
