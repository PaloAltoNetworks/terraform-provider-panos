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

func resourcePanoramaPbfRuleGroup() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaPbfRuleGroup,
		Read:   readPanoramaPbfRuleGroup,
		Update: createUpdatePanoramaPbfRuleGroup,
		Delete: deletePanoramaPbfRuleGroup,

		Schema: pbfRuleGroupSchema(true),
	}
}

func parsePanoramaPbfRuleGroup(d *schema.ResourceData) (string, string, string, int, []pbf.Entry) {
	dg := d.Get("device_group").(string)
	base := d.Get("rulebase").(string)
	oRule := d.Get("position_reference").(string)
	move := movementAtoi(d.Get("position_keyword").(string))

	rlist := d.Get("rule").([]interface{})
	list := make([]pbf.Entry, 0, len(rlist))
	for i := range rlist {
		b := rlist[i].(map[string]interface{})
		o := loadPbfEntry(b, true)

		list = append(list, o)
	}

	return dg, base, oRule, move, list
}

func parsePanoramaPbfRuleGroupId(v string) (string, string, int, string, []string) {
	t := strings.Split(v, IdSeparator)
	move, _ := strconv.Atoi(t[2])
	joined, _ := base64.StdEncoding.DecodeString(t[4])
	names := strings.Split(string(joined), "\n")
	return t[0], t[1], move, t[3], names
}

func buildPanoramaPbfRuleGroupId(a, b string, c int, d string, e []pbf.Entry) string {
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

func createUpdatePanoramaPbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, base, oRule, move, list := parsePanoramaPbfRuleGroup(d)

	if !movementIsRelative(move) && oRule != "" {
		return fmt.Errorf("'position_reference' must be empty for non-relative movement")
	}
	if err = pano.Policies.PolicyBasedForwarding.Edit(dg, base, list[0]); err != nil {
		return err
	}
	dl := make([]interface{}, len(list)-1)
	for i := 1; i < len(list); i++ {
		dl = append(dl, list[i])
	}
	_ = pano.Policies.PolicyBasedForwarding.Delete(dg, base, dl...)
	if err = pano.Policies.PolicyBasedForwarding.Set(dg, base, list[1:len(list)]...); err != nil {
		return err
	}
	if err = pano.Policies.PolicyBasedForwarding.MoveGroup(dg, base, move, oRule, list...); err != nil {
		return err
	}

	d.SetId(buildPanoramaPbfRuleGroupId(dg, base, move, oRule, list))
	return readPanoramaPbfRuleGroup(d, meta)
}

func readPanoramaPbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, base, move, oRule, names := parsePanoramaPbfRuleGroupId(d.Id())

	rules, err := pano.Policies.PolicyBasedForwarding.GetList(dg, base)
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
	d.Set("rulebase", base)
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
		o, err := pano.Policies.PolicyBasedForwarding.Get(dg, base, names[i])
		if err != nil {
			if isObjectNotFound(err) {
				break
			}
			return err
		}
		m := dumpPbfEntry(o, true)

		ilist = append(ilist, m)
	}

	if err = d.Set("rule", ilist); err != nil {
		log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
	}

	return nil
}

func deletePanoramaPbfRuleGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, base, _, _, names := parsePanoramaPbfRuleGroupId(d.Id())

	ilist := make([]interface{}, len(names))
	for i := range names {
		ilist[i] = names[i]
	}

	if err := pano.Policies.PolicyBasedForwarding.Delete(dg, base, ilist...); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
