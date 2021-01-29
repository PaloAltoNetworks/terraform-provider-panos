package bgp

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.BgpConfig namespace.
type Firewall struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get(vr string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vr), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show(vr string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vr), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vr string, e Config) error {
	return c.ns.Set(c.pather(vr), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vr string, e Config) error {
	return c.ns.Edit(c.pather(vr), e)
}

// Delete performs DELETE to remove the config.
func (c *Firewall) Delete(vr string) error {
	return c.ns.Delete(c.pather(vr), nil, nil)
}

func (c *Firewall) pather(vr string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vr)
	}
}

func (c *Firewall) xpath(vr string) ([]string, error) {
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
		"protocol",
		"bgp",
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
