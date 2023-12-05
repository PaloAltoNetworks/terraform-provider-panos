package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/tunnel/gre"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGreTunnel() *schema.Resource {
	return &schema.Resource{
		Create: createGreTunnel,
		Read:   readGreTunnel,
		Update: updateGreTunnel,
		Delete: deleteGreTunnel,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: greTunnelSchema(false),
	}
}

func greTunnelSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"interface": {
			Type:     schema.TypeString,
			Required: true,
		},
		"local_address_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      gre.LocalAddressTypeIp,
			ValidateFunc: validateStringIn(gre.LocalAddressTypeIp, gre.LocalAddressTypeFloatingIp),
		},
		"local_address_value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"peer_address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tunnel_interface": {
			Type:     schema.TypeString,
			Required: true,
		},
		"ttl": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"copy_tos": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"enable_keep_alive": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"keep_alive_interval": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"keep_alive_retry": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"keep_alive_hold_timer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"disabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	if p {
		ans["template"] = templateSchema(false)
	}

	return ans
}

func parseGreTunnel(d *schema.ResourceData) gre.Entry {
	o := loadGreTunnel(d)

	return o
}

func loadGreTunnel(d *schema.ResourceData) gre.Entry {
	return gre.Entry{
		Name:               d.Get("name").(string),
		Interface:          d.Get("interface").(string),
		LocalAddressType:   d.Get("local_address_type").(string),
		LocalAddressValue:  d.Get("local_address_value").(string),
		PeerAddress:        d.Get("peer_address").(string),
		TunnelInterface:    d.Get("tunnel_interface").(string),
		Ttl:                d.Get("ttl").(int),
		CopyTos:            d.Get("copy_tos").(bool),
		EnableKeepAlive:    d.Get("enable_keep_alive").(bool),
		KeepAliveInterval:  d.Get("keep_alive_interval").(int),
		KeepAliveRetry:     d.Get("keep_alive_retry").(int),
		KeepAliveHoldTimer: d.Get("keep_alive_hold_timer").(int),
		Disabled:           d.Get("disabled").(bool),
	}
}

func saveGreTunnel(d *schema.ResourceData, o gre.Entry) {
	d.Set("name", o.Name)
	d.Set("interface", o.Interface)
	d.Set("local_address_type", o.LocalAddressType)
	d.Set("local_address_value", o.LocalAddressValue)
	d.Set("peer_address", o.PeerAddress)
	d.Set("tunnel_interface", o.TunnelInterface)
	d.Set("ttl", o.Ttl)
	d.Set("copy_tos", o.CopyTos)
	d.Set("enable_keep_alive", o.EnableKeepAlive)
	d.Set("keep_alive_interval", o.KeepAliveInterval)
	d.Set("keep_alive_retry", o.KeepAliveRetry)
	d.Set("keep_alive_hold_timer", o.KeepAliveHoldTimer)
	d.Set("disabled", o.Disabled)
}

func createGreTunnel(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseGreTunnel(d)

	if err := fw.Network.GreTunnel.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readGreTunnel(d, meta)
}

func readGreTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Id()

	o, err := fw.Network.GreTunnel.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveGreTunnel(d, o)

	return nil
}

func updateGreTunnel(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseGreTunnel(d)

	lo, err := fw.Network.GreTunnel.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.GreTunnel.Edit(lo); err != nil {
		return err
	}

	return readGreTunnel(d, meta)
}

func deleteGreTunnel(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Id()

	err := fw.Network.GreTunnel.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
