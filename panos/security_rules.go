package panos

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceSecurityRules() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()
	s["rulebase"] = rulebaseSchema()

	return &schema.Resource{
		Read: dataSourceSecurityRulesRead,

		Schema: s,
	}
}

func dataSourceSecurityRulesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)

	id := buildSecurityPolicyId(dg, base, vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Policies.Security.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Policies.Security.GetList(dg, base)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceSecurityRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityRuleRead,

		Schema: securityRuleSchema(false, 0, []string{"position_keyword", "position_reference"}),
	}
}

func dataSourceSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o security.Entry

	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)
	d.Set("name", name)

	id := strings.Join([]string{dg, base, vsys, name}, IdSeparator)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Policies.Security.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Policies.Security.Get(dg, base, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveSecurityRules(d, []security.Entry{o})

	return nil
}

// Resource (group).
func resourceSecurityRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityRuleGroup,
		Read:   readSecurityRuleGroup,
		Update: createUpdateSecurityRuleGroup,
		Delete: deleteSecurityRuleGroup,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: securityRuleSchema(true, 1, []string{"device_group", "rulebase", "name"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: securityRuleUpgradeV0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: securityRuleSchema(true, 1, []string{"name"}),
	}
}

func securityRuleUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["rulebase"]; !ok {
		raw["rulebase"] = util.PreRulebase
	}
	if _, ok := raw["device_group"]; !ok {
		raw["device_group"] = "shared"
	}
	if _, ok := raw["vsys"]; !ok {
		raw["vsys"] = "vsys1"
	}

	return raw, nil
}

func resourcePanoramaSecurityRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityRuleGroup,
		Read:   readSecurityRuleGroup,
		Update: createUpdateSecurityRuleGroup,
		Delete: deleteSecurityRuleGroup,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: securityRuleSchema(true, 1, []string{"vsys", "name"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: securityRuleUpgradeV0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: securityRuleSchema(true, 1, []string{"name"}),
	}
}

func createUpdateSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var prevNames []string

	move := movementAtoi(d.Get("position_keyword").(string))
	oRule := d.Get("position_reference").(string)
	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)
	rules, auditComments := loadSecurityRules(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}

	d.Set("position_keyword", movementItoa(move))
	d.Set("position_reference", oRule)
	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)

	if d.Id() != "" {
		_, _, _, _, _, prevNames = parseSecurityRuleGroupId(d.Id())
	}

	id := buildSecurityRuleGroupId(dg, base, vsys, move, oRule, rules)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Security.ConfigureRules(vsys, rules, auditComments, false, move, oRule, prevNames)
	case *pango.Panorama:
		err = con.Policies.Security.ConfigureRules(dg, base, rules, auditComments, false, move, oRule, prevNames)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSecurityRuleGroup(d, meta)
}

func readSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []security.Entry

	// Migrate the ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 4 {
		d.SetId(strings.Join([]string{
			"shared", util.PreRulebase, tok[0], tok[1], tok[2], tok[3],
		}, IdSeparator))
	} else if len(tok) == 5 {
		d.SetId(strings.Join([]string{
			tok[0], tok[1], "vsys1", tok[2], tok[3], tok[4],
		}, IdSeparator))
	} else if len(tok) != 6 {
		return fmt.Errorf("Invalid ID len(%d) encountered: %s", len(tok), d.Id())
	}

	dg, base, vsys, move, oRule, names := parseSecurityRuleGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Policies.Security.GetAll(vsys)
	case *pango.Panorama:
		listing, err = con.Policies.Security.GetAll(dg, base)
	}

	if err != nil {
		d.SetId("")
		return nil
	}

	fIdx, oIdx := -1, -1
	for i := range listing {
		if listing[i].Name == names[0] {
			fIdx = i
		} else if listing[i].Name == oRule {
			oIdx = i
		}
		if fIdx != -1 && (oIdx != -1 || oRule == "") {
			break
		}
	}

	if fIdx == -1 {
		// First rule is MIA, but others may be present, so report an
		// empty ruleset to force rules to be recreated.
		d.Set("rule", nil)
		return nil
	} else if oIdx == -1 && movementIsRelative(move) {
		return fmt.Errorf("Can't position group %s %q: rule is not present", movementItoa(move), oRule)
	} else if move == util.MoveTop && fIdx != 0 {
		d.Set("position_keyword", "")
	}

	dlist := make([]security.Entry, 0, len(names))
	for i := 0; i+fIdx < len(listing) && i < len(names); i++ {
		if listing[i+fIdx].Name != names[i] {
			break
		}

		dlist = append(dlist, listing[i+fIdx])
	}

	if move == util.MoveBottom && dlist[len(dlist)-1].Name != listing[len(listing)-1].Name {
		d.Set("position_keyword", "")
	}
	saveSecurityRules(d, dlist)

	return nil
}

func deleteSecurityRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, base, vsys, _, _, names := parseSecurityRuleGroupId(d.Id())

	ilist := make([]interface{}, len(names))
	for i := range names {
		ilist[i] = names[i]
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Security.Delete(vsys, ilist...)
	case *pango.Panorama:
		err = con.Policies.Security.Delete(dg, base, ilist...)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (policy).
func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityPolicy,
		Read:   readSecurityPolicy,
		Update: createUpdateSecurityPolicy,
		Delete: deleteSecurityPolicy,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: securityRuleSchema(true, 0, []string{"position_keyword", "position_reference", "device_group", "name"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: securityRuleUpgradeV0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: securityRuleSchema(true, 0, []string{"position_keyword", "position_reference", "name"}),
	}
}

func resourcePanoramaSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSecurityPolicy,
		Read:   readSecurityPolicy,
		Update: createUpdateSecurityPolicy,
		Delete: deleteSecurityPolicy,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: securityRuleSchema(true, 0, []string{"position_keyword", "position_reference", "vsys", "name"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: securityRuleUpgradeV0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: securityRuleSchema(true, 0, []string{"position_keyword", "position_reference", "name"}),
	}
}

func createUpdateSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	vsys := d.Get("vsys").(string)
	rules, auditComments := loadSecurityRules(d)

	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("vsys", vsys)

	id := buildSecurityPolicyId(dg, base, vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Security.ConfigureRules(vsys, rules, auditComments, true, util.MoveTop, "", nil)
	case *pango.Panorama:
		err = con.Policies.Security.ConfigureRules(dg, base, rules, auditComments, true, util.MoveTop, "", nil)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSecurityPolicy(d, meta)
}

func readSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []security.Entry

	// Migrate the ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 1 {
		// Old NGFW ID format.
		d.SetId(buildSecurityPolicyId("shared", util.PreRulebase, tok[0]))
	} else if len(tok) == 2 {
		// This can be either the new NGFW format or the current Panorama format.
		switch meta.(type) {
		case *pango.Firewall:
			d.SetId(buildSecurityPolicyId("shared", tok[0], tok[1]))
		case *pango.Panorama:
			d.SetId(buildSecurityPolicyId(tok[0], tok[1], "vsys1"))
		}
	} else if len(tok) != 3 {
		// Some random incorrect format..?
		d.SetId("")
		return nil
	}

	dg, base, vsys := parseSecurityPolicyId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Policies.Security.GetAll(vsys)
	case *pango.Panorama:
		listing, err = con.Policies.Security.GetAll(dg, base)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSecurityRules(d, listing)

	return nil
}

func deleteSecurityPolicy(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, base, vsys := parseSecurityPolicyId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Security.DeleteAll(vsys)
	case *pango.Panorama:
		err = con.Policies.Security.DeleteAll(dg, base)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func securityRuleSchema(isResource bool, ruleMin int, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group":       deviceGroupSchema(),
		"rulebase":           rulebaseSchema(),
		"vsys":               vsysSchema("vsys1"),
		"position_keyword":   positionKeywordSchema(),
		"position_reference": positionReferenceSchema(),
		"name": {
			Type:        schema.TypeString,
			Description: "The rule name.",
			Required:    true,
			ForceNew:    true,
		},
		"rule": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: ruleMin,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "The name.",
						Required:    true,
					},
					"uuid": uuidSchema(),
					"type": {
						Type:         schema.TypeString,
						Description:  "Rule type.",
						Optional:     true,
						Default:      "universal",
						ValidateFunc: validateStringIn("universal", "interzone", "intrazone"),
					},
					"description": {
						Type:        schema.TypeString,
						Description: "The description.",
						Optional:    true,
					},
					"tags": tagSchema(),
					"source_zones": {
						Type:        schema.TypeSet,
						Description: "List of source zones.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"source_addresses": {
						Type:        schema.TypeSet,
						Description: "List of source addresses.",
						Required:    true,
						MinItems:    1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"negate_source": {
						Type:        schema.TypeBool,
						Description: "Negate the source addresses.",
						Optional:    true,
					},
					"source_users": {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"hip_profiles": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_zones": {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_addresses": {
						Type:     schema.TypeSet,
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
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"services": {
						Type:     schema.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"categories": {
						Type:     schema.TypeSet,
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
						ValidateFunc: validateStringIn("allow", "deny", "drop", "reset-client", "reset-server", "reset-both"),
					},
					"log_setting": {
						Type:     schema.TypeString,
						Optional: true,
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
					"group_tag": groupTagSchema(),
					"source_devices": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"destination_devices": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"target":        targetSchema(false),
					"negate_target": negateTargetSchema(),
					"audit_comment": auditCommentSchema(),
				},
			},
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"device_group", "rulebase", "vsys", "name"})
	}

	return ans
}

func loadSecurityRules(d *schema.ResourceData) ([]security.Entry, map[string]string) {
	auditComments := make(map[string]string)
	rlist := d.Get("rule").([]interface{})
	if len(rlist) == 0 {
		return nil, auditComments
	}

	ans := make([]security.Entry, 0, len(rlist))
	for i := range rlist {
		elm := rlist[i].(map[string]interface{})
		auditComments[elm["name"].(string)] = elm["audit_comment"].(string)
		ans = append(ans, security.Entry{
			Name:                            elm["name"].(string),
			Type:                            elm["type"].(string),
			Description:                     elm["description"].(string),
			Tags:                            asStringList(elm["tags"].([]interface{})),
			SourceZones:                     setAsList(elm["source_zones"].(*schema.Set)),
			SourceAddresses:                 setAsList(elm["source_addresses"].(*schema.Set)),
			NegateSource:                    elm["negate_source"].(bool),
			SourceUsers:                     setAsList(elm["source_users"].(*schema.Set)),
			HipProfiles:                     setAsList(elm["hip_profiles"].(*schema.Set)),
			DestinationZones:                setAsList(elm["destination_zones"].(*schema.Set)),
			DestinationAddresses:            setAsList(elm["destination_addresses"].(*schema.Set)),
			NegateDestination:               elm["negate_destination"].(bool),
			Applications:                    setAsList(elm["applications"].(*schema.Set)),
			Services:                        setAsList(elm["services"].(*schema.Set)),
			Categories:                      setAsList(elm["categories"].(*schema.Set)),
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
			GroupTag:                        elm["group_tag"].(string),
			SourceDevices:                   setAsList(elm["source_devices"].(*schema.Set)),
			DestinationDevices:              setAsList(elm["destination_devices"].(*schema.Set)),
			Targets:                         loadTarget(elm["target"]),
			NegateTarget:                    elm["negate_target"].(bool),
		})
	}

	return ans, auditComments
}

func saveSecurityRules(d *schema.ResourceData, rules []security.Entry) {
	if len(rules) == 0 {
		d.Set("rule", nil)
		return
	}

	list := make([]interface{}, 0, len(rules))
	for _, x := range rules {
		list = append(list, map[string]interface{}{
			"name":                               x.Name,
			"uuid":                               x.Uuid,
			"type":                               x.Type,
			"description":                        x.Description,
			"tags":                               x.Tags,
			"source_zones":                       listAsSet(x.SourceZones),
			"source_addresses":                   listAsSet(x.SourceAddresses),
			"negate_source":                      x.NegateSource,
			"source_users":                       listAsSet(x.SourceUsers),
			"hip_profiles":                       listAsSet(x.HipProfiles),
			"destination_zones":                  listAsSet(x.DestinationZones),
			"destination_addresses":              listAsSet(x.DestinationAddresses),
			"negate_destination":                 x.NegateDestination,
			"applications":                       listAsSet(x.Applications),
			"services":                           listAsSet(x.Services),
			"categories":                         listAsSet(x.Categories),
			"action":                             x.Action,
			"log_setting":                        x.LogSetting,
			"log_start":                          x.LogStart,
			"log_end":                            x.LogEnd,
			"disabled":                           x.Disabled,
			"schedule":                           x.Schedule,
			"icmp_unreachable":                   x.IcmpUnreachable,
			"disable_server_response_inspection": x.DisableServerResponseInspection,
			"group":                              x.Group,
			"virus":                              x.Virus,
			"spyware":                            x.Spyware,
			"vulnerability":                      x.Vulnerability,
			"url_filtering":                      x.UrlFiltering,
			"file_blocking":                      x.FileBlocking,
			"wildfire_analysis":                  x.WildFireAnalysis,
			"data_filtering":                     x.DataFiltering,
			"group_tag":                          x.GroupTag,
			"source_devices":                     listAsSet(x.SourceDevices),
			"destination_devices":                listAsSet(x.DestinationDevices),
			"target":                             dumpTarget(x.Targets),
			"negate_target":                      x.NegateTarget,
			"audit_comment":                      "",
		})
	}

	if err := d.Set("rule", list); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildSecurityRuleGroupId(a, b, c string, d int, e string, f []security.Entry) string {
	list := make([]string, 0, len(f))
	for _, x := range f {
		list = append(list, x.Name)
	}
	return strings.Join([]string{a, b, c, strconv.Itoa(d), e, base64Encode(list)}, IdSeparator)
}

func parseSecurityRuleGroupId(v string) (string, string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[3])
	return t[0], t[1], t[2], move, t[4], base64Decode(t[5])
}

func parseSecurityPolicyId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildSecurityPolicyId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
