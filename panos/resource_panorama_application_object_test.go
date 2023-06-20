package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/app"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaApplicationObject_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o app.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaApplicationObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaApplicationObjectConfig(dg, name, "first desc", "media", "gaming", "browser-based", `
    defaults {
        port {
            ports = ["udp/dynamic"]
        }
    }`, 2, 100, 200, 300, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationObjectExists("panos_panorama_application_object.test", &o),
					testAccCheckPanosPanoramaApplicationObjectAttributes(&o, name, "first desc", "media", "gaming", "browser-based", app.DefaultTypePort, 2, 100, 200, 300, true, false, true),
				),
			},
			{
				Config: testAccPanoramaApplicationObjectConfig(dg, name, "desc updated", "networking", "proxy", "client-server", `
    defaults {
        ip_protocol {
            value = "21"
        }
    }`, 3, 101, 202, 303, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationObjectExists("panos_panorama_application_object.test", &o),
					testAccCheckPanosPanoramaApplicationObjectAttributes(&o, name, "desc updated", "networking", "proxy", "client-server", app.DefaultTypeIpProtocol, 3, 101, 202, 303, false, true, false),
				),
			},
			{
				Config: testAccPanoramaApplicationObjectConfig(dg, name, "desc3", "collaboration", "email", "browser-based", `
    defaults {
        icmp {
            type = 7
            code = 11
        }
    }`, 1, 102, 203, 304, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationObjectExists("panos_panorama_application_object.test", &o),
					testAccCheckPanosPanoramaApplicationObjectAttributes(&o, name, "desc3", "collaboration", "email", "browser-based", app.DefaultTypeIcmp, 1, 102, 203, 304, false, false, true),
				),
			},
			{
				Config: testAccPanoramaApplicationObjectConfig(dg, name, "desc4", "collaboration", "instant-messaging", "peer-to-peer", "", 4, 103, 204, 305, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationObjectExists("panos_panorama_application_object.test", &o),
					testAccCheckPanosPanoramaApplicationObjectAttributes(&o, name, "desc4", "collaboration", "instant-messaging", "peer-to-peer", app.DefaultTypeNone, 4, 103, 204, 305, false, true, false),
				),
			},
			{
				Config: testAccPanoramaApplicationObjectConfig(dg, name, "desc final", "business-systems", "database", "client-server", `
    defaults {
        icmp6 {
            type = 8
            code = 12
        }
    }`, 5, 105, 206, 307, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationObjectExists("panos_panorama_application_object.test", &o),
					testAccCheckPanosPanoramaApplicationObjectAttributes(&o, name, "desc final", "business-systems", "database", "client-server", app.DefaultTypeIcmp6, 5, 105, 206, 307, true, false, false),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaApplicationObjectExists(n string, o *app.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, _, name := parseApplicationObjectId(rs.Primary.ID)
		v, err := pano.Objects.Application.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaApplicationObjectAttributes(o *app.Entry, name, desc, cat, subc, tech, def string, risk, tout, tcpt, udpt int, fti, vi, dpi bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.Category != cat {
			return fmt.Errorf("Category is %s, expected %s", o.Category, cat)
		}

		if o.Subcategory != subc {
			return fmt.Errorf("Subcategory is %s, expected %s", o.Subcategory, subc)
		}

		if o.Technology != tech {
			return fmt.Errorf("Technology is %s, expected %s", o.Technology, tech)
		}

		if o.DefaultType != def {
			return fmt.Errorf("Default type is %s, expected %s", o.DefaultType, def)
		}

		switch def {
		case app.DefaultTypePort:
			if len(o.DefaultPorts) != 1 || o.DefaultPorts[0] != "udp/dynamic" {
				return fmt.Errorf("Default ports is %#v, expected [udp/dynamic]", o.DefaultPorts)
			}
		case app.DefaultTypeIpProtocol:
			if o.DefaultIpProtocol != "21" {
				return fmt.Errorf("Default IP protocol is %s, expected '21'", o.DefaultIpProtocol)
			}
		case app.DefaultTypeIcmp:
			if o.DefaultIcmpType != 7 {
				return fmt.Errorf("Default icmp type is %d, not 7", o.DefaultIcmpType)
			} else if o.DefaultIcmpCode != 11 {
				return fmt.Errorf("Default ICMP code is %d, expected 11", o.DefaultIcmpCode)
			}
		case app.DefaultTypeIcmp6:
			if o.DefaultIcmpType != 8 {
				return fmt.Errorf("Default icmp type is %d, not 8", o.DefaultIcmpType)
			} else if o.DefaultIcmpCode != 12 {
				return fmt.Errorf("Default ICMP code is %d, expected 12", o.DefaultIcmpCode)
			}
		}

		if o.Risk != risk {
			return fmt.Errorf("Risk is %d, expected %d", o.Risk, risk)
		}

		if o.Timeout != tout {
			return fmt.Errorf("Timeout is %d, expected %d", o.Timeout, tout)
		}

		if o.TcpTimeout != tcpt {
			return fmt.Errorf("TCP timeout is %d, expected %d", o.TcpTimeout, tcpt)
		}

		if o.UdpTimeout != udpt {
			return fmt.Errorf("UDP timeout is %d, expected %d", o.UdpTimeout, udpt)
		}

		if o.FileTypeIdent != fti {
			return fmt.Errorf("File type ident is %t, expected %t", o.FileTypeIdent, fti)
		}

		if o.VirusIdent != vi {
			return fmt.Errorf("Virus ident is %t, expected %t", o.VirusIdent, vi)
		}

		if o.DataIdent != dpi {
			return fmt.Errorf("Data ident is %t, expected %t", o.DataIdent, dpi)
		}

		return nil
	}
}

func testAccPanosPanoramaApplicationObjectDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_application_object" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, _, name := parseApplicationObjectId(rs.Primary.ID)
			_, err := pano.Objects.Application.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaApplicationObjectConfig(dg, name, desc, cat, subc, tech, def string, risk, tout, tcpt, udpt int, fti, vi, dpi bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
    description = "application object"
}

resource "panos_panorama_application_object" "test" {
    device_group = panos_panorama_device_group.x.name
    name = %q
    description = %q
    category = %q
    subcategory = %q
    technology = %q
%s
    risk = %d
    timeout_settings {
        timeout = %d
        tcp_timeout = %d
        udp_timeout = %d
    }
    scanning {
        file_types = %t
        viruses = %t
        data_patterns = %t
    }
}
`, dg, name, desc, cat, subc, tech, def, risk, tout, tcpt, udpt, fti, vi, dpi)
}
