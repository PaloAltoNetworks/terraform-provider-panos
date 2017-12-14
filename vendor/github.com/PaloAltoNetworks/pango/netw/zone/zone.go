// Package zone is the client.Network.Zone namespace.
//
// Normalized object:  Entry
package zone

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a zone.
type Entry struct {
    Name string
    Mode string
    Interfaces []string
    ZoneProfile string
    LogSetting string
    EnableUserId bool
    IncludeAcl []string
    ExcludeAcl []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Mode = s.Mode
    o.Interfaces = s.Interfaces
    o.ZoneProfile = s.ZoneProfile
    o.LogSetting = s.LogSetting
    o.EnableUserId = s.EnableUserId
    o.IncludeAcl = s.IncludeAcl
    o.ExcludeAcl = s.ExcludeAcl
}

// Zone is a namespace struct, included as part of pango.Client.
type Zone struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *Zone) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of zones.
func (c *Zone) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of zones")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of zones.
func (c *Zone) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of zones")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given zone.
func (c *Zone) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) zone %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given zone.
func (c *Zone) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) zone %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more zones.
func (c *Zone) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given zone configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "zone"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) zones: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the zones.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to creates / updates a zone.
func (c *Zone) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) zone %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the zones.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given zone(s) from the firewall.
//
// Zones can be either a string or an Entry object.
func (c *Zone) Delete(vsys string, e ...interface{}) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    names := make([]string, len(e))
    for i := range e {
        switch v := e[i].(type) {
        case string:
            names[i] = v
        case Entry:
            names[i] = v.Name
        default:
            return fmt.Errorf("Unsupported type to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) zone(s): %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Zone struct **/

func (c *Zone) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Zone) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Zone) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "zone",
        util.AsEntryXpath(vals),
    }
}

/** Structs / functions for this namespace. **/

type normalizer interface {
    Normalize() Entry
}

type container_v1 struct {
    Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        ZoneProfile: o.Answer.Profile,
        LogSetting: o.Answer.LogSetting,
        EnableUserId: util.AsBool(o.Answer.EnableUserId),
    }
    if o.Answer.L3 != nil {
        ans.Mode = "layer3"
        ans.Interfaces = o.Answer.L3.Interfaces
    } else if o.Answer.L2 != nil {
        ans.Mode = "layer2"
        ans.Interfaces = o.Answer.L2.Interfaces
    } else if o.Answer.VWire != nil {
        ans.Mode = "virtual-wire"
        ans.Interfaces = o.Answer.VWire.Interfaces
    } else if o.Answer.Tap != nil {
        ans.Mode = "tap"
        ans.Interfaces = o.Answer.Tap.Interfaces
    } else if o.Answer.External != nil {
        ans.Mode = "external"
        ans.Interfaces = o.Answer.External.Interfaces
    }
    if o.Answer.IncludeAcl != nil {
        ans.IncludeAcl = o.Answer.IncludeAcl.Acls
    }
    if o.Answer.ExcludeAcl != nil {
        ans.ExcludeAcl = o.Answer.ExcludeAcl.Acls
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    L3 *zoneInterfaceList `xml:"network>layer3"`
    L2 *zoneInterfaceList `xml:"network>layer2"`
    VWire *zoneInterfaceList `xml:"network>virtual-wire"`
    Tap *zoneInterfaceList `xml:"network>tap"`
    External *zoneInterfaceList `xml:"network>external"`
    Profile string `xml:"network>zone-protection-profile,omitempty"`
    LogSetting string `xml:"network>log-setting,omitempty"`
    EnableUserId string `xml:"enable-user-identification"`
    IncludeAcl *aclList `xml:"user-acl>include-list"`
    ExcludeAcl *aclList `xml:"user-acl>exclude-list"`
}

type zoneInterfaceList struct {
    Interfaces []string `xml:"member"`
}

type aclList struct {
    Acls []string `xml:"member"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Profile: e.ZoneProfile,
        LogSetting: e.LogSetting,
        EnableUserId: util.YesNo(e.EnableUserId),
    }
    il := &zoneInterfaceList{e.Interfaces}
    switch e.Mode {
    case "layer2":
        ans.L2 = il
    case "layer3":
        ans.L3 = il
    case "virtual-wire":
        ans.VWire = il
    case "tap":
        ans.Tap = il
    case "external":
        ans.External = il
    }
    if len(e.IncludeAcl) > 0 {
        inu := &aclList{e.IncludeAcl}
        ans.IncludeAcl = inu
    }
    if len(e.ExcludeAcl) > 0 {
        exu := &aclList{e.ExcludeAcl}
        ans.ExcludeAcl = exu
    }

    return ans
}
