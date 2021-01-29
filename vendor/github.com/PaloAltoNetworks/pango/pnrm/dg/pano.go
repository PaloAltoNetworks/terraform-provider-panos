package dg

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Panorama.DeviceGroup namespace.
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

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll() ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(e ...Entry) error {
	return c.ns.Set(c.pather(), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(e Entry) error {
	return c.ns.Edit(c.pather(), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(), names, nErr)
}

/*
SetDeviceVsys performs a SET to add specific vsys from a device to device
group g.

If you want all vsys to be included, or the device is a virtual firewall, then
leave the vsys list empty.

The device group can be either a string or an Entry object.
*/
func (c *Panorama) SetDeviceVsys(g interface{}, d string, vsys []string) error {
	names, err := toNames([]interface{}{g})
	if err != nil {
		return err
	}

	c.ns.Client.LogAction("(set) device vsys in device group: %s", names[0])

	path, err := c.xpath(names)
	if err != nil {
		return err
	}
	path = append(path, "devices")
	m := util.MapToVsysEnt(map[string][]string{d: vsys})

	_, err = c.ns.Client.Set(path, m.Entries[0], nil, nil)
	return err
}

/*
EditDeviceVsys performs an EDIT to add specific vsys from a device to device
group g.

If you want all vsys to be included, or the device is a virtual firewall, then
leave the vsys list empty.

The device group can be either a string or an Entry object.
*/
func (c *Panorama) EditDeviceVsys(g interface{}, d string, vsys []string) error {
	names, err := toNames([]interface{}{g})
	if err != nil {
		return err
	}

	c.ns.Client.LogAction("(set) device vsys in device group: %s", names[0])

	path, err := c.xpath(names)
	if err != nil {
		return err
	}
	path = append(path, "devices", util.AsEntryXpath([]string{d}))
	m := util.MapToVsysEnt(map[string][]string{d: vsys})

	_, err = c.ns.Client.Edit(path, m.Entries[0], nil, nil)
	return err
}

/*
DeleteDeviceVsys performs a DELETE to remove specific vsys from device d from
device group g.

If you want all vsys to be removed, or the device is a virtual firewall, then
leave the vsys list empty.

The device group can be either a string or an Entry object.
*/
func (c *Panorama) DeleteDeviceVsys(g interface{}, d string, vsys []string) error {
	names, err := toNames([]interface{}{g})
	if err != nil {
		return err
	}

	c.ns.Client.LogAction("(delete) device vsys from device group: %s", names[0])

	path := make([]string, 0, 9)
	p, err := c.xpath(names)
	if err != nil {
		return err
	}
	path = append(path, p...)
	path = append(path, "devices", util.AsEntryXpath([]string{d}))
	if len(vsys) > 0 {
		path = append(path, "vsys", util.AsEntryXpath(vsys))
	}

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// GetParents returns a map where the keys are the device group's name and
// the value is the parent for that device group.
//
// An empty parent value means that the parent is the "shared" device group.
func (c *Panorama) GetParents() (map[string]string, error) {
	req := dgpReq{}
	ans := dgpResp{}

	c.ns.Client.LogOp("(op) retrieving device group parents")
	if _, err := c.ns.Client.Op(req, "", nil, &ans); err != nil {
		return nil, err
	}

	return ans.results(), nil
}

// AssignParent sets a device group's parent to `parent`.
//
// An empty string for the parent will move the device group to the
// top level (shared).
//
// This operation results in a job being submitted to the backend, which this
// function will block until the move is completed.
func (c *Panorama) AssignParent(child, parent string) error {

	req := apReq{
		Info: apInfo{
			Child:  child,
			Parent: parent,
		},
	}
	ans := util.JobResponse{}

	c.ns.Client.LogOp("(op) assigning group %q new parent: %s", child, parent)
	if _, err := c.ns.Client.Op(req, "", nil, &ans); err != nil {
		return err
	}

	return c.ns.Client.WaitForJob(ans.Id, 0, nil)
}

func (c *Panorama) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(v)
	}
}

func (c *Panorama) xpath(vals []string) ([]string, error) {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"device-group",
		util.AsEntryXpath(vals),
	}, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}

// Device group hierarchy structs.
type dgpReq struct {
	XMLName xml.Name `xml:"show"`
	Cmd     string   `xml:"dg-hierarchy"`
}

type dgpResp struct {
	Result *dgHierarchy `xml:"result>dg-hierarchy"`
}

func (o *dgpResp) results() map[string]string {
	ans := make(map[string]string)

	if o.Result != nil {
		for _, v := range o.Result.Info {
			ans[v.Name] = ""
			v.results(ans)
		}
	}

	return ans
}

type dgHierarchy struct {
	Info []dghInfo `xml:"dg"`
}

type dghInfo struct {
	Name     string    `xml:"name,attr"`
	Children []dghInfo `xml:"dg"`
}

func (o *dghInfo) results(ans map[string]string) {
	for _, v := range o.Children {
		ans[v.Name] = o.Name
		v.results(ans)
	}
}

type apReq struct {
	XMLName xml.Name `xml:"request"`
	Info    apInfo   `xml:"move-dg>entry"`
}

type apInfo struct {
	Child  string `xml:"name,attr"`
	Parent string `xml:"new-parent-dg,omitempty"`
}
