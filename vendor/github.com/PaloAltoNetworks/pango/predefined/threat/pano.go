package threat

import (
	"fmt"
	"regexp"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Predefined.Threat namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tt string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tt), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tt string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tt), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(tt, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tt), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(tt, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tt), name, ans)
	return first(ans, err)
}

// GetThreats performs a GET to retrieve a list of objects whose
// threat name matches the given regex.
func (c *Panorama) GetThreats(tt, expr string) ([]Entry, error) {
	var err error
	var re *regexp.Regexp
	if expr != "" {
		re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	}
	ans := c.container()
	err = c.ns.Objects(util.Get, c.pather(tt), ans)
	return finder(ans, re, err)
}

// ShowThreats performs a SHOW to retrieve a list of objects whose
// threat name matches the given regex.
func (c *Panorama) ShowThreats(tt, expr string) ([]Entry, error) {
	var err error
	var re *regexp.Regexp
	if expr != "" {
		re, err = regexp.Compile(expr)
		if err != nil {
			return nil, err
		}
	}
	ans := c.container()
	err = c.ns.Objects(util.Show, c.pather(tt), ans)
	return finder(ans, re, err)
}

// Making this private so we can still do unit tests.
func (c *Panorama) set(vsys string, e ...Entry) error {
	return c.ns.Set(c.pather(vsys), specifier(e...))
}

func (c *Panorama) pather(tt string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tt, v)
	}
}

func (c *Panorama) xpath(tt string, vals []string) ([]string, error) {
	switch tt {
	case Vulnerability, PhoneHome:
	default:
		return nil, fmt.Errorf("invalid threat type: %s", tt)
	}

	return []string{
		"config",
		"predefined",
		"threats",
		tt,
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
