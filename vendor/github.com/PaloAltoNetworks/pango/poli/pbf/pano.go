package pbf

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/audit"
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Policies.PolicyBasedForwarding namespace.
//
// The "dg" param in these functions is the device group.
//
// The "base" param in these functions should be one of the rulebase
// constants in the "util" package.
type Panorama struct {
	ns *namespace.Policy
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

/*
ConfigureRules configures the given rules on PAN-OS.

It does a mass SET if it can, but will EDIT any rules that are present but
differ from what is given.

Audit comments are applied only for rules which are either SET or EDIT'ed.

If isPolicy is true, then any rules not explicitly present in the rules param will
be deleted.

Params move and oRule are for moving the group into place after configuration.

Any rule name that appears in prevRules but not in the rules param will be deleted.
*/
func (c *Panorama) ConfigureRules(dg, base string, rules []Entry, auditComments map[string]string, isPolicy bool, move int, oRule string, prevNames []string) error {
	var err error
	setRules := make([]Entry, 0, len(rules))
	editRules := make([]Entry, 0, len(rules))

	curRules, err := c.GetAll(dg, base)
	if err != nil {
		return err
	}

	// Determine which can be set and which can must be edited.
	for _, x := range rules {
		var found bool
		for _, live := range curRules {
			if x.Name == live.Name {
				found = true
				if !RulesMatch(x, live) {
					editRules = append(editRules, x)
				}
				break
			}
		}
		if !found {
			setRules = append(setRules, x)
		}
	}

	// Set all rules.
	if len(setRules) > 0 {
		if err = c.Set(dg, base, setRules...); err != nil {
			return err
		}
		// Configure audit comments for each set rule.
		for _, x := range setRules {
			if comment := auditComments[x.Name]; comment != "" {
				if err = c.SetAuditComment(dg, base, x.Name, comment); err != nil {
					return err
				}
			}
		}
	}

	// Edit each rule one by one.
	for _, x := range editRules {
		if err = c.Edit(dg, base, x); err != nil {
			return err
		}
		// Configure the audit comment for each edited rule.
		if comment := auditComments[x.Name]; comment != "" {
			if err = c.SetAuditComment(dg, base, x.Name, comment); err != nil {
				return err
			}
		}
	}

	// Move the group into place.
	if err = c.MoveGroup(dg, base, move, oRule, rules...); err != nil {
		return err
	}

	// Delete rules removed from the group.
	if len(prevNames) != 0 {
		rmList := make([]interface{}, 0, len(prevNames))
		for _, name := range prevNames {
			var found bool
			for _, x := range rules {
				if x.Name == name {
					found = true
					break
				}
			}
			if !found {
				rmList = append(rmList, name)
			}
		}

		if len(rmList) != 0 {
			_ = c.Delete(dg, base, rmList...)
		}
	}

	// Optional: If this is a policy, delete everything else.
	if isPolicy {
		delRules := make([]interface{}, 0, len(curRules))
		for _, cur := range curRules {
			var found bool
			for _, x := range rules {
				if x.Name == cur.Name {
					found = true
					break
				}
			}

			if !found {
				delRules = append(delRules, cur.Name)
			}
		}

		if len(delRules) != 0 {
			if err = c.Delete(dg, base, delRules...); err != nil {
				return nil
			}
		}
	}

	return nil
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Panorama) FromPanosConfig(dg, base, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(dg, base), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Panorama) AllFromPanosConfig(dg, base string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(dg, base), ans)
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

// MoveGroup moves a logical group of policy based forwarding rules
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

// SetAuditComment sets the audit comment for the given rule.
func (c *Panorama) SetAuditComment(dg, base, rule, comment string) error {
	return c.ns.SetAuditComment(c.pather(dg, base), rule, comment)
}

// CurrentAuditComment returns the current audit comment.
func (c *Panorama) CurrentAuditComment(dg, base, rule string) (string, error) {
	return c.ns.CurrentAuditComment(c.pather(dg, base), rule)
}

// AuditCommentHistory returns a chunk of historical audit comment logs.
func (c *Panorama) AuditCommentHistory(dg, base, rule, direction string, nlogs, skip int) ([]audit.Comment, error) {
	return c.ns.AuditCommentHistory(c.pather(dg, base), rule, direction, nlogs, skip)
}

func (c *Panorama) pather(dg, base string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(dg, base, v)
	}
}

func (c *Panorama) xpath(dg, base string, vals []string) ([]string, error) {
	if err := util.ValidateRulebase(dg, base); err != nil {
		return nil, err
	}

	ans := make([]string, 0, 9)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		base,
		"pbf",
		"rules",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
