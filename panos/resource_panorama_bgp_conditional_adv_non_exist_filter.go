package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv/filter/nonexist"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaBgpConditionalAdvNonExistFilter() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpConditionalAdvNonExistFilter,
		Read:   readPanoramaBgpConditionalAdvNonExistFilter,
		Update: updatePanoramaBgpConditionalAdvNonExistFilter,
		Delete: deletePanoramaBgpConditionalAdvNonExistFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpConditionalAdvNonExistFilterSchema(true),
	}
}

func parsePanoramaBgpConditionalAdvNonExistFilter(d *schema.ResourceData) (string, string, string, string, nonexist.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, ca, o := parseBgpConditionalAdvNonExistFilter(d)

	return tmpl, ts, vr, ca, o
}

func parsePanoramaBgpConditionalAdvNonExistFilterId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaBgpConditionalAdvNonExistFilterId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createPanoramaBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, o := parsePanoramaBgpConditionalAdvNonExistFilter(d)

	if err = pano.Network.BgpConAdvNonExistFilter.Set(tmpl, ts, vr, ca, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpConditionalAdvNonExistFilterId(tmpl, ts, vr, ca, o.Name))
	return readPanoramaBgpConditionalAdvNonExistFilter(d, meta)
}

func readPanoramaBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, name := parsePanoramaBgpConditionalAdvNonExistFilterId(d.Id())

	o, err := pano.Network.BgpConAdvNonExistFilter.Get(tmpl, ts, vr, ca, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpConditionalAdvNonExistFilter(d, vr, ca, o)

	return nil
}

func updatePanoramaBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, o := parsePanoramaBgpConditionalAdvNonExistFilter(d)

	lo, err := pano.Network.BgpConAdvNonExistFilter.Get(tmpl, ts, vr, ca, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpConAdvNonExistFilter.Edit(tmpl, ts, vr, ca, lo); err != nil {
		return err
	}

	return readPanoramaBgpConditionalAdvNonExistFilter(d, meta)
}

func deletePanoramaBgpConditionalAdvNonExistFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ca, name := parsePanoramaBgpConditionalAdvNonExistFilterId(d.Id())

	if err := pano.Network.BgpConAdvNonExistFilter.Delete(tmpl, ts, vr, ca, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
