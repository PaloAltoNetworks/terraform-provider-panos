package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLicensing() *schema.Resource {
	return &schema.Resource{
		Create: createLicensing,
		Read:   readLicensing,
		Update: updateLicensing,
		Delete: deleteLicensing,

		Schema: map[string]*schema.Schema{
			"auth_codes": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"delicense": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "auto",
				ValidateFunc: validateStringIn("auto"),
			},
			"licenses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issued": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_code": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func buildLicensingId(a string, b bool) string {
	return fmt.Sprintf("%s%s%t", a, IdSeparator, b)
}

func parseLicensingId(v string) (string, bool) {
	t := strings.Split(v, IdSeparator)
	b := t[1] == "true"
	return t[0], b
}

func createLicensing(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)

	mode := d.Get("mode").(string)
	delicense := d.Get("delicense").(bool)
	codes := asStringList(d.Get("auth_codes").([]interface{}))

	for _, v := range codes {
		if err = fw.Licensing.Activate(v); err != nil {
			return fmt.Errorf("Failed to activate %q: %s", v, err)
		}
	}

	d.SetId(buildLicensingId(mode, delicense))
	return readLicensing(d, meta)
}

func updateLicensing(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)

	if d.HasChange("auth_codes") {
		o, n := d.GetChange("auth_codes")
		old := asStringList(o.([]interface{}))
		new := asStringList(n.([]interface{}))
		missing := make([]string, 0, len(new))
		for i := 0; i < len(new); i++ {
			ran := false
			for j := 0; j < len(old); j++ {
				if new[i] == old[j] {
					ran = true
					break
				}
			}

			if ran {
				continue
			}

			missing = append(missing, new[i])
		}

		for _, v := range missing {
			if err = fw.Licensing.Activate(v); err != nil {
				return fmt.Errorf("Failed to activate new code %q: %s", v, err)
			}
		}
	}

	return readLicensing(d, meta)
}

func readLicensing(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)

	mode := d.Get("mode").(string)
	delicense := d.Get("delicense").(bool)
	codes := asStringList(d.Get("auth_codes").([]interface{}))

	list, err := fw.Licensing.Current()
	if err != nil {
		return err
	} else if len(list) == 0 {
		d.SetId("")
		return nil
	}

	d.SetId(buildLicensingId(mode, delicense))
	d.Set("mode", mode)
	d.Set("delicense", delicense)
	d.Set("auth_codes", codes)

	ilist := make([]interface{}, 0, len(list))
	for i := range list {
		m := make(map[string]interface{})
		m["feature"] = list[i].Feature
		m["description"] = list[i].Description
		m["serial"] = list[i].Serial
		m["issued"] = list[i].Issued
		m["expires"] = list[i].Expires
		m["expired"] = list[i].Expired
		m["auth_code"] = list[i].AuthCode

		ilist = append(ilist, m)
	}

	if err = d.Set("licenses", ilist); err != nil {
		log.Printf("[WARN] Error setting 'licenses' param for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteLicensing(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)

	_, delicense := parseLicensingId(d.Id())

	if delicense {
		if err := fw.Licensing.Deactivate(); err != nil {
			return err
		}
	}

	return nil
}
