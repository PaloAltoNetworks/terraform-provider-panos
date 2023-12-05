package panos

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func checkDataSourceListing(x string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			fmt.Sprintf("data.%s.test", x),
			"total",
		),
	)
}

func checkDataSource(x string, keys []string) resource.TestCheckFunc {
	list := make([]resource.TestCheckFunc, 0, len(keys))

	for _, key := range keys {
		list = append(list, resource.TestCheckResourceAttrSet(
			fmt.Sprintf("data.%s.test", x),
			key,
		))
	}

	return resource.ComposeTestCheckFunc(list...)
}
