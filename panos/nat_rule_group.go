package panos

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceNatRules() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()
	s["rulebase"] = rulebaseSchema()

	return &schema.Resource{
		Read: dataSourceNatRulesRead,

		Schema: s,
	}
}

func dataSourceNatRulesRead(d *schema.ResourceData, meta interface{}) error {
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
		listing, err = con.Policies.Nat.GetList(vsys)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base}, IdSeparator)
		listing, err = con.Policies.Nat.GetList(dg, base)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceNatRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNatRuleRead,

		Schema: natRuleGroupSchema(false, nil),
	}
}

func dataSourceNatRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o nat.Entry

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
		o, err = con.Policies.Nat.Get(vsys, name)
	case *pango.Panorama:
		id = strings.Join([]string{dg, base, name}, IdSeparator)
		o, err = con.Policies.Nat.Get(dg, base, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveNatRules(d, []nat.Entry{o})

	return nil
}

// Resource.
func resourceNatRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateNatRuleGroup,
		Read:   readNatRuleGroup,
		Update: createUpdateNatRuleGroup,
		Delete: deleteNatRuleGroup,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: natRuleGroupSchema(true, []string{"device_group", "rulebase"}),
	}
}

func resourcePanoramaNatRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateNatRuleGroup,
		Read:   readNatRuleGroup,
		Update: createUpdateNatRuleGroup,
		Delete: deleteNatRuleGroup,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: natRuleGroupSchema(true, []string{"vsys"}),
	}
}

func createUpdateNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var vsys, dg, base string
	var prevNames []string

	move := movementAtoi(d.Get("position_keyword").(string))
	oRule := d.Get("position_reference").(string)
	rules, auditComments := loadNatRules(d)

	d.Set("position_keyword", movementItoa(move))
	d.Set("position_reference", oRule)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if d.Id() != "" {
			_, _, _, prevNames = parseNatRuleGroupId(d.Id())
		}
		vsys = d.Get("vsys").(string)
		d.Set("vsys", vsys)
		id = buildNatRuleGroupId(vsys, move, oRule, rules)
		err = con.Policies.Nat.ConfigureRules(vsys, rules, auditComments, false, move, oRule, prevNames)
	case *pango.Panorama:
		if d.Id() != "" {
			_, _, _, _, prevNames = parsePanoramaNatRuleGroupId(d.Id())
		}
		dg = d.Get("device_group").(string)
		base = d.Get("rulebase").(string)
		d.Set("device_group", dg)
		d.Set("rulebase", base)
		id = buildPanoramaNatRuleGroupId(dg, base, move, oRule, rules)
		err = con.Policies.Nat.ConfigureRules(dg, base, rules, auditComments, false, move, oRule, prevNames)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readNatRuleGroup(d, meta)
}

func readNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var names []string
	var vsys, dg, base, oRule string
	var listing []nat.Entry
	var move int

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, move, oRule, names = parseNatRuleGroupId(d.Id())
		listing, err = con.Policies.Nat.GetAll(vsys)
	case *pango.Panorama:
		dg, base, move, oRule, names = parsePanoramaNatRuleGroupId(d.Id())
		listing, err = con.Policies.Nat.GetAll(dg, base)
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

	dlist := make([]nat.Entry, 0, len(names))
	for i := 0; i+fIdx < len(listing) && i < len(names); i++ {
		if listing[i+fIdx].Name != names[i] {
			break
		}

		dlist = append(dlist, listing[i+fIdx])
	}

	if move == util.MoveBottom && dlist[len(dlist)-1].Name != listing[len(listing)-1].Name {
		d.Set("position_keyword", "")
	}
	saveNatRules(d, dlist)

	return nil
}

func deleteNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var vsys, dg, base string
	var names []string

	switch meta.(type) {
	case *pango.Firewall:
		vsys, _, _, names = parseNatRuleGroupId(d.Id())
	case *pango.Panorama:
		dg, base, _, _, names = parsePanoramaNatRuleGroupId(d.Id())
	}

	ilist := make([]interface{}, 0, len(names))
	for _, x := range names {
		ilist = append(ilist, x)
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Policies.Nat.Delete(vsys, ilist...)
	case *pango.Panorama:
		err = con.Policies.Nat.Delete(dg, base, ilist...)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func natRuleGroupSchema(isResource bool, rmKeys []string) map[string]*schema.Schema {
	var minRules int
	if isResource {
		minRules = 1
	}

	ans := map[string]*schema.Schema{
		"vsys":               vsysSchema("vsys1"),
		"device_group":       deviceGroupSchema(),
		"rulebase":           rulebaseSchema(),
		"position_keyword":   positionKeywordSchema(),
		"position_reference": positionReferenceSchema(),
		"rule": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: minRules,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"type": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  nat.TypeIpv4,
						ValidateFunc: validateStringIn(
							nat.TypeIpv4,
							nat.TypeNat64,
							nat.TypeNptv6,
						),
					},
					"tags": tagSchema(),
					"disabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"uuid":      uuidSchema(),
					"group_tag": groupTagSchema(),

					"original_packet": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"source_zones": {
									Type:     schema.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"destination_zone": {
									Type:     schema.TypeString,
									Required: true,
								},
								"destination_interface": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "any",
								},
								"service": {
									Type:     schema.TypeString,
									Optional: true,
									Default:  "any",
								},
								"source_addresses": {
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
							},
						},
					},

					"translated_packet": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"source": {
									Type:     schema.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"dynamic_ip_and_port": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.source.dynamic_ip",
														"rule.translated_packet.source.static_ip",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"translated_address": {
															Type:     schema.TypeList,
															Optional: true,
															//ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip_and_port.interface_address"},
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"translated_addresses": {
																		Type:     schema.TypeSet,
																		Optional: true,
																		Elem: &schema.Schema{
																			Type: schema.TypeString,
																		},
																	},
																},
															},
														},

														"interface_address": {
															Type:     schema.TypeList,
															Optional: true,
															//ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip_and_port.translated_address"},
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"interface": {
																		Type:     schema.TypeString,
																		Required: true,
																	},
																	"ip_address": {
																		Type:     schema.TypeString,
																		Optional: true,
																	},
																},
															},
														},
													},
												},
											},

											"dynamic_ip": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.source.dynamic_ip_and_port",
														"rule.translated_packet.source.static_ip",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"translated_addresses": {
															Type:     schema.TypeSet,
															Required: true,
															MinItems: 1,
															Elem: &schema.Schema{
																Type: schema.TypeString,
															},
														},
														"fallback": {
															Type:     schema.TypeList,
															Optional: true,
															MaxItems: 1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"translated_address": {
																		Type:     schema.TypeList,
																		Optional: true,
																		//ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip.fallback.interface_address"},
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"translated_addresses": {
																					Type:     schema.TypeSet,
																					Optional: true,
																					Elem: &schema.Schema{
																						Type: schema.TypeString,
																					},
																				},
																			},
																		},
																	},

																	"interface_address": {
																		Type:     schema.TypeList,
																		Optional: true,
																		//ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip.fallback.translated_address"},
																		MaxItems: 1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"interface": {
																					Type:     schema.TypeString,
																					Required: true,
																				},
																				"type": {
																					Type:         schema.TypeString,
																					Optional:     true,
																					Default:      nat.Ip,
																					ValidateFunc: validateStringIn(nat.Ip, nat.FloatingIp),
																				},
																				"ip_address": {
																					Type:     schema.TypeString,
																					Optional: true,
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

											"static_ip": {
												Type:     schema.TypeList,
												Optional: true,
												MaxItems: 1,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.source.dynamic_ip_and_port",
														"rule.translated_packet.source.dynamic_ip",
													},
												*/
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"translated_address": {
															Type:     schema.TypeString,
															Required: true,
														},
														"bi_directional": {
															Type:     schema.TypeBool,
															Optional: true,
														},
													},
												},
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
											"static": {
												Type:     schema.TypeList,
												Optional: true,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.destination.static_translation",
														"rule.translated_packet.destination.dynamic",
														"rule.translated_packet.destination.dynamic_translation",
													},
												*/
												MaxItems:   1,
												Deprecated: "Use 'static_translation' instead",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"address": {
															Type:     schema.TypeString,
															Required: true,
														},
														"port": {
															Type:     schema.TypeInt,
															Optional: true,
														},
													},
												},
											},
											"static_translation": {
												Type:     schema.TypeList,
												Optional: true,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.destination.static",
														"rule.translated_packet.destination.dynamic",
														"rule.translated_packet.destination.dynamic_translation",
													},
												*/
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"address": {
															Type:     schema.TypeString,
															Required: true,
														},
														"port": {
															Type:     schema.TypeInt,
															Optional: true,
														},
													},
												},
											},
											"dynamic": {
												Type:     schema.TypeList,
												Optional: true,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.destination.static",
														"rule.translated_packet.destination.static_translation",
														"rule.translated_packet.destination.dynamic_translation",
													},
												*/
												MaxItems:   1,
												Deprecated: "Use 'dynamic_translation' instead",
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"address": {
															Type:     schema.TypeString,
															Required: true,
														},
														"port": {
															Type:     schema.TypeInt,
															Optional: true,
														},
														"distribution": {
															Type:     schema.TypeString,
															Optional: true,
														},
													},
												},
											},
											"dynamic_translation": {
												Type:     schema.TypeList,
												Optional: true,
												/*
													ConflictsWith: []string{
														"rule.translated_packet.destination.static",
														"rule.translated_packet.destination.static_translation",
														"rule.translated_packet.destination.dynamic",
													},
												*/
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"address": {
															Type:     schema.TypeString,
															Required: true,
														},
														"port": {
															Type:     schema.TypeInt,
															Optional: true,
														},
														"distribution": {
															Type:     schema.TypeString,
															Optional: true,
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

func loadNatRules(d *schema.ResourceData) ([]nat.Entry, map[string]string) {
	auditComments := make(map[string]string)
	rlist := d.Get("rule").([]interface{})
	list := make([]nat.Entry, 0, len(rlist))
	for i := range rlist {
		x := rlist[i].(map[string]interface{})
		auditComments[x["name"].(string)] = x["audit_comment"].(string)
		list = append(list, loadNatEntry(x))
	}

	return list, auditComments
}

func loadNatEntry(b map[string]interface{}) nat.Entry {
	o := nat.Entry{
		Name:         b["name"].(string),
		GroupTag:     b["group_tag"].(string),
		Type:         b["type"].(string),
		Description:  b["description"].(string),
		Disabled:     b["disabled"].(bool),
		Tags:         asStringList(b["tags"].([]interface{})),
		Targets:      loadTarget(b["target"]),
		NegateTarget: b["negate_target"].(bool),
	}

	op := (b["original_packet"].([]interface{})[0]).(map[string]interface{})
	o.SourceZones = setAsList(op["source_zones"].(*schema.Set))
	o.DestinationZone = op["destination_zone"].(string)
	o.ToInterface = op["destination_interface"].(string)
	o.Service = op["service"].(string)
	o.SourceAddresses = setAsList(op["source_addresses"].(*schema.Set))
	o.DestinationAddresses = setAsList(op["destination_addresses"].(*schema.Set))

	tp := (b["translated_packet"].([]interface{})[0]).(map[string]interface{})

	src := asInterfaceMap(tp, "source")
	if diap := asInterfaceMap(src, "dynamic_ip_and_port"); len(diap) != 0 {
		o.SatType = nat.DynamicIpAndPort

		if s := asInterfaceMap(diap, "translated_address"); len(s) != 0 {
			o.SatAddressType = nat.TranslatedAddress

			o.SatTranslatedAddresses = setAsList(s["translated_addresses"].(*schema.Set))
		} else if s := asInterfaceMap(diap, "interface_address"); len(s) != 0 {
			o.SatAddressType = nat.InterfaceAddress

			o.SatInterface = s["interface"].(string)
			o.SatIpAddress = s["ip_address"].(string)
		}
	} else if di := asInterfaceMap(src, "dynamic_ip"); len(di) != 0 {
		o.SatType = nat.DynamicIp

		o.SatTranslatedAddresses = setAsList(di["translated_addresses"].(*schema.Set))
		if fb := asInterfaceMap(di, "fallback"); len(fb) != 0 {
			if s := asInterfaceMap(fb, "translated_address"); len(s) != 0 {
				o.SatFallbackType = nat.TranslatedAddress

				o.SatFallbackTranslatedAddresses = setAsList(s["translated_addresses"].(*schema.Set))
			} else if s := asInterfaceMap(fb, "interface_address"); len(s) != 0 {
				o.SatFallbackType = nat.InterfaceAddress

				o.SatFallbackInterface = s["interface"].(string)
				o.SatFallbackIpType = s["type"].(string)
				o.SatFallbackIpAddress = s["ip_address"].(string)
			}
		} else {
			o.SatFallbackType = nat.None
		}
	} else if s := asInterfaceMap(src, "static_ip"); len(s) != 0 {
		o.SatType = nat.StaticIp

		o.SatStaticTranslatedAddress = s["translated_address"].(string)
		o.SatStaticBiDirectional = s["bi_directional"].(bool)
	} else {
		o.SatType = nat.None
	}

	dst := asInterfaceMap(tp, "destination")
	if s := asInterfaceMap(dst, "static"); len(s) != 0 {
		o.DatType = nat.DatTypeStatic

		o.DatAddress = s["address"].(string)
		o.DatPort = s["port"].(int)
	} else if s := asInterfaceMap(dst, "static_translation"); len(s) != 0 {
		o.DatType = nat.DatTypeStatic

		o.DatAddress = s["address"].(string)
		o.DatPort = s["port"].(int)
	} else if s := asInterfaceMap(dst, "dynamic"); len(s) != 0 {
		o.DatType = nat.DatTypeDynamic

		o.DatAddress = s["address"].(string)
		o.DatPort = s["port"].(int)
		o.DatDynamicDistribution = s["distribution"].(string)
	} else if s := asInterfaceMap(dst, "dynamic_translation"); len(s) != 0 {
		o.DatType = nat.DatTypeDynamic

		o.DatAddress = s["address"].(string)
		o.DatPort = s["port"].(int)
		o.DatDynamicDistribution = s["distribution"].(string)
	}

	return o
}

func dumpNatRule(o nat.Entry) map[string]interface{} {
	m := map[string]interface{}{
		"name":          o.Name,
		"description":   o.Description,
		"type":          o.Type,
		"disabled":      o.Disabled,
		"tags":          o.Tags,
		"target":        dumpTarget(o.Targets),
		"negate_target": o.NegateTarget,
		"uuid":          o.Uuid,
		"group_tag":     o.GroupTag,
		"audit_comment": "",
	}

	op := map[string]interface{}{
		"source_zones":          listAsSet(o.SourceZones),
		"destination_zone":      o.DestinationZone,
		"destination_interface": o.ToInterface,
		"service":               o.Service,
		"source_addresses":      listAsSet(o.SourceAddresses),
		"destination_addresses": listAsSet(o.DestinationAddresses),
	}
	m["original_packet"] = []interface{}{op}

	tp := make(map[string]interface{})
	src := make(map[string]interface{})
	dst := make(map[string]interface{})
	switch o.SatType {
	case nat.DynamicIpAndPort:
		diap := make(map[string]interface{})
		switch o.SatAddressType {
		case nat.TranslatedAddress:
			diap["translated_address"] = []interface{}{
				map[string]interface{}{
					"translated_addresses": listAsSet(o.SatTranslatedAddresses),
				},
			}
		case nat.InterfaceAddress:
			diap["interface_address"] = []interface{}{
				map[string]interface{}{
					"interface":  o.SatInterface,
					"ip_address": o.SatIpAddress,
				},
			}
		}
		src["dynamic_ip_and_port"] = []interface{}{diap}
	case nat.DynamicIp:
		di := map[string]interface{}{
			"translated_addresses": listAsSet(o.SatTranslatedAddresses),
		}
		switch o.SatFallbackType {
		case nat.TranslatedAddress:
			di["fallback"] = []interface{}{
				map[string]interface{}{
					"translated_address": []interface{}{
						map[string]interface{}{
							"translated_addresses": listAsSet(o.SatFallbackTranslatedAddresses),
						},
					},
				},
			}
		case nat.InterfaceAddress:
			di["fallback"] = []interface{}{
				map[string]interface{}{
					"interface_address": []interface{}{
						map[string]interface{}{
							"interface":  o.SatFallbackInterface,
							"type":       o.SatFallbackIpType,
							"ip_address": o.SatFallbackIpAddress,
						},
					},
				},
			}
		case nat.None:
			di["fallback"] = []interface{}{}
		}
		src["dynamic_ip"] = []interface{}{di}
	case nat.StaticIp:
		src["static_ip"] = []interface{}{
			map[string]interface{}{
				"translated_address": o.SatStaticTranslatedAddress,
				"bi_directional":     o.SatStaticBiDirectional,
			},
		}
	}
	switch o.DatType {
	case nat.DatTypeStatic:
		dst["static_translation"] = []interface{}{
			map[string]interface{}{
				"address": o.DatAddress,
				"port":    o.DatPort,
			},
		}
	case nat.DatTypeDynamic:
		dst["dynamic_translation"] = []interface{}{
			map[string]interface{}{
				"address":      o.DatAddress,
				"port":         o.DatPort,
				"distribution": o.DatDynamicDistribution,
			},
		}
	}
	tp["source"] = []interface{}{src}
	tp["destination"] = []interface{}{dst}
	m["translated_packet"] = []interface{}{tp}

	return m
}

func saveNatRules(d *schema.ResourceData, rules []nat.Entry) {
	if len(rules) == 0 {
		d.Set("rule", nil)
		return
	}

	list := make([]interface{}, 0, len(rules))
	for _, x := range rules {
		list = append(list, dumpNatRule(x))
	}

	if err := d.Set("rule", list); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildNatRuleGroupId(a string, b int, c string, d []nat.Entry) string {
	names := make([]string, 0, len(d))
	for _, x := range d {
		names = append(names, x.Name)
	}
	return strings.Join([]string{a, strconv.Itoa(b), c, base64Encode(names)}, IdSeparator)
}

func buildPanoramaNatRuleGroupId(a, b string, c int, d string, e []nat.Entry) string {
	names := make([]string, 0, len(e))
	for _, x := range e {
		names = append(names, x.Name)
	}
	return strings.Join([]string{a, b, strconv.Itoa(c), d, base64Encode(names)}, IdSeparator)
}

func parseNatRuleGroupId(v string) (string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[1])
	return t[0], move, t[2], base64Decode(t[3])
}

func parsePanoramaNatRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	return t[0], t[1], move, t[3], base64Decode(t[4])
}
