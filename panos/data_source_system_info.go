package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSystemInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSystemInfoRead,

		Schema: map[string]*schema.Schema{
			"version_major": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_minor": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_patch": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"info": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSystemInfoRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*pango.Firewall)

	d.SetId(c.Hostname)
	d.Set("version_major", c.Version.Major)
	d.Set("version_minor", c.Version.Minor)
	d.Set("version_patch", c.Version.Patch)
	if err := d.Set("info", c.SystemInfo); err != nil {
		log.Printf("[WARN] Error setting 'info' field for %q: %s", d.Id(), err)
	}

	return nil
}
