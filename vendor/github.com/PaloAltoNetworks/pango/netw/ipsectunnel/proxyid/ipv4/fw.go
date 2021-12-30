package ipv4

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.IpsecTunnelProxyId namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(tun string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tun), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(tun string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tun), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(tun, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tun), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(tun, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tun), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(tun string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tun), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(tun string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tun), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(tun string, e ...Entry) error {
	return c.ns.Set(c.pather(tun), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(tun string, e Entry) error {
	return c.ns.Edit(c.pather(tun), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(tun string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(tun), names, nErr)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Firewall) FromPanosConfig(tun, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(tun), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Firewall) AllFromPanosConfig(tun string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(tun), ans)
	return all(ans, err)
}

func (c *Firewall) pather(tun string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tun, v)
	}
}

func (c *Firewall) xpath(tun string, vals []string) ([]string, error) {
	if tun == "" {
		return nil, fmt.Errorf("tun must be specified")
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"tunnel",
		"ipsec",
		util.AsEntryXpath([]string{tun}),
		"auto-key",
		"proxy-id",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
