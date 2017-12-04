package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestPanosSecurityPolicy_basic(t *testing.T) {
	var o security.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPolicyConfig(name, "first description", "10.2.2.2", "10.3.3.3", "allow", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSecurityPolicyExists("panos_security_policy.test", &o),
					testAccCheckPanosSecurityPolicyAttributes(&o, name, "first description", "10.2.2.2", "10.3.3.3", "allow", true, false),
				),
			},
			{
				Config: testAccSecurityPolicyConfig(name, "second description", "10.4.4.4", "10.5.5.5", "drop", false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSecurityPolicyExists("panos_security_policy.test", &o),
					testAccCheckPanosSecurityPolicyAttributes(&o, name, "second description", "10.4.4.4", "10.5.5.5", "drop", false, true),
				),
			},
		},
	})
}

func testAccCheckPanosSecurityPolicyExists(n string, o *security.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, rb, name := parseSecurityPolicyId(rs.Primary.ID)
		v, err := fw.Policies.Security.Get(vsys, rb, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosSecurityPolicyAttributes(o *security.Entry, n, desc, src, dst, a string, le, d bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %q, expected %q", o.Name, n)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if len(o.SourceAddress) != 1 || o.SourceAddress[0] != src {
			return fmt.Errorf("Source address is %#v, expected %#v", o.SourceAddress, []string{src})
		}

		if len(o.DestinationAddress) != 1 || o.DestinationAddress[0] != dst {
			return fmt.Errorf("Destination address is %#v, expected %#v", o.DestinationAddress, []string{dst})
		}

		if o.Action != a {
			return fmt.Errorf("Action is %s, expected %s", o.Action, a)
		}

		if o.LogEnd != le {
			return fmt.Errorf("Log end is %t, expected %t", o.LogEnd, le)
		}

		if o.Disabled != d {
			return fmt.Errorf("Disabled is %t, expected %t", o.Disabled, d)
		}

		return nil
	}
}

func testAccPanosSecurityPolicyDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_security_policy" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, rb, name := parseSecurityPolicyId(rs.Primary.ID)
			_, err := fw.Policies.Security.Get(vsys, rb, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccSecurityPolicyConfig(n, desc, src, dst, a string, le, d bool) string {
	return fmt.Sprintf(`
resource "panos_security_policy" "test" {
    name = "%s"
    description = "%s"
    source_address = ["%s"]
    destination_address = ["%s"]
    action = "%s"
    log_end = %t
    disabled = %t
}
`, n, desc, src, dst, a, le, d)
}
