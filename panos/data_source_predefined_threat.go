package panos

import (
	"log"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/predefined/threat"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePredefinedThreat() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePredefinedThreatRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specific threat ID / name",
				ConflictsWith: []string{
					"threat_regex",
					"threat_name",
				},
			},
			"threat_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "A regex to apply to the threat name",
				ValidateFunc: validateIsRegex(),
				ConflictsWith: []string{
					"name",
					"threat_name",
				},
			},
			"threat_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An exact match against the threat name",
				ConflictsWith: []string{
					"name",
					"threat_regex",
				},
			},
			"threat_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The threat type",
				Default:     threat.PhoneHome,
				ValidateFunc: validateStringIn(
					threat.Vulnerability,
					threat.PhoneHome,
				),
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of threats matching the given criteria",
			},
			"threats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The matched threats",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The threat name / ID",
						},
						"threat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The threat name",
						},
					},
				},
			},
		},
	}
}

func dataSourcePredefinedThreatRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o threat.Entry
	var ans []threat.Entry

	name := d.Get("name").(string)
	regex := d.Get("threat_regex").(string)
	threatName := d.Get("threat_name").(string)
	tt := d.Get("threat_type").(string)
	id := base64Encode([]string{
		name, regex, threatName, tt,
	})

	switch c := meta.(type) {
	case *pango.Firewall:
		switch {
		case name != "":
			o, err = c.Predefined.Threat.Get(tt, name)
			ans = []threat.Entry{o}
		case threatName != "":
			ans, err = c.Predefined.Threat.GetThreats(tt, "")
			if err == nil {
				filtered := make([]threat.Entry, 0, len(ans))
				for i := range ans {
					if ans[i].ThreatName == threatName {
						filtered = append(filtered, threat.Entry{
							Name:       ans[i].Name,
							ThreatName: ans[i].ThreatName,
						})
					}
				}
				ans = filtered
			}
		default:
			ans, err = c.Predefined.Threat.GetThreats(tt, regex)
		}
	case *pango.Panorama:
		switch {
		case name != "":
			o, err = c.Predefined.Threat.Get(tt, name)
			ans = []threat.Entry{o}
		case threatName != "":
			ans, err = c.Predefined.Threat.GetThreats(tt, "")
			if err == nil {
				filtered := make([]threat.Entry, 0, len(ans))
				for i := range ans {
					if ans[i].ThreatName == threatName {
						filtered = append(filtered, threat.Entry{
							Name:       ans[i].Name,
							ThreatName: ans[i].ThreatName,
						})
					}
				}
				ans = filtered
			}
		default:
			ans, err = c.Predefined.Threat.GetThreats(tt, regex)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	threats := make([]interface{}, 0, len(ans))
	for _, x := range ans {
		threats = append(threats, map[string]interface{}{
			"name":        x.Name,
			"threat_name": x.ThreatName,
		})
	}
	d.Set("total", len(threats))
	if err = d.Set("threats", threats); err != nil {
		log.Printf("[WARN] Error setting 'threats' field for %q: %s", d.Id(), err)
	}

	return nil
}
