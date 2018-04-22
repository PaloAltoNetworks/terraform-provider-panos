package panos

import (
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	testAccProviders                     map[string]terraform.ResourceProvider
	testAccProvider                      *schema.Provider
	testAccIsFirewall, testAccIsPanorama bool
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"panos": testAccProvider,
	}

	/*
	   We need to know whether we are running acctests for the firewall or the
	   Panorama.  Since the provider meta is nil until Configure is called on
	   it, just check here in the init function.
	*/
	con, err := pango.Connect(pango.Client{
		Hostname: os.Getenv("PANOS_HOSTNAME"),
		Username: os.Getenv("PANOS_USERNAME"),
		Password: os.Getenv("PANOS_PASSWORD"),
		Logging:  pango.LogQuiet,
	})
	if err == nil {
		switch con.(type) {
		case *pango.Firewall:
			testAccIsFirewall = true
		case *pango.Panorama:
			testAccIsPanorama = true
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PANOS_HOSTNAME") == "" {
		t.Fatal("PANOS_HOSTNAME must be set for acceptance tests")
	}
	if os.Getenv("PANOS_USERNAME") == "" {
		t.Fatal("PANOS_USERNAME must be set for acceptance tests")
	}
	if os.Getenv("PANOS_PASSWORD") == "" {
		t.Fatal("PANOS_PASSWORD must be set for acceptance tests")
	}
}
