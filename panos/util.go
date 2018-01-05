package panos

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const IdSeparator string = ":"

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
