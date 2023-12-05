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

func TestAccPanosIkeCryptoProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var mp ike.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIkeCryptoProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIkeCryptoProfileConfig(name, "group1", "md5", ike.Encryption3des, ike.TimeMinutes, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIkeCryptoProfileExists("panos_ike_crypto_profile.test", &mp),
					testAccCheckPanosIkeCryptoProfileAttributes(&mp, name, "group1", "md5", ike.Encryption3des, ike.TimeMinutes, 5),
				),
			},
			{
				Config: testAccIkeCryptoProfileConfig(name, "group5", "sha1", ike.EncryptionAes128, ike.TimeHours, 6),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIkeCryptoProfileExists("panos_ike_crypto_profile.test", &mp),
					testAccCheckPanosIkeCryptoProfileAttributes(&mp, name, "group5", "sha1", ike.EncryptionAes128, ike.TimeHours, 6),
				),
			},
			{
				Config: testAccIkeCryptoProfileConfig(name, "group14", "sha256", ike.EncryptionAes192, ike.TimeDays, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIkeCryptoProfileExists("panos_ike_crypto_profile.test", &mp),
					testAccCheckPanosIkeCryptoProfileAttributes(&mp, name, "group14", "sha256", ike.EncryptionAes192, ike.TimeDays, 7),
				),
			},
		},
	})
}

func testAccCheckPanosIkeCryptoProfileExists(n string, o *ike.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, _, name := parseIkeCryptoProfileId(rs.Primary.ID)
		v, err := fw.Network.IkeCryptoProfile.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIkeCryptoProfileAttributes(o *ike.Entry, name, dhg, auth, enc, ltt string, ltv int) resource.TestCheckFunc {
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

func testAccPanosIkeCryptoProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ike_crypto_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			_, _, name := parseIkeCryptoProfileId(rs.Primary.ID)
			_, err := fw.Network.IkeCryptoProfile.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", name)
			}
		}
		return nil
	}

	return nil
}

func testAccIkeCryptoProfileConfig(name, dhg, auth, enc, ltt string, ltv int) string {
	return fmt.Sprintf(`
resource "panos_ike_crypto_profile" "test" {
    name = %q
    dh_groups = [%q]
    authentications = [%q]
    encryptions = [%q]
    lifetime_type = %q
    lifetime_value = %d
}
`, name, dhg, auth, enc, ltt, ltv)
}
