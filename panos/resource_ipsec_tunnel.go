package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ipsectunnel"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpsecTunnel() *schema.Resource {
	return &schema.Resource{
		Create: createIpsecTunnel,
		Read:   readIpsecTunnel,
		Update: updateIpsecTunnel,
		Delete: deleteIpsecTunnel,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tunnel_interface": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"anti_replay": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_ipv6": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"copy_tos": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"copy_flow_label": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsectunnel.TypeAutoKey,
				ValidateFunc: validateStringIn(ipsectunnel.TypeAutoKey, ipsectunnel.TypeManualKey, ipsectunnel.TypeGlobalProtectSatellite),
			},
			"ak_ike_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ak_ipsec_crypto_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_spi": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_remote_spi": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_address_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_local_address_floating_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_remote_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mk_protocol": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkProtocolEsp, ipsectunnel.MkProtocolAh),
			},
			"mk_auth_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkAuthTypeMd5, ipsectunnel.MkAuthTypeSha1, ipsectunnel.MkAuthTypeSha256, ipsectunnel.MkAuthTypeSha384, ipsectunnel.MkAuthTypeSha512, ipsectunnel.MkAuthTypeNone),
			},
			"mk_auth_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mk_auth_key_enc": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"mk_esp_encryption_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsectunnel.MkEspEncryptionDes, ipsectunnel.MkEspEncryption3des, ipsectunnel.MkEspEncryptionAes128, ipsectunnel.MkEspEncryptionAes192, ipsectunnel.MkEspEncryptionAes256, ipsectunnel.MkEspEncryptionNull),
			},
			"mk_esp_encryption_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mk_esp_encryption_key_enc": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"gps_interface": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_portal_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_prefer_ipv6": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gps_interface_ip_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_ip_ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_floating_ip_ipv4": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_interface_floating_ip_ipv6": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_publish_routes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gps_publish_connected_routes": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gps_local_certificate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"gps_certificate_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_tunnel_monitor": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tunnel_monitor_destination_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_source_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tunnel_monitor_proxy_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func parseIpsecTunnel(d *schema.ResourceData) ipsectunnel.Entry {
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

	return o
}

func createIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIpsecTunnel(d)

	if err = fw.Network.IpsecTunnel.Set(o); err != nil {
		return err
	}

	eo, err := fw.Network.IpsecTunnel.Get(o.Name)
	if err != nil {
		return err
	}

	d.SetId(o.Name)
	d.Set("mk_auth_key_enc", eo.MkAuthKey)
	d.Set("mk_esp_encryption_key_enc", eo.MkEspEncryptionKey)
	return readIpsecTunnel(d, meta)
}

func readIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.IpsecTunnel.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

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

func updateIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIpsecTunnel(d)

	lo, err := fw.Network.IpsecTunnel.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.IpsecTunnel.Edit(lo); err != nil {
		return err
	}
	eo, err := fw.Network.IpsecTunnel.Get(o.Name)
	if err != nil {
		return err
	}

	d.Set("mk_auth_key_enc", eo.MkAuthKey)
	d.Set("mk_esp_encryption_key_enc", eo.MkEspEncryptionKey)
	return readIpsecTunnel(d, meta)
}

func deleteIpsecTunnel(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.IpsecTunnel.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
