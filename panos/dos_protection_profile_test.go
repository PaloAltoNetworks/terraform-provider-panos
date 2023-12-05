package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/dos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing tests.
func TestAccPanosDsDosProtectionProfileList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDosProtectionProfileConfig(name),
				Check:  checkDataSourceListing("panos_dos_protection_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsDosProtectionProfile_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDosProtectionProfileConfig(name),
				Check: checkDataSource("panos_dos_protection_profile", []string{
					"name", "description", "type",
					"syn.0.enable",
					"syn.0.action",
					"syn.0.alarm_rate",
					"syn.0.activate_rate",
					"syn.0.max_rate",
					"syn.0.block_duration",
				}),
			},
		},
	})
}

func testAccDsDosProtectionProfileConfig(name string) string {
	return fmt.Sprintf(`
data "panos_dos_protection_profiles" "test" {}

data "panos_dos_protection_profile" "test" {
    name = panos_dos_protection_profile.x.name
}

resource "panos_dos_protection_profile" "x" {
    name = %q
    description = "for dos protection profile data source acctest"
    syn {
        enable = true
        action = %q
        alarm_rate = 777
        activate_rate = 888
        max_rate = 999
        block_duration = 42
    }
}
`, name, dos.SynActionRed)
}

// Resource tests.
func TestAccPanosDosProtectionProfile_basic(t *testing.T) {
	var o dos.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	syns := []*dos.SynProtection{
		&dos.SynProtection{
			Enable:        true,
			Action:        dos.SynActionRed,
			AlarmRate:     7,
			ActivateRate:  8,
			MaxRate:       9,
			BlockDuration: 300,
		},
		&dos.SynProtection{
			Enable:       true,
			Action:       dos.SynActionCookies,
			AlarmRate:    20,
			ActivateRate: 30,
			MaxRate:      40,
		},
		&dos.SynProtection{
			Action:        dos.SynActionRed,
			AlarmRate:     200,
			ActivateRate:  300,
			MaxRate:       400,
			BlockDuration: 500,
		},
		&dos.SynProtection{
			Action:       dos.SynActionCookies,
			AlarmRate:    1000,
			ActivateRate: 2000,
			MaxRate:      3000,
		},
	}

	others := []*dos.Protection{
		&dos.Protection{
			Enable:        true,
			AlarmRate:     7,
			ActivateRate:  8,
			MaxRate:       9,
			BlockDuration: 300,
		},
		&dos.Protection{
			Enable:       true,
			AlarmRate:    20,
			ActivateRate: 30,
			MaxRate:      40,
		},
		&dos.Protection{
			AlarmRate:     200,
			ActivateRate:  300,
			MaxRate:       400,
			BlockDuration: 500,
		},
		&dos.Protection{
			AlarmRate:    1000,
			ActivateRate: 2000,
			MaxRate:      3000,
		},
		&dos.Protection{
			Enable:        true,
			AlarmRate:     2000,
			ActivateRate:  3000,
			MaxRate:       4000,
			BlockDuration: 5000,
		},
		&dos.Protection{
			AlarmRate:    101,
			ActivateRate: 201,
			MaxRate:      301,
		},
		&dos.Protection{
			Enable:        true,
			AlarmRate:     3000,
			ActivateRate:  5000,
			MaxRate:       7000,
			BlockDuration: 20,
		},
		&dos.Protection{
			AlarmRate:    111,
			ActivateRate: 222,
			MaxRate:      333,
		},
		&dos.Protection{
			Enable:        true,
			AlarmRate:     501,
			ActivateRate:  602,
			MaxRate:       703,
			BlockDuration: 804,
		},
		&dos.Protection{
			AlarmRate:    12321,
			ActivateRate: 12421,
			MaxRate:      12521,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDosProtectionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDosProtectionProfileConfig(name, "first", dos.TypeAggregate, syns[0], nil, nil, nil, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "first", dos.TypeAggregate, syns[0], nil, nil, nil, nil),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "second", dos.TypeAggregate, syns[1], nil, nil, nil, others[0]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "second", dos.TypeAggregate, syns[1], nil, nil, nil, others[0]),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "third", dos.TypeAggregate, syns[1], others[0], others[1], others[2], others[3]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "third", dos.TypeAggregate, syns[1], others[0], others[1], others[2], others[3]),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "four", dos.TypeClassified, syns[2], others[3], others[4], others[5], others[6]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "four", dos.TypeClassified, syns[2], others[3], others[4], others[5], others[6]),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "five", dos.TypeClassified, syns[3], others[4], others[5], others[6], others[7]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "five", dos.TypeClassified, syns[3], others[4], others[5], others[6], others[7]),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "six", dos.TypeClassified, nil, nil, nil, nil, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "six", dos.TypeClassified, nil, nil, nil, nil, nil),
				),
			},
			{
				Config: testAccDosProtectionProfileConfig(name, "seven", dos.TypeAggregate, syns[0], others[8], others[9], others[8], others[9]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDosProtectionProfileExists("panos_dos_protection_profile.test", &o),
					testAccCheckPanosDosProtectionProfileAttributes(&o, name, "seven", dos.TypeAggregate, syns[0], others[8], others[9], others[8], others[9]),
				),
			},
		},
	})
}

func testAccCheckPanosDosProtectionProfileExists(n string, o *dos.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v dos.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseDosProtectionProfileId(rs.Primary.ID)
			v, err = con.Objects.DosProtectionProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseDosProtectionProfileId(rs.Primary.ID)
			v, err = con.Objects.DosProtectionProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDosProtectionProfileAttributes(o *dos.Entry, name, desc, tp string, syn *dos.SynProtection, udp, icmp, icmpv6, other *dos.Protection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if o.Type != tp {
			return fmt.Errorf("Type is %q, not %q", o.Type, tp)
		}

		if syn != nil {
			if o.Syn == nil {
				return fmt.Errorf("syn is nil and should be %#v", syn)
			}

			if o.Syn.Enable != syn.Enable {
				return fmt.Errorf("syn enable is %t, not %t", o.Syn.Enable, syn.Enable)
			}

			if o.Syn.Action != syn.Action {
				return fmt.Errorf("syn action is %q, not %q", o.Syn.Action, syn.Action)
			}

			if o.Syn.AlarmRate != syn.AlarmRate {
				return fmt.Errorf("syn alarm rate is %d, not %d", o.Syn.AlarmRate, syn.AlarmRate)
			}

			if o.Syn.ActivateRate != syn.ActivateRate {
				return fmt.Errorf("syn activate rate is %d, not %d", o.Syn.ActivateRate, syn.ActivateRate)
			}

			if o.Syn.MaxRate != syn.MaxRate {
				return fmt.Errorf("syn max rate is %d, not %d", o.Syn.MaxRate, syn.MaxRate)
			}

			if o.Syn.BlockDuration != syn.BlockDuration {
				return fmt.Errorf("syn block duration is %d, not %d", o.Syn.BlockDuration, syn.BlockDuration)
			}
		}

		if udp != nil {
			if o.Udp == nil {
				return fmt.Errorf("udp is nil and should be %#v", udp)
			}

			if o.Udp.Enable != udp.Enable {
				return fmt.Errorf("udp enable is %t, not %t", o.Udp.Enable, udp.Enable)
			}

			if o.Udp.AlarmRate != udp.AlarmRate {
				return fmt.Errorf("udp alarm rate is %d, not %d", o.Udp.AlarmRate, udp.AlarmRate)
			}

			if o.Udp.ActivateRate != udp.ActivateRate {
				return fmt.Errorf("udp activate rate is %d, not %d", o.Udp.ActivateRate, udp.ActivateRate)
			}

			if o.Udp.MaxRate != udp.MaxRate {
				return fmt.Errorf("udp max rate is %d, not %d", o.Udp.MaxRate, udp.MaxRate)
			}

			if o.Udp.BlockDuration != udp.BlockDuration {
				return fmt.Errorf("udp block duration is %d, not %d", o.Udp.BlockDuration, udp.BlockDuration)
			}
		}

		if icmp != nil {
			if o.Icmp == nil {
				return fmt.Errorf("icmp is nil and should be %#v", icmp)
			}

			if o.Icmp.Enable != icmp.Enable {
				return fmt.Errorf("icmp enable is %t, not %t", o.Icmp.Enable, icmp.Enable)
			}

			if o.Icmp.AlarmRate != icmp.AlarmRate {
				return fmt.Errorf("icmp alarm rate is %d, not %d", o.Icmp.AlarmRate, icmp.AlarmRate)
			}

			if o.Icmp.ActivateRate != icmp.ActivateRate {
				return fmt.Errorf("icmp activate rate is %d, not %d", o.Icmp.ActivateRate, icmp.ActivateRate)
			}

			if o.Icmp.MaxRate != icmp.MaxRate {
				return fmt.Errorf("icmp max rate is %d, not %d", o.Icmp.MaxRate, icmp.MaxRate)
			}

			if o.Icmp.BlockDuration != icmp.BlockDuration {
				return fmt.Errorf("icmp block duration is %d, not %d", o.Icmp.BlockDuration, icmp.BlockDuration)
			}
		}

		if icmpv6 != nil {
			if o.Icmpv6 == nil {
				return fmt.Errorf("icmpv6 is nil and should be %#v", icmpv6)
			}

			if o.Icmpv6.Enable != icmpv6.Enable {
				return fmt.Errorf("icmpv6 enable is %t, not %t", o.Icmpv6.Enable, icmpv6.Enable)
			}

			if o.Icmpv6.AlarmRate != icmpv6.AlarmRate {
				return fmt.Errorf("icmpv6 alarm rate is %d, not %d", o.Icmpv6.AlarmRate, icmpv6.AlarmRate)
			}

			if o.Icmpv6.ActivateRate != icmpv6.ActivateRate {
				return fmt.Errorf("icmpv6 activate rate is %d, not %d", o.Icmpv6.ActivateRate, icmpv6.ActivateRate)
			}

			if o.Icmpv6.MaxRate != icmpv6.MaxRate {
				return fmt.Errorf("icmpv6 max rate is %d, not %d", o.Icmpv6.MaxRate, icmpv6.MaxRate)
			}

			if o.Icmpv6.BlockDuration != icmpv6.BlockDuration {
				return fmt.Errorf("icmpv6 block duration is %d, not %d", o.Icmpv6.BlockDuration, icmpv6.BlockDuration)
			}
		}

		if other != nil {
			if o.Other == nil {
				return fmt.Errorf("other is nil and should be %#v", other)
			}

			if o.Other.Enable != other.Enable {
				return fmt.Errorf("other enable is %t, not %t", o.Other.Enable, other.Enable)
			}

			if o.Other.AlarmRate != other.AlarmRate {
				return fmt.Errorf("other alarm rate is %d, not %d", o.Other.AlarmRate, other.AlarmRate)
			}

			if o.Other.ActivateRate != other.ActivateRate {
				return fmt.Errorf("other activate rate is %d, not %d", o.Other.ActivateRate, other.ActivateRate)
			}

			if o.Other.MaxRate != other.MaxRate {
				return fmt.Errorf("other max rate is %d, not %d", o.Other.MaxRate, other.MaxRate)
			}

			if o.Other.BlockDuration != other.BlockDuration {
				return fmt.Errorf("other block duration is %d, not %d", o.Other.BlockDuration, other.BlockDuration)
			}
		}

		return nil
	}
}

func testAccPanosDosProtectionProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_dos_protection_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseDosProtectionProfileId(rs.Primary.ID)
				_, err = con.Objects.DosProtectionProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseDosProtectionProfileId(rs.Primary.ID)
				_, err = con.Objects.DosProtectionProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccDosProtectionProfileConfig(name, desc, tp string, syn *dos.SynProtection, udp, icmp, icmpv6, other *dos.Protection) string {
	ans := fmt.Sprintf(`
resource "panos_dos_protection_profile" "test" {
    name = %q
    description = %q
    type = %q`, name, desc, tp)

	if syn != nil {
		ans = fmt.Sprintf(`%s
    syn {
        enable = %t
        action = %q
        alarm_rate = %d
        activate_rate = %d
        max_rate = %d
        block_duration = %d
    }`, ans, syn.Enable, syn.Action, syn.AlarmRate, syn.ActivateRate, syn.MaxRate, syn.BlockDuration)
	}

	if udp != nil {
		ans = fmt.Sprintf(`%s
    udp {
        enable = %t
        alarm_rate = %d
        activate_rate = %d
        max_rate = %d
        block_duration = %d
    }`, ans, udp.Enable, udp.AlarmRate, udp.ActivateRate, udp.MaxRate, udp.BlockDuration)
	}

	if icmp != nil {
		ans = fmt.Sprintf(`%s
    icmp {
        enable = %t
        alarm_rate = %d
        activate_rate = %d
        max_rate = %d
        block_duration = %d
    }`, ans, icmp.Enable, icmp.AlarmRate, icmp.ActivateRate, icmp.MaxRate, icmp.BlockDuration)
	}

	if icmpv6 != nil {
		ans = fmt.Sprintf(`%s
    icmpv6 {
        enable = %t
        alarm_rate = %d
        activate_rate = %d
        max_rate = %d
        block_duration = %d
    }`, ans, icmpv6.Enable, icmpv6.AlarmRate, icmpv6.ActivateRate, icmpv6.MaxRate, icmpv6.BlockDuration)
	}

	if other != nil {
		ans = fmt.Sprintf(`%s
    other {
        enable = %t
        alarm_rate = %d
        activate_rate = %d
        max_rate = %d
        block_duration = %d
    }`, ans, other.Enable, other.AlarmRate, other.ActivateRate, other.MaxRate, other.BlockDuration)
	}

	ans = fmt.Sprintf(`%s
}`, ans)

	return ans
}
