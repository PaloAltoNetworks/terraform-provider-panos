package nonexist

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.BgpConAdvNonExistFilter namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(vr, ca string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(vr, ca), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(vr, ca string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(vr, ca), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(vr, ca, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vr, ca), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(vr, ca, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vr, ca), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(vr, ca string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(vr, ca), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(vr, ca string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(vr, ca), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vr, ca string, e ...Entry) error {
	return c.ns.Set(c.pather(vr, ca), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vr, ca string, e Entry) error {
	return c.ns.Edit(c.pather(vr, ca), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vr, ca string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vr, ca), names, nErr)
}

func (c *Firewall) pather(vr, ca string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vr, ca, v)
	}
}

func (c *Firewall) xpath(vr, ca string, vals []string) ([]string, error) {
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}
	if ca == "" {
		return nil, fmt.Errorf("ca must be specified")
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath([]string{vr}),
		"protocol",
		"bgp",
		"policy",
		"conditional-advertisement",
		"policy",
		util.AsEntryXpath([]string{ca}),
		"non-exist-filters",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
