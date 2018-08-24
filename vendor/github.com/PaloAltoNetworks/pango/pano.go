package pango

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"

    // Various namespace imports.
    "github.com/PaloAltoNetworks/pango/objs"
    "github.com/PaloAltoNetworks/pango/poli"
    "github.com/PaloAltoNetworks/pango/netw"
    "github.com/PaloAltoNetworks/pango/pnrm"
    "github.com/PaloAltoNetworks/pango/licen"
    "github.com/PaloAltoNetworks/pango/userid"
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
    Licensing *licen.Licen
    UserId *userid.UserId
    Panorama *pnrm.Pnrm
    Objects *objs.PanoObjs
    Policies *poli.PanoPoli
    Network *netw.PanoNetw
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
    } else {
        c.Hostname = "localhost"
        c.ApiKey = "password"
    }
    c.initNamespaces()

    return nil
}

// CommitAll performs a Panorama commit-all.
//
// Param dg is the device group you want to commit-all on.  Note that all other
// params are ignored / unused if the device group is left empty.
//
// Param desc is the optional commit description message you want associated
// with the commit.
//
// Param serials is the list of serial numbers you want to limit the commit-all
// to that are also in the device group dg.
//
// Param tmpl should be true if you want to push template config as well.
//
// Param sync should be true if you want this function to block until the
// commit job completes.
//
// Commits result in a job being submitted to the backend.  The job ID and
// if an error was encountered or not are returned from this function.
func (c *Panorama) CommitAll(dg, desc string, serials []string, tmpl, sync bool) (uint, error) {
    c.LogAction("(commit-all) %q", desc)

    req := panoDgCommit{}
    if dg != "" {
        sp := sharedPolicy{
            Description: desc,
            WithTemplate: util.YesNo(tmpl),
            Dg: deviceGroup{
                Entry: deviceGroupEntry{
                    Name: dg,
                    Devices: util.StrToEnt(serials),
                },
            },
        }
        req.Policy = &sp
    }

    job, _, err := c.CommitConfig(req, "all", nil)
    if err != nil || !sync || job == 0 {
        return job, err
    }

    return job, c.WaitForJob(job, nil)
}

/** Private functions **/

func (c *Panorama) initNamespaces() {
    c.Licensing = &licen.Licen{}
    c.Licensing.Initialize(c)

    c.UserId = &userid.UserId{}
    c.UserId.Initialize(c)

    c.Panorama = &pnrm.Pnrm{}
    c.Panorama.Initialize(c)

    c.Objects = &objs.PanoObjs{}
    c.Objects.Initialize(c)

    c.Policies = &poli.PanoPoli{}
    c.Policies.Initialize(c)

    c.Network = &netw.PanoNetw{}
    c.Network.Initialize(c)
}

/** Internal structs / functions **/

type panoDgCommit struct {
    XMLName xml.Name `xml:"commit-all"`
    Policy *sharedPolicy `xml:"shared-policy"`
}

type sharedPolicy struct {
    Dg deviceGroup `xml:"device-group"`
    Description string `xml:"description,omitempty"`
    WithTemplate string `xml:"include-template"`
}

type deviceGroup struct {
    Entry deviceGroupEntry `xml:"entry"`
}

type deviceGroupEntry struct {
    Name string `xml:"name,attr"`
    Devices *util.EntryType `xml:"devices"`
}
