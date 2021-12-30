package decryption

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

func versioning(v version.Number) (normalizer, func(Entry) interface{}) {
	if v.Gte(version.Number{10, 0, 0, ""}) {
		return &container_v4{}, specify_v4
	} else if v.Gte(version.Number{9, 0, 0, ""}) {
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
		util.UnorderedListsMatch(a.SourceZones, b.SourceZones) &&
		util.UnorderedListsMatch(a.SourceAddresses, b.SourceAddresses) &&
		a.NegateSource == b.NegateSource &&
		util.UnorderedListsMatch(a.SourceUsers, b.SourceUsers) &&
		util.UnorderedListsMatch(a.DestinationZones, b.DestinationZones) &&
		util.UnorderedListsMatch(a.DestinationAddresses, b.DestinationAddresses) &&
		a.NegateDestination == b.NegateDestination &&
		util.OrderedListsMatch(a.Tags, b.Tags) &&
		a.Disabled == b.Disabled &&
		util.UnorderedListsMatch(a.Services, b.Services) &&
		util.UnorderedListsMatch(a.UrlCategories, b.UrlCategories) &&
		a.Action == b.Action &&
		a.DecryptionType == b.DecryptionType &&
		a.SslCertificate == b.SslCertificate &&
		a.DecryptionProfile == b.DecryptionProfile &&
		util.TargetsMatch(a.Targets, b.Targets) &&
		a.NegateTarget == b.NegateTarget &&
		a.ForwardingProfile == b.ForwardingProfile &&
		// Don't compare UUID
		a.GroupTag == b.GroupTag &&
		util.UnorderedListsMatch(a.SourceHips, b.SourceHips) &&
		util.UnorderedListsMatch(a.DestinationHips, b.DestinationHips) &&
		a.LogSuccessfulTlsHandshakes == b.LogSuccessfulTlsHandshakes &&
		a.LogFailedTlsHandshakes == b.LogFailedTlsHandshakes &&
		a.LogSetting == b.LogSetting
}
