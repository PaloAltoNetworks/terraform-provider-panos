package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaSecurityPolicy_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o1, o2 security.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	name1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	name2 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaSecurityPolicyConfig(dg, name1, "first description", "10.2.2.2", "10.3.3.3", "allow", true, false, name2, "another first", "192.168.1.1", "192.168.3.3", "deny", false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSecurityPolicyExists("panos_panorama_security_policy.test", &o1, &o2),
					testAccCheckPanosPanoramaSecurityPolicyAttributes(&o1, &o2, name1, "first description", "10.2.2.2", "10.3.3.3", "allow", true, false, name2, "another first", "192.168.1.1", "192.168.3.3", "deny", false, true),
				),
			},
			{
				Config: testAccPanoramaSecurityPolicyConfig(dg, name1, "second description", "10.4.4.4", "10.5.5.5", "drop", false, true, name2, "next description", "192.168.2.2", "192.168.4.4", "allow", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSecurityPolicyExists("panos_panorama_security_policy.test", &o1, &o2),
					testAccCheckPanosPanoramaSecurityPolicyAttributes(&o1, &o2, name1, "second description", "10.4.4.4", "10.5.5.5", "drop", false, true, name2, "next description", "192.168.2.2", "192.168.4.4", "allow", true, false),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaSecurityPolicyExists(n string, o1, o2 *security.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, rb := parsePanoramaSecurityPolicyId(rs.Primary.ID)
		list, err := pano.Policies.Security.GetList(dg, rb)
		if err != nil {
			return fmt.Errorf("Error getting list of policies: %s", err)
		} else if len(list) != 2 {
			return fmt.Errorf("Expecting 2 policies, got %d", len(list))
		}

		v1, err := pano.Policies.Security.Get(dg, rb, list[0])
		if err != nil {
			return fmt.Errorf("Error getting first policy %s: %s", list[0], err)
		}
		v2, err := pano.Policies.Security.Get(dg, rb, list[1])
		if err != nil {
			return fmt.Errorf("Error getting second policy %s: %s", list[1], err)
		}

		*o1 = v1
		*o2 = v2

		return nil
	}
}

func testAccCheckPanosPanoramaSecurityPolicyAttributes(o1, o2 *security.Entry, name1, desc1, src1, dst1, action1 string, le1, dis1 bool, name2, desc2, src2, dst2, action2 string, le2, dis2 bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o1.Name != name1 {
			return fmt.Errorf("Name is %q, expected %q", o1.Name, name1)
		}

		if o1.Description != desc1 {
			return fmt.Errorf("Description is %q, expected %q", o1.Description, desc1)
		}

		if len(o1.SourceAddresses) != 1 || o1.SourceAddresses[0] != src1 {
			return fmt.Errorf("Source address is %#v, expected %#v", o1.SourceAddresses, []string{src1})
		}

		if len(o1.DestinationAddresses) != 1 || o1.DestinationAddresses[0] != dst1 {
			return fmt.Errorf("Destination address is %#v, expected %#v", o1.DestinationAddresses, []string{dst1})
		}

		if o1.Action != action1 {
			return fmt.Errorf("Action is %s, expected %s", o1.Action, action1)
		}

		if o1.LogEnd != le1 {
			return fmt.Errorf("Log end is %t, expected %t", o1.LogEnd, le1)
		}

		if o1.Disabled != dis1 {
			return fmt.Errorf("Disabled is %t, expected %t", o1.Disabled, dis1)
		}

		if o2.Name != name2 {
			return fmt.Errorf("Name is %q, expected %q", o2.Name, name2)
		}

		if o2.Description != desc2 {
			return fmt.Errorf("Description is %q, expected %q", o2.Description, desc2)
		}

		if len(o2.SourceAddresses) != 1 || o2.SourceAddresses[0] != src2 {
			return fmt.Errorf("Source address is %#v, expected %#v", o2.SourceAddresses, []string{src2})
		}

		if len(o2.DestinationAddresses) != 1 || o2.DestinationAddresses[0] != dst2 {
			return fmt.Errorf("Destination address is %#v, expected %#v", o2.DestinationAddresses, []string{dst2})
		}

		if o2.Action != action2 {
			return fmt.Errorf("Action is %s, expected %s", o2.Action, action2)
		}

		if o2.LogEnd != le2 {
			return fmt.Errorf("Log end is %t, expected %t", o2.LogEnd, le2)
		}

		if o2.Disabled != dis2 {
			return fmt.Errorf("Disabled is %t, expected %t", o2.Disabled, dis2)
		}
		return nil
	}
}

func testAccPanosPanoramaSecurityPolicyDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_security_policy" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, rb := parsePanoramaSecurityPolicyId(rs.Primary.ID)
			list, err := pano.Policies.Security.GetList(dg, rb)
			if err != nil {
				return fmt.Errorf("Error getting list: %s", err)
			} else if len(list) != 0 {
				return fmt.Errorf("%d security policies still exist", len(list))
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaSecurityPolicyConfig(dg, name1, desc1, src1, dst1, action1 string, le1, dis1 bool, name2, desc2, src2, dst2, action2 string, le2, dis2 bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_security_policy" "test" {
    device_group = panos_panorama_device_group.x.name
    rule {
        name = "%s"
        description = "%s"
        source_addresses = ["%s"]
        destination_addresses = ["%s"]
        action = "%s"
        log_end = %t
        disabled = %t
        source_zones = ["any"]
        destination_zones = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
    }
    rule {
        name = "%s"
        description = "%s"
        source_addresses = ["%s"]
        destination_addresses = ["%s"]
        action = "%s"
        log_end = %t
        disabled = %t
        source_zones = ["any"]
        destination_zones = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
    }
}
`, dg, name1, desc1, src1, dst1, action1, le1, dis1, name2, desc2, src2, dst2, action2, le2, dis2)
}
