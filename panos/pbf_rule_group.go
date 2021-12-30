package panos

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/pbf"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourcePbfRules() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()
	s["rulebase"] = rulebaseSchema()

	return &schema.Resource{
		Read: dataSourcePbfRulesRead,

		Schema: s,
	}
}

func dataSourcePbfRulesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)

	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	d.Set("rulebase", base)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vsys
		listing, err = con.Policies.PolicyBasedForwarding.GetList(vsys)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base}, IdSeparator)
		listing, err = con.Policies.PolicyBasedForwarding.GetList(dg, base)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourcePbfRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePbfRuleRead,

		Schema: pbfRuleGroupSchema(false, nil),
	}
}

func dataSourcePbfRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o pbf.Entry

	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	name := d.Get("name").(string)

	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	d.Set("rulebase", base)
	d.Set("name", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = strings.Join([]string{vsys, name}, IdSeparator)
		o, err = con.Policies.PolicyBasedForwarding.Get(vsys, name)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base, name}, IdSeparator)
		o, err = con.Policies.PolicyBasedForwarding.Get(dg, base, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	savePbfRules(d, []pbf.Entry{o})

	return nil
}

// Resource.
func resourcePbfRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePbfRuleGroup,
		Read:   readPbfRuleGroup,
		Update: createUpdatePbfRuleGroup,
		Delete: deletePbfRuleGroup,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: pbfRuleGroupSchema(true, []string{"device_group", "rulebase"}),
	}
}

func resourcePanoramaPbfRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePbfRuleGroup,
		Read:   readPbfRuleGroup,
		Update: createUpdatePbfRuleGroup,
		Delete: deletePbfRuleGroup,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: pbfRuleGroupSchema(true, []string{"vsys"}),
	}
}

func createUpdatePbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var vsys, dg, base string
	var prevNames []string

	move := movementAtoi(d.Get("position_keyword").(string))
	oRule := d.Get("position_reference").(string)
	rules, auditComments := loadPbfRules(d)

	d.Set("position_keyword", movementItoa(move))
	d.Set("position_reference", oRule)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if d.Id() != "" {
			_, _, _, prevNames = parsePbfRuleGroupId(d.Id())
		}
		vsys = d.Get("vsys").(string)
		d.Set("vsys", vsys)
		id = buildPbfRuleGroupId(vsys, move, oRule, rules)
		err = con.Policies.PolicyBasedForwarding.ConfigureRules(vsys, rules, auditComments, false, move, oRule, prevNames)
	case *pango.Panorama:
		if d.Id() != "" {
			_, _, _, _, prevNames = parsePanoramaPbfRuleGroupId(d.Id())
		}
		dg = d.Get("device_group").(string)
		base = d.Get("rulebase").(string)
		d.Set("device_group", dg)
		d.Set("rulebase", base)
		id = buildPanoramaPbfRuleGroupId(dg, base, move, oRule, rules)
		err = con.Policies.PolicyBasedForwarding.ConfigureRules(dg, base, rules, auditComments, false, move, oRule, prevNames)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readPbfRuleGroup(d, meta)
}

func readPbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var names []string
	var vsys, dg, base, oRule string
	var listing []pbf.Entry
	var move int

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, move, oRule, names = parsePbfRuleGroupId(d.Id())
		listing, err = con.Policies.PolicyBasedForwarding.GetAll(vsys)
	case *pango.Panorama:
		dg, base, move, oRule, names = parsePanoramaPbfRuleGroupId(d.Id())
		listing, err = con.Policies.PolicyBasedForwarding.GetAll(dg, base)
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

	dlist := make([]pbf.Entry, 0, len(names))
	for i := 0; i+fIdx < len(listing) && i < len(names); i++ {
		if listing[i+fIdx].Name != names[i] {
			break
		}

		dlist = append(dlist, listing[i+fIdx])
	}

	if move == util.MoveBottom && dlist[len(dlist)-1].Name != listing[len(listing)-1].Name {
		d.Set("position_keyword", "")
	}
	savePbfRules(d, dlist)

	return nil
}

func deletePbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var vsys, dg, base string
	var names []string

	switch meta.(type) {
	case *pango.Firewall:
		vsys, _, _, names = parsePbfRuleGroupId(d.Id())
	case *pango.Panorama:
		dg, base, _, _, names = parsePanoramaPbfRuleGroupId(d.Id())
	}

	ilist := make([]interface{}, 0, len(names))
	for _, x := range names {
		ilist = append(ilist, x)
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.PolicyBasedForwarding.Delete(vsys, ilist...)
	case *pango.Panorama:
		err = con.Policies.PolicyBasedForwarding.Delete(dg, base, ilist...)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func pbfRuleGroupSchema(isResource bool, rmList []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys":               vsysSchema("vsys1"),
		"device_group":       deviceGroupSchema(),
		"rulebase":           rulebaseSchema(),
		"position_keyword":   positionKeywordSchema(),
		"position_reference": positionReferenceSchema(),
		"rule": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"tags": tagSchema(),
					"active_active_device_binding": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"schedule": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"disabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"uuid":          uuidSchema(),
					"group_tag":     groupTagSchema(),
					"target":        targetSchema(false),
					"negate_target": negateTargetSchema(),

					"source": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"zones": {
									Type:     schema.TypeSet,
									Optional: true,
									MinItems: 1,
									ConflictsWith: []string{
										"rule.source.interfaces",
									},
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"interfaces": {
									Type:     schema.TypeSet,
									Optional: true,
									MinItems: 1,
									ConflictsWith: []string{
										"rule.source.zones",
									},
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"addresses": {
									Type:     schema.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"users": {
									Type:     schema.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"negate": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},

					"destination": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"addresses": {
									Type:     schema.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
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
								"negate": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},

					"forwarding": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Type:         schema.TypeString,
									Optional:     true,
									Default:      pbf.ActionForward,
									ValidateFunc: validateStringIn(pbf.ActionForward, pbf.ActionVsysForward, pbf.ActionDiscard, pbf.ActionNoPbf),
								},
								"vsys": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"egress_interface": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"next_hop_type": {
									Type:         schema.TypeString,
									Optional:     true,
									ValidateFunc: validateStringIn(pbf.ForwardNextHopTypeIpAddress, pbf.ForwardNextHopTypeFqdn),
								},
								"next_hop_value": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"monitor": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"profile": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"ip_address": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"disable_if_unreachable": {
												Type:     schema.TypeBool,
												Optional: true,
											},
										},
									},
								},
								"symmetric_return": {
									Type:     schema.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"enable": {
												Type:     schema.TypeBool,
												Optional: true,
											},
											"addresses": {
												Type:     schema.TypeList,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
					"audit_comment": auditCommentSchema(),
				},
			},
		},
	}

	for _, key := range rmList {
		delete(ans, key)
	}

	if !isResource {
		delete(ans, "position_keyword")
		delete(ans, "position_reference")
		ans["name"] = &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		}

		computed(ans, "", []string{"device_group", "vsys", "rulebase", "name"})
	}

	return ans
}

func loadPbfRules(d *schema.ResourceData) ([]pbf.Entry, map[string]string) {
	auditComments := make(map[string]string)
	rlist := d.Get("rule").([]interface{})
	list := make([]pbf.Entry, 0, len(rlist))
	for i := range rlist {
		b := rlist[i].(map[string]interface{})
		auditComments[b["name"].(string)] = b["audit_comment"].(string)
		o := pbf.Entry{
			Name:                      b["name"].(string),
			Description:               b["description"].(string),
			Tags:                      asStringList(b["tags"].([]interface{})),
			ActiveActiveDeviceBinding: b["active_active_device_binding"].(string),
			Schedule:                  b["schedule"].(string),
			Disabled:                  b["disabled"].(bool),
			GroupTag:                  b["group_tag"].(string),
			Targets:                   loadTarget(b["target"]),
			NegateTarget:              b["negate_target"].(bool),
		}

		src := asInterfaceMap(b, "source")
		if zones := setAsList(src["zones"].(*schema.Set)); len(zones) > 0 {
			o.FromType = pbf.FromTypeZone
			o.FromValues = zones
		} else if ifaces := setAsList(src["interfaces"].(*schema.Set)); len(ifaces) > 0 {
			o.FromType = pbf.FromTypeInterface
			o.FromValues = ifaces
		}
		o.SourceAddresses = setAsList(src["addresses"].(*schema.Set))
		o.SourceUsers = setAsList(src["users"].(*schema.Set))
		o.NegateSource = src["negate"].(bool)

		dst := asInterfaceMap(b, "destination")
		o.DestinationAddresses = setAsList(dst["addresses"].(*schema.Set))
		o.Applications = setAsList(dst["applications"].(*schema.Set))
		o.Services = setAsList(dst["services"].(*schema.Set))
		o.NegateDestination = dst["negate"].(bool)

		fwd := asInterfaceMap(b, "forwarding")
		o.Action = fwd["action"].(string)
		o.ForwardVsys = fwd["vsys"].(string)
		o.ForwardEgressInterface = fwd["egress_interface"].(string)
		o.ForwardNextHopType = fwd["next_hop_type"].(string)
		o.ForwardNextHopValue = fwd["next_hop_value"].(string)
		if mon := asInterfaceMap(fwd, "monitor"); len(mon) > 0 {
			o.ForwardMonitorProfile = mon["profile"].(string)
			o.ForwardMonitorIpAddress = mon["ip_address"].(string)
			o.ForwardMonitorDisableIfUnreachable = mon["disable_if_unreachable"].(bool)
		}
		if sym := asInterfaceMap(fwd, "symmetric_return"); len(sym) > 0 {
			o.EnableEnforceSymmetricReturn = sym["enable"].(bool)
			o.SymmetricReturnAddresses = asStringList(sym["addresses"].([]interface{}))
		}

		list = append(list, o)
	}

	return list, auditComments
}

func dumpPbfRule(o pbf.Entry) map[string]interface{} {
	m := map[string]interface{}{
		"name":                         o.Name,
		"description":                  o.Description,
		"tags":                         o.Tags,
		"active_active_device_binding": o.ActiveActiveDeviceBinding,
		"schedule":                     o.Schedule,
		"disabled":                     o.Disabled,
		"uuid":                         o.Uuid,
		"group_tag":                    o.GroupTag,
		"target":                       dumpTarget(o.Targets),
		"negate_target":                o.NegateTarget,
		"audit_comment":                "",
	}

	src := map[string]interface{}{
		"addresses": listAsSet(o.SourceAddresses),
		"users":     listAsSet(o.SourceUsers),
		"negate":    o.NegateSource,
	}
	switch o.FromType {
	case pbf.FromTypeZone:
		src["zones"] = listAsSet(o.FromValues)
	case pbf.FromTypeInterface:
		src["interfaces"] = listAsSet(o.FromValues)
	}
	m["source"] = []interface{}{src}

	dst := map[string]interface{}{
		"addresses":    listAsSet(o.DestinationAddresses),
		"applications": listAsSet(o.Applications),
		"services":     listAsSet(o.Services),
		"negate":       o.NegateDestination,
	}
	m["destination"] = []interface{}{dst}

	fwd := map[string]interface{}{
		"action":           o.Action,
		"vsys":             o.ForwardVsys,
		"egress_interface": o.ForwardEgressInterface,
		"next_hop_type":    o.ForwardNextHopType,
		"next_hop_value":   o.ForwardNextHopValue,
	}
	if o.ForwardMonitorProfile != "" || o.ForwardMonitorIpAddress != "" || o.ForwardMonitorDisableIfUnreachable {
		mon := map[string]interface{}{
			"profile":                o.ForwardMonitorProfile,
			"ip_address":             o.ForwardMonitorIpAddress,
			"disable_if_unreachable": o.ForwardMonitorDisableIfUnreachable,
		}
		fwd["monitor"] = []interface{}{mon}
	} else {
		fwd["monitor"] = nil
	}
	if o.EnableEnforceSymmetricReturn || len(o.SymmetricReturnAddresses) > 0 {
		sym := map[string]interface{}{
			"enable":    o.EnableEnforceSymmetricReturn,
			"addresses": o.SymmetricReturnAddresses,
		}
		fwd["symmetric_return"] = []interface{}{sym}
	} else {
		fwd["symmetric_return"] = nil
	}
	m["forwarding"] = []interface{}{fwd}

	return m
}

func savePbfRules(d *schema.ResourceData, rules []pbf.Entry) {
	if len(rules) == 0 {
		d.Set("rule", nil)
		return
	}

	list := make([]interface{}, 0, len(rules))
	for _, x := range rules {
		list = append(list, dumpPbfRule(x))
	}

	if err := d.Set("rule", list); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildPbfRuleGroupId(a string, b int, c string, d []pbf.Entry) string {
	names := make([]string, 0, len(d))
	for _, x := range d {
		names = append(names, x.Name)
	}
	return strings.Join([]string{a, strconv.Itoa(b), c, base64Encode(names)}, IdSeparator)
}

func buildPanoramaPbfRuleGroupId(a, b string, c int, d string, e []pbf.Entry) string {
	names := make([]string, 0, len(e))
	for _, x := range e {
		names = append(names, x.Name)
	}
	return strings.Join([]string{a, b, strconv.Itoa(c), d, base64Encode(names)}, IdSeparator)
}

func parsePbfRuleGroupId(v string) (string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[1])
	return t[0], move, t[2], base64Decode(t[3])
}

func parsePanoramaPbfRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	return t[0], t[1], move, t[3], base64Decode(t[4])
}
