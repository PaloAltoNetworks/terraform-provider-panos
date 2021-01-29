package filetype

import (
	"regexp"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Predefined.TdbFileType namespace.
type Panorama struct {
	ns *namespace.Standard
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

// GetMatches performs GET to retrieve a list of objects whose full
// name matches the given regex.
func (c *Panorama) GetMatches(expr string) ([]Entry, error) {
	var err error
	var re *regexp.Regexp
	if expr != "" {
		re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	}
	ans := c.container()
	err = c.ns.Objects(util.Get, c.pather(), ans)
	return finder(ans, re, err)
}

// ShowMatches performs SHOW to retrieve a list of objects whose full
// name matches the given regex.
func (c *Panorama) ShowMatches(expr string) ([]Entry, error) {
	var err error
	var re *regexp.Regexp
	if expr != "" {
		re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	}
	ans := c.container()
	err = c.ns.Objects(util.Show, c.pather(), ans)
	return finder(ans, re, err)
}

// Making this private so we can still do unit tests.
func (c *Panorama) set(e ...Entry) error {
	return c.ns.Set(c.pather(), specifier(e...))
}

func (c *Panorama) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(v)
	}
}

func (c *Panorama) xpath(vals []string) ([]string, error) {
	return []string{
		"config",
		"predefined",
		"tdb",
		"file-type",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
