package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAdministrativeTag() *schema.Resource {
	return &schema.Resource{
		Create: createAdministrativeTag,
		Read:   readAdministrativeTag,
		Update: updateAdministrativeTag,
		Delete: deleteAdministrativeTag,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The administrative tag's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this administrative tag object in",
			},
			"color": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": &schema.Schema{
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
	fw := meta.(*pango.Firewall)
	vsys, o := parseAdministrativeTag(d)

	if err := fw.Objects.Tags.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildAdministrativeTagId(vsys, o.Name))
	return readAdministrativeTag(d, meta)
}

func readAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseAdministrativeTagId(d.Id())

	o, err := fw.Objects.Tags.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
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
	var err error
	fw := meta.(*pango.Firewall)
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
	fw := meta.(*pango.Firewall)
	vsys, name := parseAdministrativeTagId(d.Id())

	err := fw.Objects.Tags.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
