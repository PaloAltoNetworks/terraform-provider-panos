package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/auth"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaBgpAuthProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpAuthProfile,
		Read:   readPanoramaBgpAuthProfile,
		Update: updatePanoramaBgpAuthProfile,
		Delete: deletePanoramaBgpAuthProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAuthProfileSchema(true),
	}
}

func parsePanoramaBgpAuthProfile(d *schema.ResourceData) (string, string, string, auth.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgpAuthProfile(d)

	return tmpl, ts, vr, o
}

func parsePanoramaBgpAuthProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpAuthProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpAuthProfile(d)

	if err := pano.Network.BgpAuthProfile.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	lo, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpAuthProfileId(tmpl, ts, vr, o.Name))
	d.Set("secret_enc", lo.Secret)

	return readPanoramaBgpAuthProfile(d, meta)
}

func readPanoramaBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpAuthProfileId(d.Id())

	o, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpAuthProfile(d, vr, o)

	return nil
}

func updatePanoramaBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpAuthProfile(d)

	lo, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpAuthProfile.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	eo, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}

	d.Set("secret_enc", eo.Secret)
	return readPanoramaBgpAuthProfile(d, meta)
}

func deletePanoramaBgpAuthProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpAuthProfileId(d.Id())

	err := pano.Network.BgpAuthProfile.Delete(tmpl, ts, vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
