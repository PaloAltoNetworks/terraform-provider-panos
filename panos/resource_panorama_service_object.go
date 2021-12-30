package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/objs/srvc"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaServiceObject() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaServiceObject,
		Read:   readPanoramaServiceObject,
		Update: updatePanoramaServiceObject,
		Delete: deletePanoramaServiceObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: serviceObjectSchema(true),
	}
}

func parsePanoramaServiceObject(d *schema.ResourceData) (string, srvc.Entry) {
	dg := d.Get("device_group").(string)
	o := loadServiceObject(d)

	return dg, o
}

func parsePanoramaServiceObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaServiceObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_service_object")
	if err != nil {
		return err
	}
	dg, o := parsePanoramaServiceObject(d)

	if err := pano.Objects.Services.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaServiceObjectId(dg, o.Name))
	return readPanoramaServiceObject(d, meta)
}

func readPanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_service_object")
	if err != nil {
		return err
	}
	dg, name := parsePanoramaServiceObjectId(d.Id())

	o, err := pano.Objects.Services.Get(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("device_group", dg)
	saveServiceObject(d, o)

	return nil
}

func updatePanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_service_object")
	if err != nil {
		return err
	}
	dg, o := parsePanoramaServiceObject(d)

	lo, err := pano.Objects.Services.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Services.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaServiceObject(d, meta)
}

func deletePanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_service_object")
	if err != nil {
		return err
	}
	dg, name := parsePanoramaServiceObjectId(d.Id())

	err = pano.Objects.Services.Delete(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
