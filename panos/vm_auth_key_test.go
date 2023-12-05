package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source test.
func TestAccPanosDsVmAuthKey(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsVmAuthKeyConfig(),
				Check: checkDataSource("panos_vm_auth_key", []string{
					"total", "entries.0.auth_key", "entries.0.expiry",
				}),
			},
		},
	})
}

func testAccDsVmAuthKeyConfig() string {
	return `
data "panos_vm_auth_key" "test" {}
resource "panos_vm_auth_key" "x" {}
`
}

// Resource test.
func TestAccPanosVmAuthKey(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o pango.VmAuthKey
	hours := acctest.RandInt()%4 + 8

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosVmAuthKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmAuthKeyConfig(hours),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVmAuthKeyExists("panos_vm_auth_key.test", &o),
					testAccCheckPanosVmAuthKeyAttributes(&o),
				),
			},
		},
	})
}

func testAccCheckPanosVmAuthKeyExists(n string, o *pango.VmAuthKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		authKey := rs.Primary.ID
		var err error
		var list []pango.VmAuthKey

		switch con := testAccProvider.Meta().(type) {
		case *pango.Panorama:
			list, err = con.GetVmAuthKeys()
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		for _, v := range list {
			if v.AuthKey == authKey {
				*o = v
				break
			}
		}

		if o == nil {
			return fmt.Errorf("Didn't find %q", authKey)
		}

		return nil
	}
}

func testAccCheckPanosVmAuthKeyAttributes(o *pango.VmAuthKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.AuthKey == "" {
			return fmt.Errorf("Auth key is empty")
		}

		if o.Expiry == "" {
			return fmt.Errorf("Expiry is empty")
		}

		return nil
	}
}

func testAccPanosVmAuthKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_device_group_parent" {
			continue
		}

		if rs.Primary.ID != "" {
			authKey := rs.Primary.ID
			var err error
			var list []pango.VmAuthKey

			switch con := testAccProvider.Meta().(type) {
			case *pango.Panorama:
				list, err = con.GetVmAuthKeys()
				if err != nil {
					return err
				}
			}
			for _, v := range list {
				if v.AuthKey == authKey {
					return fmt.Errorf("Auth key %q still present", authKey)
				}
			}
		}
		return nil
	}

	return nil
}

func testAccVmAuthKeyConfig(hours int) string {
	return fmt.Sprintf(`
resource "panos_vm_auth_key" "test" {
    hours = %d
}
`, hours)
}
