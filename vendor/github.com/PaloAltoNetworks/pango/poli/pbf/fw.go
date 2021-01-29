package pbf

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Firewall is the client.Policies.PolicyBasedForwarding namespace.
type Firewall struct {
	ns *namespace.Standard
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

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vsys string, e Entry) error {
	return c.ns.Edit(c.pather(vsys), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vsys), names, nErr)
}

// MoveGroup moves a logical group of policy based forwarding rules somewhere
// in relation to another rule.
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
	if !c.ns.Client.Versioning().Gte(version.Number{8, 1, 0, ""}) {
		return nil, fmt.Errorf("rule hit count requires PAN-OS 8.1+")
	}

	req := util.NewHitCountRequest("pbf", vsys, rules)
	resp := util.HitCountResponse{}
	if _, err := c.ns.Client.Op(req, "", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (c *Firewall) pather(vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vsys, v)
	}
}

func (c *Firewall) xpath(vsys string, vals []string) ([]string, error) {
	if vsys == "shared" {
		return nil, fmt.Errorf("vsys cannot be 'shared'")
	} else if vsys == "" {
		return nil, fmt.Errorf("vsys must be specified")
	}

	ans := make([]string, 0, 9)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"rulebase",
		"pbf",
		"rules",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
