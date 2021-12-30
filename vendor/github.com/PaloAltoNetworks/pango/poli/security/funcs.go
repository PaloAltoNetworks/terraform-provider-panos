package security

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

func versioning(v version.Number) (normalizer, func(Entry) interface{}) {
	if v.Gte(version.Number{10, 0, 0, ""}) {
		return &container_v3{}, specify_v3
	} else if v.Gte(version.Number{9, 0, 0, ""}) {
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
		a.Type == b.Type &&
		a.Description == b.Description &&
		util.OrderedListsMatch(a.Tags, b.Tags) &&
		util.UnorderedListsMatch(a.SourceZones, b.SourceZones) &&
		util.UnorderedListsMatch(a.SourceAddresses, b.SourceAddresses) &&
		a.NegateSource == b.NegateSource &&
		util.UnorderedListsMatch(a.SourceUsers, b.SourceUsers) &&
		util.UnorderedListsMatch(a.HipProfiles, b.HipProfiles) &&
		util.UnorderedListsMatch(a.DestinationZones, b.DestinationZones) &&
		util.UnorderedListsMatch(a.DestinationAddresses, b.DestinationAddresses) &&
		a.NegateDestination == b.NegateDestination &&
		util.UnorderedListsMatch(a.Applications, b.Applications) &&
		util.UnorderedListsMatch(a.Services, b.Services) &&
		util.UnorderedListsMatch(a.Categories, b.Categories) &&
		a.Action == b.Action &&
		a.LogSetting == b.LogSetting &&
		a.LogStart == b.LogStart &&
		a.LogEnd == b.LogEnd &&
		a.Disabled == b.Disabled &&
		a.Schedule == b.Schedule &&
		a.IcmpUnreachable == b.IcmpUnreachable &&
		a.DisableServerResponseInspection == b.DisableServerResponseInspection &&
		a.Group == b.Group &&
		util.TargetsMatch(a.Targets, b.Targets) &&
		a.NegateTarget == b.NegateTarget &&
		a.Virus == b.Virus &&
		a.Spyware == b.Spyware &&
		a.Vulnerability == b.Vulnerability &&
		a.UrlFiltering == b.UrlFiltering &&
		a.FileBlocking == b.FileBlocking &&
		a.WildFireAnalysis == b.WildFireAnalysis &&
		a.DataFiltering == b.DataFiltering &&
		// Don't compare UUID.
		a.GroupTag == b.GroupTag &&
		util.UnorderedListsMatch(a.SourceDevices, b.SourceDevices) &&
		util.UnorderedListsMatch(a.DestinationDevices, b.DestinationDevices)
}
