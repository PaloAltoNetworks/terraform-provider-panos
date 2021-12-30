package util

func CopyStringSlice(v []string) []string {
	if v == nil {
		return nil
	}

	ans := make([]string, len(v))
	copy(ans, v)

	return ans
}

func CopyTargets(v map[string][]string) map[string][]string {
	if v == nil {
		return nil
	}

	ans := make(map[string][]string)
	for key, oval := range v {
		if oval == nil {
			ans[key] = nil
		} else {
			val := make([]string, len(oval))
			copy(val, oval)
			ans[key] = val
		}
	}

	return ans
}
