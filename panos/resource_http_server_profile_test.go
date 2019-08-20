package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/header"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/param"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/server"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosHttpServerProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var (
		o          http.Entry
		serverList []server.Entry
		headers    map[string][]header.Entry
		params     map[string][]param.Entry
	)

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosHttpServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHttpServerProfileConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosHttpServerProfileExists("panos_http_server_profile.test", &o, &serverList, &headers, &params),
					testAccCheckPanosHttpServerProfileAttributes(&o, &serverList, &headers, &params, name),
				),
			},
		},
	})
}

func testAccCheckPanosHttpServerProfileExists(n string, o *http.Entry, serverList *[]server.Entry, headers *map[string][]header.Entry, params *map[string][]param.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		var err error
		var list []string

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseHttpServerProfileId(rs.Primary.ID)
		v, err := fw.Device.HttpServerProfile.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		list, err = fw.Device.HttpServer.GetList(vsys, name)
		if err != nil {
			return err
		}
		entries := make([]server.Entry, 0, len(list))
		for i := range list {
			entry, err := fw.Device.HttpServer.Get(vsys, name, list[i])
			if err != nil {
				return err
			}
			entries = append(entries, entry)
		}

		*serverList = entries

		logtypes := []string{
			param.Config,
			param.System,
			param.Threat,
			param.Traffic,
			param.HipMatch,
			param.Url,
			param.Data,
			param.Wildfire,
			param.Tunnel,
			param.UserId,
			param.Gtp,
			param.Auth,
			param.Sctp,
			param.Iptag,
		}

		headerMap := make(map[string][]header.Entry)
		paramMap := make(map[string][]param.Entry)
		for _, logtype := range logtypes {
			list, err = fw.Device.HttpHeader.GetList(vsys, name, logtype)
			if err != nil {
				return err
			}
			if len(list) != 0 {
				headerList := make([]header.Entry, 0, len(list))
				for _, hdr := range list {
					entry, err := fw.Device.HttpHeader.Get(vsys, name, logtype, hdr)
					if err != nil {
						return err
					}
					headerList = append(headerList, entry)
				}
				headerMap[logtype] = headerList
			}

			list, err = fw.Device.HttpParam.GetList(vsys, name, logtype)
			if err != nil {
				return err
			}
			if len(list) != 0 {
				paramList := make([]param.Entry, 0, len(list))
				for _, prm := range list {
					entry, err := fw.Device.HttpParam.Get(vsys, name, logtype, prm)
					if err != nil {
						return err
					}
					paramList = append(paramList, entry)
				}
				paramMap[logtype] = paramList
			}
		}

		*headers = headerMap
		*params = paramMap

		return nil
	}
}

func testAccCheckPanosHttpServerProfileAttributes(o *http.Entry, serverList *[]server.Entry, headers *map[string][]header.Entry, params *map[string][]param.Entry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var hlist []header.Entry
		var plist []param.Entry

		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.ConfigName != "conf name" {
			return fmt.Errorf("Config name is %q, not 'conf name'", o.ConfigName)
		}

		if o.SystemUriFormat != "/api/incident" {
			return fmt.Errorf("System uri format is %q, not '/api/incident'", o.SystemUriFormat)
		}

		if o.ThreatPayload != "some payload" {
			return fmt.Errorf("Threat payload is %q, not 'some payload'", o.ThreatPayload)
		}

		hlist = (*headers)[param.Traffic]
		if len(hlist) != 1 || hlist[0].Name != "Content-Type" || hlist[0].Value != "text/plain" {
			return fmt.Errorf("Incorrect traffic headers: %#v", hlist)
		}

		plist = (*params)[param.HipMatch]
		if len(plist) != 1 || plist[0].Name != "type" || plist[0].Value != "security" {
			return fmt.Errorf("Incorrect hip match params: %#v", plist)
		}

		if o.UrlName != "base params" {
			return fmt.Errorf("URL name is %s", o.UrlName)
		}
		if o.UrlUriFormat != "/some/uri" {
			return fmt.Errorf("URL uri format is %q", o.UrlUriFormat)
		}
		if o.UrlPayload != "this is a test" {
			return fmt.Errorf("URL payload is %q", o.UrlPayload)
		}

		hlist = (*headers)[param.Data]
		if len(hlist) != 2 {
			return fmt.Errorf("Data headers are len %d", len(hlist))
		}
		for i := range hlist {
			if hlist[i].Name == "Secret-Id" {
				if hlist[i].Value != "swordfish" {
					return fmt.Errorf("Data header 'Secret-Id' has value %s", hlist[i].Value)
				}
			} else if hlist[i].Name == "X-Scan-Engine" {
				if hlist[i].Value != "allow" {
					return fmt.Errorf("Data header 'X-Scan-Engine' has value %s", hlist[i].Value)
				}
			} else {
				return fmt.Errorf("Bad data header: %#v", hlist[i])
			}
		}
		plist = (*params)[param.Data]
		if len(plist) != 2 {
			return fmt.Errorf("Data params are len %d", len(plist))
		}
		for i := range plist {
			if plist[i].Name == "lo" {
				if plist[i].Value != "fi" {
					return fmt.Errorf("Data param 'lo' is %s", plist[i].Value)
				}
			} else if plist[i].Name == "hip" {
				if plist[i].Value != "hop" {
					return fmt.Errorf("Data param 'hip' is %s", plist[i].Value)
				}
			} else {
				return fmt.Errorf("Bad data param: %#v", plist[i])
			}
		}

		if o.WildfireName != "all together now" {
			return fmt.Errorf("Wildfire name is %s", o.WildfireName)
		}
		if o.WildfireUriFormat != "/app/endpoint/api" {
			return fmt.Errorf("Wildfire uri format is %q", o.WildfireUriFormat)
		}
		if o.WildfirePayload != "my payload" {
			return fmt.Errorf("Wildfire payload is %q", o.WildfirePayload)
		}
		hlist = (*headers)[param.Wildfire]
		if len(hlist) != 1 {
			return fmt.Errorf("Wildfire headers are len %d", len(hlist))
		}
		if hlist[0].Name != "Content-Type" || hlist[0].Value != "application/json" {
			return fmt.Errorf("Wildfire header is %#v", hlist[0])
		}
		plist = (*params)[param.Wildfire]
		if len(plist) != 1 {
			return fmt.Errorf("Wildfire params are len %d", len(plist))
		}
		if plist[0].Name != "serial" || plist[0].Value != "$serial" {
			return fmt.Errorf("Wildfire param is %#v", plist[0])
		}

		if o.TunnelName != "tunnel format" {
			return fmt.Errorf("Tunnel name is %q", o.TunnelName)
		}
		if o.TunnelPayload != "beautiful lie" {
			return fmt.Errorf("Tunnel payload is %q", o.TunnelPayload)
		}

		if o.UserIdPayload != "esthero" {
			return fmt.Errorf("UserID payload is %q", o.UserIdPayload)
		}
		hlist = (*headers)[param.UserId]
		if len(hlist) != 1 {
			return fmt.Errorf("UserID headers are len %d", len(hlist))
		}
		if hlist[0].Name != "Beautiful" || hlist[0].Value != "Lie" {
			return fmt.Errorf("UserID header is %#v", hlist[0])
		}

		if o.GtpName != "wu tang" {
			return fmt.Errorf("Gtp name is %q", o.GtpName)
		}
		if o.GtpUriFormat != "/protect/ya/api" {
			return fmt.Errorf("Gtp uri format is %q", o.GtpUriFormat)
		}
		plist = (*params)[param.Gtp]
		for i := range plist {
			if plist[i].Name == "hostname" {
				if plist[i].Value != "$host" {
					return fmt.Errorf("Gtp hostname is %q", plist[i].Value)
				}
			} else if plist[i].Name == "type" {
				if plist[i].Value != "$type" {
					return fmt.Errorf("Gtp type is %q", plist[i].Value)
				}
			} else if plist[i].Name == "subt" {
				if plist[i].Value != "$subtype" {
					return fmt.Errorf("Gtp subt is %q", plist[i].Value)
				}
			} else {
				return fmt.Errorf("Gtp has bad param %#v", plist[i])
			}
		}
		if len(plist) != 3 {
			return fmt.Errorf("Gtp params is len %d", len(plist))
		}

		if o.AuthPayload != "<alert><sev>5</sev><msg>Fear Inoculum</msg></alert>" {
			return fmt.Errorf("Auth payload is %q", o.AuthPayload)
		}
		hlist = (*headers)[param.Auth]
		for i := range hlist {
			if hlist[i].Name == "Spiral-Out" {
				if hlist[i].Value != "Keep going" {
					return fmt.Errorf("Auth header spiral out is %s", hlist[i].Value)
				}
			} else if hlist[i].Name == "Content-type" {
				if hlist[i].Value != "application/xml" {
					return fmt.Errorf("Auth header content type is %s", hlist[i].Value)
				}
			} else {
				return fmt.Errorf("Auth has bad header %#v", hlist[i])
			}
		}
		if len(hlist) != 2 {
			return fmt.Errorf("Auth headers are len %d", len(hlist))
		}

		return nil
	}
}

func testAccPanosHttpServerProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_http_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseHttpServerProfileId(rs.Primary.ID)
			_, err := fw.Device.HttpServerProfile.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccHttpServerProfileConfig(name string) string {
	return fmt.Sprintf(`
resource "panos_http_server_profile" "test" {
    name = %q
    http_server {
        name = "internal http server"
        address = "siem.example.com"
        username = "foo"
        password = "bar"
    }
    config_format {
        name = "conf name"
    }
    system_format {
        uri_format = "/api/incident"
    }
    threat_format {
        payload = "some payload"
    }
    traffic_format {
        headers = {
            "Content-Type": "text/plain",
        }
    }
    hip_match_format {
        params = {
            "type": "security",
        }
    }
    url_format {
        name = "base params"
        uri_format = "/some/uri"
        payload = "this is a test"
    }
    data_format {
        headers = {
            "Secret-Id": "swordfish",
            "X-Scan-Engine": "allow",
        }
        params = {
            "lo": "fi",
            "hip": "hop",
        }
    }
    wildfire_format {
        name = "all together now"
        uri_format = "/app/endpoint/api"
        payload = "my payload"
        headers = {
            "Content-Type": "application/json",
        }
        params = {
            "serial": "$serial",
        }
    }
    tunnel_format {
        name = "tunnel format"
        payload = "beautiful lie"
    }
    user_id_format {
        payload = "esthero"
        headers = {
            "Beautiful": "Lie",
        }
    }
    gtp_format {
        name = "wu tang"
        uri_format = "/protect/ya/api"
        params = {
            "hostname": "$host",
            "type": "$type",
            "subt": "$subtype",
        }
    }
    auth_format {
        payload = "<alert><sev>5</sev><msg>Fear Inoculum</msg></alert>"
        headers = {
            "Spiral-Out": "Keep going",
            "Content-type": "application/xml",
        }
    }
}
`, name)
}
