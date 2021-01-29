package nat

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Policies.Nat namespace.
//
// The "dg" param in these functions is the device group.
//
// The "base" param in these functions should be one of the rulebase
// constants in the "util" package.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(dg, base string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(dg, base), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(dg, base string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(dg, base), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(dg, base, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(dg, base), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(dg, base, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(dg, base), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve information for all objects.
func (c *Panorama) GetAll(dg, base string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(dg, base), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(dg, base string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(dg, base), ans)
	return all(ans, err)
}

// Set performs SET to create / update one or more objects.
func (c *Panorama) Set(dg, base string, e ...Entry) error {
	err := c.ns.Set(c.pather(dg, base), specifier(e...))

	// On error: find the rule that's causing the error if multiple rules
	// were given.
	if err != nil && strings.Contains(err.Error(), "rules is invalid") {
		for i := 0; i < len(e); i++ {
			if e2 := c.Set(dg, base, e[i]); e2 != nil {
				return fmt.Errorf("Error with rule %d: %s", i+1, e2)
			} else {
				_ = c.Delete(dg, base, e[i])
			}
		}

		// Couldn't find it, just return the original error.
		return err
	}

	return err
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(dg, base string, e Entry) error {
	return c.ns.Edit(c.pather(dg, base), e)
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *Panorama) Delete(dg, base string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(dg, base), names, nErr)
}

// MoveGroup moves a logical group of NAT rules
// somewhere in relation to another rule.
//
// The `movement` param should be one of the Move constants in the util
// package.
//
// The `rule` param is the other rule the `movement` param is referencing.  If
// this is an empty string, then the first policy in the group isn't moved
// anywhere, but all other policies will still be moved to be grouped with the
// first one.
func (c *Panorama) MoveGroup(dg, base string, movement int, rule string, e ...Entry) error {
	lister := func() ([]string, error) {
		return c.GetList(dg, base)
	}

	ei := make([]interface{}, 0, len(e))
	for i := range e {
		ei = append(ei, e[i])
	}
	names, _ := toNames(ei)

	return c.ns.MoveGroup(c.pather(dg, base), lister, movement, rule, names)
}

func (c *Panorama) pather(dg, base string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(dg, base, v)
	}
}

func (c *Panorama) xpath(dg, base string, vals []string) ([]string, error) {
	if err := util.ValidateRulebase(base); err != nil {
		return nil, err
	}

	ans := make([]string, 0, 9)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		base,
		"nat",
		"rules",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
