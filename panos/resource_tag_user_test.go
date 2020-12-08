package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosTagUser_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	ip := fmt.Sprintf("192.168.%d.%d", (acctest.RandInt()%250)+1, (acctest.RandInt()%250 + 1))
	user := fmt.Sprintf("tf%s", acctest.RandString(6))
	o := make(map[string][]string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTagUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagUserConfig(ip, user, "tag1", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTagUserExists("panos_tag_user.test", o),
					testAccCheckPanosTagUserAttributes(o, user, "tag1", "tag2"),
				),
			},
			{
				Config: testAccTagUserConfig(ip, user, "tag2", "tag3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTagUserExists("panos_tag_user.test", o),
					testAccCheckPanosTagUserAttributes(o, user, "tag2", "tag3"),
				),
			},
		},
	})
}

func testAccCheckPanosTagUserExists(n string, o map[string][]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		vsys, user := parseTagUserId(rs.Primary.ID)
		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.UserId.GetUserTags(user, vsys)
		if err != nil {
			return err
		}

		for key := range o {
			delete(o, key)
		}
		for key, value := range v {
			o[key] = value
		}

		return nil
	}
}

func testAccCheckPanosTagUserAttributes(o map[string][]string, user, t1, t2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		info := o[user]
		if len(info) != 2 {
			return fmt.Errorf("User %q has %d/2 tags", user, len(info))
		}
		for _, v := range []string{t1, t2} {
			if v != info[0] && v != info[1] {
				return fmt.Errorf("%q not present in tags: %#v", v, info)
			}
		}

		return nil
	}
}

func testAccPanosTagUserDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_tag_user" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, user := parseTagUserId(rs.Primary.ID)
			cur, err := fw.UserId.GetUserTags(user, vsys)
			if err != nil {
				return err
			}
			if len(cur) != 0 {
				return fmt.Errorf("Found tags: %#v", cur)
			}
		}
		return nil
	}

	return nil
}

func testAccTagUserConfig(ip, user, t1, t2 string) string {
	return fmt.Sprintf(`
resource "panos_userid_login" "x" {
    ip = %q
    user = %q
}

resource "panos_tag_user" "test" {
    user = %q
    tags = [%q, %q]
}
`, ip, user, user, t1, t2)
}
