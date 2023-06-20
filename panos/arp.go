package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/arp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceArps() *schema.Resource {
	s := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(If Panorama) The template where the interface is",
		},
		"interface_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The interface type",
			Default:     arp.TypeEthernet,
			ValidateFunc: validateStringIn(
				arp.TypeEthernet,
				arp.TypeAggregate,
				arp.TypeVlan,
			),
		},
		"interface_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The interface name; leave this empty for VLAN interfaces",
		},
		"subinterface_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The subinterface name",
		},
	}

	for key, val := range listingSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceArpsRead,

		Schema: s,
	}
}

func dataSourceArpsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	iType := d.Get("interface_type").(string)
	iName := d.Get("interface_name").(string)
	subName := d.Get("subinterface_name").(string)
	tmpl := d.Get("template").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = base64Encode([]string{
			iType, iName, subName,
		})
		listing, err = con.Network.Arp.GetList(iType, iName, subName)
	case *pango.Panorama:
		id = base64Encode([]string{
			tmpl, "", iType, iName, subName,
		})
		listing, err = con.Network.Arp.GetList(tmpl, "", iType, iName, subName)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceArp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArpRead,

		Schema: arpSchema(false),
	}
}

func dataSourceArpRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o arp.Entry
	iType := d.Get("interface_type").(string)
	iName := d.Get("interface_name").(string)
	subName := d.Get("subinterface_name").(string)
	tmpl := d.Get("template").(string)
	ip := d.Get("ip").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallArpId(iType, iName, subName, ip)
		o, err = con.Network.Arp.Get(iType, iName, subName, ip)
	case *pango.Panorama:
		id = buildPanoramaArpId(tmpl, "", iType, iName, subName, ip)
		o, err = con.Network.Arp.Get(tmpl, "", iType, iName, subName, ip)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveArp(d, o)

	return nil
}

// Resource.
func resourceArp() *schema.Resource {
	return &schema.Resource{
		Create: createArp,
		Read:   readArp,
		Update: updateArp,
		Delete: deleteArp,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: arpSchema(true),
	}
}

func createArp(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	iType := d.Get("interface_type").(string)
	iName := d.Get("interface_name").(string)
	subName := d.Get("subinterface_name").(string)
	tmpl := d.Get("template").(string)
	o := loadArp(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallArpId(iType, iName, subName, o.Ip)
		err = con.Network.Arp.Set(iType, iName, subName, o)
	case *pango.Panorama:
		id = buildPanoramaArpId(tmpl, "", iType, iName, subName, o.Ip)
		err = con.Network.Arp.Set(tmpl, "", iType, iName, subName, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readArp(d, meta)
}

func readArp(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o arp.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		iType, iName, subName, ip := parseFirewallArpId(d.Id())
		o, err = con.Network.Arp.Get(iType, iName, subName, ip)
	case *pango.Panorama:
		tmpl, ts, iType, iName, subName, ip := parsePanoramaArpId(d.Id())
		o, err = con.Network.Arp.Get(tmpl, ts, iType, iName, subName, ip)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveArp(d, o)
	return nil
}

func updateArp(d *schema.ResourceData, meta interface{}) error {
	iType := d.Get("interface_type").(string)
	iName := d.Get("interface_name").(string)
	subName := d.Get("subinterface_name").(string)
	tmpl := d.Get("template").(string)
	o := loadArp(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Network.Arp.Get(iType, iName, subName, o.Ip)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.Arp.Edit(iType, iName, subName, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Network.Arp.Get(tmpl, "", iType, iName, subName, o.Ip)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.Arp.Edit(tmpl, "", iType, iName, subName, lo); err != nil {
			return err
		}
	}

	return readArp(d, meta)
}

func deleteArp(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		iType, iName, subName, ip := parseFirewallArpId(d.Id())
		err = con.Network.Arp.Delete(iType, iName, subName, ip)
	case *pango.Panorama:
		tmpl, ts, iType, iName, subName, ip := parsePanoramaArpId(d.Id())
		err = con.Network.Arp.Delete(tmpl, ts, iType, iName, subName, ip)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func arpSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(If Panorama) The template where the interface is",
			ForceNew:    true,
		},
		"interface_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The interface type",
			Default:     arp.TypeEthernet,
			ForceNew:    true,
			ValidateFunc: validateStringIn(
				arp.TypeEthernet,
				arp.TypeAggregate,
				arp.TypeVlan,
			),
		},
		"interface_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The interface name; leave this empty for VLAN interfaces",
			ForceNew:    true,
		},
		"subinterface_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The subinterface name",
			ForceNew:    true,
		},
		"ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The IP address",
			ForceNew:    true,
		},
		"mac_address": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The MAC address",
		},
		"interface": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(For interface_type = vlan) The interface",
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "interface_type", "interface_name", "subinterface_name", "ip"})
	}

	return ans
}

func loadArp(d *schema.ResourceData) arp.Entry {
	return arp.Entry{
		Ip:         d.Get("ip").(string),
		MacAddress: d.Get("mac_address").(string),
		Interface:  d.Get("interface").(string),
	}
}

func saveArp(d *schema.ResourceData, o arp.Entry) {
	d.Set("ip", o.Ip)
	d.Set("mac_address", o.MacAddress)
	d.Set("interface", o.Interface)
}

// Id functions.
func parseFirewallArpId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildFirewallArpId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaArpId(v string) (string, string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4], t[5]
}

func buildPanoramaArpId(a, b, c, d, e, f string) string {
	return strings.Join([]string{a, b, c, d, e, f}, IdSeparator)
}
