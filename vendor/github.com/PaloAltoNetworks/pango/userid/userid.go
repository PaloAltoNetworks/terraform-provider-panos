// Package userid is the client.UserId namespace, for interacting with the
// User-ID API.  This includes login/logout of a user, user/group mappings,
// and dynamic address group tags.
//
// Various features of User-ID API are supported across all versions of PANOS
// for the firewall, but User-ID API for Panorama was only added to PANOS
// version 8.0.
package userid

import (
    "fmt"
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
    "github.com/PaloAltoNetworks/pango/version"
)


// UserId is the client.UserId namespace.
type UserId struct {
    con util.XapiClient
}

// Initialize is invoked on client.Initialize().
func (c *UserId) Initialize(i util.XapiClient) {
    c.con = i
}

// Run executes the given User-Id related operations.  This allows you to
// perform the following User-Id operations:
//
//  * login users
//  * logout users
//  * register ip/tags
//  * unregister ip/tags
//
// Both logins and logouts are maps where the username is the key and the IP
// address is the value.
//
// Both reg and unreg are maps where the IP address is the key and the list
// of tags are the values.
//
// The vsys param is which vsys these operations should take place in.  If
// vsys is an empty string, vsys defaults to "vsys1".
func (c *UserId) Run(logins, logouts map[string] string, reg, unreg map[string] []string, vsys string) error {
    var i int
    if vsys == "" {
        vsys = "vsys1"
    }
    c.con.LogUid("(userid) running in %s - logins:%d logouts:%d reg:%d unreg:%d", vsys, len(logins), len(logouts), len(reg), len(unreg))

    msg := uid{Version: "1.0", Type: "update"}

    if len(logins) == 0 && len(logouts) == 0 && len(reg) == 0 && len(unreg) == 0 {
        return nil
    }

    // Login users.
    if len(logins) > 0 {
        i = 0
        msg.Payload.Login = &inOutCon{}
        msg.Payload.Login.Entry = make([]inOut, len(logins))
        for k, v := range logins {
            msg.Payload.Login.Entry[i] = inOut{Name: k, Ip: v}
            i++
        }
    }

    // Logout users.
    if len(logouts) > 0 {
        i = 0
        msg.Payload.Logout = &inOutCon{}
        msg.Payload.Logout.Entry = make([]inOut, len(logouts))
        for k, v := range logouts {
            msg.Payload.Logout.Entry[i] = inOut{Name: k, Ip: v}
            i++
        }
    }

    // Register ip/tags.
    if len(reg) > 0 {
        i = 0
        msg.Payload.Register = &regUnregCon{}
        msg.Payload.Register.Entry = make([]regUnreg, len(reg))
        for ip, tags := range reg {
            msg.Payload.Register.Entry[i] = regUnreg{Ip: ip, Tag: tags}
            i++
        }
    }

    // Unregister ip/tags.
    if len(unreg) > 0 {
        i = 0
        msg.Payload.Unregister = &regUnregCon{}
        msg.Payload.Unregister.Entry = make([]regUnreg, len(unreg))
        for ip, tags := range unreg {
            msg.Payload.Unregister.Entry[i] = regUnreg{Ip: ip, Tag: tags}
            i++
        }
    }

    _, err := c.con.Uid(msg, vsys, nil, nil)
    return err
}

// Registered returns the registered IP address / tags for the given vsys.
//
// Both the ip and tag params are server-side filters.
//
// The vsys param is which vsys these operations should take place in.  If
// vsys is an empty string, vsys defaults to "vsys1".
func (c *UserId) Registered(ip, tag, vsys string) (map[string] []string, error) {
    if vsys == "" {
        vsys = "vsys1"
    }
    c.con.LogOp("(op) getting registered ip addresses - ip:%q tag:%q vsys:%q", ip, tag, vsys)
    req := c.versioning()

    ans := make(map[string] []string)
    for {
        req.FilterOn(ip, tag, len(ans))
        resp := regResp{}

        _, err := c.con.Op(req, vsys, nil, &resp)
        if err != nil {
            return nil, err
        } else if resp.Msg != nil && resp.Msg.Outfile != "" {
            return nil, fmt.Errorf("PAN-OS returned %q instead of IP/tag mappings, please upgrade to 8.0+", resp.Msg.Outfile)
        }

        for i := range resp.Entry {
            ans[resp.Entry[i].Ip] = resp.Entry[i].Tags
        }

        if req.ShouldStop(len(resp.Entry)) {
            break
        }
    }

    return ans, nil
}

/** Internal functions for the UserId struct **/

func (c *UserId) versioning() filterer {
    v := c.con.Versioning()

    if v.Gte(version.Number{8, 0, 0, ""}) {
        return &req_v3{}
    } else if v.Gte(version.Number{6, 1, 0, ""}) {
        return &req_v2{}
    } else {
        return &req_v1{}
    }
}

/** Structs / functions for this namespace. **/

type uid struct {
    XMLName xml.Name `xml:"uid-message"`
    Version string `xml:"version"`
    Type string `xml:"type"`
    Payload payload `xml:"payload"`
}

type payload struct {
    Login *inOutCon `xml:"login"`
    Logout *inOutCon `xml:"logout"`
    Register *regUnregCon `xml:"register"`
    Unregister *regUnregCon `xml:"unregister"`
}

type inOutCon struct {
    Entry []inOut `xml:"entry"`
}

type inOut struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Ip string `xml:"ip,attr"`
}

type regUnregCon struct {
    Entry []regUnreg `xml:"entry"`
}

type regUnreg struct {
    XMLName xml.Name `xml:"entry"`
    Ip string `xml:"ip,attr"`
    Tag []string `xml:"tag>member"`
}

type filterer interface {
    FilterOn(string, string, int)
    ShouldStop(int) bool
}

type req_filter struct {
    Tag *tagFilter `xml:"tag"`
    Ip string `xml:"ip,omitempty"`
}

type tagFilter struct {
    Entry tagName `xml:"entry"`
}

type tagName struct {
    Name string `xml:"name,attr"`
}

type req_v1 struct {
    XMLName xml.Name `xml:"show"`
    Filter req_filter `xml:"object>registered-address"`
}

func (o *req_v1) FilterOn(ip, tag string, size int) {
    o.Filter.Ip = ip
    if tag != "" {
        o.Filter.Tag = &tagFilter{tagName{tag}}
    }
}

func (o *req_v1) ShouldStop(lastCount int) bool {
    return true
}

type req_v2 struct {
    XMLName xml.Name `xml:"show"`
    Filter req_filter `xml:"object>registered-ip"`
}

func (o *req_v2) FilterOn(ip, tag string, size int) {
    o.Filter.Ip = ip
    if tag != "" {
        o.Filter.Tag = &tagFilter{tagName{tag}}
    }
}

func (o *req_v2) ShouldStop(lastCount int) bool {
    return true
}

type req_v3 struct {
    XMLName xml.Name `xml:"show"`
    Filter req_filter_v2 `xml:"object>registered-ip"`
}

func (o *req_v3) FilterOn(ip, tag string, size int) {
    o.Filter.Ip = ip
    if tag != "" {
        o.Filter.Tag = &tagFilter{tagName{tag}}
    }

    o.Filter.Limit = 500
    o.Filter.Start = size + 1
}

func (o *req_v3) ShouldStop(lastCount int) bool {
    return lastCount < o.Filter.Limit
}

type req_filter_v2 struct {
    Tag *tagFilter `xml:"tag"`
    Ip string `xml:"ip,omitempty"`
    Limit int `xml:"limit"`
    Start int `xml:"start-point"`
}

type regResp struct {
    Entry []respEntry `xml:"result>entry"`
    Msg *msg `xml:"result>msg"`
}

type respEntry struct {
    Ip string `xml:"ip,attr"`
    Tags []string `xml:"tag>member"`
}

type msg struct {
    Outfile string `xml:"line>outfile"`
}
