package panos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPanosDsPredefinedThreats(t *testing.T) {
	if len(testAccPredefinedPhoneHomeThreats) == 0 {
		t.Skip("No predefined phone home threats found")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPredefinedThreatConfig(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "total"),
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "threats.0.name"),
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "threats.0.threat_name"),
				),
			},
		},
	})
}

func TestAccPanosDsPredefinedThreat(t *testing.T) {
	if len(testAccPredefinedPhoneHomeThreats) == 0 {
		t.Skip("No predefined phone home threats found")
	}

	name := testAccPredefinedPhoneHomeThreats[acctest.RandInt()%len(testAccPredefinedPhoneHomeThreats)].Name

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPredefinedThreatConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "total"),
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "threats.0.name"),
					resource.TestCheckResourceAttrSet("data.panos_predefined_threat.test", "threats.0.threat_name"),
				),
			},
		},
	})
}

func testAccPredefinedThreatConfig(name string) string {
	switch name {
	case "":
		return `
data "panos_predefined_threat" "test" {}
`
	default:
		return fmt.Sprintf(`
data "panos_predefined_threat" "test" {
    name = %q
}
`, name)
	}
}
