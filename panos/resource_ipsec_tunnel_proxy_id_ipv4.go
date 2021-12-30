package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ipsectunnel/proxyid/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIpsecTunnelProxyIdIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createIpsecTunnelProxyIdIpv4,
		Read:   readIpsecTunnelProxyIdIpv4,
		Update: updateIpsecTunnelProxyIdIpv4,
		Delete: deleteIpsecTunnelProxyIdIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipsec_tunnel": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_any": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protocol_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_tcp_local": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_tcp_remote": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_udp_local": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_udp_remote": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseIpsecTunnelProxyIdIpv4(d *schema.ResourceData) (string, ipv4.Entry) {
	tun := d.Get("ipsec_tunnel").(string)
	o := ipv4.Entry{
		Name:              d.Get("name").(string),
		Local:             d.Get("local").(string),
		Remote:            d.Get("remote").(string),
		ProtocolAny:       d.Get("protocol_any").(bool),
		ProtocolNumber:    d.Get("protocol_number").(int),
		ProtocolTcpLocal:  d.Get("protocol_tcp_local").(int),
		ProtocolTcpRemote: d.Get("protocol_tcp_remote").(int),
		ProtocolUdpLocal:  d.Get("protocol_udp_local").(int),
		ProtocolUdpRemote: d.Get("protocol_udp_remote").(int),
	}

	return tun, o
}

func parseIpsecTunnelProxyIdIpv4Id(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildIpsecTunnelProxyIdIpv4Id(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	tun, o := parseIpsecTunnelProxyIdIpv4(d)

	if err := fw.Network.IpsecTunnelProxyId.Set(tun, o); err != nil {
		return err
	}

	d.SetId(buildIpsecTunnelProxyIdIpv4Id(tun, o.Name))
	return readIpsecTunnelProxyIdIpv4(d, meta)
}

func readIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	tun, name := parseIpsecTunnelProxyIdIpv4Id(d.Id())

	o, err := fw.Network.IpsecTunnelProxyId.Get(tun, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("ipsec_tunnel", tun)
	d.Set("local", o.Local)
	d.Set("remote", o.Remote)
	d.Set("protocol_any", o.ProtocolAny)
	d.Set("protocol_number", o.ProtocolNumber)
	d.Set("protocol_tcp_local", o.ProtocolTcpLocal)
	d.Set("protocol_tcp_remote", o.ProtocolTcpRemote)
	d.Set("protocol_udp_local", o.ProtocolUdpLocal)
	d.Set("protocol_udp_remote", o.ProtocolUdpRemote)

	return nil
}

func updateIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	tun, o := parseIpsecTunnelProxyIdIpv4(d)

	lo, err := fw.Network.IpsecTunnelProxyId.Get(tun, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.IpsecTunnelProxyId.Edit(tun, lo); err != nil {
		return err
	}

	return readIpsecTunnelProxyIdIpv4(d, meta)
}

func deleteIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	tun, name := parseIpsecTunnelProxyIdIpv4Id(d.Id())

	err := fw.Network.IpsecTunnelProxyId.Delete(tun, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
