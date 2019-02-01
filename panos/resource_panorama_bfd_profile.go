package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/bfd"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaBfdProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBfdProfile,
		Read:   readPanoramaBfdProfile,
		Update: updatePanoramaBfdProfile,
		Delete: deletePanoramaBfdProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bfdProfileSchema(true),
	}
}

func parsePanoramaBfdProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaBfdProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parsePanoramaBfdProfile(d *schema.ResourceData) (string, string, bfd.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	o := parseBfdProfile(d)

	return tmpl, ts, o
}

func createPanoramaBfdProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaBfdProfile(d)

	if err := pano.Network.BfdProfile.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBfdProfileId(tmpl, ts, o.Name))
	return readPanoramaBfdProfile(d, meta)
}

func readPanoramaBfdProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaBfdProfileId(d.Id())

	o, err := pano.Network.BfdProfile.Get(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBfdProfile(d, o)

	return nil
}

func updatePanoramaBfdProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaBfdProfile(d)

	lo, err := pano.Network.BfdProfile.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BfdProfile.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaBfdProfile(d, meta)
}

func deletePanoramaBfdProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaBfdProfileId(d.Id())

	err := pano.Network.BfdProfile.Delete(tmpl, ts, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
