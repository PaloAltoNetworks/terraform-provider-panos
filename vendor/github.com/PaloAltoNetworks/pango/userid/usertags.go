package userid

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

/*
GetUserTags returns dynamic user tags.

Note:  PAN-OS 9.1+

The user param will filter on just the specified user instead of all users
and all tags.

If vsys is an empty string, then this function defaults to "vsys1".
*/
func (c *UserId) GetUserTags(user, vsys string) (map[string][]string, error) {
	if vsys == "" {
		vsys = "vsys1"
	}
	c.con.LogOp("(op) getting user tags: user:%q vsys %q", user, vsys)

	req := &gutReq{}

	ans := make(map[string][]string)
	length := 0
	for {
		req.FilterOn(user, length)
		resp := gutResp{}

		_, err := c.con.Op(req, vsys, nil, &resp)
		if err != nil {
			return nil, err
		}

		for i := range resp.Entry {
			ans[resp.Entry[i].User] = util.MemToStr(resp.Entry[i].Tags)
		}

		if req.ShouldStop(len(resp.Entry)) {
			break
		}

		length += len(resp.Entry)
	}

	return ans, nil
}

type gutReq struct {
	XMLName xml.Name `xml:"show"`
	Data    gutData  `xml:"object>registered-user"`
}

type gutData struct {
	All  *gutDataAll `xml:"all"`
	User string      `xml:"user,omitempty"`
}

type gutDataAll struct {
	Limit int `xml:"limit"`
	Start int `xml:"start-point"`
}

func (o *gutReq) FilterOn(user string, length int) {
	if user != "" {
		if o.Data.User == "" {
			o.Data.User = user
		}
	} else {
		if o.Data.All == nil {
			o.Data.All = &gutDataAll{
				Limit: 500,
			}
		}
		o.Data.All.Start = length + 1
	}
}

func (o *gutReq) ShouldStop(lastCount int) bool {
	if o.Data.All == nil {
		return true
	}

	return lastCount < o.Data.All.Limit
}

type gutResp struct {
	Entry []gutEntry `xml:"result>entry"`
}

type gutEntry struct {
	User string           `xml:"user,attr"`
	Tags *util.MemberType `xml:"tag"`
}
