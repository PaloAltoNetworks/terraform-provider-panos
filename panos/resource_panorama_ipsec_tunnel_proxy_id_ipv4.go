package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ipsectunnel/proxyid/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaIpsecTunnelProxyIdIpv4() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaIpsecTunnelProxyIdIpv4,
		Read:   readPanoramaIpsecTunnelProxyIdIpv4,
		Update: updatePanoramaIpsecTunnelProxyIdIpv4,
		Delete: deletePanoramaIpsecTunnelProxyIdIpv4,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": {
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

func parsePanoramaIpsecTunnelProxyIdIpv4(d *schema.ResourceData) (string, string, string, ipv4.Entry) {
	tmpl := d.Get("template").(string)
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

	return tmpl, "", tun, o
}

func parsePanoramaIpsecTunnelProxyIdIpv4Id(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaIpsecTunnelProxyIdIpv4Id(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func createPanoramaIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, tun, o := parsePanoramaIpsecTunnelProxyIdIpv4(d)

	if err := pano.Network.IpsecTunnelProxyId.Set(tmpl, ts, tun, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaIpsecTunnelProxyIdIpv4Id(tmpl, ts, tun, o.Name))
	return readPanoramaIpsecTunnelProxyIdIpv4(d, meta)
}

func readPanoramaIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, tun, name := parsePanoramaIpsecTunnelProxyIdIpv4Id(d.Id())

	o, err := pano.Network.IpsecTunnelProxyId.Get(tmpl, ts, tun, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
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

func updatePanoramaIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, tun, o := parsePanoramaIpsecTunnelProxyIdIpv4(d)

	lo, err := pano.Network.IpsecTunnelProxyId.Get(tmpl, ts, tun, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.IpsecTunnelProxyId.Edit(tmpl, ts, tun, lo); err != nil {
		return err
	}

	return readPanoramaIpsecTunnelProxyIdIpv4(d, meta)
}

func deletePanoramaIpsecTunnelProxyIdIpv4(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, tun, name := parsePanoramaIpsecTunnelProxyIdIpv4Id(d.Id())

	err := pano.Network.IpsecTunnelProxyId.Delete(tmpl, ts, tun, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
