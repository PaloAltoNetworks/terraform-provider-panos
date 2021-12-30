package nat

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

func versioning(v version.Number) (normalizer, func(Entry) interface{}) {
	if v.Gte(version.Number{9, 0, 0, ""}) {
		return &container_v3{}, specify_v3
	} else if v.Gte(version.Number{8, 1, 0, ""}) {
		return &container_v2{}, specify_v2
	} else {
		return &container_v1{}, specify_v1
	}
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

// FirewallNamespace returns an initialized namespace.
func FirewallNamespace(client util.XapiClient) *Firewall {
	return &Firewall{
		ns: &namespace.Policy{
			Standard: namespace.Standard{
				Common: namespace.Common{
					Singular: singular,
					Plural:   plural,
					Client:   client,
				},
			},
		},
	}
}

// PanoramaNamespace returns an initialized namespace.
func PanoramaNamespace(client util.XapiClient) *Panorama {
	return &Panorama{
		ns: &namespace.Policy{
			Standard: namespace.Standard{
				Common: namespace.Common{
					Singular: singular,
					Plural:   plural,
					Client:   client,
				},
			},
		},
	}
}

func RulesMatch(a, b Entry) bool {
	return a.Name == b.Name &&
		a.Description == b.Description &&
		a.Type == b.Type &&
		util.UnorderedListsMatch(a.SourceZones, b.SourceZones) &&
		a.DestinationZone == b.DestinationZone &&
		a.ToInterface == b.ToInterface &&
		a.Service == b.Service &&
		util.UnorderedListsMatch(a.SourceAddresses, b.SourceAddresses) &&
		util.UnorderedListsMatch(a.DestinationAddresses, b.DestinationAddresses) &&
		a.SatType == b.SatType &&
		a.SatAddressType == b.SatAddressType &&
		util.UnorderedListsMatch(a.SatTranslatedAddresses, b.SatTranslatedAddresses) &&
		a.SatInterface == b.SatInterface &&
		a.SatIpAddress == b.SatIpAddress &&
		a.SatFallbackType == b.SatFallbackType &&
		util.UnorderedListsMatch(a.SatFallbackTranslatedAddresses, b.SatFallbackTranslatedAddresses) &&
		a.SatFallbackInterface == b.SatFallbackInterface &&
		a.SatFallbackIpType == b.SatFallbackIpType &&
		a.SatFallbackIpAddress == b.SatFallbackIpAddress &&
		a.SatStaticTranslatedAddress == b.SatStaticTranslatedAddress &&
		a.SatStaticBiDirectional == b.SatStaticBiDirectional &&
		a.DatType == b.DatType &&
		a.DatAddress == b.DatAddress &&
		a.DatPort == b.DatPort &&
		a.DatDynamicDistribution == b.DatDynamicDistribution &&
		a.Disabled == b.Disabled &&
		util.TargetsMatch(a.Targets, b.Targets) &&
		a.NegateTarget == b.NegateTarget &&
		util.OrderedListsMatch(a.Tags, b.Tags)
}
