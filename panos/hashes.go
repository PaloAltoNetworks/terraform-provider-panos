package panos

import (
	"bytes"
	"fmt"
)

// From https://developer.hashicorp.com/terraform/plugin/sdkv2/guides/v2-upgrade-guide#removal-of-helper-hashcode-package

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// Strings hashes a list of strings to a unique hashcode.
func Strings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", String(buf.String()))
}

func resourceMatchAddressPrefixHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s%t", m["prefix"].(string), m["exact"].(bool)))
	return String(buf.String())
}

func resourceTargetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(m["serial"].(string))
	vl := m["vsys_list"].([]interface{})
	for i := range vl {
		buf.WriteString(vl[i].(string))
	}
	return String(buf.String())
}
