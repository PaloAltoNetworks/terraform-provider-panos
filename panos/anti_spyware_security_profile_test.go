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
func TestAccPanosDsFirewallAntiSpywareSecurityProfileList(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profiles"),
			},
		},
	})
}

func TestAccPanosDsPanoramaAntiSpywareSecurityProfileList(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsFirewallAntiSpywareSecurityProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name),
				Check: checkDataSource("panos_anti_spyware_security_profile", []string{
					"name", "description", "sinkhole_ipv4_address", "sinkhole_ipv6_address",
				}),
			},
		},
	})
}

func TestAccPanosDsPanoramaAntiSpywareSecurityProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileConfig(name),
				Check: checkDataSource("panos_anti_spyware_security_profile", []string{
					"name", "description", "sinkhole_ipv4_address", "sinkhole_ipv6_address",
				}),
			},
		},
	})
}

func testAccDsAntiSpywareSecurityProfileConfig(name string) string {
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
}
`, name)
}

// Resource tests.
func TestAccPanosFirewallAntiSpywareSecurityProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
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

func testAccCheckPanosAntiSpywareSecurityProfileAttributes(o *spyware.Entry, name, desc, sink string) resource.TestCheckFunc {
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

func testAccAntiSpywareSecurityProfileConfig(name, desc, sink string) string {
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
}
`, name, desc, sink)
}
