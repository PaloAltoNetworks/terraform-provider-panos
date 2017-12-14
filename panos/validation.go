package panos

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

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
