package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/profile/redist/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaRedistributionProfileIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaRedistributionProfileIpv4,
		Read:   readPanoramaRedistributionProfileIpv4,
		Update: updatePanoramaRedistributionProfileIpv4,
		Delete: deletePanoramaRedistributionProfileIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: redistributionProfileIpv4Schema(true),
	}
}

func parsePanoramaRedistributionProfileIpv4(d *schema.ResourceData) (string, string, string, ipv4.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	vr, o := parseRedistributionProfileIpv4(d)

	return tmpl, ts, vr, o
}

func parsePanoramaRedistributionProfileIpv4Id(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaRedistributionProfileIpv4Id(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaRedistributionProfileIpv4(d)

	if err := pano.Network.RedistributionProfile.Set(tmpl, ts, vr, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaRedistributionProfileIpv4Id(tmpl, ts, vr, o.Name))
	return readPanoramaRedistributionProfileIpv4(d, meta)
}

func readPanoramaRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaRedistributionProfileIpv4Id(d.Id())

	o, err := pano.Network.RedistributionProfile.Get(tmpl, ts, vr, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveRedistributionProfileIpv4(d, vr, o)

	return nil
}

func updatePanoramaRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, o := parsePanoramaRedistributionProfileIpv4(d)

	lo, err := pano.Network.RedistributionProfile.Get(tmpl, ts, vr, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.RedistributionProfile.Edit(tmpl, ts, vr, lo); err != nil {
		return err
	}

	return readPanoramaRedistributionProfileIpv4(d, meta)
}

func deletePanoramaRedistributionProfileIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vr, name := parsePanoramaRedistributionProfileIpv4Id(d.Id())

	err := pano.Network.RedistributionProfile.Delete(tmpl, ts, vr, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
