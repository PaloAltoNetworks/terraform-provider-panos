package cluster

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Panorama.GkeCluster namespace.
type Panorama struct {
	ns *namespace.Plugin
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(group string) ([]string, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	return c.ns.Listing(util.Get, c.pather(group), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(group string) ([]string, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	return c.ns.Listing(util.Show, c.pather(group), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(group, name string) (Entry, error) {
	ans, err := c.container()
	if err != nil {
		return Entry{}, err
	}
	err = c.ns.Object(util.Get, c.pather(group), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(group, name string) (Entry, error) {
	ans, err := c.container()
	if err != nil {
		return Entry{}, err
	}
	err = c.ns.Object(util.Show, c.pather(group), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(group string) ([]Entry, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	err = c.ns.Objects(util.Get, c.pather(group), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(group string) ([]Entry, error) {
	ans, err := c.container()
	if err != nil {
		return nil, err
	}
	err = c.ns.Objects(util.Show, c.pather(group), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(group string, e ...Entry) error {
	return c.ns.Set(c.pather(group), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(group string, e Entry) error {
	return c.ns.Edit(c.pather(group), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(group string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(group), names, nErr)
}

func (c *Panorama) pather(group string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(group, v)
	}
}

func (c *Panorama) xpath(group string, vals []string) ([]string, error) {
	if group == "" {
		return nil, fmt.Errorf("group must be specified")
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"plugins",
		"gcp",
		"gke",
		util.AsEntryXpath([]string{group}),
		"gke-cluster",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() (normalizer, error) {
	return container(c.ns.Client.Plugins())
}
