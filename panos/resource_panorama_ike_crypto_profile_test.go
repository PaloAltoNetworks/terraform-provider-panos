package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ike"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaIkeCryptoProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var mp ike.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	ts := fmt.Sprintf("tfStack%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaIkeCryptoProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaIkeCryptoProfileConfig(ts, name, "group1", "md5", ike.Encryption3des, ike.TimeMinutes, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIkeCryptoProfileExists("panos_panorama_ike_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIkeCryptoProfileAttributes(&mp, name, "group1", "md5", ike.Encryption3des, ike.TimeMinutes, 5),
				),
			},
			{
				Config: testAccPanoramaIkeCryptoProfileConfig(ts, name, "group5", "sha1", ike.EncryptionAes128, ike.TimeHours, 6),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIkeCryptoProfileExists("panos_panorama_ike_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIkeCryptoProfileAttributes(&mp, name, "group5", "sha1", ike.EncryptionAes128, ike.TimeHours, 6),
				),
			},
			{
				Config: testAccPanoramaIkeCryptoProfileConfig(ts, name, "group14", "sha256", ike.EncryptionAes192, ike.TimeDays, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIkeCryptoProfileExists("panos_panorama_ike_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIkeCryptoProfileAttributes(&mp, name, "group14", "sha256", ike.EncryptionAes192, ike.TimeDays, 7),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaIkeCryptoProfileExists(n string, o *ike.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaIkeCryptoProfileId(rs.Primary.ID)
		v, err := pano.Network.IkeCryptoProfile.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaIkeCryptoProfileAttributes(o *ike.Entry, name, dhg, auth, enc, ltt string, ltv int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.DhGroup) != 1 || o.DhGroup[0] != dhg {
			return fmt.Errorf("Dh group is %#v, expected [%s]", o.DhGroup, dhg)
		}

		if len(o.Authentication) != 1 || o.Authentication[0] != auth {
			return fmt.Errorf("Auth is %#v, expected [%s]", o.Authentication, auth)
		}

		if len(o.Encryption) != 1 || o.Encryption[0] != enc {
			return fmt.Errorf("Encryption is %#v, expected [%s]", o.Encryption, enc)
		}

		if o.LifetimeType != ltt {
			return fmt.Errorf("Lifetime type is %q, expected %q", o.LifetimeType, ltt)
		}

		if o.LifetimeValue != ltv {
			return fmt.Errorf("Lifetime value is %d, expected %d", o.LifetimeValue, ltv)
		}

		return nil
	}
}

func testAccPanosPanoramaIkeCryptoProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_ike_crypto_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaIkeCryptoProfileId(rs.Primary.ID)
			_, err := pano.Network.IkeCryptoProfile.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", name)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaIkeCryptoProfileConfig(ts, name, dhg, auth, enc, ltt string, ltv int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
}

resource "panos_panorama_ike_crypto_profile" "test" {
    name = %q
    template_stack = panos_panorama_template_stack.x.name
    dh_groups = [%q]
    authentications = [%q]
    encryptions = [%q]
    lifetime_type = %q
    lifetime_value = %d
}
`, ts, name, dhg, auth, enc, ltt, ltv)
}
