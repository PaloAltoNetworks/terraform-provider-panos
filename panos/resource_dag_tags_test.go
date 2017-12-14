package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestPanosDagTags_basic(t *testing.T) {
	o := make(map[string][]string)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccPanosDagTagsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDagTagsConfig("10.2.2.2", "tag1", "tag2", "10.3.3.3", "tag3", "10.4.4.4", "tag1", "tag3", "tag5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDagTagsExists("panos_dag_tags.test", o),
					testAccCheckPanosDagTagsAttributes(o, "10.2.2.2", []string{"tag1", "tag2"}, "10.3.3.3", []string{"tag3"}, "10.4.4.4", []string{"tag1", "tag3", "tag5"}),
				),
			},
		},
	})
}

func testAccCheckPanosDagTagsExists(n string, o map[string][]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.UserId.Registered("", "", rs.Primary.ID)
		if err != nil {
			return err
		}

		for key := range o {
			delete(o, key)
		}
		for key, value := range v {
			o[key] = value
		}

		return nil
	}
}

func testAccCheckPanosDagTagsAttributes(o map[string][]string, ip1 string, ip1t []string, ip2 string, ip2t []string, ip3 string, ip3t []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		info, ok := o[ip1]
		if !ok {
			return fmt.Errorf("IP1 %q not in results: %#v", ip1, o)
		} else {
			for _, st := range ip1t {
				found := false
				for _, tag := range info {
					if tag == st {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("IP1 tags do not include %v", ip1t)
				}
			}
		}

		info, ok = o[ip2]
		if !ok {
			return fmt.Errorf("IP2 %q not in results", ip2)
		} else {
			for _, st := range ip2t {
				found := false
				for _, tag := range info {
					if tag == st {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("IP2 tags do not include %v", ip2t)
				}
			}
		}

		info, ok = o[ip3]
		if !ok {
			return fmt.Errorf("IP3 %q not in results", ip3)
		} else {
			for _, st := range ip3t {
				found := false
				for _, tag := range info {
					if tag == st {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("IP3 tags do not include %v", ip3t)
				}
			}
		}

		return nil
	}
}

func testAccPanosDagTagsDestroy(s *terraform.State) error {
	return nil
}

func testAccDagTagsConfig(ip1, ip1t1, ip1t2, ip2, ip2t1, ip3, ip3t1, ip3t2, ip3t3 string) string {
	return fmt.Sprintf(`
resource "panos_dag_tags" "test" {
    vsys = "vsys1"
    register {
        ip = "%s"
        tags = ["%s", "%s"]
    }
    register {
        ip = "%s"
        tags = ["%s"]
    }
    register {
        ip = "%s"
        tags = ["%s", "%s", "%s"]
    }
}
`, ip1, ip1t1, ip1t2, ip2, ip2t1, ip3, ip3t1, ip3t2, ip3t3)
}
