package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/plugins/gcp/gke/cluster"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaGkeCluster_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if x := testAccPanoramaPlugins["gcp"]; x == "" {
		t.Skip("The GCP plugin must be installed to run this test")
	}

	var o cluster.Entry
	ts := fmt.Sprintf("tf%s", acctest.RandString(6))
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	ga := fmt.Sprintf("tf%s", acctest.RandString(6))
	grp := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaGkeClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaGkeClusterConfig(ts, dg, ga, grp, name, "zone1", "xGke"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGkeClusterExists("panos_panorama_gke_cluster.test", &o),
					testAccCheckPanosPanoramaGkeClusterAttributes(&o, name, "zone1", "xGke"),
				),
			},
			{
				Config: testAccPanoramaGkeClusterConfig(ts, dg, ga, grp, name, "zone2", "yGke"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGkeClusterExists("panos_panorama_gke_cluster.test", &o),
					testAccCheckPanosPanoramaGkeClusterAttributes(&o, name, "zone2", "yGke"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaGkeClusterExists(n string, o *cluster.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		grp, name := parseGkeClusterId(rs.Primary.ID)
		v, err := pano.Panorama.GkeCluster.Get(grp, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaGkeClusterAttributes(o *cluster.Entry, name, zone, cc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.GcpZone != zone {
			return fmt.Errorf("GCP zone is %s, expected %s", o.GcpZone, zone)
		}

		if o.ClusterCredential != cc {
			return fmt.Errorf("Cluster credential is %s, expected %s", o.ClusterCredential, cc)
		}

		return nil
	}
}

func testAccPanosPanoramaGkeClusterDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_gke_cluster" {
			continue
		}

		if rs.Primary.ID != "" {
			grp, name := parseGkeClusterId(rs.Primary.ID)
			_, err := pano.Panorama.GkeCluster.Get(grp, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaGkeClusterConfig(ts, dg, ga, grp, name, zone, cc string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
}

resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_gcp_account" "x" {
    name = %q
    project_id = "example-project"
    service_account_credential_type = "gcp"
    credential_file = "{\"json\": \"data\"}"
}

resource "panos_panorama_gke_cluster_group" "x" {
    name = %q
    gcp_project_credential = panos_panorama_gcp_account.x.name
    device_group = panos_panorama_device_group.x.name
    template_stack = panos_panorama_template_stack.x.name
}

resource "panos_panorama_gcp_account" "xGke" {
    name = "xGke"
    project_id = "example-project"
    service_account_credential_type = "gke"
    credential_file = "{\"json\": \"data\"}"
}

resource "panos_panorama_gcp_account" "yGke" {
    name = "yGke"
    project_id = "example-project"
    service_account_credential_type = "gke"
    credential_file = "{\"json\": \"data\"}"
}

resource "panos_panorama_gke_cluster" "test" {
    gke_cluster_group = panos_panorama_gke_cluster_group.x.name
    name = %q
    gcp_zone = %q
    cluster_credential = panos_panorama_gcp_account.%s.name
}
`, ts, dg, ga, grp, name, zone, cc)
}
