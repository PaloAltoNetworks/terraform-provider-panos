package panos

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/dhcp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDhcpRelay() *schema.Resource {
	return &schema.Resource{
		Create: createDhcpRelay,
		Read:   readDhcpRelay,
		Update: updateDhcpRelay,
		Delete: deleteDhcpRelay,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Interface to enable dhcp relay on",
			},
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "Vsys to put interface in",
			},
			"ipv4_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables dhcp relay for ipv4",
			},
			"ipv6_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables dhcp relay for ipv6",
			},
			"ipv4_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of ipv4 dhcp servers",
			},
			"ipv6_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interface": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
					Description: "List of ipv6 dhcp servers",
				},
			},
		},
	}
}

func parseDhcpRelay(d *schema.ResourceData) (string, dhcp.Entry) {
	vsys := d.Get("vsys").(string)

	if !castIntoBool(d.Get("ipv4_enabled")) && !castIntoBool(d.Get("ipv6_enabled")) {
		return "", dhcp.Entry{}
	}

	o := dhcp.Entry{
		Name: d.Get("name").(string),

		Relay: &dhcp.Relay{
			Ipv4Enabled: castIntoBool(d.Get("ipv4_enabled")),
			Ipv4Servers: flattenIPv4DhcpServers(d.Get("ipv4_servers")),
			Ipv6Enabled: castIntoBool(d.Get("ipv6_enabled")),
			Ipv6Servers: flattenIPv6DhcpServers(d.Get("ipv6_servers")),
		},
	}

	return vsys, o
}

func castIntoBool(v interface{}) bool {
	return v.(bool)
}

func flattenIPv4DhcpServers(o interface{}) []string {
	var d []string
	for _, v := range o.([]interface{}) {
		d = append(d, v.(string))
	}
	return d
}

func flattenIPv6DhcpServers(d interface{}) []dhcp.Ipv6Server {
	var o []dhcp.Ipv6Server
	for _, v := range d.([]interface{}) {
		m := v.(map[string]interface{})
		o = append(o, dhcp.Ipv6Server{
			Server:    m["server"].(string),
			Interface: m["interface"].(string),
		})
	}
	return o
}

func isEmpty(o dhcp.Entry) bool {
	return !o.Relay.Ipv4Enabled && !o.Relay.Ipv6Enabled && len(o.Relay.Ipv4Servers) ==0  && len(o.Relay.Ipv6Servers) == 0
}

func readDhcpRelay(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.Dhcp.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("ipv4_enabled", o.Relay.Ipv4Enabled)
	d.Set("ipv6_enabled", o.Relay.Ipv6Enabled)
	d.Set("ipv4_servers", o.Relay.Ipv6Enabled)
	d.Set("ipv6_servers", o.Relay.Ipv6Servers)

	return nil
}

func createDhcpRelay(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseDhcpRelay(d)

	if vsys == "" && isEmpty(o) {
		return fmt.Errorf("[ERROR] Error in creating DHCP relay: no valid configuration provided")
	}

	if err := fw.Network.Dhcp.Set(o); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())
	return d.Set("vsys", vsys)
}

func updateDhcpRelay(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	dg, err := fw.Network.Dhcp.Get(name)
	if err != nil {
		return err
	}

	vsys, o := parseDhcpRelay(d)

	if vsys == "" && isEmpty(o) {
		return fmt.Errorf("[ERROR] Error in updating DHCP relay: no valid configuration provided")
	}

	dg.Copy(o)
	if err = fw.Network.Dhcp.Edit(o); err != nil {
		return err
	}

	return readDhcpRelay(d, meta)
}

func deleteDhcpRelay(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.Dhcp.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
