package dg

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a virtual
// router.
//
// Devices is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
    Name string
    Description string
    Devices map[string] []string
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

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    Devices *util.VsysEntryType `xml:"devices"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Devices: util.MapToVsysEnt(e.Devices),
    }

    return ans
}
