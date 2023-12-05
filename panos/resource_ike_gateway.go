package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ikegw"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIkeGateway() *schema.Resource {
	return &schema.Resource{
		Create: createIkeGateway,
		Read:   readIkeGateway,
		Update: updateIkeGateway,
		Delete: deleteIkeGateway,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ikegw.Ikev1,
				ValidateFunc: validateStringIn(ikegw.Ikev1, ikegw.Ikev2, ikegw.Ikev2Preferred),
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"peer_ip_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ikegw.PeerTypeIp,
				ValidateFunc: validateStringIn(ikegw.PeerTypeIp, ikegw.PeerTypeDynamic, ikegw.PeerTypeFqdn),
			},
			"peer_ip_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local_ip_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ikegw.LocalTypeIp, ikegw.LocalTypeFloatingIp, ""),
			},
			"local_ip_address_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ikegw.AuthPreSharedKey,
				ValidateFunc: validateStringIn(ikegw.AuthPreSharedKey, ikegw.AuthCertificate),
			},
			"pre_shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"pre_shared_key_enc": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"local_id_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ikegw.IdTypeIpAddress, ikegw.IdTypeFqdn, ikegw.IdTypeUfqdn, ikegw.IdTypeKeyId, ikegw.IdTypeDn, ""),
			},
			"local_id_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_id_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ikegw.IdTypeIpAddress, ikegw.IdTypeFqdn, ikegw.IdTypeUfqdn, ikegw.IdTypeKeyId, ikegw.IdTypeDn, ""),
			},
			"peer_id_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"peer_id_check": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ikegw.PeerIdCheckExact, ikegw.PeerIdCheckWildcard),
			},
			"local_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_enable_hash_and_url": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cert_base_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_use_management_as_source": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cert_permit_payload_mismatch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cert_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_enable_strict_validation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_passive_mode": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_nat_traversal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nat_traversal_keep_alive": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nat_traversal_enable_udp_checksum": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_fragmentation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ikev1_exchange_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ikev1_crypto_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_dead_peer_detection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"dead_peer_detection_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dead_peer_detection_retry": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ikev2_crypto_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ikev2_cookie_validation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_liveness_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"liveness_check_interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseIkeGateway(d *schema.ResourceData) ikegw.Entry {
	o := ikegw.Entry{
		Name:                          d.Get("name").(string),
		Version:                       d.Get("version").(string),
		EnableIpv6:                    d.Get("enable_ipv6").(bool),
		Disabled:                      d.Get("disabled").(bool),
		PeerIpType:                    d.Get("peer_ip_type").(string),
		PeerIpValue:                   d.Get("peer_ip_value").(string),
		Interface:                     d.Get("interface").(string),
		LocalIpAddressType:            d.Get("local_ip_address_type").(string),
		LocalIpAddressValue:           d.Get("local_ip_address_value").(string),
		AuthType:                      d.Get("auth_type").(string),
		PreSharedKey:                  d.Get("pre_shared_key").(string),
		LocalIdType:                   d.Get("local_id_type").(string),
		LocalIdValue:                  d.Get("local_id_value").(string),
		PeerIdType:                    d.Get("peer_id_type").(string),
		PeerIdValue:                   d.Get("peer_id_value").(string),
		PeerIdCheck:                   d.Get("peer_id_check").(string),
		LocalCert:                     d.Get("local_cert").(string),
		CertEnableHashAndUrl:          d.Get("cert_enable_hash_and_url").(bool),
		CertBaseUrl:                   d.Get("cert_base_url").(string),
		CertUseManagementAsSource:     d.Get("cert_use_management_as_source").(bool),
		CertPermitPayloadMismatch:     d.Get("cert_permit_payload_mismatch").(bool),
		CertProfile:                   d.Get("cert_profile").(string),
		CertEnableStrictValidation:    d.Get("cert_enable_strict_validation").(bool),
		EnablePassiveMode:             d.Get("enable_passive_mode").(bool),
		EnableNatTraversal:            d.Get("enable_nat_traversal").(bool),
		NatTraversalKeepAlive:         d.Get("nat_traversal_keep_alive").(int),
		NatTraversalEnableUdpChecksum: d.Get("nat_traversal_enable_udp_checksum").(bool),
		EnableFragmentation:           d.Get("enable_fragmentation").(bool),
		Ikev1ExchangeMode:             d.Get("ikev1_exchange_mode").(string),
		Ikev1CryptoProfile:            d.Get("ikev1_crypto_profile").(string),
		EnableDeadPeerDetection:       d.Get("enable_dead_peer_detection").(bool),
		DeadPeerDetectionInterval:     d.Get("dead_peer_detection_interval").(int),
		DeadPeerDetectionRetry:        d.Get("dead_peer_detection_retry").(int),
		Ikev2CryptoProfile:            d.Get("ikev2_crypto_profile").(string),
		Ikev2CookieValidation:         d.Get("ikev2_cookie_validation").(bool),
		EnableLivenessCheck:           d.Get("enable_liveness_check").(bool),
		LivenessCheckInterval:         d.Get("liveness_check_interval").(int),
	}

	return o
}

func createIkeGateway(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIkeGateway(d)

	if err = fw.Network.IkeGateway.Set(o); err != nil {
		return err
	}

	eo, err := fw.Network.IkeGateway.Get(o.Name)
	if err != nil {
		return err
	}

	d.SetId(o.Name)
	d.Set("pre_shared_key_enc", eo.PreSharedKey)
	return readIkeGateway(d, meta)
}

func readIkeGateway(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.IkeGateway.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("version", o.Version)
	d.Set("enable_ipv6", o.EnableIpv6)
	d.Set("disabled", o.Disabled)
	d.Set("peer_ip_type", o.PeerIpType)
	d.Set("peer_ip_value", o.PeerIpValue)
	d.Set("interface", o.Interface)
	d.Set("local_ip_address_type", o.LocalIpAddressType)
	d.Set("local_ip_address_value", o.LocalIpAddressValue)
	d.Set("auth_type", o.AuthType)
	d.Set("local_id_type", o.LocalIdType)
	d.Set("local_id_value", o.LocalIdValue)
	d.Set("peer_id_type", o.PeerIdType)
	d.Set("peer_id_value", o.PeerIdValue)
	d.Set("peer_id_check", o.PeerIdCheck)
	d.Set("local_cert", o.LocalCert)
	d.Set("cert_enable_hash_and_url", o.CertEnableHashAndUrl)
	d.Set("cert_base_url", o.CertBaseUrl)
	d.Set("cert_use_management_as_source", o.CertUseManagementAsSource)
	d.Set("cert_permit_payload_mismatch", o.CertPermitPayloadMismatch)
	d.Set("cert_profile", o.CertProfile)
	d.Set("cert_enable_strict_validation", o.CertEnableStrictValidation)
	d.Set("enable_passive_mode", o.EnablePassiveMode)
	d.Set("enable_nat_traversal", o.EnableNatTraversal)
	d.Set("nat_traversal_keep_alive", o.NatTraversalKeepAlive)
	d.Set("nat_traversal_enable_udp_checksum", o.NatTraversalEnableUdpChecksum)
	d.Set("enable_fragmentation", o.EnableFragmentation)
	d.Set("ikev1_exchange_mode", o.Ikev1ExchangeMode)
	d.Set("ikev1_crypto_profile", o.Ikev1CryptoProfile)
	d.Set("enable_dead_peer_detection", o.EnableDeadPeerDetection)
	d.Set("dead_peer_detection_interval", o.DeadPeerDetectionInterval)
	d.Set("dead_peer_detection_retry", o.DeadPeerDetectionRetry)
	d.Set("ikev2_crypto_profile", o.Ikev2CryptoProfile)
	d.Set("ikev2_cookie_validation", o.Ikev2CookieValidation)
	d.Set("enable_liveness_check", o.EnableLivenessCheck)
	d.Set("liveness_check_interval", o.LivenessCheckInterval)

	if d.Get("pre_shared_key_enc").(string) != o.PreSharedKey {
		d.Set("pre_shared_key", "")
	}

	return nil
}

func updateIkeGateway(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIkeGateway(d)

	lo, err := fw.Network.IkeGateway.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.IkeGateway.Edit(lo); err != nil {
		return err
	}
	eo, err := fw.Network.IkeGateway.Get(o.Name)
	if err != nil {
		return err
	}

	d.Set("pre_shared_key_enc", eo.PreSharedKey)
	return readIkeGateway(d, meta)
}

func deleteIkeGateway(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.IkeGateway.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
