package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer/group"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaBgpPeerGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpPeerGroup,
		Read:   readPanoramaBgpPeerGroup,
		Update: updatePanoramaBgpPeerGroup,
		Delete: deletePanoramaBgpPeerGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpPeerGroupSchema(true),
	}
}

func parsePanoramaBgpPeerGroup(d *schema.ResourceData) (string, string, string, group.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseBgpPeerGroup(d)

	return tmpl, ts, vr, o
}

func parsePanoramaBgpPeerGroupId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaBgpPeerGroupId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpPeerGroup(d)

	if err := pano.Network.BgpPeerGroup.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpPeerGroupId(tmpl, ts, vr, o.Name))
	return readPanoramaBgpPeerGroup(d, meta)
}

func readPanoramaBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpPeerGroupId(d.Id())

	o, err := pano.Network.BgpPeerGroup.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveBgpPeerGroup(d, vr, o)

	return nil
}

func updatePanoramaBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaBgpPeerGroup(d)

	lo, err := pano.Network.BgpPeerGroup.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpPeerGroup.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaBgpPeerGroup(d, meta)
}

func deletePanoramaBgpPeerGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaBgpPeerGroupId(d.Id())

	err := pano.Network.BgpPeerGroup.Delete(tmpl, ts, vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
