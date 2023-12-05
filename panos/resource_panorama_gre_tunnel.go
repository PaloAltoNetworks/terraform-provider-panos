package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/tunnel/gre"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePanoramaGreTunnel() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaGreTunnel,
		Read:   readPanoramaGreTunnel,
		Update: updatePanoramaGreTunnel,
		Delete: deletePanoramaGreTunnel,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: greTunnelSchema(true),
	}
}

func buildPanoramaGreTunnelId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parsePanoramaGreTunnelId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func parsePanoramaGreTunnel(d *schema.ResourceData) (string, string, gre.Entry) {
	tmpl := d.Get("template").(string)
	ts := ""
	o := loadGreTunnel(d)

	return tmpl, ts, o
}

func createPanoramaGreTunnel(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaGreTunnel(d)

	if err := pano.Network.GreTunnel.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaGreTunnelId(tmpl, ts, o.Name))
	return readPanoramaGreTunnel(d, meta)
}

func readPanoramaGreTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaGreTunnelId(d.Id())

	o, err := pano.Network.GreTunnel.Get(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	saveGreTunnel(d, o)

	return nil
}

func updatePanoramaGreTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaGreTunnel(d)

	lo, err := pano.Network.GreTunnel.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.GreTunnel.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaGreTunnel(d, meta)
}

func deletePanoramaGreTunnel(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaGreTunnelId(d.Id())

	err := pano.Network.GreTunnel.Delete(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
