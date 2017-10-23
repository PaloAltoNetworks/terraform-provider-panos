package main

import (
    "github.com/hashicorp/terraform/helper/schema"
)


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
