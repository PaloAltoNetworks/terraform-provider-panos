package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLayer3Subinterface() *schema.Resource {
	return &schema.Resource{
		Create: createLayer3Subinterface,
		Read:   readLayer3Subinterface,
		Update: updateLayer3Subinterface,
		Delete: deleteLayer3Subinterface,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: layer3SubinterfaceSchema(false),
	}
}

func layer3SubinterfaceSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"parent_interface": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"interface_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      layer3.EthernetInterface,
			ValidateFunc: validateStringIn(layer3.EthernetInterface, layer3.AggregateInterface),
			ForceNew:     true,
		},
		"vsys": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "vsys1",
		},
		"tag": {
			Type:     schema.TypeInt,
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
		"netflow_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"comment": {
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

func parseLayer3SubinterfaceId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildLayer3SubinterfaceId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func loadLayer3Subinterface(d *schema.ResourceData) layer3.Entry {
	return layer3.Entry{
		Name:                   d.Get("name").(string),
		Tag:                    d.Get("tag").(int),
		StaticIps:              asStringList(d.Get("static_ips").([]interface{})),
		Ipv6Enabled:            d.Get("ipv6_enabled").(bool),
		Ipv6InterfaceId:        d.Get("ipv6_interface_id").(string),
		ManagementProfile:      d.Get("management_profile").(string),
		Mtu:                    d.Get("mtu").(int),
		AdjustTcpMss:           d.Get("adjust_tcp_mss").(bool),
		Ipv4MssAdjust:          d.Get("ipv4_mss_adjust").(int),
		Ipv6MssAdjust:          d.Get("ipv6_mss_adjust").(int),
		NetflowProfile:         d.Get("netflow_profile").(string),
		Comment:                d.Get("comment").(string),
		EnableDhcp:             d.Get("enable_dhcp").(bool),
		CreateDhcpDefaultRoute: d.Get("create_dhcp_default_route").(bool),
		DhcpDefaultRouteMetric: d.Get("dhcp_default_route_metric").(int),
		DhcpSendHostnameEnable: d.Get("dhcp_send_hostname_enable").(bool),
		DhcpSendHostnameValue:  d.Get("dhcp_send_hostname_value").(string),
		DecryptForward:         d.Get("decrypt_forward").(bool),
	}
}

func saveLayer3Subinterface(d *schema.ResourceData, o layer3.Entry) {
	d.Set("name", o.Name)
	d.Set("tag", o.Tag)
	if err := d.Set("static_ips", o.StaticIps); err != nil {
		log.Printf("[WARN] Error setting 'static_ips' for %q: %s", d.Id(), err)
	}
	d.Set("ipv6_enabled", o.Ipv6Enabled)
	d.Set("ipv6_interface_id", o.Ipv6InterfaceId)
	d.Set("management_profile", o.ManagementProfile)
	d.Set("mtu", o.Mtu)
	d.Set("adjust_tcp_mss", o.AdjustTcpMss)
	d.Set("ipv4_mss_adjust", o.Ipv4MssAdjust)
	d.Set("ipv6_mss_adjust", o.Ipv6MssAdjust)
	d.Set("netflow_profile", o.NetflowProfile)
	d.Set("comment", o.Comment)
	d.Set("enable_dhcp", o.EnableDhcp)
	d.Set("create_dhcp_default_route", o.CreateDhcpDefaultRoute)
	d.Set("dhcp_default_route_metric", o.DhcpDefaultRouteMetric)
	d.Set("dhcp_send_hostname_enable", o.DhcpSendHostnameEnable)
	d.Set("dhcp_send_hostname_value", o.DhcpSendHostnameValue)
	d.Set("decrypt_forward", o.DecryptForward)
}

func parseLayer3Subinterface(d *schema.ResourceData) (string, string, string, layer3.Entry) {
	eth := d.Get("parent_interface").(string)
	iType := d.Get("interface_type").(string)
	vsys := d.Get("vsys").(string)
	o := loadLayer3Subinterface(d)

	return iType, eth, vsys, o
}

func createLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	iType, eth, vsys, o := parseLayer3Subinterface(d)

	if err := fw.Network.Layer3Subinterface.Set(iType, eth, vsys, o); err != nil {
		return err
	}

	d.SetId(buildLayer3SubinterfaceId(iType, eth, vsys, o.Name))
	return readLayer3Subinterface(d, meta)
}

func readLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	iType, eth, vsys, name := parseLayer3SubinterfaceId(d.Id())

	o, err := fw.Network.Layer3Subinterface.Get(iType, eth, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}
	rv, err := fw.IsImported(util.InterfaceImport, "", "", vsys, name)
	if err != nil {
		return err
	}

	if rv {
		d.Set("vsys", vsys)
	} else {
		d.Set("vsys", fmt.Sprintf("(not %s)", vsys))
	}
	d.Set("interface_type", iType)
	d.Set("parent_interface", eth)
	saveLayer3Subinterface(d, o)

	return nil
}

func updateLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	iType, eth, vsys, o := parseLayer3Subinterface(d)

	lo, err := fw.Network.Layer3Subinterface.Get(iType, eth, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.Layer3Subinterface.Edit(iType, eth, vsys, lo); err != nil {
		return err
	}

	d.SetId(buildLayer3SubinterfaceId(iType, eth, vsys, o.Name))
	return readLayer3Subinterface(d, meta)
}

func deleteLayer3Subinterface(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	iType, eth, _, name := parseLayer3SubinterfaceId(d.Id())

	err := fw.Network.Layer3Subinterface.Delete(iType, eth, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
