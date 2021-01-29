package security

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Policies.Security namespace.
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

// VerifiableSet behaves like Set(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.  This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *Panorama) VerifiableSet(dg, base string, e ...Entry) error {
	c.ns.Client.LogAction("(set) performing verifiable set")
	again := make([]Entry, 0, len(e))

	for i := range e {
		if e[i].LogEnd {
			again = append(again, e[i])
			e[i].LogEnd = false
		}
	}

	if err := c.Set(dg, base, e...); err != nil {
		return err
	}

	if len(again) == 0 {
		return nil
	}

	return c.Set(dg, base, again...)
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(dg, base string, e Entry) error {
	return c.ns.Edit(c.pather(dg, base), e)
}

// VerifiableEdit behaves like Edit(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.  This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *Panorama) VerifiableEdit(dg, base string, e ...Entry) error {
	var err error

	c.ns.Client.LogAction("(edit) performing verifiable edit")
	again := make([]Entry, 0, len(e))

	for i := range e {
		if e[i].LogEnd {
			again = append(again, e[i])
			e[i].LogEnd = false
		}
		if err = c.Edit(dg, base, e[i]); err != nil {
			return err
		}
	}

	if len(again) == 0 {
		return nil
	}

	// It's ok to do a SET following an EDIT because we are guaranteed
	// to not have stray or conflicting config, so use SET since it
	// supports bulk operations.
	return c.Set(dg, base, again...)
}

// Delete removes the given objects.
//
// Objects can be a string or an Entry object.
func (c *Panorama) Delete(dg, base string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(dg, base), names, nErr)
}

// DeleteAll removes all security policies from the specified dg / rulebase.
func (c *Panorama) DeleteAll(dg, base string) error {
	c.ns.Client.LogAction("(delete) all security rules")
	list, err := c.GetList(dg, base)
	if err != nil || len(list) == 0 {
		return err
	}
	li := make([]interface{}, len(list))
	for i := range list {
		li[i] = list[i]
	}
	return c.Delete(dg, base, li...)
}

// MoveGroup moves a logical group of security rules
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
	/*
	   if base == "" {
	       base = util.PreRulebase
	   }
	*/
	if err := util.ValidateRulebase(base); err != nil {
		return nil, err
	}

	ans := make([]string, 0, 9)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		base,
		"security",
		"rules",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
