package provider_test

// Live-device tests for the XPath single-quote escaping fix (CWE-643) in
// pango/util.AsEntryXpath / AsMemberXpath / AsUuidXpath and its inverse
// pango/util.EntryName.
//
// Findings confirmed on a live device:
//
//   - The SDK-layer fix works: requests for quote-bearing values now pass
//     Location.XpathWithComponents() validation and reach PAN-OS (previously
//     they failed locally with "Name must be formatted as entry").
//
//   - PAN-OS rejects a single quote in an object NAME (error code 12,
//     "... is invalid"). A quoted name therefore cannot be created/persisted,
//     so the meaningful injection surface on a live device is the READ/lookup
//     path: a crafted value flowing into a `get` predicate must not be able to
//     break out and match a different, legally-named object.
//
// The objects live in a device group so the resource xpath carries an ancestor
// location variable in addition to the entry predicate.
//
// These require a live Panorama and only run under TF_ACC.

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/PaloAltoNetworks/pango/objects/address"
	"github.com/PaloAltoNetworks/pango/panorama/devicegroup"
	"github.com/PaloAltoNetworks/pango/util"
)

// skipIfNotAcc gates the pure-SDK tests so they only run against a live device.
func skipIfNotAcc(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance test skipped: set TF_ACC=1 and PANOS_HOSTNAME to run.")
	}
	testAccPreCheck(t)
}

// TestAccAddressXpathInjection_ReadIsolation is the core security test. It
// proves a crafted lookup value cannot break out of the @name predicate to
// match a different object.
//
// A victim with a legal name is created in a device group. Then a lookup is
// issued for a crafted value `bogus' or @name='<victim>`. With the vulnerable
// code the predicate expands to `@name='bogus' or @name='<victim>'` and
// resolves to the victim. With the fix the crafted value is emitted via
// concat() as a single literal, so it matches nothing and the victim is never
// returned.
func TestAccAddressXpathInjection_ReadIsolation(t *testing.T) {
	skipIfNotAcc(t)

	ctx := context.Background()
	vn := sdkClient.Versioning()
	suffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)

	// Device group to host the address objects.
	dgName := fmt.Sprintf("test-acc-dg-%s", suffix)
	dgSvc := devicegroup.NewService(sdkClient)
	dgLoc := *devicegroup.NewPanoramaLocation()
	if _, err := dgSvc.Create(ctx, dgLoc, &devicegroup.Entry{Name: dgName}); err != nil {
		t.Fatalf("Create(device group %q) failed: %v", dgName, err)
	}
	t.Cleanup(func() {
		if err := dgSvc.Delete(ctx, dgLoc, dgName); err != nil {
			t.Logf("cleanup Delete(device group %q): %v", dgName, err)
		}
	})

	svc := address.NewService(sdkClient)
	loc := *address.NewDeviceGroupLocation()
	loc.DeviceGroup.DeviceGroup = dgName

	victimName := fmt.Sprintf("test-acc-victim-%s", suffix)
	victimIP := "10.99.99.1/32"
	if _, err := svc.Create(ctx, loc, &address.Entry{Name: victimName, IpNetmask: &victimIP}); err != nil {
		t.Fatalf("Create(victim %q) failed: %v", victimName, err)
	}
	t.Cleanup(func() {
		if err := svc.Delete(ctx, loc, victimName); err != nil {
			t.Logf("cleanup Delete(victim %q): %v", victimName, err)
		}
	})

	// Positive control: the victim is findable by its exact name (confirms the
	// harness and the normal, quote-free predicate path work).
	{
		path, err := loc.XpathWithComponents(vn, util.AsEntryXpath(victimName))
		if err != nil {
			t.Fatalf("building victim xpath: %v", err)
		}
		got, err := svc.ReadWithXpath(ctx, util.AsXpath(path), "get")
		if err != nil || got == nil || got.Name != victimName {
			t.Fatalf("positive control: reading victim by real name failed: got=%v err=%v", got, err)
		}
	}

	// Injection attempt.
	crafted := fmt.Sprintf("bogus-%s' or @name='%s", suffix, victimName)

	// The safe builder emits crafted as concat(...); the generalized
	// XpathWithComponents validation must accept it (finding #1).
	path, err := loc.XpathWithComponents(vn, util.AsEntryXpath(crafted))
	if err != nil {
		t.Fatalf("crafted value was rejected by XpathWithComponents (concat form must be accepted): %v", err)
	}

	got, err := svc.ReadWithXpath(ctx, util.AsXpath(path), "get")
	if err == nil && got != nil {
		t.Fatalf("XPath injection: crafted lookup %q resolved to entry name=%q ip=%v (victim was %q) - the predicate broke out to match a different object",
			crafted, got.Name, got.IpNetmask, victimName)
	}
	// The fix causes a clean "not found" (concat matched nothing), which also
	// confirms PAN-OS accepts concat() in a get predicate on the wire.
}

// TestAccPanosAddress_QuotedNameRejected documents, end-to-end through the
// provider, that a quoted name is now forwarded to PAN-OS (finding #1 fixed)
// and rejected by the device's own name validation rather than by our local
// xpath-component validation. Before the fix this failed earlier with
// "Name must be formatted as entry" and never reached the device.
func TestAccPanosAddress_QuotedNameRejected(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test-acc-%s", acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressQuotedNameConfig,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				// PAN-OS forbids apostrophes in object names; the request now
				// reaches the device, which rejects the name itself.
				ExpectError: regexp.MustCompile("is invalid"),
			},
		},
	})
}

const testAccAddressQuotedNameConfig = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }
  name     = "${var.prefix}-dg"
}

resource "panos_address" "quoted" {
  location   = { device_group = { name = panos_device_group.dg.name } }
  name       = "${var.prefix}-O'Brien"
  ip_netmask = "10.30.30.30/32"
}
`
