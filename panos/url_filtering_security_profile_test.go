package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing tests.
func TestAccPanosDsUrlFilteringSecurityProfileList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsUrlFilteringSecurityProfileConfig(name),
				Check:  checkDataSourceListing("panos_url_filtering_security_profiles"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsUrlFilteringSecurityProfile_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsUrlFilteringSecurityProfileConfig(name),
				Check: checkDataSource("panos_url_filtering_security_profile", []string{
					"name", "description", "ucd_mode", "log_container_page_only",
					"log_http_header_xff", "log_http_header_referer",
					"log_http_header_user_agent",
					"http_header_insertion.0.name", "http_header_insertion.0.type",
					"http_header_insertion.0.http_header.0.header", "http_header_insertion.0.http_header.0.value",
					"http_header_insertion.0.http_header.1.header", "http_header_insertion.0.http_header.1.value",
				}),
			},
		},
	})
}

func testAccDsUrlFilteringSecurityProfileConfig(name string) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

data "panos_url_filtering_security_profiles" "test" {}

data "panos_url_filtering_security_profile" "test" {
    name = panos_url_filtering_security_profile.x.name
}

resource "panos_url_filtering_security_profile" "x" {
    name = %q
    description = "url filtering sec prof data source acctest"
    ucd_mode = %q
    ucd_log_severity = "${
        data.panos_system_info.x.version_major > 8 ? "medium" : ""
    }"
    log_container_page_only = true
    log_http_header_xff = true
    log_http_header_referer = true
    log_http_header_user_agent = true
    http_header_insertion {
        name = "doublelift"
        type = "Custom"
        domains = [
            "b.example.com",
            "a.example.com",
            "c.example.com",
        ]
        http_header {
            header = "X-First-Header"
            value = "alpha"
        }
        http_header {
            header = "X-Second-Header"
            value = "beta"
        }
    }
}
`, name, url.UcdModeDisabled)
}

// Resource tests.
func TestAccPanosUrlFilteringSecurityProfile_basic(t *testing.T) {
	var o url.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosUrlFilteringSecurityProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUrlFilteringSecurityProfileConfig(name, "desc one", "business-and-economy", "auctions", "command-and-control", "high-risk", "legal", url.UcdModeDisabled, "a.example.com", "X-Palo-Alto-Networks", "palto alto networks", "X-Foo", "bar", true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosUrlFilteringSecurityProfileExists("panos_url_filtering_security_profile.test", &o),
					testAccCheckPanosUrlFilteringSecurityProfileAttributes(&o, name, "desc one", "business-and-economy", "auctions", "command-and-control", "high-risk", "legal", url.UcdModeDisabled, "a.example.com", "X-Palo-Alto-Networks", "palto alto networks", "X-Foo", "bar", true, false, true, false),
				),
			},
			{
				Config: testAccUrlFilteringSecurityProfileConfig(name, "desc two", "military", "job-search", "music", "news", "shopping", url.UcdModeIpUser, "b.example.com", "X-Palo-Alto-Networks", "panw", "X-New-Header", "new value", false, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosUrlFilteringSecurityProfileExists("panos_url_filtering_security_profile.test", &o),
					testAccCheckPanosUrlFilteringSecurityProfileAttributes(&o, name, "desc two", "military", "job-search", "music", "news", "shopping", url.UcdModeIpUser, "b.example.com", "X-Palo-Alto-Networks", "panw", "X-New-Header", "new value", false, true, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosUrlFilteringSecurityProfileExists(n string, o *url.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v url.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseUrlFilteringSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.UrlFilteringProfile.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseUrlFilteringSecurityProfileId(rs.Primary.ID)
			v, err = con.Objects.UrlFilteringProfile.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosUrlFilteringSecurityProfileAttributes(o *url.Entry, name, desc, allowCat, alertCat, blockCat, continueCat, overrideCat, mode, domain, hdr1, val1, hdr2, val2 string, tcp, xff, ua, ref bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		if len(o.AllowCategories) != 1 || o.AllowCategories[0] != allowCat {
			return fmt.Errorf("Allow categories is not [%q]: %#v", allowCat, o.AllowCategories)
		}

		if len(o.AlertCategories) != 1 || o.AlertCategories[0] != alertCat {
			return fmt.Errorf("Alert categories is not [%q]: %#v", alertCat, o.AlertCategories)
		}

		if len(o.BlockCategories) != 1 || o.BlockCategories[0] != blockCat {
			return fmt.Errorf("block categories is not [%q]: %#v", blockCat, o.BlockCategories)
		}

		if len(o.ContinueCategories) != 1 || o.ContinueCategories[0] != continueCat {
			return fmt.Errorf("Continue categories is not [%q]: %#v", continueCat, o.ContinueCategories)
		}

		if len(o.OverrideCategories) != 1 || o.OverrideCategories[0] != overrideCat {
			return fmt.Errorf("Override categories is not [%q]: %#v", overrideCat, o.OverrideCategories)
		}

		if o.UcdMode != mode {
			return fmt.Errorf("UcdMode is %q, not %q", o.UcdMode, mode)
		}

		if len(o.HttpHeaderInsertions) != 1 {
			return fmt.Errorf("Http header insertions is not len 1: %#v", o.HttpHeaderInsertions)
		}

		hi := o.HttpHeaderInsertions[0]

		if len(hi.Domains) != 1 || hi.Domains[0] != domain {
			return fmt.Errorf("HHI domains is not [%q]: %#v", domain, hi.Domains)
		}

		if len(hi.HttpHeaders) != 2 {
			return fmt.Errorf("HHI headers is not len2: %#v", hi.HttpHeaders)
		}

		if hi.HttpHeaders[0].Header != hdr1 {
			return fmt.Errorf("Header 0 is %q, not %q", hi.HttpHeaders[0].Header, hdr1)
		}

		if hi.HttpHeaders[0].Value != val1 {
			return fmt.Errorf("Value 0 is %q, not %q", hi.HttpHeaders[0].Value, val1)
		}

		if hi.HttpHeaders[1].Header != hdr2 {
			return fmt.Errorf("Header 1 is %q, not %q", hi.HttpHeaders[1].Header, hdr2)
		}

		if hi.HttpHeaders[1].Value != val2 {
			return fmt.Errorf("Value 1 is %q, not %q", hi.HttpHeaders[1].Value, val2)
		}

		if o.TrackContainerPage != tcp {
			return fmt.Errorf("Track container page is %t, not %t", o.TrackContainerPage, tcp)
		}

		if o.LogHttpHeaderXff != xff {
			return fmt.Errorf("log http header xff is %t, not %t", o.LogHttpHeaderXff, xff)
		}

		if o.LogHttpHeaderUserAgent != ua {
			return fmt.Errorf("log http header user agent is %t, not %t", o.LogHttpHeaderUserAgent, ua)
		}

		if o.LogHttpHeaderReferer != ref {
			return fmt.Errorf("log http header referer is %t, not %t", o.LogHttpHeaderReferer, ref)
		}

		return nil
	}
}

func testAccPanosUrlFilteringSecurityProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_url_filtering_security_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseUrlFilteringSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.UrlFilteringProfile.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseUrlFilteringSecurityProfileId(rs.Primary.ID)
				_, err = con.Objects.UrlFilteringProfile.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccUrlFilteringSecurityProfileConfig(name, desc, allowCat, alertCat, blockCat, continueCat, overrideCat, mode, domain, hdr1, val1, hdr2, val2 string, tcp, xff, ua, ref bool) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

resource "panos_url_filtering_security_profile" "test" {
    name = %q
    description = %q
    allow_categories = [%q]
    alert_categories = [%q]
    block_categories = [%q]
    continue_categories = [%q]
    override_categories = [%q]
    ucd_mode = %q
    ucd_log_severity = "${
        data.panos_system_info.x.version_major >=8 ? "medium" : ""
    }"
    http_header_insertion {
        name = "foo"
        type = "Custom"
        domains = [%q]
        http_header {
            header = %q
            value = %q
        }
        http_header {
            header = %q
            value = %q
        }
    }
    track_container_page = %t
    log_http_header_xff = %t
    log_http_header_user_agent = %t
    log_http_header_referer = %t
}
`, name, desc, allowCat, alertCat, blockCat, continueCat, overrideCat, mode, domain, hdr1, val1, hdr2, val2, tcp, xff, ua, ref)
}
