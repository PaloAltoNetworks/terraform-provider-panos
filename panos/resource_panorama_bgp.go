package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaBgp() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgp,
		Read:   readPanoramaBgp,
		Update: updatePanoramaBgp,
		Delete: deletePanoramaBgp,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpSchema(true),
	}
}

func parsePanoramaBgpId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaBgpId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parsePanoramaBgp(d *schema.ResourceData) (string, string, string, bgp.Config) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgp(d)

	return tmpl, ts, vr, o
}

func createPanoramaBgp(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgp(d)

	if err = pano.Network.BgpConfig.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpId(tmpl, ts, vr))
	return readPanoramaBgp(d, meta)
}

func readPanoramaBgp(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr := parsePanoramaBgpId(d.Id())

	o, err := pano.Network.BgpConfig.Get(tmpl, ts, vr)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	saveBgp(d, vr, o)

	return nil
}

func updatePanoramaBgp(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgp(d)

	lo, err := pano.Network.BgpConfig.Get(tmpl, ts, vr)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpConfig.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgp(d, meta)
}

func deletePanoramaBgp(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr := parsePanoramaBgpId(d.Id())

	err := pano.Network.BgpConfig.Delete(tmpl, ts, vr)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
