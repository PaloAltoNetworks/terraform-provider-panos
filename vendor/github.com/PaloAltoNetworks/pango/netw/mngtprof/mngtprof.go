// Package mngtprof is the client.Network.ManagementProfile namespace.
//
// Normalized object:  Entry
package mngtprof

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of an interface
// management profile.
type Entry struct {
    Name string
    Ping bool
    Telnet bool
    Ssh bool
    Http bool
    HttpOcsp bool
    Https bool
    Snmp bool
    ResponsePages bool
    UseridService bool
    UseridSyslogListenerSsl bool
    UseridSyslogListenerUdp bool
    PermittedIp []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Ping = s.Ping
    o.Telnet = s.Telnet
    o.Ssh = s.Ssh
    o.Http = s.Http
    o.HttpOcsp = s.HttpOcsp
    o.Https = s.Https
    o.Snmp = s.Snmp
    o.ResponsePages = s.ResponsePages
    o.UseridService = s.UseridService
    o.UseridSyslogListenerSsl = s.UseridSyslogListenerSsl
    o.UseridSyslogListenerUdp = s.UseridSyslogListenerUdp
    o.PermittedIp = s.PermittedIp
}

// MngtProf is a namespace struct, included as part of pango.Client.
type MngtProf struct {
    con util.XapiClient
}

// Initialize is invoked when Initialize on the pango.Client is called.
func (c *MngtProf) Initialize(con util.XapiClient) {
    c.con = con
}

// GetList performs GET to retrieve a list of interface management profiles.
func (c *MngtProf) GetList() ([]string, error) {
    c.con.LogQuery("(get) list of interface management profiles")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Get, path[:len(path) - 1])
}

// ShowList performs SHOW to retrieve a list of interface management profiles.
func (c *MngtProf) ShowList() ([]string, error) {
    c.con.LogQuery("(show) list of interface management profiles")
    path := c.xpath(nil)
    return c.con.EntryListUsing(c.con.Show, path[:len(path) - 1])
}

// Get performs GET to retrieve information for the given interface management
// profile.
func (c *MngtProf) Get(name string) (Entry, error) {
    c.con.LogQuery("(get) interface management profile %q", name)
    return c.details(c.con.Get, name)
}

// Get performs SHOW to retrieve information for the given interface management
// profile.
func (c *MngtProf) Show(name string) (Entry, error) {
    c.con.LogQuery("(show) interface management profile %q", name)
    return c.details(c.con.Show, name)
}

// Set performs SET to create / update one or more interface management profiles.
func (c *MngtProf) Set(e ...Entry) error {
    var err error

    if len(e) == 0 {
        return nil
    }

    _, fn := c.versioning()
    names := make([]string, len(e))

    // Build up the struct with the given configs.
    d := util.BulkElement{XMLName: xml.Name{Local: "interface-management-profile"}}
    for i := range e {
        d.Data = append(d.Data, fn(e[i]))
        names[i] = e[i].Name
    }
    c.con.LogAction("(set) interface management profiles: %v", names)

    // Set xpath.
    path := c.xpath(names)
    if len(e) == 1 {
        path = path[:len(path) - 1]
    } else {
        path = path[:len(path) - 2]
    }

    // Create the profiles.
    _, err = c.con.Set(path, d.Config(), nil, nil)
    return err
}

// Edit performs EDIT to create / update an interface management profile.
func (c *MngtProf) Edit(e Entry) error {
    var err error

    _, fn := c.versioning()

    c.con.LogAction("(edit) interface management profile %q", e.Name)

    // Set xpath.
    path := c.xpath([]string{e.Name})

    // Edit the profile.
    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

// Delete removes the given interface management profile(s) from the firewall.
//
// Profiles can be either a string or an Entry object.
func (c *MngtProf) Delete(e ...interface{}) error {
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
    c.con.LogAction("(delete) interface management profiles: %v", names)

    path := c.xpath(names)
    _, err = c.con.Delete(path, nil, nil)
    return err
}

/** Internal functions for the MngtProf struct **/

func (c *MngtProf) versioning() (normalizer, func(Entry) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *MngtProf) details(fn util.Retriever, name string) (Entry, error) {
    path := c.xpath([]string{name})
    obj, _ := c.versioning()
    _, err := fn(path, nil, obj)
    if err != nil {
        return Entry{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *MngtProf) xpath(vals []string) []string {
    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "network",
        "profiles",
        "interface-management-profile",
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
        Ping: util.AsBool(o.Answer.Ping),
        Telnet: util.AsBool(o.Answer.Telnet),
        Ssh: util.AsBool(o.Answer.Ssh),
        Http: util.AsBool(o.Answer.Http),
        HttpOcsp: util.AsBool(o.Answer.HttpOcsp),
        Https: util.AsBool(o.Answer.Https),
        Snmp: util.AsBool(o.Answer.Snmp),
        ResponsePages: util.AsBool(o.Answer.ResponsePages),
        UseridService: util.AsBool(o.Answer.UseridService),
        UseridSyslogListenerSsl: util.AsBool(o.Answer.UseridSyslogListenerSsl),
        UseridSyslogListenerUdp: util.AsBool(o.Answer.UseridSyslogListenerUdp),
        PermittedIp: util.EntToStr(o.Answer.PermittedIp),
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Ping string `xml:"ping"`
    Telnet string `xml:"telnet"`
    Ssh string `xml:"ssh"`
    Http string `xml:"http"`
    HttpOcsp string `xml:"http-ocsp"`
    Https string `xml:"https"`
    Snmp string `xml:"snmp"`
    ResponsePages string `xml:"response-pages"`
    UseridService string `xml:"userid-service"`
    UseridSyslogListenerSsl string `xml:"userid-syslog-listener-ssl"`
    UseridSyslogListenerUdp string `xml:"userid-syslog-listener-udp"`
    PermittedIp *util.Entry `xml:"permitted-ip"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Ping: util.YesNo(e.Ping),
        Telnet: util.YesNo(e.Telnet),
        Ssh: util.YesNo(e.Ssh),
        Http: util.YesNo(e.Http),
        HttpOcsp: util.YesNo(e.HttpOcsp),
        Https: util.YesNo(e.Https),
        Snmp: util.YesNo(e.Snmp),
        ResponsePages: util.YesNo(e.ResponsePages),
        UseridService: util.YesNo(e.UseridService),
        UseridSyslogListenerSsl: util.YesNo(e.UseridSyslogListenerSsl),
        UseridSyslogListenerUdp: util.YesNo(e.UseridSyslogListenerUdp),
        PermittedIp: util.StrToEnt(e.PermittedIp),
    }

    return ans
}
