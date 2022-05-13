package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAdministrativeTag() *schema.Resource {
	return &schema.Resource{
		Create: createAdministrativeTag,
		Read:   readAdministrativeTag,
		Update: updateAdministrativeTag,
		Delete: deleteAdministrativeTag,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The administrative tag's name",
			},
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this administrative tag object in",
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parseAdministrativeTag(d *schema.ResourceData) (string, tags.Entry) {
	vsys := d.Get("vsys").(string)
	o := tags.Entry{
		Name:    d.Get("name").(string),
		Color:   d.Get("color").(string),
		Comment: d.Get("comment").(string),
	}

	return vsys, o
}

func parseAdministrativeTagId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildAdministrativeTagId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_administrative_tag")
	if err != nil {
		return err
	}

	vsys, o := parseAdministrativeTag(d)

	if err = fw.Objects.Tags.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildAdministrativeTagId(vsys, o.Name))
	return readAdministrativeTag(d, meta)
}

func readAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_administrative_tag")
	if err != nil {
		return err
	}

	vsys, name := parseAdministrativeTagId(d.Id())

	o, err := fw.Objects.Tags.Get(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("color", o.Color)
	d.Set("comment", o.Comment)

	return nil
}

func updateAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_administrative_tag")
	if err != nil {
		return err
	}

	vsys, o := parseAdministrativeTag(d)

	lo, err := fw.Objects.Tags.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.Tags.Edit(vsys, lo); err != nil {
		return err
	}

	return readAdministrativeTag(d, meta)
}

func deleteAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_administrative_tag")
	if err != nil {
		return err
	}

	vsys, name := parseAdministrativeTagId(d.Id())

	if err = fw.Objects.Tags.Delete(vsys, name); err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}
