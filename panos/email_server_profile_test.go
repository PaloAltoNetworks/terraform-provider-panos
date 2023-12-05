package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing test.
func TestAccPanosDsEmailServerProfileList(t *testing.T) {
	var tmpl string

	if testAccIsPanorama {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
	}

	o := email.Entry{
		Name:     fmt.Sprintf("tf%s", acctest.RandString(6)),
		HipMatch: "$serial $severity",
		Config:   "$from $rule",
		Servers: []email.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				DisplayName:  "foobar",
				From:         "from@example.com",
				To:           "to@example.com",
				AlsoTo:       "archive@example.com",
				EmailGateway: "smtp.example.com",
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailServerProfileConfig(tmpl, o),
				Check:  checkDataSourceListing("panos_email_server_profiles"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsEmailServerProfile(t *testing.T) {
	var tmpl string

	if testAccIsPanorama {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
	}

	o := email.Entry{
		Name:     fmt.Sprintf("tf%s", acctest.RandString(6)),
		HipMatch: "$serial $severity",
		Config:   "$from $rule",
		Servers: []email.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				DisplayName:  "foobar",
				From:         "from@example.com",
				To:           "to@example.com",
				AlsoTo:       "archive@example.com",
				EmailGateway: "smtp.example.com",
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailServerProfileConfig(tmpl, o),
				Check: checkDataSource(
					"panos_email_server_profile",
					[]string{"name", "hip_match_format", "config_format"},
				),
			},
		},
	})
}

// Resource tests.
func TestAccPanosEmailServerProfile_basic(t *testing.T) {
	var tmpl, rsName string

	if testAccIsFirewall {
		rsName = "panos_email_server_profile"
	} else {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
		rsName = "panos_panorama_email_server_profile"
	}

	var o email.Entry
	one := email.Entry{
		Name:     fmt.Sprintf("tf%s", acctest.RandString(6)),
		HipMatch: "$serial $severity",
		Config:   "$from $rule",
		Servers: []email.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				DisplayName:  "foobar",
				From:         "from@example.com",
				To:           "to@example.com",
				AlsoTo:       "archive@example.com",
				EmailGateway: "smtp.example.com",
			},
		},
	}
	two := email.Entry{
		Name:     fmt.Sprintf("tf%s", acctest.RandString(6)),
		HipMatch: "$from $severity",
		Config:   "$serial $rule",
		Servers: []email.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				DisplayName:  "foobar",
				From:         "from@example.com",
				To:           "to@example.com",
				AlsoTo:       "archive@example.com",
				EmailGateway: "smtp.example.com",
			},
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				DisplayName:  "barbaz",
				From:         "john@example.com",
				To:           "jacob@example.com",
				AlsoTo:       "jingleheimer@example.com",
				EmailGateway: "smtp2.example.com",
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosEmailServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailServerProfileConfig(tmpl, one),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEmailServerProfileExists(rsName, &o),
					testAccCheckPanosEmailServerProfileAttributes(&o, &one),
				),
			},
			{
				Config: testAccEmailServerProfileConfig(tmpl, two),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEmailServerProfileExists(rsName, &o),
					testAccCheckPanosEmailServerProfileAttributes(&o, &two),
				),
			},
		},
	})
}

func testAccCheckPanosEmailServerProfileExists(n string, o *email.Entry) resource.TestCheckFunc {
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
		var v email.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseEmailServerProfileId(rs.Primary.ID)
			v, err = con.Device.EmailServerProfile.Get(vsys, name)
		case *pango.Panorama:
			tmpl, ts, vsys, name := parsePanoramaEmailServerProfileId(rs.Primary.ID)
			v, err = con.Device.EmailServerProfile.Get(tmpl, ts, vsys, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosEmailServerProfileAttributes(o, conf *email.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != conf.Name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, conf.Name)
		}

		if o.HipMatch != conf.HipMatch {
			return fmt.Errorf("HIP match format is %q, not %q", o.HipMatch, conf.HipMatch)
		}

		if o.Config != conf.Config {
			return fmt.Errorf("Config format is %q, not %q", o.Config, conf.Config)
		}

		if len(o.Servers) != len(conf.Servers) {
			return fmt.Errorf("Server list len mismatch: %d vs %d", len(o.Servers), len(conf.Servers))
		}

		for i := range o.Servers {
			if o.Servers[i].Name != conf.Servers[i].Name {
				return fmt.Errorf("Server name is %s, not %s", o.Servers[i].Name, conf.Servers[i].Name)
			}

			if o.Servers[i].DisplayName != conf.Servers[i].DisplayName {
				return fmt.Errorf("Server DisplayName is %q, not %q", o.Servers[i].DisplayName, conf.Servers[i].DisplayName)
			}

			if o.Servers[i].From != conf.Servers[i].From {
				return fmt.Errorf("Server from email is %s, not %s", o.Servers[i].From, conf.Servers[i].From)
			}

			if o.Servers[i].To != conf.Servers[i].To {
				return fmt.Errorf("Server to email is %s, not %s", o.Servers[i].To, conf.Servers[i].To)
			}

			if o.Servers[i].AlsoTo != conf.Servers[i].AlsoTo {
				return fmt.Errorf("Server also to email is %s, not %s", o.Servers[i].AlsoTo, conf.Servers[i].AlsoTo)
			}

			if o.Servers[i].EmailGateway != conf.Servers[i].EmailGateway {
				return fmt.Errorf("Server email gateway is %s, not %s", o.Servers[i].EmailGateway, conf.Servers[i].EmailGateway)
			}
		}

		return nil
	}
}

func testAccPanosEmailServerProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_email_server_profile" && rs.Type != "panos_panorama_email_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseEmailServerProfileId(rs.Primary.ID)
				_, err = con.Device.EmailServerProfile.Get(vsys, name)
			case *pango.Panorama:
				tmpl, ts, vsys, name := parsePanoramaEmailServerProfileId(rs.Primary.ID)
				_, err = con.Device.EmailServerProfile.Get(tmpl, ts, vsys, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccEmailServerProfileConfig(tmpl string, o email.Entry) string {
	var b strings.Builder
	var rsName, tmplSpec, tmplAttr string

	if tmpl == "" {
		rsName = "panos_email_server_profile"
	} else {
		rsName = "panos_panorama_email_server_profile"
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
    email_server {
        name = %q
        display_name = %q
        from_email = %q
        to_email = %q
        also_to_email = %q
        email_gateway = %q
    }
`, x.Name, x.DisplayName, x.From, x.To, x.AlsoTo, x.EmailGateway))
	}

	return fmt.Sprintf(`
%s

resource %q "test" {
%s
    vsys = "shared"
    name = %q
    hip_match_format = %q
    config_format = %q
%s
}

data "panos_email_server_profile" "test" {
%s
    vsys = "shared"
    name = %s.test.name
}

data "panos_email_server_profiles" "test" {
%s
    vsys = "shared"
}
`, tmplSpec, rsName, tmplAttr, o.Name, o.HipMatch, o.Config, b.String(), tmplAttr, rsName, tmplAttr)
}
