package email

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Device.EmailServerProfile namespace.
type Panorama struct {
	ns *namespace.Standard
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

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Panorama) FromPanosConfig(tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Panorama) AllFromPanosConfig(tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(tmpl, ts, vsys), ans)
	return all(ans, err)
}

func (c *Panorama) pather(tmpl, ts, vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vsys, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, vsys string, vals []string) ([]string, error) {
	var ans []string

	if tmpl != "" || ts != "" {
		if vsys == "" {
			vsys = "shared"
		}

		ans = make([]string, 0, 13)
		ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
		ans = append(ans, util.VsysXpathPrefix(vsys)...)
	} else {
		ans = make([]string, 0, 5)
		ans = append(ans, "config", "panorama")
	}

	ans = append(ans,
		"log-settings",
		"email",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
