package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaBgpConditionalAdv() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpConditionalAdv,
		Read:   readPanoramaBgpConditionalAdv,
		Update: updatePanoramaBgpConditionalAdv,
		Delete: deletePanoramaBgpConditionalAdv,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvSchema(true),
	}
}

func parsePanoramaBgpConditionalAdvId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpConditionalAdvId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaBgpConditionalAdv(d *schema.ResourceData) (string, string, string, conadv.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgpConditionalAdv(d)

	return tmpl, ts, vr, o
}

func createPanoramaBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpConditionalAdv(d)

	if err = pano.Network.BgpConditionalAdv.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpConditionalAdvId(tmpl, ts, vr, o.Name))
	return readPanoramaBgpConditionalAdv(d, meta)
}

func readPanoramaBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpConditionalAdvId(d.Id())

	o, err := pano.Network.BgpConditionalAdv.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpConditionalAdv(d, vr, o)

	return nil
}

func updatePanoramaBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpConditionalAdv(d)

	lo, err := pano.Network.BgpConditionalAdv.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpConditionalAdv.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgpConditionalAdv(d, meta)
}

func deletePanoramaBgpConditionalAdv(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpConditionalAdvId(d.Id())

	err := pano.Network.BgpConditionalAdv.Delete(tmpl, ts, vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
