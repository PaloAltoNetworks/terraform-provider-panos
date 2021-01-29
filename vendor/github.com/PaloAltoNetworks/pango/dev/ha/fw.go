package ha

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Device.HaConfig namespace.
type Firewall struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get() (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show() (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(e Config) error {
	return c.ns.Set(c.pather(), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(e Config) error {
	return c.ns.Edit(c.pather(), e)
}

// Delete performs DELETE to remove the config.
func (c *Firewall) Delete() error {
	return c.ns.Delete(c.pather(), nil, nil)
}

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
		"high-availability",
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
