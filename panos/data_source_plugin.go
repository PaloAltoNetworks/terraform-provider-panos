package panos

import (
	"log"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/plugin"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginRead,

		Schema: map[string]*schema.Schema{
			"installed": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release_note_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package_file": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"installed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"downloaded": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePluginRead(d *schema.ResourceData, meta interface{}) error {
	var id string
	var list []plugin.Info

	switch v := meta.(type) {
	case *pango.Panorama:
		id = v.Hostname
		list = v.Plugin
	case *pango.Firewall:
		id = v.Hostname
		list = v.Plugin
	}

	details := make([]interface{}, 0, len(list))
	installed := make([]string, 0, len(list))
	for _, pinfo := range list {
		entry := map[string]interface{}{
			"name":             pinfo.Name,
			"version":          pinfo.Version,
			"release_date":     pinfo.ReleaseDate,
			"release_note_url": pinfo.ReleaseNoteUrl,
			"package_file":     pinfo.PackageFile,
			"size":             pinfo.Size,
			"platform":         pinfo.Platform,
			"installed":        pinfo.Installed,
			"downloaded":       pinfo.Downloaded,
		}

		if pinfo.Installed == "yes" {
			installed = append(installed, pinfo.Name)
		}

		details = append(details, entry)
	}

	d.SetId(id)
	if err := d.Set("details", details); err != nil {
		log.Printf("[WARN] Error setting 'info' for %q: %s", d.Id(), err)
	}
	if err := d.Set("installed", installed); err != nil {
		log.Printf("[WARN] Error setting 'installed_plugins' for %q: %s", d.Id(), err)
	}
	d.Set("total", len(list))

	return nil
}
