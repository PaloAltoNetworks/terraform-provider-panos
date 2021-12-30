package account

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Panorama.GcpAccount namespace.
type Panorama struct {
	ns *namespace.Plugin
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList() ([]string, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	return c.ns.Listing(util.Get, c.pather(), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList() ([]string, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	return c.ns.Listing(util.Show, c.pather(), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(name string) (Entry, error) {
	ans, err := c.container()
	if err != nil {
		return Entry{}, err
	}
	err = c.ns.Object(util.Get, c.pather(), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(name string) (Entry, error) {
	ans, err := c.container()
	if err != nil {
		return Entry{}, err
	}
	err = c.ns.Object(util.Show, c.pather(), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll() ([]Entry, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	err = c.ns.Objects(util.Get, c.pather(), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll() ([]Entry, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	err = c.ns.Objects(util.Show, c.pather(), ans)
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
		"plugins",
		"gcp",
		"setup",
		"service-act-credentials",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() (normalizer, error) {
	return container(c.ns.Client.Plugins())
}
