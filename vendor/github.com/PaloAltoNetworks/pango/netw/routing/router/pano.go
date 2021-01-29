package router

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.VirtualRouter namespace.
type Panorama struct {
	ns *namespace.Importable
}

/*
SetInterface performs a SET to add an interface to a virtual router.

The virtual router can be either a string or an Entry object.
*/
func (c *Panorama) SetInterface(tmpl, ts string, vr interface{}, iface string) error {
	names, err := toNames([]interface{}{vr})
	if err != nil {
		return err
	}
	name := names[0]

	c.ns.Client.LogAction("(set) interface for %s %q: %s", singular, name, iface)

	path, err := c.xpath(tmpl, ts, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "interface")

	_, err = c.ns.Client.Set(path, util.Member{Value: iface}, nil, nil)
	return err
}

/*
DeleteInterface performs a DELETE to remove an interface from a virtual router.

The virtual router can be either a string or an Entry object.
*/
func (c *Panorama) DeleteInterface(tmpl, ts string, vr interface{}, iface string) error {
	names, err := toNames([]interface{}{vr})
	if err != nil {
		return err
	}
	name := names[0]

	c.ns.Client.LogAction("(delete) interface for %s %q: %s", singular, name, iface)

	path, err := c.xpath(tmpl, ts, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "interface", util.AsMemberXpath([]string{iface}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts), ans)
}

// ShowList performs a SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts), ans)
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve all objects configured.
func (c *Panorama) ShowAll(tmpl, ts string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vsys string, e ...Entry) error {
	return c.ns.Set(tmpl, ts, vsys, c.pather(tmpl, ts), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vsys string, e Entry) error {
	return c.ns.Edit(tmpl, ts, vsys, c.pather(tmpl, ts), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(tmpl, ts, c.pather(tmpl, ts), names, nErr)
}

// CleanupDefault clears the `default` route configuration instead of deleting
// it outright.  This involves unimporting the route "default" from the given
// vsys, then performing an `EDIT` with an empty router.Entry object.
func (c *Panorama) CleanupDefault(tmpl, ts string) error {
	c.ns.Client.LogAction("(action) cleaning up %s: default", c.ns.Singular)

	// Cleanup the interfaces the virtual router refers to.
	info := Entry{Name: "default"}
	return c.Edit(tmpl, ts, "", info)
}

func (c *Panorama) pather(tmpl, ts string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, v)
	}
}

func (c *Panorama) xpath(tmpl, ts string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}

	ans := make([]string, 0, 11)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
