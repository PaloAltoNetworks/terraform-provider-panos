package panos

import (
    "os"
    "testing"

    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/terraform"
)


var testAccProviders map[string] terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
    testAccProvider = Provider().(*schema.Provider)
    testAccProviders = map[string] terraform.ResourceProvider{
        "panos": testAccProvider,
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
