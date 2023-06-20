package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/aggregate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaBgpAggregate() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpAggregate,
		Read:   readPanoramaBgpAggregate,
		Update: updatePanoramaBgpAggregate,
		Delete: deletePanoramaBgpAggregate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpAggregateSchema(true),
	}
}

func parsePanoramaBgpAggregate(d *schema.ResourceData) (string, string, string, aggregate.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vr, o := parseBgpAggregate(d)

	return tmpl, ts, vr, o
}

func parsePanoramaBgpAggregateId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpAggregateId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpAggregate(d)

	if err = pano.Network.BgpAggregate.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpAggregateId(tmpl, ts, vr, o.Name))
	return readPanoramaBgpAggregate(d, meta)
}

func readPanoramaBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpAggregateId(d.Id())

	o, err := pano.Network.BgpAggregate.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	saveBgpAggregate(d, vr, o)

	return nil
}

func updatePanoramaBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpAggregate(d)

	lo, err := pano.Network.BgpAggregate.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpAggregate.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgpAggregate(d, meta)
}

func deletePanoramaBgpAggregate(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpAggregateId(d.Id())

	if err := pano.Network.BgpAggregate.Delete(tmpl, ts, vr, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
