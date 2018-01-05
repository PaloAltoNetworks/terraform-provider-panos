package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/eth"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceEthernetInterface() *schema.Resource {
	return &schema.Resource{
		Create: createEthernetInterface,
		Read:   readEthernetInterface,
		Update: updateEthernetInterface,
		Delete: deleteEthernetInterface,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ethernet interface's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The vsys to import this ethernet interface into",
			},
			"mode": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The interface mode (layer3, layer2, virtual-wire, tap, ha, decrypt-mirror, aggregate-group)",
				ValidateFunc: validateStringIn("layer3", "layer2", "virtual-wire", "tap", "ha", "decrypt-mirror", "aggregate-group"),
			},
			"static_ips": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of static IP addresses",
			},
			"enable_dhcp": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_dhcp_default_route": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dhcp_default_route_metric": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"management_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"adjust_tcp_mss": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"netflow_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"lldp_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"lldp_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"link_speed": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("10", "100", "1000", "auto"),
			},
			"link_duplex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("full", "half", "auto"),
			},
			"link_state": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("up", "down", "auto"),
			},
			"aggregate_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_mss_adjust": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ipv6_mss_adjust": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseEthernetInterfaceId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildEthernetInterfaceId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func parseEthernetInterface(d *schema.ResourceData) (string, eth.Entry) {
	vsys := d.Get("vsys").(string)
	o := eth.Entry{
		Name:                   d.Get("name").(string),
		Mode:                   d.Get("mode").(string),
		StaticIps:              asStringList(d.Get("static_ips").([]interface{})),
		EnableDhcp:             d.Get("enable_dhcp").(bool),
		CreateDhcpDefaultRoute: d.Get("create_dhcp_default_route").(bool),
		DhcpDefaultRouteMetric: d.Get("dhcp_default_route_metric").(int),
		Ipv6Enabled:            d.Get("ipv6_enabled").(bool),
		ManagementProfile:      d.Get("management_profile").(string),
		Mtu:                    d.Get("mtu").(int),
		AdjustTcpMss:           d.Get("adjust_tcp_mss").(bool),
		NetflowProfile:         d.Get("netflow_profile").(string),
		LldpEnabled:            d.Get("lldp_enabled").(bool),
		LldpProfile:            d.Get("lldp_profile").(string),
		LinkSpeed:              d.Get("link_speed").(string),
		LinkDuplex:             d.Get("link_duplex").(string),
		LinkState:              d.Get("link_state").(string),
		AggregateGroup:         d.Get("aggregate_group").(string),
		Comment:                d.Get("comment").(string),
		Ipv4MssAdjust:          d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:          d.Get("ipv6_mss_adjust").(int),
	}

	return vsys, o
}

func createEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseEthernetInterface(d)

	if err := fw.Network.EthernetInterface.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildEthernetInterfaceId(vsys, o.Name))
	return readEthernetInterface(d, meta)
}

func readEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseEthernetInterfaceId(d.Id())

	o, err := fw.Network.EthernetInterface.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("mode", o.Mode)
	if err = d.Set("static_ips", o.StaticIps); err != nil {
		log.Printf("[WARN] Error setting 'static_ips' for %q: %s", d.Id(), err)
	}
	d.Set("enable_dhcp", o.EnableDhcp)
	d.Set("create_dhcp_default_route", o.CreateDhcpDefaultRoute)
	d.Set("dhcp_default_route_metric", o.DhcpDefaultRouteMetric)
	d.Set("ipv6_enabled", o.Ipv6Enabled)
	d.Set("management_profile", o.ManagementProfile)
	d.Set("mtu", o.Mtu)
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("netflow_profile", o.NetflowProfile)
	d.Set("lldp_enabled", o.LldpEnabled)
	d.Set("lldp_profile", o.LldpProfile)
	d.Set("link_speed", o.LinkSpeed)
	d.Set("link_duplex", o.LinkDuplex)
	d.Set("link_state", o.LinkState)
	d.Set("aggregate_group", o.AggregateGroup)
	d.Set("comment", o.Comment)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)

	return nil
}

func updateEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseEthernetInterface(d)

	lo, err := fw.Network.EthernetInterface.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.EthernetInterface.Edit(vsys, lo); err != nil {
		return err
	}

	return readEthernetInterface(d, meta)
}

func deleteEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	err := fw.Network.EthernetInterface.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
