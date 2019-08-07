package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/group"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaApplicationGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaApplicationGroup,
		Read:   readPanoramaApplicationGroup,
		Update: updatePanoramaApplicationGroup,
		Delete: deletePanoramaApplicationGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationGroupSchema(true),
	}
}

func parsePanoramaApplicationGroup(d *schema.ResourceData) (string, group.Entry) {
	dg := d.Get("device_group").(string)
	o := loadApplicationGroup(d)

	return dg, o
}

func parsePanoramaApplicationGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaApplicationGroupId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createPanoramaApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaApplicationGroup(d)

	if err := pano.Objects.AppGroup.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaApplicationGroupId(dg, o.Name))
	return readPanoramaApplicationGroup(d, meta)
}

func readPanoramaApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaApplicationGroupId(d.Id())

	o, err := pano.Objects.AppGroup.Get(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("device_group", dg)
	saveApplicationGroup(d, o)

	return nil
}

func updatePanoramaApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaApplicationGroup(d)

	lo, err := pano.Objects.AppGroup.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.AppGroup.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaApplicationGroup(d, meta)
}

func deletePanoramaApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaApplicationGroupId(d.Id())

	err := pano.Objects.AppGroup.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
