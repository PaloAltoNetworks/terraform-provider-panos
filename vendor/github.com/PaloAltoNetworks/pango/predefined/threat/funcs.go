package threat

import (
	"fmt"
	"regexp"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

func versioning(v version.Number) (normalizer, func(Entry) interface{}) {
	return &container_v1{}, specify_v1
}

func specifier(e ...Entry) []namespace.Specifier {
	ans := make([]namespace.Specifier, 0, len(e))

	var val namespace.Specifier
	for _, x := range e {
		val = x
		ans = append(ans, val)
	}

	return ans
}

func container(v version.Number) normalizer {
	r, _ := versioning(v)
	return r
}

func first(ans normalizer, err error) (Entry, error) {
	if err != nil {
		return Entry{}, err
	}

	return ans.Normalize()[0], nil
}

func finder(ans normalizer, re *regexp.Regexp, err error) ([]Entry, error) {
	if err != nil {
		return nil, err
	}

	all := ans.Normalize()
	if re == nil {
		return all, nil
	}

	filtered := make([]Entry, 0, len(all))
	for _, x := range all {
		if re.MatchString(x.ThreatName) {
			filtered = append(filtered, x)
		}
	}

	return filtered, nil
}

func toNames(e []interface{}) ([]string, error) {
	ans := make([]string, len(e))
	for i := range e {
		switch v := e[i].(type) {
		case string:
			ans[i] = v
		case Entry:
			ans[i] = v.Name
		default:
			return nil, fmt.Errorf("invalid type: %s", v)
		}
	}

	return ans, nil
}

// FirewallNamespace returns an initialized namespace.
func FirewallNamespace(client util.XapiClient) *Firewall {
	return &Firewall{
		ns: &namespace.Standard{
			Common: namespace.Common{
				Singular: singular,
				Plural:   plural,
				Client:   client,
			},
		},
	}
}

// PanoramaNamespace returns an initialized namespace.
func PanoramaNamespace(client util.XapiClient) *Panorama {
	return &Panorama{
		ns: &namespace.Standard{
			Common: namespace.Common{
				Singular: singular,
				Plural:   plural,
				Client:   client,
			},
		},
	}
}
