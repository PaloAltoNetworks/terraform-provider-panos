package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	agg "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		"vsys": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "vsys1",
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
		"lacp_enable": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"lacp_fast_failover": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"lacp_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", agg.LacpModePassive, agg.LacpModeActive),
		},
		"lacp_transmission_rate": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", agg.LacpTransmissionRateFast, agg.LacpTransmissionRateSlow),
		},
		"lacp_system_priority": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"lacp_max_ports": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"lacp_ha_passive_pre_negotiation": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"lacp_ha_enable_same_system_mac": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"lacp_ha_same_system_mac_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lldp_enable": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"lldp_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lldp_ha_passive_pre_negotiation": {
			Type:     schema.TypeBool,
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
		Name:                        d.Get("name").(string),
		Mode:                        d.Get("mode").(string),
		NetflowProfile:              d.Get("netflow_profile").(string),
		Mtu:                         d.Get("mtu").(int),
		AdjustTcpMss:                d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:               d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:               d.Get("ipv6_mss_adjust").(int),
		EnableUntaggedSubinterface:  d.Get("enable_untagged_subinterface").(bool),
		StaticIps:                   asStringList(d.Get("static_ips").([]interface{})),
		Ipv6Enabled:                 d.Get("ipv6_enabled").(bool),
		Ipv6InterfaceId:             d.Get("ipv6_interface_id").(string),
		ManagementProfile:           d.Get("management_profile").(string),
		EnableDhcp:                  d.Get("enable_dhcp").(bool),
		CreateDhcpDefaultRoute:      d.Get("create_dhcp_default_route").(bool),
		DhcpDefaultRouteMetric:      d.Get("dhcp_default_route_metric").(int),
		LacpEnable:                  d.Get("lacp_enable").(bool),
		LacpFastFailover:            d.Get("lacp_fast_failover").(bool),
		LacpMode:                    d.Get("lacp_mode").(string),
		LacpTransmissionRate:        d.Get("lacp_transmission_rate").(string),
		LacpSystemPriority:          d.Get("lacp_system_priority").(int),
		LacpMaxPorts:                d.Get("lacp_max_ports").(int),
		LacpHaPassivePreNegotiation: d.Get("lacp_ha_passive_pre_negotiation").(bool),
		LacpHaEnableSameSystemMac:   d.Get("lacp_ha_enable_same_system_mac").(bool),
		LacpHaSameSystemMacAddress:  d.Get("lacp_ha_same_system_mac_address").(string),
		LldpEnable:                  d.Get("lldp_enable").(bool),
		LldpProfile:                 d.Get("lldp_profile").(string),
		LldpHaPassivePreNegotiation: d.Get("lldp_ha_passive_pre_negotiation").(bool),
		Comment:                     d.Get("comment").(string),
		DecryptForward:              d.Get("decrypt_forward").(bool),
		DhcpSendHostnameEnable:      d.Get("dhcp_send_hostname_enable").(bool),
		DhcpSendHostnameValue:       d.Get("dhcp_send_hostname_value").(string),
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
	d.Set("lacp_enable", o.LacpEnable)
	d.Set("lacp_fast_failover", o.LacpFastFailover)
	d.Set("lacp_mode", o.LacpMode)
	d.Set("lacp_transmission_rate", o.LacpTransmissionRate)
	d.Set("lacp_system_priority", o.LacpSystemPriority)
	d.Set("lacp_max_ports", o.LacpMaxPorts)
	d.Set("lacp_ha_passive_pre_negotiation", o.LacpHaPassivePreNegotiation)
	d.Set("lacp_ha_enable_same_system_mac", o.LacpHaEnableSameSystemMac)
	d.Set("lacp_ha_same_system_mac_address", o.LacpHaSameSystemMacAddress)
	d.Set("lldp_enable", o.LldpEnable)
	d.Set("lldp_profile", o.LldpProfile)
	d.Set("lldp_ha_passive_pre_negotiation", o.LldpHaPassivePreNegotiation)
	d.Set("decrypt_forward", o.DecryptForward)
	d.Set("dhcp_send_hostname_enable", o.DhcpSendHostnameEnable)
	d.Set("dhcp_send_hostname_value", o.DhcpSendHostnameValue)
}

func parseAggregateInterface(d *schema.ResourceData) (string, agg.Entry) {
	o := loadAggregateInterface(d)
	vsys := d.Get("vsys").(string)

	return vsys, o
}

func buildAggregateInterfaceId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseAggregateInterfaceId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func createAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseAggregateInterface(d)

	if err := fw.Network.AggregateInterface.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildAggregateInterfaceId(vsys, o.Name))
	return readAggregateInterface(d, meta)
}

func readAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseAggregateInterfaceId(d.Id())

	o, err := fw.Network.AggregateInterface.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := fw.IsImported(util.InterfaceImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	saveAggregateInterface(d, o)
	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}

	return nil
}

func updateAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseAggregateInterface(d)

	lo, err := fw.Network.AggregateInterface.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.AggregateInterface.Edit(vsys, lo); err != nil {
		return err
	}

	return readAggregateInterface(d, meta)
}

func deleteAggregateInterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	_, name := parseAggregateInterfaceId(d.Id())

	err := fw.Network.AggregateInterface.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
