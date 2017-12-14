// Package nat is the client.Policies.Nat namespace.
//
// Normalized object:  Entry
package nat

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Values for Entry.SatType.
const (
    DynamicIpAndPort = "dynamic-ip-and-port"
    DynamicIp = "dynamic-ip"
    StaticIp = "static-ip"
)

// Values for Entry.SatAddressType.
const (
    InterfaceAddress = "interface-address"
    TranslatedAddress = "translated-address"
)

// None is a valid value for both Entry.SatType and Entry.SatAddressType.
const None = "none"

// These are the valid settings for Entry.SatFallbackIpType.
const (
    Ip = "ip"
    FloatingIp = "floating"
)

// Entry is a normalized, version independent representation of a NAT
// policy.  The prefix "Sat" stands for "Source Address Translation" while
// the prefix "Dat" stands for "Destination Address Translation".
//
// The following Sat params are linked:
//
// SatType = nat.DynamicIpAndPort && SatAddressType = nat.TranslatedAddress:
//
//      * SatTranslatedAddress
//
// SatType = nat.DynamicIpAndPort && SatAddressType = nat.InterfaceAddress:
//
//      * SatInterface
//      * SatIpAddress
//
// For ALL SatType = nat.DynamicIp:
//
//      * SatTranslatedAddress
//
// For ALL SatType = nat.DynamicIp and SatFallbackType = nat.InterfaceAddress:
//
//      * SatFallbackInterface
//
// SatType = nat.DynamicIp && SatFallbackType = nat.InterfaceAddress && SatFallbackIpType = nat.Ip:
//
//      * SatFallbackIpAddress
//
// SatType = nat.DynamicIp && SatFallbackType = nat.InterfaceAddress && SatFallbackIpType = nat.FloatingIp:
//
//      * SatFallbackIpAddress
//
// SatType = nat.DynamicIp and SatFallbackType = nat.TranslatedAddress:
//
//      * SatFallbackTranslatedAddress
//
// SatType = nat.StaticIp:
//
//      * SatStaticTranslatedAddress
//      * SatStaticBiDirectional
//
// If both DatAddress and DatPort are unintialized, then no destination
// address translation will be enabled.
type Entry struct {
    Name string
    Description string
    Type string
    SourceZone []string
    DestinationZone string
    ToInterface string
    Service string
    SourceAddress []string
    DestinationAddress []string
    SatType string
    SatAddressType string
    SatTranslatedAddress []string
    SatInterface string
    SatIpAddress string
    SatFallbackType string
    SatFallbackTranslatedAddress []string
    SatFallbackInterface string
    SatFallbackIpType string
    SatFallbackIpAddress string
    SatStaticTranslatedAddress string
    SatStaticBiDirectional bool
    DatAddress string
    DatPort int
    Disabled bool
    Target []string
    NegateTarget bool
    Tag []string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * Type: "ipv4"
//      * ToInterface: "any"
//      * Service: "any"
//      * SourceAddress: ["any"]
//      * DestinationAddress: ["any"]
//      * SatType: None
func (o *Entry) Defaults() {
    if o.Type == "" {
        o.Type = "ipv4"
    }

    if o.ToInterface == "" {
        o.ToInterface = "any"
    }

    if o.Service == "" {
        o.Service = "any"
    }

    if len(o.SourceAddress) == 0 {
        o.SourceAddress = []string{"any"}
    }

    if len(o.DestinationAddress) == 0 {
        o.DestinationAddress = []string{"any"}
    }

    if o.SatType == "" {
        o.SatType = None
    }
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Type = s.Type
    o.SourceZone = s.SourceZone
    o.DestinationZone = s.DestinationZone
    o.ToInterface = s.ToInterface
    o.Service = s.Service
    o.SourceAddress = s.SourceAddress
    o.DestinationAddress = s.DestinationAddress
    o.SatType = s.SatType
    o.SatAddressType = s.SatAddressType
    o.SatTranslatedAddress = s.SatTranslatedAddress
    o.SatInterface = s.SatInterface
    o.SatIpAddress = s.SatIpAddress
    o.SatFallbackType = s.SatFallbackType
    o.SatFallbackTranslatedAddress = s.SatFallbackTranslatedAddress
    o.SatFallbackInterface = s.SatFallbackInterface
    o.SatFallbackIpType = s.SatFallbackIpType
    o.SatFallbackIpAddress = s.SatFallbackIpAddress
    o.SatStaticTranslatedAddress = s.SatStaticTranslatedAddress
    o.SatStaticBiDirectional = s.SatStaticBiDirectional
    o.DatAddress = s.DatAddress
    o.DatPort = s.DatPort
    o.Disabled = s.Disabled
    o.Target = s.Target
    o.NegateTarget = s.NegateTarget
    o.Tag = s.Tag
}

// Nat is the client.Policies.Nat namespace.
type Nat struct {
    con util.XapiClient
}

// Initialize is invoed by client.Initialize().
func (c *Nat) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of NAT policies.
func (c *Nat) GetList(vsys, base string) ([]string, error) {
    c.con.LogQuery("(get) list of NAT policies")
    path := c.xpath(vsys, base, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of NAT policies.
func (c *Nat) ShowList(vsys, base string) ([]string, error) {
    c.con.LogQuery("(show) list of NAT policies")
    path := c.xpath(vsys, base, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given NAT policy.
func (c *Nat) Get(vsys, base, name string) (Entry, error) {
    c.con.LogQuery("(get) NAT policy %q", name)
    return c.details(c.con.Get, vsys, base, name)
}

// Get performs SHOW to retrieve information for the given NAT policy.
func (c *Nat) Show(vsys, base, name string) (Entry, error) {
    c.con.LogQuery("(show) NAT policy %q", name)
    return c.details(c.con.Show, vsys, base, name)
}

// Set performs SET to create / update one or more NAT policies.
func (c *Nat) Set(vsys, base string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "rules"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) NAT policies: %v", names)

    // Set xpath.
    path := c.xpath(vsys, base, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the NAT policies.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update a NAT policy.
func (c *Nat) Edit(vsys, base string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) NAT policy %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, base, []string{e.Name})

    // Edit the NAT policy.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given NAT policies.
//
// NAT policies can be either a string or an Entry object.
func (c *Nat) Delete(vsys, base string, e ...interface{}) error {
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
    c.con.LogAction("(delete) NAT policies: %v", names)

    path := c.xpath(vsys, base, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the Zone struct **/

func (c *Nat) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Nat) details(fn util.Retriever, vsys, base, name string) (Entry, error) {
    path := c.xpath(vsys, base, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Nat) xpath(vsys, base string, vals []string) []string {
    if vsys == "" {
        vsys = "vsys1"
    }
    if base == "" {
        base = util.Rulebase
    }

    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        base,
        "nat",
        "rules",
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
        Type: o.Answer.Type,
        SourceZone: util.MemToStr(o.Answer.SourceZone),
        DestinationZone: o.Answer.DestinationZone,
        ToInterface: o.Answer.ToInterface,
        Service: o.Answer.Service,
        SourceAddress: util.MemToStr(o.Answer.SourceAddress),
        DestinationAddress: util.MemToStr(o.Answer.DestinationAddress),
        Disabled: util.AsBool(o.Answer.Disabled),
        Tag: util.MemToStr(o.Answer.Tag),
    }

    if o.Answer.Sat == nil {
        ans.SatType = None
    } else {
        switch {
        case o.Answer.Sat.Diap != nil:
            ans.SatType = DynamicIpAndPort
            if o.Answer.Sat.Diap.InterfaceAddress != nil {
                ans.SatAddressType = InterfaceAddress
                ans.SatInterface = o.Answer.Sat.Diap.InterfaceAddress.Interface
                ans.SatIpAddress = o.Answer.Sat.Diap.InterfaceAddress.Ip
            } else {
                ans.SatAddressType = TranslatedAddress
                ans.SatTranslatedAddress = util.MemToStr(o.Answer.Sat.Diap.TranslatedAddress)
            }
        case o.Answer.Sat.Di != nil:
            ans.SatType = DynamicIp
            ans.SatTranslatedAddress = util.MemToStr(o.Answer.Sat.Di.TranslatedAddress)
            if o.Answer.Sat.Di.Fallback == nil {
                ans.SatFallbackType = None
            } else if o.Answer.Sat.Di.Fallback.TranslatedAddress != nil {
                ans.SatFallbackType = TranslatedAddress
                ans.SatFallbackTranslatedAddress = util.MemToStr(o.Answer.Sat.Di.Fallback.TranslatedAddress)
            } else if o.Answer.Sat.Di.Fallback.InterfaceAddress != nil {
                ans.SatFallbackType = InterfaceAddress
                ans.SatFallbackInterface = o.Answer.Sat.Di.Fallback.InterfaceAddress.Interface
                if o.Answer.Sat.Di.Fallback.InterfaceAddress.Ip != "" {
                    ans.SatFallbackIpType = Ip
                    ans.SatFallbackIpAddress = o.Answer.Sat.Di.Fallback.InterfaceAddress.Ip
                } else if o.Answer.Sat.Di.Fallback.InterfaceAddress.FloatingIp != "" {
                    ans.SatFallbackIpType = FloatingIp
                    ans.SatFallbackIpAddress = o.Answer.Sat.Di.Fallback.InterfaceAddress.FloatingIp
                }
            }
        case o.Answer.Sat.Static != nil:
            ans.SatType = StaticIp
            ans.SatStaticTranslatedAddress = o.Answer.Sat.Static.Address
            ans.SatStaticBiDirectional = util.AsBool(o.Answer.Sat.Static.BiDirectional)
        }
    }

    if o.Answer.Dat != nil {
        ans.DatAddress = o.Answer.Dat.Address
        ans.DatPort = o.Answer.Dat.Port
    }

    if o.Answer.Target != nil {
        ans.Target = util.EntToStr(o.Answer.Target.Target)
        ans.NegateTarget = util.AsBool(o.Answer.Target.NegateTarget)
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    Type string `xml:"nat-type"`
    SourceZone *util.Member `xml:"from"`
    DestinationZone string `xml:"to>member"`
    ToInterface string `xml:"to-interface"`
    Service string `xml:"service"`
    SourceAddress *util.Member `xml:"source"`
    DestinationAddress *util.Member `xml:"destination"`
    Sat *srcXlate `xml:"source-translation"`
    Dat *dstXlate `xml:"destination-translation"`
    Disabled string `xml:"disabled"`
    Target *targetInfo `xml:"target"`
    Tag *util.Member `xml:"tag"`
}

type dstXlate struct {
    Address string `xml:"translated-address,omitempty"`
    Port int `xml:"translated-port,omitempty"`
}

type srcXlate struct {
    Diap *srcXlateDiap `xml:"dynamic-ip-and-port"`
    Di *srcXlateDi `xml:"dynamic-ip"`
    Static *srcXlateStatic `xml:"static-ip"`
}

type srcXlateDiap struct {
    TranslatedAddress *util.Member `xml:"translated-address"`
    InterfaceAddress *srcXlateDiapIa `xml:"interface-address"`
}

type srcXlateDiapIa struct {
    Interface string `xml:"interface"`
    Ip string `xml:"ip,omitempty"`
}

type srcXlateDi struct {
    TranslatedAddress *util.Member `xml:"translated-address"`
    Fallback *fallback `xml:"fallback"`
}

type fallback struct {
    TranslatedAddress *util.Member `xml:"translated-address"`
    InterfaceAddress *fallbackIface `xml:"interface-address"`
}

type fallbackIface struct {
    Ip string `xml:"ip,omitempty"`
    Interface string `xml:"interface,omitempty"`
    FloatingIp string `xml:"floating-ip,omitempty"`
}

type srcXlateStatic struct {
    Address string `xml:"translated-address"`
    BiDirectional string `xml:"bi-directional"`
}

type targetInfo struct {
    Target *util.Entry `xml:"devices"`
    NegateTarget string `xml:"negate,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Type: e.Type,
        SourceZone: util.StrToMem(e.SourceZone),
        DestinationZone: e.DestinationZone,
        ToInterface: e.ToInterface,
        Service: e.Service,
        SourceAddress: util.StrToMem(e.SourceAddress),
        DestinationAddress: util.StrToMem(e.DestinationAddress),
        Disabled: util.YesNo(e.Disabled),
        Tag: util.StrToMem(e.Tag),
    }

    var sv *srcXlate
    switch e.SatType {
    case DynamicIpAndPort:
        sv = &srcXlate{
            Diap: &srcXlateDiap{},
        }
        switch e.SatAddressType {
        case TranslatedAddress:
            sv.Diap.TranslatedAddress = util.StrToMem(e.SatTranslatedAddress)
        case InterfaceAddress:
            sv.Diap.InterfaceAddress = &srcXlateDiapIa{
                Interface: e.SatInterface,
                Ip: e.SatIpAddress,
            }
        }
    case DynamicIp:
        sv = &srcXlate{
            Di: &srcXlateDi{
                TranslatedAddress: util.StrToMem(e.SatTranslatedAddress),
            },
        }
        switch e.SatFallbackType {
        case InterfaceAddress:
            sv.Di.Fallback = &fallback{
                InterfaceAddress: &fallbackIface{
                    Interface: e.SatFallbackInterface,
                },
            }
            switch e.SatFallbackIpType {
            case Ip:
                sv.Di.Fallback.InterfaceAddress.Ip = e.SatFallbackIpAddress
            case FloatingIp:
                sv.Di.Fallback.InterfaceAddress.FloatingIp = e.SatFallbackIpAddress
            }
        case TranslatedAddress:
            sv.Di.Fallback = &fallback{TranslatedAddress: util.StrToMem(e.SatFallbackTranslatedAddress)}
        }
    case StaticIp:
        sv = &srcXlate{
            Static: &srcXlateStatic{
                e.SatStaticTranslatedAddress,
                util.YesNo(e.SatStaticBiDirectional),
            },
        }
    }
    ans.Sat = sv

    if e.DatAddress != "" || e.DatPort != 0 {
        ans.Dat = &dstXlate{
            e.DatAddress,
            e.DatPort,
        }
    }

    if len(e.Target) != 0 || e.NegateTarget {
        ans.Target = &targetInfo{
            Target: util.StrToEnt(e.Target),
            NegateTarget: util.YesNo(e.NegateTarget),
        }
    }

    return ans
}
