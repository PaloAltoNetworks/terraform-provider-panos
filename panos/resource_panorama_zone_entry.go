package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaZoneEntry() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaZoneEntry,
		Read:   readPanoramaZoneEntry,
		Delete: deletePanoramaZoneEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      zone.ModeL3,
				ValidateFunc: validateStringIn(zone.ModeL3, zone.ModeL2, zone.ModeVirtualWire, zone.ModeTap, zone.ModeExternal),
				ForceNew:     true,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func parsePanoramaZoneEntryId(v string) (string, string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4], t[5]
}

func buildPanoramaZoneEntryId(a, b, c, d, e, f string) string {
	return strings.Join([]string{a, b, c, d, e, f}, IdSeparator)
}

func createPanoramaZoneEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl := d.Get("template").(string)
	ts := ""
	vsys := d.Get("vsys").(string)
	zone_name := d.Get("zone").(string)
	mode := d.Get("mode").(string)
	iface := d.Get("interface").(string)

	if err := pano.Network.Zone.SetInterface(tmpl, ts, vsys, zone_name, mode, iface); err != nil {
		return err
	}

	d.SetId(buildPanoramaZoneEntryId(tmpl, ts, vsys, zone_name, mode, iface))
	return readPanoramaZoneEntry(d, meta)
}

func readPanoramaZoneEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, zone_name, mode, iface := parsePanoramaZoneEntryId(d.Id())

	// Three possibilities:  either the zone itself doesn't exist, the
	// interface isn't present, or the mode is incorrect.
	o, err := pano.Network.Zone.Get(tmpl, ts, vsys, zone_name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	if o.Mode != mode {
		d.SetId("")
		return nil
	}

	for i := range o.Interfaces {
		if o.Interfaces[i] == iface {
			d.Set("template", tmpl)
			d.Set("vsys", vsys)
			d.Set("zone", zone_name)
			d.Set("mode", mode)
			d.Set("interface", iface)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deletePanoramaZoneEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, zone_name, mode, iface := parsePanoramaZoneEntryId(d.Id())

	if err := pano.Network.Zone.DeleteInterface(tmpl, ts, vsys, zone_name, mode, iface); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
