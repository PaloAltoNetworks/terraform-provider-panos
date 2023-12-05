package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/vlan"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVlan() *schema.Resource {
	return &schema.Resource{
		Create: createVlan,
		Read:   readVlan,
		Update: updateVlan,
		Delete: deleteVlan,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: vlanSchema(false),
	}
}

func vlanSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys": vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"vlan_interface": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"interfaces": {
			Type:     schema.TypeSet,
			Computed: true,
			Optional: true,
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

func parseVlan(d *schema.ResourceData) (string, vlan.Entry) {
	vsys := d.Get("vsys").(string)
	o := loadVlan(d)

	return vsys, o
}

func loadVlan(d *schema.ResourceData) vlan.Entry {
	return vlan.Entry{
		Name:          d.Get("name").(string),
		VlanInterface: d.Get("vlan_interface").(string),
		Interfaces:    setAsList(d.Get("interfaces").(*schema.Set)),
	}
}

func saveVlan(d *schema.ResourceData, o vlan.Entry) {
	d.Set("name", o.Name)
	d.Set("vlan_interface", o.VlanInterface)
	if err := d.Set("interfaces", listAsSet(o.Interfaces)); err != nil {
		log.Printf("[WARN] Error setting 'interfaces' param for %q: %s", d.Id(), err)
	}
}

func parseVlanId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildVlanId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createVlan(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseVlan(d)

	if err := fw.Network.Vlan.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildVlanId(vsys, o.Name))
	return readVlan(d, meta)
}

func readVlan(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseVlanId(d.Id())

	o, err := fw.Network.Vlan.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	rv, err := fw.IsImported(util.VlanImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	saveVlan(d, o)

	return nil
}

func updateVlan(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseVlan(d)

	lo, err := fw.Network.Vlan.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o, false)
	if err = fw.Network.Vlan.Edit(vsys, lo); err != nil {
		return err
	}

	return readVlan(d, meta)
}

func deleteVlan(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	_, name := parseVlanId(d.Id())

	err := fw.Network.Vlan.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
