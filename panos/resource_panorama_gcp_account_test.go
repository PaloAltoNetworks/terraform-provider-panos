package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/account"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaGcpAccount_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if x := testAccPlugins["gcp"]; x == "" {
		t.Skip("The GCP plugin must be installed to run this test")
	} else if !strings.HasPrefix(testAccPlugins["gcp"], "1.") {
		t.Skip("GCP Plugin should be version 1 for this acctest")
	}

	var o account.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaGcpAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaGcpAccountConfig(name, "desc1", "project-1", account.Project, `{"foo": "bar", "bar": "baz"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGcpAccountExists("panos_panorama_gcp_account.test", &o),
					testAccCheckPanosPanoramaGcpAccountAttributes(&o, name, "desc1", "project-1", account.Project),
				),
			},
			{
				Config: testAccPanoramaGcpAccountConfig(name, "desc2", "project-followup", account.Gke, `{"terraform": "acceptance", "acceptance": "test"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGcpAccountExists("panos_panorama_gcp_account.test", &o),
					testAccCheckPanosPanoramaGcpAccountAttributes(&o, name, "desc2", "project-followup", account.Gke),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaGcpAccountExists(n string, o *account.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		name := rs.Primary.ID
		v, err := pano.Panorama.GcpAccount.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaGcpAccountAttributes(o *account.Entry, name, desc, pi, sact string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.ProjectId != pi {
			return fmt.Errorf("Project ID is %s, expected %s", o.ProjectId, pi)
		}

		if o.ServiceAccountCredentialType != sact {
			return fmt.Errorf("Service account credential type is %s, expected %s", o.ServiceAccountCredentialType, sact)
		}

		return nil
	}
}

func testAccPanosPanoramaGcpAccountDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_gcp_account" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := pano.Panorama.GcpAccount.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaGcpAccountConfig(name, desc, pi, sact, cf string) string {
	return fmt.Sprintf(`
resource "panos_panorama_gcp_account" "test" {
    name = %q
    description = %q
    project_id = %q
    service_account_credential_type = %q
    credential_file = %q
}
`, name, desc, pi, sact, cf)
}
