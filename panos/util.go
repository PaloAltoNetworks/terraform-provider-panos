package panos

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

const IdSeparator string = ":"

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
	ans := make([]string, len(mm))

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
		ans[i] = v[i].(string)
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
