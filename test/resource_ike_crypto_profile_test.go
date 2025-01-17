package provider_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/profiles/ikecrypto"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccIkeCryptoProfile_1(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccIkeCryptoProfileDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: ikeCryptoProfile1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("authentication_multiple"),
						knownvalue.Int64Exact(50),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("dh_group"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("group1"),
							knownvalue.StringExact("group2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("encryption"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("3des"),
							knownvalue.StringExact("aes-256-gcm"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("hash"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("md5"),
							knownvalue.StringExact("sha256"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile1",
						tfjsonpath.New("lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"days":    knownvalue.Null(),
							"hours":   knownvalue.Null(),
							"minutes": knownvalue.Null(),
							"seconds": knownvalue.Int64Exact(3600),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile2",
						tfjsonpath.New("lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"days":    knownvalue.Int64Exact(7),
							"hours":   knownvalue.Null(),
							"minutes": knownvalue.Null(),
							"seconds": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile3",
						tfjsonpath.New("lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"days":    knownvalue.Null(),
							"hours":   knownvalue.Int64Exact(20),
							"minutes": knownvalue.Null(),
							"seconds": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_crypto_profile.profile4",
						tfjsonpath.New("lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"days":    knownvalue.Null(),
							"hours":   knownvalue.Null(),
							"minutes": knownvalue.Int64Exact(15),
							"seconds": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const ikeCryptoProfile1 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }
  name     = local.template_name
}


resource "panos_ike_crypto_profile" "profile1" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-profile1", var.prefix)

  authentication_multiple = 50
  dh_group = ["group1", "group2"]
  encryption = ["3des", "aes-256-gcm"]
  hash = ["md5", "sha256"]
  lifetime = {
    seconds = 3600
  }
}

resource "panos_ike_crypto_profile" "profile2" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-profile2", var.prefix)

  lifetime = {
    days = 7
  }
}

resource "panos_ike_crypto_profile" "profile3" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-profile3", var.prefix)

  lifetime = {
    hours = 20
  }
}

resource "panos_ike_crypto_profile" "profile4" {
  location = { template = { name = panos_template.template.name } }

  name    = format("%s-profile4", var.prefix)

  lifetime = {
    minutes = 15
  }
}
`

func testAccIkeCryptoProfileDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		entry := fmt.Sprintf("%s-profile", prefix)
		api := ikecrypto.NewService(sdkClient)
		ctx := context.TODO()

		location := ikecrypto.NewTemplateLocation()
		location.Template.Template = fmt.Sprintf("%s-tmpl", prefix)

		reply, err := api.Read(ctx, *location, entry, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading ethernet entry via sdk: %v", err)
		}

		if reply != nil {
			err := fmt.Errorf("terraform didn't delete the server entry properly")
			delErr := api.Delete(ctx, *location, entry)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
			return err
		}

		return nil
	}
}
