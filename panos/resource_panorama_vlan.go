package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/vlan"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaVlan() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaVlan,
		Read:   readPanoramaVlan,
		Update: updatePanoramaVlan,
		Delete: deletePanoramaVlan,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: vlanSchema(true),
	}
}

func parsePanoramaVlan(d *schema.ResourceData) (string, string, string, vlan.Entry) {
	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := ""
	o := loadVlan(d)

	return tmpl, ts, vsys, o
}

func parsePanoramaVlanId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaVlanId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaVlan(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVlan(d)

	if err := pano.Network.Vlan.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaVlanId(tmpl, ts, vsys, o.Name))
	return readPanoramaVlan(d, meta)
}

func readPanoramaVlan(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaVlanId(d.Id())

	o, err := pano.Network.Vlan.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	rv, err := pano.IsImported(util.VlanImport, tmpl, ts, vsys, name)
	if err != nil {
		return err
	}

	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	saveVlan(d, o)

	return nil
}

func updatePanoramaVlan(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaVlan(d)

	lo, err := pano.Network.Vlan.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o, false)
	if err = pano.Network.Vlan.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaVlan(d, meta)
}

func deletePanoramaVlan(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaVlanId(d.Id())

	err := pano.Network.Vlan.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
