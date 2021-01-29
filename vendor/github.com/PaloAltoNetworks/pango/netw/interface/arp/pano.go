package arp

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.Arp namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list object names.
func (c *Panorama) GetList(tmpl, ts, iType, iName, subName string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts, iType, iName, subName), ans)
}

// ShowList performs SHOW to retrieve a list of object names.
func (c *Panorama) ShowList(tmpl, ts, iType, iName, subName string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts, iType, iName, subName), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(tmpl, ts, iType, iName, subName, ip string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, iType, iName, subName), ip, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(tmpl, ts, iType, iName, subName, ip string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, iType, iName, subName), ip, ans)
	return first(ans, err)
}

// GetAll performs a GET to retrieve information for all objects.
func (c *Panorama) GetAll(tmpl, ts, iType, iName, subName string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, iType, iName, subName), ans)
	return all(ans, err)
}

// ShowAll performs a SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(tmpl, ts, iType, iName, subName string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, iType, iName, subName), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, iType, iName, subName string, e ...Entry) error {
	return c.ns.Set(c.pather(tmpl, ts, iType, iName, subName), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, iType, iName, subName string, e Entry) error {
	return c.ns.Edit(c.pather(tmpl, ts, iType, iName, subName), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts, iType, iName, subName string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(tmpl, ts, iType, iName, subName), names, nErr)
}

func (c *Panorama) pather(tmpl, ts, iType, iName, subName string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, iType, iName, subName, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, iType, iName, subName string, vals []string) ([]string, error) {
	// Sanity checks.
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	switch iType {
	case "":
		return nil, fmt.Errorf("iType must be specified")
	case TypeVlan:
		if iName != "" {
			return nil, fmt.Errorf("iName should be an empty string for %s types", iType)
		}
	case TypeEthernet, TypeAggregate:
		if iName == "" {
			return nil, fmt.Errorf("iName must be specified")
		}
	default:
		return nil, fmt.Errorf("unknown iType value: %s", iType)
	}

	ans := make([]string, 0, 17)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"interface",
		iType,
	)

	if iType != TypeVlan {
		ans = append(ans, util.AsEntryXpath([]string{iName}), "layer3")
	}

	if subName != "" {
		ans = append(ans, "units", util.AsEntryXpath([]string{subName}))
	}

	ans = append(ans, "arp", util.AsEntryXpath(vals))

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
