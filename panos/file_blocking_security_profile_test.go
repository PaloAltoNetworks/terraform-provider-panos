package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/file"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsFileBlockingSecurityProfileList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsFileBlockingSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_file_blocking_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsFileBlockingSecurityProfile_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsFileBlockingSecurityProfileConfig(name),
				Check: checkDataSource("panos_file_blocking_security_profile", []string{
					"name", "description",
					"rule.0.name", "rule.0.direction", "rule.0.action",
					"rule.0.applications.0", "rule.0.file_types.0",
				}),
			},
		},
	})
}

func testAccDsFileBlockingSecurityProfileConfig(name string) string {
	return fmt.Sprintf(`
data "panos_file_blocking_security_profiles" "test" {}

data "panos_file_blocking_security_profile" "test" {
    name = panos_file_blocking_security_profile.x.name
}

resource "panos_file_blocking_security_profile" "x" {
    name = %q
    description = "file blocking sec prof acctest"
    rule {
        name = "foo"
        applications = ["bbc-streaming"]
        file_types = ["ogg"]
    }
}
`, name)
}

// Resource tests.
func TestAccPanosFileBlockingSecurityProfile_basic(t *testing.T) {
	var o file.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosFileBlockingSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFileBlockingSecurityProfileConfig(name, "desc one", "first", "bbc-streaming", "facebook-mail", "jpeg", "gif", file.DirectionDownload, file.ActionContinue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosFileBlockingSecurityProfileExists("panos_file_blocking_security_profile.test", &o),
					testAccCheckPanosFileBlockingSecurityProfileAttributes(&o, name, "desc one", "first", "bbc-streaming", "facebook-mail", "jpeg", "gif", file.DirectionDownload, file.ActionContinue),
				),
			},
			{
				Config: testAccFileBlockingSecurityProfileConfig(name, "desc two", "second", "pop3", "dell-update", "ogg", "mp3", file.DirectionUpload, file.ActionAlert),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosFileBlockingSecurityProfileExists("panos_file_blocking_security_profile.test", &o),
					testAccCheckPanosFileBlockingSecurityProfileAttributes(&o, name, "desc two", "second", "pop3", "dell-update", "ogg", "mp3", file.DirectionUpload, file.ActionAlert),
				),
			},
		},
	})
}

func testAccCheckPanosFileBlockingSecurityProfileExists(n string, o *file.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v file.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseFileBlockingSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.FileBlockingProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseFileBlockingSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.FileBlockingProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosFileBlockingSecurityProfileAttributes(o *file.Entry, name, desc, rName, app1, app2, ft1, ft2, dir, action string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if len(o.Rules) != 1 {
			return fmt.Errorf("rules is not length 1: %#v", o.Rules)
		}

		r := o.Rules[0]
		if r.Name != rName {
			return fmt.Errorf("rule name is %q, expected %q", r.Name, rName)
		}

		if len(r.Applications) != 2 {
			return fmt.Errorf("rule apps is not len 2: %#v", r.Applications)
		}

		if r.Applications[0] != app1 {
			return fmt.Errorf("rule app1 is %q, not %q", r.Applications[0], app1)
		}

		if r.Applications[1] != app2 {
			return fmt.Errorf("rule app2 is %q, not %q", r.Applications[1], app2)
		}

		if len(r.FileTypes) != 2 {
			return fmt.Errorf("rule file types is not len 2: %#v", r.FileTypes)
		}

		if r.FileTypes[0] != ft1 {
			return fmt.Errorf("rule ft1 is %q, not %q", r.FileTypes[0], ft1)
		}

		if r.FileTypes[1] != ft2 {
			return fmt.Errorf("rule ft2 is %q, not %q", r.FileTypes[1], ft2)
		}

		if r.Direction != dir {
			return fmt.Errorf("rule direction is %q, not %q", r.Direction, dir)
		}

		if r.Action != action {
			return fmt.Errorf("rule action is %q, not %q", r.Action, action)
		}

		return nil
	}
}

func testAccPanosFileBlockingSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_file_blocking_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseFileBlockingSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.FileBlockingProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseFileBlockingSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.FileBlockingProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccFileBlockingSecurityProfileConfig(name, desc, rName, app1, app2, ft1, ft2, dir, action string) string {
	return fmt.Sprintf(`
resource "panos_file_blocking_security_profile" "test" {
    name = %q
    description = %q
    rule {
        name = %q
        applications = [%q, %q]
        file_types = [%q, %q]
        direction = %q
        action = %q
    }
}
`, name, desc, rName, app1, app2, ft1, ft2, dir, action)
}
