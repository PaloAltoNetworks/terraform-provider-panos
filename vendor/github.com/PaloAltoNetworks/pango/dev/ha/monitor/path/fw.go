package path

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Device.HaPathMonitorGroup namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(gType string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(gType), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(gType string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(gType), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(gType, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(gType), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(gType, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(gType), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(gType string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(gType), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(gType string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(gType), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(gType string, e ...Entry) error {
	return c.ns.Set(c.pather(gType), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(gType string, e Entry) error {
	return c.ns.Edit(c.pather(gType), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(gType string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(gType), names, nErr)
}

func (c *Firewall) pather(gType string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(gType, v)
	}
}

func (c *Firewall) xpath(gType string, vals []string) ([]string, error) {
	switch gType {
	case "":
		return nil, fmt.Errorf("gType must be specified")
	case VirtualWire, Vlan, VirtualRouter, LogicalRouter:
	default:
		return nil, fmt.Errorf("unknown gType value: %s", gType)
	}
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"deviceconfig",
		"high-availability",
		"group",
		"monitoring",
		"path-monitoring",
		"path-group",
		gType,
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
