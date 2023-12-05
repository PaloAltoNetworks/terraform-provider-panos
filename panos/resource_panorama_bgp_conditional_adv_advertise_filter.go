package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv/filter/advertise"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaBgpConditionalAdvAdvertiseFilter() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpConditionalAdvAdvertiseFilter,
		Read:   readPanoramaBgpConditionalAdvAdvertiseFilter,
		Update: updatePanoramaBgpConditionalAdvAdvertiseFilter,
		Delete: deletePanoramaBgpConditionalAdvAdvertiseFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvAdvertiseFilterSchema(true),
	}
}

func parsePanoramaBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData) (string, string, string, string, advertise.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, ca, o := parseBgpConditionalAdvAdvertiseFilter(d)

	return tmpl, ts, vr, ca, o
}

func parsePanoramaBgpConditionalAdvAdvertiseFilterId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaBgpConditionalAdvAdvertiseFilterId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createPanoramaBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, o := parsePanoramaBgpConditionalAdvAdvertiseFilter(d)

	if err = pano.Network.BgpConAdvAdvertiseFilter.Set(tmpl, ts, vr, ca, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpConditionalAdvAdvertiseFilterId(tmpl, ts, vr, ca, o.Name))
	return readPanoramaBgpConditionalAdvAdvertiseFilter(d, meta)
}

func readPanoramaBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, name := parsePanoramaBgpConditionalAdvAdvertiseFilterId(d.Id())

	o, err := pano.Network.BgpConAdvAdvertiseFilter.Get(tmpl, ts, vr, ca, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpConditionalAdvAdvertiseFilter(d, vr, ca, o)

	return nil
}

func updatePanoramaBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, o := parsePanoramaBgpConditionalAdvAdvertiseFilter(d)

	lo, err := pano.Network.BgpConAdvAdvertiseFilter.Get(tmpl, ts, vr, ca, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpConAdvAdvertiseFilter.Edit(tmpl, ts, vr, ca, lo); err != nil {
		return err
	}

	return readPanoramaBgpConditionalAdvAdvertiseFilter(d, meta)
}

func deletePanoramaBgpConditionalAdvAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, name := parsePanoramaBgpConditionalAdvAdvertiseFilterId(d.Id())

	if err := pano.Network.BgpConAdvAdvertiseFilter.Delete(tmpl, ts, vr, ca, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
