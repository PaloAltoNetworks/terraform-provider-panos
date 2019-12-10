package panos

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/pbf"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePbfRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePbfRuleGroup,
		Read:   readPbfRuleGroup,
		Update: createUpdatePbfRuleGroup,
		Delete: deletePbfRuleGroup,

		Schema: pbfRuleGroupSchema(false),
	}
}

func pbfRuleGroupSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
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
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},

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
				},
			},
		},
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
		ans["rulebase"] = rulebaseSchema()

		r := ans["rule"].Elem.(*schema.Resource)
		r.Schema["target"] = targetSchema()
		r.Schema["negate_target"] = negateTargetSchema()
	} else {
		ans["vsys"] = vsysSchema()
	}

	return ans
}

func parsePbfRuleGroup(d *schema.ResourceData) (string, string, int, []pbf.Entry) {
	vsys := d.Get("vsys").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))

	rlist := d.Get("rule").([]interface{})
	list := make([]pbf.Entry, 0, len(rlist))
	for i := range rlist {
		b := rlist[i].(map[string]interface{})
		o := loadPbfEntry(b, false)

		list = append(list, o)
	}

	return vsys, oRule, move, list
}

func loadPbfEntry(b map[string]interface{}, p bool) pbf.Entry {
	o := pbf.Entry{
		Name:                      b["name"].(string),
		Description:               b["description"].(string),
		Tags:                      asStringList(b["tags"].([]interface{})),
		ActiveActiveDeviceBinding: b["active_active_device_binding"].(string),
		Schedule:                  b["schedule"].(string),
		Disabled:                  b["disabled"].(bool),
		Uuid:                      b["uuid"].(string),
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

	if p {
		o.Targets = parseTarget(b["target"])
		o.NegateTarget = b["negate_target"].(bool)
	}

	return o
}

func dumpPbfEntry(o pbf.Entry, p bool) map[string]interface{} {
	m := map[string]interface{}{
		"name":                         o.Name,
		"description":                  o.Description,
		"tags":                         o.Tags,
		"active_active_device_binding": o.ActiveActiveDeviceBinding,
		"schedule":                     o.Schedule,
		"disabled":                     o.Disabled,
		"uuid":                         o.Uuid,
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

	if p {
		m["target"] = buildTarget(o.Targets)
		m["negate_target"] = o.NegateTarget
	}

	return m
}

func parsePbfRuleGroupId(v string) (string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[1])
	joined, _ := base64.StdEncoding.DecodeString(t[3])
	names := strings.Split(string(joined), "\n")
	return t[0], move, t[2], names
}

func buildPbfRuleGroupId(a string, b int, c string, d []pbf.Entry) string {
	var buf bytes.Buffer
	for i := range d {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(d[i].Name)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())

	return strings.Join([]string{a, strconv.Itoa(b), c, enc}, IdSeparator)
}

func createUpdatePbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, oRule, move, list := parsePbfRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = fw.Policies.PolicyBasedForwarding.Edit(vsys, list[0]); err != nil {
		return err
	}
	dl := make([]interface{}, len(list)-1)
	for i := 1; i < len(list); i++ {
		dl = append(dl, list[i])
	}
	_ = fw.Policies.PolicyBasedForwarding.Delete(vsys, dl...)
	if err = fw.Policies.PolicyBasedForwarding.Set(vsys, list[1:len(list)]...); err != nil {
		return err
	}
	if err = fw.Policies.PolicyBasedForwarding.MoveGroup(vsys, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildPbfRuleGroupId(vsys, move, oRule, list))
	return readPbfRuleGroup(d, meta)
}

func readPbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, move, oRule, names := parsePbfRuleGroupId(d.Id())

	rules, err := fw.Policies.PolicyBasedForwarding.GetList(vsys)
	if err != nil {
		return err
	}

	fIdx, oIdx := -1, -1
	for i := range rules {
		if rules[i] == names[0] {
			fIdx = i
		} else if rules[i] == oRule {
			oIdx = i
		}
		if fIdx != -1 && oIdx != -1 {
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
	}

	d.Set("vsys", vsys)
	d.Set("position_keyword", movementItoa(move))
	if groupPositionIsOk(move, fIdx, oIdx, rules, names) {
		d.Set("position_reference", oRule)
	} else {
		d.Set("position_reference", "(incorrect group positioning)")
	}

	ilist := make([]interface{}, 0, len(names))
	for i := 0; i+fIdx < len(rules) && i < len(names); i++ {
		if rules[i+fIdx] != names[i] {
			// Must be contiguous.
			break
		}
		o, err := fw.Policies.PolicyBasedForwarding.Get(vsys, names[i])
		if err != nil {
			if isObjectNotFound(err) {
				break
			}
			return err
		}
		m := dumpPbfEntry(o, false)

		ilist = append(ilist, m)
	}

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}

	return nil
}

func deletePbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, _, _, names := parsePbfRuleGroupId(d.Id())

	ilist := make([]interface{}, len(names))
	for i := range names {
		ilist[i] = names[i]
	}

	if err := fw.Policies.PolicyBasedForwarding.Delete(vsys, ilist...); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
