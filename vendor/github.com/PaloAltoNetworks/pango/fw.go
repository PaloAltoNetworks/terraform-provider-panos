package pango

import (
    // Various namespace imports.
    "github.com/PaloAltoNetworks/pango/netw"
    "github.com/PaloAltoNetworks/pango/dev"
    "github.com/PaloAltoNetworks/pango/poli"
    "github.com/PaloAltoNetworks/pango/objs"
    "github.com/PaloAltoNetworks/pango/licen"
    "github.com/PaloAltoNetworks/pango/userid"
)


// Firewall is a firewall specific client, providing version safe functions
// for the PAN-OS Xpath API methods.  After creating the object, invoke
// Initialize() to prepare it for use.
//
// It has the following namespaces:
//      * Network
//      * Device
//      * Policies
//      * Objects
//      * Licensing
//      * UserId
type Firewall struct {
    Client

    // Namespaces
    Network *netw.Netw
    Device *dev.Dev
    Policies *poli.FwPoli
    Objects *objs.FwObjs
    Licensing *licen.Licen
    UserId *userid.UserId
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
    } else {
        c.Hostname = "localhost"
        c.ApiKey = "password"
    }
    c.initNamespaces()

    return nil
}

/** Private functions **/

func (c *Firewall) initNamespaces() {
    c.Network = &netw.Netw{}
    c.Network.Initialize(c)

    c.Device = &dev.Dev{}
    c.Device.Initialize(c)

    c.Policies = &poli.FwPoli{}
    c.Policies.Initialize(c)

    c.Objects = &objs.FwObjs{}
    c.Objects.Initialize(c)

    c.Licensing = &licen.Licen{}
    c.Licensing.Initialize(c)

    c.UserId = &userid.UserId{}
    c.UserId.Initialize(c)
}
