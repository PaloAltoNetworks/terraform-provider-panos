package panos

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	SkipPanoramaAccTest = "Skipping panorama test"
	SkipFirewallAccTest = "Skipping firewall test"
	SkipL2AccTest       = "Skipping L2 test for PAN-OS model that does not have L2 support"
	SkipAggregateTest   = "Skipping test as aggregate ethernet interfaces are not supported by PAN-OS"
)

func addStringInSliceValidation(desc string, values []string) string {
	var b strings.Builder
	b.Grow(len(desc) + 20*len(values))
	b.WriteString(desc)

	if len(values) > 0 {
		b.WriteString(" Valid values are")

		for i := range values {
			if i != 0 && len(values) > 2 {
				b.WriteString(",")
			}
			b.WriteString(" ")
			if i == len(values)-1 {
				b.WriteString("or ")
			}
			b.WriteString("`")
			b.WriteString(values[i])
			b.WriteString("`")
		}
		b.WriteString(".")
	}

	return b.String()
}

func validateStringIn(vals ...string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		ok := false
		for i := range vals {
			if vals[i] == value {
				ok = true
				break
			}
		}

		if !ok {
			errors = append(errors, fmt.Errorf("%q (%q) not in %#v", k, value, vals))
		}

		return
	}
}

func validateIntInRange(low, high int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < low || value > high {
			errors = append(errors, fmt.Errorf("%q (%d) not in range [%d, %d]", k, value, low, high))
		}

		return
	}
}

func validateSetKeyIsUnique(key string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		counts := make(map[string]int)
		list := v.(*schema.Set).List()
		for i := range list {
			group := list[i].(map[string]interface{})
			val := group[key].(string)
			counts[val] = counts[val] + 1
		}

		for ck, cv := range counts {
			if cv > 1 {
				errors = append(errors, fmt.Errorf("%q (%s) is not unique - repeated %d times", k, ck, cv))
			}
		}

		return
	}
}

func validateStringHasPrefix(p string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		val := v.(string)
		if !strings.HasPrefix(val, p) {
			errors = append(errors, fmt.Errorf("Param value must start with %q", p))
		}

		return
	}
}

func validateIsRegex() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		val := v.(string)
		if _, err := regexp.Compile(val); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
