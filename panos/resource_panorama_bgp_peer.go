package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaBgpPeer,
		Read:   readPanoramaBgpPeer,
		Update: updatePanoramaBgpPeer,
		Delete: deletePanoramaBgpPeer,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpPeerSchema(true),
	}
}

func parsePanoramaBgpPeer(d *schema.ResourceData) (string, string, string, string, peer.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, pg, o := parseBgpPeer(d)

	return tmpl, ts, vr, pg, o
}

func parsePanoramaBgpPeerId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaBgpPeerId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createPanoramaBgpPeer(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, pg, o := parsePanoramaBgpPeer(d)

	if err := pano.Network.BgpPeer.Set(tmpl, ts, vr, pg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaBgpPeerId(tmpl, ts, vr, pg, o.Name))
	return readPanoramaBgpPeer(d, meta)
}

func readPanoramaBgpPeer(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, pg, name := parsePanoramaBgpPeerId(d.Id())

	o, err := pano.Network.BgpPeer.Get(tmpl, ts, vr, pg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	saveBgpPeer(d, vr, pg, o)

	return nil
}

func updatePanoramaBgpPeer(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, pg, o := parsePanoramaBgpPeer(d)

	lo, err := pano.Network.BgpPeer.Get(tmpl, ts, vr, pg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.BgpPeer.Edit(tmpl, ts, vr, pg, lo); err != nil {
		return err
	}

	return readPanoramaBgpPeer(d, meta)
}

func deletePanoramaBgpPeer(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, pg, name := parsePanoramaBgpPeerId(d.Id())

	err := pano.Network.BgpPeer.Delete(tmpl, ts, vr, pg, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
