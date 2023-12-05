package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source test.
func TestAccPanosDsIpTag(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	ip := fmt.Sprintf("10.1.59.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsIpTagConfig(ip),
				Check: checkDataSource("panos_ip_tag", []string{
					"ip", "entries.0.ip",
				}),
			},
		},
	})
}

func testAccDsIpTagConfig(ip string) string {
	return fmt.Sprintf(`
data "panos_ip_tag" "test" {
    ip = panos_ip_tag.x.ip
}

resource "panos_ip_tag" "x" {
    ip = %q
    tags = [
        "terraform",
        "datasource",
        "acctest",
    ]
}
`, ip)
}

// Resource test.
func TestAccPanosIpTag(t *testing.T) {
	var o map[string][]string
	ip := fmt.Sprintf("10.1.59.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIpTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpTagConfig(ip, "tag1", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpTagExists("panos_ip_tag.test", &o),
					testAccCheckPanosIpTagAttributes(&o, ip, "tag1", "tag2"),
				),
			},
			{
				Config: testAccIpTagConfig(ip, "tag3", "tag2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpTagExists("panos_ip_tag.test", &o),
					testAccCheckPanosIpTagAttributes(&o, ip, "tag3", "tag2"),
				),
			},
		},
	})
}

func testAccCheckPanosIpTagExists(n string, o *map[string][]string) resource.TestCheckFunc {
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
			vsys, ip, _ := parseIpTagId(rs.Primary.ID)
			v, err = con.UserId.GetIpTags(ip, "", vsys)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIpTagAttributes(o *map[string][]string, name, tag1, tag2 string) resource.TestCheckFunc {
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

func testAccPanosIpTagDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ip_tag" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error
			var ip string
			var curTags map[string][]string

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, ip, _ := parseIpTagId(rs.Primary.ID)
				curTags, err = con.UserId.GetIpTags(ip, "", vsys)
				if err != nil {
					return err
				}
			}
			if len(curTags[ip]) != 0 {
				return fmt.Errorf("User %q still has tags: %#v", ip, curTags[ip])
			}
		}
		return nil
	}

	return nil
}

func testAccIpTagConfig(ip, tag1, tag2 string) string {
	return fmt.Sprintf(`
resource "panos_ip_tag" "test" {
    ip = %q
    tags = [%q, %q]
}
`, ip, tag1, tag2)
}
