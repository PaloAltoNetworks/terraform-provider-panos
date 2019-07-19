package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePanoramaPlugin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePanoramaPluginRead,

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

func dataSourcePanoramaPluginRead(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)

	details := make([]interface{}, 0, len(pano.Plugin))
	installed := make([]string, 0, len(pano.Plugin))
	for _, pinfo := range pano.Plugin {
		entry := map[string]interface{}{
			"name":             pinfo["name"],
			"version":          pinfo["version"],
			"release_date":     pinfo["release-date"],
			"release_note_url": pinfo["release-note-url"],
			"package_file":     pinfo["package-file"],
			"size":             pinfo["size"],
			"platform":         pinfo["platform"],
			"installed":        pinfo["installed"],
			"downloaded":       pinfo["downloaded"],
		}

		if pinfo["installed"] == "yes" {
			installed = append(installed, pinfo["name"])
		}

		details = append(details, entry)
	}

	d.SetId(pano.Hostname)
	if err := d.Set("details", details); err != nil {
		log.Printf("[WARN] Error setting 'info' for %q: %s", d.Id(), err)
	}
	if err := d.Set("installed", installed); err != nil {
		log.Printf("[WARN] Error setting 'installed_plugins' for %q: %s", d.Id(), err)
	}
	d.Set("total", len(pano.Plugin))

	return nil
}
