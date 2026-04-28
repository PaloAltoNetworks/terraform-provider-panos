package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccPredefinedDlpFileType_Basic reads a well-known predefined DLP file type
// ("pdf") and verifies that its name and file_property list are correctly populated.
// The "pdf" entry is guaranteed to exist on PAN-OS devices and its file-property
// entries (e.g. panav-rsp-pdf-dlp-author) have stable, device-supplied labels.
// ListPartial with index 0 asserts the list is non-empty and the first element is
// a well-formed object; exact entry names are not pinned because ordering may vary
// across device versions.
func TestAccPredefinedDlpFileType_Basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: predefinedDlpFileType_Basic_Tmpl,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pdf"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("file_property"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"label": knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}

// TestAccPredefinedDlpFileType_FileProperties verifies that the file_property
// sub-entries for the "pdf" file type are read with both name and label fields
// populated. The full, ordered list of file-property sub-entries is asserted
// against the values returned by Panorama 11.2. The device returns 8 properties
// for "pdf" in the order: title, author, subject, comments, keywords,
// titus-corp-sensitivity, titus-corp-classification, titus-GUID.
func TestAccPredefinedDlpFileType_FileProperties(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: predefinedDlpFileType_FileProperties_Tmpl,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pdf"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("file_property"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-title"),
								"label": knownvalue.StringExact("Title"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-author"),
								"label": knownvalue.StringExact("Author"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-subject"),
								"label": knownvalue.StringExact("Subject"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-comments"),
								"label": knownvalue.StringExact("Comments"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-keywords"),
								"label": knownvalue.StringExact("Keywords"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-titus-corp-sensitivity"),
								"label": knownvalue.StringExact("Sensitivity"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-titus-corp-classification"),
								"label": knownvalue.StringExact("Classification"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("panav-rsp-pdf-dlp-titus-GUID"),
								"label": knownvalue.StringExact("TITUS GUID"),
							}),
						}),
					),
				},
			},
		},
	})
}

// TestAccPredefinedDlpFileType_MultipleEntries reads several well-known predefined
// DLP file types in a single Terraform configuration and verifies that each one
// resolves independently with the correct name. This exercises the data source
// for the full set of known entries (docx, pptx, xlsx, pdf) and confirms that
// multiple data source instances can coexist in one plan.
func TestAccPredefinedDlpFileType_MultipleEntries(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: predefinedDlpFileType_MultipleEntries_Tmpl,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.docx",
						tfjsonpath.New("name"),
						knownvalue.StringExact("docx"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.docx",
						tfjsonpath.New("file_property"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"label": knownvalue.NotNull(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pptx",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pptx"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pptx",
						tfjsonpath.New("file_property"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"label": knownvalue.NotNull(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.xlsx",
						tfjsonpath.New("name"),
						knownvalue.StringExact("xlsx"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.xlsx",
						tfjsonpath.New("file_property"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"label": knownvalue.NotNull(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pdf"),
					),
					statecheck.ExpectKnownValue(
						"data.panos_predefined_dlp_file_type.pdf",
						tfjsonpath.New("file_property"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":  knownvalue.NotNull(),
								"label": knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}

// predefinedDlpFileType_Basic_Tmpl reads the "pdf" predefined DLP file type
// from the device using the predefined location. No supporting resources are
// required because the entry is device-supplied and read-only.
const predefinedDlpFileType_Basic_Tmpl = `
data "panos_predefined_dlp_file_type" "pdf" {
  location = { predefined = {} }
  name     = "pdf"
}
`

// predefinedDlpFileType_FileProperties_Tmpl reads "pdf" without specifying an
// explicit file_property list so the provider returns all entries from the device
// in the server-defined order. The ListExact assertion in the test captures the
// full set of 8 properties returned by Panorama 11.2.
const predefinedDlpFileType_FileProperties_Tmpl = `
data "panos_predefined_dlp_file_type" "pdf" {
  location = { predefined = {} }
  name     = "pdf"
}
`

// TestAccPredefinedDlpFileType_CustomDataPatternLookup mirrors the v1 provider
// pattern where a user defines custom data patterns with file_property entries
// that reference file types by name and properties by human-readable label:
//
//	custom_data_pattern:
//	  - name: "test3"
//	    type: "file-properties"
//	    file_property:
//	      - name: "blah2"
//	        file_type: "rtf"
//	        file_property: "Keywords/Tags"
//	        property_value: "foo"
//
// The test uses a local variable to model this input, looks up each unique file
// type via the data source, resolves label → internal property name in locals,
// and then constructs the final resolved list that would feed into a
// panos_custom_data_object resource. Outputs verify the resolved values.
func TestAccPredefinedDlpFileType_CustomDataPatternLookup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: predefinedDlpFileType_CustomDataPatternLookup_Tmpl,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"resolved_file_properties",
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":           knownvalue.StringExact("blah2"),
								"file_type":      knownvalue.StringExact("rtf"),
								"file_property":  knownvalue.StringExact("panav-rsp-rtf-dlp-keywords"),
								"property_value": knownvalue.StringExact("foo"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":           knownvalue.StringExact("blah3"),
								"file_type":      knownvalue.StringExact("pdf"),
								"file_property":  knownvalue.StringExact("panav-rsp-pdf-dlp-author"),
								"property_value": knownvalue.StringExact("bar"),
							}),
						}),
					),
				},
			},
		},
	})
}

// predefinedDlpFileType_CustomDataPatternLookup_Tmpl models the full v1 workflow:
//
//  1. A local variable defines custom data pattern entries (mimicking user input)
//  2. Data sources look up each unique file type
//  3. A local resolves each entry's label to the internal property name
//  4. The resolved list is output — in production it would feed into
//     panos_custom_data_object file_properties pattern entries
const predefinedDlpFileType_CustomDataPatternLookup_Tmpl = `
# Step 1: User-defined custom data pattern input (mirrors v1 variable structure)
locals {
  custom_data_pattern = {
    name = "test3"
    type = "file-properties"
    file_property = [
      {
        name           = "blah2"
        file_type      = "rtf"
        file_property  = "Keywords/Tags"
        property_value = "foo"
      },
      {
        name           = "blah3"
        file_type      = "pdf"
        file_property  = "Author"
        property_value = "bar"
      },
    ]
  }

  # Collect unique file types from the input
  unique_file_types = distinct([
    for fp in local.custom_data_pattern.file_property : fp.file_type
  ])
}

# Step 2: Look up each unique file type (one data source per file type)
data "panos_predefined_dlp_file_type" "lookup_rtf" {
  location = { predefined = {} }
  name     = "rtf"
}

data "panos_predefined_dlp_file_type" "lookup_pdf" {
  location = { predefined = {} }
  name     = "pdf"
}

locals {
  # Build a map of file_type -> (label -> internal_name) for easy lookup
  file_type_data = {
    "rtf" = data.panos_predefined_dlp_file_type.lookup_rtf.file_property
    "pdf" = data.panos_predefined_dlp_file_type.lookup_pdf.file_property
  }

  # Step 3: Resolve each file_property entry's label to internal property name
  resolved_file_properties = [
    for fp in local.custom_data_pattern.file_property : {
      name           = fp.name
      file_type      = fp.file_type
      file_property  = one([
        for p in local.file_type_data[fp.file_type]
        : p.name if p.label == fp.file_property
      ])
      property_value = fp.property_value
    }
  ]
}

# Step 4: Output the resolved list (in production, this feeds into
# panos_custom_data_object pattern_type.file_properties.pattern entries)
output "resolved_file_properties" {
  value = local.resolved_file_properties
}
`

// predefinedDlpFileType_MultipleEntries_Tmpl reads the four known file types
// (docx, pptx, xlsx, pdf) in one configuration to exercise independent data
// source instances sharing the same predefined location.
const predefinedDlpFileType_MultipleEntries_Tmpl = `
data "panos_predefined_dlp_file_type" "docx" {
  location = { predefined = {} }
  name     = "docx"
}

data "panos_predefined_dlp_file_type" "pptx" {
  location = { predefined = {} }
  name     = "pptx"
}

data "panos_predefined_dlp_file_type" "xlsx" {
  location = { predefined = {} }
  name     = "xlsx"
}

data "panos_predefined_dlp_file_type" "pdf" {
  location = { predefined = {} }
  name     = "pdf"
}
`
