package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsAntiSpywareSecurityProfileList(t *testing.T) {
	if len(testAccPredefinedPhoneHomeThreats) == 0 {
		t.Skip("No predefined phone home threats present")
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tName := testAccPredefinedPhoneHomeThreats[acctest.RandInt()%len(testAccPredefinedPhoneHomeThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name, tName),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsAntiSpywareSecurityProfile_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tName := testAccPredefinedPhoneHomeThreats[acctest.RandInt()%len(testAccPredefinedPhoneHomeThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name, tName),
				Check: checkDataSource("panos_anti_spyware_security_profile", []string{
					"name", "description", "sinkhole_ipv4_address", "sinkhole_ipv6_address",
				}),
			},
		},
	})
}

func testAccDsAntiSpywareSecurityProfileConfig(name, tName string) string {
	return fmt.Sprintf(`
data "panos_anti_spyware_security_profiles" "test" {}

data "panos_anti_spyware_security_profile" "test" {
    name = panos_anti_spyware_security_profile.x.name
}

resource "panos_anti_spyware_security_profile" "x" {
    name = %q
    description = "anti_spyware sec prof acctest"
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
    rule {
        name = "foo"
        threat_name = "any"
        category = "adware"
        action = "alert"
        packet_capture = "disable"
        severities = ["any"]
    }
    exception {
        name = %q
        action = "default"
    }
}
`, name, tName)
}

// Resource tests.
func TestAccPanosAntiSpywareSecurityProfile_basic(t *testing.T) {
	if len(testAccPredefinedPhoneHomeThreats) == 0 {
		t.Skip("No predefined phone home threats present")
	}

	var o spyware.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	rName := fmt.Sprintf("tf%s", acctest.RandString(6))
	eName := testAccPredefinedPhoneHomeThreats[acctest.RandInt()%len(testAccPredefinedPhoneHomeThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntiSpywareSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiSpywareSecurityProfileConfig(name, "desc one", "ip1.example.com", rName, "any", "any", spyware.ActionResetClient, spyware.SinglePacket, "critical", "low", eName, spyware.Disable, spyware.ActionAllow, "192.168.55.55", "192.168.44.44"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileExists("panos_anti_spyware_security_profile.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileAttributes(&o, name, "desc one", "ip1.example.com", rName, "any", "any", spyware.ActionResetClient, spyware.SinglePacket, "critical", "low", eName, spyware.Disable, spyware.ActionAllow, "192.168.55.55", "192.168.44.44"),
				),
			},
			{
				Config: testAccAntiSpywareSecurityProfileConfig(name, "desc two", "ip2.example.com", rName, "foo", "adware", spyware.ActionAllow, spyware.Disable, "medium", "high", eName, spyware.ExtendedCapture, spyware.ActionDrop, "192.168.55.55", "192.168.66.66"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileExists("panos_anti_spyware_security_profile.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileAttributes(&o, name, "desc two", "ip2.example.com", rName, "foo", "adware", spyware.ActionAllow, spyware.Disable, "medium", "high", eName, spyware.ExtendedCapture, spyware.ActionDrop, "192.168.55.55", "192.168.66.66"),
				),
			},
		},
	})
}

func testAccCheckPanosAntiSpywareSecurityProfileExists(n string, o *spyware.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v spyware.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseAntiSpywareSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.AntiSpywareProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseAntiSpywareSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.AntiSpywareProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAntiSpywareSecurityProfileAttributes(o *spyware.Entry, name, desc, sink, rName, tn, cat, rAction, rCap, sev1, sev2, eName, eCap, eAction, e1, e2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if o.SinkholeIpv4Address != sink {
			return fmt.Errorf("IPv4 sinkhole is %q, expected %q", o.SinkholeIpv4Address, sink)
		}

		if len(o.Rules) != 1 {
			return fmt.Errorf("rules is not length 1: %#v", o.Rules)
		}

		r := o.Rules[0]
		if r.Name != rName {
			return fmt.Errorf("rule name is %q, expected %q", r.Name, rName)
		}

		if r.ThreatName != tn {
			return fmt.Errorf("rule threat name is %q, not %q", r.ThreatName, tn)
		}

		if r.Category != cat {
			return fmt.Errorf("rule category is %q, not %q", r.Category, cat)
		}

		if r.Action != rAction {
			return fmt.Errorf("rule action is %q, not %q", r.Action, rAction)
		}

		if r.PacketCapture != rCap {
			return fmt.Errorf("rule packet capture is %q, not %q", r.PacketCapture, rCap)
		}

		if len(r.Severities) != 2 {
			return fmt.Errorf("rule severities is not 2: %#v", r.Severities)
		}

		if r.Severities[0] != sev1 {
			return fmt.Errorf("rule sev1 is %q, not %q", r.Severities[0], sev1)
		}

		if r.Severities[1] != sev2 {
			return fmt.Errorf("rule sev2 is %q, not %q", r.Severities[1], sev2)
		}

		if len(o.Exceptions) != 1 {
			return fmt.Errorf("Exceptions is not len 1: %#v", o.Exceptions)
		}

		e := o.Exceptions[0]
		if e.Name != eName {
			return fmt.Errorf("exception name is %q, not %q", e.Name, eName)
		}

		if e.PacketCapture != eCap {
			return fmt.Errorf("exception packet capture is %q, not %q", e.PacketCapture, eCap)
		}

		if e.Action != eAction {
			return fmt.Errorf("exception action is %q, not %q", e.Action, eAction)
		}

		if len(e.ExemptIps) != 2 {
			return fmt.Errorf("exception exempt ips is not len 2: %#v", e.ExemptIps)
		}

		if e.ExemptIps[0] != e1 {
			return fmt.Errorf("exception exempt ip1 is %q, not %q", e.ExemptIps[0], e1)
		}

		if e.ExemptIps[1] != e2 {
			return fmt.Errorf("exception exempt ip2 is %q, not %q", e.ExemptIps[1], e2)
		}

		return nil
	}
}

func testAccPanosAntiSpywareSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_anti_spyware_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseAntiSpywareSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.AntiSpywareProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseAntiSpywareSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.AntiSpywareProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAntiSpywareSecurityProfileConfig(name, desc, sink, rName, tn, rCat, rAction, rCap, sev1, sev2, eName, eCap, eAction, e1, e2 string) string {
	return fmt.Sprintf(`
resource "panos_anti_spyware_security_profile" "test" {
    name = %q
    description = %q
    sinkhole_ipv4_address = %q
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
    rule {
        name = %q
        threat_name = %q
        category = %q
        action = %q
        packet_capture = %q
        severities = [%q, %q]
    }
    exception {
        name = %q
        packet_capture = %q
        action = %q
        exempt_ips = [%q, %q]
    }
}
`, name, desc, sink, rName, tn, rCat, rAction, rCap, sev1, sev2, eName, eCap, eAction, e1, e2)
}
