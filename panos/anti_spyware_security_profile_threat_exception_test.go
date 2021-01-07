package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	tex "github.com/PaloAltoNetworks/pango/objs/profile/security/spyware/texception"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsAntiSpywareSecurityProfileThreatExceptionList(t *testing.T) {
	if len(testAccPredefinedThreats) == 0 {
		t.Skip("No predefined threats present")
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := testAccPredefinedThreats[acctest.RandInt()%len(testAccPredefinedThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileThreatExceptionConfig(prof, name),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profile_threat_exceptions"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsAntiSpywareSecurityProfileThreatException(t *testing.T) {
	if len(testAccPredefinedThreats) == 0 {
		t.Skip("No predefined threats present")
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := testAccPredefinedThreats[acctest.RandInt()%len(testAccPredefinedThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileThreatExceptionConfig(prof, name),
				Check: checkDataSource("panos_anti_spyware_security_profile_threat_exception", []string{
					"name", "packet_capture", "action",
				}),
			},
		},
	})
}

func testAccDsAntiSpywareSecurityProfileThreatExceptionConfig(prof, name string) string {
	return fmt.Sprintf(`
data "panos_anti_spyware_security_profile_threat_exceptions" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
}

data "panos_anti_spyware_security_profile_threat_exception" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    name = panos_anti_spyware_security_profile_threat_exception.x.name
}

resource "panos_anti_spyware_security_profile" "x" {
    name = %q
    description = "anti_spyware sec prof threat exception acctest"
    sinkhole_ipv4_address = "pan-sinkhole-default-ip"
    sinkhole_ipv6_address = "::1"
    botnet_list {
        name = "default-paloalto-dns"
        action = "sinkhole"
        packet_capture = "disable"
    }
    botnet_list {
        name = "default-paloalto-cloud"
        action = "allow"
        packet_capture = "disable"
    }
}

resource "panos_anti_spyware_security_profile_threat_exception" "x" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    name = %q
    action = "default"
}
`, prof, name)
}

// Resource tests.
func TestAccPanosFirewallAntiSpywareSecurityProfileThreatException_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if len(testAccPredefinedThreats) == 0 {
		t.Skip("No predefined threats present")
	}

	var o tex.Entry
	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := testAccPredefinedThreats[acctest.RandInt()%len(testAccPredefinedThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntiSpywareSecurityProfileThreatExceptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiSpywareSecurityProfileThreatExceptionConfig(prof, name, tex.Disable, tex.ActionAllow, "192.168.55.55", "192.168.44.44"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionExists("panos_anti_spyware_security_profile_threat_exception.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionAttributes(&o, name, tex.Disable, tex.ActionAllow, "192.168.55.55", "192.168.44.44"),
				),
			},
			{
				Config: testAccAntiSpywareSecurityProfileThreatExceptionConfig(prof, name, tex.ExtendedCapture, tex.ActionDrop, "192.168.55.55", "192.168.66.66"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionExists("panos_anti_spyware_security_profile_threat_exception.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionAttributes(&o, name, tex.ExtendedCapture, tex.ActionDrop, "192.168.55.55", "192.168.66.66"),
				),
			},
		},
	})
}

/*
func TestAccPanosPanoramaAntiSpywareSecurityProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o spyware.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntiSpywareSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiSpywareSecurityProfileConfig(name, "desc one", "ip1.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileExists("panos_anti_spyware_security_profile.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileAttributes(&o, name, "desc one", "ip1.example.com"),
				),
			},
			{
				Config: testAccAntiSpywareSecurityProfileConfig(name, "desc two", "ip2.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileExists("panos_anti_spyware_security_profile.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileAttributes(&o, name, "desc two", "ip2.example.com"),
				),
			},
		},
	})
}
*/

func testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionExists(n string, o *tex.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v tex.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(rs.Primary.ID)
			v, err = con.Objects.AntiSpywareThreatException.Get(vsys, prof, name)
		case *pango.Panorama:
			dg, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(rs.Primary.ID)
			v, err = con.Objects.AntiSpywareThreatException.Get(dg, prof, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAntiSpywareSecurityProfileThreatExceptionAttributes(o *tex.Entry, name, cap, action, ip1, ip2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.PacketCapture != cap {
			return fmt.Errorf("Packet capture is %q, expected %q", o.PacketCapture, cap)
		}

		if o.Action != action {
			return fmt.Errorf("Action is %q, expected %q", o.Action, action)
		}

		if len(o.ExemptIps) != 2 {
			return fmt.Errorf("Exempt IPs is %#v not len 2", o.ExemptIps)
		}

		if o.ExemptIps[0] != ip1 {
			return fmt.Errorf("Exempt IP1 is %q, not %q", o.ExemptIps[0], ip1)
		}

		if o.ExemptIps[1] != ip2 {
			return fmt.Errorf("Exempt IP2 is %q, not %q", o.ExemptIps[1], ip2)
		}

		return nil
	}
}

func testAccPanosAntiSpywareSecurityProfileThreatExceptionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_anti_spyware_security_profile_threat_exception" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(rs.Primary.ID)
				_, err = con.Objects.AntiSpywareThreatException.Get(vsys, prof, name)
			case *pango.Panorama:
				dg, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(rs.Primary.ID)
				_, err = con.Objects.AntiSpywareThreatException.Get(dg, prof, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAntiSpywareSecurityProfileThreatExceptionConfig(prof, name, cap, action, ip1, ip2 string) string {
	return fmt.Sprintf(`
resource "panos_anti_spyware_security_profile" "x" {
    name = %q
    description = "for resource threat exception acctest"
    sinkhole_ipv4_address = "pan-sinkhole-default-ip"
    sinkhole_ipv6_address = "::1"
    botnet_list {
        name = "default-paloalto-dns"
        action = "sinkhole"
        packet_capture = "single-packet"
    }
    botnet_list {
        name = "default-paloalto-cloud"
        action = "allow"
        packet_capture = "extended-capture"
    }
}

resource "panos_anti_spyware_security_profile_threat_exception" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    name = %q
    packet_capture = %q
    action = %q
    exempt_ips = [%q, %q]
}
`, prof, name, cap, action, ip1, ip2)
}
