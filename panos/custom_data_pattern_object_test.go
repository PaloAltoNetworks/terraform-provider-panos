package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/custom/data"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsCustomDataPatternObjectList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsCustomDataPatternObjectConfig(name, data.TypePredefined),
				Check:  checkDataSourceListing("panos_custom_data_pattern_objects"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsCustomDataPatternObjectPredefinedPattern(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsCustomDataPatternObjectConfig(name, data.TypePredefined),
				Check: checkDataSource("panos_custom_data_pattern_object", []string{
					"name", "description", "type",
					"predefined_pattern.0.name",
					"predefined_pattern.0.file_types.0",
					"predefined_pattern.0.file_types.1",
				}),
			},
		},
	})
}

func TestAccPanosDsCustomDataPatternObjectRegex(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsCustomDataPatternObjectConfig(name, data.TypeRegex),
				Check: checkDataSource("panos_custom_data_pattern_object", []string{
					"name", "description", "type",
					"regex.0.name",
					"regex.0.regex",
					"regex.0.file_types.0",
					"regex.0.file_types.1",
					"regex.0.file_types.2",
				}),
			},
		},
	})
}

func TestAccPanosDsCustomDataPatternObjectFileProperties(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsCustomDataPatternObjectConfig(name, data.TypeFileProperties),
				Check: checkDataSource("panos_custom_data_pattern_object", []string{
					"name", "description", "type",
					"file_property.0.name",
					"file_property.0.file_type",
					"file_property.0.file_property",
					"file_property.0.property_value",
				}),
			},
		},
	})
}

func testAccDsCustomDataPatternObjectConfig(name, tp string) string {
	ans := fmt.Sprintf(`
data "panos_custom_data_pattern_objects" "test" {}

data "panos_custom_data_pattern_object" "test" {
    name = panos_custom_data_pattern_object.x.name
}

resource "panos_custom_data_pattern_object" "x" {
    name = %q
    description = "custom data pattern object ds acctest"
    type = %q`, name, tp)

	switch tp {
	case data.TypePredefined:
		ans = fmt.Sprintf(`%s
    predefined_pattern {
        name = "social-security-numbers"
        file_types = ["docx", "xlsx"]`, ans)
	case data.TypeRegex:
		ans = fmt.Sprintf(`%s
    regex {
        name = "blah"
        file_types = ["docx", "doc", "text/html"]
        regex = "shin megami tensei"`, ans)
	case data.TypeFileProperties:
		ans = fmt.Sprintf(`%s
    file_property {
        name = "blah"
        file_type = "pdf"
        file_property = "panav-rsp-pdf-dlp-keywords"
        property_value = "foo"`, ans)
	}

	return fmt.Sprintf(`%s
    }
}`, ans)
}

// Resource tests.
func TestAccPanosCustomDataPatternObjectPredefinedPattern_basic(t *testing.T) {
	var o data.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	pp := &data.PredefinedPattern{
		Name:      "social-security-numbers",
		FileTypes: []string{"xlsx", "doc"},
	}
	r := &data.Regex{
		Name:      "second",
		FileTypes: []string{"text/html", "docx"},
		Regex:     "shin megami tensei",
	}
	fp := &data.FileProperty{
		Name:          "third",
		FileType:      "pdf",
		FileProperty:  "panav-rsp-pdf-dlp-keywords",
		PropertyValue: "foo",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosCustomDataPatternObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomDataPatternObjectConfig(name, "descOne", data.TypePredefined, pp, nil, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosCustomDataPatternObjectExists("panos_custom_data_pattern_object.test", &o),
					testAccCheckPanosCustomDataPatternObjectAttributes(&o, name, "descOne", data.TypePredefined, pp, nil, nil),
				),
			},
			{
				Config: testAccCustomDataPatternObjectConfig(name, "descTwo", data.TypeRegex, nil, r, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosCustomDataPatternObjectExists("panos_custom_data_pattern_object.test", &o),
					testAccCheckPanosCustomDataPatternObjectAttributes(&o, name, "descTwo", data.TypeRegex, nil, r, nil),
				),
			},
			{
				Config: testAccCustomDataPatternObjectConfig(name, "descThree", data.TypeFileProperties, nil, nil, fp),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosCustomDataPatternObjectExists("panos_custom_data_pattern_object.test", &o),
					testAccCheckPanosCustomDataPatternObjectAttributes(&o, name, "descThree", data.TypeFileProperties, nil, nil, fp),
				),
			},
		},
	})
}

func testAccCheckPanosCustomDataPatternObjectExists(n string, o *data.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v data.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseCustomDataPatternObjectId(rs.Primary.ID)
			v, err = con.Objects.DataPattern.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseCustomDataPatternObjectId(rs.Primary.ID)
			v, err = con.Objects.DataPattern.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosCustomDataPatternObjectAttributes(o *data.Entry, name, desc, tp string, pp *data.PredefinedPattern, r *data.Regex, fp *data.FileProperty) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, not %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.Type != tp {
			return fmt.Errorf("Type is %s, not %s", o.Type, tp)
		}

		if pp == nil {
			if len(o.PredefinedPatterns) != 0 {
				return fmt.Errorf("predefined patterns is not nil: %#v", o.PredefinedPatterns)
			}
		} else if len(o.PredefinedPatterns) != 1 {
			return fmt.Errorf("predefined patterns should be 1: %#v", o.PredefinedPatterns)
		} else {
			x := o.PredefinedPatterns[0]

			if x.Name != pp.Name {
				return fmt.Errorf("pp.Name is %q, not %q", x.Name, pp.Name)
			}

			if len(x.FileTypes) != len(pp.FileTypes) {
				return fmt.Errorf("pp.FileTypes is not %d: %#v", len(pp.FileTypes), x.FileTypes)
			}

			for i := 0; i < len(pp.FileTypes); i++ {
				if x.FileTypes[i] != pp.FileTypes[i] {
					return fmt.Errorf("pp.FileTypes[%d] is %q, not %q", i, x.FileTypes[i], pp.FileTypes[i])
				}
			}
		}

		if r == nil {
			if len(o.Regexes) != 0 {
				return fmt.Errorf("regexes is not nil: %#v", o.Regexes)
			}
		} else if len(o.Regexes) != 1 {
			return fmt.Errorf("regexes should be 1: %#v", o.Regexes)
		} else {
			x := o.Regexes[0]

			if x.Name != r.Name {
				return fmt.Errorf("r.Name is %q, not %q", x.Name, pp.Name)
			}

			if len(x.FileTypes) != len(r.FileTypes) {
				return fmt.Errorf("r.FileTypes is not %d: %#v", len(r.FileTypes), x.FileTypes)
			}

			for i := 0; i < len(r.FileTypes); i++ {
				if x.FileTypes[i] != r.FileTypes[i] {
					return fmt.Errorf("r.FileTypes[%d] is %q, not %q", i, x.FileTypes[i], pp.FileTypes[i])
				}
			}

			if x.Regex != r.Regex {
				return fmt.Errorf("r.Regex is %q, not %q", x.Regex, r.Regex)
			}
		}

		if fp == nil {
			if len(o.FileProperties) != 0 {
				return fmt.Errorf("regexes is not nil: %#v", o.FileProperties)
			}
		} else if len(o.FileProperties) != 1 {
			return fmt.Errorf("regexes should be 1: %#v", o.FileProperties)
		} else {
			x := o.FileProperties[0]

			if x.Name != fp.Name {
				return fmt.Errorf("fp.Name is %q, not %q", x.Name, pp.Name)
			}

			if x.FileType != fp.FileType {
				return fmt.Errorf("fp.FileType is %q, not %q", x.FileType, fp.FileType)
			}

			if x.FileProperty != fp.FileProperty {
				return fmt.Errorf("fp.FileProperty is %q, not %q", x.FileProperty, fp.FileProperty)
			}

			if x.PropertyValue != fp.PropertyValue {
				return fmt.Errorf("fp.PropertyValue is %q, not %q", x.PropertyValue, fp.PropertyValue)
			}
		}

		return nil
	}
}

func testAccPanosCustomDataPatternObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_custom_data_pattern_object" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseCustomDataPatternObjectId(rs.Primary.ID)
				_, err = con.Objects.DataPattern.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseCustomDataPatternObjectId(rs.Primary.ID)
				_, err = con.Objects.DataPattern.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccCustomDataPatternObjectConfig(name, desc, tp string, pp *data.PredefinedPattern, r *data.Regex, fp *data.FileProperty) string {
	prefix := fmt.Sprintf(`
resource "panos_custom_data_pattern_object" "test" {
    name = %q
    description = %q
    type = %q
`, name, desc, tp)

	switch {
	case pp != nil:
		return fmt.Sprintf(`%s
    predefined_pattern {
        name = %q
        file_types = [%q, %q]
    }
}`, prefix, pp.Name, pp.FileTypes[0], pp.FileTypes[1])
	case r != nil:
		return fmt.Sprintf(`%s
    regex {
        name = %q
        file_types = [%q, %q]
        regex = %q
    }
}`, prefix, r.Name, r.FileTypes[0], r.FileTypes[1], r.Regex)
	case fp != nil:
		return fmt.Sprintf(`%s
    file_property {
        name = %q
        file_type = %q
        file_property = %q
        property_value = %q
    }
}`, prefix, fp.Name, fp.FileType, fp.FileProperty, fp.PropertyValue)
	}

	return prefix
}
