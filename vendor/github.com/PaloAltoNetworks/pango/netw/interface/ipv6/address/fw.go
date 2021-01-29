package address

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Network.Ipv6Address namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(iType, iName, subName string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(iType, iName, subName), ans)
}

// ShowList performs a SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(iType, iName, subName string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(iType, iName, subName), ans)
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get(iType, iName, subName, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(iType, iName, subName), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show(iType, iName, subName, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(iType, iName, subName), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(iType, iName, subName string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(iType, iName, subName), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve all objects configured.
func (c *Firewall) ShowAll(iType, iName, subName string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(iType, iName, subName), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(iType, iName, subName string, e ...Entry) error {
	return c.ns.Set(c.pather(iType, iName, subName), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(iType, iName, subName string, e Entry) error {
	return c.ns.Edit(c.pather(iType, iName, subName), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(iType, iName, subName string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(iType, iName, subName), names, nErr)
}

func (c *Firewall) pather(iType, iName, subName string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(iType, iName, subName, v)
	}
}

func (c *Firewall) xpath(iType, iName, subName string, vals []string) ([]string, error) {
	// Sanity checks.
	switch iType {
	case "":
		return nil, fmt.Errorf("iType must be specified")
	case TypeVlan, TypeTunnel, TypeLoopback:
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

	ans := make([]string, 0, 13)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"interface",
		iType,
	)

	if iType == TypeEthernet || iType == TypeAggregate {
		ans = append(ans, util.AsEntryXpath([]string{iName}), "layer3")
	}

	if subName != "" {
		ans = append(ans, "units", util.AsEntryXpath([]string{subName}))
	}

	ans = append(ans,
		"ipv6",
		"address",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
