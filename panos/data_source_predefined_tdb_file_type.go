package panos

import (
	"fmt"
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/predefined/tdb/filetype"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePredefinedTdbFileType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePredefinedTdbFileTypeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specific file type",
				ConflictsWith: []string{
					"full_name",
					"full_name_regex",
				},
			},
			"full_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The full name",
				ConflictsWith: []string{
					"name",
					"full_name_regex",
				},
			},
			"full_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A regex to match against the full name",
				ValidateFunc: validateIsRegex(),
				ConflictsWith: []string{
					"name",
					"full_name",
				},
			},
			"data_ident_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Limit results to those with data_ident = true",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of file types",
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
						"file_type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID",
						},
						"threat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The threat name",
						},
						"full_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The full name",
						},
						"data_ident": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Data ident",
						},
						"file_type_ident": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "File type ident",
						},
					},
				},
			},
		},
	}
}

func dataSourcePredefinedTdbFileTypeRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o filetype.Entry
	var ans []filetype.Entry

	name := d.Get("name").(string)
	fullName := d.Get("full_name").(string)
	fullNameRegex := d.Get("full_name_regex").(string)
	dio := d.Get("data_ident_only").(bool)
	id := base64Encode([]string{
		name, fullName, fullNameRegex, fmt.Sprintf("%t", dio),
	})

	switch c := meta.(type) {
	case *pango.Firewall:
		switch {
		case name != "":
			o, err = c.Predefined.TdbFileType.Get(name)
			ans = []filetype.Entry{o}
		case fullName != "":
			ans, err := c.Predefined.TdbFileType.GetMatches("")
			if err == nil {
				filtered := make([]filetype.Entry, 0, len(ans))
				for _, x := range ans {
					if x.FullName == fullName {
						filtered = append(filtered, x)
					}
				}
				ans = filtered
			}
		default:
			ans, err = c.Predefined.TdbFileType.GetMatches(fullNameRegex)
		}
	case *pango.Panorama:
		switch {
		case name != "":
			o, err = c.Predefined.TdbFileType.Get(name)
			ans = []filetype.Entry{o}
		case fullName != "":
			ans, err := c.Predefined.TdbFileType.GetMatches("")
			if err == nil {
				filtered := make([]filetype.Entry, 0, len(ans))
				for _, x := range ans {
					if x.FullName == fullName {
						filtered = append(filtered, x)
					}
				}
				ans = filtered
			}
		default:
			ans, err = c.Predefined.TdbFileType.GetMatches(fullNameRegex)
		}
	}

	if err != nil {
		return err
	}

	if dio {
		filtered := make([]filetype.Entry, 0, len(ans))
		for _, x := range ans {
			if x.DataIdent {
				filtered = append(filtered, x)
			}
		}
		ans = filtered
	}

	d.SetId(id)
	types := make([]interface{}, 0, len(ans))
	for _, x := range ans {
		types = append(types, map[string]interface{}{
			"name":            x.Name,
			"file_type_id":    x.Id,
			"threat_name":     x.ThreatName,
			"full_name":       x.FullName,
			"data_ident":      x.DataIdent,
			"file_type_ident": x.FileTypeIdent,
		})
	}
	d.Set("total", len(types))
	if err = d.Set("file_types", types); err != nil {
		log.Printf("[WARN] Error setting 'threats' field for %q: %s", d.Id(), err)
	}

	return nil
}
