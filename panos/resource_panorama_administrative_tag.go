package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaAdministrativeTag() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaAdministrativeTag,
		Read:   readPanoramaAdministrativeTag,
		Update: updatePanoramaAdministrativeTag,
		Delete: deletePanoramaAdministrativeTag,

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
			"device_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "shared",
				ForceNew:    true,
				Description: "The device group to put this administrative tag object in",
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

func parsePanoramaAdministrativeTag(d *schema.ResourceData) (string, tags.Entry) {
	dg := d.Get("device_group").(string)
	o := tags.Entry{
		Name:    d.Get("name").(string),
		Color:   d.Get("color").(string),
		Comment: d.Get("comment").(string),
	}

	return dg, o
}

func parsePanoramaAdministrativeTagId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaAdministrativeTagId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_administrative_tag")
	if err != nil {
		return err
	}

	dg, o := parsePanoramaAdministrativeTag(d)

	if err = pano.Objects.Tags.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaAdministrativeTagId(dg, o.Name))
	return readPanoramaAdministrativeTag(d, meta)
}

func readPanoramaAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_administrative_tag")
	if err != nil {
		return err
	}

	dg, name := parsePanoramaAdministrativeTagId(d.Id())

	o, err := pano.Objects.Tags.Get(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("device_group", dg)
	d.Set("color", o.Color)
	d.Set("comment", o.Comment)

	return nil
}

func updatePanoramaAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_administrative_tag")
	if err != nil {
		return err
	}

	dg, o := parsePanoramaAdministrativeTag(d)

	lo, err := pano.Objects.Tags.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Tags.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaAdministrativeTag(d, meta)
}

func deletePanoramaAdministrativeTag(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_administrative_tag")
	if err != nil {
		return err
	}

	dg, name := parsePanoramaAdministrativeTagId(d.Id())

	if err = pano.Objects.Tags.Delete(dg, name); err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}
