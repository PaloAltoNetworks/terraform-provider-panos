package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/suppress"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaBgpAggregateSuppressFilter() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpAggregateSuppressFilter,
		Read:   readPanoramaBgpAggregateSuppressFilter,
		Update: updatePanoramaBgpAggregateSuppressFilter,
		Delete: deletePanoramaBgpAggregateSuppressFilter,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAggregateSuppressFilterSchema(true),
	}
}

func parsePanoramaBgpAggregateSuppressFilter(d *schema.ResourceData) (string, string, string, string, suppress.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, ag, o := parseBgpAggregateSuppressFilter(d)

	return tmpl, ts, vr, ag, o
}

func parsePanoramaBgpAggregateSuppressFilterId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaBgpAggregateSuppressFilterId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createPanoramaBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, o := parsePanoramaBgpAggregateSuppressFilter(d)

	if err = pano.Network.BgpAggSuppressFilter.Set(tmpl, ts, vr, ag, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpAggregateSuppressFilterId(tmpl, ts, vr, ag, o.Name))
	return readPanoramaBgpAggregateSuppressFilter(d, meta)
}

func readPanoramaBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateSuppressFilterId(d.Id())

	o, err := pano.Network.BgpAggSuppressFilter.Get(tmpl, ts, vr, ag, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpAggregateSuppressFilter(d, vr, ag, o)

	return nil
}

func updatePanoramaBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, o := parsePanoramaBgpAggregateSuppressFilter(d)

	lo, err := pano.Network.BgpAggSuppressFilter.Get(tmpl, ts, vr, ag, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpAggSuppressFilter.Edit(tmpl, ts, vr, ag, lo); err != nil {
		return err
	}

	return readPanoramaBgpAggregateSuppressFilter(d, meta)
}

func deletePanoramaBgpAggregateSuppressFilter(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateSuppressFilterId(d.Id())

	err := pano.Network.BgpAggSuppressFilter.Delete(tmpl, ts, vr, ag, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
