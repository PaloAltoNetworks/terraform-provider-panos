package exp

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network.BgpExport namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(tmpl, ts, vr string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(tmpl, ts, vr), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(tmpl, ts, vr string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(tmpl, ts, vr), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(tmpl, ts, vr, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, vr), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(tmpl, ts, vr, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, vr), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(tmpl, ts, vr string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(tmpl, ts, vr), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(tmpl, ts, vr string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(tmpl, ts, vr), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vr string, e ...Entry) error {
	err := c.ns.Set(c.pather(tmpl, ts, vr), specifier(e...))

	// On error: find the rule that's causing the error if multiple rules
	// were given.
	if err != nil && strings.Contains(err.Error(), "rules is invalid") {
		for i := 0; i < len(e); i++ {
			if e2 := c.Set(tmpl, ts, vr, e[i]); e2 != nil {
				return fmt.Errorf("Error with rule %d: %s", i+1, e2)
			} else {
				_ = c.Delete(tmpl, ts, vr, e[i])
			}
		}

		// Couldn't find it, just return the original error.
		return err
	}

	return err
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vr string, e Entry) error {
	return c.ns.Edit(c.pather(tmpl, ts, vr), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(tmpl, ts, vr string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(tmpl, ts, vr), names, nErr)
}

// MoveGroup moves a logical group of bgp export rules somewhere in relation
// to another security policy.
//
// The `movement` param should be one of the Move constants in the util
// package.
//
// The `rule` param is the other rule the `movement` param is referencing.  If
// this is an empty string, then the first policy in the group isn't moved
// anywhere, but all other policies will still be moved to be grouped with the
// first one.
func (c *Panorama) MoveGroup(tmpl, ts, vr string, movement int, rule string, e ...Entry) error {
	lister := func() ([]string, error) {
		return c.GetList(tmpl, ts, vr)
	}

	ei := make([]interface{}, 0, len(e))
	for i := range e {
		ei = append(ei, e[i])
	}
	names, _ := toNames(ei)

	return c.ns.MoveGroup(c.pather(tmpl, ts, vr), lister, movement, rule, names)
}

func (c *Panorama) pather(tmpl, ts, vr string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vr, v)
	}
}

func (c *Panorama) xpath(tmpl, ts, vr string, vals []string) ([]string, error) {
	if tmpl == "" && ts == "" {
		return nil, fmt.Errorf("tmpl or ts must be specified")
	}
	if vr == "" {
		return nil, fmt.Errorf("vr must be specified")
	}

	ans := make([]string, 0, 17)
	ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"network",
		"virtual-router",
		util.AsEntryXpath([]string{vr}),
		"protocol",
		"bgp",
		"policy",
		"export",
		"rules",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
