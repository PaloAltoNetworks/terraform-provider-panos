package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPanosDsPredefinedTdbFileType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPredefinedTdbFileTypeConfig(),
				Check: checkDataSource("panos_predefined_tdb_file_type", []string{
					"total",
					"file_types.0.data_ident",
					"file_types.0.file_type_id",
					"file_types.0.name",
					"file_types.0.full_name",
				}),
			},
		},
	})
}

func testAccPredefinedTdbFileTypeConfig() string {
	return `
data "panos_predefined_tdb_file_type" "test" {
    data_ident_only = true
}
`
}
