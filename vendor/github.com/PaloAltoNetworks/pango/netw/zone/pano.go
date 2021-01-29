package zone

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.Zone namespace.
type Panorama struct {
	ns *namespace.Standard
}

/*
SetInterface performs a SET to add an interface to a zone.

The zone can be either a string or an Entry object.
*/
func (c *Panorama) SetInterface(tmpl, ts, vsys string, zone interface{}, mode, iface string) error {
	names, err := toNames([]interface{}{zone})
	if err != nil {
		return err
	}
	name := names[0]

	switch mode {
	case ModeL2, ModeL3, ModeVirtualWire, ModeTap, ModeExternal:
	default:
		return fmt.Errorf("unknown mode value: %s", mode)
	}

	c.ns.Client.LogAction("(set) %s interface: %s", singular, name)

	path, err := c.xpath(tmpl, ts, vsys, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "network", mode)

	_, err = c.ns.Client.Set(path, util.Member{Value: iface}, nil, nil)
	return err
}

/*
DeleteInterface performs a DELETE to remove the interface from the zone.

The zone can be either a string or an Entry object.
*/
func (c *Panorama) DeleteInterface(tmpl, ts, vsys string, zone interface{}, mode, iface string) error {
	names, err := toNames([]interface{}{zone})
	if err != nil {
		return err
	}
	name := names[0]

	switch mode {
	case ModeL2, ModeL3, ModeVirtualWire, ModeTap, ModeExternal:
	default:
		return fmt.Errorf("unknown mode value: %s", mode)
	}

	c.ns.Client.LogAction("(delete) %s interface: %s", singular, name)

	path, err := c.xpath(tmpl, ts, vsys, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "network", mode, util.AsMemberXpath([]string{iface}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts, vsys), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts, vsys), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, vsys), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts, vsys), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vsys string, e ...Entry) error {
	return c.ns.Set(c.pather(tmpl, ts, vsys), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vsys string, e Entry) error {
	return c.ns.Edit(c.pather(tmpl, ts, vsys), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts, vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(tmpl, ts, vsys), names, nErr)
}

func (c *Panorama) pather(tmpl, ts, vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vsys, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, vsys string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	if vsys == "shared" {
		return nil, fmt.Errorf("vsys cannot be 'shared'")
	}

	ans := make([]string, 0, 12)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"zone",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
