package panos

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaSecurityRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaSecurityRuleGroup,
		Read:   readPanoramaSecurityRuleGroup,
		Update: createUpdatePanoramaSecurityRuleGroup,
		Delete: deletePanoramaSecurityRuleGroup,

		Schema: map[string]*schema.Schema{
			"device_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
				ForceNew: true,
			},
			"rulebase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      util.PreRulebase,
				ForceNew:     true,
				ValidateFunc: validateStringIn(util.Rulebase, util.PreRulebase, util.PostRulebase),
			},
			"position_keyword": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validateStringIn(movementKeywords()...),
				ForceNew:     true,
			},
			"position_reference": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"rule": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "universal",
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
							ValidateFunc: validateStringIn("allow", "deny", "drop", "reset-client", "reset-server", "reset-both"),
						},
						"log_setting": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
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
						"target": &schema.Schema{
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
						"negate_target": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func parsePanoramaSecurityRuleGroup(d *schema.ResourceData) (string, string, string, int, []security.Entry) {
	dg := d.Get("device_group").(string)
	rb := d.Get("rulebase").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))

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

	return dg, rb, oRule, move, ans
}

func parsePanoramaSecurityRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	joined, _ := base64.StdEncoding.DecodeString(t[4])
	names := strings.Split(string(joined), "\n")
	return t[0], t[1], move, t[3], names
}

func buildPanoramaSecurityRuleGroupId(a, b string, c int, d string, e []security.Entry) string {
	var buf bytes.Buffer
	for i := range e {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(e[i].Name)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())

	return fmt.Sprintf("%s%s%s%s%d%s%s%s%s", a, IdSeparator, b, IdSeparator, c, IdSeparator, d, IdSeparator, enc)
}

func createUpdatePanoramaSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, oRule, move, list := parsePanoramaSecurityRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = pano.Policies.Security.VerifiableEdit(dg, rb, list...); err != nil {
		return err
	}
	if err = pano.Policies.Security.MoveGroup(dg, rb, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildPanoramaSecurityRuleGroupId(dg, rb, move, oRule, list))
	return readPanoramaSecurityRuleGroup(d, meta)
}

func readPanoramaSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, move, oRule, policies := parsePanoramaSecurityRuleGroupId(d.Id())

	list, err := pano.Policies.Security.GetList(dg, rb)
	if err != nil {
		return err
	}

	fIdx, oIdx := -1, -1
	for i := range list {
		if list[i] == policies[0] {
			fIdx = i
		} else if list[i] == oRule {
			oIdx = i
		}
		if fIdx != -1 && oIdx != -1 {
			break
		}
	}

	if fIdx == -1 {
		// First policy is MIA, but others may be present, so report an
		// empty ruleset to force rules to be recreated.
		d.Set("rule", nil)
		return nil
	} else if oIdx == -1 && movementIsRelative(move) {
		return fmt.Errorf("Can't position group %s %q: rule is not present", movementItoa(move), oRule)
	}

	d.Set("device_group", dg)
	d.Set("rulebase", rb)
	d.Set("position_keyword", movementItoa(move))
	if groupPositionIsOk(move, fIdx, oIdx, list, policies) {
		d.Set("position_reference", oRule)
	} else {
		d.Set("position_reference", "(incorrect group positioning)")
	}

	ts2 := d.Get("rule").([]interface{})
	elm := ts2[0].(map[string]interface{})
	ts := elm["target"].(*schema.Set)

	ilist := make([]interface{}, 0, len(policies))
	for i := 0; i+fIdx < len(list) && i < len(policies); i++ {
		if list[i+fIdx] != policies[i] {
			// Policies must be contiguous.
			break
		}
		o, err := pano.Policies.Security.Get(dg, rb, policies[i])
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

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deletePanoramaSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, rb, _, _, list := parsePanoramaSecurityRuleGroupId(d.Id())

	ilist := make([]interface{}, len(list))
	for i := range list {
		ilist[i] = list[i]
	}

	if err := pano.Policies.Security.Delete(dg, rb, ilist...); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
