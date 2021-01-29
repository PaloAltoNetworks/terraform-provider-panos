package pbf

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a policy
// based forwarding rule.
//
// Targets is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
	Name                               string
	Description                        string
	Tags                               []string // ordered
	FromType                           string
	FromValues                         []string // unordered
	SourceAddresses                    []string // unordered
	SourceUsers                        []string // unordered
	NegateSource                       bool
	DestinationAddresses               []string // unordered
	NegateDestination                  bool
	Applications                       []string // unordered
	Services                           []string // unordered
	Schedule                           string
	Disabled                           bool
	Action                             string
	ForwardVsys                        string
	ForwardEgressInterface             string
	ForwardNextHopType                 string
	ForwardNextHopValue                string
	ForwardMonitorProfile              string
	ForwardMonitorIpAddress            string
	ForwardMonitorDisableIfUnreachable bool
	EnableEnforceSymmetricReturn       bool
	SymmetricReturnAddresses           []string // ordered
	ActiveActiveDeviceBinding          string
	Targets                            map[string][]string
	NegateTarget                       bool
	Uuid                               string // 9.0+
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
	o.Targets = s.Targets
	o.NegateTarget = s.NegateTarget
	o.Uuid = s.Uuid
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                      o.Name,
		SourceAddresses:           util.MemToStr(o.SourceAddresses),
		SourceUsers:               util.MemToStr(o.SourceUsers),
		NegateSource:              util.AsBool(o.NegateSource),
		DestinationAddresses:      util.MemToStr(o.DestinationAddresses),
		NegateDestination:         util.AsBool(o.NegateDestination),
		Applications:              util.MemToStr(o.Applications),
		Services:                  util.MemToStr(o.Services),
		Schedule:                  o.Schedule,
		Tags:                      util.MemToStr(o.Tags),
		Disabled:                  util.AsBool(o.Disabled),
		Description:               o.Description,
		ActiveActiveDeviceBinding: o.ActiveActiveDeviceBinding,
	}

	if o.TargetInfo != nil {
		ans.NegateTarget = util.AsBool(o.TargetInfo.NegateTarget)
		ans.Targets = util.VsysEntToMap(o.TargetInfo.Targets)
	}

	switch {
	case o.FromZones != nil:
		ans.FromType = FromTypeZone
		ans.FromValues = util.MemToStr(o.FromZones)
	case o.FromInterfaces != nil:
		ans.FromType = FromTypeInterface
		ans.FromValues = util.MemToStr(o.FromInterfaces)
	}

	switch {
	case o.Action.Forward != nil:
		ans.Action = ActionForward
		ans.ForwardEgressInterface = o.Action.Forward.ForwardEgressInterface

		if o.Action.Forward.NextHop != nil {
			if o.Action.Forward.NextHop.IpAddress != "" {
				ans.ForwardNextHopType = ForwardNextHopTypeIpAddress
				ans.ForwardNextHopValue = o.Action.Forward.NextHop.IpAddress
			}
		}

		if o.Action.Forward.Monitor != nil {
			ans.ForwardMonitorProfile = o.Action.Forward.Monitor.ForwardMonitorProfile
			ans.ForwardMonitorIpAddress = o.Action.Forward.Monitor.ForwardMonitorIpAddress
			ans.ForwardMonitorDisableIfUnreachable = util.AsBool(o.Action.Forward.Monitor.ForwardMonitorDisableIfUnreachable)
		}
	case o.Action.ForwardVsys != nil:
		ans.Action = ActionVsysForward
		ans.ForwardVsys = *o.Action.ForwardVsys
	case o.Action.Discard != nil:
		ans.Action = ActionDiscard
	case o.Action.NoPbf != nil:
		ans.Action = ActionNoPbf
	}

	if o.Symmetric != nil {
		ans.EnableEnforceSymmetricReturn = util.AsBool(o.Symmetric.EnableEnforceSymmetricReturn)
		ans.SymmetricReturnAddresses = util.EntToStr(o.Symmetric.SymmetricReturnAddresses)
	}

	return ans
}

type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                      o.Name,
		SourceAddresses:           util.MemToStr(o.SourceAddresses),
		SourceUsers:               util.MemToStr(o.SourceUsers),
		NegateSource:              util.AsBool(o.NegateSource),
		DestinationAddresses:      util.MemToStr(o.DestinationAddresses),
		NegateDestination:         util.AsBool(o.NegateDestination),
		Applications:              util.MemToStr(o.Applications),
		Services:                  util.MemToStr(o.Services),
		Schedule:                  o.Schedule,
		Tags:                      util.MemToStr(o.Tags),
		Disabled:                  util.AsBool(o.Disabled),
		Description:               o.Description,
		ActiveActiveDeviceBinding: o.ActiveActiveDeviceBinding,
		Uuid:                      o.Uuid,
	}

	if o.TargetInfo != nil {
		ans.NegateTarget = util.AsBool(o.TargetInfo.NegateTarget)
		ans.Targets = util.VsysEntToMap(o.TargetInfo.Targets)
	}

	switch {
	case o.FromZones != nil:
		ans.FromType = FromTypeZone
		ans.FromValues = util.MemToStr(o.FromZones)
	case o.FromInterfaces != nil:
		ans.FromType = FromTypeInterface
		ans.FromValues = util.MemToStr(o.FromInterfaces)
	}

	switch {
	case o.Action.Forward != nil:
		ans.Action = ActionForward
		ans.ForwardEgressInterface = o.Action.Forward.ForwardEgressInterface

		if o.Action.Forward.NextHop != nil {
			if o.Action.Forward.NextHop.IpAddress != "" {
				ans.ForwardNextHopType = ForwardNextHopTypeIpAddress
				ans.ForwardNextHopValue = o.Action.Forward.NextHop.IpAddress
			} else if o.Action.Forward.NextHop.Fqdn != "" {
				ans.ForwardNextHopType = ForwardNextHopTypeFqdn
				ans.ForwardNextHopValue = o.Action.Forward.NextHop.Fqdn
			}
		}

		if o.Action.Forward.Monitor != nil {
			ans.ForwardMonitorProfile = o.Action.Forward.Monitor.ForwardMonitorProfile
			ans.ForwardMonitorIpAddress = o.Action.Forward.Monitor.ForwardMonitorIpAddress
			ans.ForwardMonitorDisableIfUnreachable = util.AsBool(o.Action.Forward.Monitor.ForwardMonitorDisableIfUnreachable)
		}
	case o.Action.ForwardVsys != nil:
		ans.Action = ActionVsysForward
		ans.ForwardVsys = *o.Action.ForwardVsys
	case o.Action.Discard != nil:
		ans.Action = ActionDiscard
	case o.Action.NoPbf != nil:
		ans.Action = ActionNoPbf
	}

	if o.Symmetric != nil {
		ans.EnableEnforceSymmetricReturn = util.AsBool(o.Symmetric.EnableEnforceSymmetricReturn)
		ans.SymmetricReturnAddresses = util.EntToStr(o.Symmetric.SymmetricReturnAddresses)
	}

	return ans
}

type entry_v1 struct {
	XMLName                   xml.Name         `xml:"entry"`
	Name                      string           `xml:"name,attr"`
	FromZones                 *util.MemberType `xml:"from>zone"`
	FromInterfaces            *util.MemberType `xml:"from>interface"`
	SourceAddresses           *util.MemberType `xml:"source"`
	SourceUsers               *util.MemberType `xml:"source-user"`
	NegateSource              string           `xml:"negate-source"`
	DestinationAddresses      *util.MemberType `xml:"destination"`
	NegateDestination         string           `xml:"negate-destination"`
	Applications              *util.MemberType `xml:"application"`
	Services                  *util.MemberType `xml:"service"`
	Schedule                  string           `xml:"schedule,omitempty"`
	Tags                      *util.MemberType `xml:"tag"`
	Disabled                  string           `xml:"disabled"`
	Description               string           `xml:"description,omitempty"`
	Action                    act_v1           `xml:"action"`
	Symmetric                 *sym             `xml:"enforce-symmetric-return"`
	ActiveActiveDeviceBinding string           `xml:"active-active-device-binding,omitempty"`
	TargetInfo                *targetInfo      `xml:"target"`
}

type act_v1 struct {
	Forward     *fwd_v1 `xml:"forward"`
	ForwardVsys *string `xml:"forward-to-vsys"`
	Discard     *string `xml:"discard"`
	NoPbf       *string `xml:"no-pbf"`
}

type fwd_v1 struct {
	ForwardEgressInterface string      `xml:"egress-interface"`
	NextHop                *nextHop_v1 `xml:"nexthop"`
	Monitor                *fwdMonitor `xml:"monitor"`
}

type nextHop_v1 struct {
	IpAddress string `xml:"ip-address,omitempty"`
}

type fwdMonitor struct {
	ForwardMonitorProfile              string `xml:"profile"`
	ForwardMonitorIpAddress            string `xml:"ip-address,omitempty"`
	ForwardMonitorDisableIfUnreachable string `xml:"disable-if-unreachable"`
}

type sym struct {
	EnableEnforceSymmetricReturn string          `xml:"enabled"`
	SymmetricReturnAddresses     *util.EntryType `xml:"nexthop-address-list"`
}

type targetInfo struct {
	Targets      *util.VsysEntryType `xml:"devices"`
	NegateTarget string              `xml:"negate,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                      e.Name,
		SourceAddresses:           util.StrToMem(e.SourceAddresses),
		SourceUsers:               util.StrToMem(e.SourceUsers),
		NegateSource:              util.YesNo(e.NegateSource),
		DestinationAddresses:      util.StrToMem(e.DestinationAddresses),
		NegateDestination:         util.YesNo(e.NegateDestination),
		Applications:              util.StrToMem(e.Applications),
		Services:                  util.StrToMem(e.Services),
		Schedule:                  e.Schedule,
		Tags:                      util.StrToMem(e.Tags),
		Disabled:                  util.YesNo(e.Disabled),
		Description:               e.Description,
		ActiveActiveDeviceBinding: e.ActiveActiveDeviceBinding,
	}

	if e.Targets != nil || e.NegateTarget {
		ans.TargetInfo = &targetInfo{
			Targets:      util.MapToVsysEnt(e.Targets),
			NegateTarget: util.YesNo(e.NegateTarget),
		}
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
				ForwardMonitorProfile:              e.ForwardMonitorProfile,
				ForwardMonitorIpAddress:            e.ForwardMonitorIpAddress,
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
			SymmetricReturnAddresses:     util.StrToEnt(e.SymmetricReturnAddresses),
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName                   xml.Name         `xml:"entry"`
	Name                      string           `xml:"name,attr"`
	Uuid                      string           `xml:"uuid,attr,omitempty"`
	FromZones                 *util.MemberType `xml:"from>zone"`
	FromInterfaces            *util.MemberType `xml:"from>interface"`
	SourceAddresses           *util.MemberType `xml:"source"`
	SourceUsers               *util.MemberType `xml:"source-user"`
	NegateSource              string           `xml:"negate-source"`
	DestinationAddresses      *util.MemberType `xml:"destination"`
	NegateDestination         string           `xml:"negate-destination"`
	Applications              *util.MemberType `xml:"application"`
	Services                  *util.MemberType `xml:"service"`
	Schedule                  string           `xml:"schedule,omitempty"`
	Tags                      *util.MemberType `xml:"tag"`
	Disabled                  string           `xml:"disabled"`
	Description               string           `xml:"description,omitempty"`
	Action                    act_v2           `xml:"action"`
	Symmetric                 *sym             `xml:"enforce-symmetric-return"`
	ActiveActiveDeviceBinding string           `xml:"active-active-device-binding,omitempty"`
	TargetInfo                *targetInfo      `xml:"target"`
}

type act_v2 struct {
	Forward     *fwd_v2 `xml:"forward"`
	ForwardVsys *string `xml:"forward-to-vsys"`
	Discard     *string `xml:"discard"`
	NoPbf       *string `xml:"no-pbf"`
}

type fwd_v2 struct {
	ForwardEgressInterface string      `xml:"egress-interface"`
	NextHop                *nextHop_v2 `xml:"nexthop"`
	Monitor                *fwdMonitor `xml:"monitor"`
}

type nextHop_v2 struct {
	IpAddress string `xml:"ip-address,omitempty"`
	Fqdn      string `xml:"fqdn,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                      e.Name,
		SourceAddresses:           util.StrToMem(e.SourceAddresses),
		SourceUsers:               util.StrToMem(e.SourceUsers),
		NegateSource:              util.YesNo(e.NegateSource),
		DestinationAddresses:      util.StrToMem(e.DestinationAddresses),
		NegateDestination:         util.YesNo(e.NegateDestination),
		Applications:              util.StrToMem(e.Applications),
		Services:                  util.StrToMem(e.Services),
		Schedule:                  e.Schedule,
		Tags:                      util.StrToMem(e.Tags),
		Disabled:                  util.YesNo(e.Disabled),
		Description:               e.Description,
		ActiveActiveDeviceBinding: e.ActiveActiveDeviceBinding,
		Uuid:                      e.Uuid,
	}

	if e.Targets != nil || e.NegateTarget {
		ans.TargetInfo = &targetInfo{
			Targets:      util.MapToVsysEnt(e.Targets),
			NegateTarget: util.YesNo(e.NegateTarget),
		}
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
				ForwardMonitorProfile:              e.ForwardMonitorProfile,
				ForwardMonitorIpAddress:            e.ForwardMonitorIpAddress,
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
			SymmetricReturnAddresses:     util.StrToEnt(e.SymmetricReturnAddresses),
		}
	}

	return ans
}
