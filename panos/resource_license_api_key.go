package panos

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLicenseApiKey() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateLicenseApiKey,
		Read:   readLicenseApiKey,
		Update: createUpdateLicenseApiKey,
		Delete: deleteLicenseApiKey,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"retain_key": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func createUpdateLicenseApiKey(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)

	key := d.Get("key").(string)
	keep := d.Get("retain_key").(bool)

	if err := fw.Licensing.SetApiKey(key); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%t", keep))
	return readLicenseApiKey(d, meta)
}

func readLicenseApiKey(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	keep := d.Get("retain_key").(bool)

	key, err := fw.Licensing.GetApiKey()
	if err != nil {
		return err
	} else if key == "" {
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%t", keep))
	d.Set("key", key)
	d.Set("retain_key", keep)

	return nil
}

func deleteLicenseApiKey(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	keep := d.Id()

	if keep == "false" {
		if err := fw.Licensing.DeleteApiKey(); err != nil {
			return err
		}
	}

	return nil
}
