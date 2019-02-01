// Package util contains various shared structs and functions used across
// the pango package.
package util


import (
    "bytes"
    "encoding/xml"
    "fmt"
    "regexp"
    "strings"

    "github.com/PaloAltoNetworks/pango/version"
)

// Retriever is a type that is intended to act as a stand-in for using
// either the Get or Show pango Client functions.
type Retriever func(interface{}, interface{}, interface{}) ([]byte, error)

// Rulebase constants for various policies.
const (
    Rulebase = "rulebase"
    PreRulebase = "pre-rulebase"
    PostRulebase = "post-rulebase"
)

// Valid values to use for VsysImport() or VsysUnimport().
const (
    InterfaceImport = "interface"
    VirtualRouterImport = "virtual-router"
    VirtualWireImport = "virtual-wire"
    VlanImport = "vlan"
)

// XapiClient is the interface that describes an pango.Client.
type XapiClient interface {
    String() string
    Versioning() version.Number
    LogAction(string, ...interface{})
    LogQuery(string, ...interface{})
    LogOp(string, ...interface{})
    LogUid(string, ...interface{})
    Op(interface{}, string, interface{}, interface{}) ([]byte, error)
    Show(interface{}, interface{}, interface{}) ([]byte, error)
    Get(interface{}, interface{}, interface{}) ([]byte, error)
    Delete(interface{}, interface{}, interface{}) ([]byte, error)
    Set(interface{}, interface{}, interface{}, interface{}) ([]byte, error)
    Edit(interface{}, interface{}, interface{}, interface{}) ([]byte, error)
    Move(interface{}, string, string, interface{}, interface{}) ([]byte, error)
    Uid(interface{}, string, interface{}, interface{}) ([]byte, error)
    EntryListUsing(Retriever, []string) ([]string, error)
    MemberListUsing(Retriever, []string) ([]string, error)
    RequestPasswordHash(string) (string, error)
    VsysImport(string, string, string, string, []string) error
    VsysUnimport(string, string, string, []string) error
    WaitForJob(uint, interface{}) error
    Commit(string, bool, bool, bool, bool) (uint, error)
    PositionFirstEntity(int, string, string, []string, []string) error
}

// BulkElement is a generic bulk container for bulk operations.
type BulkElement struct {
    XMLName xml.Name
    Data []interface{}
}

// Config returns an interface to be Marshaled.
func (o BulkElement) Config() interface{} {
    if len(o.Data) == 1 {
        return o.Data[0]
    }
    return o
}

// MemberType defines a member config node used for sending and receiving XML
// from PAN-OS.
type MemberType struct {
    Members []Member `xml:"member"`
}

// Member defines a member config node used for sending and receiving XML
// from PANOS.
type Member struct {
    XMLName xml.Name `xml:"member"`
    Value string `xml:",chardata"`
}

// MemToStr normalizes a MemberType pointer into a list of strings.
func MemToStr(e *MemberType) []string {
    if e == nil {
        return nil
    }

    ans := make([]string, len(e.Members))
    for i := range e.Members {
        ans[i] = e.Members[i].Value
    }

    return ans
}

// StrToMem converts a list of strings into a MemberType pointer.
func StrToMem(e []string) *MemberType {
    if e == nil {
        return nil
    }

    ans := make([]Member, len(e))
    for i := range e {
        ans[i] = Member{Value: e[i]}
    }

    return &MemberType{ans}
}

// MemToOneStr normalizes a MemberType pointer for a max_items=1 XML node
// into a string.
func MemToOneStr(e *MemberType) string {
    if e == nil || len(e.Members) == 0 {
        return ""
    }

    return e.Members[0].Value
}

// OneStrToMem converts a string into a MemberType pointer for a max_items=1
// XML node.
func OneStrToMem(e string) *MemberType {
    if e == "" {
        return nil
    }

    return &MemberType{[]Member{
        {Value: e},
    }}
}

// EntryType defines an entry config node used for sending and receiving XML
// from PAN-OS.
type EntryType struct {
    Entries []Entry `xml:"entry"`
}

// Entry is a standalone entry struct.
type Entry struct {
    XMLName xml.Name `xml:"entry"`
    Value string `xml:"name,attr"`
}

// EntToStr normalizes an EntryType pointer into a list of strings.
func EntToStr(e *EntryType) []string {
    if e == nil {
        return nil
    }

    ans := make([]string, len(e.Entries))
    for i := range e.Entries {
        ans[i] = e.Entries[i].Value
    }

    return ans
}

// StrToEnt converts a list of strings into an EntryType pointer.
func StrToEnt(e []string) *EntryType {
    if e == nil {
        return nil
    }

    ans := make([]Entry, len(e))
    for i := range e {
        ans[i] = Entry{Value: e[i]}
    }

    return &EntryType{ans}
}

// EntToOneStr normalizes an EntryType pointer for a max_items=1 XML node
// into a string.
func EntToOneStr(e *EntryType) string {
    if e == nil || len(e.Entries) == 0 {
        return ""
    }

    return e.Entries[0].Value
}

// OneStrToEnt converts a string into an EntryType pointer for a max_items=1
// XML node.
func OneStrToEnt(e string) *EntryType {
    if e == "" {
        return nil
    }

    return &EntryType{[]Entry{
        {Value: e},
    }}
}

// VsysEntryType defines an entry config node with vsys entries underneath.
type VsysEntryType struct {
    Entries []VsysEntry `xml:"entry"`
}

// VsysEntry defines the "vsys" xpath node under a VsysEntryType config node.
type VsysEntry struct {
    XMLName xml.Name `xml:"entry"`
    Serial string `xml:"name,attr"`
    Vsys *EntryType `xml:"vsys"`
}

// VsysEntToMap normalizes a VsysEntryType pointer into a map.
func VsysEntToMap(ve *VsysEntryType) (map[string] []string) {
    if ve == nil {
        return nil
    }

    ans := make(map[string] []string)
    for i := range ve.Entries {
        ans[ve.Entries[i].Serial] = EntToStr(ve.Entries[i].Vsys)
    }

    return ans
}

// MapToVsysEnt converts a map into a VsysEntryType pointer.
//
// This struct is used for "Target" information on Panorama when dealing with
// various policies.  Maps are unordered, but FWICT Panorama doesn't seem to
// order anything anyways when doing things in the GUI, so hopefully this is
// ok...?
func MapToVsysEnt(e map[string] []string) *VsysEntryType {
    if len(e) == 0 {
        return nil
    }

    i := 0
    ve := make([]VsysEntry, len(e))
    for key := range e {
        ve[i].Serial = key
        ve[i].Vsys = StrToEnt(e[key])
        i++
    }

    return &VsysEntryType{ve}
}

// YesNo returns "yes" on true, "no" on false.
func YesNo(v bool) string {
    if v {
        return "yes"
    }
    return "no"
}

// AsBool returns true on yes, else false.
func AsBool(val string) bool {
    if val == "yes" {
        return true
    }
    return false
}

// AsXpath makes an xpath out of the given interface.
func AsXpath(i interface{}) string {
    switch val := i.(type) {
    case string:
        return val
    case []string:
        return fmt.Sprintf("/%s", strings.Join(val, "/"))
    default:
        return ""
    }
}

// AsEntryXpath returns the given values as an entry xpath segment.
func AsEntryXpath(vals []string) string {
    if len(vals) == 0 || (len(vals) == 1 && vals[0] == "") {
        return "entry"
    }

    var buf bytes.Buffer

    buf.WriteString("entry[")
    for i := range vals {
        if i != 0 {
            buf.WriteString(" or ")
        }
        buf.WriteString("@name='")
        buf.WriteString(vals[i])
        buf.WriteString("'")
    }
    buf.WriteString("]")

    return buf.String()
}

// AsMemberXpath returns the given values as a member xpath segment.
func AsMemberXpath(vals []string) string {
    var buf bytes.Buffer

    buf.WriteString("member[")
    for i := range vals {
        if i != 0 {
            buf.WriteString(" or ")
        }
        buf.WriteString("text()='")
        buf.WriteString(vals[i])
        buf.WriteString("'")
    }

    buf.WriteString("]")

    return buf.String()
}

// TemplateXpath returns the template xpath prefix of the given template name.
func TemplateXpathPrefix(tmpl, ts string) []string {
    if tmpl != "" {
        return []string{
            "config",
            "devices",
            AsEntryXpath([]string{"localhost.localdomain"}),
            "template",
            AsEntryXpath([]string{tmpl}),
        }
    }

    return []string{
        "config",
        "devices",
        AsEntryXpath([]string{"localhost.localdomain"}),
        "template-stack",
        AsEntryXpath([]string{ts}),
    }
}

// License defines a license entry.
type License struct {
    XMLName xml.Name `xml:"entry"`
    Feature string `xml:"feature"`
    Description string `xml:"description"`
    Serial string `xml:"serial"`
    Issued string `xml:"issued"`
    Expires string `xml:"expires"`
    Expired string `xml:"expired"`
    AuthCode string `xml:"authcode"`
}

// Lock represents either a config lock or a commit lock.
type Lock struct {
    XMLName xml.Name `xml:"entry"`
    Owner string `xml:"name,attr"`
    Name string `xml:"name"`
    Type string `xml:"type"`
    LoggedIn string `xml:"loggedin"`
    Comment CdataText `xml:"comment"`
}

// CdataText is for getting CDATA contents of XML docs.
type CdataText struct {
    Text string `xml:",cdata"`
}

// RawXml is what allows the use of Edit commands on a XPATH without
// truncating any other child objects that may be attached to it.
type RawXml struct {
    Text string `xml:",innerxml"`
}

// CleanRawXml removes extra XML attributes from RawXml objects without
// requiring us to have to parse everything.
func CleanRawXml(v string) string {
    re := regexp.MustCompile(` admin="\S+" dirtyId="\d+" time="\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}"`)
    return re.ReplaceAllString(v, "")
}

// JobResponse parses a XML response that includes a job ID.
type JobResponse struct {
    XMLName xml.Name `xml:"response"`
    Id uint `xml:"result>job"`
}

// BasicJob is a struct for parsing minimal information about a submitted
// job to PANOS.
type BasicJob struct {
    XMLName xml.Name `xml:"response"`
    Result string `xml:"result>job>result"`
    Progress uint `xml:"result>job>progress"`
    Details []string `xml:"result>job>details>line"`
    Devices []devJob `xml:"result>job>devices>entry"`
}

// Internally used by BasicJob for panorama commit-all.
type devJob struct {
    Serial string `xml:"serial-no"`
    Result string `xml:"result"`
}

// These constants are valid move locations to pass to various movement
// functions (aka - policy management).
const (
    MoveSkip = iota
    MoveBefore
    MoveDirectlyBefore
    MoveAfter
    MoveDirectlyAfter
    MoveTop
    MoveBottom
)
