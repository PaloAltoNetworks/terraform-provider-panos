package panos

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const IdSeparator string = ":"
const WrongPanosWithoutAltError string = "This is a %s resource, but encountered a %s system"
const WrongPanosWithAltError string = "This is a %s resource, but encountered a %s system - Please use %s instead"

func getMovementMap() map[int]string {
	return map[int]string{
		util.MoveSkip:           "",
		util.MoveBefore:         "before",
		util.MoveDirectlyBefore: "directly before",
		util.MoveAfter:          "after",
		util.MoveDirectlyAfter:  "directly after",
		util.MoveTop:            "top",
		util.MoveBottom:         "bottom",
	}
}

func movementKeywords() []string {
	mm := getMovementMap()
	ans := make([]string, 0, len(mm))

	for _, v := range mm {
		ans = append(ans, v)
	}

	return ans
}

func movementItoa(v int) string {
	mm := getMovementMap()
	return mm[v]
}

func movementAtoi(v string) int {
	mm := getMovementMap()

	for k, s := range mm {
		if s == v {
			return k
		}
	}

	return util.MoveSkip
}

func movementIsRelative(v int) bool {
	switch v {
	case util.MoveBefore, util.MoveDirectlyBefore, util.MoveAfter, util.MoveDirectlyAfter:
		return true
	default:
		return false
	}
}

func groupIndexes(rules, names []string, move int, oRule string) (int, int, error) {
	var err error
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

	if oIdx == -1 && movementIsRelative(move) {
		err = fmt.Errorf("Can't verify positioning as position_reference %q is not present", oRule)
	}

	return fIdx, oIdx, err
}

func groupPositionIsOk(move, fIdx, oIdx int, list, grp []string) bool {
	switch move {
	case util.MoveSkip:
		return true
	case util.MoveTop:
		if list[0] == grp[0] {
			return true
		}
	case util.MoveBottom:
		if len(grp) <= len(list) && list[len(list)-len(grp)] == grp[0] {
			return true
		}
	case util.MoveBefore:
		if fIdx < oIdx {
			return true
		}
	case util.MoveDirectlyBefore:
		if fIdx+1 == oIdx {
			return true
		}
	case util.MoveAfter:
		if fIdx > oIdx {
			return true
		}
	case util.MoveDirectlyAfter:
		if fIdx == oIdx+1 {
			return true
		}
	}

	return false
}

func asStringList(v []interface{}) []string {
	if len(v) == 0 {
		return nil
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case string:
			ans[i] = x
		case nil:
			ans[i] = ""
		}
	}

	return ans
}

func setAsList(d *schema.Set) []string {
	list := d.List()
	ans := make([]string, len(list))
	for i := range list {
		ans[i] = list[i].(string)
	}

	if len(list) == 0 {
		return nil
	} else {
		return ans
	}
}

func listAsSet(list []string) *schema.Set {
	items := make([]interface{}, len(list))
	for i := range list {
		items[i] = list[i]
	}

	return schema.NewSet(schema.HashString, items)
}

func isObjectNotFound(e error) bool {
	e2, ok := e.(errors.Panos)
	if ok && e2.ObjectNotFound() {
		return true
	}

	return false
}

func asInterfaceMap(m map[string]interface{}, k string) map[string]interface{} {
	if _, ok := m[k]; ok {
		v1, ok := m[k].([]interface{})
		if !ok || len(v1) == 0 {
			return map[string]interface{}{}
		}

		v2, ok := v1[0].(map[string]interface{})
		if !ok || v2 == nil {
			return map[string]interface{}{}
		}

		return v2
	}

	return map[string]interface{}{}
}

func configFolder(d *schema.ResourceData, key string) map[string]interface{} {
	if clist, ok := d.Get(key).([]interface{}); ok {
		if clist != nil && len(clist) == 1 {
			ans, ok := clist[0].(map[string]interface{})
			if ok && ans != nil {
				return ans
			}
		}
	}

	return nil
}

func loadTarget(v interface{}) map[string][]string {
	if v == nil {
		return nil
	}

	ans := make(map[string][]string)
	sl := v.(*schema.Set).List()

	if len(sl) == 0 {
		return nil
	}

	for i := range sl {
		dev := sl[i].(map[string]interface{})
		key := dev["serial"].(string)
		value := setAsList(dev["vsys_list"].(*schema.Set))
		ans[key] = value
	}

	return ans
}

func dumpTarget(m map[string][]string) *schema.Set {
	var items []interface{}

	if len(m) > 0 {
		items = make([]interface{}, 0, len(m))
		for key := range m {
			items = append(items, map[string]interface{}{
				"serial":    key,
				"vsys_list": listAsSet(m[key]),
			})
		}
	}

	return schema.NewSet(
		schema.HashResource(
			targetSchema(false).Elem.(*schema.Resource),
		),
		items,
	)
}

func firewall(meta interface{}, alt string) (*pango.Firewall, error) {
	if fw, ok := meta.(*pango.Firewall); ok {
		return fw, nil
	}

	if alt != "" {
		return nil, fmt.Errorf(WrongPanosWithAltError, "firewall", "Panorama", alt)
	}

	return nil, fmt.Errorf(WrongPanosWithoutAltError, "firewall", "Panorama")
}

func panorama(meta interface{}, alt string) (*pango.Panorama, error) {
	if p, ok := meta.(*pango.Panorama); ok {
		return p, nil
	}

	if alt != "" {
		return nil, fmt.Errorf(WrongPanosWithAltError, "Panorama", "firewall", alt)
	}

	return nil, fmt.Errorf(WrongPanosWithoutAltError, "Panorama", "firewall")
}

func computed(sm map[string]*schema.Schema, parent string, omits []string) {
	for key, s := range sm {
		stop := false
		for _, o := range omits {
			if parent == "" {
				if o == key {
					stop = true
					break
				}
			} else if o == parent+"."+key {
				stop = true
				break
			}
		}
		if stop {
			continue
		}
		s.Computed = true
		s.Required = false
		s.Optional = false
		s.MinItems = 0
		s.MaxItems = 0
		s.Default = nil
		s.DiffSuppressFunc = nil
		s.DefaultFunc = nil
		s.ConflictsWith = nil
		s.ExactlyOneOf = nil
		s.AtLeastOneOf = nil
		s.ValidateFunc = nil
		//s.RequiredWith = nil
		if s.Type == schema.TypeList || s.Type == schema.TypeSet {
			switch et := s.Elem.(type) {
			case *schema.Resource:
				var path string
				if parent == "" {
					path = key
				} else {
					path = parent + "." + key
				}
				computed(et.Schema, path, omits)
			}
		}
	}
}

func base64Encode(v []string) string {
	var buf bytes.Buffer

	for i := range v {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(v[i])
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func base64Decode(v string) []string {
	joined, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil
	}

	return strings.Split(string(joined), "\n")
}
