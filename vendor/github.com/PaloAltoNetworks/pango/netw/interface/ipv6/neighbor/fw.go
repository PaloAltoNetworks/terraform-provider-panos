package neighbor

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.Ipv6NeighborDiscovery namespace.
type Firewall struct {
	ns *namespace.Standard
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get(iType, iName, subName string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(iType, iName, subName), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show(iType, iName, subName string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(iType, iName, subName), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(iType, iName, subName string, e Config) error {
	return c.ns.Set(c.pather(iType, iName, subName), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(iType, iName, subName string, e Config) error {
	return c.ns.Edit(c.pather(iType, iName, subName), e)
}

// Delete performs DELETE to remove the config.
func (c *Firewall) Delete(iType, iName, subName string) error {
	return c.ns.Delete(c.pather(iType, iName, subName), nil, nil)
}

func (c *Firewall) pather(iType, iName, subName string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(iType, iName, subName)
	}
}

func (c *Firewall) xpath(iType, iName, subName string) ([]string, error) {
	// Sanity checks.
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

	ans := make([]string, 0, 12)
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

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
