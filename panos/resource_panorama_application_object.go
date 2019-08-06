package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaApplicationObject() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaApplicationObject,
		Read:   readPanoramaApplicationObject,
		Update: updatePanoramaApplicationObject,
		Delete: deletePanoramaApplicationObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationObjectSchema(true),
	}
}

func parsePanoramaApplicationObject(d *schema.ResourceData) (string, app.Entry) {
	dg := d.Get("device_group").(string)
	o := loadApplicationObject(d)

	return dg, o
}

func parsePanoramaApplicationObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaApplicationObjectId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createPanoramaApplicationObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaApplicationObject(d)

	if err := pano.Objects.Application.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaApplicationObjectId(dg, o.Name))
	return readPanoramaApplicationObject(d, meta)
}

func readPanoramaApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaApplicationObjectId(d.Id())

	o, err := pano.Objects.Application.Get(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("device_group", dg)
	saveApplicationObject(d, o)

	return nil
}

func updatePanoramaApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaApplicationObject(d)

	lo, err := pano.Objects.Application.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Application.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaApplicationObject(d, meta)
}

func deletePanoramaApplicationObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaApplicationObjectId(d.Id())

	err := pano.Objects.Application.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
