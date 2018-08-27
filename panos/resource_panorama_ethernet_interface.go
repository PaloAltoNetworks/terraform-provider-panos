package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaEthernetInterface() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaEthernetInterface,
		Read:   readPanoramaEthernetInterface,
		Update: updatePanoramaEthernetInterface,
		Delete: deletePanoramaEthernetInterface,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template_stack"},
			},
			"template_stack": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template"},
			},
			"vsys": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
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

func parsePanoramaEthernetInterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaEthernetInterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parsePanoramaEthernetInterface(d *schema.ResourceData) (string, string, string, eth.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
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

	return tmpl, ts, vsys, o
}

func createPanoramaEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaEthernetInterface(d)

	if err := pano.Network.EthernetInterface.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaEthernetInterfaceId(tmpl, ts, vsys, o.Name))
	return readPanoramaEthernetInterface(d, meta)
}

func readPanoramaEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaEthernetInterfaceId(d.Id())

	o, err := pano.Network.EthernetInterface.Get(tmpl, ts, name)
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

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("name", o.Name)
	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
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

func updatePanoramaEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaEthernetInterface(d)

	lo, err := pano.Network.EthernetInterface.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.EthernetInterface.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaEthernetInterface(d, meta)
}

func deletePanoramaEthernetInterface(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, _, name := parsePanoramaEthernetInterfaceId(d.Id())

	err := pano.Network.EthernetInterface.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
