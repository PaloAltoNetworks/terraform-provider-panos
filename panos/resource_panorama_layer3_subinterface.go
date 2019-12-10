package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaLayer3Subinterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaLayer3Subinterface,
		Read:   readPanoramaLayer3Subinterface,
		Update: updatePanoramaLayer3Subinterface,
		Delete: deletePanoramaLayer3Subinterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: layer3SubinterfaceSchema(true),
	}
}

func parsePanoramaLayer3SubinterfaceId(v string) (string, string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4], t[5]
}

func buildPanoramaLayer3SubinterfaceId(a, b, c, d, e, f string) string {
	return strings.Join([]string{a, b, c, d, e, f}, IdSeparator)
}

func parsePanoramaLayer3Subinterface(d *schema.ResourceData) (string, string, string, string, string, layer3.Entry) {
	tmpl := d.Get("template").(string)
	ts := ""
	iType := d.Get("interface_type").(string)
	eth := d.Get("parent_interface").(string)
	vsys := d.Get("vsys").(string)
	o := loadLayer3Subinterface(d)

	return tmpl, ts, iType, eth, vsys, o
}

func createPanoramaLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, vsys, o := parsePanoramaLayer3Subinterface(d)

	if err := pano.Network.Layer3Subinterface.Set(tmpl, ts, iType, eth, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaLayer3SubinterfaceId(tmpl, ts, iType, eth, vsys, o.Name))
	return readPanoramaLayer3Subinterface(d, meta)
}

func readPanoramaLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, vsys, name := parsePanoramaLayer3SubinterfaceId(d.Id())

	o, err := pano.Network.Layer3Subinterface.Get(tmpl, ts, iType, eth, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := pano.IsImported(util.InterfaceImport, tmpl, ts, vsys, name)
	if err != nil {
		return err
	}

	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("template", tmpl)
	d.Set("interface_type", iType)
	d.Set("parent_interface", eth)
	saveLayer3Subinterface(d, o)

	return nil
}

func updatePanoramaLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, vsys, o := parsePanoramaLayer3Subinterface(d)

	lo, err := pano.Network.Layer3Subinterface.Get(tmpl, ts, iType, eth, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.Layer3Subinterface.Edit(tmpl, ts, iType, eth, vsys, lo); err != nil {
		return err
	}

	d.SetId(buildPanoramaLayer3SubinterfaceId(tmpl, ts, iType, eth, vsys, o.Name))
	return readPanoramaLayer3Subinterface(d, meta)
}

func deletePanoramaLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, _, name := parsePanoramaLayer3SubinterfaceId(d.Id())

	err := pano.Network.Layer3Subinterface.Delete(tmpl, ts, iType, eth, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
