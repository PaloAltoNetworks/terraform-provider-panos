package panos

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNatRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateNatRuleGroup,
		Read:   readNatRuleGroup,
		Update: createUpdateNatRuleGroup,
		Delete: deleteNatRuleGroup,

		Schema: natRuleGroupSchema(false),
	}
}

func natRuleGroupSchema(p bool) map[string]*schema.Schema {
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
					"type": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      nat.TypeIpv4,
						ValidateFunc: validateStringIn(nat.TypeIpv4, nat.TypeNat64, nat.TypeNptv6),
					},
					"tags": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"disabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},

					"original_packet": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"source_zones": {
									Type:     schema.TypeList,
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
												ConflictsWith: []string{
													"rule.translated_packet.source.dynamic_ip",
													"rule.translated_packet.source.static_ip",
												},
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"translated_address": {
															Type:          schema.TypeList,
															Optional:      true,
															ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip_and_port.interface_address"},
															MaxItems:      1,
															Elem: &schema.Resource{
																Schema: map[string]*schema.Schema{
																	"translated_addresses": {
																		Type:     schema.TypeList,
																		Optional: true,
																		Elem: &schema.Schema{
																			Type: schema.TypeString,
																		},
																	},
																},
															},
														},

														"interface_address": {
															Type:          schema.TypeList,
															Optional:      true,
															ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip_and_port.translated_address"},
															MaxItems:      1,
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
												ConflictsWith: []string{
													"rule.translated_packet.source.dynamic_ip_and_port",
													"rule.translated_packet.source.static_ip",
												},
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"translated_addresses": {
															Type:     schema.TypeList,
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
																		Type:          schema.TypeList,
																		Optional:      true,
																		ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip.fallback.interface_address"},
																		MaxItems:      1,
																		Elem: &schema.Resource{
																			Schema: map[string]*schema.Schema{
																				"translated_addresses": {
																					Type:     schema.TypeList,
																					Optional: true,
																					Elem: &schema.Schema{
																						Type: schema.TypeString,
																					},
																				},
																			},
																		},
																	},

																	"interface_address": {
																		Type:          schema.TypeList,
																		Optional:      true,
																		ConflictsWith: []string{"rule.translated_packet.source.dynamic_ip.fallback.translated_address"},
																		MaxItems:      1,
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
												ConflictsWith: []string{
													"rule.translated_packet.source.dynamic_ip_and_port",
													"rule.translated_packet.source.dynamic_ip",
												},
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
												Type:          schema.TypeList,
												Optional:      true,
												ConflictsWith: []string{"rule.translated_packet.destination.dynamic"},
												MaxItems:      1,
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
												Type:          schema.TypeList,
												Optional:      true,
												ConflictsWith: []string{"rule.translated_packet.destination.static"},
												MaxItems:      1,
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

func parseNatRuleGroup(d *schema.ResourceData) (string, string, int, []nat.Entry) {
	vsys := d.Get("vsys").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))

	rlist := d.Get("rule").([]interface{})
	list := make([]nat.Entry, 0, len(rlist))
	for i := range rlist {
		b := rlist[i].(map[string]interface{})
		o := loadNatEntry(b)

		list = append(list, o)
	}

	return vsys, oRule, move, list
}

func loadNatEntry(b map[string]interface{}) nat.Entry {
	o := nat.Entry{
		Name:        b["name"].(string),
		Type:        b["type"].(string),
		Description: b["description"].(string),
		Disabled:    b["disabled"].(bool),
		Tags:        asStringList(b["tags"].([]interface{})),
	}

	op := (b["original_packet"].([]interface{})[0]).(map[string]interface{})
	o.SourceZones = asStringList(op["source_zones"].([]interface{}))
	o.DestinationZone = op["destination_zone"].(string)
	o.ToInterface = op["destination_interface"].(string)
	o.Service = op["service"].(string)
	o.SourceAddresses = asStringList(op["source_addresses"].([]interface{}))
	o.DestinationAddresses = asStringList(op["destination_addresses"].([]interface{}))

	tp := (b["translated_packet"].([]interface{})[0]).(map[string]interface{})

	src := asInterfaceMap(tp, "source")
	if diap := asInterfaceMap(src, "dynamic_ip_and_port"); len(diap) != 0 {
		o.SatType = nat.DynamicIpAndPort

		if s := asInterfaceMap(diap, "translated_address"); len(s) != 0 {
			o.SatAddressType = nat.TranslatedAddress

			o.SatTranslatedAddresses = asStringList(s["translated_addresses"].([]interface{}))
		} else if s := asInterfaceMap(diap, "interface_address"); len(s) != 0 {
			o.SatAddressType = nat.InterfaceAddress

			o.SatInterface = s["interface"].(string)
			o.SatIpAddress = s["ip_address"].(string)
		}
	} else if di := asInterfaceMap(src, "dynamic_ip"); len(di) != 0 {
		o.SatType = nat.DynamicIp

		o.SatTranslatedAddresses = asStringList(di["translated_addresses"].([]interface{}))
		if fb := asInterfaceMap(di, "fallback"); len(fb) != 0 {
			if s := asInterfaceMap(fb, "translated_address"); len(s) != 0 {
				o.SatFallbackType = nat.TranslatedAddress

				o.SatFallbackTranslatedAddresses = asStringList(s["translated_addresses"].([]interface{}))
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
	} else if s := asInterfaceMap(dst, "dynamic"); len(s) != 0 {
		o.DatType = nat.DatTypeDynamic

		o.DatAddress = s["address"].(string)
		o.DatPort = s["port"].(int)
		o.DatDynamicDistribution = s["distribution"].(string)
	}

	return o
}

func dumpNatEntry(o nat.Entry) map[string]interface{} {
	m := map[string]interface{}{
		"name":        o.Name,
		"description": o.Description,
		"type":        o.Type,
		"disabled":    o.Disabled,
		"tags":        o.Tags,
	}

	op := map[string]interface{}{
		"source_zones":          o.SourceZones,
		"destination_zone":      o.DestinationZone,
		"destination_interface": o.ToInterface,
		"service":               o.Service,
		"source_addresses":      o.SourceAddresses,
		"destination_addresses": o.DestinationAddresses,
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
					"translated_addresses": o.SatTranslatedAddresses,
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
			"translated_addresses": o.SatTranslatedAddresses,
		}
		switch o.SatFallbackType {
		case nat.TranslatedAddress:
			di["fallback"] = []interface{}{
				map[string]interface{}{
					"translated_address": []interface{}{
						map[string]interface{}{
							"translated_addresses": o.SatFallbackTranslatedAddresses,
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
		dst["static"] = []interface{}{
			map[string]interface{}{
				"address": o.DatAddress,
				"port":    o.DatPort,
			},
		}
	case nat.DatTypeDynamic:
		dst["dynamic"] = []interface{}{
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

func parseNatRuleGroupId(v string) (string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[1])
	joined, _ := base64.StdEncoding.DecodeString(t[3])
	names := strings.Split(string(joined), "\n")
	return t[0], move, t[2], names
}

func buildNatRuleGroupId(a string, b int, c string, d []nat.Entry) string {
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

func createUpdateNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, oRule, move, list := parseNatRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = fw.Policies.Nat.Edit(vsys, list[0]); err != nil {
		return err
	}
	dl := make([]interface{}, len(list)-1)
	for i := 1; i < len(list); i++ {
		dl = append(dl, list[i])
	}
	_ = fw.Policies.Nat.Delete(vsys, dl...)
	if err = fw.Policies.Nat.Set(vsys, list[1:len(list)]...); err != nil {
		return err
	}
	if err = fw.Policies.Nat.MoveGroup(vsys, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildNatRuleGroupId(vsys, move, oRule, list))
	return readNatRuleGroup(d, meta)
}

func readNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, move, oRule, names := parseNatRuleGroupId(d.Id())

	rules, err := fw.Policies.Nat.GetList(vsys)
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
		o, err := fw.Policies.Nat.Get(vsys, names[i])
		if err != nil {
			if isObjectNotFound(err) {
				break
			}
			return err
		}
		m := dumpNatEntry(o)

		ilist = append(ilist, m)
	}

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, _, _, names := parseNatRuleGroupId(d.Id())

	ilist := make([]interface{}, len(names))
	for i := range names {
		ilist[i] = names[i]
	}

	if err := fw.Policies.Nat.Delete(vsys, ilist...); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
