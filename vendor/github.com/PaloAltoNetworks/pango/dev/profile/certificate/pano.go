package certificate

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Device.CertificateProfile namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(shared bool, tmpl, ts, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(shared, tmpl, ts, vsys), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(shared bool, tmpl, ts, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(shared, tmpl, ts, vsys), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(shared bool, tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(shared, tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(shared bool, tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(shared, tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(shared bool, tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(shared, tmpl, ts, vsys), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(shared bool, tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(shared, tmpl, ts, vsys), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(shared bool, tmpl, ts, vsys string, e ...Entry) error {
	return c.ns.Set(c.pather(shared, tmpl, ts, vsys), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(shared bool, tmpl, ts, vsys string, e Entry) error {
	return c.ns.Edit(c.pather(shared, tmpl, ts, vsys), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(shared bool, tmpl, ts, vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(shared, tmpl, ts, vsys), names, nErr)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Panorama) FromPanosConfig(shared bool, tmpl, ts, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(shared, tmpl, ts, vsys), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Panorama) AllFromPanosConfig(shared bool, tmpl, ts, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(shared, tmpl, ts, vsys), ans)
	return all(ans, err)
}

func (c *Panorama) pather(shared bool, tmpl, ts, vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(shared, tmpl, ts, vsys, v)
	}
}

func (c *Panorama) xpath(shared bool, tmpl, ts, vsys string, vals []string) ([]string, error) {
	// Sanity check input.
	if tmpl == "" && ts == "" && vsys != "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}

	var ans []string
	if tmpl != "" || ts != "" {
		ans = make([]string, 0, 12)
		ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
		ans = append(ans, util.VsysXpathPrefix(vsys)...)
	} else {
		ans = make([]string, 0, 4)
		if shared {
			ans = append(ans,
				"config",
				"shared",
			)
		} else {
			ans = append(ans,
				"config",
				"panorama",
			)
		}
	}

	ans = append(ans,
		"certificate-profile",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
