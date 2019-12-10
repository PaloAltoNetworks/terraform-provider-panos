package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/redist"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaBgpRedistRule() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpRedistRule,
		Read:   readPanoramaBgpRedistRule,
		Update: updatePanoramaBgpRedistRule,
		Delete: deletePanoramaBgpRedistRule,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpRedistRuleSchema(true),
	}
}

func parsePanoramaBgpRedistRule(d *schema.ResourceData) (string, string, string, redist.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgpRedistRule(d)

	return tmpl, ts, vr, o
}

func parsePanoramaBgpRedistRuleId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpRedistRuleId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpRedistRule(d)

	if err = pano.Network.BgpRedistRule.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpRedistRuleId(tmpl, ts, vr, o.Name))
	return readPanoramaBgpRedistRule(d, meta)
}

func readPanoramaBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpRedistRuleId(d.Id())

	o, err := pano.Network.BgpRedistRule.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpRedistRule(d, vr, o)

	return nil
}

func updatePanoramaBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpRedistRule(d)

	lo, err := pano.Network.BgpRedistRule.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpRedistRule.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgpRedistRule(d, meta)
}

func deletePanoramaBgpRedistRule(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpRedistRuleId(d.Id())

	if err := pano.Network.BgpRedistRule.Delete(tmpl, ts, vr, name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
