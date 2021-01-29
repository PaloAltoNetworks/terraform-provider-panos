package userid

import (
	"encoding/xml"
)

/*
GetLogins returns a list of IP/user mappings.

If `ip' is not an empty string, filter on the given IP/netmask.

If `lType' is not an empty string and `ip' is specified, then filter on the
given login type.  This can be any of the following:

* AD - Active directory
* CP - Captive Portal
* EDIR - eDirectory
* GP - Global Protect
* GP-CLIENTLESSVPN - Global Protect Clientless VPN
* SSO - SSO
* SYSLOG - Syslog
* UIA - User-ID Agent
* UNKNOWN - Unknown
* XMLAPI - XML API
*/
func (c *UserId) GetLogins(ip, lType, vsys string) ([]LoginInfo, error) {
	if vsys == "" {
		vsys = "vsys1"
	}

	c.con.LogOp("(op) getting ip/user mappings - ip:%s vsys:%s", ip, vsys)

	req := &ipUserReq{
		Data: ipUserData{
			Ip: ip,
		},
	}
	if ip == "" {
		req.Data.All = &ipUserAll{
			Type: lType,
		}
	}

	resp := ipUserResp{}

	if _, err := c.con.Op(req, vsys, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Entries, nil
}

type ipUserReq struct {
	XMLName xml.Name   `xml:"show"`
	Data    ipUserData `xml:"user>ip-user-mapping"`
}

type ipUserData struct {
	All *ipUserAll `xml:"all"`
	Ip  string     `xml:"ip,omitempty"`
}

type ipUserAll struct {
	Type string `xml:"type,omitempty"`
}

type ipUserResp struct {
	Entries []LoginInfo `xml:"result>entry"`
}

// LoginInfo is the structure returned from GetLogins().
type LoginInfo struct {
	Ip          string `xml:"ip"`
	Vsys        string `xml:"vsys"`
	Type        string `xml:"type"`
	User        string `xml:"user"`
	IdleTimeout int    `xml:"idle_timeout"`
	Timeout     int    `xml:"timeout"`
}
