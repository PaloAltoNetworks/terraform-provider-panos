package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaSecurityPolicy,
		Read:   readPanoramaSecurityPolicy,
		Update: createUpdatePanoramaSecurityPolicy,
		Delete: deletePanoramaSecurityPolicy,

		Schema: map[string]*schema.Schema{
			"device_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
				ForceNew: true,
			},
			"rulebase": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      util.PreRulebase,
				ForceNew:     true,
				ValidateFunc: validateStringIn(util.Rulebase, util.PreRulebase, util.PostRulebase),
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "universal",
							Description:  "Security rule type (default: universal, interzone, intrazone)",
							ValidateFunc: validateStringIn("universal", "interzone", "intrazone"),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_zones": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"source_addresses": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"negate_source": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"source_users": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"hip_profiles": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_zones": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"destination_addresses": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"negate_destination": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"applications": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"categories": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "allow",
							Description:  "Action (default: allow, deny, drop, reset-client, reset-server, reset-both)",
							ValidateFunc: validateStringIn("allow", "deny", "drop", "reset-client", "reset-server", "reset-both"),
						},
						"log_setting": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log forwarding profile",
						},
						"log_start": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"log_end": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp_unreachable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"disable_server_response_inspection": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"virus": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"spyware": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vulnerability": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url_filtering": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"file_blocking": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"wildfire_analysis": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"data_filtering": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target": {
							Type:     schema.TypeSet,
							Optional: true,
							// TODO(gfreeman): Uncomment once ValidateFunc is supported for TypeSet.
							//ValidateFunc: validateSetKeyIsUnique("serial"),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"serial": {
										Type:     schema.TypeString,
										Required: true,
									},
									"vsys_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"negate_target": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func parsePanoramaSecurityPolicy(d *schema.ResourceData) (string, string, []security.Entry) {
	dg := d.Get("device_group").(string)
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
			Group:                           elm["group"].(string),
			Virus:                           elm["virus"].(string),
			Spyware:                         elm["spyware"].(string),
			Vulnerability:                   elm["vulnerability"].(string),
			UrlFiltering:                    elm["url_filtering"].(string),
			FileBlocking:                    elm["file_blocking"].(string),
			WildFireAnalysis:                elm["wildfire_analysis"].(string),
			DataFiltering:                   elm["data_filtering"].(string),
			NegateTarget:                    elm["negate_target"].(bool),
		}

		m := make(map[string][]string)
		tl := elm["target"].(*schema.Set).List()
		for i := range tl {
			device := tl[i].(map[string]interface{})
			key := device["serial"].(string)
			value := asStringList(device["vsys_list"].(*schema.Set).List())
			m[key] = value
		}
		o.Targets = m

		ans = append(ans, o)
	}

	return dg, rb, ans
}

func parsePanoramaSecurityPolicyId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaSecurityPolicyId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createUpdatePanoramaSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, list := parsePanoramaSecurityPolicy(d)

	if err = pano.Policies.Security.DeleteAll(dg, rb); err != nil {
		return err
	}
	if err = pano.Policies.Security.VerifiableSet(dg, rb, list...); err != nil {
		return err
	}

	d.SetId(buildPanoramaSecurityPolicyId(dg, rb))
	return readPanoramaSecurityPolicy(d, meta)
}

func readPanoramaSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb := parsePanoramaSecurityPolicyId(d.Id())

	list, err := pano.Policies.Security.GetList(dg, rb)
	if err != nil {
		return err
	}

	ts2 := d.Get("rule").([]interface{})
	elm := ts2[0].(map[string]interface{})
	ts := elm["target"].(*schema.Set)

	ilist := make([]interface{}, 0, len(list))
	for i := range list {
		o, err := pano.Policies.Security.Get(dg, rb, list[i])
		if err != nil {
			return err
		}

		s := &schema.Set{F: ts.F}
		for key := range o.Targets {
			sg := make(map[string]interface{})
			sg["serial"] = key
			sg["vsys_list"] = listAsSet(o.Targets[key])
			s.Add(sg)
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
		m["target"] = s
		m["negate_target"] = o.NegateTarget
		ilist = append(ilist, m)
	}

	d.Set("device_group", dg)
	d.Set("rulebase", rb)
	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deletePanoramaSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, rb := parsePanoramaSecurityPolicyId(d.Id())

	if err := pano.Policies.Security.DeleteAll(dg, rb); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
