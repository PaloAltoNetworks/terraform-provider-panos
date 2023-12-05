package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPanosDsPredefinedDlpFileType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPredefinedDlpFileTypeConfig(),
				Check: checkDataSource("panos_predefined_dlp_file_type", []string{
					"total",
					"file_types.0.name",
					"file_types.0.properties.0.name",
					"file_types.0.properties.0.label",
					"file_types.1.name",
					"file_types.1.properties.0.name",
					"file_types.1.properties.0.label",
					"file_types.2.name",
					"file_types.2.properties.0.name",
					"file_types.2.properties.0.label",
				}),
			},
		},
	})
}

func testAccPredefinedDlpFileTypeConfig() string {
	return `
data "panos_predefined_dlp_file_type" "test" {}
`
}
