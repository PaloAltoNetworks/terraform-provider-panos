package manager

import (
	"fmt"
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/util"
	"strings"
)

// checks version of the install plugins and returns a normalizer "container" struct depending of the version installed.
func versioning(pkgs []plugin.Info) (normalizer, func(Entry) interface{}, error) {
	pluginName := "sw_fw_license"

	for _, pkg := range pkgs {
		if pkg.Name == pluginName {
			if pkg.Installed != "yes" {
				return nil, nil, fmt.Errorf("plugin not installed: %s", pluginName)
			}
			if !strings.HasPrefix(pkg.Version, "1.") {
				return nil, nil, fmt.Errorf("need %q plugin version 1, but %q is installed", pluginName, pkg.Version)
			}
			return &container_v1{}, specify_v1, nil
		}
	}

	return nil, nil, fmt.Errorf("plugin %q not found", pluginName)
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

// returns the container struct from the versioning function dependent of plugin version
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
