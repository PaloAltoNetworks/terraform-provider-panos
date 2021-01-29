package userid

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/version"
)

// GetIpTags returns the registered IP address / tags for the given vsys.
//
// Both the ip and tag params are server-side filters.
//
// The vsys param is which vsys these operations should take place in.  If
// vsys is an empty string, vsys defaults to "vsys1".
func (c *UserId) GetIpTags(ip, tag, vsys string) (map[string][]string, error) {
	if vsys == "" {
		vsys = "vsys1"
	}
	c.con.LogOp("(op) getting registered ip addresses - ip:%q tag:%q vsys:%q", ip, tag, vsys)
	req := c.versioning()

	ans := make(map[string][]string)
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

type filterer interface {
	FilterOn(string, string, int)
	ShouldStop(int) bool
}

type req_filter struct {
	Tag *tagFilter `xml:"tag"`
	Ip  string     `xml:"ip,omitempty"`
}

type tagFilter struct {
	Entry tagName `xml:"entry"`
}

type tagName struct {
	Name string `xml:"name,attr"`
}

type req_v1 struct {
	XMLName xml.Name   `xml:"show"`
	Filter  req_filter `xml:"object>registered-address"`
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
	XMLName xml.Name   `xml:"show"`
	Filter  req_filter `xml:"object>registered-ip"`
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
	XMLName xml.Name      `xml:"show"`
	Filter  req_filter_v2 `xml:"object>registered-ip"`
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
	Tag   *tagFilter `xml:"tag"`
	Ip    string     `xml:"ip,omitempty"`
	Limit int        `xml:"limit"`
	Start int        `xml:"start-point"`
}

type regResp struct {
	Entry []respEntry `xml:"result>entry"`
	Msg   *msgResp    `xml:"result>msg"`
}

type respEntry struct {
	Ip   string   `xml:"ip,attr"`
	Tags []string `xml:"tag>member"`
}

type msgResp struct {
	Outfile string `xml:"line>outfile"`
}
