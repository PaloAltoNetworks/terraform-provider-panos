package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	av "github.com/PaloAltoNetworks/pango/objs/profile/security/virus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccDsPanosFirewallAntivirusSecurityProfileList(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntivirusSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_antivirus_security_profiles"),
			},
		},
	})
}

func TestAccDsPanosPanoramaAntivirusSecurityProfileList(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntivirusSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_antivirus_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccDsPanosFirewallAntivirusSecurityProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntivirusSecurityProfileConfig(name),
				Check: checkDataSource("panos_antivirus_security_profile", []string{
					"name", "description", "packet_capture",
				}),
			},
		},
	})
}

func TestAccDsPanosPanoramaAntivirusSecurityProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntivirusSecurityProfileConfig(name),
				Check: checkDataSource("panos_antivirus_security_profile", []string{
					"name", "description", "packet_capture",
				}),
			},
		},
	})
}

func testAccDsAntivirusSecurityProfileConfig(name string) string {
	return fmt.Sprintf(`
data "panos_antivirus_security_profiles" "test" {}

data "panos_antivirus_security_profile" "test" {
    name = panos_antivirus_security_profile.x.name
}

resource "panos_antivirus_security_profile" "x" {
    name = %q
    description = "antivirus sec prof acctest"
    decoder { name = "smtp" }
    decoder { name = "smb" }
    decoder { name = "pop3" }
    decoder { name = "imap" }
    decoder { name = "http2" }
    decoder { name = "http" }
    decoder { name = "ftp" }
    application_exception {
        application = "hotmail"
        action = "alert"
    }
}
`, name)
}

// Resource tests.
func TestAccPanosFirewallAntivirusSecurityProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o av.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntivirusSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntivirusSecurityProfileConfig(name, "desc one", av.Allow, av.Alert, "hotmail", av.Alert, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntivirusSecurityProfileExists("panos_antivirus_security_profile.test", &o),
					testAccCheckPanosAntivirusSecurityProfileAttributes(&o, name, "desc one", av.Allow, av.Alert, "hotmail", av.Alert, true),
				),
			},
			{
				Config: testAccAntivirusSecurityProfileConfig(name, "desc two", av.ResetClient, av.ResetServer, "aim-mail", av.ResetBoth, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntivirusSecurityProfileExists("panos_antivirus_security_profile.test", &o),
					testAccCheckPanosAntivirusSecurityProfileAttributes(&o, name, "desc two", av.ResetClient, av.ResetServer, "aim-mail", av.ResetBoth, false),
				),
			},
		},
	})
}

func TestAccPanosPanoramaAntivirusSecurityProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o av.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntivirusSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntivirusSecurityProfileConfig(name, "desc one", av.Allow, av.Alert, "hotmail", av.Alert, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntivirusSecurityProfileExists("panos_antivirus_security_profile.test", &o),
					testAccCheckPanosAntivirusSecurityProfileAttributes(&o, name, "desc one", av.Allow, av.Alert, "hotmail", av.Alert, true),
				),
			},
			{
				Config: testAccAntivirusSecurityProfileConfig(name, "desc two", av.ResetClient, av.ResetServer, "aim-mail", av.ResetBoth, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntivirusSecurityProfileExists("panos_antivirus_security_profile.test", &o),
					testAccCheckPanosAntivirusSecurityProfileAttributes(&o, name, "desc two", av.ResetClient, av.ResetServer, "aim-mail", av.ResetBoth, false),
				),
			},
		},
	})
}

func testAccCheckPanosAntivirusSecurityProfileExists(n string, o *av.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v av.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseAntivirusSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.AntivirusProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseAntivirusSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.AntivirusProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAntivirusSecurityProfileAttributes(o *av.Entry, name, desc, smtpAction, imapAction, app, appAction string, cap bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if o.PacketCapture != cap {
			return fmt.Errorf("Packet capture is %t, expected %t", o.PacketCapture, cap)
		}

		if len(o.ApplicationExceptions) != 1 {
			return fmt.Errorf("Expected 1 application exception: %#v", o.ApplicationExceptions)
		}

		if o.ApplicationExceptions[0].Application != app {
			return fmt.Errorf("App exception application is %q, expected %q", o.ApplicationExceptions[0].Application, app)
		}

		if o.ApplicationExceptions[0].Action != appAction {
			return fmt.Errorf("App exception action is %q, expected %q", o.ApplicationExceptions[0].Action, appAction)
		}

		var imap string
		var smtp string
		for _, dec := range o.Decoders {
			switch dec.Name {
			case "imap":
				imap = dec.Action
			case "smtp":
				smtp = dec.Action
			}
		}

		switch imap {
		case "":
			return fmt.Errorf("imap action was not present")
		case imapAction:
		default:
			return fmt.Errorf("imap action is %q, expected %q", imap, imapAction)
		}

		switch smtp {
		case "":
			return fmt.Errorf("smtp action was not present")
		case smtpAction:
		default:
			return fmt.Errorf("smtp action is %q, expected %q", smtp, smtpAction)
		}

		return nil
	}
}

func testAccPanosAntivirusSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_antivirus_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseAntivirusSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.AntivirusProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseAntivirusSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.AntivirusProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAntivirusSecurityProfileConfig(name, desc, smtpAction, imapAction, app, appAction string, cap bool) string {
	return fmt.Sprintf(`
resource "panos_antivirus_security_profile" "test" {
    name = %q
    description = %q
    packet_capture = %t
    decoder {
        name = "smtp"
        action = %q
    }
    decoder {
        name = "imap"
        action = %q
    }
    decoder { name = "smb" }
    decoder { name = "pop3" }
    decoder { name = "http2" }
    decoder { name = "http" }
    decoder { name = "ftp" }
    application_exception {
        application = %q
        action = %q
    }
}
`, name, desc, cap, smtpAction, imapAction, app, appAction)
}
