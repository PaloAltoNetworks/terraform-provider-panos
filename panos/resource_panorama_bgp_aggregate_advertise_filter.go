package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/advertise"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaBgpAggregateAdvertiseFilter() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpAggregateAdvertiseFilter,
		Read:   readPanoramaBgpAggregateAdvertiseFilter,
		Update: updatePanoramaBgpAggregateAdvertiseFilter,
		Delete: deletePanoramaBgpAggregateAdvertiseFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAggregateAdvertiseFilterSchema(true),
	}
}

func parsePanoramaBgpAggregateAdvertiseFilter(d *schema.ResourceData) (string, string, string, string, advertise.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, ag, o := parseBgpAggregateAdvertiseFilter(d)

	return tmpl, ts, vr, ag, o
}

func parsePanoramaBgpAggregateAdvertiseFilterId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaBgpAggregateAdvertiseFilterId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createPanoramaBgpAggregateAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, o := parsePanoramaBgpAggregateAdvertiseFilter(d)

	if err = pano.Network.BgpAggAdvertiseFilter.Set(tmpl, ts, vr, ag, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpAggregateAdvertiseFilterId(tmpl, ts, vr, ag, o.Name))
	return readPanoramaBgpAggregateAdvertiseFilter(d, meta)
}

func readPanoramaBgpAggregateAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateAdvertiseFilterId(d.Id())

	o, err := pano.Network.BgpAggAdvertiseFilter.Get(tmpl, ts, vr, ag, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpAggregateAdvertiseFilter(d, vr, ag, o)

	return nil
}

func updatePanoramaBgpAggregateAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, o := parsePanoramaBgpAggregateAdvertiseFilter(d)

	lo, err := pano.Network.BgpAggAdvertiseFilter.Get(tmpl, ts, vr, ag, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpAggAdvertiseFilter.Edit(tmpl, ts, vr, ag, lo); err != nil {
		return err
	}

	return readPanoramaBgpAggregateAdvertiseFilter(d, meta)
}

func deletePanoramaBgpAggregateAdvertiseFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateAdvertiseFilterId(d.Id())

	err := pano.Network.BgpAggAdvertiseFilter.Delete(tmpl, ts, vr, ag, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
