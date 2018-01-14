// Package srvc is the client.Objects.Services namespace.
//
// Normalized object:  Entry
package srvc

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a service
// object.
//
// Protocol should be either "tcp" or "udp".
type Entry struct {
    Name string
    Description string
    Protocol string
    SourcePort string
    DestinationPort string
    Tags []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Protocol = s.Protocol
    o.SourcePort = s.SourcePort
    o.DestinationPort = s.DestinationPort
    o.Tags = s.Tags
}

// Srvc is a namespace struct, included as part of pango.Client.
type Srvc struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *Srvc) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of service objects.
func (c *Srvc) GetList(vsys string) ([]string, error) {
    c.con.LogQuery("(get) list of service objects")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of service objects.
func (c *Srvc) ShowList(vsys string) ([]string, error) {
    c.con.LogQuery("(show) list of service objects")
    path := c.xpath(vsys, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given service object.
func (c *Srvc) Get(vsys, name string) (Entry, error) {
    c.con.LogQuery("(get) service object %q", name)
    return c.details(c.con.Get, vsys, name)
}

// Get performs SHOW to retrieve information for the given service object.
func (c *Srvc) Show(vsys, name string) (Entry, error) {
    c.con.LogQuery("(show) service object %q", name)
    return c.details(c.con.Show, vsys, name)
}

// Set performs SET to create / update one or more service objects.
func (c *Srvc) Set(vsys string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "service"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) service objects: %v", names)

    // Set xpath.
    path := c.xpath(vsys, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the objects.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update a service object.
func (c *Srvc) Edit(vsys string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) service object %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, []string{e.Name})

    // Create the object.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given service objects from the firewall.
//
// Service objects can be either a string or an Entry object.
func (c *Srvc) Delete(vsys string, e ...interface{}) error {
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
    c.con.LogAction("(delete) service objects: %v", names)

    path := c.xpath(vsys, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Srvc struct **/

func (c *Srvc) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Srvc) details(fn util.Retriever, vsys, name string) (Entry, error) {
    path := c.xpath(vsys, []string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Srvc) xpath(vsys string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }

    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "service",
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
        Description: o.Answer.Description,
        Tags: util.MemToStr(o.Answer.Tags),
    }
    switch {
    case o.Answer.TcpProto != nil:
        ans.Protocol = "tcp"
        ans.SourcePort = o.Answer.TcpProto.SourcePort
        ans.DestinationPort = o.Answer.TcpProto.DestinationPort
    case o.Answer.UdpProto != nil:
        ans.Protocol = "udp"
        ans.SourcePort = o.Answer.UdpProto.SourcePort
        ans.DestinationPort = o.Answer.UdpProto.DestinationPort
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TcpProto *protoDef `xml:"protocol>tcp"`
    UdpProto *protoDef `xml:"protocol>udp"`
    Description string `xml:"description"`
    Tags *util.Member `xml:"tag"`
}

type protoDef struct {
    SourcePort string `xml:"source-port,omitempty"`
    DestinationPort string `xml:"port"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
    }
    switch e.Protocol {
    case "tcp":
        ans.TcpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    case "udp":
        ans.UdpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    }

    return ans
}
