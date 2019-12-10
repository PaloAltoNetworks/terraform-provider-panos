package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	agg "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaAggregateInterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaAggregateInterface,
		Read:   readPanoramaAggregateInterface,
		Update: updatePanoramaAggregateInterface,
		Delete: deletePanoramaAggregateInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: aggregateInterfaceSchema(true),
	}
}

func parsePanoramaAggregateInterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaAggregateInterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaAggregateInterface(d *schema.ResourceData) (string, string, string, agg.Entry) {
	tmpl := d.Get("template").(string)
	ts := ""
	vsys := d.Get("vsys").(string)

	o := loadAggregateInterface(d)

	return tmpl, ts, vsys, o
}

func createPanoramaAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaAggregateInterface(d)

	if err := pano.Network.AggregateInterface.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaAggregateInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaAggregateInterface(d, meta)
}

func readPanoramaAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaAggregateInterfaceId(d.Id())

	o, err := pano.Network.AggregateInterface.Get(tmpl, ts, name)
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
	saveAggregateInterface(d, o)

	return nil
}

func updatePanoramaAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaAggregateInterface(d)

	lo, err := pano.Network.AggregateInterface.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.AggregateInterface.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	d.SetId(buildPanoramaAggregateInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaAggregateInterface(d, meta)
}

func deletePanoramaAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaAggregateInterfaceId(d.Id())

	err := pano.Network.AggregateInterface.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
