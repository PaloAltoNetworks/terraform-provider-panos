package neighbor

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.Ipv6NeighborDiscovery namespace.
type Panorama struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, iType, iName, subName string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, iType, iName, subName), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, iType, iName, subName string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, iType, iName, subName), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, iType, iName, subName string, e Config) error {
	return c.ns.Set(c.pather(tmpl, ts, iType, iName, subName), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, iType, iName, subName string, e Config) error {
	return c.ns.Edit(c.pather(tmpl, ts, iType, iName, subName), e)
}

// Delete performs DELETE to remove the config.
func (c *Panorama) Delete(tmpl, ts, iType, iName, subName string) error {
	return c.ns.Delete(c.pather(tmpl, ts, iType, iName, subName), nil, nil)
}

func (c *Panorama) pather(tmpl, ts, iType, iName, subName string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, iType, iName, subName)
	}
}

func (c *Panorama) xpath(tmpl, ts, iType, iName, subName string) ([]string, error) {
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
		} else if subName == "" {
			return nil, fmt.Errorf("subName must be specified")
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

	ans = append(ans,
		"ipv6",
		"neighbor-discovery",
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
