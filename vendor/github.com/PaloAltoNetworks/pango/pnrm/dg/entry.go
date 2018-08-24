package dg

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a device group.
//
// Devices is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
    Name string
    Description string
    Devices map[string] []string

    raw map[string] string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Devices = s.Devices
}

/** Structs / functions for normalization. **/

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
        Devices: util.VsysEntToMap(o.Answer.Devices),
    }

    ans.raw = make(map[string] string)

    if o.Answer.Ao != nil {
        ans.raw["ao"] = util.CleanRawXml(o.Answer.Ao.Text)
    }
    if o.Answer.Ag != nil {
        ans.raw["ag"] = util.CleanRawXml(o.Answer.Ag.Text)
    }
    if o.Answer.App != nil {
        ans.raw["app"] = util.CleanRawXml(o.Answer.App.Text)
    }
    if o.Answer.Afil != nil {
        ans.raw["afil"] = util.CleanRawXml(o.Answer.Afil.Text)
    }
    if o.Answer.Appg != nil {
        ans.raw["appg"] = util.CleanRawXml(o.Answer.Appg.Text)
    }
    if o.Answer.As != nil {
        ans.raw["as"] = util.CleanRawXml(o.Answer.As.Text)
    }
    if o.Answer.At != nil {
        ans.raw["at"] = util.CleanRawXml(o.Answer.At.Text)
    }
    if o.Answer.Aobj != nil {
        ans.raw["aobj"] = util.CleanRawXml(o.Answer.Aobj.Text)
    }
    if o.Answer.Acode != nil {
        ans.raw["acode"] = util.CleanRawXml(o.Answer.Acode.Text)
    }
    if o.Answer.Email != nil {
        ans.raw["email"] = util.CleanRawXml(o.Answer.Email.Text)
    }
    if o.Answer.Edl != nil {
        ans.raw["edl"] = util.CleanRawXml(o.Answer.Edl.Text)
    }
    if o.Answer.Ls != nil {
        ans.raw["ls"] = util.CleanRawXml(o.Answer.Ls.Text)
    }
    if o.Answer.Master != nil {
        ans.raw["master"] = util.CleanRawXml(o.Answer.Master.Text)
    }
    if o.Answer.Pdf != nil {
        ans.raw["pdf"] = util.CleanRawXml(o.Answer.Pdf.Text)
    }
    if o.Answer.Postrb != nil {
        ans.raw["postrb"] = util.CleanRawXml(o.Answer.Postrb.Text)
    }
    if o.Answer.Prerb != nil {
        ans.raw["prerb"] = util.CleanRawXml(o.Answer.Prerb.Text)
    }
    if o.Answer.Profg != nil {
        ans.raw["profg"] = util.CleanRawXml(o.Answer.Profg.Text)
    }
    if o.Answer.Profs != nil {
        ans.raw["profs"] = util.CleanRawXml(o.Answer.Profs.Text)
    }
    if o.Answer.Reg != nil {
        ans.raw["reg"] = util.CleanRawXml(o.Answer.Reg.Text)
    }
    if o.Answer.Repg != nil {
        ans.raw["repg"] = util.CleanRawXml(o.Answer.Repg.Text)
    }
    if o.Answer.Rep != nil {
        ans.raw["rep"] = util.CleanRawXml(o.Answer.Rep.Text)
    }
    if o.Answer.Schd != nil {
        ans.raw["schd"] = util.CleanRawXml(o.Answer.Schd.Text)
    }
    if o.Answer.Srv != nil {
        ans.raw["srv"] = util.CleanRawXml(o.Answer.Srv.Text)
    }
    if o.Answer.Srvg != nil {
        ans.raw["srvg"] = util.CleanRawXml(o.Answer.Srvg.Text)
    }
    if o.Answer.Tag != nil {
        ans.raw["tag"] = util.CleanRawXml(o.Answer.Tag.Text)
    }
    if o.Answer.Thr != nil {
        ans.raw["thr"] = util.CleanRawXml(o.Answer.Thr.Text)
    }
    if o.Answer.Tsv != nil {
        ans.raw["tsv"] = util.CleanRawXml(o.Answer.Tsv.Text)
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    Devices *util.VsysEntryType `xml:"devices"`

    Ao *util.RawXml `xml:"address"`
    Ag *util.RawXml `xml:"address-group"`
    App *util.RawXml `xml:"application"`
    Afil *util.RawXml `xml:"application-filter"`
    Appg *util.RawXml `xml:"application-group"`
    As *util.RawXml `xml:"application-status"`
    At *util.RawXml `xml:"application-tag"`
    Aobj *util.RawXml `xml:"authentication-object"`
    Acode *util.RawXml `xml:"authorization-code"`
    Email *util.RawXml `xml:"email-scheduler"`
    Edl *util.RawXml `xml:"external-list"`
    Ls *util.RawXml `xml:"log-settings"`
    Master *util.RawXml `xml:"master-device"`
    Pdf *util.RawXml `xml:"pdf-summary-report"`
    Postrb *util.RawXml `xml:"post-rulebase"`
    Prerb *util.RawXml `xml:"pre-rulebase"`
    Profg *util.RawXml `xml:"profile-group"`
    Profs *util.RawXml `xml:"profiles"`
    Reg *util.RawXml `xml:"region"`
    Repg *util.RawXml `xml:"report-group"`
    Rep *util.RawXml `xml:"reports"`
    Schd *util.RawXml `xml:"schedule"`
    Srv *util.RawXml `xml:"service"`
    Srvg *util.RawXml `xml:"service-group"`
    Tag *util.RawXml `xml:"tag"`
    Thr *util.RawXml `xml:"threats"`
    Tsv *util.RawXml `xml:"to-sw-version"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Devices: util.MapToVsysEnt(e.Devices),
    }

    if t, p := e.raw["ao"]; p {
        ans.Ao = &util.RawXml{t}
    }

    if t, p := e.raw["ag"]; p {
        ans.Ag = &util.RawXml{t}
    }

    if t, p := e.raw["app"]; p {
        ans.App = &util.RawXml{t}
    }

    if t, p := e.raw["afil"]; p {
        ans.Afil = &util.RawXml{t}
    }

    if t, p := e.raw["appg"]; p {
        ans.Appg = &util.RawXml{t}
    }

    if t, p := e.raw["as"]; p {
        ans.As = &util.RawXml{t}
    }

    if t, p := e.raw["at"]; p {
        ans.At = &util.RawXml{t}
    }

    if t, p := e.raw["aobj"]; p {
        ans.Aobj = &util.RawXml{t}
    }

    if t, p := e.raw["acode"]; p {
        ans.Acode = &util.RawXml{t}
    }

    if t, p := e.raw["email"]; p {
        ans.Email = &util.RawXml{t}
    }

    if t, p := e.raw["edl"]; p {
        ans.Edl = &util.RawXml{t}
    }

    if t, p := e.raw["ls"]; p {
        ans.Ls = &util.RawXml{t}
    }

    if t, p := e.raw["master"]; p {
        ans.Master = &util.RawXml{t}
    }

    if t, p := e.raw["pdf"]; p {
        ans.Pdf = &util.RawXml{t}
    }

    if t, p := e.raw["postrb"]; p {
        ans.Postrb = &util.RawXml{t}
    }

    if t, p := e.raw["prerb"]; p {
        ans.Prerb = &util.RawXml{t}
    }

    if t, p := e.raw["profg"]; p {
        ans.Profg = &util.RawXml{t}
    }

    if t, p := e.raw["profs"]; p {
        ans.Profs = &util.RawXml{t}
    }

    if t, p := e.raw["reg"]; p {
        ans.Reg = &util.RawXml{t}
    }

    if t, p := e.raw["repg"]; p {
        ans.Repg = &util.RawXml{t}
    }

    if t, p := e.raw["rep"]; p {
        ans.Rep = &util.RawXml{t}
    }

    if t, p := e.raw["schd"]; p {
        ans.Schd = &util.RawXml{t}
    }

    if t, p := e.raw["srv"]; p {
        ans.Srv = &util.RawXml{t}
    }

    if t, p := e.raw["srvg"]; p {
        ans.Srvg = &util.RawXml{t}
    }

    if t, p := e.raw["tag"]; p {
        ans.Tag = &util.RawXml{t}
    }

    if t, p := e.raw["thr"]; p {
        ans.Thr = &util.RawXml{t}
    }

    if t, p := e.raw["tsv"]; p {
        ans.Tsv = &util.RawXml{t}
    }

    return ans
}
