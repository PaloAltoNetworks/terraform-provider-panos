package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/profile/dampening"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaBgpDampeningProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpDampeningProfile,
		Read:   readPanoramaBgpDampeningProfile,
		Update: updatePanoramaBgpDampeningProfile,
		Delete: deletePanoramaBgpDampeningProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpDampeningProfileSchema(true),
	}
}

func parsePanoramaBgpDampeningProfile(d *schema.ResourceData) (string, string, string, dampening.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgpDampeningProfile(d)

	return tmpl, ts, vr, o
}

func parsePanoramaBgpDampeningProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpDampeningProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpDampeningProfile(d)

	if err := pano.Network.BgpDampeningProfile.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpDampeningProfileId(tmpl, ts, vr, o.Name))
	return readPanoramaBgpDampeningProfile(d, meta)
}

func readPanoramaBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpDampeningProfileId(d.Id())

	o, err := pano.Network.BgpDampeningProfile.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpDampeningProfile(d, vr, o)

	return nil
}

func updatePanoramaBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpDampeningProfile(d)

	lo, err := pano.Network.BgpDampeningProfile.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpDampeningProfile.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgpDampeningProfile(d, meta)
}

func deletePanoramaBgpDampeningProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpDampeningProfileId(d.Id())

	err := pano.Network.BgpDampeningProfile.Delete(tmpl, ts, vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
