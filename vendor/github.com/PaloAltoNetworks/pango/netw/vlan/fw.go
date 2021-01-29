package vlan

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.Vlan namespace.
type Firewall struct {
	ns *namespace.Importable
}

/*
SetInterface performs a SET to add an interface to a VLAN.

The VLAN can be either a string or an Entry object.
The iface variable is the interface.
The rmMacs and addMacs params are MAC addresses to remove/add that
will reference the iface interface.
*/
func (c *Firewall) SetInterface(vlan interface{}, iface string, rmMacs, addMacs []string) error {
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

	basePath, err := c.xpath([]string{name})
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
func (c *Firewall) DeleteInterface(vlan interface{}, iface string) error {
	var (
		name string
		err  error
	)

	names, err := toNames([]interface{}{vlan})
	if err != nil {
		return err
	}
	name = names[0]

	o, err := c.Get(name)
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

	basePath, err := c.xpath([]string{name})
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
func (c *Firewall) GetList() ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(), ans)
}

// ShowList performs a SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList() ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(), ans)
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get(name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show(name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve all objects configured.
func (c *Firewall) ShowAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vsys string, e ...Entry) error {
	return c.ns.Set("", "", vsys, c.pather(), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vsys string, e Entry) error {
	return c.ns.Edit("", "", vsys, c.pather(), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete("", "", c.pather(), names, nErr)
}

func (c *Firewall) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(v)
	}
}

func (c *Firewall) xpath(vals []string) ([]string, error) {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"vlan",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
