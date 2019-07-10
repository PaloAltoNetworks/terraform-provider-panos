package logfwd

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a log forwarding profile.
//
// PAN-OS 8.0+.
type Entry struct {
    Name string
    Description string
    EnhancedLogging bool

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.EnhancedLogging = s.EnhancedLogging
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
    }

    if o.Answer.MatchList != nil {
        ans.raw = map[string] string {
            "ml": util.CleanRawXml(o.Answer.MatchList.Text),
        }
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Description: o.Answer.Description,
        EnhancedLogging: util.AsBool(o.Answer.EnhancedLogging),
    }

    if o.Answer.MatchList != nil {
        ans.raw = map[string] string {
            "ml": util.CleanRawXml(o.Answer.MatchList.Text),
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description,omitempty"`
    MatchList *util.RawXml `xml:"match-list"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
    }

    if text := e.raw["ml"]; text != "" {
        ans.MatchList = &util.RawXml{text}
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description,omitempty"`
    EnhancedLogging string `xml:"enhanced-application-logging"`
    MatchList *util.RawXml `xml:"match-list"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Description: e.Description,
        EnhancedLogging: util.YesNo(e.EnhancedLogging),
    }

    if text := e.raw["ml"]; text != "" {
        ans.MatchList = &util.RawXml{text}
    }

    return ans
}
