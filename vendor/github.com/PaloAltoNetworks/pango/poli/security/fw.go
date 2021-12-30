package security

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/audit"
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Policies.PolicyBasedForwarding namespace.
type Firewall struct {
	ns *namespace.Policy
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(vsys), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(vsys), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vsys), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vsys), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(vsys), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(vsys), ans)
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
func (c *Firewall) ConfigureRules(vsys string, rules []Entry, auditComments map[string]string, isPolicy bool, move int, oRule string, prevNames []string) error {
	var err error
	setRules := make([]Entry, 0, len(rules))
	editRules := make([]Entry, 0, len(rules))

	curRules, err := c.GetAll(vsys)
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
		if err = c.Set(vsys, setRules...); err != nil {
			return err
		}
		// Configure audit comments for each set rule.
		for _, x := range setRules {
			if comment := auditComments[x.Name]; comment != "" {
				if err = c.SetAuditComment(vsys, x.Name, comment); err != nil {
					return err
				}
			}
		}
	}

	// Edit each rule one by one.
	for _, x := range editRules {
		if err = c.Edit(vsys, x); err != nil {
			return err
		}
		// Configure the audit comment for each edited rule.
		if comment := auditComments[x.Name]; comment != "" {
			if err = c.SetAuditComment(vsys, x.Name, comment); err != nil {
				return err
			}
		}
	}

	// Move the group into place.
	if err = c.MoveGroup(vsys, move, oRule, rules...); err != nil {
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
			_ = c.Delete(vsys, rmList...)
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
			if err = c.Delete(vsys, delRules...); err != nil {
				return nil
			}
		}
	}

	return nil
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Firewall) FromPanosConfig(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(vsys), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Firewall) AllFromPanosConfig(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(vsys), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vsys string, e ...Entry) error {
	err := c.ns.Set(c.pather(vsys), specifier(e...))

	// On error: find the rule that's causing the error if multiple rules
	// were given.
	if err != nil && strings.Contains(err.Error(), "rules is invalid") {
		for i := 0; i < len(e); i++ {
			if e2 := c.Set(vsys, e[i]); e2 != nil {
				return fmt.Errorf("Error with rule %d: %s", i+1, e2)
			} else {
				_ = c.Delete(vsys, e[i])
			}
		}

		// Couldn't find it, just return the original error.
		return err
	}

	return err
}

// VerifiableSet behaves like Set(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.
//
// NOTE:  Custom XML unmarshaling is now implemented, making this function
// unnecessary.
//
// This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *Firewall) VerifiableSet(vsys string, e ...Entry) error {
	c.ns.Client.LogAction("(set) performing verifiable set")
	again := make([]Entry, 0, len(e))

	for i := range e {
		if e[i].LogEnd {
			again = append(again, e[i])
			e[i].LogEnd = false
		}
	}

	if err := c.Set(vsys, e...); err != nil {
		return err
	}

	if len(again) == 0 {
		return nil
	}

	return c.Set(vsys, again...)
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vsys string, e Entry) error {
	return c.ns.Edit(c.pather(vsys), e)
}

// VerifiableEdit behaves like Edit(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.
//
// NOTE:  Custom XML unmarshaling is now implemented, making this function
// unnecessary.
//
// This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *Firewall) VerifiableEdit(vsys string, e ...Entry) error {
	var err error

	c.ns.Client.LogAction("(edit) performing verifiable edit")
	again := make([]Entry, 0, len(e))

	for i := range e {
		if e[i].LogEnd {
			again = append(again, e[i])
			e[i].LogEnd = false
		}
		if err = c.Edit(vsys, e[i]); err != nil {
			return err
		}
	}

	if len(again) == 0 {
		return nil
	}

	// It's ok to do a SET following an EDIT because we are guaranteed
	// to not have stray or conflicting config, so use SET since it
	// supports bulk operations.
	return c.Set(vsys, again...)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vsys), names, nErr)
}

// DeleteAll removes all security policies from the specified vsys.
func (c *Firewall) DeleteAll(vsys string) error {
	c.ns.Client.LogAction("(delete) all security policies")
	list, err := c.GetList(vsys)
	if err != nil || len(list) == 0 {
		return err
	}
	li := make([]interface{}, len(list))
	for i := range list {
		li[i] = list[i]
	}
	return c.Delete(vsys, li...)
}

// MoveGroup moves a logical group of security rules somewhere in relation
// to another security policy.
//
// The `movement` param should be one of the Move constants in the util
// package.
//
// The `rule` param is the other rule the `movement` param is referencing.  If
// this is an empty string, then the first policy in the group isn't moved
// anywhere, but all other policies will still be moved to be grouped with the
// first one.
func (c *Firewall) MoveGroup(vsys string, movement int, rule string, e ...Entry) error {
	lister := func() ([]string, error) {
		return c.GetList(vsys)
	}

	ei := make([]interface{}, 0, len(e))
	for i := range e {
		ei = append(ei, e[i])
	}
	names, _ := toNames(ei)

	return c.ns.MoveGroup(c.pather(vsys), lister, movement, rule, names)
}

// HitCount gets the rule hit count for the given rules.
//
// If the rules param is nil, then the hit count for all rules is returned.
func (c *Firewall) HitCount(vsys string, rules []string) ([]util.HitCount, error) {
	return c.ns.HitCount("security", vsys, rules)
}

// SetAuditComment sets the audit comment for the given rule.
func (c *Firewall) SetAuditComment(vsys, rule, comment string) error {
	return c.ns.SetAuditComment(c.pather(vsys), rule, comment)
}

// CurrentAuditComment returns the current audit comment.
func (c *Firewall) CurrentAuditComment(vsys, rule string) (string, error) {
	return c.ns.CurrentAuditComment(c.pather(vsys), rule)
}

// AuditCommentHistory returns a chunk of historical audit comment logs.
func (c *Firewall) AuditCommentHistory(vsys, rule, direction string, nlogs, skip int) ([]audit.Comment, error) {
	return c.ns.AuditCommentHistory(c.pather(vsys), rule, direction, nlogs, skip)
}

func (c *Firewall) pather(vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vsys, v)
	}
}

func (c *Firewall) xpath(vsys string, vals []string) ([]string, error) {
	if vsys == "" {
		vsys = "vsys1"
	} else if vsys == "shared" {
		return nil, fmt.Errorf("vsys must be specified")
	}

	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"vsys",
		util.AsEntryXpath([]string{vsys}),
		"rulebase",
		"security",
		"rules",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
