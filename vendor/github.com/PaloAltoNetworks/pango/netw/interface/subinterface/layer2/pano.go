package layer2

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.Layer2Subinterface namespace.
type Panorama struct {
	ns *namespace.Importable
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts, iType, eth, mType string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts, iType, eth, mType), ans)
}

// ShowList performs a SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts, iType, eth, mType string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts, iType, eth, mType), ans)
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, iType, eth, mType, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, iType, eth, mType), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, iType, eth, mType, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, iType, eth, mType), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts, iType, eth, mType string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, iType, eth, mType), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve all objects configured.
func (c *Panorama) ShowAll(tmpl, ts, iType, eth, mType string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts, iType, eth, mType), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, iType, eth, mType, vsys string, e ...Entry) error {
	return c.ns.Set(tmpl, ts, vsys, c.pather(tmpl, ts, iType, eth, mType), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, iType, eth, mType, vsys string, e Entry) error {
	return c.ns.Edit(tmpl, ts, vsys, c.pather(tmpl, ts, iType, eth, mType), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts, iType, eth, mType string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(tmpl, ts, c.pather(tmpl, ts, iType, eth, mType), names, nErr)
}

func (c *Panorama) pather(tmpl, ts, iType, eth, mType string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, iType, eth, mType, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, iType, eth, mType string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	switch iType {
	case "":
		return nil, fmt.Errorf("iType must be specified")
	case EthernetInterface, AggregateInterface:
	default:
		return nil, fmt.Errorf("unknown iType value: %s", iType)
	}
	if eth == "" {
		return nil, fmt.Errorf("eth must be specified")
	}
	switch mType {
	case "":
		return nil, fmt.Errorf("mType must be specified")
	case VirtualWire, Layer2:
	default:
		return nil, fmt.Errorf("unknown mType value: %s", mType)
	}

	ans := make([]string, 0, 15)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"interface",
		iType,
		util.AsEntryXpath([]string{eth}),
		mType,
		"units",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
