package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/profile/security/wildfire"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsWildfireAnalysisSecurityProfileList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsWildfireAnalysisSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_wildfire_analysis_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsWildfireAnalysisSecurityProfile_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsWildfireAnalysisSecurityProfileConfig(name),
				Check: checkDataSource("panos_wildfire_analysis_security_profile", []string{
					"name", "description",
					"rule.0.name", "rule.0.direction", "rule.0.analysis",
					"rule.0.applications.0", "rule.0.file_types.0",
				}),
			},
		},
	})
}

func testAccDsWildfireAnalysisSecurityProfileConfig(name string) string {
	return fmt.Sprintf(`
data "panos_wildfire_analysis_security_profiles" "test" {}

data "panos_wildfire_analysis_security_profile" "test" {
    name = panos_wildfire_analysis_security_profile.x.name
}

resource "panos_wildfire_analysis_security_profile" "x" {
    name = %q
    description = "url filtering sec prof data source acctest"
    rule {
        name = "foo"
        applications = ["pop3"]
        file_types = ["pdf"]
    }
}
`, name)
}

// Resource tests.
func TestAccPanosWildfireAnalysisSecurityProfile_basic(t *testing.T) {
	var o wildfire.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosWildfireAnalysisSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWildfireAnalysisSecurityProfileConfig(name, "desc one", "first", "pop3", "aim-mail", "script", "pdf", wildfire.DirectionDownload, wildfire.AnalysisPrivateCloud),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosWildfireAnalysisSecurityProfileExists("panos_wildfire_analysis_security_profile.test", &o),
					testAccCheckPanosWildfireAnalysisSecurityProfileAttributes(&o, name, "desc one", "first", "pop3", "aim-mail", "script", "pdf", wildfire.DirectionDownload, wildfire.AnalysisPrivateCloud),
				),
			},
			{
				Config: testAccWildfireAnalysisSecurityProfileConfig(name, "desc two", "second", "aim-mail", "bbc-streaming", "ms-office", "pdf", wildfire.DirectionBoth, wildfire.AnalysisPublicCloud),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosWildfireAnalysisSecurityProfileExists("panos_wildfire_analysis_security_profile.test", &o),
					testAccCheckPanosWildfireAnalysisSecurityProfileAttributes(&o, name, "desc two", "second", "aim-mail", "bbc-streaming", "ms-office", "pdf", wildfire.DirectionBoth, wildfire.AnalysisPublicCloud),
				),
			},
		},
	})
}

func testAccCheckPanosWildfireAnalysisSecurityProfileExists(n string, o *wildfire.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v wildfire.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseWildfireAnalysisSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.WildfireAnalysisProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseWildfireAnalysisSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.WildfireAnalysisProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosWildfireAnalysisSecurityProfileAttributes(o *wildfire.Entry, name, desc, rName, app1, app2, ft1, ft2, dir, ana string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if len(o.Rules) != 1 {
			return fmt.Errorf("Rules is not len 1: %#v", o.Rules)
		}

		r := o.Rules[0]

		if len(r.Applications) != 2 {
			return fmt.Errorf("Rule applicationss is not len2: %#v", r.Applications)
		}

		if r.Applications[0] != app1 {
			return fmt.Errorf("Rule app1 is %q, not %q", r.Applications[0], app1)
		}

		if r.Applications[1] != app2 {
			return fmt.Errorf("Rule app2 is %q, not %q", r.Applications[1], app2)
		}

		if len(r.FileTypes) != 2 {
			return fmt.Errorf("Rule file types is not len2: %#v", r.FileTypes)
		}

		if r.FileTypes[0] != ft1 {
			return fmt.Errorf("Rule ft1 is %q, not %q", r.FileTypes[0], ft1)
		}

		if r.FileTypes[1] != ft2 {
			return fmt.Errorf("Rule ft2 is %q, not %q", r.FileTypes[1], ft2)
		}

		if r.Direction != dir {
			return fmt.Errorf("Rule direction is %q, not %q", r.Direction, dir)
		}

		if r.Analysis != ana {
			return fmt.Errorf("Rule analysis is %q, not %q", r.Analysis, ana)
		}

		return nil
	}
}

func testAccPanosWildfireAnalysisSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_wildfire_analysis_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseWildfireAnalysisSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.WildfireAnalysisProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseWildfireAnalysisSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.WildfireAnalysisProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccWildfireAnalysisSecurityProfileConfig(name, desc, rName, app1, app2, ft1, ft2, dir, ana string) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

resource "panos_wildfire_analysis_security_profile" "test" {
    name = %q
    description = %q
    rule {
        name = %q
        applications = [%q, %q]
        file_types = [%q, %q]
        direction = %q
        analysis = %q
    }
}
`, name, desc, rName, app1, app2, ft1, ft2, dir, ana)
}
