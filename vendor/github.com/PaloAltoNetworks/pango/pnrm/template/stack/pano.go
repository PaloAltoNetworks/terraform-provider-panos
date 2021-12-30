package stack

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Panorama.TemplateStack namespace.
type Panorama struct {
	ns *namespace.Standard
}

/*
SetDevice performs a SET to add specific device to template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Panorama) SetDevice(st interface{}, d string) error {
	var name string

	switch v := st.(type) {
	case string:
		name = v
	case Entry:
		name = v.Name
	default:
		return fmt.Errorf("Unknown type sent to set device: %s", v)
	}

	c.ns.Client.LogAction("(set) device in template stack: %s", name)

	path, err := c.xpath([]string{name})
	if err != nil {
		return err
	}
	path = append(path, "devices")

	_, err = c.ns.Client.Set(path, util.Entry{Value: d}, nil, nil)
	return err
}

/*
EditDevice performs an EDIT to add specific device to template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Panorama) EditDevice(st interface{}, d string) error {
	var name string

	switch v := st.(type) {
	case string:
		name = v
	case Entry:
		name = v.Name
	default:
		return fmt.Errorf("Unknown type sent to edit device: %s", v)
	}

	c.ns.Client.LogAction("(edit) device in template stack: %s", name)

	path, err := c.xpath([]string{name})
	if err != nil {
		return err
	}
	path = append(path, "devices", util.AsEntryXpath([]string{d}))

	_, err = c.ns.Client.Edit(path, util.Entry{Value: d}, nil, nil)
	return err
}

/*
DeleteDevice performs a DELETE to remove specific device d from template stack st.

The template stack can be either a string or an Entry object.
*/
func (c *Panorama) DeleteDevice(st interface{}, d string) error {
	var name string

	switch v := st.(type) {
	case string:
		name = v
	case Entry:
		name = v.Name
	default:
		return fmt.Errorf("Unknown type sent to delete device: %s", v)
	}

	c.ns.Client.LogAction("(delete) device from template stack: %s", name)

	path, err := c.xpath([]string{name})
	if err != nil {
		return err
	}
	path = append(path, "devices", util.AsEntryXpath([]string{d}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList() ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList() ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(e ...Entry) error {
	return c.ns.Set(c.pather(), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(e Entry) error {
	return c.ns.Edit(c.pather(), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(), names, nErr)
}

func (c *Panorama) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(v)
	}
}

func (c *Panorama) xpath(vals []string) ([]string, error) {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"template-stack",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
