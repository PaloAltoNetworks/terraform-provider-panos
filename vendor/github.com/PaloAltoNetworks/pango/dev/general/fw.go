package general

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is a namespace struct, included as part of pango.Client.
type Firewall struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve the device's general settings.
func (c *Firewall) Get() (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve the device's general settings.
func (c *Firewall) Show() (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(), "", ans)
	return first(ans, err)
}

// Set performs SET to create / update the device's general settings.
func (c *Firewall) Set(e Config) error {
	return c.ns.Set(c.pather(), specifier(e))
}

// Edit performs EDIT to update the device's general settings.
func (c *Firewall) Edit(e Config) error {
	return c.ns.Edit(c.pather(), e)
}

/** Internal functions for the Firewall struct **/
func (c *Firewall) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath()
	}
}

func (c *Firewall) xpath() ([]string, error) {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"deviceconfig",
		"system",
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
