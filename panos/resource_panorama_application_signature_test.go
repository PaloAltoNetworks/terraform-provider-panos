package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/signature"
	"github.com/PaloAltoNetworks/pango/objs/app/signature/andcond"
	"github.com/PaloAltoNetworks/pango/objs/app/signature/orcond"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaApplicationSignature_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o signature.Entry
	var andList []andcond.Entry
	var orMap map[string][]orcond.Entry

	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	app := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaApplicationSignatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaApplicationSignatureConfig(dg, app, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationSignatureExists("panos_panorama_application_signature.test", &o, &andList, &orMap),
					testAccCheckPanosPanoramaApplicationSignatureAttributes(&o, name, &andList, &orMap),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaApplicationSignatureExists(n string, o *signature.Entry, andList *[]andcond.Entry, orMap *map[string][]orcond.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, app, name := parsePanoramaApplicationSignatureId(rs.Primary.ID)
		v, err := pano.Objects.AppSignature.Get(dg, app, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		aList, err := pano.Objects.AppSigAndCond.GetList(dg, app, name)
		if err != nil {
			return err
		}
		andObjList := make([]andcond.Entry, 0, len(aList))
		orObjMap := make(map[string][]orcond.Entry)
		for i := range aList {
			andObj, err := pano.Objects.AppSigAndCond.Get(dg, app, name, aList[i])
			if err != nil {
				return err
			}
			andObjList = append(andObjList, andObj)
			oList, err := pano.Objects.AppSigOrCond.GetList(dg, app, name, andObj.Name)
			if err != nil {
				return err
			}
			orObjList := make([]orcond.Entry, 0, len(oList))
			for j := range oList {
				orObj, err := pano.Objects.AppSigOrCond.Get(dg, app, name, andObj.Name, oList[j])
				if err != nil {
					return err
				}
				orObjList = append(orObjList, orObj)
			}
			orObjMap[andObj.Name] = orObjList
		}

		*o = v
		*andList = andObjList
		*orMap = orObjMap

		return nil
	}
}

func testAccCheckPanosPanoramaApplicationSignatureAttributes(o *signature.Entry, name string, andListP *[]andcond.Entry, orMapP *map[string][]orcond.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		andList := *andListP
		orMap := *orMapP

		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Comment != "my sig comment" {
			return fmt.Errorf("Comment is %q, expected 'my sig comment'", o.Comment)
		}
		if o.OrderFree {
			return fmt.Errorf("Order free is %t, expected false", o.OrderFree)
		}

		if len(andList) != 2 {
			return fmt.Errorf("And cond list is len %d, not 2", len(andList))
		}

		var ae andcond.Entry
		var oe orcond.Entry
		var ol []orcond.Entry

		ae = andList[0]
		if ae.Name != "And Condition 1" {
			return fmt.Errorf("Andcond 1 name is %q", ae.Name)
		}

		ol = orMap[ae.Name]
		if len(ol) != 3 {
			return fmt.Errorf("First or cond list is len %d, not 3", len(ol))
		}

		oe = ol[0]
		if oe.Name != "Or Condition 1" {
			return fmt.Errorf("AC1-OC1 name is %s", oe.Name)
		}
		if oe.Operator != orcond.OperatorPatternMatch {
			return fmt.Errorf("AC1-OC1 operator is %s", oe.Operator)
		}
		if oe.Context != "http-req-headers" {
			return fmt.Errorf("AC1-OC1 context is %s", oe.Context)
		}
		if oe.Pattern != "firstpattern" {
			return fmt.Errorf("AC1-OC1 pattern is %s", oe.Pattern)
		}

		oe = ol[1]
		if oe.Name != "Or Condition 2" {
			return fmt.Errorf("AC1-OC2 name is %s", oe.Name)
		}
		if oe.Operator != orcond.OperatorGreaterThan {
			return fmt.Errorf("AC1-OC2 operator is %s", oe.Operator)
		}
		if oe.Context != "cotp-req-x420-message-size" {
			return fmt.Errorf("AC1-OC2 context is %s", oe.Context)
		}
		if oe.Value != "123456" {
			return fmt.Errorf("AC1-OC2 value is %s", oe.Value)
		}

		oe = ol[2]
		if oe.Name != "Or Condition 3" {
			return fmt.Errorf("AC1-OC3 name is %s", oe.Name)
		}
		if oe.Operator != orcond.OperatorLessThan {
			return fmt.Errorf("AC1-OC3 operator is %s", oe.Operator)
		}
		if oe.Context != "cotp-req-x420-message-size" {
			return fmt.Errorf("AC1-OC3 context is %s", oe.Context)
		}
		if oe.Value != "42" {
			return fmt.Errorf("AC1-OC3 value is %s", oe.Value)
		}

		ae = andList[1]
		if ae.Name != "And Condition 2" {
			return fmt.Errorf("AC2 name is %s", ae.Name)
		}

		ol = orMap[ae.Name]
		if len(ol) != 1 {
			return fmt.Errorf("First or cond list is len %d, not 1", len(ol))
		}

		oe = ol[0]
		if oe.Name != "Or Condition 1" {
			return fmt.Errorf("AC2-OC1 name is %s", oe.Name)
		}
		if oe.Operator != orcond.OperatorEqualTo {
			return fmt.Errorf("AC2-OC1 operator is %s", oe.Operator)
		}
		if oe.Context != "unknown-req-tcp" {
			return fmt.Errorf("AC2-OC1 context is %s", oe.Context)
		}
		if oe.Position != "first-4bytes" {
			return fmt.Errorf("AC2-OC1 position is %s", oe.Position)
		}
		if oe.Mask != "0Xff112345" {
			return fmt.Errorf("AC2-OC1 mask is %s", oe.Mask)
		}
		if oe.Value != "0X11bb33dd" {
			return fmt.Errorf("AC2-OC1 value is %s", oe.Value)
		}

		return nil
	}
}

func testAccPanosPanoramaApplicationSignatureDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_application_signature" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, app, name := parsePanoramaApplicationSignatureId(rs.Primary.ID)
			_, err := pano.Objects.AppSignature.Get(dg, app, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaApplicationSignatureConfig(dg, app, name string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
    description = "for app sig test"
}

resource "panos_panorama_application_object" "x" {
    device_group = panos_panorama_device_group.x.name
    name = %q
    description = "application sig test"
    category = "media"
    subcategory = "gaming"
    technology = "client-server"
    risk = 5
}

resource "panos_panorama_application_signature" "test" {
    device_group = panos_panorama_device_group.x.name
    application_object = panos_panorama_application_object.x.name
    name = %q
    comment = "my sig comment"
    ordered_match = true
    and_condition {
        or_condition {
            pattern_match {
                context = "http-req-headers"
                pattern = "firstpattern"
                qualifiers = {
                    "http-method": "COPY",
                    "req-hdr-type": "HOST",
                }
            }
        }
        or_condition {
            greater_than {
                context = "cotp-req-x420-message-size"
                value = "123456"
            }
        }
        or_condition {
            less_than {
                context = "cotp-req-x420-message-size"
                value = "42"
            }
        }
    }
    and_condition {
        or_condition {
            equal_to {
                context = "unknown-req-tcp"
                position = "first-4bytes"
                mask = "0Xff112345"
                value = "0X11bb33dd"
            }
        }
    }
}
`, dg, app, name)
}
