package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer2"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaLayer2Subinterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaLayer2Subinterface,
		Read:   readPanoramaLayer2Subinterface,
		Update: updatePanoramaLayer2Subinterface,
		Delete: deletePanoramaLayer2Subinterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: layer2SubinterfaceSchema(true),
	}
}

func parsePanoramaLayer2SubinterfaceId(v string) (string, string, string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4], t[5], t[6]
}

func buildPanoramaLayer2SubinterfaceId(a, b, c, d, e, f, g string) string {
	return strings.Join([]string{a, b, c, d, e, f, g}, IdSeparator)
}

func parsePanoramaLayer2Subinterface(d *schema.ResourceData) (string, string, string, string, string, string, layer2.Entry) {
	tmpl := d.Get("template").(string)
	ts := ""
	iType := d.Get("interface_type").(string)
	eth := d.Get("parent_interface").(string)
	mType := d.Get("parent_mode").(string)
	vsys := d.Get("vsys").(string)
	o := loadLayer2Subinterface(d)

	return tmpl, ts, iType, eth, mType, vsys, o
}

func createPanoramaLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, mType, vsys, o := parsePanoramaLayer2Subinterface(d)

	if err := pano.Network.Layer2Subinterface.Set(tmpl, ts, iType, eth, mType, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaLayer2SubinterfaceId(tmpl, ts, iType, eth, mType, vsys, o.Name))
	return readPanoramaLayer2Subinterface(d, meta)
}

func readPanoramaLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, mType, vsys, name := parsePanoramaLayer2SubinterfaceId(d.Id())

	o, err := pano.Network.Layer2Subinterface.Get(tmpl, ts, iType, eth, mType, name)
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
	d.Set("parent_mode", mType)
	saveLayer2Subinterface(d, o)

	return nil
}

func updatePanoramaLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, mType, vsys, o := parsePanoramaLayer2Subinterface(d)

	lo, err := pano.Network.Layer2Subinterface.Get(tmpl, ts, iType, eth, mType, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.Layer2Subinterface.Edit(tmpl, ts, iType, eth, mType, vsys, lo); err != nil {
		return err
	}

	d.SetId(buildPanoramaLayer2SubinterfaceId(tmpl, ts, iType, eth, mType, vsys, o.Name))
	return readPanoramaLayer2Subinterface(d, meta)
}

func deletePanoramaLayer2Subinterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, iType, eth, mType, _, name := parsePanoramaLayer2SubinterfaceId(d.Id())

	err := pano.Network.Layer2Subinterface.Delete(tmpl, ts, iType, eth, mType, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
