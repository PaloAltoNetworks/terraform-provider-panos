/*
Package pango is a golang cross version mechanism for interacting with Palo Alto
Networks devices (including physical and virtualized Next-generation Firewalls
and Panorama).  Versioning support is in place for PANOS 6.1 to 8.0.

To start, create a client connection with the desired parameters and then
initialize the connection:

    package main
    
    import (
        "log"
        "github.com/PaloAltoNetworks/pango"
    )
    
    func main() {
        var err error
        c := pango.Firewall{Client: pango.Client{
            Hostname: "127.0.0.1",
            Username: "admin",
            Password: "admin",
            Logging: pango.LogAction | pango.LogOp,
        }}
        if err = c.Initialize(); err != nil {
            log.Printf("Failed to initialize client: %s", err)
            return
        }
        log.Printf("Initialize ok")
    }

Initializing the connection creates the API key (if it was not already
specified), then performs "show system info" to get the PANOS version.  Once
the firewall client is created, you can query and configure the Palo
Alto Networks device from the functions inside the various namespaces of the
client connection.  Namespaces correspond to the various configuration areas
available in the GUI.  For example:

    err = c.Network.EthernetInterface.Set(...)
    myPolicies, err := c.Policies.Security.GetList(...)

Generally speaking, there are the following functions inside each namespace:

    * GetList
    * ShowList
    * Get
    * Show
    * Set
    * Edit
    * Delete

These functions correspond with PANOS Get, Show, Set, Edit, and
Delete API calls.  Get(), Set(), and Edit() take and return normalized,
version independent objects.  These version safe objects are typically named
Entry, which corresponds to how the object is placed in the PANOS XPATH.

Some Entry objects have a special function, Defaults().  Invoking this
function will initialize the object with some default values.  Each Entry
that implements Defaults() calls out in its documentation what parameters
are affected by this, and what the defaults are.

For any version safe object, attempting to configure a parameter that your
PANOS doesn't support will be safely ignored in the resultant XML sent to the
firewall / Panorama.

Using Edit Functions

The PANOS XML API Edit command can be used to both create as well as update
existing config, however it can also truncate config for the given XPATH.  Due
to this, if you want to use Edit(), you need to make sure that you perform
either a Get() or a Show() first, make your modification, then invoke
Edit() using that object.  If you don't do this, you will truncate any sub
config.

To learn more about PANOS XML API, please refer to the Palo Alto Netowrks
API documentation.

Examples

The following program will create ethernet1/7 as a DHCP interface and import
it into vsys1 if it isn't already present:

    package main
    
    import (
        "log"
        "github.com/PaloAltoNetworks/pango"
        "github.com/PaloAltoNetworks/pango/netw/eth"
    )
    
    func main() {
        var err error
    
        c := &pango.Firewall{Client: pango.Client{
            Hostname: "127.0.0.1",
            Username: "admin",
            Password: "admin",
            Logging: pango.LogAction | pango.LogOp,
        }}
        if err = c.Initialize(); err != nil {
            log.Printf("Failed to initialize client: %s", err)
            return
        }
    
        e := eth.Entry{
            Name: "ethernet1/7",
            Mode: "layer3",
            EnableDhcp: true,
            CreateDhcpDefaultRoute: true,
        }
    
        interfaces, err := c.Network.EthernetInterface.GetList()
        if err != nil {
            log.Printf("Failed to get data interfaces: %s", err)
            return
        }
        for i := range interfaces {
            if e.Name == interfaces[i] {
                log.Printf("%s already exists", e.Name)
                return
            }
        }
    
        err = c.Network.EthernetInterface.Set("vsys1", e)
        if err != nil {
            log.Printf("Failed to create %s: %s", e.Name, err)
            return
        }
        log.Printf("Created %s ok", e.Name)
    }

*/
package pango

import (
    "crypto/tls"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "time"

    "github.com/PaloAltoNetworks/pango/version"
    "github.com/PaloAltoNetworks/pango/util"

    // Various namespace imports.
    "github.com/PaloAltoNetworks/pango/netw"
    "github.com/PaloAltoNetworks/pango/dev"
    "github.com/PaloAltoNetworks/pango/poli"
    "github.com/PaloAltoNetworks/pango/objs"
    "github.com/PaloAltoNetworks/pango/licen"
    "github.com/PaloAltoNetworks/pango/userid"
)


// These bit flags control what is logged by client connections.  Of the flags
// available for use, LogSend and LogReceive will log ALL communication between
// the connection object and the PANOS XML API.  The API key being used for
// communication will be blanked out, but no other sensitive data will be.  As
// such, those two flags should be considered for debugging only.  To disable
// all logging, set the logging level as LogQuiet.
//
// The bit-wise flags are as follows:
//
//      * LogQuiet: disables all logging
//      * LogAction: action being performed (Set / Delete functions)
//      * LogQuery: queries being run (Get / Show functions)
//      * LogOp: operation commands (Op functions)
//      * LogUid: User-Id commands (Uid functions)
//      * LogXpath: the resultant xpath
//      * LogSend: xml docuemnt being sent
//      * LogReceive: xml responses being received
const (
    LogQuiet = 1 << (iota + 1)
    LogAction
    LogQuery
    LogOp
    LogUid
    LogXpath
    LogSend
    LogReceive
)

// Client is a generic connector struct.  It provides wrapper functions for
// invoking the various PANOS XPath API methods.  After creating the client,
// invoke Initialize() to prepare it for use.
type Client struct {
    // Connection properties.
    Hostname string
    Username string
    Password string
    ApiKey string
    Protocol string
    Port uint
    Timeout int
    Target string

    // Variables determined at runtime.
    Version version.Number
    SystemInfo map[string] string

    // Logging level.
    Logging uint32

    // Internal variables.
    con *http.Client
    api_url string

    // Variables for testing, response bytes and response index.
    rp []url.Values
    rb [][]byte
    ri int
}

// Firewall is a firewall specific client, providing version safe functions
// for the PANOS Xpath API methods.  After creating the object, invoke
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
    Policies *poli.Poli
    Objects *objs.Objs
    Licensing *licen.Licen
    UserId *userid.UserId
}

// Panorama is a panorama specific client, providing version safe functions
// for the PANOS Xpath API methods.  After creating the object, invoke
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
}

// String is the string representation of a client connection.  Both the
// password and API key are replaced with stars, if set, making it safe
// to print the client connection in log messages.
func (c *Client) String() string {
    var passwd string
    var api_key string

    if c.Password == "" {
        passwd = ""
    } else {
        passwd = "********"
    }

    if c.ApiKey == "" {
        api_key = ""
    } else {
        api_key = "********"
    }

    return fmt.Sprintf(
        "{Hostname:%s Username:%s Password:%s ApiKey:%s Protocol:%s Port:%d Timeout:%d Logging:%d}",
        c.Hostname, c.Username, passwd, api_key, c.Protocol, c.Port, c.Timeout, c.Logging)
}

// Versioning returns the client version number.
func (c *Client) Versioning() version.Number {
    return c.Version
}

// Initialize does some initial setup of the Client connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PANOS version.  The full results are saved into the
// client's SystemInfo map.
//
// If not specified, the following is assumed:
//  * Protocol: https
//  * Port: (unspecified)
//  * Timeout: 10
//  * Logging: LogAction | LogUid
func (c *Client) Initialize() error {
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

    return nil
}

// Initialize does some initial setup of the Firewall connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PANOS version.  The full results are saved into the
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

// Initialize does some initial setup of the Panorama connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PANOS version.  The full results are saved into the
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

// RetrieveApiKey retrieves the API key, which will require that both the
// username and password are defined.
//
// The currently set ApiKey is forgotten when invoking this function.
func (c *Client) RetrieveApiKey() error {
    c.LogAction("%s: Retrieving API key", c.Hostname)

    type key_gen_ans struct {
        Key string `xml:"result>key"`
    }

    c.ApiKey = ""
    ans := key_gen_ans{}
    data := url.Values{}
    data.Add("user", c.Username)
    data.Add("password", c.Password)
    data.Add("type", "keygen")

    _, err := c.Communicate(data, &ans)
    if err != nil {
        return err
    }

    c.ApiKey = ans.Key

    return nil
}

// EntryListUsing retrieves an list of entries using the given function, either
// Get or Show.
func (c *Client) EntryListUsing(fn util.Retriever, path []string) ([]string, error) {
    var err error
    type Entry struct {
        Name string `xml:"name,attr"`
    }

    type resp_struct struct {
        Entries []Entry `xml:"result>entry"`
    }

    if path == nil {
        return nil, fmt.Errorf("xpath is empty")
    }
    path = append(path, "entry", "@name")
    resp := resp_struct{}

    _, err = fn(path, nil, &resp)
    if err != nil {
        e2, ok := err.(PanosError)
        if ok && e2.ObjectNotFound() {
            return nil, nil
        }
        return nil, err
    }

    ans := make([]string, len(resp.Entries))
    for i := range resp.Entries {
        ans[i] = resp.Entries[i].Name
    }

    return ans, nil
}

// MemberListUsing retrieves an list of members using the given function, either
// Get or Show.
func (c *Client) MemberListUsing(fn util.Retriever, path []string) ([]string, error) {
    type resp_struct struct {
        Members []string `xml:"result>member"`
    }

    if path == nil {
        return nil, fmt.Errorf("xpath is empty")
    }
    path = append(path, "member")
    resp := resp_struct{}

    _, err := fn(path, nil, &resp)
    if err != nil {
        e2, ok := err.(PanosError)
        if ok && e2.ObjectNotFound() {
            return nil, nil
        }
        return nil, err
    }

    return resp.Members, nil
}

// RequestPasswordHash requests a password hash of the given string.
func (c *Client) RequestPasswordHash(val string) (string, error) {
    c.LogOp("(op) creating password hash")
    type phash_req struct {
        XMLName xml.Name `xml:"request"`
        Val string `xml:"password-hash>password"`
    }

    type phash_ans struct {
        Hash string `xml:"result>phash"`
    }

    req := phash_req{Val: val}
    ans := phash_ans{}

    if _, err := c.Op(req, "", nil, &ans); err != nil {
        return "", err
    }

    return ans.Hash, nil
}

// ImportVlans imports VLANs into the vsys specified.  VLANs must be
// imported into a vsys before they can be used.
//
// This is invoked automatically when creating VLANs as long as the vsys given
// is not an empty string.
func (c *Client) ImportVlans(vsys string, names []string) error {
    return c.vsysImport("vlan", vsys, names)
}

// UnimportVlans unimports VLANs from the vsys specified.  VLANs that are
// imported into an vsys cannot be deleted.
//
// This is invoked automatically when deleting VLANs as long as the vsys given
// is not an empty string.
func (c *Client) UnimportVlans(vsys string, names []string) error {
    return c.vsysUnimport("vlan", vsys, names)
}

// ImportInterfaces imports interfaces into the vsys specified.  Interfaces
// must be imported into a vsys before they can be used.
//
// This is invoked automatically when creating interfaces as long as the
// vsys given is not an empty string.
func (c *Client) ImportInterfaces(vsys string, names []string) error {
    return c.vsysImport("interface", vsys, names)
}

// UnimportInterfaces unimports interfaces from the vsys specified.  Interfaces
// that are imported into an vsys cannot be deleted.
//
// This is invoked automatically when deleting interfaces as long as the
// vsys given is not an empty string.
func (c *Client) UnimportInterfaces(vsys string, names []string) error {
    return c.vsysUnimport("interface", vsys, names)
}

// ImportVirtualRouters imports virtual routers into the vsys specified.
// Virtual routers that are imported into a vsys cannot be deleted.
//
// This is invoked automatically when creating virtual routers as long as the
// vsys given is not an empty string.
func (c *Client) ImportVirtualRouters(vsys string, names []string) error {
    return c.vsysImport("virtual-router", vsys, names)
}

// UnimportVirtualRouters unimports virtual routers from the vsys specified.
// Virtual routers that are imported into an vsys cannot be deleted.
//
// This is invoked automatically when deleting virtual routers as long as the
// vsys given is not an empty string.
func (c *Client) UnimportVirtualRouters(vsys string, names []string) error {
    return c.vsysUnimport("virtual-router", vsys, names)
}

// ValidateConfig performs a commit config validation check.
//
// Setting sync to true means that this function will block until the job
// finishes.
//
// This function returns the job ID and if any errors were encountered.
func (c *Client) ValidateConfig(sync bool) (uint, error) {
    var err error

    c.LogOp("(op) validating config")
    type op_req struct {
        XMLName xml.Name `xml:"validate"`
        Cmd string `xml:"full"`
    }
    job_ans := util.JobResponse{}
    _, err = c.Op(op_req{}, "", nil, &job_ans)
    if err != nil {
        return 0, err
    }

    id := job_ans.Id
    if !sync {
        return id, nil
    }

    return id, c.WaitForJob(id, nil)
}

// RevertToRunningConfig discards any changes made and reverts to the last
// config committed.
func (c *Client) RevertToRunningConfig() error {
    c.LogOp("(op) reverting to running config")
    _, err := c.Op("<load><config><from>running-config.xml</from></config></load>", "", nil, nil)
    return err
}

// ConfigLocks returns any config locks that are currently in place.
//
// If vsys is an empty string, then the vsys will default to "shared".
func (c *Client) ConfigLocks(vsys string) ([]util.Lock, error) {
    if vsys == "" {
        vsys = "shared"
    }

    c.LogOp("(op) getting config locks for scope %q", vsys)
    ans := configLocks{}
    _, err := c.Op("<show><config-locks /></show>", vsys, nil, &ans)
    if err != nil {
        return nil, err
    }
    return ans.Locks, nil
}

// LockConfig locks the config for the given scope with the given comment.
//
// If vsys is an empty string, the scope defaults to "shared".
func (c *Client) LockConfig(vsys, comment string) error {
    if vsys == "" {
        vsys = "shared"
    }
    c.LogOp("(op) locking config for scope %q", vsys)

    var inner string
    if comment == "" {
        inner = "<add />"
    } else {
        inner = fmt.Sprintf("<add><comment>%s</comment></add>", comment)
    }
    cmd := fmt.Sprintf("<request><config-lock>%s</config-lock></request>", inner)

    _, err := c.Op(cmd, vsys, nil, nil)
    return err
}

// UnlockConfig removes the config lock on the given scope.
//
// If vsys is an empty string, the scope defaults to "shared".
func (c *Client) UnlockConfig(vsys string) error {
    if vsys == "" {
        vsys = "shared"
    }

    type cmd struct {
        XMLName xml.Name `xml:"request"`
        Cmd string `xml:"config-lock>remove"`
    }

    c.LogOp("(op) unlocking config for scope %q", vsys)
    _, err := c.Op(cmd{}, vsys, nil, nil)
    return err
}

// CommitLocks returns any commit locks that are currently in place.
//
// If vsys is an empty string, then the vsys will default to "shared".
func (c *Client) CommitLocks(vsys string) ([]util.Lock, error) {
    if vsys == "" {
        vsys = "shared"
    }

    c.LogOp("(op) getting commit locks for scope %q", vsys)
    ans := commitLocks{}
    _, err := c.Op("<show><commit-locks /></show>", vsys, nil, &ans)
    if err != nil {
        return nil, err
    }
    return ans.Locks, nil
}

// LockCommits locks commits for the given scope with the given comment.
//
// If vsys is an empty string, the scope defaults to "shared".
func (c *Client) LockCommits(vsys, comment string) error {
    if vsys == "" {
        vsys = "shared"
    }
    c.LogOp("(op) locking commits for scope %q", vsys)

    var inner string
    if comment == "" {
        inner = "<add />"
    } else {
        inner = fmt.Sprintf("<add><comment>%s</comment></add>", comment)
    }
    cmd := fmt.Sprintf("<request><commit-lock>%s</commit-lock></request>", inner)

    _, err := c.Op(cmd, vsys, nil, nil)
    return err
}

// UnlockCommits removes the commit lock on the given scope owned by the given
// admin, if this admin is someone other than the current acting admin.
//
// If vsys is an empty string, the scope defaults to "shared".
func (c *Client) UnlockCommits(vsys, admin string) error {
    if vsys == "" {
        vsys = "shared"
    }

    type cmd struct {
        XMLName xml.Name `xml:"request"`
        Admin string `xml:"commit-lock>remove>admin,omitempty"`
    }

    c.LogOp("(op) unlocking commits for scope %q", vsys)
    _, err := c.Op(cmd{Admin: admin}, vsys, nil, nil)
    return err
}

// WaitForJob polls the device, waiting for the specified job to finish.
//
// If you want to unmarshal the response into a struct, then pass in a
// pointer to the struct for the "resp" param.  If you just want to know if
// the job completed with a status other than "FAIL", you only need to check
// the returned error message.
//
// In the case that there are multiple errors returned from the job, the first
// error is returned as the error string, and no unmarshaling is attempted.
func (c *Client) WaitForJob(id uint, resp interface{}) error {
    var err error
    var prev uint
    var data []byte

    c.LogOp("(op) waiting for job %d", id)
    type op_req struct {
        XMLName xml.Name `xml:"show"`
        Id uint `xml:"jobs>id"`
    }
    req := op_req{Id: id}

    ans := util.BasicJob{}
    for ans.Progress != 100 {
        // Get current percent complete.
        data, err = c.Op(req, "", nil, &ans)
        if err != nil {
            return err
        }
        // Output percent complete if it's new.
        if ans.Progress != prev {
            prev = ans.Progress
            c.LogOp("(op) job %d: %d percent complete", id, prev)
        }
    }

    if ans.Result == "FAIL" {
        if len(ans.Details) > 0 {
            return fmt.Errorf(ans.Details[0])
        } else {
            return fmt.Errorf("Job %d has failed to complete successfully", id)
        }
    }

    if resp == nil {
        return nil
    }

    return xml.Unmarshal(data, resp)
}

// Commit performs a Firewall commit.
//
// Param desc is the optional commit description message you want associated
// with the commit.
//
// Params dan and pao are advanced options for doing partial commits.  Setting
// param dan to false excludes the Device and Network configuration, while
// setting param pao to false excludes the Policy and Object configuration.
//
// Param force is if you want to force a commit even if no changes are
// required.
//
// Param sync should be true if you want this function to block until the
// commit job completes.
//
// Commits result in a job being submitted to the backend.  The job ID and
// if an error was encountered or not are returned from this function.
func (c *Firewall) Commit(desc string, dan, pao, force, sync bool) (uint, error) {
    c.LogAction("(commit) %q", desc)

    req := fwCommit{Description: desc}
    if !dan || !pao {
        req.Partial = &fwCommitPartial{}
        if !dan {
            req.Partial.Dan = "excluded"
        }
        if !pao {
            req.Partial.Pao = "excluded"
        }
    }
    if force {
        req.Force = ""
    }

    job, _, err := c.CommitConfig(req, "", nil)
    if err != nil || !sync || job == 0 {
        return job, err
    }

    return job, c.WaitForJob(job, nil)
}

// LogAction writes a log message for SET/DELETE operations if LogAction is set.
func (c *Client) LogAction(msg string, i ...interface{}) {
    if c.Logging & LogAction == LogAction {
        log.Printf(msg, i...)
    }
}

// LogQuery writes a log message for GET/SHOW operations if LogQuery is set.
func (c *Client) LogQuery(msg string, i ...interface{}) {
    if c.Logging & LogQuery == LogQuery {
        log.Printf(msg, i...)
    }
}

// LogOp writes a log message for OP operations if LogOp is set.
func (c *Client) LogOp(msg string, i ...interface{}) {
    if c.Logging & LogOp == LogOp {
        log.Printf(msg, i...)
    }
}

// LogUid writes a log message for User-Id operations if LogUid is set.
func (c *Client) LogUid(msg string, i ...interface{}) {
    if c.Logging & LogUid == LogUid {
        log.Printf(msg, i...)
    }
}

// Communicate sends the given data to PANOS.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
//
// Even if an answer struct is given, we first check for known error formats.  If
// a known error format is detected, unmarshalling into the answer struct is not
// performed.
//
// If the API key is set, but not present in the given data, then it is added in.
func (c *Client) Communicate(data url.Values, ans interface{}) ([]byte, error) {
    if c.ApiKey != "" && data.Get("key") == "" {
        data.Set("key", c.ApiKey)
    }

    if c.Logging & LogSend == LogSend {
        old_key := data.Get("key")
        if old_key != "" {
            data.Set("key", "########")
        }
        log.Printf("Sending data: %#v", data)
        if old_key != "" {
            data.Set("key", old_key)
        }
    }

    body, err := c.post(data)
    if err != nil {
        return nil, err
    }
    if c.Logging & LogReceive == LogReceive {
        log.Printf("Response = %s", body)
    }

    // Check for errors first
    errType1 := &panosErrorResponseWithoutLine{}
    err = xml.Unmarshal(body, errType1)
    // At this point, we make use of the shared error error checking that exists
    // between response types.  If the first response is not an error type, we
    // don't have to check the others.  We can get some modest speed gains as
    // a result.
    if errType1.Failed() {
        if err == nil && errType1.Error() != "" {
            return body, PanosError{errType1.Error(), errType1.ResponseCode}
        }
        errType2 := panosErrorResponseWithLine{}
        err = xml.Unmarshal(body, &errType2)
        if err == nil && errType2.Error() != "" {
            return body, PanosError{errType2.Error(), errType2.ResponseCode}
        }
        // Still an error, but some unknown format.
        return body, fmt.Errorf("Unknown error format: %s", body)
    }

    // Return the body string if we weren't given something to unmarshal into
    if ans == nil {
        return body, nil
    }

    // Unmarshal using the struct passed in
    err = xml.Unmarshal(body, ans)
    if err != nil {
        return body, fmt.Errorf("Error unmarshaling into provided interface: %s", err)
    }

    return body, nil
}

// Op runs an operational or "op" type command.
//
// The req param can be either a properly formatted XML string or a struct
// that can be marshalled into XML.
//
// The vsys param is the vsys the op command should be executed in, if any.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Op(req interface{}, vsys string, extras, ans interface{}) ([]byte, error) {
    var err error
    data := url.Values{}
    data.Set("type", "op")

    if err = addToData("cmd", req, true, &data); err != nil {
        return nil, err
    }

    if vsys != "" {
        data.Set("vsys", vsys)
    }

    if c.Target != "" {
        data.Set("target", c.Target)
    }

    if err = mergeUrlValues(&data, extras); err != nil {
        return nil, err
    }

    return c.Communicate(data, ans)
}

// Show runs a "show" type command.
//
// The path param should be either a string or a slice of strings.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Show(path, extras, ans interface{}) ([]byte, error) {
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    return c.typeConfig("show", data, extras, ans)
}

// Get runs a "get" type command.
//
// The path param should be either a string or a slice of strings.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Get(path, extras, ans interface{}) ([]byte, error) {
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    return c.typeConfig("get", data, extras, ans)
}

// Delete runs a "delete" type command, removing the supplied xpath and
// everything underneath it.
//
// The path param should be either a string or a slice of strings.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Delete(path, extras, ans interface{}) ([]byte, error) {
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    return c.typeConfig("delete", data, extras, ans)
}

// Set runs a "set" type command, creating the element at the given xpath.
//
// The path param should be either a string or a slice of strings.
//
// The element param can be either a string of properly formatted XML to send
// or a struct which can be marshaled into a string.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Set(path, element, extras, ans interface{}) ([]byte, error) {
    var err error
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    if err = addToData("element", element, true, &data); err != nil {
        return nil, err
    }

    return c.typeConfig("set", data, extras, ans)
}

// Edit runs a "edit" type command, modifying what is at the given xpath
// with the supplied element.
//
// The path param should be either a string or a slice of strings.
//
// The element param can be either a string of properly formatted XML to send
// or a struct which can be marshaled into a string.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Edit(path, element, extras, ans interface{}) ([]byte, error) {
    var err error
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    if err = addToData("element", element, true, &data); err != nil {
        return nil, err
    }

    return c.typeConfig("edit", data, extras, ans)
}

// Move does a "move" type command.
func (c *Client) Move(path interface{}, where, dst string, extras, ans interface{}) ([]byte, error) {
    data := url.Values{}
    xp := util.AsXpath(path)
    c.logXpath(xp)
    data.Set("xpath", xp)

    if where != "" {
        data.Set("where", where)
    }

    if dst != "" {
        data.Set("dst", dst)
    }

    return c.typeConfig("move", data, extras, ans)
}

// Uid performs User-ID API calls.
func (c *Client) Uid(cmd interface{}, vsys string, extras, ans interface{}) ([]byte, error) {
    var err error
    data := url.Values{}
    data.Set("type", "user-id")

    if err = addToData("cmd", cmd, true, &data); err != nil {
        return nil, err
    }

    if vsys != "" {
        data.Set("vsys", vsys)
    }

    if c.Target != "" {
        data.Set("target", c.Target)
    }

    if err = mergeUrlValues(&data, extras); err != nil {
        return nil, err
    }

    return c.Communicate(data, ans)
}

// CommitConfig performs PANOS commits.  This is the underlying function
// invoked by Firewall.Commit() and Panorama.Commit().
//
// The cmd param can be either a properly formatted XML string or a struct
// that can be marshalled into XML.
//
// The action param is the commit action to be taken, if any (e.g. - "all").
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// Commits result in a job being submitted to the backend.  The job ID, assuming
// the commit action was successfully submitted, the response from the server,
// and if an error was encountered or not are all returned from this function.
func (c *Client) CommitConfig(cmd interface{}, action string, extras interface{}) (uint, []byte, error) {
    var err error
    data := url.Values{}
    data.Set("type", "commit")

    if err = addToData("cmd", cmd, true, &data); err != nil {
        return 0, nil, err
    }

    if action != "" {
        data.Set("action", action)
    }

    if c.Target != "" {
        data.Set("target", c.Target)
    }

    if err = mergeUrlValues(&data, extras); err != nil {
        return 0, nil, err
    }

    ans := util.JobResponse{}
    b, err := c.Communicate(data, &ans)
    return ans.Id, b, err
}

/*** Internal functions ***/

func (c *Client) initCon() error {
    var tout time.Duration

    // Sets the logging level.
    if c.Logging == 0 {
        c.Logging = LogAction | LogUid
    }

    // Set the timeout
    if c.Timeout == 0 {
        c.Timeout = 10
    } else if c.Timeout > 60 {
        return fmt.Errorf("Timeout for %q is %d, expecting a number between [0, 60]", c.Hostname, c.Timeout)
    }
    tout = time.Duration(time.Duration(c.Timeout) * time.Second)

    // Set the protocol
    if c.Protocol == "" {
        c.Protocol = "https"
    } else if c.Protocol != "http" && c.Protocol != "https" {
        return fmt.Errorf("Invalid protocol %q.  Must be \"http\" or \"https\"", c.Protocol)
    }

    // Check port number
    if c.Port > 65535 {
        return fmt.Errorf("Port %d is out of bounds", c.Port)
    }

    // Setup the https client
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    c.con = &http.Client{
        Transport: tr,
        Timeout: tout,
    }

    // Configure the api url
    if c.Port == 0 {
        c.api_url = fmt.Sprintf("%s://%s/api", c.Protocol, c.Hostname)
    } else {
        c.api_url = fmt.Sprintf("%s://%s:%d/api", c.Protocol, c.Hostname, c.Port)
    }

    return nil
}

func (c *Client) initApiKey() error {
    if c.ApiKey != "" {
        return nil
    }

    return c.RetrieveApiKey()
}

func (c *Client) initSystemInfo() error {
    var err error
    c.LogOp("(op) show system info")

    // Run "show system info"
    type system_info_req struct {
        XMLName xml.Name `xml:"show"`
        Cmd string `xml:"system>info"`
    }

    type tagVal struct {
        XMLName xml.Name
        Value string `xml:",chardata"`
    }

    type sysTag struct {
        XMLName xml.Name `xml:"system"`
        Tag []tagVal `xml:",any"`
    }

    type system_info_ans struct {
        System sysTag `xml:"result>system"`
    }

    req := system_info_req{}
    ans := system_info_ans{}

    _, err = c.Op(req, "", nil, &ans)
    if err != nil {
        return fmt.Errorf("Error getting system info: %s", err)
    }

    c.SystemInfo = make(map[string] string, len(ans.System.Tag))
    for i := range ans.System.Tag {
        c.SystemInfo[ans.System.Tag[i].XMLName.Local] = ans.System.Tag[i].Value
        if ans.System.Tag[i].XMLName.Local == "sw-version" {
            c.Version, err = version.New(ans.System.Tag[i].Value)
            if err != nil {
                return fmt.Errorf("Error parsing version %s: %s", ans.System.Tag[i].Value, err)
            }
        }
    }

    return nil
}

func (c *Firewall) initNamespaces() {
    c.Network = &netw.Netw{}
    c.Network.Initialize(c)

    c.Device = &dev.Dev{}
    c.Device.Initialize(c)

    c.Policies = &poli.Poli{}
    c.Policies.Initialize(c)

    c.Objects = &objs.Objs{}
    c.Objects.Initialize(c)

    c.Licensing = &licen.Licen{}
    c.Licensing.Initialize(c)

    c.UserId = &userid.UserId{}
    c.UserId.Initialize(c)
}

func (c *Panorama) initNamespaces() {
    c.Licensing = &licen.Licen{}
    c.Licensing.Initialize(c)

    c.UserId = &userid.UserId{}
    c.UserId.Initialize(c)
}

func (c *Client) typeConfig(action string, data url.Values, extras, ans interface{}) ([]byte, error) {
    var err error

    data.Set("type", "config")
    data.Set("action", action)
    if c.Target != "" {
        data.Set("target", c.Target)
    }

    if err = mergeUrlValues(&data, extras); err != nil {
        return nil, err
    }

    return c.Communicate(data, ans)
}

func (c *Client) logXpath(p string) {
    if c.Logging & LogXpath == LogXpath {
        log.Printf("(xpath) %s", p)
    }
}

func (c *Client) vsysImport(loc, vsys string, names []string) error {
    path := c.xpathImport(vsys)
    if len(names) == 0 || vsys == "" {
        return nil
    } else if len(names) == 1 {
        path = append(path, loc)
    }

    obj := util.BulkElement{XMLName: xml.Name{Local: loc}}
    for i := range names {
        obj.Data = append(obj.Data, vis{xml.Name{Local: "member"}, names[i]})
    }

    _, err := c.Set(path, obj.Config(), nil, nil)
    return err
}

func (c *Client) vsysUnimport(loc, vsys string, names []string) error {
    if len(names) == 0 || vsys == "" {
        return nil
    }

    path := c.xpathImport(vsys)
    path = append(path, loc, util.AsMemberXpath(names))

    _, err := c.Delete(path, nil, nil)
    return err
}

func (c *Client) xpathImport(vsys string) ([]string) {
    return []string {
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "vsys",
        util.AsEntryXpath([]string{vsys}),
        "import",
        "network",
    }
}

func (c *Client) post(data url.Values) ([]byte, error) {
    if len(c.rb) == 0 {
        r, err := c.con.PostForm(c.api_url, data)
        if err != nil {
            return nil, err
        }

        defer r.Body.Close()
        return ioutil.ReadAll(r.Body)
    } else {
        if c.ri < len(c.rb) {
            c.rp = append(c.rp, data)
        }
        body := c.rb[c.ri % len(c.rb)]
        c.ri++
        return body, nil
    }
}

/** Non-struct private functions **/

func mergeUrlValues(data *url.Values, extras interface{}) error {
    if extras == nil {
        return nil
    }

    ev, ok := extras.(url.Values)
    if !ok {
        return fmt.Errorf("extras needs to be of type url.Values or nil")
    }

    for key := range ev {
        data.Set(key, ev.Get(key))
    }

    return nil
}

func addToData(key string, i interface{}, attemptMarshal bool, data *url.Values) error {
    if i == nil {
        return nil
    }

    val, err := asString(i, attemptMarshal)
    if err != nil {
        return err
    }

    data.Set(key, val)
    return nil
}

func asString(i interface{}, attemptMarshal bool) (string, error) {
    switch val := i.(type) {
    case string:
        return val, nil
    case fmt.Stringer:
        return val.String(), nil
    case nil:
        return "", fmt.Errorf("nil encountered")
    default:
        if !attemptMarshal {
            return "", fmt.Errorf("value must be string or fmt.Stringer")
        }

        rb, err := xml.Marshal(val)
        if err != nil {
            return "", err
        }
        return string(rb), nil
    }
}

// PanosError is the error struct returned from the Communicate method.
type PanosError struct {
    Msg string
    Code int
}

// Error returns the error message.
func (e PanosError) Error() string {
    return e.Msg
}

// ObjectNotFound returns true on missing object error.
func (e PanosError) ObjectNotFound() bool {
    return e.Code == 7
}

/*
// Code returns the error code.
func (e PanosError) Code() int {
    return e.ErrCode
}
*/

type panosStatus struct {
    ResponseStatus string `xml:"status,attr"`
    ResponseCode int `xml:"code,attr"`
}

// Failed checks for a status of "failed" or "error".
func (e panosStatus) Failed() bool {
    if e.ResponseStatus == "failed" || e.ResponseStatus == "error" {
        return true
    } else if e.ResponseCode == 0 || e.ResponseCode == 19 || e.ResponseCode == 20 {
        return false
    } else {
        return true
    }
}

func (e panosStatus) codeError() string {
    switch e.ResponseCode {
    case 1:
        return "Unknown command"
    case 2, 3, 4, 5, 11:
        return fmt.Sprintf("Internal error (%d) encountered", e.ResponseCode)
    case 6:
        return "Bad Xpath"
    case 7:
        return "Object not found"
    case 8:
        return "Object not unique"
    case 10:
        return "Reference count not zero"
    case 12:
        return "Invalid object"
    case 14:
        return "Operation not possible"
    case 15:
        return "Operation denied"
    case 16:
        return "Unauthorized"
    case 17:
        return "Invalid command"
    case 18:
        return "Malformed command"
    case 0, 19, 20:
        return ""
    case 22:
        return "Session timed out"
    default:
        return fmt.Sprintf("(%d) Unknown failure code, operation failed", e.ResponseCode)
    }
}

// panosErrorResponseWithLine is one of a few known error formats that PANOS
// outputs.  This has to be split from the other error struct because the
// the XML unmarshaler doesn't like a single struct to have overlapping
// definitions (the msg>line part).
type panosErrorResponseWithLine struct {
    XMLName xml.Name `xml:"response"`
    panosStatus
    ResponseMsg string `xml:"msg>line"`
}

// Error retrieves the parsed error message.
func (e panosErrorResponseWithLine) Error() string {
    if e.ResponseMsg != "" {
        return e.ResponseMsg
    } else {
        return e.codeError()
    }
}


// panosErrorResponseWithoutLine is one of a few known error formats that PANOS
// outputs.  It checks two locations that the error could be, and returns the
// one that was discovered in its Error().
type panosErrorResponseWithoutLine struct {
    XMLName xml.Name `xml:"response"`
    panosStatus
    ResponseMsg1 string `xml:"result>msg"`
    ResponseMsg2 string `xml:"msg"`
}

// Error retrieves the parsed error message.
func (e panosErrorResponseWithoutLine) Error() string {
    if e.ResponseMsg1 != "" {
        return e.ResponseMsg1
    } else {
        return e.ResponseMsg2
    }
}

// vis is a vsys import struct.
type vis struct {
    XMLName xml.Name
    Text string `xml:",chardata"`
}

type configLocks struct {
    Locks []util.Lock `xml:"result>config-locks>entry"`
}

type commitLocks struct {
    Locks []util.Lock `xml:"result>commit-locks>entry"`
}

type fwCommit struct {
    XMLName xml.Name `xml:"commit"`
    Description string `xml:"description,omitempty"`
    Partial *fwCommitPartial `xml:"partial"`
    Force interface{} `xml:"force"`
}

type fwCommitPartial struct {
    Dan string `xml:"device-and-network,omitempty"`
    Pao string `xml:"policy-and-objects,omitempty"`
}
