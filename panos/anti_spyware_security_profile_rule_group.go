package panos

import (
	"log"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware/rule"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceAntiSpywareSecurityProfileRules() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema()
	s["device_group"] = deviceGroupSchema()
	s["anti_spyware_security_profile"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The anti-spyware security profile name",
	}

	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfileRulesRead,

		Schema: s,
	}
}

func dataSourceAntiSpywareSecurityProfileRulesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	prof := d.Get("anti_spyware_security_profile").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = strings.Join([]string{vsys, prof}, IdSeparator)
		listing, err = con.Objects.AntiSpywareRule.GetList(vsys, prof)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = strings.Join([]string{dg, prof}, IdSeparator)
		listing, err = con.Objects.AntiSpywareRule.GetList(dg, prof)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceAntiSpywareSecurityProfileRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfileRuleRead,

		Schema: antiSpywareSecurityProfileRuleGroupSchema(false),
	}
}

func dataSourceAntiSpywareSecurityProfileRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o rule.Entry
	prof := d.Get("anti_spyware_security_profile").(string)
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileRuleGroupId(vsys, prof, 0, "", []interface{}{name})
		o, err = con.Objects.AntiSpywareRule.Get(vsys, prof, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileRuleGroupId(dg, prof, 0, "", []interface{}{name})
		o, err = con.Objects.AntiSpywareRule.Get(dg, prof, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveAntiSpywareSecurityProfileRuleGroup(d, []rule.Entry{o})

	return nil
}

// Resource.
func resourceAntiSpywareSecurityProfileRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateAntiSpywareSecurityProfileRuleGroup,
		Read:   readAntiSpywareSecurityProfileRuleGroup,
		Update: createUpdateAntiSpywareSecurityProfileRuleGroup,
		Delete: deleteAntiSpywareSecurityProfileRuleGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: antiSpywareSecurityProfileRuleGroupSchema(true),
	}
}

func createUpdateAntiSpywareSecurityProfileRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	move, oRule, prof, rules := loadAntiSpywareSecurityProfileRuleGroup(d)

	names := make([]interface{}, 0, len(rules))
	for i := range rules {
		names = append(names, rules[i].Name)
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileRuleGroupId(vsys, prof, move, oRule, names)
		_ = con.Objects.AntiSpywareRule.Delete(vsys, prof, names...)
		if err = con.Objects.AntiSpywareRule.Set(vsys, prof, rules...); err != nil {
			return err
		}
		err = con.Objects.AntiSpywareRule.MoveGroup(vsys, prof, move, oRule, rules...)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileRuleGroupId(dg, prof, move, oRule, names)
		_ = con.Objects.AntiSpywareRule.Delete(dg, prof, names...)
		if err = con.Objects.AntiSpywareRule.Set(dg, prof, rules...); err != nil {
			return err
		}
		err = con.Objects.AntiSpywareRule.MoveGroup(dg, prof, move, oRule, rules...)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAntiSpywareSecurityProfileRuleGroup(d, meta)
}

func readAntiSpywareSecurityProfileRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var rules, names []string
	var listing []rule.Entry
	var fIdx, oIdx, move int

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, prof, move, oRule, names := parseAntiSpywareSecurityProfileRuleGroupId(d.Id())
		listing = make([]rule.Entry, 0, len(names))
		rules, err = con.Objects.AntiSpywareRule.GetList(vsys, prof)
		if err != nil {
			return err
		}

		fIdx, oIdx, err = groupIndexes(rules, names, move, oRule)
		if fIdx == -1 {
			d.SetId("")
			return nil
		} else if err != nil {
			return err
		}

		for i := 0; i+fIdx < len(rules) && i < len(names); i++ {
			if rules[i+fIdx] != names[i] {
				// Must be contiguous.
				break
			}
			o, err := con.Objects.AntiSpywareRule.Get(vsys, prof, names[i])
			if err != nil {
				if isObjectNotFound(err) {
					break
				}
				return err
			}
			listing = append(listing, o)
		}
	case *pango.Panorama:
		dg, prof, move, oRule, names := parseAntiSpywareSecurityProfileRuleGroupId(d.Id())
		listing = make([]rule.Entry, 0, len(names))
		rules, err = con.Objects.AntiSpywareRule.GetList(dg, prof)
		if err != nil {
			return err
		}

		fIdx, oIdx, err = groupIndexes(rules, names, move, oRule)
		if fIdx == -1 {
			d.SetId("")
			return nil
		} else if err != nil {
			return err
		}

		for i := 0; i+fIdx < len(rules) && i < len(names); i++ {
			if rules[i+fIdx] != names[i] {
				// Must be contiguous.
				break
			}
			o, err := con.Objects.AntiSpywareRule.Get(dg, prof, names[i])
			if err != nil {
				if isObjectNotFound(err) {
					break
				}
				return err
			}
			listing = append(listing, o)
		}
	}

	if !groupPositionIsOk(move, fIdx, oIdx, rules, names) {
		d.Set("position_reference", "(incorrect group positioning)")
	}
	saveAntiSpywareSecurityProfileRuleGroup(d, listing)
	return nil
}

func deleteAntiSpywareSecurityProfileRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, prof, _, _, names := parseAntiSpywareSecurityProfileRuleGroupId(d.Id())
		ilist := make([]interface{}, len(names))
		for i := range names {
			ilist[i] = names[i]
		}
		err = con.Objects.AntiSpywareRule.Delete(vsys, prof, ilist...)
	case *pango.Panorama:
		dg, prof, _, _, names := parseAntiSpywareSecurityProfileRuleGroupId(d.Id())
		ilist := make([]interface{}, len(names))
		for i := range names {
			ilist[i] = names[i]
		}
		err = con.Objects.AntiSpywareRule.Delete(dg, prof, ilist...)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func antiSpywareSecurityProfileRuleGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group":       deviceGroupSchema(),
		"vsys":               vsysSchema(),
		"position_keyword":   positionKeywordSchema(),
		"position_reference": positionReferenceSchema(),
		"anti_spyware_security_profile": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The anti-spyware security policy name",
			ForceNew:    true,
		},
		"rule": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Rule specs",
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"threat_name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Threat name",
					},
					"category": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The category",
					},
					"severities": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of severities",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"packet_capture": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Packet capture setting",
						ValidateFunc: validateStringIn(
							"", rule.Disable, rule.SinglePacket, rule.ExtendedCapture,
						),
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Action to take",
						ValidateFunc: validateStringIn(
							"",
							rule.ActionDefault,
							rule.ActionAllow,
							rule.ActionAlert,
							rule.ActionDrop,
							rule.ActionResetClient,
							rule.ActionResetServer,
							rule.ActionResetBoth,
							rule.ActionBlockIp,
						),
					},
					"block_ip_track_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "(For action = block-ip) The track by setting",
						ValidateFunc: validateStringIn(
							"",
							rule.TrackBySource,
							rule.TrackBySourceAndDestination,
						),
					},
					"block_ip_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "(For action = block-ip) The duration",
					},
				},
			},
		},
	}

	if !isResource {
		ans["name"] = &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The rule name",
		}

		computed(ans, "", []string{
			"vsys", "device_group", "anti_spyware_security_profile", "name",
		})
	}

	return ans
}

func loadAntiSpywareSecurityProfileRuleGroup(d *schema.ResourceData) (int, string, string, []rule.Entry) {
	prof := d.Get("anti_spyware_security_profile").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))
	var rules []rule.Entry

	list := d.Get("rule").([]interface{})
	if len(list) > 0 {
		rules = make([]rule.Entry, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			rules = append(rules, rule.Entry{
				Name:            elm["name"].(string),
				ThreatName:      elm["threat_name"].(string),
				Category:        elm["category"].(string),
				Severities:      asStringList(elm["severities"].([]interface{})),
				PacketCapture:   elm["packet_capture"].(string),
				Action:          elm["action"].(string),
				BlockIpTrackBy:  elm["block_ip_track_by"].(string),
				BlockIpDuration: elm["block_ip_duration"].(int),
			})
		}
	}

	return move, oRule, prof, rules
}

func saveAntiSpywareSecurityProfileRuleGroup(d *schema.ResourceData, listing []rule.Entry) {
	if len(listing) == 0 {
		d.Set("rule", nil)
		return
	}

	data := make([]interface{}, 0, len(listing))
	for _, o := range listing {
		data = append(data, map[string]interface{}{
			"name":              o.Name,
			"threat_name":       o.ThreatName,
			"category":          o.Category,
			"severities":        o.Severities,
			"packet_capture":    o.PacketCapture,
			"action":            o.Action,
			"block_ip_track_by": o.BlockIpTrackBy,
			"block_ip_duration": o.BlockIpDuration,
		})
	}

	if err := d.Set("rule", data); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildAntiSpywareSecurityProfileRuleGroupId(a, b string, c int, d string, e []interface{}) string {
	return strings.Join([]string{a, b, strconv.Itoa(c), d, base64Encode(e)}, IdSeparator)
}

func parseAntiSpywareSecurityProfileRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	return t[0], t[1], move, t[3], base64Decode(t[4])
}
