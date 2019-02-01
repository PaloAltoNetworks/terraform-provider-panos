package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/profile/redist/ipv4"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosRedistributionProfileIpv4(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o ipv4.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosRedistributionProfileIpv4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRedistributionProfileIpv4Config(vr, name, 1, ipv4.ActionRedist, ipv4.TypeBgp, ipv4.TypeConnect),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosRedistributionProfileIpv4Exists("panos_redistribution_profile_ipv4.test", &o),
					testAccCheckPanosRedistributionProfileIpv4Attributes(&o, name, 1, ipv4.ActionRedist, ipv4.TypeBgp, ipv4.TypeConnect),
				),
			},
			{
				Config: testAccRedistributionProfileIpv4Config(vr, name, 2, ipv4.ActionRedist, ipv4.TypeOspf, ipv4.TypeRip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosRedistributionProfileIpv4Exists("panos_redistribution_profile_ipv4.test", &o),
					testAccCheckPanosRedistributionProfileIpv4Attributes(&o, name, 2, ipv4.ActionRedist, ipv4.TypeOspf, ipv4.TypeRip),
				),
			},
			{
				Config: testAccRedistributionProfileIpv4Config(vr, name, 3, ipv4.ActionNoRedist, ipv4.TypeBgp, ipv4.TypeStatic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosRedistributionProfileIpv4Exists("panos_redistribution_profile_ipv4.test", &o),
					testAccCheckPanosRedistributionProfileIpv4Attributes(&o, name, 3, ipv4.ActionNoRedist, ipv4.TypeBgp, ipv4.TypeStatic),
				),
			},
		},
	})
}

func testAccCheckPanosRedistributionProfileIpv4Exists(n string, o *ipv4.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseRedistributionProfileIpv4Id(rs.Primary.ID)
		v, err := fw.Network.RedistributionProfile.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosRedistributionProfileIpv4Attributes(o *ipv4.Entry, name string, pri int, act, t1, t2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Priority != pri {
			return fmt.Errorf("Priority is %d, expected %d", o.Priority, pri)
		}

		if o.Action != act {
			return fmt.Errorf("Action is %s, expected %s", o.Action, act)
		}

		if len(o.Types) != 2 {
			return fmt.Errorf("Types is %#v, expected 2 entries", o.Types)
		}

		if o.Types[0] != t1 {
			return fmt.Errorf("Types[0] is %s, expected %s", o.Types[0], t1)
		}

		if o.Types[1] != t2 {
			return fmt.Errorf("Types[1] is %s, expected %s", o.Types[1], t2)
		}

		return nil
	}
}

func testAccPanosRedistributionProfileIpv4Destroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_redistribution_profile_ipv4" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseRedistributionProfileIpv4Id(rs.Primary.ID)
			_, err := fw.Network.RedistributionProfile.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccRedistributionProfileIpv4Config(vr, name string, pri int, act, t1, t2 string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "eth6" {
    name = "ethernet1/6"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
}

resource "panos_virtual_router" "vr" {
    name = %q
    interfaces = ["${panos_ethernet_interface.eth6.name}"]
}

resource "panos_redistribution_profile_ipv4" "test" {
    name = %q
    virtual_router = "${panos_virtual_router.vr.name}"
    priority = %d
    action = %q
    types = [%q, %q]
    interfaces = ["ethernet1/6"]
}
`, vr, name, pri, act, t1, t2)
}
