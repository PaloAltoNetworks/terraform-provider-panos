package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZoneEntry() *schema.Resource {
	return &schema.Resource{
		Create: createZoneEntry,
		Read:   readZoneEntry,
		Delete: deleteZoneEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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

func parseZoneEntryId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildZoneEntryId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createZoneEntry(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)
	zone_name := d.Get("zone").(string)
	mode := d.Get("mode").(string)
	iface := d.Get("interface").(string)

	if err := fw.Network.Zone.SetInterface(vsys, zone_name, mode, iface); err != nil {
		return err
	}

	d.SetId(buildZoneEntryId(vsys, zone_name, mode, iface))
	return readZoneEntry(d, meta)
}

func readZoneEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, zone_name, mode, iface := parseZoneEntryId(d.Id())

	// Three possibilities:  either the zone itself doesn't exist, the
	// interface isn't present, or the mode is incorrect.
	o, err := fw.Network.Zone.Get(vsys, zone_name)
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

func deleteZoneEntry(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, zone_name, mode, iface := parseZoneEntryId(d.Id())

	if err := fw.Network.Zone.DeleteInterface(vsys, zone_name, mode, iface); err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
