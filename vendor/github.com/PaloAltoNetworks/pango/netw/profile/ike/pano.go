package ike

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// PanoIke is a namespace struct, included as part of pango.Client.
type PanoIke struct {
	con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *PanoIke) Initialize(con util.XapiClient) {
	c.con = con
}

// GetList performs GET to retrieve a list of IKE crypto profiles.
func (c *PanoIke) GetList(tmpl, ts string) ([]string, error) {
	c.con.LogQuery("(get) list of ike crypto profiles")
	path := c.xpath(tmpl, ts, nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// ShowList performs SHOW to retrieve a list of IKE crypto profiles.
func (c *PanoIke) ShowList(tmpl, ts string) ([]string, error) {
	c.con.LogQuery("(show) list of ike crypto profiles")
	path := c.xpath(tmpl, ts, nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given IKE crypto
// profile.
func (c *PanoIke) Get(tmpl, ts, name string) (Entry, error) {
	c.con.LogQuery("(get) ike crypto profile %q", name)
	return c.details(c.con.Get, tmpl, ts, name)
}

// Get performs SHOW to retrieve information for the given IKE crypto
// profile.
func (c *PanoIke) Show(tmpl, ts, name string) (Entry, error) {
	c.con.LogQuery("(show) ike crypto profile %q", name)
	return c.details(c.con.Show, tmpl, ts, name)
}

// Set performs SET to create / update one or more IKE crypto profiles.
func (c *PanoIke) Set(tmpl, ts string, e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	} else if tmpl == "" && ts == "" {
		return fmt.Errorf("tmpl or ts must be specified")
	}

	_, fn, vint := c.versioning()
	names := make([]string, len(e))

	// Build up the struct with the given configs.
	se := make([]Entry, len(e))
	d := util.BulkElement{XMLName: xml.Name{Local: "ike-crypto-profiles"}}
	for i := range e {
		se[i].Name = e[i].Name
		se[i].Copy(e[i])
		se[i].SpecifyEncryption(vint)
		d.Data = append(d.Data, fn(se[i]))
		names[i] = se[i].Name
	}
	c.con.LogAction("(set) ike crypto profiles: %v", names)

	// Set xpath.
	path := c.xpath(tmpl, ts, names)
	if len(se) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the profiles.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update an IKE crypto profile.
func (c *PanoIke) Edit(tmpl, ts string, e Entry) error {
	var err error

	if tmpl == "" && ts == "" {
		return fmt.Errorf("tmpl or ts must be specified")
	}

	_, fn, vint := c.versioning()

	c.con.LogAction("(edit) ike crypto profile %q", e.Name)

	se := Entry{Name: e.Name}
	se.Copy(e)
	se.SpecifyEncryption(vint)

	// Set xpath.
	path := c.xpath(tmpl, ts, []string{se.Name})

	// Edit the profile.
	_, err = c.con.Edit(path, fn(se), nil, nil)
	return err
}

// Delete removes the given IKE crypto profile(s) from the firewall.
//
// Profiles can be either a string or an Entry object.
func (c *PanoIke) Delete(tmpl, ts string, e ...interface{}) error {
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
			return fmt.Errorf("Unsupported type to delete: %s", v)
		}
	}
	c.con.LogAction("(delete) ike crypto profiles: %v", names)

	path := c.xpath(tmpl, ts, names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *PanoIke) versioning() (normalizer, func(Entry) interface{}, int) {
	v := c.con.Versioning()

	if v.Gte(version.Number{7, 0, 0, ""}) {
		return &container_v2{}, specify_v2, 2
	} else {
		return &container_v1{}, specify_v1, 1
	}
}

func (c *PanoIke) details(fn util.Retriever, tmpl, ts, name string) (Entry, error) {
	path := c.xpath(tmpl, ts, []string{name})
	obj, _, _ := c.versioning()
	_, err := fn(path, nil, obj)
	if err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()
	ans.NormalizeEncryption()

	return ans, nil
}

func (c *PanoIke) xpath(tmpl, ts string, vals []string) []string {
	ans := make([]string, 0, 13)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"ike",
		"crypto-profiles",
		"ike-crypto-profiles",
		util.AsEntryXpath(vals),
	)

	return ans
}
