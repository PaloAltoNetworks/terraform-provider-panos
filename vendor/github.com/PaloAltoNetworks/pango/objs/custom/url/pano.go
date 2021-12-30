package url

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Objects.CustomUrlCategory namespace.
type Panorama struct {
	ns *namespace.Standard
}

// SetSite performs a SET to add a site to the custom URL category.
func (c *Panorama) SetSite(dg, name, site string) error {
	c.ns.Client.LogAction("(set) site for %s: %s", name, site)

	path, err := c.xpath(dg, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "list")

	_, err = c.ns.Client.Set(path, util.Member{Value: site}, nil, nil)
	return err
}

// DeleteSite performs a DELETE to remove a site from the custom URL category.
func (c *Panorama) DeleteSite(dg, name, site string) error {
	c.ns.Client.LogAction("(delete) site from %s: %s", name, site)

	path, err := c.xpath(dg, []string{name})
	if err != nil {
		return err
	}
	path = append(path, "list", util.AsMemberXpath([]string{site}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(dg string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(dg), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(dg string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(dg), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(dg, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(dg), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(dg, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(dg), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(dg string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(dg), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(dg string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(dg), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(dg string, e ...Entry) error {
	return c.ns.Set(c.pather(dg), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(dg string, e Entry) error {
	return c.ns.Edit(c.pather(dg), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(dg string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(dg), names, nErr)
}

func (c *Panorama) pather(dg string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(dg, v)
	}
}

func (c *Panorama) xpath(dg string, vals []string) ([]string, error) {
	ans := make([]string, 0, 8)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		"profiles",
		"custom-url-category",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
