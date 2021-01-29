package userid

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/PaloAltoNetworks/pango/util"
)

/*
GetGroups returns the list of groups defined.

The style param can be used to limit the groups returned to the specified kind.  If
style is an empty string, all groups are returned.

The vsys will default to "vsys1" if left as an empty string.
*/
func (c *UserId) GetGroups(style, vsys string) ([]string, error) {
	if vsys == "" {
		vsys = "vsys1"
	}
	req := groupListReq{}
	if style != "" {
		req.Style.Entry = &util.Entry{Value: style}
	}
	c.con.LogOp("(op) getting %q groups", style)

	resp := groupResp{}

	_, err := c.con.Op(req, vsys, nil, &resp)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(resp.Result.Text, "\n")
	ans := make([]string, 0, len(lines))
	for _, line := range lines {
		val := strings.TrimSpace(line)
		if val == "" {
			continue
		} else if strings.HasPrefix(val, "Total: ") {
			break
		}
		ans = append(ans, val)
	}

	return ans, nil
}

/*
Examples of things returned from PAN-OS when getting group members:

<response status="success"><result>User group 'static_group' does not exist or does not have members
]]></result></response>

<response status="success"><result><![CDATA[

source type: xmlapi

[1     ] sg21
[2     ] sg22

]]></result></response>
*/

/*
GetGroupsMembers returns the list of users in the given group.

The vsys will default to "vsys1" if left as an empty string.
*/
func (c *UserId) GetGroupMembers(group, vsys string) ([]string, error) {
	if vsys == "" {
		vsys = "vsys1"
	}
	req := groupMembersReq{Group: group}
	c.con.LogOp("(op) getting group members: %s", group)

	resp := groupResp{}

	b, err := c.con.Op(req, vsys, nil, &resp)
	if err != nil {
		if bytes.Contains(b, []byte("does not exist or does not have members")) {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(resp.Result.Text, "\n")
	ans := make([]string, 0, len(lines))
	for _, line := range lines {
		tokens := strings.Split(line, "]")
		if len(tokens) == 2 {
			ans = append(ans, strings.TrimSpace(tokens[1]))
		}
	}

	return ans, nil
}

type groupListReq struct {
	XMLName xml.Name       `xml:"show"`
	Style   groupListStyle `xml:"user>group>list"`
}

type groupListStyle struct {
	Entry *util.Entry `xml:"entry"`
}

type groupResp struct {
	XMLName xml.Name       `xml:"response"`
	Result  util.CdataText `xml:"result"`
}

type groupMembersReq struct {
	XMLName xml.Name `xml:"show"`
	Group   string   `xml:"user>group>name"`
}
