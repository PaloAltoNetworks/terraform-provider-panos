package bgp

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.BgpConfig namespace.
type Panorama struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, vr string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, vr), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, vr string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, vr), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vr string, e Config) error {
	return c.ns.Set(c.pather(tmpl, ts, vr), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vr string, e Config) error {
	return c.ns.Edit(c.pather(tmpl, ts, vr), e)
}

// Delete performs DELETE to remove the config.
func (c *Panorama) Delete(tmpl, ts, vr string) error {
	return c.ns.Delete(c.pather(tmpl, ts, vr), nil, nil)
}

func (c *Panorama) pather(tmpl, ts, vr string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vr)
	}
}

func (c *Panorama) xpath(tmpl, ts, vr string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}

	ans := make([]string, 0, 13)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath([]string{vr}),
		"protocol",
		"bgp",
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
