package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/profile/ipsec"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosIpsecCryptoProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var mp ipsec.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIpsecCryptoProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecCryptoProfileConfig(name, "group1", "md5", ipsec.Encryption3des, ipsec.TimeMinutes, 5, ipsec.SizeKb, 15),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecCryptoProfileExists("panos_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosIpsecCryptoProfileAttributes(&mp, name, "group1", "md5", ipsec.Encryption3des, ipsec.TimeMinutes, 5, ipsec.SizeKb, 15),
				),
			},
			{
				Config: testAccIpsecCryptoProfileConfig(name, "group5", "sha1", ipsec.EncryptionAes128, ipsec.TimeHours, 6, ipsec.SizeMb, 16),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecCryptoProfileExists("panos_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosIpsecCryptoProfileAttributes(&mp, name, "group5", "sha1", ipsec.EncryptionAes128, ipsec.TimeHours, 6, ipsec.SizeMb, 16),
				),
			},
			{
				Config: testAccIpsecCryptoProfileConfig(name, "group14", "sha256", ipsec.EncryptionAes192, ipsec.TimeDays, 7, ipsec.SizeGb, 17),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecCryptoProfileExists("panos_ipsec_crypto_profile.test", &mp),
					testAccCheckPanosIpsecCryptoProfileAttributes(&mp, name, "group14", "sha256", ipsec.EncryptionAes192, ipsec.TimeDays, 7, ipsec.SizeGb, 17),
				),
			},
		},
	})
}

func testAccCheckPanosIpsecCryptoProfileExists(n string, o *ipsec.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		name := rs.Primary.ID
		v, err := fw.Network.IpsecCryptoProfile.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIpsecCryptoProfileAttributes(o *ipsec.Entry, name, dhg, auth, enc, ltt string, ltv int, lst string, lsv int) resource.TestCheckFunc {
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

func testAccPanosIpsecCryptoProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ipsec_crypto_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.IpsecCryptoProfile.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccIpsecCryptoProfileConfig(name, dhg, auth, enc, ltt string, ltv int, lst string, lsv int) string {
	return fmt.Sprintf(`
resource "panos_ipsec_crypto_profile" "test" {
    name = %q
    dh_group = %q
    authentications = [%q]
    encryptions = [%q]
    lifetime_type = %q
    lifetime_value = %d
    lifesize_type = %q
    lifesize_value = %d
}
`, name, dhg, auth, enc, ltt, ltv, lst, lsv)
}
