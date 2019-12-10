package panos

import (
	"fmt"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

/*
In order to run this test you'll need:

    * licensing api key (PANOS_LICENSE_API_KEY)

License API keys can be downloaded from the support site as per this:

    https://www.paloaltonetworks.com/documentation/71/virtualization/virtualization/license-the-vm-series-firewall/install-a-license-deactivation-api-key#46350

*/

type Ans struct {
	Value string
}

func TestAccPanosLicenseApiKey_basic(t *testing.T) {
	key := os.Getenv("PANOS_LICENSE_API_KEY")

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if key == "" {
		t.Skip("PANOS_LICENSE_API_KEY must be set")
	}

	o := Ans{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosLicenseApiKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLicenseApiKeyConfig(key, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLicenseApiKeyExists("panos_license_api_key.test", &o),
					testAccCheckPanosLicenseApiKeyAttributes(&o, key),
				),
			},
			{
				Config: testAccLicenseApiKeyConfig(key, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLicenseApiKeyExists("panos_license_api_key.test", &o),
					testAccCheckPanosLicenseApiKeyAttributes(&o, key),
				),
			},
		},
	})
}

func testAccCheckPanosLicenseApiKeyExists(n string, o *Ans) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		key, err := fw.Licensing.GetApiKey()
		if err != nil {
			return fmt.Errorf("Error in get api key: %s", err)
		} else if key == "" {
			return fmt.Errorf("No licensing api key set")
		}

		o.Value = key
		return nil
	}
}

func testAccCheckPanosLicenseApiKeyAttributes(o *Ans, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Value != key {
			return fmt.Errorf("License key is %q, not %q", o.Value, key)
		}
		return nil
	}
}

func testAccPanosLicenseApiKeyDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_license_api_key" {
			continue
		}

		if rs.Primary.ID != "" {
			key, err := fw.Licensing.GetApiKey()
			if err != nil {
				return err
			} else if key == "" && rs.Primary.ID == "true" {
				return fmt.Errorf("Licensing API key was incorrectly deleted")
			} else if key != "" && rs.Primary.ID == "false" {
				return fmt.Errorf("Licensing API key was incorrectly retained")
			}
		}
		return nil
	}

	return nil
}

func testAccLicenseApiKeyConfig(key string, keep bool) string {
	return fmt.Sprintf(`
resource "panos_license_api_key" "test" {
    key = %q
    retain_key = %t
}
`, key, keep)
}
