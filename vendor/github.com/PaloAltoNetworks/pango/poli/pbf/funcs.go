package pbf

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

func versioning(v version.Number) (normalizer, func(Entry) interface{}) {
	if v.Gte(version.Number{9, 0, 0, ""}) {
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
		util.OrderedListsMatch(a.Tags, b.Tags) &&
		a.FromType == b.FromType &&
		util.UnorderedListsMatch(a.FromValues, b.FromValues) &&
		util.UnorderedListsMatch(a.SourceAddresses, b.SourceAddresses) &&
		util.UnorderedListsMatch(a.SourceUsers, b.SourceUsers) &&
		a.NegateSource == b.NegateSource &&
		util.UnorderedListsMatch(a.DestinationAddresses, b.DestinationAddresses) &&
		a.NegateDestination == b.NegateDestination &&
		util.UnorderedListsMatch(a.Applications, b.Applications) &&
		util.UnorderedListsMatch(a.Services, b.Services) &&
		a.Schedule == b.Schedule &&
		a.Disabled == b.Disabled &&
		a.Action == b.Action &&
		a.ForwardVsys == b.ForwardVsys &&
		a.ForwardEgressInterface == b.ForwardEgressInterface &&
		a.ForwardNextHopType == b.ForwardNextHopType &&
		a.ForwardNextHopValue == b.ForwardNextHopValue &&
		a.ForwardMonitorProfile == b.ForwardMonitorProfile &&
		a.ForwardMonitorIpAddress == b.ForwardMonitorIpAddress &&
		a.ForwardMonitorDisableIfUnreachable == b.ForwardMonitorDisableIfUnreachable &&
		a.EnableEnforceSymmetricReturn == b.EnableEnforceSymmetricReturn &&
		util.OrderedListsMatch(a.SymmetricReturnAddresses, b.SymmetricReturnAddresses) &&
		a.ActiveActiveDeviceBinding == b.ActiveActiveDeviceBinding &&
		util.TargetsMatch(a.Targets, b.Targets) &&
		// Don't compare UUID.
		a.NegateTarget == b.NegateTarget
}
