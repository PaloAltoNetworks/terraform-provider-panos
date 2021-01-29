package vlan

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.Vlan namespace.
type Panorama struct {
	ns *namespace.Importable
}

/*
SetInterface performs a SET to add an interface to a VLAN.

The VLAN can be either a string or an Entry object.
The iface variable is the interface.
The rmMacs and addMacs params are MAC addresses to remove/add that
will reference the iface interface.
*/
func (c *Panorama) SetInterface(tmpl, ts string, vlan interface{}, iface string, rmMacs, addMacs []string) error {
	var (
		name string
		err  error
	)

	names, err := toNames([]interface{}{vlan})
	if err != nil {
		return err
	}
	name = names[0]

	c.ns.Client.LogAction("(set) interface for %s %q: %s", singular, name, iface)

	basePath, err := c.xpath(tmpl, ts, []string{name})
	if err != nil {
		return err
	}
	iPath := append(basePath, "interface")

	if _, err = c.ns.Client.Set(iPath, util.Member{Value: iface}, nil, nil); err != nil {
		return err
	}

	if len(rmMacs) > 0 {
		c.ns.Client.LogAction("(delete) removing %q mac addresses: %#v", name, rmMacs)
		rPath := append(basePath, "mac", util.AsEntryXpath(rmMacs))
		if _, err = c.ns.Client.Delete(rPath, nil, nil); err != nil {
			return err
		}
	}

	if len(addMacs) > 0 {
		c.ns.Client.LogAction("(set) adding %q mac addresses: %#v", name, addMacs)
		d := util.BulkElement{XMLName: xml.Name{Local: "mac"}}
		for i := range addMacs {
			d.Data = append(d.Data, macList{Mac: addMacs[i], Interface: iface})
		}
		aPath := make([]string, 0, len(basePath)+1)
		aPath = append(aPath, basePath...)
		if len(addMacs) == 1 {
			aPath = append(aPath, "mac")
		}
		if _, err = c.ns.Client.Set(aPath, d.Config(), nil, nil); err != nil {
			return err
		}
	}

	return nil
}

/*
DeleteInterface performs a DELETE to remove an interface from a VLAN.

The VLAN can be either a string or an Entry object.

All MAC addresses referencing this interface are deleted.
*/
func (c *Panorama) DeleteInterface(tmpl, ts string, vlan interface{}, iface string) error {
	var (
		name string
		err  error
	)

	names, err := toNames([]interface{}{vlan})
	if err != nil {
		return err
	}
	name = names[0]

	o, err := c.Get(tmpl, ts, name)
	if err != nil {
		return err
	}
	rmMacs := make([]string, 0, len(o.StaticMacs))
	for k, v := range o.StaticMacs {
		if v == iface {
			rmMacs = append(rmMacs, k)
		}
	}

	c.ns.Client.LogAction("(delete) interface for %s %q: %s", singular, name, iface)

	basePath, err := c.xpath(tmpl, ts, []string{name})
	if err != nil {
		return err
	}
	mPath := append(basePath, "mac", util.AsEntryXpath(rmMacs))
	iPath := append(basePath, "interface", util.AsMemberXpath([]string{iface}))

	if len(rmMacs) > 0 {
		c.ns.Client.LogAction("(delete) removing %q mac addresses: %#v", iface, rmMacs)
		if _, err = c.ns.Client.Delete(mPath, nil, nil); err != nil {
			return err
		}
	}

	_, err = c.ns.Client.Delete(iPath, nil, nil)
	return err
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts), ans)
}

// ShowList performs a SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts), ans)
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve all objects configured.
func (c *Panorama) ShowAll(tmpl, ts string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vsys string, e ...Entry) error {
	return c.ns.Set(tmpl, ts, vsys, c.pather(tmpl, ts), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vsys string, e Entry) error {
	return c.ns.Edit(tmpl, ts, vsys, c.pather(tmpl, ts), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(tmpl, ts, c.pather(tmpl, ts), names, nErr)
}

func (c *Panorama) pather(tmpl, ts string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, v)
	}
}

func (c *Panorama) xpath(tmpl, ts string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}

	ans := make([]string, 0, 11)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"vlan",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
