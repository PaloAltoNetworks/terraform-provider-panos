package ipsec

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// FwIpsec is a namespace struct, included as part of pango.Client.
type FwIpsec struct {
	con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *FwIpsec) Initialize(con util.XapiClient) {
	c.con = con
}

// GetList performs GET to retrieve a list of IPSec crypto profiles.
func (c *FwIpsec) GetList() ([]string, error) {
	c.con.LogQuery("(get) list of ipsec crypto profiles")
	path := c.xpath(nil)
	return c.con.EntryListUsing(c.con.Get, path[:len(path)-1])
}

// ShowList performs SHOW to retrieve a list of IPSec crypto profiles.
func (c *FwIpsec) ShowList() ([]string, error) {
	c.con.LogQuery("(show) list of ipsec crypto profiles")
	path := c.xpath(nil)
	return c.con.EntryListUsing(c.con.Show, path[:len(path)-1])
}

// Get performs GET to retrieve information for the given IPSec crypto
// profile.
func (c *FwIpsec) Get(name string) (Entry, error) {
	c.con.LogQuery("(get) ipsec crypto profile %q", name)
	return c.details(c.con.Get, name)
}

// Get performs SHOW to retrieve information for the given IPSec crypto
// profile.
func (c *FwIpsec) Show(name string) (Entry, error) {
	c.con.LogQuery("(show) ipsec crypto profile %q", name)
	return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more IPSec crypto profiles.
func (c *FwIpsec) Set(e ...Entry) error {
	var err error

	if len(e) == 0 {
		return nil
	}

	_, fn, vint := c.versioning()
	names := make([]string, len(e))

	// Build up the struct with the given configs.
	se := make([]Entry, len(e))
	d := util.BulkElement{XMLName: xml.Name{Local: "ipsec-crypto-profiles"}}
	for i := range e {
		se[i].Name = e[i].Name
		se[i].Copy(e[i])
		se[i].SpecifyEncryption(vint)
		d.Data = append(d.Data, fn(se[i]))
		names[i] = se[i].Name
	}
	c.con.LogAction("(set) ipsec crypto profiles: %v", names)

	// Set xpath.
	path := c.xpath(names)
	if len(se) == 1 {
		path = path[:len(path)-1]
	} else {
		path = path[:len(path)-2]
	}

	// Create the profiles.
	_, err = c.con.Set(path, d.Config(), nil, nil)
	return err
}

// Edit performs EDIT to create / update an IPSec crypto profile.
func (c *FwIpsec) Edit(e Entry) error {
	var err error

	_, fn, vint := c.versioning()

	c.con.LogAction("(edit) ipsec crypto profile %q", e.Name)

	se := Entry{Name: e.Name}
	se.Copy(e)
	se.SpecifyEncryption(vint)

	// Set xpath.
	path := c.xpath([]string{se.Name})

	// Edit the profile.
	_, err = c.con.Edit(path, fn(se), nil, nil)
	return err
}

// Delete removes the given IPSec crypto profile(s) from the firewall.
//
// Profiles can be either a string or an Entry object.
func (c *FwIpsec) Delete(e ...interface{}) error {
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
	c.con.LogAction("(delete) ipsec crypto profiles: %v", names)

	path := c.xpath(names)
	_, err = c.con.Delete(path, nil, nil)
	return err
}

/** Internal functions for this namespace struct **/

func (c *FwIpsec) versioning() (normalizer, func(Entry) interface{}, int) {
	v := c.con.Versioning()

	var vint int

	if v.Gte(version.Number{7, 0, 0, ""}) {
		vint = 2
	} else {
		vint = 1
	}

	return &container_v1{}, specify_v1, vint
}

func (c *FwIpsec) details(fn util.Retriever, name string) (Entry, error) {
	path := c.xpath([]string{name})
	obj, _, _ := c.versioning()
	_, err := fn(path, nil, obj)
	if err != nil {
		return Entry{}, err
	}
	ans := obj.Normalize()
	ans.NormalizeEncryption()

	return ans, nil
}

func (c *FwIpsec) xpath(vals []string) []string {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"ike",
		"crypto-profiles",
		"ipsec-crypto-profiles",
		util.AsEntryXpath(vals),
	}
}
