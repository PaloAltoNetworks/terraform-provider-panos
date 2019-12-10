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

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaNatRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaNatRuleGroup,
		Read:   readPanoramaNatRuleGroup,
		Update: createUpdatePanoramaNatRuleGroup,
		Delete: deletePanoramaNatRuleGroup,

		Schema: natRuleGroupSchema(true),
	}
}

func parsePanoramaNatRuleGroup(d *schema.ResourceData) (string, string, int, string, []nat.Entry) {
	dg := d.Get("device_group").(string)
	rb := d.Get("rulebase").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))

	rlist := d.Get("rule").([]interface{})
	list := make([]nat.Entry, 0, len(rlist))
	for i := range rlist {
		b := rlist[i].(map[string]interface{})
		o := loadNatEntry(b)
		o.Targets = parseTarget(b["target"])
		o.NegateTarget = b["negate_target"].(bool)

		list = append(list, o)
	}

	return dg, rb, move, oRule, list
}

func parsePanoramaNatRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	joined, _ := base64.StdEncoding.DecodeString(t[4])
	names := strings.Split(string(joined), "\n")
	return t[0], t[1], move, t[3], names
}

func buildPanoramaNatRuleGroupId(a, b string, c int, d string, e []nat.Entry) string {
	var buf bytes.Buffer
	for i := range e {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(e[i].Name)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())

	return strings.Join([]string{a, b, strconv.Itoa(c), d, enc}, IdSeparator)
}

func createUpdatePanoramaNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, move, oRule, list := parsePanoramaNatRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = pano.Policies.Nat.Edit(dg, rb, list[0]); err != nil {
		return err
	}
	dl := make([]interface{}, len(list)-1)
	for i := 1; i < len(list); i++ {
		dl = append(dl, list[i])
	}
	_ = pano.Policies.Nat.Delete(dg, rb, dl...)
	if err = pano.Policies.Nat.Set(dg, rb, list[1:len(list)]...); err != nil {
		return err
	}
	if err = pano.Policies.Nat.MoveGroup(dg, rb, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildPanoramaNatRuleGroupId(dg, rb, move, oRule, list))
	return readPanoramaNatRuleGroup(d, meta)
}

func readPanoramaNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, rb, move, oRule, names := parsePanoramaNatRuleGroupId(d.Id())

	rules, err := pano.Policies.Nat.GetList(dg, rb)
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

	d.Set("device_group", dg)
	d.Set("rulebase", rb)
	d.Set("position_keyword", movementItoa(move))
	if groupPositionIsOk(move, fIdx, oIdx, rules, names) {
		d.Set("position_reference", oRule)
	} else {
		d.Set("position_reference", "(incorrect group positioning)")
	}

	schemaStyle := natRuleGroupSchemaStyle(d)

	ilist := make([]interface{}, 0, len(names))
	for i := 0; i+fIdx < len(rules) && i < len(names); i++ {
		if rules[i+fIdx] != names[i] {
			// Must be contiguous.
			break
		}
		o, err := pano.Policies.Nat.Get(dg, rb, names[i])
		if err != nil {
			if isObjectNotFound(err) {
				break
			}
			return err
		}
		m := dumpNatEntry(o, schemaStyle)
		m["target"] = buildTarget(o.Targets)
		m["negate_target"] = o.NegateTarget

		ilist = append(ilist, m)
	}

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}

	return nil
}

func deletePanoramaNatRuleGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, rb, _, _, names := parsePanoramaNatRuleGroupId(d.Id())

	ilist := make([]interface{}, len(names))
	for i := range names {
		ilist[i] = names[i]
	}

	if err := pano.Policies.Nat.Delete(dg, rb, ilist...); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
