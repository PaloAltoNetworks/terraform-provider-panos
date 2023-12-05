package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/dug"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing tests.
func TestAccPanosDsDynamicUserGroupList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDynamicUserGroupConfig(name),
				Check:  checkDataSourceListing("panos_dynamic_user_groups"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsDynamicUserGroup_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDynamicUserGroupConfig(name),
				Check: checkDataSource("panos_dynamic_user_group", []string{
					"name", "description", "filter",
				}),
			},
		},
	})
}

func testAccDsDynamicUserGroupConfig(name string) string {
	return fmt.Sprintf(`
data "panos_dynamic_user_groups" "test" {}

data "panos_dynamic_user_group" "test" {
    name = panos_dynamic_user_group.x.name
}

resource "panos_dynamic_user_group" "x" {
    name = %q
    description = "dug ds acctest"
    filter = "'tomato'"
}
`, name)
}

// Resource tests.
func TestAccPanosDynamicUserGroup_basic(t *testing.T) {
	var o dug.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDynamicUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDynamicUserGroupConfig(name, "first", "'foo' or 'bar'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDynamicUserGroupExists("panos_dynamic_user_group.test", &o),
					testAccCheckPanosDynamicUserGroupAttributes(&o, name, "first", "'foo' or 'bar'"),
				),
			},
			{
				Config: testAccDynamicUserGroupConfig(name, "second", "'wu' and 'tang' and 'clan'"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDynamicUserGroupExists("panos_dynamic_user_group.test", &o),
					testAccCheckPanosDynamicUserGroupAttributes(&o, name, "second", "'wu' and 'tang' and 'clan'"),
				),
			},
		},
	})
}

func testAccCheckPanosDynamicUserGroupExists(n string, o *dug.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v dug.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseDynamicUserGroupId(rs.Primary.ID)
			v, err = con.Objects.DynamicUserGroup.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseDynamicUserGroupId(rs.Primary.ID)
			v, err = con.Objects.DynamicUserGroup.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDynamicUserGroupAttributes(o *dug.Entry, name, desc, filter string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, not %q", o.Description, desc)
		}

		if o.Filter != filter {
			return fmt.Errorf("Filter is %q, not %q", o.Filter, filter)
		}

		return nil
	}
}

func testAccPanosDynamicUserGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_dynamic_user_group" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseDynamicUserGroupId(rs.Primary.ID)
				_, err = con.Objects.DynamicUserGroup.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseDynamicUserGroupId(rs.Primary.ID)
				_, err = con.Objects.DynamicUserGroup.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccDynamicUserGroupConfig(name, desc, filter string) string {
	return fmt.Sprintf(`
resource "panos_dynamic_user_group" "test" {
    name = %q
    description = %q
    filter = %q
}
`, name, desc, filter)
}
