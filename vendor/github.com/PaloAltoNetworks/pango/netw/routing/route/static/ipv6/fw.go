package ipv6

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.Ipv6StaticRoute namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(vr string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(vr), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(vr string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(vr), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(vr, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vr), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(vr, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vr), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(vr string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(vr), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(vr string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(vr), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vr string, e ...Entry) error {
	return c.ns.Set(c.pather(vr), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vr string, e Entry) error {
	return c.ns.Edit(c.pather(vr), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vr string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vr), names, nErr)
}

func (c *Firewall) pather(vr string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vr, v)
	}
}

func (c *Firewall) xpath(vr string, vals []string) ([]string, error) {
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath([]string{vr}),
		"routing-table",
		"ipv6",
		"static-route",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
