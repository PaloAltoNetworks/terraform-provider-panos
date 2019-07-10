package pbf

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a peer.
type Entry struct {
    Name string
    Description string
    Tags []string // ordered
    FromType string
    FromValues []string // unordered
    SourceAddresses []string // unordered
    SourceUsers []string // unordered
    NegateSource bool
    DestinationAddresses []string // unordered
    NegateDestination bool
    Applications []string // unordered
    Services []string // unordered
    Schedule string
    Disabled bool
    Action string
    ForwardVsys string
    ForwardEgressInterface string
    ForwardNextHopType string
    ForwardNextHopValue string
    ForwardMonitorProfile string
    ForwardMonitorIpAddress string
    ForwardMonitorDisableIfUnreachable bool
    EnableEnforceSymmetricReturn bool
    SymmetricReturnAddresses []string // ordered
    ActiveActiveDeviceBinding string
    Uuid string // 9.0+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Tags = s.Tags
    o.FromType = s.FromType
    o.FromValues = s.FromValues
    o.SourceAddresses = s.SourceAddresses
    o.SourceUsers = s.SourceUsers
    o.NegateSource = s.NegateSource
    o.DestinationAddresses = s.DestinationAddresses
    o.NegateDestination = s.NegateDestination
    o.Applications = s.Applications
    o.Services = s.Services
    o.Schedule = s.Schedule
    o.Disabled = s.Disabled
    o.Action = s.Action
    o.ForwardVsys = s.ForwardVsys
    o.ForwardEgressInterface = s.ForwardEgressInterface
    o.ForwardNextHopType = s.ForwardNextHopType
    o.ForwardNextHopValue = s.ForwardNextHopValue
    o.ForwardMonitorProfile = s.ForwardMonitorProfile
    o.ForwardMonitorIpAddress = s.ForwardMonitorIpAddress
    o.ForwardMonitorDisableIfUnreachable = s.ForwardMonitorDisableIfUnreachable
    o.EnableEnforceSymmetricReturn = s.EnableEnforceSymmetricReturn
    o.SymmetricReturnAddresses = s.SymmetricReturnAddresses
    o.ActiveActiveDeviceBinding = s.ActiveActiveDeviceBinding
    o.Uuid = s.Uuid
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
        SourceAddresses: util.MemToStr(o.Answer.SourceAddresses),
        SourceUsers: util.MemToStr(o.Answer.SourceUsers),
        NegateSource: util.AsBool(o.Answer.NegateSource),
        DestinationAddresses: util.MemToStr(o.Answer.DestinationAddresses),
        NegateDestination: util.AsBool(o.Answer.NegateDestination),
        Applications: util.MemToStr(o.Answer.Applications),
        Services: util.MemToStr(o.Answer.Services),
        Schedule: o.Answer.Schedule,
        Tags: util.MemToStr(o.Answer.Tags),
        Disabled: util.AsBool(o.Answer.Disabled),
        Description: o.Answer.Description,
        ActiveActiveDeviceBinding: o.Answer.ActiveActiveDeviceBinding,
    }

    switch {
    case o.Answer.FromZones != nil:
        ans.FromType = FromTypeZone
        ans.FromValues = util.MemToStr(o.Answer.FromZones)
    case o.Answer.FromInterfaces != nil:
        ans.FromType = FromTypeInterface
        ans.FromValues = util.MemToStr(o.Answer.FromInterfaces)
    }

    switch {
    case o.Answer.Action.Forward != nil:
        ans.Action = ActionForward
        ans.ForwardEgressInterface = o.Answer.Action.Forward.ForwardEgressInterface

        if o.Answer.Action.Forward.NextHop != nil {
            if o.Answer.Action.Forward.NextHop.IpAddress != "" {
                ans.ForwardNextHopType = ForwardNextHopTypeIpAddress
                ans.ForwardNextHopValue = o.Answer.Action.Forward.NextHop.IpAddress
            }
        }

        if o.Answer.Action.Forward.Monitor != nil {
            ans.ForwardMonitorProfile = o.Answer.Action.Forward.Monitor.ForwardMonitorProfile
            ans.ForwardMonitorIpAddress = o.Answer.Action.Forward.Monitor.ForwardMonitorIpAddress
            ans.ForwardMonitorDisableIfUnreachable = util.AsBool(o.Answer.Action.Forward.Monitor.ForwardMonitorDisableIfUnreachable)
        }
    case o.Answer.Action.ForwardVsys != nil:
        ans.Action = ActionVsysForward
        ans.ForwardVsys = *o.Answer.Action.ForwardVsys
    case o.Answer.Action.Discard != nil:
        ans.Action = ActionDiscard
    case o.Answer.Action.NoPbf != nil:
        ans.Action = ActionNoPbf
    }

    if o.Answer.Symmetric != nil {
        ans.EnableEnforceSymmetricReturn = util.AsBool(o.Answer.Symmetric.EnableEnforceSymmetricReturn)
        ans.SymmetricReturnAddresses = util.EntToStr(o.Answer.Symmetric.SymmetricReturnAddresses)
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        SourceAddresses: util.MemToStr(o.Answer.SourceAddresses),
        SourceUsers: util.MemToStr(o.Answer.SourceUsers),
        NegateSource: util.AsBool(o.Answer.NegateSource),
        DestinationAddresses: util.MemToStr(o.Answer.DestinationAddresses),
        NegateDestination: util.AsBool(o.Answer.NegateDestination),
        Applications: util.MemToStr(o.Answer.Applications),
        Services: util.MemToStr(o.Answer.Services),
        Schedule: o.Answer.Schedule,
        Tags: util.MemToStr(o.Answer.Tags),
        Disabled: util.AsBool(o.Answer.Disabled),
        Description: o.Answer.Description,
        ActiveActiveDeviceBinding: o.Answer.ActiveActiveDeviceBinding,
        Uuid: o.Answer.Uuid,
    }

    switch {
    case o.Answer.FromZones != nil:
        ans.FromType = FromTypeZone
        ans.FromValues = util.MemToStr(o.Answer.FromZones)
    case o.Answer.FromInterfaces != nil:
        ans.FromType = FromTypeInterface
        ans.FromValues = util.MemToStr(o.Answer.FromInterfaces)
    }

    switch {
    case o.Answer.Action.Forward != nil:
        ans.Action = ActionForward
        ans.ForwardEgressInterface = o.Answer.Action.Forward.ForwardEgressInterface

        if o.Answer.Action.Forward.NextHop != nil {
            if o.Answer.Action.Forward.NextHop.IpAddress != "" {
                ans.ForwardNextHopType = ForwardNextHopTypeIpAddress
                ans.ForwardNextHopValue = o.Answer.Action.Forward.NextHop.IpAddress
            } else if o.Answer.Action.Forward.NextHop.Fqdn != "" {
                ans.ForwardNextHopType = ForwardNextHopTypeFqdn
                ans.ForwardNextHopValue = o.Answer.Action.Forward.NextHop.Fqdn
            }
        }

        if o.Answer.Action.Forward.Monitor != nil {
            ans.ForwardMonitorProfile = o.Answer.Action.Forward.Monitor.ForwardMonitorProfile
            ans.ForwardMonitorIpAddress = o.Answer.Action.Forward.Monitor.ForwardMonitorIpAddress
            ans.ForwardMonitorDisableIfUnreachable = util.AsBool(o.Answer.Action.Forward.Monitor.ForwardMonitorDisableIfUnreachable)
        }
    case o.Answer.Action.ForwardVsys != nil:
        ans.Action = ActionVsysForward
        ans.ForwardVsys = *o.Answer.Action.ForwardVsys
    case o.Answer.Action.Discard != nil:
        ans.Action = ActionDiscard
    case o.Answer.Action.NoPbf != nil:
        ans.Action = ActionNoPbf
    }

    if o.Answer.Symmetric != nil {
        ans.EnableEnforceSymmetricReturn = util.AsBool(o.Answer.Symmetric.EnableEnforceSymmetricReturn)
        ans.SymmetricReturnAddresses = util.EntToStr(o.Answer.Symmetric.SymmetricReturnAddresses)
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    FromZones *util.MemberType `xml:"from>zone"`
    FromInterfaces *util.MemberType `xml:"from>interface"`
    SourceAddresses *util.MemberType `xml:"source"`
    SourceUsers *util.MemberType `xml:"source-user"`
    NegateSource string `xml:"negate-source"`
    DestinationAddresses *util.MemberType `xml:"destination"`
    NegateDestination string `xml:"negate-destination"`
    Applications *util.MemberType `xml:"application"`
    Services *util.MemberType `xml:"service"`
    Schedule string `xml:"schedule,omitempty"`
    Tags *util.MemberType `xml:"tag"`
    Disabled string `xml:"disabled"`
    Description string `xml:"description,omitempty"`
    Action act_v1 `xml:"action"`
    Symmetric *sym `xml:"enforce-symmetric-return"`
    ActiveActiveDeviceBinding string `xml:"active-active-device-binding,omitempty"`
}

type act_v1 struct {
    Forward *fwd_v1 `xml:"forward"`
    ForwardVsys *string `xml:"forward-to-vsys"`
    Discard *string `xml:"discard"`
    NoPbf *string `xml:"no-pbf"`
}

type fwd_v1 struct {
    ForwardEgressInterface string `xml:"egress-interface"`
    NextHop *nextHop_v1 `xml:"nexthop"`
    Monitor *fwdMonitor `xml:"monitor"`
}

type nextHop_v1 struct {
    IpAddress string `xml:"ip-address,omitempty"`
}

type fwdMonitor struct {
    ForwardMonitorProfile string `xml:"profile"`
    ForwardMonitorIpAddress string `xml:"ip-address,omitempty"`
    ForwardMonitorDisableIfUnreachable string `xml:"disable-if-unreachable"`
}

type sym struct {
    EnableEnforceSymmetricReturn string `xml:"enabled"`
    SymmetricReturnAddresses *util.EntryType `xml:"nexthop-address-list"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        SourceAddresses: util.StrToMem(e.SourceAddresses),
        SourceUsers: util.StrToMem(e.SourceUsers),
        NegateSource: util.YesNo(e.NegateSource),
        DestinationAddresses: util.StrToMem(e.DestinationAddresses),
        NegateDestination: util.YesNo(e.NegateDestination),
        Applications: util.StrToMem(e.Applications),
        Services: util.StrToMem(e.Services),
        Schedule: e.Schedule,
        Tags: util.StrToMem(e.Tags),
        Disabled: util.YesNo(e.Disabled),
        Description: e.Description,
        ActiveActiveDeviceBinding: e.ActiveActiveDeviceBinding,
    }

    switch e.FromType {
    case FromTypeZone:
        ans.FromZones = util.StrToMem(e.FromValues)
    case FromTypeInterface:
        ans.FromInterfaces = util.StrToMem(e.FromValues)
    }

    switch e.Action {
    case ActionForward:
        ans.Action.Forward = &fwd_v1{
            ForwardEgressInterface: e.ForwardEgressInterface,
        }

        switch e.ForwardNextHopType {
        case ForwardNextHopTypeIpAddress:
            ans.Action.Forward.NextHop = &nextHop_v1{
                IpAddress: e.ForwardNextHopValue,
            }
        }

        if e.ForwardMonitorProfile != "" || e.ForwardMonitorIpAddress != "" || e.ForwardMonitorDisableIfUnreachable {
            ans.Action.Forward.Monitor = &fwdMonitor{
                ForwardMonitorProfile: e.ForwardMonitorProfile,
                ForwardMonitorIpAddress: e.ForwardMonitorIpAddress,
                ForwardMonitorDisableIfUnreachable: util.YesNo(e.ForwardMonitorDisableIfUnreachable),
            }
        }
    case ActionVsysForward:
        ans.Action.ForwardVsys = &e.ForwardVsys
    case ActionDiscard:
        s := ""
        ans.Action.Discard = &s
    case ActionNoPbf:
        s := ""
        ans.Action.NoPbf = &s
    }

    if e.EnableEnforceSymmetricReturn || len(e.SymmetricReturnAddresses) > 0 {
        ans.Symmetric = &sym{
            EnableEnforceSymmetricReturn: util.YesNo(e.EnableEnforceSymmetricReturn),
            SymmetricReturnAddresses: util.StrToEnt(e.SymmetricReturnAddresses),
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Uuid string `xml:"uuid,attr,omitempty"`
    FromZones *util.MemberType `xml:"from>zone"`
    FromInterfaces *util.MemberType `xml:"from>interface"`
    SourceAddresses *util.MemberType `xml:"source"`
    SourceUsers *util.MemberType `xml:"source-user"`
    NegateSource string `xml:"negate-source"`
    DestinationAddresses *util.MemberType `xml:"destination"`
    NegateDestination string `xml:"negate-destination"`
    Applications *util.MemberType `xml:"application"`
    Services *util.MemberType `xml:"service"`
    Schedule string `xml:"schedule,omitempty"`
    Tags *util.MemberType `xml:"tag"`
    Disabled string `xml:"disabled"`
    Description string `xml:"description,omitempty"`
    Action act_v2 `xml:"action"`
    Symmetric *sym `xml:"enforce-symmetric-return"`
    ActiveActiveDeviceBinding string `xml:"active-active-device-binding,omitempty"`
}

type act_v2 struct {
    Forward *fwd_v2 `xml:"forward"`
    ForwardVsys *string `xml:"forward-to-vsys"`
    Discard *string `xml:"discard"`
    NoPbf *string `xml:"no-pbf"`
}

type fwd_v2 struct {
    ForwardEgressInterface string `xml:"egress-interface"`
    NextHop *nextHop_v2 `xml:"nexthop"`
    Monitor *fwdMonitor `xml:"monitor"`
}

type nextHop_v2 struct {
    IpAddress string `xml:"ip-address,omitempty"`
    Fqdn string `xml:"fqdn,omitempty"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        SourceAddresses: util.StrToMem(e.SourceAddresses),
        SourceUsers: util.StrToMem(e.SourceUsers),
        NegateSource: util.YesNo(e.NegateSource),
        DestinationAddresses: util.StrToMem(e.DestinationAddresses),
        NegateDestination: util.YesNo(e.NegateDestination),
        Applications: util.StrToMem(e.Applications),
        Services: util.StrToMem(e.Services),
        Schedule: e.Schedule,
        Tags: util.StrToMem(e.Tags),
        Disabled: util.YesNo(e.Disabled),
        Description: e.Description,
        ActiveActiveDeviceBinding: e.ActiveActiveDeviceBinding,
        Uuid: e.Uuid,
    }

    switch e.FromType {
    case FromTypeZone:
        ans.FromZones = util.StrToMem(e.FromValues)
    case FromTypeInterface:
        ans.FromInterfaces = util.StrToMem(e.FromValues)
    }

    switch e.Action {
    case ActionForward:
        ans.Action.Forward = &fwd_v2{
            ForwardEgressInterface: e.ForwardEgressInterface,
        }

        switch e.ForwardNextHopType {
        case ForwardNextHopTypeIpAddress:
            ans.Action.Forward.NextHop = &nextHop_v2{
                IpAddress: e.ForwardNextHopValue,
            }
        case ForwardNextHopTypeFqdn:
            ans.Action.Forward.NextHop = &nextHop_v2{
                Fqdn: e.ForwardNextHopValue,
            }
        }

        if e.ForwardMonitorProfile != "" || e.ForwardMonitorIpAddress != "" || e.ForwardMonitorDisableIfUnreachable {
            ans.Action.Forward.Monitor = &fwdMonitor{
                ForwardMonitorProfile: e.ForwardMonitorProfile,
                ForwardMonitorIpAddress: e.ForwardMonitorIpAddress,
                ForwardMonitorDisableIfUnreachable: util.YesNo(e.ForwardMonitorDisableIfUnreachable),
            }
        }
    case ActionVsysForward:
        ans.Action.ForwardVsys = &e.ForwardVsys
    case ActionDiscard:
        s := ""
        ans.Action.Discard = &s
    case ActionNoPbf:
        s := ""
        ans.Action.NoPbf = &s
    }

    if e.EnableEnforceSymmetricReturn || len(e.SymmetricReturnAddresses) > 0 {
        ans.Symmetric = &sym{
            EnableEnforceSymmetricReturn: util.YesNo(e.EnableEnforceSymmetricReturn),
            SymmetricReturnAddresses: util.StrToEnt(e.SymmetricReturnAddresses),
        }
    }

    return ans
}
