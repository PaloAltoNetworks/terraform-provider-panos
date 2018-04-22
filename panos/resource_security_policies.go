package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourceSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityPolicies,
		Read:   readSecurityPolicies,
		Update: createUpdateSecurityPolicies,
		Delete: deleteSecurityPolicies,

		SchemaVersion: 1,
		MigrateState:  migrateResourceSecurityPolicies,

		Schema: map[string]*schema.Schema{
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this object in (default: vsys1)",
			},
			"rulebase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "The Panorama rulebase",
				Deprecated:   "This parameter is not really used in a firewall context.  Simply remove this setting from your plan file, as it will be removed later.",
				ValidateFunc: validateStringIn(util.Rulebase, util.PreRulebase, util.PostRulebase),
			},
			"rule": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "universal",
							Description:  "Security rule type (default: universal, interzone, intrazone)",
							ValidateFunc: validateStringIn("universal", "interzone", "intrazone"),
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_zones": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_addresses": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"negate_source": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_users": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hip_profiles": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_zones": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_addresses": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"negate_destination": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"applications": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"categories": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"action": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "allow",
							Description:  "Action (default: allow, deny, drop, reset-client, reset-server, reset-both)",
							ValidateFunc: validateStringIn("allow", "deny", "drop", "reset-client", "reset-server", "reset-both"),
						},
						"log_setting": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log forwarding profile",
						},
						"log_start": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_end": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"disabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"schedule": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_unreachable": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_server_response_inspection": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"group": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"virus": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"spyware": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"vulnerability": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"url_filtering": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_blocking": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"wildfire_analysis": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_filtering": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func migrateResourceSecurityPolicies(ov int, s *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if ov == 0 {
		t := strings.Split(s.ID, IdSeparator)
		if len(t) != 2 {
			return nil, fmt.Errorf("ID is malformed")
		} else if t[1] != util.Rulebase {
			return nil, fmt.Errorf("Rulebase is %q, not %q", t[1])
		}
		s.ID = t[0]

		ov = 1
	}

	return s, nil
}

func parseSecurityPolicies(d *schema.ResourceData) (string, string, []security.Entry) {
	vsys := d.Get("vsys").(string)
	rb := d.Get("rulebase").(string)

	rlist := d.Get("rule").([]interface{})
	ans := make([]security.Entry, 0, len(rlist))
	for i := range rlist {
		elm := rlist[i].(map[string]interface{})
		o := security.Entry{
			Name:                            elm["name"].(string),
			Type:                            elm["type"].(string),
			Description:                     elm["description"].(string),
			Tags:                            setAsList(elm["tags"].(*schema.Set)),
			SourceZones:                     asStringList(elm["source_zones"].([]interface{})),
			SourceAddresses:                 asStringList(elm["source_addresses"].([]interface{})),
			NegateSource:                    elm["negate_source"].(bool),
			SourceUsers:                     asStringList(elm["source_users"].([]interface{})),
			HipProfiles:                     asStringList(elm["hip_profiles"].([]interface{})),
			DestinationZones:                asStringList(elm["destination_zones"].([]interface{})),
			DestinationAddresses:            asStringList(elm["destination_addresses"].([]interface{})),
			NegateDestination:               elm["negate_destination"].(bool),
			Applications:                    asStringList(elm["applications"].([]interface{})),
			Services:                        asStringList(elm["services"].([]interface{})),
			Categories:                      asStringList(elm["categories"].([]interface{})),
			Action:                          elm["action"].(string),
			LogSetting:                      elm["log_setting"].(string),
			LogStart:                        elm["log_start"].(bool),
			LogEnd:                          elm["log_end"].(bool),
			Disabled:                        elm["disabled"].(bool),
			Schedule:                        elm["schedule"].(string),
			IcmpUnreachable:                 elm["icmp_unreachable"].(bool),
			DisableServerResponseInspection: elm["disable_server_response_inspection"].(bool),
			Group:            elm["group"].(string),
			Virus:            elm["virus"].(string),
			Spyware:          elm["spyware"].(string),
			Vulnerability:    elm["vulnerability"].(string),
			UrlFiltering:     elm["url_filtering"].(string),
			FileBlocking:     elm["file_blocking"].(string),
			WildFireAnalysis: elm["wildfire_analysis"].(string),
			DataFiltering:    elm["data_filtering"].(string),
		}
		ans = append(ans, o)
	}

	return vsys, rb, ans
}

func createUpdateSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, _, list := parseSecurityPolicies(d)

	if err = fw.Policies.Security.DeleteAll(vsys); err != nil {
		return err
	}
	if err = fw.Policies.Security.VerifiableSet(vsys, list...); err != nil {
		return err
	}

	d.SetId(vsys)
	return readSecurityPolicies(d, meta)
}

func readSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys := d.Id()

	list, err := fw.Policies.Security.GetList(vsys)
	if err != nil {
		return err
	}

	ilist := make([]interface{}, 0, len(list))
	for i := range list {
		o, err := fw.Policies.Security.Get(vsys, list[i])
		if err != nil {
			return err
		}
		m := make(map[string]interface{})
		m["name"] = o.Name
		m["type"] = o.Type
		m["description"] = o.Description
		m["tags"] = listAsSet(o.Tags)
		m["source_zones"] = o.SourceZones
		m["source_addresses"] = o.SourceAddresses
		m["negate_source"] = o.NegateSource
		m["source_users"] = o.SourceUsers
		m["hip_profiles"] = o.HipProfiles
		m["destination_zones"] = o.DestinationZones
		m["destination_addresses"] = o.DestinationAddresses
		m["negate_destination"] = o.NegateDestination
		m["applications"] = o.Applications
		m["services"] = o.Services
		m["categories"] = o.Categories
		m["action"] = o.Action
		m["log_setting"] = o.LogSetting
		m["log_start"] = o.LogStart
		m["log_end"] = o.LogEnd
		m["disabled"] = o.Disabled
		m["schedule"] = o.Schedule
		m["icmp_unreachable"] = o.IcmpUnreachable
		m["disable_server_response_inspection"] = o.DisableServerResponseInspection
		m["group"] = o.Group
		m["virus"] = o.Virus
		m["spyware"] = o.Spyware
		m["vulnerability"] = o.Vulnerability
		m["url_filtering"] = o.UrlFiltering
		m["file_blocking"] = o.FileBlocking
		m["wildfire_analysis"] = o.WildFireAnalysis
		m["data_filtering"] = o.DataFiltering
		ilist = append(ilist, m)
	}

	d.Set("vsys", vsys)
	d.Set("rulebase", util.Rulebase)
	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Id()

	if err := fw.Policies.Security.DeleteAll(vsys); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
