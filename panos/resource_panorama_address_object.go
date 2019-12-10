package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaAddressObject() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaAddressObject,
		Read:   readPanoramaAddressObject,
		Update: updatePanoramaAddressObject,
		Delete: deletePanoramaAddressObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: addressObjectSchema(true),
	}
}

func parsePanoramaAddressObject(d *schema.ResourceData) (string, addr.Entry) {
	dg := d.Get("device_group").(string)
	o := loadAddressObject(d)

	return dg, o
}

func parsePanoramaAddressObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaAddressObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaAddressObject(d)

	if err := pano.Objects.Address.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaAddressObjectId(dg, o.Name))
	return readPanoramaAddressObject(d, meta)
}

func readPanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaAddressObjectId(d.Id())

	o, err := pano.Objects.Address.Get(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("device_group", dg)
	saveAddressObject(d, o)

	return nil
}

func updatePanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaAddressObject(d)

	lo, err := pano.Objects.Address.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Address.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaAddressObject(d, meta)
}

func deletePanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaAddressObjectId(d.Id())

	err := pano.Objects.Address.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
