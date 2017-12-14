// Package router is the client.Network.VirtualRouter namespace.
//
// Normalized object:  Entry
package router

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a virtual
// router.
type Entry struct {
    Name string
    Interfaces []string
    StaticDist int
    StaticIpv6Dist int
    OspfIntDist int
    OspfExtDist int
    Ospfv3IntDist int
    Ospfv3ExtDist int
    IbgpDist int
    EbgpDist int
    RipDist int

    raw map[string] string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * StaticDist: 10
//      * StaticIpv6Dist: 10
//      * OspfIntDist: 30
//      * OspfExtDist: 110
//      * Ospfv3IntDist: 30
//      * Ospfv3ExtDist: 110
//      * IbgpDist: 200
//      * EbgpDist: 20
//      * RipDist: 120
func (o *Entry) Defaults() {
    if o.StaticDist == 0 {
        o.StaticDist = 10
    }

    if o.StaticIpv6Dist == 0 {
        o.StaticIpv6Dist = 10
    }

    if o.OspfIntDist == 0 {
        o.OspfIntDist = 30
    }

    if o.OspfExtDist == 0 {
        o.OspfExtDist = 110
    }

    if o.Ospfv3IntDist == 0 {
        o.Ospfv3IntDist = 30
    }

    if o.Ospfv3ExtDist == 0 {
        o.Ospfv3ExtDist = 110
    }

    if o.IbgpDist == 0 {
        o.IbgpDist = 200
    }

    if o.EbgpDist == 0 {
        o.EbgpDist = 20
    }

    if o.RipDist == 0 {
        o.RipDist = 120
    }
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Interfaces = s.Interfaces
    o.StaticDist = s.StaticDist
    o.StaticIpv6Dist = s.StaticIpv6Dist
    o.OspfIntDist = s.OspfIntDist
    o.OspfExtDist = s.OspfExtDist
    o.Ospfv3IntDist = s.Ospfv3IntDist
    o.Ospfv3ExtDist = s.Ospfv3ExtDist
    o.IbgpDist = s.IbgpDist
    o.EbgpDist = s.EbgpDist
    o.RipDist = s.RipDist
}

// Router is the client.Network.VirtualRouter namespace.
type Router struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *Router) Initialize(con util.XapiClient) {
    c.con = con
}

// ShowList performs SHOW to retrieve a list of virtual routers.
func (c *Router) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of virtual routeres")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// GetList performs GET to retrieve a list of virtual routers.
func (c *Router) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of virtual routers")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given virtual router.
func (c *Router) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) virtual router %q", name)
    return c.details(c.con.Get, name)
}

// Show performs SHOW to retrieve information for the given virtual router.
func (c *Router) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) virtual router %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more virtual routers.
//
// Specify a non-empty vsys to import the virtual routers into the given vsys
// after creating, allowing the vsys to use them.
func (c *Router) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given router configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "virtual-router"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) virtual routers: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the virtual routers.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    if err != nil {
        return err
    }

    // Perform vsys import next.
    if vsys == "" {
        return nil
    }
    return c.con.ImportVirtualRouters(vsys, names)
}

// Edit performs EDIT to create / update a virtual router.
//
// Specify a non-empty vsys to import the virtual router into the given vsys
// after creating, allowing the vsys to use them.
func (c *Router) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) virtual router %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the virtual router.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    if err != nil {
        return err
    }

    // Perform vsys import next.
    if vsys == "" {
        return nil
    }
    return c.con.ImportVirtualRouters(vsys, []string{e.Name})
}

// Delete removes the given virtual routers from the firewall.
//
// Specify a non-empty vsys to have this function remove the virtual routers
// from the vsys prior to deleting them.
//
// Virtual routers can be a string or an Entry object.
func (c *Router) Delete(vsys string, e ...interface{}) error {
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
            return fmt.Errorf("Unknown type sent to delete: %s", v)
        }
    }
    c.con.LogAction("(delete) virtual routers: %v", names)

    // Unimport virtual routers from the given vsys.
    err = c.con.UnimportVirtualRouters(vsys, names)
    if err != nil {
        return err
    }

    // Remove virtual routers next.
    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

// CleanupDefault clears the `default` route configuration instead of deleting
// it outright.  This involves unimporting the route "default" from the given
// vsys, then performing an `EDIT` with an empty router.Entry object.
func (c *Router) CleanupDefault(vsys string) error {
    var err error

    c.con.LogAction("(action) cleaning up default route")

    // Unimport the default virtual router.
    if err = c.con.UnimportVirtualRouters(vsys, []string{"default"}); err != nil {
        return err
    }

    // Cleanup the interfaces the virtual router refers to.
    info := Entry{Name: "default"}
    if err = c.Edit("", info); err != nil {
        return err
    }

    return nil
}

/** Internal functions for the Router struct **/

func (c *Router) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Router) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Router) xpath(vals []string) []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "virtual-router",
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
        Interfaces: util.MemToStr(o.Answer.Interfaces),
        StaticDist: o.Answer.Dist.StaticDist,
        StaticIpv6Dist: o.Answer.Dist.StaticIpv6Dist,
        OspfIntDist: o.Answer.Dist.OspfIntDist,
        OspfExtDist: o.Answer.Dist.OspfExtDist,
        Ospfv3IntDist: o.Answer.Dist.Ospfv3IntDist,
        Ospfv3ExtDist: o.Answer.Dist.Ospfv3ExtDist,
        IbgpDist: o.Answer.Dist.IbgpDist,
        EbgpDist: o.Answer.Dist.EbgpDist,
        RipDist: o.Answer.Dist.RipDist,
    }
    ans.raw = make(map[string] string)
    if o.Answer.Ecmp != nil {
        ans.raw["ecmp"] = util.CleanRawXml(o.Answer.Ecmp.Text)
    }
    if o.Answer.Multicast != nil {
        ans.raw["multicast"] = util.CleanRawXml(o.Answer.Multicast.Text)
    }
    if o.Answer.Protocol != nil {
        ans.raw["protocol"] = util.CleanRawXml(o.Answer.Protocol.Text)
    }
    if o.Answer.Routing != nil {
        ans.raw["routing"] = util.CleanRawXml(o.Answer.Routing.Text)
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Interfaces *util.Member `xml:"interface"`
    Dist dist `xml:"admin-dists"`
    Ecmp *util.RawXml `xml:"ecmp"`
    Multicast *util.RawXml `xml:"multicast"`
    Protocol *util.RawXml `xml:"protocol"`
    Routing *util.RawXml `xml:"routing-table"`
}

type dist struct {
    StaticDist int `xml:"static,omitempty"`
    StaticIpv6Dist int `xml:"static-ipv6,omitempty"`
    OspfIntDist int `xml:"ospf-int,omitempty"`
    OspfExtDist int `xml:"ospf-ext,omitempty"`
    Ospfv3IntDist int `xml:"ospfv3-int,omitempty"`
    Ospfv3ExtDist int `xml:"ospfv3-ext,omitempty"`
    IbgpDist int `xml:"ibgp,omitempty"`
    EbgpDist int `xml:"ebgp,omitempty"`
    RipDist int `xml:"rip,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Interfaces: util.StrToMem(e.Interfaces),
        Dist: dist{
            StaticDist: e.StaticDist,
            StaticIpv6Dist: e.StaticIpv6Dist,
            OspfIntDist: e.OspfIntDist,
            OspfExtDist: e.OspfExtDist,
            Ospfv3IntDist: e.Ospfv3IntDist,
            Ospfv3ExtDist: e.Ospfv3ExtDist,
            IbgpDist: e.IbgpDist,
            EbgpDist: e.EbgpDist,
            RipDist: e.RipDist,
        },
    }
    if text, present := e.raw["ecmp"]; present {
        ans.Ecmp = &util.RawXml{text}
    }
    if text, present := e.raw["multicast"]; present {
        ans.Multicast = &util.RawXml{text}
    }
    if text, present := e.raw["protocol"]; present {
        ans.Protocol = &util.RawXml{text}
    }
    if text, present := e.raw["routing"]; present {
        ans.Routing = &util.RawXml{text}
    }

    return ans
}
