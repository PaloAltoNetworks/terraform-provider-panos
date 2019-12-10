package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ipsectunnel"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaIpsecTunnel() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaIpsecTunnel,
		Read:   readPanoramaIpsecTunnel,
		Update: updatePanoramaIpsecTunnel,
		Delete: deletePanoramaIpsecTunnel,

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
			"tunnel_interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"anti_replay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"copy_tos": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"copy_flow_label": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsectunnel.TypeAutoKey,
				ValidateFunc: validateStringIn(ipsectunnel.TypeAutoKey, ipsectunnel.TypeManualKey, ipsectunnel.TypeGlobalProtectSatellite),
			},
			"ak_ike_gateway": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ak_ipsec_crypto_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_spi": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_remote_spi": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_address_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_address_floating_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_remote_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkProtocolEsp, ipsectunnel.MkProtocolAh),
			},
			"mk_auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkAuthTypeMd5, ipsectunnel.MkAuthTypeSha1, ipsectunnel.MkAuthTypeSha256, ipsectunnel.MkAuthTypeSha384, ipsectunnel.MkAuthTypeSha512, ipsectunnel.MkAuthTypeNone),
			},
			"mk_auth_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mk_auth_key_enc": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"mk_esp_encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkEspEncryptionDes, ipsectunnel.MkEspEncryption3des, ipsectunnel.MkEspEncryptionAes128, ipsectunnel.MkEspEncryptionAes192, ipsectunnel.MkEspEncryptionAes256, ipsectunnel.MkEspEncryptionNull),
			},
			"mk_esp_encryption_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mk_esp_encryption_key_enc": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"gps_interface": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_portal_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_prefer_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gps_interface_ip_ipv4": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_ip_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_floating_ip_ipv4": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_floating_ip_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_publish_routes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gps_publish_connected_routes": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gps_local_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_certificate_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_tunnel_monitor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tunnel_monitor_destination_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_proxy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parsePanoramaIpsecTunnel(d *schema.ResourceData) (string, string, ipsectunnel.Entry) {
	tmpl := d.Get("template").(string)

	o := ipsectunnel.Entry{
		Name:                       d.Get("name").(string),
		TunnelInterface:            d.Get("tunnel_interface").(string),
		AntiReplay:                 d.Get("anti_replay").(bool),
		EnableIpv6:                 d.Get("enable_ipv6").(bool),
		CopyTos:                    d.Get("copy_tos").(bool),
		CopyFlowLabel:              d.Get("copy_flow_label").(bool),
		Disabled:                   d.Get("disabled").(bool),
		Type:                       d.Get("type").(string),
		AkIkeGateway:               d.Get("ak_ike_gateway").(string),
		AkIpsecCryptoProfile:       d.Get("ak_ipsec_crypto_profile").(string),
		MkLocalSpi:                 d.Get("mk_local_spi").(string),
		MkRemoteSpi:                d.Get("mk_remote_spi").(string),
		MkInterface:                d.Get("mk_interface").(string),
		MkLocalAddressIp:           d.Get("mk_local_address_ip").(string),
		MkLocalAddressFloatingIp:   d.Get("mk_local_address_floating_ip").(string),
		MkRemoteAddress:            d.Get("mk_remote_address").(string),
		MkProtocol:                 d.Get("mk_protocol").(string),
		MkAuthType:                 d.Get("mk_auth_type").(string),
		MkAuthKey:                  d.Get("mk_auth_key").(string),
		MkEspEncryptionType:        d.Get("mk_esp_encryption_type").(string),
		MkEspEncryptionKey:         d.Get("mk_esp_encryption_key").(string),
		GpsInterface:               d.Get("gps_interface").(string),
		GpsPortalAddress:           d.Get("gps_portal_address").(string),
		GpsPreferIpv6:              d.Get("gps_prefer_ipv6").(bool),
		GpsInterfaceIpIpv4:         d.Get("gps_interface_ip_ipv4").(string),
		GpsInterfaceIpIpv6:         d.Get("gps_interface_ip_ipv6").(string),
		GpsInterfaceFloatingIpIpv4: d.Get("gps_interface_floating_ip_ipv4").(string),
		GpsInterfaceFloatingIpIpv6: d.Get("gps_interface_floating_ip_ipv6").(string),
		GpsPublishConnectedRoutes:  d.Get("gps_publish_connected_routes").(bool),
		GpsPublishRoutes:           asStringList(d.Get("gps_publish_routes").([]interface{})),
		GpsLocalCertificate:        d.Get("gps_local_certificate").(string),
		GpsCertificateProfile:      d.Get("gps_certificate_profile").(string),
		EnableTunnelMonitor:        d.Get("enable_tunnel_monitor").(bool),
		TunnelMonitorDestinationIp: d.Get("tunnel_monitor_destination_ip").(string),
		TunnelMonitorSourceIp:      d.Get("tunnel_monitor_source_ip").(string),
		TunnelMonitorProfile:       d.Get("tunnel_monitor_profile").(string),
		TunnelMonitorProxyId:       d.Get("tunnel_monitor_proxy_id").(string),
	}

	return tmpl, "", o
}

func buildPanoramaIpsecTunnelId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parsePanoramaIpsecTunnelId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func createPanoramaIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIpsecTunnel(d)

	if err = pano.Network.IpsecTunnel.Set(tmpl, ts, o); err != nil {
		return err
	}

	eo, err := pano.Network.IpsecTunnel.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}

	d.SetId(buildPanoramaIpsecTunnelId(tmpl, ts, o.Name))
	d.Set("mk_auth_key_enc", eo.MkAuthKey)
	d.Set("mk_esp_encryption_key_enc", eo.MkEspEncryptionKey)
	return readPanoramaIpsecTunnel(d, meta)
}

func readPanoramaIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIpsecTunnelId(d.Id())

	o, err := pano.Network.IpsecTunnel.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("tunnel_interface", o.TunnelInterface)
	d.Set("anti_replay", o.AntiReplay)
	d.Set("enable_ipv6", o.EnableIpv6)
	d.Set("copy_tos", o.CopyTos)
	d.Set("copy_flow_label", o.CopyFlowLabel)
	d.Set("disabled", o.Disabled)
	d.Set("type", o.Type)
	d.Set("ak_ike_gateway", o.AkIkeGateway)
	d.Set("ak_ipsec_crypto_profile", o.AkIpsecCryptoProfile)
	d.Set("mk_local_spi", o.MkLocalSpi)
	d.Set("mk_remote_spi", o.MkRemoteSpi)
	d.Set("mk_interface", o.MkInterface)
	d.Set("mk_local_address_ip", o.MkLocalAddressIp)
	d.Set("mk_local_address_floating_ip", o.MkLocalAddressFloatingIp)
	d.Set("mk_remote_address", o.MkRemoteAddress)
	d.Set("mk_protocol", o.MkProtocol)
	d.Set("mk_auth_type", o.MkAuthType)
	d.Set("mk_esp_encryption_type", o.MkEspEncryptionType)
	d.Set("gps_interface", o.GpsInterface)
	d.Set("gps_portal_address", o.GpsPortalAddress)
	d.Set("gps_prefer_ipv6", o.GpsPreferIpv6)
	d.Set("gps_interface_ip_ipv4", o.GpsInterfaceIpIpv4)
	d.Set("gps_interface_ip_ipv6", o.GpsInterfaceIpIpv6)
	d.Set("gps_interface_floating_ip_ipv6", o.GpsInterfaceFloatingIpIpv6)
	d.Set("gps_interface_floating_ip_ipv4", o.GpsInterfaceFloatingIpIpv4)
	d.Set("gps_publish_connected_routes", o.GpsPublishConnectedRoutes)
	if err = d.Set("gps_publish_routes", o.GpsPublishRoutes); err != nil {
		log.Printf("[WARN] Error setting 'gps_publish_routes' field for %q: %s", d.Id(), err)
	}
	d.Set("gps_local_certificate", o.GpsLocalCertificate)
	d.Set("gps_certificate_profile", o.GpsCertificateProfile)
	d.Set("enable_tunnel_monitor", o.EnableTunnelMonitor)
	d.Set("tunnel_monitor_destination_ip", o.TunnelMonitorDestinationIp)
	d.Set("tunnel_monitor_source_ip", o.TunnelMonitorSourceIp)
	d.Set("tunnel_monitor_profile", o.TunnelMonitorProfile)
	d.Set("tunnel_monitor_proxy_id", o.TunnelMonitorProxyId)

	if d.Get("mk_auth_key_enc").(string) != o.MkAuthKey {
		d.Set("mk_auth_key", "")
	}
	if d.Get("mk_esp_encryption_key_enc").(string) != o.MkEspEncryptionKey {
		d.Set("mk_esp_encryption_key", "")
	}

	return nil
}

func updatePanoramaIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIpsecTunnel(d)

	lo, err := pano.Network.IpsecTunnel.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.IpsecTunnel.Edit(tmpl, ts, lo); err != nil {
		return err
	}
	eo, err := pano.Network.IpsecTunnel.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}

	d.Set("mk_auth_key_enc", eo.MkAuthKey)
	d.Set("mk_esp_encryption_key_enc", eo.MkEspEncryptionKey)
	return readPanoramaIpsecTunnel(d, meta)
}

func deletePanoramaIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIpsecTunnelId(d.Id())

	err := pano.Network.IpsecTunnel.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
