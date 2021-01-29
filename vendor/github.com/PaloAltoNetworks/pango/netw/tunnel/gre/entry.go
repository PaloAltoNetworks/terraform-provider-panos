package gre

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a peer.
type Entry struct {
	Name               string
	Interface          string
	LocalAddressType   string
	LocalAddressValue  string
	PeerAddress        string
	TunnelInterface    string
	Ttl                int
	CopyTos            bool
	EnableKeepAlive    bool
	KeepAliveInterval  int
	KeepAliveRetry     int
	KeepAliveHoldTimer int
	Disabled           bool
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Interface = s.Interface
	o.LocalAddressType = s.LocalAddressType
	o.LocalAddressValue = s.LocalAddressValue
	o.PeerAddress = s.PeerAddress
	o.TunnelInterface = s.TunnelInterface
	o.Ttl = s.Ttl
	o.CopyTos = s.CopyTos
	o.EnableKeepAlive = s.EnableKeepAlive
	o.KeepAliveInterval = s.KeepAliveInterval
	o.KeepAliveRetry = s.KeepAliveRetry
	o.KeepAliveHoldTimer = s.KeepAliveHoldTimer
	o.Disabled = s.Disabled
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
		Name:            o.Answer.Name,
		Interface:       o.Answer.Local.Interface,
		PeerAddress:     o.Answer.Peer.PeerAddress,
		TunnelInterface: o.Answer.TunnelInterface,
		Ttl:             o.Answer.Ttl,
		CopyTos:         util.AsBool(o.Answer.CopyTos),
		Disabled:        util.AsBool(o.Answer.Disabled),
	}

	if o.Answer.Local.Ip != "" {
		ans.LocalAddressType = LocalAddressTypeIp
		ans.LocalAddressValue = o.Answer.Local.Ip
	} else if o.Answer.Local.FloatingIp != "" {
		ans.LocalAddressType = LocalAddressTypeFloatingIp
		ans.LocalAddressValue = o.Answer.Local.FloatingIp
	}

	if o.Answer.KeepAlive != nil {
		ans.EnableKeepAlive = util.AsBool(o.Answer.KeepAlive.EnableKeepAlive)
		ans.KeepAliveInterval = o.Answer.KeepAlive.KeepAliveInterval
		ans.KeepAliveRetry = o.Answer.KeepAlive.KeepAliveRetry
		ans.KeepAliveHoldTimer = o.Answer.KeepAlive.KeepAliveHoldTimer
	}

	return ans
}

type entry_v1 struct {
	XMLName         xml.Name     `xml:"entry"`
	Name            string       `xml:"name,attr"`
	Local           localAddress `xml:"local-address"`
	Peer            peerAddress  `xml:"peer-address"`
	TunnelInterface string       `xml:"tunnel-interface"`
	Ttl             int          `xml:"ttl,omitempty"`
	CopyTos         string       `xml:"copy-tos"`
	KeepAlive       *ka          `xml:"keep-alive"`
	Disabled        string       `xml:"disabled"`
}

type localAddress struct {
	Interface  string `xml:"interface"`
	Ip         string `xml:"ip,omitempty"`
	FloatingIp string `xml:"floating-ip,omitempty"`
}

type peerAddress struct {
	PeerAddress string `xml:"ip"`
}

type ka struct {
	EnableKeepAlive    string `xml:"enable"`
	KeepAliveInterval  int    `xml:"interval,omitempty"`
	KeepAliveRetry     int    `xml:"retry,omitempty"`
	KeepAliveHoldTimer int    `xml:"hold-timer,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
		Local: localAddress{
			Interface: e.Interface,
		},
		Peer: peerAddress{
			PeerAddress: e.PeerAddress,
		},
		TunnelInterface: e.TunnelInterface,
		Ttl:             e.Ttl,
		CopyTos:         util.YesNo(e.CopyTos),
		Disabled:        util.YesNo(e.Disabled),
	}

	switch e.LocalAddressType {
	case LocalAddressTypeIp:
		ans.Local.Ip = e.LocalAddressValue
	case LocalAddressTypeFloatingIp:
		ans.Local.FloatingIp = e.LocalAddressValue
	}

	if e.EnableKeepAlive || e.KeepAliveInterval != 0 || e.KeepAliveRetry != 0 || e.KeepAliveHoldTimer != 0 {
		ans.KeepAlive = &ka{
			EnableKeepAlive:    util.YesNo(e.EnableKeepAlive),
			KeepAliveInterval:  e.KeepAliveInterval,
			KeepAliveRetry:     e.KeepAliveRetry,
			KeepAliveHoldTimer: e.KeepAliveHoldTimer,
		}
	}

	return ans
}
