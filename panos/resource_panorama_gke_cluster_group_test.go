package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/pnrm/plugins/gcp/gke/cluster/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaGkeClusterGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if x := testAccPlugins["gcp"]; x == "" {
		t.Skip("The GCP plugin must be installed to run this test")
	} else if !strings.HasPrefix(testAccPlugins["gcp"], "1.") {
		t.Skip("GCP Plugin should be version 1 for this acctest")
	}

	var o group.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	ts := fmt.Sprintf("tf%s", acctest.RandString(6))
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaGkeClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaGkeClusterGroupConfig(ts, dg, name, "desc1", "xAcc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGkeClusterGroupExists("panos_panorama_gke_cluster_group.test", &o),
					testAccCheckPanosPanoramaGkeClusterGroupAttributes(&o, name, "desc1", "xAcc"),
				),
			},
			{
				Config: testAccPanoramaGkeClusterGroupConfig(ts, dg, name, "desc2", "yAcc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGkeClusterGroupExists("panos_panorama_gke_cluster_group.test", &o),
					testAccCheckPanosPanoramaGkeClusterGroupAttributes(&o, name, "desc2", "yAcc"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaGkeClusterGroupExists(n string, o *group.Entry) resource.TestCheckFunc {
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
		v, err := pano.Panorama.GkeClusterGroup.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaGkeClusterGroupAttributes(o *group.Entry, name, desc, ga string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.GcpProjectCredential != ga {
			return fmt.Errorf("GCP project credential is %s, expected %s", o.GcpProjectCredential, ga)
		}

		return nil
	}
}

func testAccPanosPanoramaGkeClusterGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_gke_cluster_group" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := pano.Panorama.GkeClusterGroup.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaGkeClusterGroupConfig(ts, dg, name, desc, ga string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
}

resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_gcp_account" "xAcc" {
    name = "xAcc"
    project_id = "example-project"
    service_account_credential_type = "gcp"
    credential_file = "{\"json\": \"data\"}"
}

resource "panos_panorama_gcp_account" "yAcc" {
    name = "yAcc"
    project_id = "example-project"
    service_account_credential_type = "gcp"
    credential_file = "{\"json\": \"data\"}"
}

resource "panos_panorama_gke_cluster_group" "test" {
    name = %q
    description = %q
    gcp_project_credential = panos_panorama_gcp_account.%s.name
    device_group = panos_panorama_device_group.x.name
    template_stack = panos_panorama_template_stack.x.name
}
`, ts, dg, name, desc, ga)
}
