package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/arp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsArpEthernetList(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%4+1)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpEthernetConfig(tmpl, iName, ip),
				Check:  checkDataSourceListing("panos_arps"),
			},
		},
	})
}

func TestAccPanosDsArpVlanList(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	subName := fmt.Sprintf("vlan.%d", acctest.RandInt()%5+5)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpVlanConfig(tmpl, subName, ip),
				Check:  checkDataSourceListing("panos_arps"),
			},
		},
	})
}

func TestAccPanosDsArpAggregateList(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	} else if !testAccSupportsAggregateInterfaces {
		t.Skip(SkipAggregateTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ae%d", acctest.RandInt()%3+2)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpAggregateConfig(tmpl, iName, ip),
				Check:  checkDataSourceListing("panos_arps"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsArpEthernet_basic(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%4+1)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpEthernetConfig(tmpl, iName, ip),
				Check: checkDataSource("panos_arp", []string{
					"ip", "mac_address",
				}),
			},
		},
	})
}

func TestAccPanosDsArpVlan_basic(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	subName := fmt.Sprintf("vlan.%d", acctest.RandInt()%5+5)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpVlanConfig(tmpl, subName, ip),
				Check: checkDataSource("panos_arp", []string{
					"ip", "mac_address",
				}),
			},
		},
	})
}

func TestAccPanosDsArpAggregate_basic(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	} else if !testAccSupportsAggregateInterfaces {
		t.Skip(SkipAggregateTest)
	}

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ae%d", acctest.RandInt()%3+2)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsArpAggregateConfig(tmpl, iName, ip),
				Check: checkDataSource("panos_arp", []string{
					"ip", "mac_address",
				}),
			},
		},
	})
}

func testAccDsArpEthernetConfig(tmpl, iName, ip string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_arps" "test" {
    template = panos_panorama_template.x.name
    interface_type = "ethernet"
    interface_name = panos_panorama_ethernet_interface.x.name
}

data "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = "ethernet"
    interface_name = panos_panorama_ethernet_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    template = panos_panorama_template.x.name
    interface_type = "ethernet"
    interface_name = panos_panorama_ethernet_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, ip, tmpl, iName)
	}

	return fmt.Sprintf(`
data "panos_arps" "test" {
    interface_type = "ethernet"
    interface_name = panos_ethernet_interface.x.name
}

data "panos_arp" "test" {
    interface_type = "ethernet"
    interface_name = panos_ethernet_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    interface_type = "ethernet"
    interface_name = panos_ethernet_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_ethernet_interface" "x" {
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, ip, iName)
}

func testAccDsArpAggregateConfig(tmpl, iName, ip string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_arps" "test" {
    template = panos_panorama_template.x.name
    interface_type = "aggregate-ethernet"
    interface_name = panos_panorama_aggregate_interface.x.name
}

data "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = "aggregate-ethernet"
    interface_name = panos_panorama_aggregate_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    template = panos_panorama_template.x.name
    interface_type = "aggregate-ethernet"
    interface_name = panos_panorama_aggregate_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_aggregate_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, ip, tmpl, iName)
	}

	return fmt.Sprintf(`
data "panos_arps" "test" {
    interface_type = "aggregate-ethernet"
    interface_name = panos_aggregate_interface.x.name
}

data "panos_arp" "test" {
    interface_type = "aggregate-ethernet"
    interface_name = panos_aggregate_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    interface_type = "aggregate-ethernet"
    interface_name = panos_aggregate_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_aggregate_interface" "x" {
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, ip, iName)
}

func testAccDsArpVlanConfig(tmpl, subName, ip string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_arps" "test" {
    template = panos_panorama_template.x.name
    interface_type = "vlan"
    subinterface_name = panos_panorama_vlan_interface.x.name
}

data "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = "vlan"
    subinterface_name = panos_panorama_vlan_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    template = panos_panorama_template.x.name
    interface_type = "vlan"
    subinterface_name = panos_panorama_vlan_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_vlan_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
}
`, ip, tmpl, subName)
	}

	return fmt.Sprintf(`
data "panos_arps" "test" {
    interface_type = "vlan"
    subinterface_name = panos_vlan_interface.x.name
}

data "panos_arp" "test" {
    interface_type = "vlan"
    subinterface_name = panos_vlan_interface.x.name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    interface_type = "vlan"
    subinterface_name = panos_vlan_interface.x.name
    ip = %q
    mac_address = "00:30:48:52:ab:cd"
}

resource "panos_vlan_interface" "x" {
    name = %q
    vsys = "vsys1"
}
`, ip, subName)
}

// Resource tests.
func TestAccPanosArpEthernet(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o arp.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%4+1)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosArpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArpEthernetConfig(tmpl, iName, ip, "00:30:48:52:00:01"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:01", ""),
				),
			},
			{
				Config: testAccArpEthernetConfig(tmpl, iName, ip, "00:30:48:52:00:02"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:02", ""),
				),
			},
		},
	})
}

func TestAccPanosArpVlan(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o arp.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	subName := fmt.Sprintf("vlan.%d", acctest.RandInt()%5+5)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)
	l2i := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%4+1)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosArpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArpVlanConfig(tmpl, subName, ip, "00:30:48:52:00:01", l2i, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:01", ""),
				),
			},
			{
				Config: testAccArpVlanConfig(tmpl, subName, ip, "00:30:48:52:00:02", l2i, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:02", l2i),
				),
			},
		},
	})
}

func TestAccPanosArpAggregate(t *testing.T) {
	if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	} else if !testAccSupportsAggregateInterfaces {
		t.Skip(SkipAggregateTest)
	}

	var o arp.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	iName := fmt.Sprintf("ae%d", acctest.RandInt()%3+2)
	ip := fmt.Sprintf("10.1.1.%d", acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosArpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccArpAggregateConfig(tmpl, iName, ip, "00:30:48:52:00:01"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:01", ""),
				),
			},
			{
				Config: testAccArpAggregateConfig(tmpl, iName, ip, "00:30:48:52:00:02"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosArpExists("panos_arp.test", &o),
					testAccCheckPanosArpAttributes(&o, ip, "00:30:48:52:00:02", ""),
				),
			},
		},
	})
}

func testAccCheckPanosArpExists(n string, o *arp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v arp.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			iType, iName, subName, ip := parseFirewallArpId(rs.Primary.ID)
			v, err = con.Network.Arp.Get(iType, iName, subName, ip)
		case *pango.Panorama:
			tmpl, ts, iType, iName, subName, ip := parsePanoramaArpId(rs.Primary.ID)
			v, err = con.Network.Arp.Get(tmpl, ts, iType, iName, subName, ip)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosArpAttributes(o *arp.Entry, ip, mac, iface string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Ip != ip {
			return fmt.Errorf("Ip is %q, not %q", o.Ip, ip)
		}

		if o.MacAddress != mac {
			return fmt.Errorf("Mac address is %q, not %q", o.MacAddress, mac)
		}

		if o.Interface != iface {
			return fmt.Errorf("Interface is %q, not %q", o.Interface, iface)
		}

		return nil
	}
}

func testAccPanosArpDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_address_object" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				iType, iName, subName, ip := parseFirewallArpId(rs.Primary.ID)
				_, err = con.Network.Arp.Get(iType, iName, subName, ip)
			case *pango.Panorama:
				tmpl, ts, iType, iName, subName, ip := parsePanoramaArpId(rs.Primary.ID)
				_, err = con.Network.Arp.Get(tmpl, ts, iType, iName, subName, ip)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccArpEthernetConfig(tmpl, iName, ip, mac string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = %q
    interface_name = panos_panorama_ethernet_interface.x.name
    ip = %q
    mac_address = %q
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, arp.TypeEthernet, ip, mac, tmpl, iName)
	}

	return fmt.Sprintf(`
resource "panos_arp" "test" {
    interface_type = %q
    interface_name = panos_ethernet_interface.x.name
    ip = %q
    mac_address = %q
}

resource "panos_ethernet_interface" "x" {
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, arp.TypeEthernet, ip, mac, iName)
}

func testAccArpAggregateConfig(tmpl, iName, ip, mac string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = %q
    interface_name = panos_panorama_aggregate_interface.x.name
    ip = %q
    mac_address = %q
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_aggregate_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, arp.TypeAggregate, ip, mac, tmpl, iName)
	}

	return fmt.Sprintf(`
resource "panos_arp" "test" {
    interface_type = %q
    interface_name = panos_aggregate_interface.x.name
    ip = %q
    mac_address = %q
}

resource "panos_aggregate_interface" "x" {
    name = %q
    vsys = "vsys1"
    mode = "layer3"
}
`, arp.TypeAggregate, ip, mac, iName)
}

func testAccArpVlanConfig(tmpl, subName, ip, mac, l2i string, link bool) string {
	var interface_line string

	if testAccIsPanorama {
		if link {
			interface_line = "    interface = panos_panorama_ethernet_interface.x.name"
		}
		return fmt.Sprintf(`
resource "panos_arp" "test" {
    template = panos_panorama_template.x.name
    interface_type = %q
    subinterface_name = panos_panorama_vlan_interface.x.name
    ip = %q
    mac_address = %q
%s
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_vlan_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    vsys = "vsys1"
    mode = "layer2"
}
`, arp.TypeVlan, ip, mac, interface_line, tmpl, subName, l2i)
	}

	if link {
		interface_line = "    interface = panos_ethernet_interface.x.name"
	}
	return fmt.Sprintf(`
resource "panos_arp" "test" {
    interface_type = %q
    subinterface_name = panos_vlan_interface.x.name
    ip = %q
    mac_address = %q
%s
}

resource "panos_vlan_interface" "x" {
    name = %q
    vsys = "vsys1"
}

resource "panos_ethernet_interface" "x" {
    name = %q
    vsys = "vsys1"
    mode = "layer2"
}
`, arp.TypeVlan, ip, mac, interface_line, subName, l2i)
}
