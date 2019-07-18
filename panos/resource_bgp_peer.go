package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: createBgpPeer,
		Read:   readBgpPeer,
		Update: updateBgpPeer,
		Delete: deleteBgpPeer,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: bgpPeerSchema(false),
	}
}

func bgpPeerSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"virtual_router": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"bgp_peer_group": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"enable": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"peer_as": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"local_address_interface": {
			Type:     schema.TypeString,
			Required: true,
		},
		"local_address_ip": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"peer_address_ip": {
			Type:     schema.TypeString,
			Required: true,
		},
		"reflector_client": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", peer.ReflectorClientNonClient, peer.ReflectorClientClient, peer.ReflectorClientMeshedClient),
		},
		"peering_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      peer.PeeringTypeUnspecified,
			ValidateFunc: validateStringIn(peer.PeeringTypeUnspecified, peer.PeeringTypeBilateral),
		},
		"max_prefixes": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "5000",
		},
		"auth_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"keep_alive_interval": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  30,
		},
		"multi_hop": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"open_delay_time": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"hold_time": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  90,
		},
		"idle_hold_time": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		"allow_incoming_connections": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"incoming_connections_remote_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"allow_outgoing_connections": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"outgoing_connections_local_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"bfd_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enable_mp_bgp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"address_family_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", peer.AddressFamilyTypeIpv4, peer.AddressFamilyTypeIpv6),
		},
		"subsequent_address_family_unicast": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"subsequent_address_family_multicast": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enable_sender_side_loop_detection": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"min_route_advertisement_interval": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseBgpPeer(d *schema.ResourceData) (string, string, peer.Entry) {
	vr := d.Get("virtual_router").(string)
	pg := d.Get("bgp_peer_group").(string)

	o := peer.Entry{
		Name:                             d.Get("name").(string),
		Enable:                           d.Get("enable").(bool),
		PeerAs:                           d.Get("peer_as").(string),
		LocalAddressInterface:            d.Get("local_address_interface").(string),
		LocalAddressIp:                   d.Get("local_address_ip").(string),
		PeerAddressIp:                    d.Get("peer_address_ip").(string),
		ReflectorClient:                  d.Get("reflector_client").(string),
		PeeringType:                      d.Get("peering_type").(string),
		MaxPrefixes:                      d.Get("max_prefixes").(string),
		AuthProfile:                      d.Get("auth_profile").(string),
		KeepAliveInterval:                d.Get("keep_alive_interval").(int),
		MultiHop:                         d.Get("multi_hop").(int),
		OpenDelayTime:                    d.Get("open_delay_time").(int),
		HoldTime:                         d.Get("hold_time").(int),
		IdleHoldTime:                     d.Get("idle_hold_time").(int),
		AllowIncomingConnections:         d.Get("allow_incoming_connections").(bool),
		IncomingConnectionsRemotePort:    d.Get("incoming_connections_remote_port").(int),
		AllowOutgoingConnections:         d.Get("allow_outgoing_connections").(bool),
		OutgoingConnectionsLocalPort:     d.Get("outgoing_connections_local_port").(int),
		BfdProfile:                       d.Get("bfd_profile").(string),
		EnableMpBgp:                      d.Get("enable_mp_bgp").(bool),
		AddressFamilyType:                d.Get("address_family_type").(string),
		SubsequentAddressFamilyUnicast:   d.Get("subsequent_address_family_unicast").(bool),
		SubsequentAddressFamilyMulticast: d.Get("subsequent_address_family_multicast").(bool),
		EnableSenderSideLoopDetection:    d.Get("enable_sender_side_loop_detection").(bool),
		MinRouteAdvertisementInterval:    d.Get("min_route_advertisement_interval").(int),
	}

	return vr, pg, o
}

func saveBgpPeer(d *schema.ResourceData, vr, pg string, o peer.Entry) {
	d.Set("virtual_router", vr)
	d.Set("bgp_peer_group", pg)

	d.Set("name", o.Name)
	d.Set("enable", o.Enable)
	d.Set("peer_as", o.PeerAs)
	d.Set("local_address_interface", o.LocalAddressInterface)
	d.Set("local_address_ip", o.LocalAddressIp)
	d.Set("peer_address_ip", o.PeerAddressIp)
	d.Set("reflector_client", o.ReflectorClient)
	d.Set("peering_type", o.PeeringType)
	d.Set("max_prefixes", o.MaxPrefixes)
	d.Set("auth_profile", o.AuthProfile)
	d.Set("keep_alive_interval", o.KeepAliveInterval)
	d.Set("multi_hop", o.MultiHop)
	d.Set("open_delay_time", o.OpenDelayTime)
	d.Set("hold_time", o.HoldTime)
	d.Set("idle_hold_time", o.IdleHoldTime)
	d.Set("allow_incoming_connections", o.AllowIncomingConnections)
	d.Set("incoming_connections_remote_port", o.IncomingConnectionsRemotePort)
	d.Set("allow_outgoing_connections", o.AllowOutgoingConnections)
	d.Set("outgoing_connections_local_port", o.OutgoingConnectionsLocalPort)
	d.Set("bfd_profile", o.BfdProfile)
	d.Set("enable_mp_bgp", o.EnableMpBgp)
	d.Set("address_family_type", o.AddressFamilyType)
	d.Set("subsequent_address_family_unicast", o.SubsequentAddressFamilyUnicast)
	d.Set("subsequent_address_family_multicast", o.SubsequentAddressFamilyMulticast)
	d.Set("enable_sender_side_loop_detection", o.EnableSenderSideLoopDetection)
	d.Set("min_route_advertisement_interval", o.MinRouteAdvertisementInterval)
}

func parseBgpPeerId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildBgpPeerId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createBgpPeer(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, pg, o := parseBgpPeer(d)

	if err := fw.Network.BgpPeer.Set(vr, pg, o); err != nil {
		return err
	}

	d.SetId(buildBgpPeerId(vr, pg, o.Name))
	return readBgpPeer(d, meta)
}

func readBgpPeer(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, pg, name := parseBgpPeerId(d.Id())

	o, err := fw.Network.BgpPeer.Get(vr, pg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveBgpPeer(d, vr, pg, o)

	return nil
}

func updateBgpPeer(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vr, pg, o := parseBgpPeer(d)

	lo, err := fw.Network.BgpPeer.Get(vr, pg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.BgpPeer.Edit(vr, pg, lo); err != nil {
		return err
	}

	return readBgpPeer(d, meta)
}

func deleteBgpPeer(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vr, pg, name := parseBgpPeerId(d.Id())

	err := fw.Network.BgpPeer.Delete(vr, pg, name)
	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
