package suppress

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.BgpAggSuppressFilter namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts, vr, ag string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts, vr, ag), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts, vr, ag string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts, vr, ag), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(tmpl, ts, vr, ag, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, vr, ag), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(tmpl, ts, vr, ag, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, vr, ag), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts, vr, ag string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, vr, ag), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(tmpl, ts, vr, ag string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts, vr, ag), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vr, ag string, e ...Entry) error {
	return c.ns.Set(c.pather(tmpl, ts, vr, ag), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vr, ag string, e Entry) error {
	return c.ns.Edit(c.pather(tmpl, ts, vr, ag), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts, vr, ag string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(tmpl, ts, vr, ag), names, nErr)
}

func (c *Panorama) pather(tmpl, ts, vr, ag string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vr, ag, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, vr, ag string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}
	if ag == "" {
		return nil, fmt.Errorf("ag must be specified")
	}

	ans := make([]string, 0, 19)
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
		"policy",
		"aggregation",
		"address",
		util.AsEntryXpath([]string{ag}),
		"suppress-filters",
		util.AsEntryXpath(vals),
	)
	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
