package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	agg "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAggregateInterface() *schema.Resource {
	return &schema.Resource{
		Create: createAggregateInterface,
		Read:   readAggregateInterface,
		Update: updateAggregateInterface,
		Delete: deleteAggregateInterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: aggregateInterfaceSchema(false),
	}
}

func aggregateInterfaceSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      agg.ModeLayer3,
			ValidateFunc: validateStringIn(agg.ModeHa, agg.ModeDecryptMirror, agg.ModeVirtualWire, agg.ModeLayer2, agg.ModeLayer3),
		},
		"netflow_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mtu": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"adjust_tcp_mss": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ipv4_mss_adjust": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ipv6_mss_adjust": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"enable_untagged_subinterface": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"static_ips": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of static IP addresses",
		},
		"ipv6_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ipv6_interface_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"management_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enable_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"create_dhcp_default_route": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dhcp_default_route_metric": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"decrypt_forward": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dhcp_send_hostname_enable": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"dhcp_send_hostname_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(false)
	}

	return ans
}

func loadAggregateInterface(d *schema.ResourceData) agg.Entry {
	return agg.Entry{
		Name:                       d.Get("name").(string),
		Mode:                       d.Get("mode").(string),
		NetflowProfile:             d.Get("netflow_profile").(string),
		Mtu:                        d.Get("mtu").(int),
		AdjustTcpMss:               d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:              d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:              d.Get("ipv6_mss_adjust").(int),
		EnableUntaggedSubinterface: d.Get("enable_untagged_subinterface").(bool),
		StaticIps:                  asStringList(d.Get("static_ips").([]interface{})),
		Ipv6Enabled:                d.Get("ipv6_enabled").(bool),
		Ipv6InterfaceId:            d.Get("ipv6_interface_id").(string),
		ManagementProfile:          d.Get("management_profile").(string),
		EnableDhcp:                 d.Get("enable_dhcp").(bool),
		CreateDhcpDefaultRoute:     d.Get("create_dhcp_default_route").(bool),
		DhcpDefaultRouteMetric:     d.Get("dhcp_default_route_metric").(int),
		Comment:                    d.Get("comment").(string),
		DecryptForward:             d.Get("decrypt_forward").(bool),
		DhcpSendHostnameEnable:     d.Get("dhcp_send_hostname_enable").(bool),
		DhcpSendHostnameValue:      d.Get("dhcp_send_hostname_value").(string),
	}
}

func saveAggregateInterface(d *schema.ResourceData, o agg.Entry) {
	d.Set("name", o.Name)
	d.Set("mode", o.Mode)
	d.Set("netflow_profile", o.NetflowProfile)
	d.Set("mtu", o.Mtu)
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv4MssAdjust)
	d.Set("enable_untagged_subinterface", o.EnableUntaggedSubinterface)
	if err := d.Set("static_ips", o.StaticIps); err != nil {
		log.Printf("[WARN] Error setting 'static_ips' for %q: %s", d.Id(), err)
	}
	d.Set("ipv6_enabled", o.Ipv6Enabled)
	d.Set("ipv6_interface_id", o.Ipv6InterfaceId)
	d.Set("management_profile", o.ManagementProfile)
	d.Set("enable_dhcp", o.EnableDhcp)
	d.Set("create_dhcp_default_route", o.CreateDhcpDefaultRoute)
	d.Set("dhcp_default_route_metric", o.DhcpDefaultRouteMetric)
	d.Set("comment", o.Comment)
	d.Set("decrypt_forward", o.DecryptForward)
	d.Set("dhcp_send_hostname_enable", o.DhcpSendHostnameEnable)
	d.Set("dhcp_send_hostname_value", o.DhcpSendHostnameValue)
}

func parseAggregateInterface(d *schema.ResourceData) agg.Entry {
	o := loadAggregateInterface(d)

	return o
}

func createAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseAggregateInterface(d)

	if err := fw.Network.AggregateInterface.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readAggregateInterface(d, meta)
}

func readAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Id()

	o, err := fw.Network.AggregateInterface.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAggregateInterface(d, o)

	return nil
}

func updateAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseAggregateInterface(d)

	lo, err := fw.Network.AggregateInterface.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.AggregateInterface.Edit(lo); err != nil {
		return err
	}

	return readAggregateInterface(d, meta)
}

func deleteAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Id()

	err := fw.Network.AggregateInterface.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
