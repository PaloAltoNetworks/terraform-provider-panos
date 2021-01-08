package panos

import (
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	agg "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"
	"github.com/PaloAltoNetworks/pango/netw/interface/vlan"
	"github.com/PaloAltoNetworks/pango/pnrm/template"
	"github.com/PaloAltoNetworks/pango/predefined/threat"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccProviders                                      map[string]terraform.ResourceProvider
	testAccProvider                                       *schema.Provider
	testAccIsFirewall, testAccIsPanorama                  bool
	testAccSupportsL2, testAccSupportsAggregateInterfaces bool
	testAccPanosVersion                                   version.Number
	testAccPlugins                                        map[string]string
	testAccPredefinedPhoneHomeThreats                     []threat.Entry
	testAccPredefinedVulnerabilityThreats                 []threat.Entry
)

func init() {
	var err error

	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"panos": testAccProvider,
	}

	/*
	   We need to know whether we are running acctests for the firewall or the
	   Panorama.  Since the provider meta is nil until Configure is called on
	   it, just check here in the init function.
	*/
	con, err := pango.Connect(pango.Client{
		Hostname: os.Getenv("PANOS_HOSTNAME"),
		Username: os.Getenv("PANOS_USERNAME"),
		Password: os.Getenv("PANOS_PASSWORD"),
		Logging:  pango.LogQuiet,
	})
	if err == nil {
		vt := vlan.Entry{
			Name:    "vlan.42",
			Comment: "acctest l2 check",
		}
		pt := template.Entry{
			Name:        "accL2Chk",
			Description: "acctest l2 check",
		}
		ai := agg.Entry{
			Name: "ae3",
			Mode: agg.ModeLayer3,
		}

		switch c := con.(type) {
		case *pango.Firewall:
			testAccIsFirewall = true
			testAccPanosVersion = c.Versioning()

			testAccPlugins = make(map[string]string)
			for _, v := range c.Plugin {
				if v.Installed == "yes" {
					testAccPlugins[v.Name] = v.Version
				}
			}

			if err = c.Network.VlanInterface.Set("", vt); err == nil {
				c.Network.VlanInterface.Delete(vt)
				testAccSupportsL2 = true
			}

			if err = c.Network.AggregateInterface.Edit("vsys1", ai); err == nil {
				c.Network.AggregateInterface.Delete(ai)
				testAccSupportsAggregateInterfaces = true
			}

			testAccPredefinedPhoneHomeThreats, _ = c.Predefined.Threat.GetThreats(threat.PhoneHome, "Phishing")
			testAccPredefinedVulnerabilityThreats, _ = c.Predefined.Threat.GetThreats(threat.Vulnerability, "Overflow")
		case *pango.Panorama:
			testAccIsPanorama = true
			testAccPanosVersion = c.Versioning()

			testAccPlugins = make(map[string]string)
			for _, v := range c.Plugin {
				if v.Installed == "yes" {
					testAccPlugins[v.Name] = v.Version
				}
			}

			if err = c.Panorama.Template.Set(pt); err == nil {
				if err = c.Network.VlanInterface.Set(pt.Name, "", "vsys1", vt); err == nil {
					c.Network.VlanInterface.Delete(pt.Name, "", vt)
					testAccSupportsL2 = true
				}

				if err = c.Network.AggregateInterface.Edit(pt.Name, "", "vsys1", ai); err == nil {
					c.Network.AggregateInterface.Delete(pt.Name, "", ai)
					testAccSupportsAggregateInterfaces = true
				}
				c.Panorama.Template.Delete(pt)
			}

			testAccPredefinedPhoneHomeThreats, _ = c.Predefined.Threat.GetThreats(threat.PhoneHome, "Phishing")
			testAccPredefinedVulnerabilityThreats, _ = c.Predefined.Threat.GetThreats(threat.Vulnerability, "Overflow")
		}
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PANOS_HOSTNAME") == "" {
		t.Fatal("PANOS_HOSTNAME must be set for acceptance tests")
	}
	if os.Getenv("PANOS_USERNAME") == "" {
		t.Fatal("PANOS_USERNAME must be set for acceptance tests")
	}
	if os.Getenv("PANOS_PASSWORD") == "" {
		t.Fatal("PANOS_PASSWORD must be set for acceptance tests")
	}
}
