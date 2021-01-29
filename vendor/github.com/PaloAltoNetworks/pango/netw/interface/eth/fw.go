package eth

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.EthernetInterface namespace.
type Firewall struct {
	ns *namespace.Importable
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
		"interface",
		"ethernet",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
