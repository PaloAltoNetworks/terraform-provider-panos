package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// The IPSec tunnel manual_key AH/ESP keys are hex key material. These tests
// verify that each key leaf behaves as a hashing.type: solo field: the
// configured plaintext round-trips exactly in state, and re-applying the
// identical config produces no diff. A key that is only sensitive (not solo)
// would fail the ExpectEmptyPlan step because PAN-OS returns a hashed/encrypted
// value on read, causing a perpetual diff.
//
// Hex key lengths follow the PAN-OS manual-key requirements per algorithm
// (md5=32, sha1=40, sha256=64, sha384=96, sha512=128 hex digits; aes-256-cbc=64).
// The local_spi/remote_spi values must be within the PAN-OS hex SPI range
// (00001000..1FFFFFFF). Adjust these values if a specific PAN-OS build rejects
// them.
const (
	ipsecKeyMd5    = "00112233445566778899aabbccddeeff"
	ipsecKeySha1   = "00112233445566778899aabbccddeeff00112233"
	ipsecKeySha256 = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	ipsecKeySha384 = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	ipsecKeySha512 = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	ipsecKeyEnc256 = "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
)

func TestAccPanosIpsecTunnel_ManualKey(t *testing.T) {
	t.Parallel()

	base := tfjsonpath.New("manual_key")

	cases := []struct {
		name          string
		manualKeyAttr string
		keyPath       tfjsonpath.Path
		keyValue      string
	}{
		{
			name:          "Ah_Md5",
			manualKeyAttr: `ah = { md5 = { key = "` + ipsecKeyMd5 + `" } }`,
			keyPath:       base.AtMapKey("ah").AtMapKey("md5").AtMapKey("key"),
			keyValue:      ipsecKeyMd5,
		},
		{
			name:          "Ah_Sha1",
			manualKeyAttr: `ah = { sha1 = { key = "` + ipsecKeySha1 + `" } }`,
			keyPath:       base.AtMapKey("ah").AtMapKey("sha1").AtMapKey("key"),
			keyValue:      ipsecKeySha1,
		},
		{
			name:          "Ah_Sha256",
			manualKeyAttr: `ah = { sha256 = { key = "` + ipsecKeySha256 + `" } }`,
			keyPath:       base.AtMapKey("ah").AtMapKey("sha256").AtMapKey("key"),
			keyValue:      ipsecKeySha256,
		},
		{
			name:          "Ah_Sha384",
			manualKeyAttr: `ah = { sha384 = { key = "` + ipsecKeySha384 + `" } }`,
			keyPath:       base.AtMapKey("ah").AtMapKey("sha384").AtMapKey("key"),
			keyValue:      ipsecKeySha384,
		},
		{
			name:          "Ah_Sha512",
			manualKeyAttr: `ah = { sha512 = { key = "` + ipsecKeySha512 + `" } }`,
			keyPath:       base.AtMapKey("ah").AtMapKey("sha512").AtMapKey("key"),
			keyValue:      ipsecKeySha512,
		},
		{
			name: "Esp_Auth_Md5",
			manualKeyAttr: `esp = {
      authentication = { md5 = { key = "` + ipsecKeyMd5 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("authentication").AtMapKey("md5").AtMapKey("key"),
			keyValue: ipsecKeyMd5,
		},
		{
			name: "Esp_Auth_Sha1",
			manualKeyAttr: `esp = {
      authentication = { sha1 = { key = "` + ipsecKeySha1 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("authentication").AtMapKey("sha1").AtMapKey("key"),
			keyValue: ipsecKeySha1,
		},
		{
			name: "Esp_Auth_Sha256",
			manualKeyAttr: `esp = {
      authentication = { sha256 = { key = "` + ipsecKeySha256 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("authentication").AtMapKey("sha256").AtMapKey("key"),
			keyValue: ipsecKeySha256,
		},
		{
			name: "Esp_Auth_Sha384",
			manualKeyAttr: `esp = {
      authentication = { sha384 = { key = "` + ipsecKeySha384 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("authentication").AtMapKey("sha384").AtMapKey("key"),
			keyValue: ipsecKeySha384,
		},
		{
			name: "Esp_Auth_Sha512",
			manualKeyAttr: `esp = {
      authentication = { sha512 = { key = "` + ipsecKeySha512 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("authentication").AtMapKey("sha512").AtMapKey("key"),
			keyValue: ipsecKeySha512,
		},
		{
			name: "Esp_Encryption",
			manualKeyAttr: `esp = {
      authentication = { md5 = { key = "` + ipsecKeyMd5 + `" } }
      encryption     = { algorithm = "aes-256-cbc", key = "` + ipsecKeyEnc256 + `" }
    }`,
			keyPath:  base.AtMapKey("esp").AtMapKey("encryption").AtMapKey("key"),
			keyValue: ipsecKeyEnc256,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
			prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
			cfg := testAccIpsecTunnelManualKeyConfig(tc.manualKeyAttr)
			vars := map[string]config.Variable{"prefix": config.StringVariable(prefix)}

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:          cfg,
						ConfigVariables: vars,
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(
								"panos_ipsec_tunnel.test",
								tc.keyPath,
								knownvalue.StringExact(tc.keyValue),
							),
						},
					},
					{
						// Re-applying the identical config must produce no diff.
						Config:          cfg,
						ConfigVariables: vars,
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectEmptyPlan(),
							},
						},
					},
				},
			})
		})
	}
}

func testAccIpsecTunnelManualKeyConfig(manualKeyAttr string) string {
	return fmt.Sprintf(`
variable "prefix" { type = string }

resource "panos_template" "test" {
  location = { panorama = {} }
  name     = var.prefix
}

resource "panos_tunnel_interface" "test" {
  location = { template = { name = panos_template.test.name } }
  name     = "tunnel.1"
}

# The manual_key local_address.ip must reference an IP actually configured on
# the terminating interface, so create an L3 interface that owns it.
resource "panos_ethernet_interface" "test" {
  location = { template = { vsys = "vsys1", name = panos_template.test.name } }
  name     = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
  }
}

resource "panos_ipsec_tunnel" "test" {
  location = { template = { name = panos_template.test.name } }

  name             = var.prefix
  tunnel_interface = panos_tunnel_interface.test.name

  manual_key = {
    local_spi  = "00001000"
    remote_spi = "00001001"
    # local_address requires exactly one of ip / floating_ip; ip references an
    # address configured on the interface.
    local_address = {
      interface = panos_ethernet_interface.test.name
      ip        = "10.0.0.1/24"
    }
    peer_address = { ip = "10.0.0.2" }
    %s
  }
}
`, manualKeyAttr)
}
