package userid

import (
	"bytes"
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// UserId is the client.UserId namespace.
type UserId struct {
	con util.XapiClient
}

// Initialize is invoked on client.Initialize().
func (c *UserId) Initialize(i util.XapiClient) {
	c.con = i
}

/*
Run executes the given User-Id message.  This allows you to
perform User-Id operations, such as logging in users, tagging
IP addresses, and setting group memberhsip.

Please refer to the Message class for the details.

The vsys param is which vsys these operations should take place in.  If
vsys is an empty string, vsys defaults to "vsys1".
*/
func (c *UserId) Run(msg *Message, vsys string) error {
	req, desc := encode(msg)
	if req == nil {
		return nil
	}

	if vsys == "" {
		vsys = "vsys1"
	}

	c.con.LogUid("(userid) running in %s -%s", vsys, desc)

	b, err := c.con.Uid(req, vsys, nil, nil)
	if bytes.Contains(b, []byte("<uid-response>")) {
		resp := uidResponse{}
		if e2 := xml.Unmarshal(b, &resp); e2 == nil {
			e3 := resp.GetError()
			if e3 != nil {
				return e3
			}
		}
	}

	return err
}
