package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/app/signature"
	"github.com/fpluchorg/pango/objs/app/signature/andcond"
	"github.com/fpluchorg/pango/objs/app/signature/orcond"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApplicationSignature() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateApplicationSignature,
		Read:   readApplicationSignature,
		Update: createUpdateApplicationSignature,
		Delete: deleteApplicationSignature,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationSignatureSchema(false),
	}
}

func applicationSignatureSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"application_object": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"scope": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      signature.ScopeTransaction,
			ValidateFunc: validateStringIn(signature.ScopeTransaction, signature.ScopeSession),
		},
		"ordered_match": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"and_condition": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"or_condition": {
						Type:     schema.TypeList,
						MinItems: 1,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"pattern_match": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									/*
										ConflictsWith: []string{
											"and_condition.or_condition.greater_than",
											"and_condition.or_condition.less_than",
											"and_condition.or_condition.equal_to",
										},
									*/
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"context": {
												Type:     schema.TypeString,
												Required: true,
											},
											"pattern": {
												Type:     schema.TypeString,
												Required: true,
											},
											"qualifiers": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"greater_than": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									/*
										ConflictsWith: []string{
											"and_condition.or_condition.pattern_match",
											"and_condition.or_condition.less_than",
											"and_condition.or_condition.equal_to",
										},
									*/
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"context": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:     schema.TypeString,
												Required: true,
											},
											"qualifiers": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"less_than": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									/*
										ConflictsWith: []string{
											"and_condition.or_condition.pattern_match",
											"and_condition.or_condition.greater_than",
											"and_condition.or_condition.equal_to",
										},
									*/
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"context": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:     schema.TypeString,
												Required: true,
											},
											"qualifiers": {
												Type:     schema.TypeMap,
												Optional: true,
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"equal_to": {
									Type:     schema.TypeList,
									MaxItems: 1,
									Optional: true,
									/*
										ConflictsWith: []string{
											"and_condition.or_condition.pattern_match",
											"and_condition.or_condition.greater_than",
											"and_condition.or_condition.less_than",
										},
									*/
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"context": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:     schema.TypeString,
												Required: true,
											},
											"position": {
												Type:     schema.TypeString,
												Optional: true,
											},
											"mask": {
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
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema("vsys1")
	}

	return ans
}

func parseApplicationSignature(d *schema.ResourceData) (string, string, signature.Entry, []andcond.Entry, map[string][]orcond.Entry) {
	vsys := d.Get("vsys").(string)
	ao := d.Get("application_object").(string)
	o, andList, orMap := loadApplicationSignature(d)

	return vsys, ao, o, andList, orMap
}

func loadApplicationSignature(d *schema.ResourceData) (signature.Entry, []andcond.Entry, map[string][]orcond.Entry) {
	o := signature.Entry{
		Name:      d.Get("name").(string),
		Comment:   d.Get("comment").(string),
		Scope:     d.Get("scope").(string),
		OrderFree: !d.Get("ordered_match").(bool),
	}

	al := d.Get("and_condition").([]interface{})
	if len(al) == 0 || (len(al) == 1 && al[0] == nil) {
		return o, nil, nil
	}

	andList := make([]andcond.Entry, 0, len(al))
	orMap := make(map[string][]orcond.Entry)
	for i := range al {
		andBlock := al[i].(map[string]interface{})
		andList = append(andList, andcond.Entry{
			Name: fmt.Sprintf("And Condition %d", i+1),
		})
		ol := andBlock["or_condition"].([]interface{})
		orList := make([]orcond.Entry, 0, len(ol))
		for j := range ol {
			orBlock := ol[j].(map[string]interface{})
			orList = append(orList, orcond.Entry{
				Name: fmt.Sprintf("Or Condition %d", j+1),
			})
			if x := asInterfaceMap(orBlock, "pattern_match"); len(x) != 0 {
				orList[j].Operator = orcond.OperatorPatternMatch
				orList[j].Context = x["context"].(string)
				orList[j].Pattern = x["pattern"].(string)
				qual := x["qualifiers"].(map[string]interface{})
				if len(qual) != 0 {
					orList[j].Qualifiers = make(map[string]string)
					for k, v := range qual {
						orList[j].Qualifiers[k] = v.(string)
					}
				}
			} else if x := asInterfaceMap(orBlock, "greater_than"); len(x) != 0 {
				orList[j].Operator = orcond.OperatorGreaterThan
				orList[j].Context = x["context"].(string)
				orList[j].Value = x["value"].(string)
				qual := x["qualifiers"].(map[string]interface{})
				if len(qual) != 0 {
					orList[j].Qualifiers = make(map[string]string)
					for k, v := range qual {
						orList[j].Qualifiers[k] = v.(string)
					}
				}
			} else if x := asInterfaceMap(orBlock, "less_than"); len(x) != 0 {
				orList[j].Operator = orcond.OperatorLessThan
				orList[j].Context = x["context"].(string)
				orList[j].Value = x["value"].(string)
				qual := x["qualifiers"].(map[string]interface{})
				if len(qual) != 0 {
					orList[j].Qualifiers = make(map[string]string)
					for k, v := range qual {
						orList[j].Qualifiers[k] = v.(string)
					}
				}
			} else if x := asInterfaceMap(orBlock, "equal_to"); len(x) != 0 {
				orList[j].Operator = orcond.OperatorEqualTo
				orList[j].Context = x["context"].(string)
				orList[j].Value = x["value"].(string)
				orList[j].Position = x["position"].(string)
				orList[j].Mask = x["mask"].(string)
			}
		}
		orMap[andList[i].Name] = orList
	}

	return o, andList, orMap
}

func saveApplicationSignature(d *schema.ResourceData, o signature.Entry, andList []andcond.Entry, orMap map[string][]orcond.Entry) {
	d.Set("name", o.Name)
	d.Set("comment", o.Comment)
	d.Set("scope", o.Scope)
	d.Set("ordered_match", !o.OrderFree)

	if len(andList) == 0 {
		d.Set("and_condition", nil)
		return
	}

	ac := make([]interface{}, 0, len(andList))
	for i := range andList {
		acEntry := make(map[string]interface{})
		acEntry["name"] = andList[i].Name

		orObjects := orMap[andList[i].Name]
		ocList := make([]interface{}, 0, len(orObjects))
		for j := range orObjects {
			oe := orObjects[j]
			orEntry := make(map[string]interface{})
			orEntry["name"] = oe.Name
			switch oe.Operator {
			case orcond.OperatorPatternMatch:
				orEntry["pattern_match"] = []interface{}{
					map[string]interface{}{
						"context":    oe.Context,
						"pattern":    oe.Pattern,
						"qualifiers": oe.Qualifiers,
					},
				}
			case orcond.OperatorGreaterThan:
				orEntry["greater_than"] = []interface{}{
					map[string]interface{}{
						"context":    oe.Context,
						"value":      oe.Value,
						"qualifiers": oe.Qualifiers,
					},
				}
			case orcond.OperatorLessThan:
				orEntry["less_than"] = []interface{}{
					map[string]interface{}{
						"context":    oe.Context,
						"value":      oe.Value,
						"qualifiers": oe.Qualifiers,
					},
				}
			case orcond.OperatorEqualTo:
				orEntry["equal_to"] = []interface{}{
					map[string]interface{}{
						"context":  oe.Context,
						"value":    oe.Value,
						"position": oe.Position,
						"mask":     oe.Mask,
					},
				}
			}
			ocList = append(ocList, orEntry)
		}
		acEntry["or_condition"] = ocList

		ac = append(ac, acEntry)
	}
	if err := d.Set("and_condition", ac); err != nil {
		log.Printf("[WARN] Error setting 'and_condition' for %q: %s", d.Id(), err)
	}
}

func parseApplicationSignatureId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildApplicationSignatureId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createUpdateApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, app, o, andList, orMap := parseApplicationSignature(d)

	if err := fw.Objects.AppSignature.Edit(vsys, app, o); err != nil {
		return err
	}

	if err := fw.Objects.AppSigAndCond.Set(vsys, app, o.Name, andList...); err != nil {
		return err
	}

	for i := range andList {
		orList := orMap[andList[i].Name]
		if err := fw.Objects.AppSigOrCond.Set(vsys, app, o.Name, andList[i].Name, orList...); err != nil {
			return err
		}
	}

	d.SetId(buildApplicationSignatureId(vsys, app, o.Name))
	return readApplicationSignature(d, meta)
}

func readApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, app, name := parseApplicationSignatureId(d.Id())

	o, err := fw.Objects.AppSignature.Get(vsys, app, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	andNames, err := fw.Objects.AppSigAndCond.GetList(vsys, app, name)
	andList := make([]andcond.Entry, 0, len(andNames))
	orMap := make(map[string][]orcond.Entry)
	for _, andName := range andNames {
		andEntry, err := fw.Objects.AppSigAndCond.Get(vsys, app, name, andName)
		if err != nil {
			return err
		}
		andList = append(andList, andEntry)
		orNames, err := fw.Objects.AppSigOrCond.GetList(vsys, app, name, andName)
		orList := make([]orcond.Entry, 0, len(orNames))
		if err != nil {
			return err
		}
		for _, orName := range orNames {
			orEntry, err := fw.Objects.AppSigOrCond.Get(vsys, app, name, andName, orName)
			if err != nil {
				return err
			}
			orList = append(orList, orEntry)
		}
		orMap[andEntry.Name] = orList
	}

	d.Set("vsys", vsys)
	d.Set("application_object", app)
	saveApplicationSignature(d, o, andList, orMap)

	return nil
}

func deleteApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, app, name := parseApplicationSignatureId(d.Id())

	err := fw.Objects.AppSignature.Delete(vsys, app, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
