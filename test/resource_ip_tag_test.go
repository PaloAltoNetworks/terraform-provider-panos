package provider_test

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango/xmlapi"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// The IPs we register tags against. 192.0.2.0/24 is TEST-NET-1 (RFC 5737),
// reserved for documentation/testing, so it is safe to use on a live device.
// Each test uses a distinct address so the parallel suites don't collide: the
// data source reports every tag present on an IP, so a shared address would let
// one test observe another's registrations.
const (
	testAccIpTagAddress           = "192.0.2.1"
	testAccIpTagDataSourceAddress = "192.0.2.2"
)

func testAccIpTagPreCheck(t *testing.T) {
	if os.Getenv("PANOS_HOSTNAME") == "" {
		t.Fatal("PANOS_HOSTNAME must be set for acceptance tests")
	}
}

// TestAccIpTag exercises the panorama location: tags registered directly on
// Panorama's own User-ID table (no target firewall).
func TestAccIpTag(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	webTag := fmt.Sprintf("test-acc-%s-web", suffix)
	dbTag := fmt.Sprintf("test-acc-%s-db", suffix)
	prodTag := fmt.Sprintf("test-acc-%s-prod", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccIpTagPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccIpTagCheckDestroy(testAccIpTagAddress, []string{webTag, dbTag, prodTag}),
		Steps: []resource.TestStep{
			{
				// Create: register {web, db}.
				Config: testAccIpTagTmpl,
				ConfigVariables: config.Variables{
					"ip": config.StringVariable(testAccIpTagAddress),
					"tags": config.SetVariable(
						config.StringVariable(webTag),
						config.StringVariable(dbTag),
					),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ip_tag.test",
						tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(webTag),
							knownvalue.StringExact(dbTag),
						}),
					),
				},
			},
			{
				// Update: {web, db} -> {web, prod}. Registers prod, unregisters db.
				Config: testAccIpTagTmpl,
				ConfigVariables: config.Variables{
					"ip": config.StringVariable(testAccIpTagAddress),
					"tags": config.SetVariable(
						config.StringVariable(webTag),
						config.StringVariable(prodTag),
					),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ip_tag.test",
						tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(webTag),
							knownvalue.StringExact(prodTag),
						}),
					),
				},
			},
		},
	})
}

const testAccIpTagTmpl = `
variable "ip" { type = string }
variable "tags" { type = set(string) }

resource "panos_ip_tag" "test" {
  location = { panorama = {} }

  ip   = var.ip
  tags = var.tags
}
`

// TestAccIpTagDataSource registers tags via the resource and reads them back
// through the data source, asserting the data source reports every tag present
// on the IP.
func TestAccIpTagDataSource(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	webTag := fmt.Sprintf("test-acc-%s-web", suffix)
	dbTag := fmt.Sprintf("test-acc-%s-db", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccIpTagPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccIpTagCheckDestroy(testAccIpTagDataSourceAddress, []string{webTag, dbTag}),
		Steps: []resource.TestStep{
			{
				Config: testAccIpTagDataSourceTmpl,
				ConfigVariables: config.Variables{
					"ip": config.StringVariable(testAccIpTagDataSourceAddress),
					"tags": config.SetVariable(
						config.StringVariable(webTag),
						config.StringVariable(dbTag),
					),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.panos_ip_tag.test",
						tfjsonpath.New("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(webTag),
							knownvalue.StringExact(dbTag),
						}),
					),
				},
			},
		},
	})
}

const testAccIpTagDataSourceTmpl = `
variable "ip" { type = string }
variable "tags" { type = set(string) }

resource "panos_ip_tag" "test" {
  location = { panorama = {} }

  ip   = var.ip
  tags = var.tags
}

data "panos_ip_tag" "test" {
  location = { panorama = {} }

  ip = panos_ip_tag.test.ip
}
`

// testAccIpTagCheckDestroy queries Panorama's registered-ip table and fails if
// any of the test's managed tags are still present on the given IP.
func testAccIpTagCheckDestroy(ip string, managedTags []string) func(*terraform.State) error {
	return func(_ *terraform.State) error {
		type respEntry struct {
			Ip   string   `xml:"ip,attr"`
			Tags []string `xml:"tag>member"`
		}
		type showFilter struct {
			Ip    string `xml:"ip,omitempty"`
			Limit int    `xml:"limit"`
			Start int    `xml:"start-point"`
		}
		type showReq struct {
			XMLName xml.Name   `xml:"show"`
			Filter  showFilter `xml:"object>registered-ip"`
		}
		type showResp struct {
			Entries []respEntry `xml:"result>entry"`
		}

		cmd := &xmlapi.Op{
			Command: showReq{Filter: showFilter{Ip: ip, Limit: 500, Start: 1}},
		}

		var resp showResp
		if _, _, err := sdkClient.Communicate(context.TODO(), cmd, false, &resp); err != nil {
			return fmt.Errorf("failed to query registered-ip table: %w", err)
		}

		managed := make(map[string]struct{}, len(managedTags))
		for _, tag := range managedTags {
			managed[tag] = struct{}{}
		}

		for _, entry := range resp.Entries {
			if entry.Ip != ip {
				continue
			}
			for _, tag := range entry.Tags {
				if _, ok := managed[tag]; ok {
					return fmt.Errorf("tag %q is still registered on %s after destroy", tag, ip)
				}
			}
		}

		return nil
	}
}
