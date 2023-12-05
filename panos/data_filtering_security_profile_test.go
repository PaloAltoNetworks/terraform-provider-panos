package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/data"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing tests.
func TestAccPanosDsDataFilteringSecurityProfileList(t *testing.T) {
	cdp := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDataFilteringSecurityProfileConfig(cdp, name),
				Check:  checkDataSourceListing("panos_data_filtering_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsDataFilteringSecurityProfile_basic(t *testing.T) {
	cdp := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDataFilteringSecurityProfileConfig(cdp, name),
				Check: checkDataSource("panos_data_filtering_security_profile", []string{
					"name", "description",
					"rule.0.name",
					"rule.0.data_pattern",
					"rule.0.applications.0",
					"rule.0.file_types.0",
					"rule.0.direction",
				}),
			},
		},
	})
}

func testAccDsDataFilteringSecurityProfileConfig(cdp, name string) string {
	return fmt.Sprintf(`
data "panos_data_filtering_security_profiles" "test" {}

data "panos_data_filtering_security_profile" "test" {
    name = panos_data_filtering_security_profile.x.name
}

resource "panos_custom_data_pattern_object" "x" {
    name = %q
    description = "for data filtering security profile ds acctest"
    type = "regex"
    regex {
        name = "my regex"
        file_types = ["any"]
        regex = "this is my regex"
    }
}

resource "panos_data_filtering_security_profile" "x" {
    name = %q
    description = "data filtering sec prof data source acctest"
    rule {
        data_pattern = panos_custom_data_pattern_object.x.name
        applications = ["any"]
        file_types = ["any"]
    }
}
`, cdp, name)
}

// Resource tests.
func TestAccPanosDataFilteringSecurityProfile_basic(t *testing.T) {
	var o data.Entry
	cdp := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDataFilteringSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFilteringSecurityProfileConfig(cdp, name, "desc one", "bugzilla", "ask.fm", "xls", "doc", data.DirectionUpload, true, 11, 12),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDataFilteringSecurityProfileExists("panos_data_filtering_security_profile.test", &o),
					testAccCheckPanosDataFilteringSecurityProfileAttributes(&o, name, "desc one", "bugzilla", "ask.fm", "xls", "doc", data.DirectionUpload, true, 11, 12),
				),
			},
			{
				Config: testAccDataFilteringSecurityProfileConfig(cdp, name, "desc two", "gist", "foursquare", "pdf", "docx", data.DirectionBoth, false, 13, 14),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDataFilteringSecurityProfileExists("panos_data_filtering_security_profile.test", &o),
					testAccCheckPanosDataFilteringSecurityProfileAttributes(&o, name, "desc two", "gist", "foursquare", "pdf", "docx", data.DirectionBoth, false, 13, 14),
				),
			},
		},
	})
}

func testAccCheckPanosDataFilteringSecurityProfileExists(n string, o *data.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v data.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseDataFilteringSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.DataFilteringProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseDataFilteringSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.DataFilteringProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDataFilteringSecurityProfileAttributes(o *data.Entry, name, desc, app1, app2, ft1, ft2, dir string, cap bool, at, bt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if o.DataCapture != cap {
			return fmt.Errorf("Data capture is %t, not %t", o.DataCapture, cap)
		}

		if len(o.Rules) != 1 {
			return fmt.Errorf("Rules is not len 1: %#v", o.Rules)
		}

		r := o.Rules[0]

		if r.Name == "" {
			return fmt.Errorf("Rule name is not set")
		}

		if len(r.Applications) != 2 || r.Applications[0] != app1 || r.Applications[1] != app2 {
			return fmt.Errorf("Rule app is %#v, not [%q, %q]", r.Applications, app1, app2)
		}

		if len(r.FileTypes) != 2 || r.FileTypes[0] != ft1 || r.FileTypes[1] != ft2 {
			return fmt.Errorf("Rule file type is %#v, not [%q, %q]", r.FileTypes, ft1, ft2)
		}

		if r.Direction != dir {
			return fmt.Errorf("Rule direction is %q, not %q", r.Direction, dir)
		}

		if r.AlertThreshold != at {
			return fmt.Errorf("Rule alert threshold is %d, not %d", r.AlertThreshold, at)
		}

		if r.BlockThreshold != bt {
			return fmt.Errorf("Rule block threshold is %d, not %d", r.BlockThreshold, bt)
		}

		return nil
	}
}

func testAccPanosDataFilteringSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_data_filtering_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseDataFilteringSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.DataFilteringProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseDataFilteringSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.DataFilteringProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccDataFilteringSecurityProfileConfig(cdp, name, desc, app1, app2, ft1, ft2, dir string, cap bool, at, bt int) string {
	return fmt.Sprintf(`
resource "panos_custom_data_pattern_object" "x" {
    name = %q
    description = "for data filtering security profile acctest"
    type = "regex"
    regex {
        name = "my regex"
        file_types = ["any"]
        regex = "this is my regex"
    }
}

resource "panos_data_filtering_security_profile" "test" {
    name = %q
    description = %q
    data_capture = %t
    rule {
        data_pattern = panos_custom_data_pattern_object.x.name
        applications = [%q, %q]
        file_types = [%q, %q]
        direction = %q
        alert_threshold = %d
        block_threshold = %d
    }
}
`, cdp, name, desc, cap, app1, app2, ft1, ft2, dir, at, bt)
}
