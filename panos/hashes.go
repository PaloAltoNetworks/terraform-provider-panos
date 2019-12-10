package panos

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

func resourceMatchAddressPrefixHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s%t", m["prefix"].(string), m["exact"].(bool)))
	return hashcode.String(buf.String())
}

func resourceTargetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["serial"].(string))
	vl := m["vsys_list"].([]interface{})
	for i := range vl {
		buf.WriteString(vl[i].(string))
	}
	return hashcode.String(buf.String())
}
