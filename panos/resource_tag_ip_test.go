package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosTagIp_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	ip := "192.168.55.3"
	o := make(map[string][]string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTagIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagIpConfig(ip, "tag1", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTagIpExists("panos_tag_ip.test", o),
					testAccCheckPanosTagIpAttributes(o, ip, "tag1", "tag2"),
				),
			},
			{
				Config: testAccTagIpConfig(ip, "tag2", "tag3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTagIpExists("panos_tag_ip.test", o),
					testAccCheckPanosTagIpAttributes(o, ip, "tag2", "tag3"),
				),
			},
		},
	})
}

func testAccCheckPanosTagIpExists(n string, o map[string][]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		vsys, ip := parseTagIpId(rs.Primary.ID)
		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.UserId.GetIpTags(ip, "", vsys)
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

func testAccCheckPanosTagIpAttributes(o map[string][]string, ip, t1, t2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		info := o[ip]
		if len(info) != 2 {
			return fmt.Errorf("IP %q has %d/2 tags", ip, len(info))
		}
		for _, v := range []string{t1, t2} {
			if v != info[0] && v != info[1] {
				return fmt.Errorf("%q not present in tags: %#v", v, info)
			}
		}

		return nil
	}
}

func testAccPanosTagIpDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_tag_ip" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, ip := parseTagIpId(rs.Primary.ID)
			cur, err := fw.UserId.GetIpTags(ip, "", vsys)
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

func testAccTagIpConfig(ip, t1, t2 string) string {
	return fmt.Sprintf(`
resource "panos_tag_ip" "test" {
    ip = %q
    tags = [%q, %q]
}
`, ip, t1, t2)
}
