package variable

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
)

// Variable is the client.Panorama.TemplateVariable namespace.
type Variable struct {
	con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Variable) Initialize(con util.XapiClient) {
	c.con = con
}

// ShowList performs SHOW to retrieve a list of variables.
func (c *Variable) ShowList(tmpl, ts string) ([]string, error) {
	c.con.LogQuery("(show) list of template variables")
	path := c.xpath(tmpl, ts, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// GetList performs GET to retrieve a list of variables.
func (c *Variable) GetList(tmpl, ts string) ([]string, error) {
	c.con.LogQuery("(get) list of template variables")
	path := c.xpath(tmpl, ts, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given variable.
func (c *Variable) Get(tmpl, ts, name string) (Entry, error) {
	c.con.LogQuery("(get) template variable %q", name)
	return c.details(c.con.Get, tmpl, ts, name)
}

// Show performs SHOW to retrieve information for the given variable.
func (c *Variable) Show(tmpl, ts, name string) (Entry, error) {
	c.con.LogQuery("(show) template variable %q", name)
	return c.details(c.con.Show, tmpl, ts, name)
}

// Set performs SET to create / update one or more template variables.
func (c *Variable) Set(tmpl, ts string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if tmpl == "" && ts == "" {
		return fmt.Errorf("tmpl or ts must be specified")
	}

	_, fn := c.versioning()
	names := make([]string, len(e))

	// Build up the struct with the given configs.
	d := util.BulkElement{XMLName: xml.Name{Local: "variable"}}
	for i := range e {
		d.Data = append(d.Data, fn(e[i]))
		names[i] = e[i].Name
	}
	c.con.LogAction("(set) template variables: %v", names)

	// Set xpath.
	path := c.xpath(tmpl, ts, names)
	if len(e) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the template variables.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update a template.
func (c *Variable) Edit(tmpl, ts string, e Entry) error {
	var err error

	if tmpl == "" && ts == "" {
		return fmt.Errorf("tmpl or ts must be specified")
	}

	_, fn := c.versioning()

	c.con.LogAction("(edit) template variable %q", e.Name)

	// Set xpath.
	path := c.xpath(tmpl, ts, []string{e.Name})

	// Edit the template.
	_, err = c.con.Edit(path, fn(e), nil, nil)
	return err
}

// Delete removes the given template variables from the firewall.
//
// Variables can be a string or an Entry object.
func (c *Variable) Delete(tmpl, ts string, e ...interface{}) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if tmpl == "" && ts == "" {
		return fmt.Errorf("tmpl or ts must be specified")
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
	c.con.LogAction("(delete) template variables: %v", names)

	// Remove the template variables.
	path := c.xpath(tmpl, ts, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *Variable) versioning() (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func (c *Variable) details(fn util.Retriever, tmpl, ts, name string) (Entry, error) {
	path := c.xpath(tmpl, ts, []string{name})
	obj, _ := c.versioning()
	if _, err := fn(path, nil, obj); err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()

	return ans, nil
}

func (c *Variable) xpath(tmpl, ts string, vals []string) []string {
	ans := make([]string, 0, 7)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"variable",
		util.AsEntryXpath(vals),
	)

	return ans
}
