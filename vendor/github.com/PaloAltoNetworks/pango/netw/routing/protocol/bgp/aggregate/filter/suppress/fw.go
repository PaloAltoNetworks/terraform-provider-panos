package suppress

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.BgpAggSuppressFilter namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(vr, ag string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(vr, ag), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(vr, ag string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(vr, ag), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(vr, ag, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vr, ag), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(vr, ag, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vr, ag), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(vr, ag string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(vr, ag), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(vr, ag string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(vr, ag), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vr, ag string, e ...Entry) error {
	return c.ns.Set(c.pather(vr, ag), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vr, ag string, e Entry) error {
	return c.ns.Edit(c.pather(vr, ag), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vr, ag string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vr, ag), names, nErr)
}

func (c *Firewall) pather(vr, ag string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vr, ag, v)
	}
}

func (c *Firewall) xpath(vr, ag string, vals []string) ([]string, error) {
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}
	if ag == "" {
		return nil, fmt.Errorf("ag must be specified")
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
		"aggregation",
		"address",
		util.AsEntryXpath([]string{ag}),
		"suppress-filters",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
