package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/dev/profile/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Resource tests.
func TestAccPanosHttpServerProfile_basic(t *testing.T) {
	var tmpl, rsName string

	if testAccIsFirewall {
		rsName = "panos_http_server_profile"
	} else {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
		rsName = "panos_panorama_http_server_profile"
	}

	var o http.Entry
	one := http.Entry{
		Name:   fmt.Sprintf("tf%s", acctest.RandString(6)),
		Config: &http.PayloadFormat{Name: "conf name"},
		System: &http.PayloadFormat{UriFormat: "/api/incident"},
		Threat: &http.PayloadFormat{Payload: "some payload"},
		Traffic: &http.PayloadFormat{Headers: []http.Header{{
			Name:  "Content-Type",
			Value: "text/plain",
		}}},
		HipMatch: &http.PayloadFormat{Parameters: []http.Parameter{{
			Name:  "type",
			Value: "security",
		}}},
		Url: &http.PayloadFormat{
			Name:      "base params",
			UriFormat: "/some/uri",
			Payload:   "this is a test",
		},
		Data: &http.PayloadFormat{
			Headers: []http.Header{
				{Name: "Secret-Id", Value: "swordfish"},
				{Name: "X-Scan-Engine", Value: "allow"},
			},
			Parameters: []http.Parameter{
				{Name: "lo", Value: "fi"},
				{Name: "hip", Value: "hop"},
			},
		},
		Servers: []http.Server{{
			Name:       "first server",
			Address:    "siem.example.com",
			HttpMethod: "POST",
			Username:   "foo",
			Password:   "bar",
		}},
	}
	two := http.Entry{
		Name:            one.Name,
		TagRegistration: true,
		Wildfire: &http.PayloadFormat{
			Name:      "all together now",
			UriFormat: "/app/endpoint/api",
			Payload:   "my payload",
			Headers: []http.Header{{
				Name:  "Content-Type",
				Value: "application/json",
			}},
			Parameters: []http.Parameter{{
				Name:  "serial",
				Value: "$serial",
			}},
		},
		Tunnel: &http.PayloadFormat{
			Name:    "tunnel format",
			Payload: "tool",
		},
		UserId: &http.PayloadFormat{
			Payload: "esthero",
			Headers: []http.Header{{
				Name:  "Beautiful",
				Value: "lie",
			}},
		},
		Gtp: &http.PayloadFormat{
			Name:      "wu tang",
			UriFormat: "/protect/ya/api",
			Parameters: []http.Parameter{
				{Name: "hostname", Value: "$host"},
				{Name: "type", Value: "$type"},
				{Name: "subt", Value: "$subtype"},
			},
		},
		Auth: &http.PayloadFormat{
			Payload: "<alert><sev>5</sev><msg>Fear Inoculum</msg></alert>",
			Headers: []http.Header{
				{Name: "Spiral-Out", Value: "Keep going"},
				{Name: "Content-Type", Value: "application/xml"},
			},
		},
		Servers: []http.Server{{
			Name:       one.Servers[0].Name,
			Address:    one.Servers[0].Address,
			HttpMethod: "GET",
			Username:   "blues",
			Password:   "brothers",
		}},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosHttpServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHttpServerProfileConfig(tmpl, one),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosHttpServerProfileExists(rsName, &o),
					testAccCheckPanosHttpServerProfileAttributes(&o, &one),
				),
			},
			{
				Config: testAccHttpServerProfileConfig(tmpl, two),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosHttpServerProfileExists(rsName, &o),
					testAccCheckPanosHttpServerProfileAttributes(&o, &two),
				),
			},
		},
	})
}

func testAccCheckPanosHttpServerProfileExists(n string, o *http.Entry) resource.TestCheckFunc {
	name := fmt.Sprintf("%s.test", n)
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s / %#v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v http.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseHttpServerProfileId(rs.Primary.ID)
			v, err = con.Device.HttpServerProfile.Get(vsys, name)
		case *pango.Panorama:
			tmpl, ts, vsys, name := parsePanoramaHttpServerProfileId(rs.Primary.ID)
			v, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosHttpServerPayload(desc string, v, conf *http.PayloadFormat) error {
	if v == nil && conf == nil {
		return nil
	} else if v == nil {
		return fmt.Errorf("%s is nil but should be %#v", desc, conf)
	} else if conf == nil {
		return fmt.Errorf("%s is %#v but should be nil", desc, v)
	}

	if v.Name != conf.Name {
		return fmt.Errorf("%s name is %s, not %s", desc, v.Name, conf.Name)
	}

	if v.UriFormat != conf.UriFormat {
		return fmt.Errorf("%s uri format is %s, not %s", desc, v.UriFormat, conf.UriFormat)
	}

	if v.Payload != conf.Payload {
		return fmt.Errorf("%s payload is %s, not %s", desc, v.Payload, conf.Payload)
	}

	if len(v.Headers) != len(conf.Headers) {
		return fmt.Errorf("%s headers length mismatch: %d not %d", desc, len(v.Headers), len(conf.Headers))
	}

	for i := range v.Headers {
		if v.Headers[i].Name != conf.Headers[i].Name {
			return fmt.Errorf("%s header[%d] name is %s, not %s", desc, i, v.Headers[i].Name, conf.Headers[i].Name)
		}

		if v.Headers[i].Value != conf.Headers[i].Value {
			return fmt.Errorf("%s header[%d] value is %s, not %s", desc, i, v.Headers[i].Value, conf.Headers[i].Value)
		}
	}

	if len(v.Parameters) != len(conf.Parameters) {
		return fmt.Errorf("%s parameters length mismatch: %d not %d", desc, len(v.Parameters), len(conf.Parameters))
	}

	for i := range v.Parameters {
		if v.Parameters[i].Name != conf.Parameters[i].Name {
			return fmt.Errorf("%s parameter[%d] name is %s, not %s", desc, i, v.Parameters[i].Name, conf.Parameters[i].Name)
		}

		if v.Parameters[i].Value != conf.Parameters[i].Value {
			return fmt.Errorf("%s parameter[%d] value is %s, not %s", desc, i, v.Parameters[i].Value, conf.Parameters[i].Value)
		}
	}

	return nil
}

func testAccCheckPanosHttpServerProfileAttributes(o, conf *http.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		if o.Name != conf.Name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, conf.Name)
		}

		if o.TagRegistration != conf.TagRegistration {
			return fmt.Errorf("Tag registration is %t not %t", o.TagRegistration, conf.TagRegistration)
		}

		if err = testAccCheckPanosHttpServerPayload("config", o.Config, conf.Config); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("system", o.System, conf.System); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("threat", o.Threat, conf.Threat); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("traffic", o.Traffic, conf.Traffic); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("hip match", o.HipMatch, conf.HipMatch); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("url", o.Url, conf.Url); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("data", o.Data, conf.Data); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("wildfire", o.Wildfire, conf.Wildfire); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("tunnel", o.Tunnel, conf.Tunnel); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("user id", o.UserId, conf.UserId); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("gtp", o.Gtp, conf.Gtp); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("auth", o.Auth, conf.Auth); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("sctp", o.Sctp, conf.Sctp); err != nil {
			return nil
		}

		if err = testAccCheckPanosHttpServerPayload("iptag", o.Iptag, conf.Iptag); err != nil {
			return nil
		}

		for i := range o.Servers {
			if o.Servers[i].Name != conf.Servers[i].Name {
				return fmt.Errorf("Server name is %s, not %s", o.Servers[i].Name, conf.Servers[i].Name)
			}

			if o.Servers[i].Address != conf.Servers[i].Address {
				return fmt.Errorf("Server address is %s, not %s", o.Servers[i].Address, conf.Servers[i].Address)
			}

			if o.Servers[i].Username != conf.Servers[i].Username {
				return fmt.Errorf("Server username is %s, not %s", o.Servers[i].Username, conf.Servers[i].Username)
			}
		}

		return nil
	}
}

func testAccPanosHttpServerProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_http_server_profile" && rs.Type != "panos_panorama_http_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseHttpServerProfileId(rs.Primary.ID)
				_, err = con.Device.HttpServerProfile.Get(vsys, name)
			case *pango.Panorama:
				tmpl, ts, vsys, name := parsePanoramaHttpServerProfileId(rs.Primary.ID)
				_, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccHttpServerProfileConfig(tmpl string, o http.Entry) string {
	var b strings.Builder
	var rsName, tmplSpec, tmplAttr string

	if tmpl == "" {
		rsName = "panos_http_server_profile"
	} else {
		rsName = "panos_panorama_http_server_profile"
		tmplSpec = fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}
`, tmpl)
		tmplAttr = `
    template = panos_panorama_template.x.name`
	}

	b.Grow(200 * len(o.Servers))
	for _, x := range o.Servers {
		b.WriteString(fmt.Sprintf(`
    http_server {
        name = %q
        address = %q
        http_method = %q
        username = %q
        password = %q
        certificate_profile = data.panos_system_info.x.version_major >= 9 ? "None" : ""
        tls_version = data.panos_system_info.x.version_major >= 9 ? "1.2" : ""
    }`, x.Name, x.Address, x.HttpMethod, x.Username, x.Password))
	}

	return fmt.Sprintf(`
%s

data "panos_system_info" "x" {}

resource %q "test" {
%s
    vsys = "shared"
    name = %q
    tag_registration = %t
%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s
}
`,
		tmplSpec, rsName, tmplAttr, o.Name, o.TagRegistration, b.String(),
		testAccHttpPayloadConfig("config_format", o.Config),
		testAccHttpPayloadConfig("system_format", o.System),
		testAccHttpPayloadConfig("threat_format", o.Threat),
		testAccHttpPayloadConfig("traffic_format", o.Traffic),
		testAccHttpPayloadConfig("hip_match_format", o.HipMatch),
		testAccHttpPayloadConfig("url_format", o.Url),
		testAccHttpPayloadConfig("data_format", o.Data),
		testAccHttpPayloadConfig("wildfire_format", o.Wildfire),
		testAccHttpPayloadConfig("tunnel_format", o.Tunnel),
		testAccHttpPayloadConfig("user_id_format", o.UserId),
		testAccHttpPayloadConfig("gtp_format", o.Gtp),
		testAccHttpPayloadConfig("auth_format", o.Auth),
		testAccHttpPayloadConfig("sctp_format", o.Sctp),
		testAccHttpPayloadConfig("iptag_format", o.Iptag),
	)
}

func testAccHttpPayloadConfig(desc string, v *http.PayloadFormat) string {
	if v == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
    %s {`, desc))

	if v.Name != "" {
		b.WriteString(fmt.Sprintf(`
        name = %q`, v.Name))
	}

	if v.UriFormat != "" {
		b.WriteString(fmt.Sprintf(`
        uri_format = %q`, v.UriFormat))
	}

	if v.Payload != "" {
		b.WriteString(fmt.Sprintf(`
        payload = %q`, v.Payload))
	}

	if len(v.Headers) > 0 {
		b.WriteString(`
        headers = {`)
		for _, x := range v.Headers {
			b.WriteString(fmt.Sprintf(`
            %q: %q,`, x.Name, x.Value))
		}
		b.WriteString(`
        }`)
	}

	if len(v.Parameters) > 0 {
		b.WriteString(`
        params = {`)
		for _, x := range v.Parameters {
			b.WriteString(fmt.Sprintf(`
            %q: %q,`, x.Name, x.Value))
		}
		b.WriteString(`
        }`)
	}

	b.WriteString(`
    }`)

	return b.String()
}
