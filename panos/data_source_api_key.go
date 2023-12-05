package panos

import (
	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiKeyRead,

		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The API key associated with the provider's authentication credentials",
			},
		},
	}
}

func dataSourceApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	switch c := meta.(type) {
	case *pango.Firewall:
		d.SetId(base64Encode([]string{
			c.Hostname, c.Username,
		}))
		d.Set("api_key", c.ApiKey)
	case *pango.Panorama:
		d.SetId(base64Encode([]string{
			c.Hostname, c.Username,
		}))
		d.Set("api_key", c.ApiKey)
	}

	return nil
}
