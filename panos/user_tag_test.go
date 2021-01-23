package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test.
func TestAccPanosDsUserTag(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsUserTagConfig(name),
				Check: checkDataSource("panos_user_tag", []string{
					"user", "users.0.user",
				}),
			},
		},
	})
}

func testAccDsUserTagConfig(name string) string {
	return fmt.Sprintf(`
data "panos_user_tag" "test" {
    user = panos_user_tag.x.user
}

resource "panos_user_tag" "x" {
    user = panos_userid_login.x.user
    tags = [
        "terraform",
        "acctest",
    ]
}

resource "panos_userid_login" "x" {
    ip = "10.5.6.89"
    user = %q
}
`, name)
}

// Resource test.
func TestAccPanosUserTag(t *testing.T) {
	var o map[string][]string
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosUserTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserTagConfig(name, "tag1", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosUserTagExists("panos_user_tag.test", &o),
					testAccCheckPanosUserTagAttributes(&o, name, "tag1", "tag2"),
				),
			},
			{
				Config: testAccUserTagConfig(name, "tag3", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosUserTagExists("panos_user_tag.test", &o),
					testAccCheckPanosUserTagAttributes(&o, name, "tag3", "tag2"),
				),
			},
		},
	})
}

func testAccCheckPanosUserTagExists(n string, o *map[string][]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v map[string][]string

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name, _ := parseUserTagId(rs.Primary.ID)
			v, err = con.UserId.GetUserTags(name, vsys)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosUserTagAttributes(o *map[string][]string, name, tag1, tag2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var found1 bool
		var found2 bool

		tags := (*o)[name]

		for _, tag := range tags {
			if tag == tag1 {
				found1 = true
			} else if tag == tag2 {
				found2 = true
			}
			if found1 && found2 {
				break
			}
		}

		if !found1 {
			return fmt.Errorf("%q not in tags: %#v", tag1, tags)
		}

		if !found2 {
			return fmt.Errorf("%q not in tags: %#v", tag2, tags)
		}

		return nil
	}
}

func testAccPanosUserTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_user_tag" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error
			var name string
			var curTags map[string][]string

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name, _ := parseUserTagId(rs.Primary.ID)
				curTags, err = con.UserId.GetUserTags(name, vsys)
				if err != nil {
					return err
				}
			}
			if len(curTags[name]) != 0 {
				return fmt.Errorf("User %q still has tags: %#v", name, curTags[name])
			}
		}
		return nil
	}

	return nil
}

func testAccUserTagConfig(name, tag1, tag2 string) string {
	return fmt.Sprintf(`
resource "panos_userid_login" "x" {
    ip = "10.20.59.77"
    user = %q
}

resource "panos_user_tag" "test" {
    user = panos_userid_login.x.user
    tags = [%q, %q]
}
`, name, tag1, tag2)
}
