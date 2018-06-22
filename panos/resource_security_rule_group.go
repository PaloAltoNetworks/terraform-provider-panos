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

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityRuleGroup,
		Read:   readSecurityRuleGroup,
		Update: createUpdateSecurityRuleGroup,
		Delete: deleteSecurityRuleGroup,

		Schema: map[string]*schema.Schema{
			"vsys": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
				ForceNew: true,
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
					},
				},
			},
		},
	}
}

func parseSecurityRuleGroup(d *schema.ResourceData) (string, string, int, []security.Entry) {
	vsys := d.Get("vsys").(string)
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

	return vsys, oRule, move, ans
}

func parseSecurityRuleGroupId(v string) (string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[1])
	joined, _ := base64.StdEncoding.DecodeString(t[3])
	names := strings.Split(string(joined), "\n")
	return t[0], move, t[2], names
}

func buildSecurityRuleGroupId(a string, b int, c string, d []security.Entry) string {
	var buf bytes.Buffer
	for i := range d {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(d[i].Name)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())

	return fmt.Sprintf("%s%s%d%s%s%s%s", a, IdSeparator, b, IdSeparator, c, IdSeparator, enc)
}

func createUpdateSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, oRule, move, list := parseSecurityRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = fw.Policies.Security.VerifiableEdit(vsys, list...); err != nil {
		return err
	}
	if err = fw.Policies.Security.MoveGroup(vsys, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildSecurityRuleGroupId(vsys, move, oRule, list))
	return readSecurityRuleGroup(d, meta)
}

func readSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, move, oRule, policies := parseSecurityRuleGroupId(d.Id())

	list, err := fw.Policies.Security.GetList(vsys)
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

	d.Set("vsys", vsys)
	d.Set("position_keyword", movementItoa(move))
	if groupPositionIsOk(move, fIdx, oIdx, list, policies) {
		d.Set("position_reference", oRule)
	} else {
		d.Set("position_reference", "(incorrect group positioning)")
	}

	ilist := make([]interface{}, 0, len(policies))
	for i := 0; i+fIdx < len(list) && i < len(policies); i++ {
		if list[i+fIdx] != policies[i] {
			// Policies must be contiguous.
			break
		}
		o, err := fw.Policies.Security.Get(vsys, policies[i])
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

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, _, _, list := parseSecurityRuleGroupId(d.Id())

	ilist := make([]interface{}, len(list))
	for i := range list {
		ilist[i] = list[i]
	}

	if err := fw.Policies.Security.Delete(vsys, ilist...); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
