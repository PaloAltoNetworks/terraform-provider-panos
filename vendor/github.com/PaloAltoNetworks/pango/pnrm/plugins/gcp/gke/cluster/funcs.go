package cluster

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/util"
)

func versioning(pkgs []plugin.Info) (normalizer, func(Entry) interface{}, error) {
	name := "gcp"

	for _, pkg := range pkgs {
		if pkg.Name == name {
			if pkg.Installed != "yes" {
				return nil, nil, fmt.Errorf("Plugin not installed: %s", name)
			}
			if !strings.HasPrefix(pkg.Version, "1.") {
				return nil, nil, fmt.Errorf("Need %q plugin version 1, but %q is installed", name, pkg.Version)
			}
			return &container_v1{}, specify_v1, nil
		}
	}

	return nil, nil, fmt.Errorf("Plugin %q not found", name)
}

func specifier(e ...Entry) []namespace.PluginSpecifier {
	ans := make([]namespace.PluginSpecifier, 0, len(e))

	var val namespace.PluginSpecifier
	for _, x := range e {
		val = x
		ans = append(ans, val)
	}

	return ans
}

func container(pkgs []plugin.Info) (normalizer, error) {
	r, _, err := versioning(pkgs)
	return r, err
}

func first(ans normalizer, err error) (Entry, error) {
	if err != nil {
		return Entry{}, err
	}

	return ans.Normalize()[0], nil
}

func all(ans normalizer, err error) ([]Entry, error) {
	if err != nil {
		return nil, err
	}

	return ans.Normalize(), nil
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

// PanoramaNamespace returns an initialized namespace.
func PanoramaNamespace(client util.XapiClient) *Panorama {
	return &Panorama{
		ns: &namespace.Plugin{
			Common: namespace.Common{
				Singular: singular,
				Plural:   plural,
				Client:   client,
			},
		},
	}
}
