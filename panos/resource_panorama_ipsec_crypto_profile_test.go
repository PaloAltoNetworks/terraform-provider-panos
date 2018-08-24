package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ipsec"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaIpsecCryptoProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var mp ipsec.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	ts := fmt.Sprintf("tfStack%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaIpsecCryptoProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaIpsecCryptoProfileConfig(ts, name, "group1", "md5", ipsec.Encryption3des, ipsec.TimeMinutes, 5, ipsec.SizeKb, 15),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIpsecCryptoProfileExists("panos_panorama_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIpsecCryptoProfileAttributes(&mp, name, "group1", "md5", ipsec.Encryption3des, ipsec.TimeMinutes, 5, ipsec.SizeKb, 15),
				),
			},
			{
				Config: testAccPanoramaIpsecCryptoProfileConfig(ts, name, "group5", "sha1", ipsec.EncryptionAes128, ipsec.TimeHours, 6, ipsec.SizeMb, 16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIpsecCryptoProfileExists("panos_panorama_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIpsecCryptoProfileAttributes(&mp, name, "group5", "sha1", ipsec.EncryptionAes128, ipsec.TimeHours, 6, ipsec.SizeMb, 16),
				),
			},
			{
				Config: testAccPanoramaIpsecCryptoProfileConfig(ts, name, "group14", "sha256", ipsec.EncryptionAes192, ipsec.TimeDays, 7, ipsec.SizeGb, 17),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIpsecCryptoProfileExists("panos_panorama_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosPanoramaIpsecCryptoProfileAttributes(&mp, name, "group14", "sha256", ipsec.EncryptionAes192, ipsec.TimeDays, 7, ipsec.SizeGb, 17),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaIpsecCryptoProfileExists(n string, o *ipsec.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaIpsecCryptoProfileId(rs.Primary.ID)
		v, err := pano.Network.IpsecCryptoProfile.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaIpsecCryptoProfileAttributes(o *ipsec.Entry, name, dhg, auth, enc, ltt string, ltv int, lst string, lsv int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.DhGroup != dhg {
			return fmt.Errorf("Dh group is %q, expected %q", o.DhGroup, dhg)
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

		if o.LifesizeType != lst {
			return fmt.Errorf("Lifesize type is %q, expected %q", o.LifesizeType, lst)
		}

		if o.LifesizeValue != lsv {
			return fmt.Errorf("Lifesize value is %d, expected %d", o.LifesizeValue, lsv)
		}

		return nil
	}
}

func testAccPanosPanoramaIpsecCryptoProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_ipsec_crypto_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaIpsecCryptoProfileId(rs.Primary.ID)
			_, err := pano.Network.IpsecCryptoProfile.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaIpsecCryptoProfileConfig(ts, name, dhg, auth, enc, ltt string, ltv int, lst string, lsv int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
}

resource "panos_panorama_ipsec_crypto_profile" "test" {
    name = %q
    template_stack = "${panos_panorama_template_stack.x.name}"
    dh_group = %q
    authentications = [%q]
    encryptions = [%q]
    lifetime_type = %q
    lifetime_value = %d
    lifesize_type = %q
    lifesize_value = %d
}
`, ts, name, dhg, auth, enc, ltt, ltv, lst, lsv)
}
