package action

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a log forwarding profile match list action.
//
// PAN-OS 8.0+.
type Entry struct {
	Name         string
	ActionType   string
	Action       string
	Target       string
	Registration string
	HttpProfile  string
	Tags         []string // ordered
	Timeout      int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.ActionType = s.ActionType
	o.Action = s.Action
	o.Target = s.Target
	o.Registration = s.Registration
	o.HttpProfile = s.HttpProfile
	o.Tags = s.Tags
	o.Timeout = s.Timeout
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
		Name:       o.Answer.Name,
		ActionType: ActionTypeTagging,
		Action:     o.Answer.Type.Tagging.Action,
		Target:     o.Answer.Type.Tagging.Target,
		Tags:       util.MemToStr(o.Answer.Type.Tagging.Tags),
	}

	if o.Answer.Type.Tagging.Reg.Local != nil {
		ans.Registration = RegistrationLocal
	} else if o.Answer.Type.Tagging.Reg.Panorama != nil {
		ans.Registration = RegistrationPanorama
	} else if o.Answer.Type.Tagging.Reg.Remote != nil {
		ans.Registration = RegistrationRemote
		ans.HttpProfile = o.Answer.Type.Tagging.Reg.Remote.HttpProfile
	}

	return ans
}

type container_v2 struct {
	Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	if o.Answer.Type.Tagging != nil {
		ans.ActionType = ActionTypeTagging
		ans.Action = o.Answer.Type.Tagging.Action
		ans.Target = o.Answer.Type.Tagging.Target
		ans.Tags = util.MemToStr(o.Answer.Type.Tagging.Tags)

		if o.Answer.Type.Tagging.Reg.Local != nil {
			ans.Registration = RegistrationLocal
		} else if o.Answer.Type.Tagging.Reg.Panorama != nil {
			ans.Registration = RegistrationPanorama
		} else if o.Answer.Type.Tagging.Reg.Remote != nil {
			ans.Registration = RegistrationRemote
			ans.HttpProfile = o.Answer.Type.Tagging.Reg.Remote.HttpProfile
		}
	} else if o.Answer.Type.Integration != nil {
		ans.ActionType = ActionTypeIntegration
		ans.Action = o.Answer.Type.Integration.Action
	}

	return ans
}

type container_v3 struct {
	Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	if o.Answer.Type.Tagging != nil {
		ans.ActionType = ActionTypeTagging
		ans.Action = o.Answer.Type.Tagging.Action
		ans.Target = o.Answer.Type.Tagging.Target
		ans.Tags = util.MemToStr(o.Answer.Type.Tagging.Tags)
		ans.Timeout = o.Answer.Type.Tagging.Timeout

		if o.Answer.Type.Tagging.Reg.Local != nil {
			ans.Registration = RegistrationLocal
		} else if o.Answer.Type.Tagging.Reg.Panorama != nil {
			ans.Registration = RegistrationPanorama
		} else if o.Answer.Type.Tagging.Reg.Remote != nil {
			ans.Registration = RegistrationRemote
			ans.HttpProfile = o.Answer.Type.Tagging.Reg.Remote.HttpProfile
		}
	} else if o.Answer.Type.Integration != nil {
		ans.ActionType = ActionTypeIntegration
		ans.Action = o.Answer.Type.Integration.Action
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name      `xml:"entry"`
	Name    string        `xml:"name,attr"`
	Type    actionType_v1 `xml:"type"`
}

type actionType_v1 struct {
	Tagging tagging_v1 `xml:"tagging"`
}

type tagging_v1 struct {
	Action string           `xml:"action"`
	Target string           `xml:"target"`
	Reg    reg              `xml:"registration"`
	Tags   *util.MemberType `xml:"tags"`
}

type reg struct {
	Local    *string    `xml:"localhost"`
	Panorama *string    `xml:"panorama"`
	Remote   *regRemote `xml:"remote"`
}

type regRemote struct {
	HttpProfile string `xml:"http-profile"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
		Type: actionType_v1{
			Tagging: tagging_v1{
				Action: e.Action,
				Target: e.Target,
				Tags:   util.StrToMem(e.Tags),
			},
		},
	}

	s := ""
	switch e.Registration {
	case RegistrationLocal:
		ans.Type.Tagging.Reg.Local = &s
	case RegistrationPanorama:
		ans.Type.Tagging.Reg.Panorama = &s
	case RegistrationRemote:
		ans.Type.Tagging.Reg.Remote = &regRemote{e.HttpProfile}
	}

	return ans
}

type entry_v2 struct {
	XMLName xml.Name      `xml:"entry"`
	Name    string        `xml:"name,attr"`
	Type    actionType_v2 `xml:"type"`
}

type actionType_v2 struct {
	Tagging     *tagging_v1  `xml:"tagging"`
	Integration *integration `xml:"integration"`
}

type integration struct {
	Action string `xml:"action"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	switch e.ActionType {
	case ActionTypeTagging:
		ans.Type.Tagging = &tagging_v1{
			Action: e.Action,
			Target: e.Target,
			Tags:   util.StrToMem(e.Tags),
		}

		s := ""
		switch e.Registration {
		case RegistrationLocal:
			ans.Type.Tagging.Reg.Local = &s
		case RegistrationPanorama:
			ans.Type.Tagging.Reg.Panorama = &s
		case RegistrationRemote:
			ans.Type.Tagging.Reg.Remote = &regRemote{e.HttpProfile}
		}
	case ActionTypeIntegration:
		ans.Type.Integration = &integration{
			Action: e.Action,
		}
	}

	return ans
}

type entry_v3 struct {
	XMLName xml.Name      `xml:"entry"`
	Name    string        `xml:"name,attr"`
	Type    actionType_v3 `xml:"type"`
}

type actionType_v3 struct {
	Tagging     *tagging_v2  `xml:"tagging"`
	Integration *integration `xml:"integration"`
}

type tagging_v2 struct {
	Action  string           `xml:"action"`
	Target  string           `xml:"target"`
	Reg     reg              `xml:"registration"`
	Tags    *util.MemberType `xml:"tags"`
	Timeout int              `xml:"timeout,omitempty"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name: e.Name,
	}

	switch e.ActionType {
	case ActionTypeTagging:
		ans.Type.Tagging = &tagging_v2{
			Action:  e.Action,
			Target:  e.Target,
			Tags:    util.StrToMem(e.Tags),
			Timeout: e.Timeout,
		}

		s := ""
		switch e.Registration {
		case RegistrationLocal:
			ans.Type.Tagging.Reg.Local = &s
		case RegistrationPanorama:
			ans.Type.Tagging.Reg.Panorama = &s
		case RegistrationRemote:
			ans.Type.Tagging.Reg.Remote = &regRemote{e.HttpProfile}
		}
	case ActionTypeIntegration:
		ans.Type.Integration = &integration{
			Action: e.Action,
		}
	}

	return ans
}
