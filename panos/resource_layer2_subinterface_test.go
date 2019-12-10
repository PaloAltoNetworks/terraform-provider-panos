package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosLayer2Subinterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o layer2.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/5.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosLayer2SubinterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLayer2SubinterfaceConfig(name, "desc1", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLayer2SubinterfaceExists("panos_layer2_subinterface.test", &o),
					testAccCheckPanosLayer2SubinterfaceAttributes(&o, name, "desc1", 5),
				),
			},
			{
				Config: testAccLayer2SubinterfaceConfig(name, "desc2", 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLayer2SubinterfaceExists("panos_layer2_subinterface.test", &o),
					testAccCheckPanosLayer2SubinterfaceAttributes(&o, name, "desc2", 7),
				),
			},
		},
	})
}

func testAccCheckPanosLayer2SubinterfaceExists(n string, o *layer2.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		iType, eth, mType, _, name := parseLayer2SubinterfaceId(rs.Primary.ID)
		v, err := fw.Network.Layer2Subinterface.Get(iType, eth, mType, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosLayer2SubinterfaceAttributes(o *layer2.Entry, name, com string, tag int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Comment != com {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, com)
		}

		if o.Tag != tag {
			return fmt.Errorf("Tag is %d, not %d", o.Tag, tag)
		}

		return nil
	}
}

func testAccPanosLayer2SubinterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_layer2_subinterface" {
			continue
		}

		if rs.Primary.ID != "" {
			iType, eth, mType, _, name := parseLayer2SubinterfaceId(rs.Primary.ID)
			_, err := fw.Network.Layer2Subinterface.Get(iType, eth, mType, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccLayer2SubinterfaceConfig(name, com string, tag int) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "x" {
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer2"
    comment = "for layer2 test"
}

resource "panos_layer2_subinterface" "test" {
    name = %q
    parent_interface = panos_ethernet_interface.x.name
    parent_mode = panos_ethernet_interface.x.mode
    comment = %q
    tag = %d
}
`, name, com, tag)
}
