package panos

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const IdSeparator string = ":"

func asStringList(d *schema.ResourceData, key string) []string {
	if d.Get(key) == nil || len(d.Get(key).([]interface{})) == 0 {
		return nil
	}

	list := d.Get(key).([]interface{})
	ans := make([]string, len(list))
	for i := range list {
		ans[i] = list[i].(string)
	}

	return ans
}

func setAsList(d *schema.ResourceData, key string) []string {
	list := d.Get(key).(*schema.Set).List()
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
