// Package security is the client.Policies.Security namespace.
//
// Normalized object:  Entry
package security

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a security
// rule.
type Entry struct {
    Name string
    Type string
    Description string
    Tags []string
    SourceZone []string
    SourceAddress []string
    NegateSource bool
    SourceUser []string
    HipProfile []string
    DestinationZone []string
    DestinationAddress []string
    NegateDestination bool
    Application []string
    Service []string
    Category []string
    Action string
    LogSetting string
    LogStart bool
    LogEnd bool
    Disabled bool
    Schedule string
    IcmpUnreachable bool
    DisableServerResponseInspection bool
    Group string
    Target []string
    NegateTarget bool
    Virus string
    Spyware string
    Vulnerability string
    UrlFiltering string
    FileBlocking string
    WildFireAnalysis string
    DataFiltering string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * Type: "universal"
//      * SourceZone: ["any"]
//      * SourceAddress: ["any"]
//      * SourceUser: ["any"]
//      * HipProfile: ["any"]
//      * DestinationZone: ["any"]
//      * DestinationAddress: ["any"]
//      * Application: ["any"]
//      * Service: ["application-default"]
//      * Category: ["any"]
//      * Action: "allow"
//      * LogEnd: true
func (o *Entry) Defaults() {
    if o.Type == "" {
        o.Type = "universal"
    }

    if len(o.SourceZone) == 0 {
        o.SourceZone = []string{"any"}
    }

    if len(o.DestinationZone) == 0 {
        o.DestinationZone = []string{"any"}
    }

    if len(o.SourceAddress) == 0 {
        o.SourceAddress = []string{"any"}
    }

    if len(o.SourceUser) == 0 {
        o.SourceUser = []string{"any"}
    }

    if len(o.HipProfile) == 0 {
        o.HipProfile = []string{"any"}
    }

    if len(o.DestinationAddress) == 0 {
        o.DestinationAddress = []string{"any"}
    }

    if len(o.Application) == 0 {
        o.Application = []string{"any"}
    }

    if len(o.Service) == 0 {
        o.Service = []string{"application-default"}
    }

    if len(o.Category) == 0 {
        o.Category = []string{"any"}
    }

    if o.Action == "" {
        o.Action = "allow"
    }

    if !o.LogEnd {
        o.LogEnd = true
    }
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Type = s.Type
    o.Description = s.Description
    o.Tags = s.Tags
    o.SourceZone = s.SourceZone
    o.SourceAddress = s.SourceAddress
    o.NegateSource = s.NegateSource
    o.SourceUser = s.SourceUser
    o.HipProfile = s.HipProfile
    o.DestinationZone = s.DestinationZone
    o.DestinationAddress = s.DestinationAddress
    o.NegateDestination = s.NegateDestination
    o.Application = s.Application
    o.Service = s.Service
    o.Category = s.Category
    o.Action = s.Action
    o.LogSetting = s.LogSetting
    o.LogStart = s.LogStart
    o.LogEnd = s.LogEnd
    o.Disabled = s.Disabled
    o.Schedule = s.Schedule
    o.IcmpUnreachable = s.IcmpUnreachable
    o.DisableServerResponseInspection = s.DisableServerResponseInspection
    o.Group = s.Group
    o.Target = s.Target
    o.NegateTarget = s.NegateTarget
    o.Virus = s.Virus
    o.Spyware = s.Spyware
    o.Vulnerability = s.Vulnerability
    o.UrlFiltering = s.UrlFiltering
    o.FileBlocking = s.FileBlocking
    o.WildFireAnalysis = s.WildFireAnalysis
    o.DataFiltering = s.DataFiltering
}

// Security is the client.Policies.Security namespace.
type Security struct {
    con util.XapiClient
}

// Initialize is invoed by client.Initialize().
func (c *Security) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of security policies.
func (c *Security) GetList(vsys, base string) ([]string, error) {
    c.con.LogQuery("(get) list of security policies")
    path := c.xpath(vsys, base, nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of security policies.
func (c *Security) ShowList(vsys, base string) ([]string, error) {
    c.con.LogQuery("(show) list of security policies")
    path := c.xpath(vsys, base, nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given security policy.
func (c *Security) Get(vsys, base, name string) (Entry, error) {
    c.con.LogQuery("(get) security policy %q", name)
    return c.details(c.con.Get, vsys, base, name)
}

// Get performs SHOW to retrieve information for the given security policy.
func (c *Security) Show(vsys, base, name string) (Entry, error) {
    c.con.LogQuery("(show) security policy %q", name)
    return c.details(c.con.Show, vsys, base, name)
}

// Set performs SET to create / update one or more security policies.
func (c *Security) Set(vsys, base string, e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given security policy configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "rules"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) security policies: %v", names)

    // Set xpath.
    path := c.xpath(vsys, base, names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the security policies.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// VerifiableSet behaves like Set(), except policies with LogEnd as true
// will first be created with LogEnd as false, and then a second Set() is
// performed which will do LogEnd as true.  This is due to the unique
// combination of being a boolean value that is true by default, the XML
// returned from querying the rule details will omit the LogEnd setting,
// which will be interpreted as false, when in fact it is true.  We can
// get around this by setting the value to a non-standard value, then back
// again, in which case it will properly show up in the returned XML.
func (c *Security) VerifiableSet(vsys, base string, e ...Entry) error {
    c.con.LogAction("(set) performing verifiable set")
    again := make([]Entry, 0, len(e))

    for i := range e {
        if e[i].LogEnd {
            again = append(again, e[i])
            e[i].LogEnd = false
        }
    }

    if err := c.Set(vsys, base, e...); err != nil {
        return err
    }

    if len(again) == 0 {
        return nil
    }

    return c.Set(vsys, base, again...)
}

// Edit performs EDIT to create / update a security policy.
func (c *Security) Edit(vsys, base string, e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) security policy %q", e.Name)

    // Set xpath.
    path := c.xpath(vsys, base, []string{e.Name})

    // Edit the security policy.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given security policies.
//
// Security policies can be either a string or an Entry object.
func (c *Security) Delete(vsys, base string, e ...interface{}) error {
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
    c.con.LogAction("(delete) security policies: %v", names)

    path := c.xpath(vsys, base, names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

// DeleteAll removes all security policies from the specified vsys / rulebase.
func (c *Security) DeleteAll(vsys, base string) error {
    c.con.LogAction("(delete) all security policies")
    list, err := c.GetList(vsys, base)
    if err != nil || len(list) == 0 {
        return err
    }
    li := make([]interface{}, len(list))
    for i := range list {
        li[i] = list[i]
    }
    return c.Delete(vsys, base, li...)
}

/** Internal functions for the Zone struct **/

func (c *Security) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *Security) details(fn util.Retriever, vsys, base, name string) (Entry, error) {
    path := c.xpath(vsys, base, []string{name})
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *Security) xpath(vsys, base string, vals []string) []string {
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
        "security",
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
        Type: o.Answer.Type,
        Description: o.Answer.Description,
        Tags: util.MemToStr(o.Answer.Tags),
        SourceZone: util.MemToStr(o.Answer.SourceZone),
        DestinationZone: util.MemToStr(o.Answer.DestinationZone),
        SourceAddress: util.MemToStr(o.Answer.SourceAddress),
        NegateSource: util.AsBool(o.Answer.NegateSource),
        SourceUser: util.MemToStr(o.Answer.SourceUser),
        HipProfile: util.MemToStr(o.Answer.HipProfile),
        DestinationAddress: util.MemToStr(o.Answer.DestinationAddress),
        NegateDestination: util.AsBool(o.Answer.NegateDestination),
        Application: util.MemToStr(o.Answer.Application),
        Service: util.MemToStr(o.Answer.Service),
        Category: util.MemToStr(o.Answer.Category),
        Action: o.Answer.Action,
        LogSetting: o.Answer.LogSetting,
        LogStart: util.AsBool(o.Answer.LogStart),
        LogEnd: util.AsBool(o.Answer.LogEnd),
        Disabled: util.AsBool(o.Answer.Disabled),
        Schedule: o.Answer.Schedule,
        IcmpUnreachable: util.AsBool(o.Answer.IcmpUnreachable),
    }
    if o.Answer.Options != nil {
        ans.DisableServerResponseInspection = util.AsBool(o.Answer.Options.DisableServerResponseInspection)
    }
    if o.Answer.TargetInfo != nil {
        ans.NegateTarget = util.AsBool(o.Answer.TargetInfo.NegateTarget)
        ans.Target = util.EntToStr(o.Answer.TargetInfo.Target)
    }
    if o.Answer.ProfileSettings != nil {
        ans.Group = util.MemToOneStr(o.Answer.ProfileSettings.Group)
        if o.Answer.ProfileSettings.Profiles != nil {
            ans.Virus = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.Virus)
            ans.Spyware = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.Spyware)
            ans.Vulnerability = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.Vulnerability)
            ans.UrlFiltering = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.UrlFiltering)
            ans.FileBlocking = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.FileBlocking)
            ans.WildFireAnalysis = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.WildFireAnalysis)
            ans.DataFiltering = util.MemToOneStr(o.Answer.ProfileSettings.Profiles.DataFiltering)
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Type string `xml:"rule-type"`
    Description string `xml:"description"`
    Tags *util.Member `xml:"tag"`
    SourceZone *util.Member `xml:"from"`
    DestinationZone *util.Member `xml:"to"`
    SourceAddress *util.Member `xml:"source"`
    NegateSource string `xml:"negate-source"`
    SourceUser *util.Member `xml:"source-user"`
    HipProfile *util.Member `xml:"hip-profiles"`
    DestinationAddress *util.Member `xml:"destination"`
    NegateDestination string `xml:"negate-destination"`
    Application *util.Member `xml:"application"`
    Service *util.Member `xml:"service"`
    Category *util.Member `xml:"category"`
    Action string `xml:"action"`
    LogSetting string `xml:"log-setting,omitempty"`
    LogStart string `xml:"log-start"`
    LogEnd string `xml:"log-end"`
    Disabled string `xml:"disabled"`
    Schedule string `xml:"schedule,omitempty"`
    IcmpUnreachable string `xml:"icmp-unreachable"`
    Options *secOptions `xml:"option"`
    TargetInfo *targetInfo `xml:"target"`
    ProfileSettings *profileSettings `xml:"profile-setting"`
}

type secOptions struct {
    DisableServerResponseInspection string `xml:"disable-server-response-inspection,omitempty"`
}

type targetInfo struct {
    Target *util.Entry `xml:"devices"`
    NegateTarget string `xml:"negate,omitempty"`
}

type profileSettings struct {
    Group *util.Member `xml:"group"`
    Profiles *profileSettingsProfile `xml:"profiles"`
}

type profileSettingsProfile struct {
    Virus *util.Member `xml:"virus"`
    Spyware *util.Member `xml:"spyware"`
    Vulnerability *util.Member `xml:"vulnerability"`
    UrlFiltering *util.Member `xml:"url-filtering"`
    FileBlocking *util.Member `xml:"file-blocking"`
    WildFireAnalysis *util.Member `xml:"wildfire-analysis"`
    DataFiltering *util.Member `xml:"data-filtering"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Type: e.Type,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
        SourceZone: util.StrToMem(e.SourceZone),
        DestinationZone: util.StrToMem(e.DestinationZone),
        SourceAddress: util.StrToMem(e.SourceAddress),
        NegateSource: util.YesNo(e.NegateSource),
        SourceUser: util.StrToMem(e.SourceUser),
        HipProfile: util.StrToMem(e.HipProfile),
        DestinationAddress: util.StrToMem(e.DestinationAddress),
        NegateDestination: util.YesNo(e.NegateDestination),
        Application: util.StrToMem(e.Application),
        Service: util.StrToMem(e.Service),
        Category: util.StrToMem(e.Category),
        Action: e.Action,
        LogSetting: e.LogSetting,
        LogStart: util.YesNo(e.LogStart),
        LogEnd: util.YesNo(e.LogEnd),
        Disabled: util.YesNo(e.Disabled),
        Schedule: e.Schedule,
        IcmpUnreachable: util.YesNo(e.IcmpUnreachable),
        Options: &secOptions{util.YesNo(e.DisableServerResponseInspection)},
    }
    if e.Target != nil || e.NegateTarget {
        nfo := &targetInfo{
            Target: util.StrToEnt(e.Target),
            NegateTarget: util.YesNo(e.NegateTarget),
        }
        ans.TargetInfo = nfo
    }
    gs := e.Virus != "" || e.Spyware != "" || e.Vulnerability != "" || e.UrlFiltering != "" || e.FileBlocking != "" || e.WildFireAnalysis != "" || e.DataFiltering != ""
    if e.Group != "" || gs {
        ps := &profileSettings{
            Group: util.OneStrToMem(e.Group),
        }
        if gs {
            ps.Profiles = &profileSettingsProfile{
                util.OneStrToMem(e.Virus),
                util.OneStrToMem(e.Spyware),
                util.OneStrToMem(e.Vulnerability),
                util.OneStrToMem(e.UrlFiltering),
                util.OneStrToMem(e.FileBlocking),
                util.OneStrToMem(e.WildFireAnalysis),
                util.OneStrToMem(e.DataFiltering),
            }
        }
        ans.ProfileSettings = ps
    }

    return ans
}
