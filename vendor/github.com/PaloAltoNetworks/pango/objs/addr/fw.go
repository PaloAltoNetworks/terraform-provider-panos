package addr

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Objects.Address namespace.
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
	return c.ns.Set(c.pather(vsys), specifier(e...))
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

/*
ConfigureGroup configures the given address objects on PAN-OS.

Due to caching in the GUI, objects configured via a bulk SET will not show
up in the dropdowns in the GUI (aka - Security Rules).  Restarting the
management plane will force the cache to refresh.
*/
func (c *Firewall) ConfigureGroup(vsys string, objs []Entry, prevNames []string) error {
	var err error
	setList := make([]Entry, 0, len(objs))
	editList := make([]Entry, 0, len(objs))

	curObjs, err := c.GetAll(vsys)
	if err != nil {
		return err
	}

	// Determine which can be set and which must be edited.
	for _, x := range objs {
		var found bool
		for _, live := range curObjs {
			if x.Name == live.Name {
				found = true
				if !ObjectsMatch(x, live) {
					editList = append(editList, x)
				}
				break
			}
		}
		if !found {
			setList = append(setList, x)
		}
	}

	// Set all objects.
	if len(setList) > 0 {
		if err = c.Set(vsys, setList...); err != nil {
			return err
		}
	}

	// Edit each object one by one.
	for _, x := range editList {
		if err = c.Edit(vsys, x); err != nil {
			return err
		}
	}

	// Delete rules removed from the group.
	if len(prevNames) != 0 {
		rmList := make([]interface{}, 0, len(prevNames))
		for _, name := range prevNames {
			var found bool
			for _, x := range objs {
				if x.Name == name {
					found = true
					break
				}
			}
			if !found {
				rmList = append(rmList, name)
			}
		}

		if err = c.Delete(vsys, rmList...); err != nil {
			return err
		}
	}

	return nil
}

func (c *Firewall) pather(vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vsys, v)
	}
}

func (c *Firewall) xpath(vsys string, vals []string) ([]string, error) {
	ans := make([]string, 0, 7)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"address",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
