package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaVlanEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaVlanEntry,
		Read:   readPanoramaVlanEntry,
		Update: createUpdatePanoramaVlanEntry,
		Delete: deletePanoramaVlanEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: vlanEntrySchema(true),
	}
}

func parsePanoramaVlanEntry(d *schema.ResourceData) (string, string, string, string, []string, []string) {
	tmpl := d.Get("template").(string)
	ts := ""
	vlan, iface, rmMacs, addMacs := parseVlanEntry(d)

	return tmpl, ts, vlan, iface, rmMacs, addMacs
}

func parsePanoramaVlanEntryId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaVlanEntryId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createUpdatePanoramaVlanEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vlan, iface, rmMacs, addMacs := parsePanoramaVlanEntry(d)

	if err := pano.Network.Vlan.SetInterface(tmpl, ts, vlan, iface, rmMacs, addMacs); err != nil {
		return err
	}

	d.SetId(buildPanoramaVlanEntryId(tmpl, ts, vlan, iface))
	return readPanoramaVlanEntry(d, meta)
}

func readPanoramaVlanEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vlan, iface := parsePanoramaVlanEntryId(d.Id())

	// Two possibilities:  either the router itself doesn't exist or the
	// interface isn't present.
	o, err := pano.Network.Vlan.Get(tmpl, ts, vlan)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	found := false
	for i := range o.Interfaces {
		if o.Interfaces[i] == iface {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	macs := make([]string, 0, len(o.StaticMacs))
	for k, v := range o.StaticMacs {
		if v == iface {
			macs = append(macs, k)
		}
	}

	saveVlanEntry(d, vlan, iface, macs)
	d.Set("template", tmpl)

	return nil
}

func deletePanoramaVlanEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vlan, iface := parsePanoramaVlanEntryId(d.Id())

	if err := pano.Network.VirtualRouter.DeleteInterface(tmpl, ts, vlan, iface); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
