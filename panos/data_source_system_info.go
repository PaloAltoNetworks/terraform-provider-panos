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
			"version_major": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_minor": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_patch": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"info": {
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
	switch c := meta.(type) {
	case *pango.Firewall:
		d.SetId(c.Hostname)
		d.Set("version_major", c.Version.Major)
		d.Set("version_minor", c.Version.Minor)
		d.Set("version_patch", c.Version.Patch)
		if err := d.Set("info", c.SystemInfo); err != nil {
			log.Printf("[WARN] Error setting 'info' field for %q: %s", d.Id(), err)
		}
	case *pango.Panorama:
		d.SetId(c.Hostname)
		d.Set("version_major", c.Version.Major)
		d.Set("version_minor", c.Version.Minor)
		d.Set("version_patch", c.Version.Patch)
		if err := d.Set("info", c.SystemInfo); err != nil {
			log.Printf("[WARN] Error setting 'info' field for %q: %s", d.Id(), err)
		}
	}

	return nil
}
