package pango

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"

	// Various namespace imports.
	"github.com/PaloAltoNetworks/pango/dev"
	"github.com/PaloAltoNetworks/pango/licen"
	"github.com/PaloAltoNetworks/pango/netw"
	"github.com/PaloAltoNetworks/pango/objs"
	"github.com/PaloAltoNetworks/pango/panosplugin"
	"github.com/PaloAltoNetworks/pango/poli"
	"github.com/PaloAltoNetworks/pango/predefined"
	"github.com/PaloAltoNetworks/pango/userid"
	"github.com/PaloAltoNetworks/pango/vsys"
)

// Firewall is a firewall specific client, providing version safe functions
// for the PAN-OS Xpath API methods.  After creating the object, invoke
// Initialize() to prepare it for use.
//
// It has the following namespaces:
//      * Predefined
//      * Network
//      * Device
//      * Policies
//      * Objects
//      * Licensing
//      * UserId
type Firewall struct {
	Client

	// Namespaces
	Predefined  *predefined.Firewall
	Network     *netw.Firewall
	Device      *dev.Firewall
	Policies    *poli.Firewall
	Objects     *objs.FwObjs
	Licensing   *licen.Licen
	UserId      *userid.UserId
	Vsys        *vsys.Firewall
	PanosPlugin *panosplugin.Firewall
}

// Initialize does some initial setup of the Firewall connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PAN-OS version.  The full results are saved into the
// client's SystemInfo map.
//
// If not specified, the following is assumed:
//  * Protocol: https
//  * Port: (unspecified)
//  * Timeout: 10
//  * Logging: LogAction | LogUid
func (c *Firewall) Initialize() error {
	if len(c.rb) == 0 {
		var e error

		if e = c.initCon(); e != nil {
			return e
		} else if e = c.initApiKey(); e != nil {
			return e
		} else if e = c.initSystemInfo(); e != nil {
			return e
		}
		if c.Version.Gte(version.Number{9, 0, 0, ""}) {
			c.initPlugins()
		}
	} else {
		c.Hostname = "localhost"
		c.ApiKey = "password"
	}
	c.initNamespaces()

	return nil
}

// InitializeUsing does Initialize(), but takes in a filename that contains
// fallback authentication credentials if they aren't specified.
//
// The order of preference for auth / connection settings is:
//
// * explicitly set
// * environment variable (set chkenv to true to enable this)
// * json file
func (c *Firewall) InitializeUsing(filename string, chkenv bool) error {
	c.CheckEnvironment = chkenv
	c.credsFile = filename

	return c.Initialize()
}

// GetDhcpInfo returns the DHCP client information about the given interface.
func (c *Firewall) GetDhcpInfo(i string) (map[string]string, error) {
	c.LogOp("(op) show dhcp client state %q", i)

	type ireq struct {
		XMLName xml.Name `xml:"show"`
		Val     string   `xml:"dhcp>client>state"`
	}

	type ireq_ans struct {
		Interface  string `xml:"result>entry>interface"`
		State      string `xml:"result>entry>state"`
		Ip         string `xml:"result>entry>ip"`
		Gateway    string `xml:"result>entry>gw"`
		Server     string `xml:"result>entry>server"`
		ServerId   string `xml:"result>entry>server-id"`
		Dns1       string `xml:"result>entry>dns1"`
		Dns2       string `xml:"result>entry>dns2"`
		Wins1      string `xml:"result>entry>wins1"`
		Wins2      string `xml:"result>entry>wins2"`
		Nis1       string `xml:"result>entry>nis1"`
		Nis2       string `xml:"result>entry>nis2"`
		Ntp1       string `xml:"result>entry>ntp1"`
		Ntp2       string `xml:"result>entry>ntp2"`
		Pop3Server string `xml:"result>entry>pop3"`
		SmtpServer string `xml:"result>entry>smtp"`
		DnsSuffix  string `xml:"result>entry>dns-suffix"`
	}

	req := ireq{Val: i}
	ans := ireq_ans{}

	if _, err := c.Op(req, "", nil, &ans); err != nil {
		return nil, err
	}

	return map[string]string{
		"interface":      ans.Interface,
		"state":          ans.State,
		"ip":             ans.Ip,
		"gateway":        ans.Gateway,
		"server":         ans.Server,
		"server_id":      ans.ServerId,
		"primary_dns":    ans.Dns1,
		"secondary_dns":  ans.Dns2,
		"primary_wins":   ans.Wins1,
		"secondary_wins": ans.Wins2,
		"primary_nis":    ans.Nis1,
		"secondary_nis":  ans.Nis2,
		"primary_ntp":    ans.Ntp1,
		"secondary_ntp":  ans.Ntp2,
		"pop3_server":    ans.Pop3Server,
		"smtp_server":    ans.SmtpServer,
		"dns_suffix":     ans.DnsSuffix,
	}, nil
}

/** Private functions **/

func (c *Firewall) initNamespaces() {
	c.Predefined = predefined.FirewallNamespace(c)

	c.Network = netw.FirewallNamespace(c)

	c.Device = dev.FirewallNamespace(c)

	c.Policies = poli.FirewallNamespace(c)

	c.Objects = &objs.FwObjs{}
	c.Objects.Initialize(c)

	c.Licensing = &licen.Licen{}
	c.Licensing.Initialize(c)

	c.UserId = &userid.UserId{}
	c.UserId.Initialize(c)

	c.Vsys = vsys.FirewallNamespace(c)
	c.PanosPlugin = panosplugin.FirewallNamespace(c)
}
