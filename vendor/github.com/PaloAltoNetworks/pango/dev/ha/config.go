package ha

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of a device
// HA configuration.
type Config struct {
	Enable                 bool
	GroupId                int
	Description            string
	Mode                   string
	PeerHa1IpAddress       string
	BackupPeerHa1IpAddress string
	ConfigSyncEnable       bool

	Ha1       *Ha1Interface
	Ha1Backup *Ha1BackupInterface
	Ha2       *Ha2Interface
	Ha2Backup *Ha2Interface
	Ha3       *Ha3Interface
	Ha4       *Ha4Interface
	Ha4Backup *Ha4Interface

	// active passive
	ApPassiveLinkState        string
	ApMonitorFailHoldDownTime int

	// active active
	AaDeviceId string
	// HA3 packet forwarding
	AaTentativeHoldTime     string
	AaSyncVirtualRouter     bool
	AaSyncQos               bool
	AaSessionOwnerSelection string // primary-device, first-packet
	// session owner selection: first-packet
	AaFpSessionSetup string // primary-device, first-packet, ip-modulo, ip-hash
	// session setup: ip-hash
	AaFpSessionSetupIpHashKey  string
	AaFpSessionSetupIpHashSeed string

	ElectionDevicePriority  string
	ElectionPreemptive      bool
	ElectionHeartBeatBackup bool

	ElectionTimersMode                          string
	ElectionTimersAdvPromotionHoldTime          string
	ElectionTimersAdvHelloInterval              int
	ElectionTimersAdvHeartBeatInterval          int
	ElectionTimersAdvFlapMax                    string
	ElectionTimersAdvPreemptionHoldTime         int
	ElectionTimersAdvMonitorFailHoldUpTime      string
	ElectionTimersAdvAdditionalMasterHoldUpTime string

	Ha2StateSyncEnable             bool
	Ha2StateSyncTransport          string
	Ha2StateSyncKeepAliveEnable    bool
	Ha2StateSyncKeepAliveAction    string
	Ha2StateSyncKeepAliveThreshold int

	LinkMonitorEnable           bool
	LinkMonitorFailureCondition string

	raw map[string]string
}

type Ha1Interface struct {
	Port             string
	IpAddress        string
	Netmask          string
	Gateway          string
	EncryptionEnable bool
	MonitorHoldTime  int
}

type Ha1BackupInterface struct {
	Port      string
	IpAddress string
	Netmask   string
	Gateway   string
}

type Ha2Interface struct {
	Port      string
	IpAddress string
	Netmask   string
	Gateway   string
}

type Ha3Interface struct {
	Port string
}

type Ha4Interface struct {
	Port      string
	IpAddress string
	Netmask   string
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.Enable = s.Enable
	o.GroupId = s.GroupId
	o.Description = s.Description
	o.Mode = s.Mode
	o.PeerHa1IpAddress = s.PeerHa1IpAddress
	o.BackupPeerHa1IpAddress = s.BackupPeerHa1IpAddress
	o.ConfigSyncEnable = s.ConfigSyncEnable

	if s.Ha1 != nil {
		o.Ha1 = &Ha1Interface{
			Port:             s.Ha1.Port,
			IpAddress:        s.Ha1.IpAddress,
			Netmask:          s.Ha1.Netmask,
			Gateway:          s.Ha1.Gateway,
			EncryptionEnable: s.Ha1.EncryptionEnable,
			MonitorHoldTime:  s.Ha1.MonitorHoldTime,
		}
	} else {
		o.Ha1 = nil
	}
	if s.Ha1Backup != nil {
		o.Ha1Backup = &Ha1BackupInterface{
			Port:      s.Ha1Backup.Port,
			IpAddress: s.Ha1Backup.IpAddress,
			Netmask:   s.Ha1Backup.Netmask,
			Gateway:   s.Ha1Backup.Gateway,
		}
	} else {
		o.Ha1Backup = nil
	}
	if s.Ha2 != nil {
		o.Ha2 = &Ha2Interface{
			Port:      s.Ha2.Port,
			IpAddress: s.Ha2.IpAddress,
			Netmask:   s.Ha2.Netmask,
			Gateway:   s.Ha2.Gateway,
		}
	} else {
		o.Ha2 = nil
	}
	if s.Ha2Backup != nil {
		o.Ha2Backup = &Ha2Interface{
			Port:      s.Ha2Backup.Port,
			IpAddress: s.Ha2Backup.IpAddress,
			Netmask:   s.Ha2Backup.Netmask,
			Gateway:   s.Ha2Backup.Gateway,
		}
	} else {
		o.Ha2Backup = nil
	}
	if s.Ha3 != nil {
		o.Ha3 = &Ha3Interface{
			Port: s.Ha3.Port,
		}
	} else {
		o.Ha3 = nil
	}
	if s.Ha4 != nil {
		o.Ha4 = &Ha4Interface{
			Port:      s.Ha4.Port,
			IpAddress: s.Ha4.IpAddress,
			Netmask:   s.Ha4.Netmask,
		}
	} else {
		o.Ha4 = nil
	}
	if s.Ha4Backup != nil {
		o.Ha4Backup = &Ha4Interface{
			Port:      s.Ha4Backup.Port,
			IpAddress: s.Ha4Backup.IpAddress,
			Netmask:   s.Ha4Backup.Netmask,
		}
	} else {
		o.Ha4Backup = nil
	}

	o.ApPassiveLinkState = s.ApPassiveLinkState
	o.ApMonitorFailHoldDownTime = s.ApMonitorFailHoldDownTime

	o.AaDeviceId = s.AaDeviceId
	o.AaTentativeHoldTime = s.AaTentativeHoldTime
	o.AaSyncVirtualRouter = s.AaSyncVirtualRouter
	o.AaSyncQos = s.AaSyncQos
	o.AaSessionOwnerSelection = s.AaSessionOwnerSelection
	o.AaFpSessionSetup = s.AaFpSessionSetup
	o.AaFpSessionSetupIpHashKey = s.AaFpSessionSetupIpHashKey
	o.AaFpSessionSetupIpHashSeed = s.AaFpSessionSetupIpHashSeed

	o.ElectionDevicePriority = s.ElectionDevicePriority
	o.ElectionPreemptive = s.ElectionPreemptive
	o.ElectionHeartBeatBackup = s.ElectionHeartBeatBackup

	o.ElectionTimersMode = s.ElectionTimersMode
	o.ElectionTimersAdvPromotionHoldTime = s.ElectionTimersAdvPromotionHoldTime
	o.ElectionTimersAdvHelloInterval = s.ElectionTimersAdvHelloInterval
	o.ElectionTimersAdvHeartBeatInterval = s.ElectionTimersAdvHeartBeatInterval
	o.ElectionTimersAdvFlapMax = s.ElectionTimersAdvFlapMax
	o.ElectionTimersAdvPreemptionHoldTime = s.ElectionTimersAdvPreemptionHoldTime
	o.ElectionTimersAdvMonitorFailHoldUpTime = s.ElectionTimersAdvMonitorFailHoldUpTime
	o.ElectionTimersAdvAdditionalMasterHoldUpTime = s.ElectionTimersAdvAdditionalMasterHoldUpTime

	o.Ha2StateSyncEnable = s.Ha2StateSyncEnable
	o.Ha2StateSyncTransport = s.Ha2StateSyncTransport
	o.Ha2StateSyncKeepAliveEnable = s.Ha2StateSyncKeepAliveEnable
	o.Ha2StateSyncKeepAliveAction = s.Ha2StateSyncKeepAliveAction
	o.Ha2StateSyncKeepAliveThreshold = s.Ha2StateSyncKeepAliveThreshold

	o.LinkMonitorEnable = s.LinkMonitorEnable
	o.LinkMonitorFailureCondition = s.LinkMonitorFailureCondition
}

/** Structs / functions for this namespace. **/

func (o Config) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return "", fn(o)
}

type normalizer interface {
	Normalize() []Config
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"high-availability"`
}

func (o *container_v1) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	return nil
}

func (o *entry_v1) normalize() Config {
	ans := Config{
		Enable:                 util.AsBool(o.Enable),
		GroupId:                o.GroupId,
		Description:            o.Description,
		PeerHa1IpAddress:       o.PeerHa1IpAddress,
		BackupPeerHa1IpAddress: o.BackupPeerHa1IpAddress,
		ConfigSyncEnable:       util.AsBool(o.ConfigSyncEnable),
	}

	raw := make(map[string]string)

	if o.Mode != nil {
		switch {
		case o.Mode.ActivePassive != nil:
			ans.Mode = ModeActivePassive
			ans.ApPassiveLinkState = o.Mode.ActivePassive.PassiveLinkState
			ans.ApMonitorFailHoldDownTime = o.Mode.ActivePassive.MonitorFailHoldDownTime
		case o.Mode.ActiveActive != nil:
			ans.Mode = ModeActiveActive
			ans.AaDeviceId = o.Mode.ActiveActive.DeviceId
			ans.AaTentativeHoldTime = o.Mode.ActiveActive.TentativeHoldTime
			if o.Mode.ActiveActive.NetConfSync != nil {
				x := o.Mode.ActiveActive.NetConfSync
				ans.AaSyncVirtualRouter = util.AsBool(x.VirtualRouter)
				ans.AaSyncQos = util.AsBool(x.Qos)
			}
			if o.Mode.ActiveActive.SessionOwnerSelection != nil {
				x := o.Mode.ActiveActive.SessionOwnerSelection
				if x.PrimaryDevice != nil {
					ans.AaSessionOwnerSelection = AaSessionOwnerSelectionPrimaryDevice
				} else if x.FirstPacket != nil {
					ans.AaSessionOwnerSelection = AaSessionOwnerSelectionFirstPacket
					switch {
					case x.FirstPacket.PrimaryDevice != nil:
						ans.AaFpSessionSetup = AaFpSessionSetupPrimaryDevice
					case x.FirstPacket.FirstPacket != nil:
						ans.AaFpSessionSetup = AaFpSessionSetupFirstPacket
					case x.FirstPacket.IpModulo != nil:
						ans.AaFpSessionSetup = AaFpSessionSetupIpModulo
					case x.FirstPacket.IpHash != nil:
						ans.AaFpSessionSetup = AaFpSessionSetupIpHash
						ans.AaFpSessionSetupIpHashKey =
							x.FirstPacket.IpHash.HashKey
						ans.AaFpSessionSetupIpHashSeed =
							x.FirstPacket.IpHash.HashSeed
					}
				}
			}
			if o.Mode.ActiveActive.VirtualAddress != nil {
				raw["vaddr"] = util.CleanRawXml(o.Mode.ActiveActive.VirtualAddress.Text)
			}
		}
	}

	if o.Interfaces != nil {
		x := o.Interfaces
		if x.Ha1 != nil {
			ans.Ha1 = &Ha1Interface{
				Port:             x.Ha1.Port,
				IpAddress:        x.Ha1.IpAddress,
				Netmask:          x.Ha1.Netmask,
				Gateway:          x.Ha1.Gateway,
				EncryptionEnable: util.AsBool(x.Ha1.EncryptionEnable),
				MonitorHoldTime:  x.Ha1.MonitorHoldTime,
			}
		}
		if x.Ha1Backup != nil {
			ans.Ha1Backup = &Ha1BackupInterface{
				Port:      x.Ha1Backup.Port,
				IpAddress: x.Ha1Backup.IpAddress,
				Netmask:   x.Ha1Backup.Netmask,
				Gateway:   x.Ha1Backup.Gateway,
			}
		}
		if x.Ha2 != nil {
			ans.Ha2 = &Ha2Interface{
				Port:      x.Ha2.Port,
				IpAddress: x.Ha2.IpAddress,
				Netmask:   x.Ha2.Netmask,
				Gateway:   x.Ha2.Gateway,
			}
		}
		if x.Ha2Backup != nil {
			ans.Ha2Backup = &Ha2Interface{
				Port:      x.Ha2Backup.Port,
				IpAddress: x.Ha2Backup.IpAddress,
				Netmask:   x.Ha2Backup.Netmask,
				Gateway:   x.Ha2Backup.Gateway,
			}
		}
		if x.Ha3 != nil {
			ans.Ha3 = &Ha3Interface{
				Port: x.Ha3.Port,
			}
		}
		if x.Ha4 != nil {
			ans.Ha4 = &Ha4Interface{
				Port:      x.Ha4.Port,
				IpAddress: x.Ha4.IpAddress,
				Netmask:   x.Ha4.Netmask,
			}
		}
		if x.Ha4Backup != nil {
			ans.Ha4Backup = &Ha4Interface{
				Port:      x.Ha4Backup.Port,
				IpAddress: x.Ha4Backup.IpAddress,
				Netmask:   x.Ha4Backup.Netmask,
			}
		}
	}

	if o.ElectionOption != nil {
		ans.ElectionDevicePriority = o.ElectionOption.DevicePriority
		ans.ElectionPreemptive = util.AsBool(o.ElectionOption.Preemptive)
		ans.ElectionHeartBeatBackup = util.AsBool(o.ElectionOption.HeartBeatBackup)

		if o.ElectionOption.Timers != nil {
			switch {
			case o.ElectionOption.Timers.Recommended != nil:
				ans.ElectionTimersMode = ElectionTimersModeRecommended
			case o.ElectionOption.Timers.Aggressive != nil:
				ans.ElectionTimersMode = ElectionTimersModeAggressive
			case o.ElectionOption.Timers.Advanced != nil:
				ans.ElectionTimersMode = ElectionTimersModeAdvanced
				x := o.ElectionOption.Timers.Advanced
				ans.ElectionTimersAdvPromotionHoldTime = x.PromotionHoldTime
				ans.ElectionTimersAdvHelloInterval = x.HelloInterval
				ans.ElectionTimersAdvHeartBeatInterval = x.HeartBeatInterval
				ans.ElectionTimersAdvFlapMax = x.FlapMax
				ans.ElectionTimersAdvPreemptionHoldTime = x.PreemptionHoldTime
				ans.ElectionTimersAdvMonitorFailHoldUpTime = x.MonitorFailHoldUpTime
				ans.ElectionTimersAdvAdditionalMasterHoldUpTime = x.AdditionalMasterHoldUpTime
			}
		}
	}

	if o.StateSync != nil {
		ans.Ha2StateSyncEnable = util.AsBool(o.StateSync.Enable)
		ans.Ha2StateSyncTransport = o.StateSync.Transport
		if o.StateSync.Ha2KeepAlive != nil {
			ans.Ha2StateSyncKeepAliveEnable = util.AsBool(o.StateSync.Ha2KeepAlive.Enable)
			ans.Ha2StateSyncKeepAliveAction = o.StateSync.Ha2KeepAlive.Action
			ans.Ha2StateSyncKeepAliveThreshold = o.StateSync.Ha2KeepAlive.Threshold
		}
	}

	if o.LinkMonitor != nil {
		ans.LinkMonitorEnable = util.AsBool(o.LinkMonitor.Enable)
		ans.LinkMonitorFailureCondition = o.LinkMonitor.FailureCondition
		if o.LinkMonitor.LinkGroup != nil {
			raw["linkgroup"] = util.CleanRawXml(o.LinkMonitor.LinkGroup.Text)
		}
	}

	if o.Cluster != nil {
		raw["cluster"] = util.CleanRawXml(o.Cluster.Text)
	}
	if o.PathMonitor != nil {
		raw["pathmonitor"] = util.CleanRawXml(o.PathMonitor.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName                xml.Name        `xml:"high-availability"`
	Enable                 string          `xml:"enabled"`
	Interfaces             *interfaces     `xml:"interface"`
	GroupId                int             `xml:"group>group-id,omitempty"`
	Description            string          `xml:"group>description,omitempty"`
	Mode                   *mode           `xml:"group>mode"`
	PeerHa1IpAddress       string          `xml:"group>peer-ip,omitempty"`
	BackupPeerHa1IpAddress string          `xml:"group>peer-ip-backup,omitempty"`
	ConfigSyncEnable       string          `xml:"group>configuration-synchronization>enabled"`
	ElectionOption         *electionOption `xml:"group>election-option"`
	StateSync              *stateSync      `xml:"group>state-synchronization"`
	LinkMonitor            *linkMonitor    `xml:"group>monitoring>link-monitoring"`

	PathMonitor *util.RawXml `xml:"group>monitoring>path-monitoring"`
	Cluster     *util.RawXml `xml:"cluster"`
}

func (e *entry_v1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type localEntry_v1 entry_v1
	ans := localEntry_v1{
		StateSync: &stateSync{
			Enable: util.YesNo(true),
		},
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v1(ans)
	return nil
}

type interfaces struct {
	Ha1       *ha1Interface       `xml:"ha1"`
	Ha1Backup *ha1BackupInterface `xml:"ha1-backup"`
	Ha2       *ha2Interface       `xml:"ha2"`
	Ha2Backup *ha2Interface       `xml:"ha2-backup"`
	Ha3       *ha3Interface       `xml:"ha3"`
	Ha4       *ha4Interface       `xml:"ha4"`
	Ha4Backup *ha4Interface       `xml:"ha4-backup"`
}

type ha1Interface struct {
	Port             string `xml:"port,omitempty"`
	IpAddress        string `xml:"ip-address,omitempty"`
	Netmask          string `xml:"netmask,omitempty"`
	Gateway          string `xml:"gateway,omitempty"`
	EncryptionEnable string `xml:"encryption>enabled"`
	MonitorHoldTime  int    `xml:"monitor-hold-time,omitempty"`
}

type ha1BackupInterface struct {
	Port      string `xml:"port,omitempty"`
	IpAddress string `xml:"ip-address,omitempty"`
	Netmask   string `xml:"netmask,omitempty"`
	Gateway   string `xml:"gateway,omitempty"`
}

type ha2Interface struct {
	Port      string `xml:"port,omitempty"`
	IpAddress string `xml:"ip-address,omitempty"`
	Netmask   string `xml:"netmask,omitempty"`
	Gateway   string `xml:"gateway,omitempty"`
}

type ha3Interface struct {
	Port string `xml:"port,omitempty"`
}

type ha4Interface struct {
	Port      string `xml:"port,omitempty"`
	IpAddress string `xml:"ip-address,omitempty"`
	Netmask   string `xml:"netmask,omitempty"`
}

type electionOption struct {
	DevicePriority  string `xml:"device-priority,omitempty"`
	Preemptive      string `xml:"preemptive"`
	HeartBeatBackup string `xml:"heartbeat-backup"`

	Timers *optionTimers `xml:"timers"`
}

type optionTimers struct {
	Recommended *string   `xml:"recommended"`
	Aggressive  *string   `xml:"aggressive"`
	Advanced    *advanced `xml:"advanced"`
}

type advanced struct {
	PromotionHoldTime          string `xml:"promotion-hold-time,omitempty"`
	HelloInterval              int    `xml:"hello-interval,omitempty"`
	HeartBeatInterval          int    `xml:"heartbeat-interval,omitempty"`
	FlapMax                    string `xml:"flap-max,omitempty"`
	PreemptionHoldTime         int    `xml:"preemption-hold-time,omitempty"`
	MonitorFailHoldUpTime      string `xml:"monitor-fail-hold-up-time,omitempty"`
	AdditionalMasterHoldUpTime string `xml:"additional-master-hold-up-time,omitempty"`
}

type stateSync struct {
	Enable       string        `xml:"enabled"`
	Transport    string        `xml:"transport,omitempty"`
	Ha2KeepAlive *ha2KeepAlive `xml:"ha2-keep-alive"`
}

type ha2KeepAlive struct {
	Enable    string `xml:"enabled"`
	Action    string `xml:"action,omitempty"`
	Threshold int    `xml:"threshold,omitempty"`
}

type mode struct {
	ActivePassive *activePassive `xml:"active-passive"`
	ActiveActive  *activeActive  `xml:"active-active"`
}

type activePassive struct {
	PassiveLinkState        string `xml:"passive-link-state,omitempty"`
	MonitorFailHoldDownTime int    `xml:"monitor-fail-hold-down-time,omitempty"`
}

type activeActive struct {
	DeviceId              string                 `xml:"device-id"`
	TentativeHoldTime     string                 `xml:"tentative-hold-time,omitempty"`
	NetConfSync           *netConfSync           `xml:"network-configuration>sync"`
	SessionOwnerSelection *sessionOwnerSelection `xml:"session-owner-selection"`
	VirtualAddress        *util.RawXml           `xml:"virtual-address"`
}

type netConfSync struct {
	VirtualRouter string `xml:"virtual-router"`
	Qos           string `xml:"qos"`
}

type sessionOwnerSelection struct {
	PrimaryDevice *string             `xml:"primary-device"`
	FirstPacket   *firstPacketSession `xml:"first-packet>session-setup"`
}

type firstPacketSession struct {
	PrimaryDevice *string `xml:"primary-device"`
	FirstPacket   *string `xml:"first-packet"`
	IpModulo      *string `xml:"ip-modulo"`
	IpHash        *ipHash `xml:"ip-hash"`
}

type ipHash struct {
	HashKey  string `xml:"hash-key,omitempty"`
	HashSeed string `xml:"hash-seed,omitempty"`
}

type linkMonitor struct {
	Enable           string       `xml:"enabled"`
	FailureCondition string       `xml:"failure-condition,omitempty"`
	LinkGroup        *util.RawXml `xml:"link-group"`
}

func specify_v1(e Config) interface{} {
	ans := entry_v1{
		Enable:                 util.YesNo(e.Enable),
		GroupId:                e.GroupId,
		Description:            e.Description,
		PeerHa1IpAddress:       e.PeerHa1IpAddress,
		BackupPeerHa1IpAddress: e.BackupPeerHa1IpAddress,
		ConfigSyncEnable:       util.YesNo(e.ConfigSyncEnable),
	}

	s := ""

	if e.Mode != "" {
		ans.Mode = &mode{}
		switch e.Mode {
		case ModeActivePassive:
			ans.Mode.ActivePassive = &activePassive{
				PassiveLinkState:        e.ApPassiveLinkState,
				MonitorFailHoldDownTime: e.ApMonitorFailHoldDownTime,
			}
		case ModeActiveActive:
			ans.Mode.ActiveActive = &activeActive{
				DeviceId:          e.AaDeviceId,
				TentativeHoldTime: e.AaTentativeHoldTime,
			}
			if e.AaSyncVirtualRouter || e.AaSyncQos {
				ans.Mode.ActiveActive.NetConfSync = &netConfSync{
					VirtualRouter: util.YesNo(e.AaSyncVirtualRouter),
					Qos:           util.YesNo(e.AaSyncQos),
				}
			}
			switch e.AaSessionOwnerSelection {
			case AaSessionOwnerSelectionPrimaryDevice:
				ans.Mode.ActiveActive.SessionOwnerSelection = &sessionOwnerSelection{
					PrimaryDevice: &s,
				}
			case AaSessionOwnerSelectionFirstPacket:
				x := &sessionOwnerSelection{
					FirstPacket: &firstPacketSession{},
				}
				switch e.AaFpSessionSetup {
				case AaFpSessionSetupPrimaryDevice:
					x.FirstPacket.PrimaryDevice = &s
				case AaFpSessionSetupFirstPacket:
					x.FirstPacket.FirstPacket = &s
				case AaFpSessionSetupIpModulo:
					x.FirstPacket.IpModulo = &s
				case AaFpSessionSetupIpHash:
					x.FirstPacket.IpHash = &ipHash{
						HashKey:  e.AaFpSessionSetupIpHashKey,
						HashSeed: e.AaFpSessionSetupIpHashSeed,
					}
				}
				ans.Mode.ActiveActive.SessionOwnerSelection = x
			}
			if text, present := e.raw["vaddr"]; present {
				ans.Mode.ActiveActive.VirtualAddress = &util.RawXml{text}
			}
		}
	}

	ans.Interfaces = &interfaces{} // optional="no"
	if e.Ha1 != nil {
		ans.Interfaces.Ha1 = &ha1Interface{
			Port:             e.Ha1.Port,
			IpAddress:        e.Ha1.IpAddress,
			Netmask:          e.Ha1.Netmask,
			Gateway:          e.Ha1.Gateway,
			EncryptionEnable: util.YesNo(e.Ha1.EncryptionEnable),
			MonitorHoldTime:  e.Ha1.MonitorHoldTime,
		}
	}
	if e.Ha1Backup != nil {
		ans.Interfaces.Ha1Backup = &ha1BackupInterface{
			Port:      e.Ha1Backup.Port,
			IpAddress: e.Ha1Backup.IpAddress,
			Netmask:   e.Ha1Backup.Netmask,
			Gateway:   e.Ha1Backup.Gateway,
		}
	}
	if e.Ha2 != nil {
		ans.Interfaces.Ha2 = &ha2Interface{
			Port:      e.Ha2.Port,
			IpAddress: e.Ha2.IpAddress,
			Netmask:   e.Ha2.Netmask,
			Gateway:   e.Ha2.Gateway,
		}
	}
	if e.Ha2Backup != nil {
		ans.Interfaces.Ha2Backup = &ha2Interface{
			Port:      e.Ha2Backup.Port,
			IpAddress: e.Ha2Backup.IpAddress,
			Netmask:   e.Ha2Backup.Netmask,
			Gateway:   e.Ha2Backup.Gateway,
		}
	}
	if e.Ha3 != nil {
		ans.Interfaces.Ha3 = &ha3Interface{
			Port: e.Ha3.Port,
		}
	}
	if e.Ha4 != nil {
		ans.Interfaces.Ha4 = &ha4Interface{
			Port:      e.Ha4.Port,
			IpAddress: e.Ha4.IpAddress,
			Netmask:   e.Ha4.Netmask,
		}
	}
	if e.Ha4Backup != nil {
		ans.Interfaces.Ha4Backup = &ha4Interface{
			Port:      e.Ha4Backup.Port,
			IpAddress: e.Ha4Backup.IpAddress,
			Netmask:   e.Ha4Backup.Netmask,
		}
	}

	if e.ElectionDevicePriority != "" || e.ElectionPreemptive ||
		e.ElectionHeartBeatBackup {
		ans.ElectionOption = &electionOption{
			DevicePriority:  e.ElectionDevicePriority,
			Preemptive:      util.YesNo(e.ElectionPreemptive),
			HeartBeatBackup: util.YesNo(e.ElectionHeartBeatBackup),
		}
	}

	if e.ElectionTimersMode != "" {
		if ans.ElectionOption == nil {
			ans.ElectionOption = &electionOption{}
		}
		switch e.ElectionTimersMode {
		case ElectionTimersModeRecommended:
			ans.ElectionOption.Timers = &optionTimers{Recommended: &s}
		case ElectionTimersModeAggressive:
			ans.ElectionOption.Timers = &optionTimers{Aggressive: &s}
		case ElectionTimersModeAdvanced:
			x := advanced{
				PromotionHoldTime:          e.ElectionTimersAdvPromotionHoldTime,
				HelloInterval:              e.ElectionTimersAdvHelloInterval,
				HeartBeatInterval:          e.ElectionTimersAdvHeartBeatInterval,
				FlapMax:                    e.ElectionTimersAdvFlapMax,
				PreemptionHoldTime:         e.ElectionTimersAdvPreemptionHoldTime,
				MonitorFailHoldUpTime:      e.ElectionTimersAdvMonitorFailHoldUpTime,
				AdditionalMasterHoldUpTime: e.ElectionTimersAdvAdditionalMasterHoldUpTime,
			}
			ans.ElectionOption.Timers = &optionTimers{Advanced: &x}
		}
	}

	if !e.Ha2StateSyncEnable || e.Ha2StateSyncTransport != "" || e.Ha2StateSyncKeepAliveEnable ||
		e.Ha2StateSyncKeepAliveAction != "" || e.Ha2StateSyncKeepAliveThreshold != 0 {
		ans.StateSync = &stateSync{
			Enable:    util.YesNo(e.Ha2StateSyncEnable),
			Transport: e.Ha2StateSyncTransport,
			Ha2KeepAlive: &ha2KeepAlive{
				Enable:    util.YesNo(e.Ha2StateSyncKeepAliveEnable),
				Action:    e.Ha2StateSyncKeepAliveAction,
				Threshold: e.Ha2StateSyncKeepAliveThreshold,
			},
		}
	}

	if text, present := e.raw["linkgroup"]; present || !e.LinkMonitorEnable ||
		e.LinkMonitorFailureCondition != "" {
		ans.LinkMonitor = &linkMonitor{
			Enable:           util.YesNo(e.LinkMonitorEnable),
			FailureCondition: e.LinkMonitorFailureCondition,
		}
		if present {
			ans.LinkMonitor.LinkGroup = &util.RawXml{text}
		}
	}

	if text, present := e.raw["cluster"]; present {
		ans.Cluster = &util.RawXml{text}
	}
	if text, present := e.raw["pathmonitor"]; present {
		ans.PathMonitor = &util.RawXml{text}
	}

	return ans
}
