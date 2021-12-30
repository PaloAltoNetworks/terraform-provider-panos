package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/predefined/dlp/filetype"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePredefinedDlpFileType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePredefinedDlpFileTypeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specific file type",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specific label that you want for the given file type",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of file types",
			},
			"file_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The results",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file type",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of property specs",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The DLP property name",
									},
									"label": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The DLP property label",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePredefinedDlpFileTypeRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o filetype.Entry
	var ans []filetype.Entry

	name := d.Get("name").(string)
	label := d.Get("label").(string)
	id := base64Encode([]string{
		name, label,
	})

	switch c := meta.(type) {
	case *pango.Firewall:
		switch {
		case name != "":
			o, err = c.Predefined.DlpFileType.Get(name)
			ans = []filetype.Entry{o}
		default:
			ans, err = c.Predefined.DlpFileType.GetAll()
		}
	case *pango.Panorama:
		switch {
		case name != "":
			o, err = c.Predefined.DlpFileType.Get(name)
			ans = []filetype.Entry{o}
		default:
			ans, err = c.Predefined.DlpFileType.GetAll()
		}
	}

	if err != nil {
		return err
	}

	if label != "" {
		filtered := make([]filetype.Entry, 0, len(ans))
		for _, x := range ans {
			fx := filetype.Entry{
				Name:       x.Name,
				Properties: make([]filetype.Property, 0, len(x.Properties)),
			}
			for _, prop := range x.Properties {
				if prop.Label == label {
					fx.Properties = append(fx.Properties, filetype.Property{
						Name:  prop.Name,
						Label: prop.Label,
					})
				}
			}
			if len(fx.Properties) > 0 {
				filtered = append(filtered, fx)
			}
		}
		ans = filtered
	}

	d.SetId(id)
	types := make([]interface{}, 0, len(ans))
	for _, x := range ans {
		props := make([]interface{}, 0, len(x.Properties))
		for _, prop := range x.Properties {
			props = append(props, map[string]interface{}{
				"name":  prop.Name,
				"label": prop.Label,
			})
		}
		types = append(types, map[string]interface{}{
			"name":       x.Name,
			"properties": props,
		})
	}
	d.Set("total", len(types))
	if err = d.Set("file_types", types); err != nil {
		log.Printf("[WARN] Error setting 'file_types' field for %q: %s", d.Id(), err)
	}

	return nil
}
