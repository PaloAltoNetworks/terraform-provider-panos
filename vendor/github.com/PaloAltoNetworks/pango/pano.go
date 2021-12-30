package pango

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	// Various namespace imports.
	"github.com/PaloAltoNetworks/pango/dev"
	"github.com/PaloAltoNetworks/pango/licen"
	"github.com/PaloAltoNetworks/pango/netw"
	"github.com/PaloAltoNetworks/pango/objs"
	"github.com/PaloAltoNetworks/pango/pnrm"
	"github.com/PaloAltoNetworks/pango/poli"
	"github.com/PaloAltoNetworks/pango/predefined"
	"github.com/PaloAltoNetworks/pango/userid"
	"github.com/PaloAltoNetworks/pango/vsys"

	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is a panorama specific client, providing version safe functions
// for the PAN-OS Xpath API methods.  After creating the object, invoke
// Initialize() to prepare it for use.
//
// It has the following namespaces:
//      * Licensing
//      * UserId
type Panorama struct {
	Client

	// Namespaces
	Predefined *predefined.Panorama
	Device     *dev.Panorama
	Licensing  *licen.Licen
	UserId     *userid.UserId
	Panorama   *pnrm.Panorama
	Objects    *objs.PanoObjs
	Policies   *poli.Panorama
	Network    *netw.Panorama
	Vsys       *vsys.Panorama
}

// Initialize does some initial setup of the Panorama connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PAN-OS version.  The full results are saved into the
// client's SystemInfo map.
//
// If not specified, the following is assumed:
//  * Protocol: https
//  * Port: (unspecified)
//  * Timeout: 10
//  * Logging: LogAction | LogUid
func (c *Panorama) Initialize() error {
	if len(c.rb) == 0 {
		var e error

		if e = c.initCon(); e != nil {
			return e
		} else if e = c.initApiKey(); e != nil {
			return e
		} else if e = c.initSystemInfo(); e != nil {
			return e
		}
		c.initPlugins()
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
func (c *Panorama) InitializeUsing(filename string, chkenv bool) error {
	c.CheckEnvironment = chkenv
	c.credsFile = filename

	return c.Initialize()
}

// CreateVmAuthKey creates a VM auth key to bootstrap a VM-Series firewall.
//
// VM auth keys are only valid for the number of hours specified.
func (c *Panorama) CreateVmAuthKey(hours int) (VmAuthKey, error) {
	clock, err := c.Clock()
	if err != nil {
		c.LogOp("(op) Failed to get/parse system time: %s", err)
	}

	type ak_req struct {
		XMLName  xml.Name `xml:"request"`
		Duration int      `xml:"bootstrap>vm-auth-key>generate>lifetime"`
	}

	type ak_resp struct {
		Msg string `xml:"result"`
	}

	req := ak_req{Duration: hours}
	ans := ak_resp{}

	c.LogOp("(op) generating a vm auth code")
	if b, err := c.Op(req, "", nil, &ans); err != nil {
		return VmAuthKey{}, err
	} else if ans.Msg == "" {
		return VmAuthKey{}, fmt.Errorf("No msg: %s", b)
	} else if !strings.HasPrefix(ans.Msg, "VM auth key ") {
		return VmAuthKey{}, fmt.Errorf("Wrong resp prefix: %s", b)
	}

	tokens := strings.Fields(ans.Msg)
	if len(tokens) != 9 {
		return VmAuthKey{}, fmt.Errorf("Got %d of 9 fields from: %s", len(tokens), ans.Msg)
	}

	key := VmAuthKey{
		AuthKey: tokens[3],
		Expiry:  strings.Join(tokens[7:], " "),
	}
	key.ParseExpires(clock)

	return key, nil
}

// GetVmAuthKeys gets the list of VM auth keys.
func (c *Panorama) GetVmAuthKeys() ([]VmAuthKey, error) {
	clock, err := c.Clock()
	if err != nil {
		c.LogOp("(op) Failed to get/parse system time: %s", err)
	}

	type l_req struct {
		XMLName xml.Name `xml:"request"`
		Msg     string   `xml:"bootstrap>vm-auth-key>show"`
	}

	type l_resp struct {
		List []VmAuthKey `xml:"result>bootstrap-vm-auth-keys>entry"`
	}

	req := l_req{}
	ans := l_resp{}

	c.LogOp("(op) listing vm auth codes")
	if _, err := c.Op(req, "", nil, &ans); err != nil {
		return nil, err
	}

	for i := range ans.List {
		ans.List[i].ParseExpires(clock)
	}

	return ans.List, nil
}

// RemoveVmAuthKey revokes a VM auth key.
func (c *Panorama) RevokeVmAuthKey(key string) error {
	type rreq struct {
		XMLName xml.Name `xml:"request"`
		Key     string   `xml:"bootstrap>vm-auth-key>revoke>vm-auth-key"`
	}

	req := rreq{
		Key: key,
	}

	c.LogOp("(op) revoking vm auth code: %s", key)

	_, err := c.Op(req, "", nil, nil)
	return err
}

/** Public structs **/

// VmAuthKey is a VM auth key paired with when it expires.
//
// The Expiry field is the string returned from PAN-OS, while the Expires
// field is an attempt at parsing the Expiry field.
type VmAuthKey struct {
	AuthKey string `xml:"vm-auth-key"`
	Expiry  string `xml:"expiry-time"`
	Expires time.Time
}

// ParseExpires sets Expires from the Expiry field.
//
// Since PAN-OS does not output timezone information with the expirations,
// the current PAN-OS time is retrieved, which does contain timezone
// information.  Then in the string parsing for Expires, the location
// information of the system clock is applied.
func (o *VmAuthKey) ParseExpires(clock time.Time) {
	if t, err := time.ParseInLocation(util.PanosTimeWithoutTimezoneFormat, o.Expiry, clock.Location()); err == nil {
		o.Expires = t
	}
}

/** Private functions **/

func (c *Panorama) initNamespaces() {
	c.Predefined = predefined.PanoramaNamespace(c)

	c.Device = dev.PanoramaNamespace(c)

	c.Licensing = &licen.Licen{}
	c.Licensing.Initialize(c)

	c.UserId = &userid.UserId{}
	c.UserId.Initialize(c)

	c.Panorama = pnrm.PanoramaNamespace(c)

	c.Objects = &objs.PanoObjs{}
	c.Objects.Initialize(c)

	c.Policies = poli.PanoramaNamespace(c)

	c.Network = netw.PanoramaNamespace(c)

	c.Vsys = vsys.PanoramaNamespace(c)
}
