package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityPolicies,
		Read:   readSecurityPolicies,
		Update: createUpdateSecurityPolicies,
		Delete: deleteSecurityPolicies,

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
				Default:      "rulebase",
				ForceNew:     true,
				Description:  "The rulebase (default: rulebase, pre-rulebase, post-rulebase)",
				ValidateFunc: validateStringIn("rulebase", "pre-rulebase", "post-rulebase"),
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
						"source_zone": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_address": &schema.Schema{
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
						"source_user": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hip_profile": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_zone": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_address": &schema.Schema{
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
						"application": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"service": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"category": &schema.Schema{
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
			SourceZone:                      asStringList(elm["source_zone"].([]interface{})),
			SourceAddress:                   asStringList(elm["source_address"].([]interface{})),
			NegateSource:                    elm["negate_source"].(bool),
			SourceUser:                      asStringList(elm["source_user"].([]interface{})),
			HipProfile:                      asStringList(elm["hip_profile"].([]interface{})),
			DestinationZone:                 asStringList(elm["destination_zone"].([]interface{})),
			DestinationAddress:              asStringList(elm["destination_address"].([]interface{})),
			NegateDestination:               elm["negate_destination"].(bool),
			Application:                     asStringList(elm["application"].([]interface{})),
			Service:                         asStringList(elm["service"].([]interface{})),
			Category:                        asStringList(elm["category"].([]interface{})),
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

func parseSecurityPoliciesId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildSecurityPoliciesId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createUpdateSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, rb, list := parseSecurityPolicies(d)

	if err = fw.Policies.Security.DeleteAll(vsys, rb); err != nil {
		return err
	}
	if err = fw.Policies.Security.VerifiableSet(vsys, rb, list...); err != nil {
		return err
	}

	d.SetId(buildSecurityPoliciesId(vsys, rb))
	return readSecurityPolicies(d, meta)
}

func readSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, rb := parseSecurityPoliciesId(d.Id())

	list, err := fw.Policies.Security.GetList(vsys, rb)
	if err != nil {
		return err
	}

	ilist := make([]interface{}, 0, len(list))
	for i := range list {
		o, err := fw.Policies.Security.Get(vsys, rb, list[i])
		if err != nil {
			return err
		}
		m := make(map[string]interface{})
		m["name"] = o.Name
		m["type"] = o.Type
		m["description"] = o.Description
		m["tags"] = listAsSet(o.Tags)
		m["source_zone"] = o.SourceZone
		m["source_address"] = o.SourceAddress
		m["negate_source"] = o.NegateSource
		m["source_user"] = o.SourceUser
		m["hip_profile"] = o.HipProfile
		m["destination_zone"] = o.DestinationZone
		m["destination_address"] = o.DestinationAddress
		m["negate_destination"] = o.NegateDestination
		m["application"] = o.Application
		m["service"] = o.Service
		m["category"] = o.Category
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
	d.Set("rulebase", rb)
	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteSecurityPolicies(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, rb := parseSecurityPoliciesId(d.Id())

	if err := fw.Policies.Security.DeleteAll(vsys, rb); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
