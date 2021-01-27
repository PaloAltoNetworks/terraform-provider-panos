package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/profile/auth"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test (listing).
func TestAccPanosDsOspfAuthProfileList(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAuthProfileConfig(tmpl, vr, name),
				Check:  checkDataSourceListing("panos_ospf_auth_profiles"),
			},
		},
	})
}

func testAccDsOspfAuthProfileConfig(tmpl, vr, name string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf_auth_profiles" "test" {
    template = panos_ospf_auth_profile.x.template
    virtual_router = panos_ospf_auth_profile.x.virtual_router
}

resource "panos_panorama_template" "x" {
    name = %q
    description = "for ospf auth profile data source acctest"
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = %q
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

resource "panos_ospf_auth_profile" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    password = "secret"
}
`, tmpl, vr, name)
	}

	return fmt.Sprintf(`
data "panos_ospf_auth_profiles" "test" {
    virtual_router = panos_ospf_auth_profile.x.virtual_router
}

resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_auth_profile" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    password = "secret"
}
`, vr, name)
}

// Resource tests.
func TestAccPanosOspfAuthProfile(t *testing.T) {
	var o auth.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	md5s := []auth.Md5Key{
		{
			KeyId:     1,
			Key:       "alpha key",
			Preferred: true,
		},
		{
			KeyId:     2,
			Key:       "key two",
			Preferred: false,
		},
		{
			KeyId:     3,
			Key:       "key trois",
			Preferred: true,
		},
		{
			KeyId:     4,
			Key:       "key 4",
			Preferred: false,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfAuthProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfAuthProfileConfig(tmpl, vr, name, auth.AuthTypePassword, "secret1", nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAuthProfileExists("panos_ospf_auth_profile.test", &o),
					testAccCheckPanosOspfAuthProfileAttributes(&o, name, auth.AuthTypePassword, "secret1", nil),
				),
			},
			{
				Config: testAccOspfAuthProfileConfig(tmpl, vr, name, auth.AuthTypePassword, "second", nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAuthProfileExists("panos_ospf_auth_profile.test", &o),
					testAccCheckPanosOspfAuthProfileAttributes(&o, name, auth.AuthTypePassword, "second", nil),
				),
			},
			{
				Config: testAccOspfAuthProfileConfig(tmpl, vr, name, auth.AuthTypeMd5, "", []auth.Md5Key{md5s[0], md5s[1]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAuthProfileExists("panos_ospf_auth_profile.test", &o),
					testAccCheckPanosOspfAuthProfileAttributes(&o, name, auth.AuthTypeMd5, "", []auth.Md5Key{md5s[0], md5s[1]}),
				),
			},
			{
				Config: testAccOspfAuthProfileConfig(tmpl, vr, name, auth.AuthTypeMd5, "", []auth.Md5Key{md5s[0], md5s[2]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAuthProfileExists("panos_ospf_auth_profile.test", &o),
					testAccCheckPanosOspfAuthProfileAttributes(&o, name, auth.AuthTypeMd5, "", []auth.Md5Key{md5s[0], md5s[2]}),
				),
			},
			{
				Config: testAccOspfAuthProfileConfig(tmpl, vr, name, auth.AuthTypeMd5, "", md5s),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAuthProfileExists("panos_ospf_auth_profile.test", &o),
					testAccCheckPanosOspfAuthProfileAttributes(&o, name, auth.AuthTypeMd5, "", md5s),
				),
			},
		},
	})
}

func testAccCheckPanosOspfAuthProfileExists(n string, o *auth.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v auth.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr, name := parseFirewallOspfAuthProfileId(rs.Primary.ID)
			v, err = con.Network.OspfAuthProfile.Get(vr, name)
		case *pango.Panorama:
			tmpl, ts, vr, name := parsePanoramaOspfAuthProfileId(rs.Primary.ID)
			v, err = con.Network.OspfAuthProfile.Get(tmpl, ts, vr, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfAuthProfileAttributes(o *auth.Entry, name, authType, pwd string, md5s []auth.Md5Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.AuthType != authType {
			return fmt.Errorf("Auth type is %q, not %q", o.AuthType, authType)
		}

		if len(o.Md5Keys) != len(md5s) {
			return fmt.Errorf("md5s is len %d, not %d", len(o.Md5Keys), len(md5s))
		}

		for i := range o.Md5Keys {
			if o.Md5Keys[i].KeyId != md5s[i].KeyId {
				return fmt.Errorf("md5s[%d] key id is %d, not %d", i, o.Md5Keys[i].KeyId, md5s[i].KeyId)
			}

			if o.Md5Keys[i].Preferred != md5s[i].Preferred {
				return fmt.Errorf("md5s[%d] prerferred is %t, not %t", i, o.Md5Keys[i].Preferred, md5s[i].Preferred)
			}
		}

		return nil
	}
}

func testAccPanosOspfAuthProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf_auth_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr, name := parseFirewallOspfAuthProfileId(rs.Primary.ID)
				_, err = con.Network.OspfAuthProfile.Get(vr, name)
			case *pango.Panorama:
				tmpl, ts, vr, name := parsePanoramaOspfAuthProfileId(rs.Primary.ID)
				_, err = con.Network.OspfAuthProfile.Get(tmpl, ts, vr, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfAuthProfileConfig(tmpl, vr, name, authType, pwd string, md5s []auth.Md5Key) string {
	var b strings.Builder
	b.WriteString(`
resource "panos_ospf_auth_profile" "test" {`)
	if testAccIsPanorama {
		b.WriteString(`
    template = panos_ospf.x.template`)
	}
	fmt.Fprintf(&b, `
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    auth_type = %q`, name, authType)
	switch authType {
	case auth.AuthTypePassword:
		fmt.Fprintf(&b, `
    password = %q`, pwd)
	case auth.AuthTypeMd5:
		for _, x := range md5s {
			fmt.Fprintf(&b, `
    md5_key {
        key_id = %d
        key = %q
        preferred = %t
    }`, x.KeyId, x.Key, x.Preferred)
		}
	}
	b.WriteString(`
}`)

	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "for ospf auth profile resource acctest"
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = %q
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

%s
`, tmpl, vr, b.String())
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

%s
`, vr, b.String())
}
