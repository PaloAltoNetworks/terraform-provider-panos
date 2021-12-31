package pango

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// These bit flags control what is logged by client connections.  Of the flags
// available for use, LogSend and LogReceive will log ALL communication between
// the connection object and the PAN-OS XML API.  The API key being used for
// communication will be blanked out, but no other sensitive data will be.  As
// such, those two flags should be considered for debugging only.  To disable
// all logging, set the logging level as LogQuiet.
//
// As of right now, pango is not officially supported by Palo Alto Networks TAC,
// however using the API itself via cURL is.  If you run into an issue and you believe
// it to be a PAN-OS problem, you can enable a cURL output logging style to have pango
// output an equivalent cURL command to use when interfacing with TAC.
//
// If you want to get the cURL command so that you can run it yourself, then set
// the LogCurlWithPersonalData flag, which will output your real API key, hostname,
// and any custom headers you have configured the client to send to PAN-OS.
//
// The bit-wise flags are as follows:
//
//      * LogQuiet: disables all logging
//      * LogAction: action being performed (Set / Edit / Delete functions)
//      * LogQuery: queries being run (Get / Show functions)
//      * LogOp: operation commands (Op functions)
//      * LogUid: User-Id commands (Uid functions)
//      * LogLog: log retrieval commands
//      * LogExport: log export commands
//      * LogXpath: the resultant xpath
//      * LogSend: xml docuemnt being sent
//      * LogReceive: xml responses being received
//      * LogOsxCurl: output an OSX cURL command for the data being sent in
//      * LogCurlWithPersonalData: If doing a curl style logging, then include
//        personal data in the curl command instead of tokens.
const (
	LogQuiet = 1 << (iota + 1)
	LogAction
	LogQuery
	LogOp
	LogUid
	LogLog
	LogExport
	LogImport
	LogXpath
	LogSend
	LogReceive
	LogOsxCurl
	LogCurlWithPersonalData
)

// Client is a generic connector struct.  It provides wrapper functions for
// invoking the various PAN-OS XPath API methods.  After creating the client,
// invoke Initialize() to prepare it for use.
//
// Many of the functions attached to this struct will take a param named
// `extras`.  Under normal circumstances this will just be nil, but if you have
// some extra values you need to send in with your request you can specify them
// here.
//
// Likewise, a lot of these functions will return a slice of bytes.  Under normal
// circumstances, you don't need to do anything with this, but sometimes you do,
// so you can find the raw XML returned from PAN-OS there.
type Client struct {
	// Connection properties.
	Hostname string            `json:"hostname"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	ApiKey   string            `json:"api_key"`
	Protocol string            `json:"protocol"`
	Port     uint              `json:"port"`
	Timeout  int               `json:"timeout"`
	Target   string            `json:"target"`
	Headers  map[string]string `json:"headers"`

	// Set to true if you want to check environment variables
	// for auth and connection properties.
	CheckEnvironment bool `json:"-"`

	// HTTP transport options.  Note that the VerifyCertificate setting is
	// only used if you do not specify a HTTP transport yourself.
	VerifyCertificate bool            `json:"verify_certificate"`
	Transport         *http.Transport `json:"-"`

	// Variables determined at runtime.
	Version        version.Number    `json:"-"`
	SystemInfo     map[string]string `json:"-"`
	Plugin         []plugin.Info     `json:"-"`
	MultiConfigure *MultiConfigure   `json:"-"`

	// Logging level.
	Logging               uint32   `json:"-"`
	LoggingFromInitialize []string `json:"logging"`

	// Internal variables.
	credsFile  string
	con        *http.Client
	api_url    string
	configTree *util.XmlNode

	// Variables for testing, response bytes, headers, and response index.
	rp              []url.Values
	rb              [][]byte
	rh              []http.Header
	ri              int
	authFileContent []byte
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

// Plugins returns the plugin information.
func (c *Client) Plugins() []plugin.Info {
	return c.Plugin
}

// Initialize does some initial setup of the Client connection, retrieves
// the API key if it was not already present, then performs "show system
// info" to get the PAN-OS version.  The full results are saved into the
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

// InitializeUsing does Initialize(), but takes in a filename that contains
// fallback authentication credentials if they aren't specified.
//
// The order of preference for auth / connection settings is:
//
// * explicitly set
// * environment variable (set chkenv to true to enable this)
// * json file
func (c *Client) InitializeUsing(filename string, chkenv bool) error {
	c.CheckEnvironment = chkenv
	c.credsFile = filename

	return c.Initialize()
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

	_, _, err := c.Communicate(data, &ans)
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
		e2, ok := err.(errors.Panos)
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
		e2, ok := err.(errors.Panos)
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
		Val     string   `xml:"password-hash>password"`
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

// ValidateConfig performs a commit config validation check.
//
// Setting sync to true means that this function will block until the job
// finishes.
//
//
// The sleep param is an optional sleep duration to wait between polling for
// job completion.  This param is only used if sync is set to true.
//
// This function returns the job ID and if any errors were encountered.
func (c *Client) ValidateConfig(sync bool, sleep time.Duration) (uint, error) {
	var err error

	c.LogOp("(op) validating config")
	type op_req struct {
		XMLName xml.Name `xml:"validate"`
		Cmd     string   `xml:"full"`
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

	return id, c.WaitForJob(id, sleep, nil, nil)
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
	var err error
	var cmd string
	ans := configLocks{}

	if vsys == "" {
		vsys = "shared"
	}

	if c.Version.Gte(version.Number{9, 1, 0, ""}) {
		var tgt string
		if vsys == "shared" {
			tgt = "all"
		} else {
			tgt = vsys
		}
		cmd = fmt.Sprintf("<show><config-locks><vsys>%s</vsys></config-locks></show>", tgt)
	} else {
		cmd = "<show><config-locks /></show>"
	}

	c.LogOp("(op) getting config locks for scope %q", vsys)
	_, err = c.Op(cmd, vsys, nil, &ans)
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
		Cmd     string   `xml:"config-lock>remove"`
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
		Admin   string   `xml:"commit-lock>remove>admin,omitempty"`
	}

	c.LogOp("(op) unlocking commits for scope %q", vsys)
	_, err := c.Op(cmd{Admin: admin}, vsys, nil, nil)
	return err
}

// WaitForJob polls the device, waiting for the specified job to finish.
//
// The sleep param is the length of time to wait between polling for job
// completion.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// If you want to unmarshal the response into a struct, then pass in a
// pointer to the struct for the "resp" param.  If you just want to know if
// the job completed with a status other than "FAIL", you only need to check
// the returned error message.
//
// In the case that there are multiple errors returned from the job, the first
// error is returned as the error string, and no unmarshaling is attempted.
func (c *Client) WaitForJob(id uint, sleep time.Duration, extras, resp interface{}) error {
	var err error
	var prev uint
	var data []byte
	dp := false
	all_ok := true

	c.LogOp("(op) waiting for job %d", id)
	type op_req struct {
		XMLName xml.Name `xml:"show"`
		Id      uint     `xml:"jobs>id"`
	}
	req := op_req{Id: id}

	var ans util.BasicJob
	for {
		// We need to zero out the response each iteration because the slices
		// of strings append to each other instead of zeroing out.
		ans = util.BasicJob{}

		// Get current percent complete.
		data, err = c.Op(req, "", extras, &ans)
		if err != nil {
			return err
		}

		// Output percent complete if it's new.
		if ans.Progress != prev {
			prev = ans.Progress
			c.LogOp("(op) job %d: %d percent complete", id, prev)
		}

		// Check for device commits.
		all_done := true
		for _, d := range ans.Devices {
			c.LogOp("%q result: %s", d.Serial, d.Result)
			if d.Result == "PEND" {
				all_done = false
				break
			} else if d.Result != "OK" && all_ok {
				all_ok = false
			}
		}

		// Check for end condition.
		if ans.Progress == 100 {
			if all_done {
				break
			} else if !dp {
				c.LogOp("(op) Waiting for %d device commits ...", len(ans.Devices))
				dp = true
			}
		}

		if sleep > 0 {
			time.Sleep(sleep)
		}
	}

	// Check the results for a failed commit.
	if ans.Result == "FAIL" {
		if len(ans.Details.Lines) > 0 {
			return fmt.Errorf(ans.Details.String())
		} else {
			return fmt.Errorf("Job %d has failed to complete successfully", id)
		}
	} else if !all_ok {
		return fmt.Errorf("Commit failed on one or more devices")
	}

	if resp == nil {
		return nil
	}

	return xml.Unmarshal(data, resp)
}

// LogAction writes a log message for SET/EDIT/DELETE operations if LogAction is set.
func (c *Client) LogAction(msg string, i ...interface{}) {
	if c.Logging&LogAction == LogAction {
		log.Printf(msg, i...)
	}
}

// LogQuery writes a log message for GET/SHOW operations if LogQuery is set.
func (c *Client) LogQuery(msg string, i ...interface{}) {
	if c.Logging&LogQuery == LogQuery {
		log.Printf(msg, i...)
	}
}

// LogOp writes a log message for OP operations if LogOp is set.
func (c *Client) LogOp(msg string, i ...interface{}) {
	if c.Logging&LogOp == LogOp {
		log.Printf(msg, i...)
	}
}

// LogUid writes a log message for User-Id operations if LogUid is set.
func (c *Client) LogUid(msg string, i ...interface{}) {
	if c.Logging&LogUid == LogUid {
		log.Printf(msg, i...)
	}
}

// LogLog writes a log message for LOG operations if LogLog is set.
func (c *Client) LogLog(msg string, i ...interface{}) {
	if c.Logging&LogLog == LogLog {
		log.Printf(msg, i...)
	}
}

// LogExport writes a log message for EXPORT operations if LogExport is set.
func (c *Client) LogExport(msg string, i ...interface{}) {
	if c.Logging&LogExport == LogExport {
		log.Printf(msg, i...)
	}
}

// LogImport writes a log message for IMPORT operations if LogImport is set.
func (c *Client) LogImport(msg string, i ...interface{}) {
	if c.Logging&LogImport == LogImport {
		log.Printf(msg, i...)
	}
}

// Communicate sends the given data to PAN-OS.
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
func (c *Client) Communicate(data url.Values, ans interface{}) ([]byte, http.Header, error) {
	if c.ApiKey != "" && data.Get("key") == "" {
		data.Set("key", c.ApiKey)
	}

	c.logSend(data)

	body, hdrs, err := c.post(data)
	if err != nil {
		return body, hdrs, err
	}

	return body, hdrs, c.endCommunication(body, ans)
}

// CommunicateFile does a file upload to PAN-OS.
//
// The content param is the content of the file you want to upload.
//
// The filename param is the basename of the file you want to specify in the
// multipart form upload.
//
// The fp param is the name of the param for the file upload.
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
func (c *Client) CommunicateFile(content, filename, fp string, data url.Values, ans interface{}) ([]byte, http.Header, error) {
	var err error

	if c.ApiKey != "" && data.Get("key") == "" {
		data.Set("key", c.ApiKey)
	}

	c.logSend(data)

	buf := bytes.Buffer{}
	w := multipart.NewWriter(&buf)

	for k := range data {
		w.WriteField(k, data.Get(k))
	}

	w2, err := w.CreateFormFile(fp, filename)
	if err != nil {
		return nil, nil, err
	}

	if _, err = io.Copy(w2, strings.NewReader(content)); err != nil {
		return nil, nil, err
	}

	w.Close()

	req, err := http.NewRequest("POST", c.api_url, &buf)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	res, err := c.con.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return body, res.Header, err
	}

	return body, res.Header, c.endCommunication(body, ans)
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

	b, _, err := c.Communicate(data, ans)
	return b, err
}

// Log submits a "log" command.
//
// Use `WaitForLogs` to get the results of the log command.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Log(logType, action, query, dir string, nlogs, skip int, extras, ans interface{}) ([]byte, error) {
	data := url.Values{}
	data.Set("type", "log")

	if logType != "" {
		data.Set("log-type", logType)
	}

	if action != "" {
		data.Set("action", action)
	}

	if query != "" {
		data.Set("query", query)
	}

	if dir != "" {
		data.Set("dir", dir)
	}

	if nlogs != 0 {
		data.Set("nlogs", strconv.Itoa(nlogs))
	}

	if skip != 0 {
		data.Set("skip", strconv.Itoa(skip))
	}

	if err := mergeUrlValues(&data, extras); err != nil {
		return nil, err
	}

	b, _, err := c.Communicate(data, ans)
	return b, err
}

// WaitForLogs performs repeated log retrieval operations until the log job is complete
// or the timeout is reached.
//
// Specify a timeout of zero to wait indefinitely.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) WaitForLogs(id uint, sleep, timeout time.Duration, ans interface{}) ([]byte, error) {
	var err error
	var data []byte
	var prev string
	start := time.Now()
	end := start.Add(timeout)
	extras := url.Values{}
	extras.Set("job-id", fmt.Sprintf("%d", id))

	c.LogLog("(log) waiting for logs: %d", id)

	var resp util.BasicJob
	for {
		resp = util.BasicJob{}

		data, err = c.Log("", "get", "", "", 0, 0, extras, &resp)
		if err != nil {
			return data, err
		}

		if resp.Status != prev {
			prev = resp.Status
			c.LogLog("(log) job %d status: %s", id, prev)
		}

		if resp.Status == "FIN" {
			break
		}

		if timeout > 0 && end.After(time.Now()) {
			return data, fmt.Errorf("timeout")
		}

		if sleep > 0 {
			time.Sleep(sleep)
		}
	}

	if resp.Result == "FAIL" {
		if len(resp.Details.Lines) > 0 {
			return data, fmt.Errorf(resp.Details.String())
		} else {
			return data, fmt.Errorf("Job %d has failed to complete successfully", id)
		}
	}

	if ans == nil {
		return data, nil
	}

	err = xml.Unmarshal(data, ans)
	return data, err
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

	return c.typeConfig("show", data, nil, extras, ans)
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

	return c.typeConfig("get", data, nil, extras, ans)
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

	return c.typeConfig("delete", data, nil, extras, ans)
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
	data := url.Values{}
	xp := util.AsXpath(path)
	c.logXpath(xp)
	data.Set("xpath", xp)

	return c.typeConfig("set", data, element, extras, ans)
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
	data := url.Values{}
	xp := util.AsXpath(path)
	c.logXpath(xp)
	data.Set("xpath", xp)

	return c.typeConfig("edit", data, element, extras, ans)
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

	return c.typeConfig("move", data, nil, extras, ans)
}

// Rename does a "rename" type command.
func (c *Client) Rename(path interface{}, newname string, extras, ans interface{}) ([]byte, error) {
	data := url.Values{}
	xp := util.AsXpath(path)
	c.logXpath(xp)
	data.Set("xpath", xp)
	data.Set("newname", newname)

	return c.typeConfig("rename", data, nil, extras, ans)
}

// MultiConfig does a "multi-config" type command.
//
// Param strict should be true if you want strict transactional support.
//
// Note that the error returned from this function is only if there was an error
// unmarshaling the response into the the multi config response struct.  If the
// multi config itself failed, then the reason can be found in its results.
func (c *Client) MultiConfig(element MultiConfigure, strict bool, extras interface{}) ([]byte, MultiConfigureResponse, error) {
	data := url.Values{}
	if strict {
		data.Set("strict-transactional", "yes")
	}

	text, _ := c.typeConfig("multi-config", data, element, extras, nil)

	resp := MultiConfigureResponse{}
	err := xml.Unmarshal(text, &resp)
	return text, resp, err
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

	b, _, err := c.Communicate(data, ans)
	return b, err
}

// Import performs an import type command.
//
// The cat param is the category.
//
// The content param is the content of the file you want to upload.
//
// The filename param is the basename of the file you want to specify in the
// multipart form upload.
//
// The fp param is the name of the param for the file upload.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
func (c *Client) Import(cat, content, filename, fp string, timeout time.Duration, extras, ans interface{}) ([]byte, error) {
	if timeout < 0 {
		return nil, fmt.Errorf("timeout cannot be negative")
	} else if timeout > 0 {
		defer func(c *Client, v time.Duration) {
			c.con.Timeout = v
		}(c, c.con.Timeout)
		c.con.Timeout = timeout
	}

	data := url.Values{}
	data.Set("type", "import")
	data.Set("category", cat)

	if err := mergeUrlValues(&data, extras); err != nil {
		return nil, err
	}

	b, _, err := c.CommunicateFile(content, filename, fp, data, ans)
	return b, err
}

// Commit performs PAN-OS commits.
//
// The cmd param can be a properly formatted XML string, a struct that can
// be marshalled into XML, or one of the commit types that can be found in the
// commit package.
//
// The action param is the commit action to be taken.  If you are using one of the
// commit structs as the `cmd` param and the action param is an empty string, then
// the action is taken from the commit struct passed in.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// Commits result in a job being submitted to the backend.  The job ID, assuming
// the commit action was successfully submitted, the response from the server,
// and if an error was encountered or not are all returned from this function.
func (c *Client) Commit(cmd interface{}, action string, extras interface{}) (uint, []byte, error) {
	var err error
	data := url.Values{}
	data.Set("type", "commit")

	if err = addToData("cmd", cmd, true, &data); err != nil {
		return 0, nil, err
	}

	if action != "" {
		data.Set("action", action)
	} else if ca, ok := cmd.(util.Actioner); ok && ca.Action() != "" {
		data.Set("action", ca.Action())
	}

	if c.Target != "" {
		data.Set("target", c.Target)
	}

	if err = mergeUrlValues(&data, extras); err != nil {
		return 0, nil, err
	}

	ans := util.JobResponse{}
	b, _, err := c.Communicate(data, &ans)
	return ans.Id, b, err
}

// Export runs an "export" type command.
//
// The category param specifies the desired file type to export.
//
// The extras param should be either nil or a url.Values{} to be mixed in with
// the constructed request.
//
// The ans param should be a pointer to a struct to unmarshal the response
// into or nil.
//
// Any response received from the server is returned, along with any errors
// encountered.
//
// If the export invoked results in a file being downloaded from PAN-OS, then
// the string returned is the name of the remote file that is retrieved,
// otherwise it's just an empty string.
func (c *Client) Export(category string, timeout time.Duration, extras, ans interface{}) (string, []byte, error) {
	if timeout < 0 {
		return "", nil, fmt.Errorf("timeout cannot be negative")
	} else if timeout > 0 {
		defer func(c *Client, v time.Duration) {
			c.con.Timeout = v
		}(c, c.con.Timeout)
		c.con.Timeout = timeout
	}

	data := url.Values{}
	data.Set("type", "export")

	if category != "" {
		data.Set("category", category)
	}

	if err := mergeUrlValues(&data, extras); err != nil {
		return "", nil, err
	}

	var filename string
	b, hdrs, err := c.Communicate(data, ans)
	if err == nil && hdrs != nil {
		// Check and see if there's a filename in the content disposition.
		mediatype, params, err := mime.ParseMediaType(hdrs.Get("Content-Disposition"))
		if err == nil && mediatype == "attachment" {
			filename = params["filename"]
		}
	}

	return filename, b, err
}

/*** Internal functions ***/

func (c *Client) initCon() error {
	var tout time.Duration

	// Load up the JSON config file.
	json_client := &Client{}
	if c.credsFile != "" {
		var (
			b   []byte
			err error
		)
		if len(c.rb) == 0 {
			b, err = ioutil.ReadFile(c.credsFile)
		} else {
			b, err = c.authFileContent, nil
		}

		if err != nil {
			return err
		}

		if err = json.Unmarshal(b, &json_client); err != nil {
			return err
		}
	}

	// Hostname.
	if c.Hostname == "" {
		if val := os.Getenv("PANOS_HOSTNAME"); c.CheckEnvironment && val != "" {
			c.Hostname = val
		} else {
			c.Hostname = json_client.Hostname
		}
	}

	// Username.
	if c.Username == "" {
		if val := os.Getenv("PANOS_USERNAME"); c.CheckEnvironment && val != "" {
			c.Username = val
		} else {
			c.Username = json_client.Username
		}
	}

	// Password.
	if c.Password == "" {
		if val := os.Getenv("PANOS_PASSWORD"); c.CheckEnvironment && val != "" {
			c.Password = val
		} else {
			c.Password = json_client.Password
		}
	}

	// API key.
	if c.ApiKey == "" {
		if val := os.Getenv("PANOS_API_KEY"); c.CheckEnvironment && val != "" {
			c.ApiKey = val
		} else {
			c.ApiKey = json_client.ApiKey
		}
	}

	// Protocol.
	if c.Protocol == "" {
		if val := os.Getenv("PANOS_PROTOCOL"); c.CheckEnvironment && val != "" {
			c.Protocol = val
		} else if json_client.Protocol != "" {
			c.Protocol = json_client.Protocol
		} else {
			c.Protocol = "https"
		}
	}
	if c.Protocol != "http" && c.Protocol != "https" {
		return fmt.Errorf("Invalid protocol %q.  Must be \"http\" or \"https\"", c.Protocol)
	}

	// Port.
	if c.Port == 0 {
		if val := os.Getenv("PANOS_PORT"); c.CheckEnvironment && val != "" {
			if cp, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("Failed to parse the env port number: %s", err)
			} else {
				c.Port = uint(cp)
			}
		} else if json_client.Port != 0 {
			c.Port = json_client.Port
		}
	}
	if c.Port > 65535 {
		return fmt.Errorf("Port %d is out of bounds", c.Port)
	}

	// Timeout.
	if c.Timeout == 0 {
		if val := os.Getenv("PANOS_TIMEOUT"); c.CheckEnvironment && val != "" {
			if ival, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("Failed to parse timeout env var as int: %s", err)
			} else {
				c.Timeout = ival
			}
		} else if json_client.Timeout != 0 {
			c.Timeout = json_client.Timeout
		} else {
			c.Timeout = 10
		}
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("Timeout for %q must be a positive int", c.Hostname)
	}
	tout = time.Duration(time.Duration(c.Timeout) * time.Second)

	// Target.
	if c.Target == "" {
		if val := os.Getenv("PANOS_TARGET"); c.CheckEnvironment && val != "" {
			c.Target = val
		} else {
			c.Target = json_client.Target
		}
	}

	// Headers.
	if len(c.Headers) == 0 {
		if val := os.Getenv("PANOS_HEADERS"); c.CheckEnvironment && val != "" {
			if err := json.Unmarshal([]byte(val), &c.Headers); err != nil {
				return err
			}
		}
		if len(c.Headers) == 0 && len(json_client.Headers) > 0 {
			c.Headers = make(map[string]string)
			for k, v := range json_client.Headers {
				c.Headers[k] = v
			}
		}
	}

	// Verify cert.
	if !c.VerifyCertificate {
		if val := os.Getenv("PANOS_VERIFY_CERTIFICATE"); c.CheckEnvironment && val != "" {
			if vcb, err := strconv.ParseBool(val); err != nil {
				return err
			} else if vcb {
				c.VerifyCertificate = vcb
			}
		}
		if !c.VerifyCertificate && json_client.VerifyCertificate {
			c.VerifyCertificate = json_client.VerifyCertificate
		}
	}

	// Logging.
	if c.Logging == 0 {
		var ll []string
		if val := os.Getenv("PANOS_LOGGING"); c.CheckEnvironment && val != "" {
			ll = strings.Split(val, ",")
		} else {
			ll = json_client.LoggingFromInitialize
		}
		if len(ll) > 0 {
			var lv uint32
			for _, x := range ll {
				switch x {
				case "quiet":
					lv |= LogQuiet
				case "action":
					lv |= LogAction
				case "query":
					lv |= LogQuery
				case "op":
					lv |= LogOp
				case "uid":
					lv |= LogUid
				case "xpath":
					lv |= LogXpath
				case "send":
					lv |= LogSend
				case "receive":
					lv |= LogReceive
				default:
					return fmt.Errorf("Unknown logging requested: %s", x)
				}
			}
			c.Logging = lv
		} else {
			c.Logging = LogAction | LogUid
		}
	}

	// Setup the https client.
	if c.Transport == nil {
		c.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !c.VerifyCertificate,
			},
		}
	}
	c.con = &http.Client{
		Transport: c.Transport,
		Timeout:   tout,
	}

	// Sanity check.
	if c.Hostname == "" {
		return fmt.Errorf("No hostname specified")
	} else if c.ApiKey == "" && (c.Username == "" && c.Password == "") {
		return fmt.Errorf("No username/password or API key given")
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
		Cmd     string   `xml:"system>info"`
	}

	type tagVal struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}

	type sysTag struct {
		XMLName xml.Name `xml:"system"`
		Tag     []tagVal `xml:",any"`
	}

	type system_info_ans struct {
		System sysTag `xml:"result>system"`
	}

	req := system_info_req{}
	ans := system_info_ans{}

	_, err = c.Op(req, "", nil, &ans)
	if err != nil {
		return err
	}

	c.SystemInfo = make(map[string]string, len(ans.System.Tag))
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

func (c *Client) initPlugins() {
	c.LogOp("(op) getting plugin info")

	var req plugin.GetPlugins
	var ans plugin.PackageListing

	if _, err := c.Op(req, "", nil, &ans); err != nil {
		c.LogAction("WARNING: Failed to get plugin info: %s", err)
		return
	}

	c.Plugin = ans.Listing()
}

func (c *Client) typeConfig(action string, data url.Values, element, extras, ans interface{}) ([]byte, error) {
	var err error

	if c.MultiConfigure != nil && (action == "set" ||
		action == "edit" ||
		action == "delete") {
		r := MultiConfigureRequest{
			XMLName: xml.Name{Local: action},
			Xpath:   data.Get("xpath"),
		}
		if element != nil {
			r.Data = element
		}
		c.MultiConfigure.Reqs = append(c.MultiConfigure.Reqs, r)
		return nil, nil
	}

	data.Set("type", "config")
	data.Set("action", action)

	if element != nil {
		if err = addToData("element", element, true, &data); err != nil {
			return nil, err
		}
	}

	if c.Target != "" {
		data.Set("target", c.Target)
	}

	if err = mergeUrlValues(&data, extras); err != nil {
		return nil, err
	}

	b, _, err := c.Communicate(data, ans)
	return b, err
}

func (c *Client) logXpath(p string) {
	if c.Logging&LogXpath == LogXpath {
		log.Printf("(xpath) %s", p)
	}
}

// VsysImport imports the given names into the specified template / vsys.
func (c *Client) VsysImport(loc, tmpl, ts, vsys string, names []string) error {
	path := c.xpathImport(tmpl, ts, vsys)
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

// VsysUnimport removes the given names from all (template, optional) vsys.
func (c *Client) VsysUnimport(loc, tmpl, ts string, names []string) error {
	if len(names) == 0 {
		return nil
	}

	path := make([]string, 0, 14)
	path = append(path, c.xpathImport(tmpl, ts, "")...)
	path = append(path, loc, util.AsMemberXpath(names))

	_, err := c.Delete(path, nil, nil)
	if err != nil {
		e2, ok := err.(errors.Panos)
		if ok && e2.ObjectNotFound() {
			return nil
		}
	}
	return err
}

// IsImported checks if the importable object is actually imported in the
// specified location.
func (c *Client) IsImported(loc, tmpl, ts, vsys, name string) (bool, error) {
	path := make([]string, 0, 14)
	path = append(path, c.xpathImport(tmpl, ts, vsys)...)
	path = append(path, loc, util.AsMemberXpath([]string{name}))

	_, err := c.Get(path, nil, nil)
	if err == nil {
		if vsys != "" {
			return true, nil
		} else {
			return false, nil
		}
	}

	e2, ok := err.(errors.Panos)
	if ok && e2.ObjectNotFound() {
		if vsys != "" {
			return false, nil
		} else {
			return true, nil
		}
	}

	return false, err
}

func (c *Client) xpathImport(tmpl, ts, vsys string) []string {
	ans := make([]string, 0, 12)
	if tmpl != "" || ts != "" {
		ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
	}
	ans = append(ans,
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"vsys",
		util.AsEntryXpath([]string{vsys}),
		"import",
		"network",
	)

	return ans
}

func (c *Client) post(data url.Values) ([]byte, http.Header, error) {
	if len(c.rb) == 0 {
		req, err := http.NewRequest("POST", c.api_url, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}

		r, err := c.con.Do(req)
		if err != nil {
			return nil, nil, err
		}

		defer r.Body.Close()
		ans, err := ioutil.ReadAll(r.Body)
		return ans, r.Header, err
	} else {
		if c.ri < len(c.rb) {
			c.rp = append(c.rp, data)
		}
		body := c.rb[c.ri%len(c.rb)]
		var hdr http.Header
		if len(c.rh) > 0 {
			hdr = c.rh[c.ri%len(c.rh)]
		}
		c.ri++
		return body, hdr, nil
	}
}

func (c *Client) endCommunication(body []byte, ans interface{}) error {
	var err error

	if c.Logging&LogReceive == LogReceive {
		log.Printf("Response = %s", body)
	}

	// Check for errors first
	if err = errors.Parse(body); err != nil {
		return err
	}

	// Return the body string if we weren't given something to unmarshal into
	if ans == nil {
		return nil
	}

	// Unmarshal using the struct passed in
	err = xml.Unmarshal(body, ans)
	if err != nil {
		return fmt.Errorf("Error unmarshaling into provided interface: %s", err)
	}

	return nil
}

/*
PositionFirstEntity moves an element before another one using the Move API command.

Param `mvt` is a util.Move* constant.

Param `rel` is the relative entity that `mvt` is in relation to.

Param `ent` is the entity that is to be positioned.

Param `path` is the XPATH of `ent`.

Param `elms` is the ordered list of entities that should include both
`rel` and `ent`.
be found.
*/
func (c *Client) PositionFirstEntity(mvt int, rel, ent string, path, elms []string) error {
	// Sanity checks.
	if rel == ent {
		return fmt.Errorf("Can't position %q in relation to itself", rel)
	} else if mvt < util.MoveSkip && mvt > util.MoveBottom {
		return fmt.Errorf("Invalid position int given: %d", mvt)
	} else if (mvt == util.MoveBefore || mvt == util.MoveDirectlyBefore || mvt == util.MoveAfter || mvt == util.MoveDirectlyAfter) && rel == "" {
		return fmt.Errorf("Specify 'ref' in order to perform relative group positioning")
	}

	var err error
	fIdx := -1
	oIdx := -1

	switch mvt {
	case util.MoveSkip:
		return nil
	case util.MoveTop:
		_, em := c.Move(path, "top", "", nil, nil)
		if em != nil && em.Error() != "already at the top" {
			err = em
		}
	case util.MoveBottom:
		_, em := c.Move(path, "bottom", "", nil, nil)
		if em != nil && em.Error() != "already at the bottom" {
			err = em
		}
	default:
		// Find the indexes of the first rule and the ref rule.
		for i, v := range elms {
			if v == ent {
				fIdx = i
			} else if v == rel {
				oIdx = i
			}
			if fIdx != -1 && oIdx != -1 {
				break
			}
		}

		// Sanity check: both rules should be present.
		if fIdx == -1 {
			return fmt.Errorf("Entity to be moved %q does not exist", ent)
		} else if oIdx == -1 {
			return fmt.Errorf("Reference entity %q does not exist", rel)
		}

		// Move the first element, if needed.
		if (mvt == util.MoveBefore && fIdx > oIdx) || (mvt == util.MoveDirectlyBefore && fIdx+1 != oIdx) {
			_, err = c.Move(path, "before", rel, nil, nil)
		} else if (mvt == util.MoveAfter && fIdx < oIdx) || (mvt == util.MoveDirectlyAfter && fIdx != oIdx+1) {
			_, err = c.Move(path, "after", rel, nil, nil)
		}
	}

	return err
}

// Clock gets the time on the PAN-OS appliance.
func (c *Client) Clock() (time.Time, error) {
	type t_req struct {
		XMLName xml.Name `xml:"show"`
		Cmd     string   `xml:"clock"`
	}

	type t_resp struct {
		Result string `xml:"result"`
	}

	req := t_req{}
	ans := t_resp{}

	c.LogOp("(op) getting system time")
	if _, err := c.Op(req, "", nil, &ans); err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.UnixDate+"\n", ans.Result)
}

// PrepareMultiConfigure will start a multi config command.
//
// Capacity is the initial capacity of the requests to be sent.
func (c *Client) PrepareMultiConfigure(capacity int) {
	c.MultiConfigure = &MultiConfigure{
		Reqs: make([]MultiConfigureRequest, 0, capacity),
	}
}

// SendMultiConfigure will send the accumulated multi configure request.
//
// Param strict should be true if you want strict transactional support.
//
// Note that the error returned from this function is only if there was an error
// unmarshaling the response into the the multi config response struct.  If the
// multi config itself failed, then the reason can be found in its results.
func (c *Client) SendMultiConfigure(strict bool) (MultiConfigureResponse, error) {
	if c.MultiConfigure == nil {
		return MultiConfigureResponse{}, nil
	}

	mc := c.MultiConfigure
	c.MultiConfigure = nil

	_, ans, err := c.MultiConfig(*mc, strict, nil)
	return ans, err
}

// GetTechSupportFile returns the tech support .tgz file.
//
// This function returns the name of the tech support file, the file
// contents, and an error if one occurred.
//
// The timeout param is the new timeout (in seconds) to temporarily assign to
// client connections to allow for the successful download of the tech support
// file.  If the timeout is zero, then pango.Client.Timeout is the timeout for
// tech support file retrieval.
func (c *Client) GetTechSupportFile(timeout time.Duration) (string, []byte, error) {
	if timeout < 0 {
		return "", nil, fmt.Errorf("timeout cannot be negative")
	}

	var err error
	var resp util.JobResponse
	cmd := "tech-support"

	c.LogExport("(export) tech support file")

	// Request the tech support file creation.
	_, _, err = c.Export(cmd, 0, nil, &resp)
	if err != nil {
		return "", nil, err
	}
	if resp.Id == 0 {
		return "", nil, fmt.Errorf("Job ID was not found")
	}

	extras := url.Values{}
	extras.Set("action", "status")
	extras.Set("job-id", fmt.Sprintf("%d", resp.Id))

	// Poll the job until it's done.
	var pr util.BasicJob
	var prev uint
	for {
		_, _, err = c.Export(cmd, 0, extras, &pr)
		if err != nil {
			return "", nil, err
		}

		// The progress is not an uint when the job completes, so don't print
		// the progress as 0 when the job is actually complete.
		if pr.Progress != prev && pr.Progress != 0 {
			prev = pr.Progress
			c.LogExport("(export) tech support job %d: %d percent complete", resp.Id, prev)
		}

		if pr.Status == "FIN" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	if pr.Result == "FAIL" {
		return "", nil, fmt.Errorf(pr.Details.String())
	}

	extras.Set("action", "get")
	return c.Export(cmd, timeout, extras, nil)
}

// RetrievePanosConfig retrieves either the running config, candidate config,
// or the specified saved config file, then does `LoadPanosConfig()` to save it.
//
// After the config is loaded, config can be queried and retrieved using
// any `FromPanosConfig()` methods.
//
// Param `value` can be the word "candidate" to load candidate config or
// `running` to load running config.  If the value is neither of those, it
// is assumed to be the name of a saved config and that is loaded.
func (c *Client) RetrievePanosConfig(value string) error {
	type getConfig struct {
		XMLName   xml.Name `xml:"show"`
		Running   *string  `xml:"config>running"`
		Candidate *string  `xml:"config>candidate"`
		Saved     *string  `xml:"config>saved"`
	}

	type data struct {
		Data []byte `xml:",innerxml"`
	}

	type resp struct {
		XMLName xml.Name `xml:"response"`
		Result  data     `xml:"result"`
	}

	s := ""
	req := getConfig{}
	switch value {
	case "candidate":
		req.Candidate = &s
	case "running":
		req.Running = &s
	default:
		req.Saved = &value
	}
	ans := resp{}

	if _, err := c.Op(req, "", nil, &ans); err != nil {
		return err
	}

	return c.LoadPanosConfig(ans.Result.Data)
}

// LoadPanosConfig stores the given XML document into the local client instance.
//
// The `config` can either be `<config>...</config>` or something that contians
// only the config document (such as `<result ...><config>...</config></result>`).
//
// After the config is loaded, config can be queried and retrieved using
// any `FromPanosConfig()` methods.
func (c *Client) LoadPanosConfig(config []byte) error {
	log.Printf("load panos config")
	if err := xml.Unmarshal(config, &c.configTree); err != nil {
		return err
	}

	if c.configTree.XMLName.Local == "config" {
		// Add a place holder parent util.XmlNode.
		c.configTree = &util.XmlNode{
			XMLName: xml.Name{
				Local: "a",
			},
			Nodes: []util.XmlNode{
				*c.configTree,
			},
		}
		return nil
	}

	if len(c.configTree.Nodes) == 1 && c.configTree.Nodes[0].XMLName.Local == "config" {
		// Already has a place holder parent.
		return nil
	}

	c.configTree = nil
	return fmt.Errorf("doesn't seem to be a config tree")
}

// ConfigTree returns the configuration tree that was loaded either via
// `RetrievePanosConfig()` or `LoadPanosConfig()`.
func (c *Client) ConfigTree() *util.XmlNode {
	return c.configTree
}

func (c *Client) logSend(data url.Values) {
	var b strings.Builder

	// Traditional send logging.
	if c.Logging&LogSend == LogSend {
		if b.Len() > 0 {
			fmt.Fprintf(&b, "\n")
		}
		realKey := data.Get("key")
		if realKey != "" {
			data.Set("key", "########")
		}
		fmt.Fprintf(&b, "Sending data: %#v", data)
		if realKey != "" {
			data.Set("key", realKey)
		}
	}

	// Log the send data as an OSX curl command.
	if c.Logging&LogOsxCurl == LogOsxCurl {
		if b.Len() > 0 {
			fmt.Fprintf(&b, "\n")
		}
		special := map[string]string{
			"key":     "",
			"element": "",
		}
		ev := url.Values{}
		for k := range data {
			var isSpecial bool
			for sk := range special {
				if sk == k {
					isSpecial = true
					special[k] = data.Get(k)
					break
				}
			}
			if !isSpecial {
				ev[k] = make([]string, 0, len(data[k]))
				for i := range data[k] {
					ev[k] = append(ev[k], data[k][i])
				}
			}
		}

		// Build up the curl command.
		fmt.Fprintf(&b, "curl")
		// Verify cert.
		if !c.VerifyCertificate {
			fmt.Fprintf(&b, " -k")
		}
		// Headers.
		if len(c.Headers) > 0 && c.Logging&LogCurlWithPersonalData == LogCurlWithPersonalData {
			for k, v := range c.Headers {
				if v != "" {
					fmt.Fprintf(&b, " --header '%s: %s'", k, v)
				} else {
					fmt.Fprintf(&b, " --header '%s;'", k)
				}
			}
		}
		// Add URL encoded values.
		if special["key"] != "" {
			if c.Logging&LogCurlWithPersonalData == LogCurlWithPersonalData {
				ev.Set("key", special["key"])
			} else {
				ev.Set("key", "APIKEY")
			}
		}
		// Add in the element, if present.
		if special["element"] != "" {
			fmt.Fprintf(&b, " --data-urlencode element@element.xml")
		}
		// URL.
		fmt.Fprintf(&b, " '%s://", c.Protocol)
		if c.Logging&LogCurlWithPersonalData == LogCurlWithPersonalData {
			fmt.Fprintf(&b, "%s", c.Hostname)
		} else {
			fmt.Fprintf(&b, "HOST")
		}
		if c.Port != 0 {
			fmt.Fprintf(&b, ":%d", c.Port)
		}
		fmt.Fprintf(&b, "/api")
		if len(ev) > 0 {
			fmt.Fprintf(&b, "?%s", ev.Encode())
		}
		fmt.Fprintf(&b, "'")
		// Data.
		if special["element"] != "" {
			fmt.Fprintf(&b, "\nelement.xml:\n%s", special["element"])
		}
	}

	if b.Len() > 0 {
		log.Printf("%s", b.String())
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
	if a, ok := i.(fmt.Stringer); ok {
		return a.String(), nil
	}

	if b, ok := i.(util.Elementer); ok {
		i = b.Element()
	}

	switch val := i.(type) {
	case nil:
		return "", fmt.Errorf("nil encountered")
	case string:
		return val, nil
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

// vis is a vsys import struct.
type vis struct {
	XMLName xml.Name
	Text    string `xml:",chardata"`
}

type configLocks struct {
	Locks []util.Lock `xml:"result>config-locks>entry"`
}

type commitLocks struct {
	Locks []util.Lock `xml:"result>commit-locks>entry"`
}
