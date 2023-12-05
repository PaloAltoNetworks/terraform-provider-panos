package panos

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVlanEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateVlanEntry,
		Read:   readVlanEntry,
		Update: createUpdateVlanEntry,
		Delete: deleteVlanEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: vlanEntrySchema(false),
	}
}

func vlanEntrySchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vlan": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"interface": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"mac_addresses": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"live_mac_addresses": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if p {
		ans["template"] = templateSchema(false)
	}

	return ans
}

func parseVlanEntry(d *schema.ResourceData) (string, string, []string, []string) {
	vlan, iface, macs, liveMacs := loadVlanEntry(d)

	wmm := make(map[string]bool)
	for _, v := range macs {
		wmm[v] = true
	}

	lmm := make(map[string]bool)
	for _, v := range liveMacs {
		lmm[v] = true
	}

	rmMacs := make([]string, 0, len(liveMacs))
	for _, mac := range liveMacs {
		if wanted := wmm[mac]; !wanted {
			rmMacs = append(rmMacs, mac)
		}
	}

	addMacs := make([]string, 0, len(macs))
	for _, mac := range macs {
		if present := lmm[mac]; !present {
			addMacs = append(addMacs, mac)
		}
	}

	return vlan, iface, rmMacs, addMacs
}

func loadVlanEntry(d *schema.ResourceData) (string, string, []string, []string) {
	vlan := d.Get("vlan").(string)
	iface := d.Get("interface").(string)
	macs := setAsList(d.Get("mac_addresses").(*schema.Set))
	liveMacs := setAsList(d.Get("live_mac_addresses").(*schema.Set))

	return vlan, iface, macs, liveMacs
}

func saveVlanEntry(d *schema.ResourceData, vlan, iface string, macs []string) {
	d.Set("vlan", vlan)
	d.Set("interface", iface)
	if err := d.Set("mac_addresses", listAsSet(macs)); err != nil {
		log.Printf("[WARN] Error setting 'mac_addresses' param for %q: %s", d.Id(), err)
	}
	if err := d.Set("live_mac_addresses", listAsSet(macs)); err != nil {
		log.Printf("[WARN] Error setting 'live_mac_addresses' param for %q: %s", d.Id(), err)
	}
}

func parseVlanEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildVlanEntryId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateVlanEntry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_vlan_entry")
	if err != nil {
		return err
	}
	vlan, iface, rmMacs, addMacs := parseVlanEntry(d)

	if err = fw.Network.Vlan.SetInterface(vlan, iface, rmMacs, addMacs); err != nil {
		return err
	}

	d.SetId(buildVlanEntryId(vlan, iface))
	return readVlanEntry(d, meta)
}

func readVlanEntry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_vlan_entry")
	if err != nil {
		return err
	}

	vlan, iface := parseVlanEntryId(d.Id())

	// Two possibilities:  either the router itself doesn't exist or the
	// interface isn't present.
	o, err := fw.Network.Vlan.Get(vlan)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	found := false
	for _, x := range o.Interfaces {
		if x == iface {
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
	return nil
}

func deleteVlanEntry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_vlan_entry")
	if err != nil {
		return err
	}
	vlan, iface := parseVlanEntryId(d.Id())

	if err = fw.Network.Vlan.DeleteInterface(vlan, iface); err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}
