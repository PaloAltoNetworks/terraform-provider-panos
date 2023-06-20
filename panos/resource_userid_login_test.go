package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/userid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosUseridLogin_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	ip := fmt.Sprintf("192.168.%d.%d", (acctest.RandInt()%250)+1, (acctest.RandInt()%250 + 1))
	user := fmt.Sprintf("tf%s", acctest.RandString(6))
	var o userid.LoginInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosUseridLoginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUseridLoginConfig(ip, user),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosUseridLoginExists("panos_userid_login.test", &o),
					testAccCheckPanosUseridLoginAttributes(&o, ip, user),
				),
			},
		},
	})
}

func testAccCheckPanosUseridLoginExists(n string, o *userid.LoginInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		vsys, ip, _ := parseUseridLoginId(rs.Primary.ID)
		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.UserId.GetLogins(ip, "", vsys)
		if err != nil {
			return err
		}

		if len(v) != 1 {
			return fmt.Errorf("Got %d results not 1 back from checkfunc", len(v))
		}

		*o = v[0]

		return nil
	}
}

func testAccCheckPanosUseridLoginAttributes(o *userid.LoginInfo, ip, user string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Ip != ip {
			return fmt.Errorf("Ip is %q expected %q", o.Ip, ip)
		}

		if o.User != user {
			return fmt.Errorf("User is %q expected %q", o.User, user)
		}

		return nil
	}
}

func testAccPanosUseridLoginDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_userid_login" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, ip, _ := parseUseridLoginId(rs.Primary.ID)
			cur, err := fw.UserId.GetLogins(ip, "", vsys)
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

func testAccUseridLoginConfig(ip, user string) string {
	return fmt.Sprintf(`
resource "panos_userid_login" "test" {
    ip = %q
    user = %q
}
`, ip, user)
}
