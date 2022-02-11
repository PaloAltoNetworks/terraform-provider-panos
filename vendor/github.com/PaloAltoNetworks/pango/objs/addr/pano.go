package addr

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Objects.Address namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(dg string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(dg), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(dg string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(dg), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(dg, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(dg), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(dg, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(dg), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(dg string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(dg), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(dg string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(dg), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(dg string, e ...Entry) error {
	return c.ns.Set(c.pather(dg), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(dg string, e Entry) error {
	return c.ns.Edit(c.pather(dg), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(dg string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(dg), names, nErr)
}

/*
ConfigureGroup configures the given address objects on PAN-OS.

Due to caching in the GUI, objects configured via a bulk SET will not show
up in the dropdowns in the GUI (aka - Security Rules).  Restarting the
management plane will force the cache to refresh.
*/
func (c *Panorama) ConfigureGroup(dg string, objs []Entry, prevNames []string) error {
	var err error
	setList := make([]Entry, 0, len(objs))
	editList := make([]Entry, 0, len(objs))

	curObjs, err := c.GetAll(dg)
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
		if err = c.Set(dg, setList...); err != nil {
			return err
		}
	}

	// Edit each object one by one.
	for _, x := range editList {
		if err = c.Edit(dg, x); err != nil {
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

		if err = c.Delete(dg, rmList...); err != nil {
			return err
		}
	}

	return nil
}

func (c *Panorama) pather(dg string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(dg, v)
	}
}

func (c *Panorama) xpath(dg string, vals []string) ([]string, error) {
	ans := make([]string, 0, 7)
	ans = append(ans, util.DeviceGroupXpathPrefix(dg)...)
	ans = append(ans,
		"address",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
