package cloudwatch

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/util"
)

func versioning(pkgs []plugin.Info) (normalizer, func(Config) interface{}, error) {
	pluginName := "vm_series"

	for _, pkg := range pkgs {
		// There seem to be a _lot_ of vm_series plugins around, both versioned and
		// unversioned.  So just do a prefix search for one that's installed for now.
		if strings.HasPrefix(pkg.Name, pluginName) {
			if pkg.Installed != "yes" {
				continue
			}

			// As mentioned above, not sure if there are versioning requirements on this
			// just yet, so just comment version checking out until we know for sure.

			/*
				if !strings.HasPrefix(pkg.Version, "1.") {
					return nil, nil, fmt.Errorf("need %q plugin version 1, but %q is installed", pluginName, pkg.Version)
				}
			*/
			return &container_v1{}, specify_v1, nil
		}
	}

	return nil, nil, fmt.Errorf("plugin %q not found", pluginName)
}

func specifier(e Config) []namespace.PluginSpecifier {
	return []namespace.PluginSpecifier{e}
}

func container(pkgs []plugin.Info) (normalizer, error) {
	r, _, err := versioning(pkgs)
	return r, err
}

func first(ans normalizer, err error) (Config, error) {
	if err != nil {
		return Config{}, err
	}

	return ans.Normalize()[0], nil
}

// FirewallNamespace returns an initialized namespace.
func FirewallNamespace(client util.XapiClient) *Firewall {
	return &Firewall{
		ns: &namespace.Plugin{
			Common: namespace.Common{
				Singular: singular,
				Client:   client,
			},
		},
	}
}
