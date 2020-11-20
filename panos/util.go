package panos

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	e2, ok := e.(pango.PanosError)
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

func parseTarget(v interface{}) map[string][]string {
	ans := make(map[string][]string)
	sl := v.(*schema.Set).List()

	for i := range sl {
		dev := sl[i].(map[string]interface{})
		key := dev["serial"].(string)
		value := asStringList(dev["vsys_list"].(*schema.Set).List())
		ans[key] = value
	}

	return ans
}

func buildTarget(m map[string][]string) *schema.Set {
	ans := &schema.Set{
		F: resourceTargetHash,
	}

	for k, v := range m {
		ans.Add(map[string]interface{}{
			"serial":    k,
			"vsys_list": listAsSet(v),
		})
	}

	return ans
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
