package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ipsec"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIpsecCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createIpsecCryptoProfile,
		Read:   readIpsecCryptoProfile,
		Update: updateIpsecCryptoProfile,
		Delete: deleteIpsecCryptoProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsec.ProtocolEsp,
				ValidateFunc: validateStringIn(ipsec.ProtocolEsp, ipsec.ProtocolAh),
			},
			"authentications": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"encryptions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateStringIn(ipsec.EncryptionDes, ipsec.Encryption3des, ipsec.EncryptionAes128, ipsec.EncryptionAes192, ipsec.EncryptionAes256, ipsec.EncryptionAes128Gcm, ipsec.EncryptionAes256Gcm, ipsec.EncryptionNull),
				},
			},
			"dh_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringHasPrefix("group"),
			},
			"lifetime_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsec.TimeHours,
				ValidateFunc: validateStringIn(ipsec.TimeSeconds, ipsec.TimeMinutes, ipsec.TimeHours, ipsec.TimeDays),
			},
			"lifetime_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lifesize_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsec.SizeKb, ipsec.SizeMb, ipsec.SizeGb, ipsec.SizeTb),
			},
			"lifesize_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseIpsecCryptoProfile(d *schema.ResourceData) ipsec.Entry {
	o := ipsec.Entry{
		Name:           d.Get("name").(string),
		Protocol:       d.Get("protocol").(string),
		Authentication: asStringList(d.Get("authentications").([]interface{})),
		Encryption:     asStringList(d.Get("encryptions").([]interface{})),
		DhGroup:        d.Get("dh_group").(string),
		LifetimeType:   d.Get("lifetime_type").(string),
		LifetimeValue:  d.Get("lifetime_value").(int),
		LifesizeType:   d.Get("lifesize_type").(string),
		LifesizeValue:  d.Get("lifesize_value").(int),
	}

	return o
}

func createIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseIpsecCryptoProfile(d)

	if err := fw.Network.IpsecCryptoProfile.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readIpsecCryptoProfile(d, meta)
}

func readIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.IpsecCryptoProfile.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	if err = d.Set("authentications", o.Authentication); err != nil {
		log.Printf("[WARN] Error setting 'authentications' for %q: %s", d.Id(), err)
	}
	if err = d.Set("encryptions", o.Encryption); err != nil {
		log.Printf("[WARN] Error setting 'encryptions' for %q: %s", d.Id(), err)
	}
	d.Set("dh_group", o.DhGroup)
	d.Set("lifetime_type", o.LifetimeType)
	d.Set("lifetime_value", o.LifetimeValue)
	d.Set("lifesize_type", o.LifesizeType)
	d.Set("lifesize_value", o.LifesizeValue)

	return nil
}

func updateIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIpsecCryptoProfile(d)

	lo, err := fw.Network.IpsecCryptoProfile.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.IpsecCryptoProfile.Edit(lo); err != nil {
		return err
	}

	return readIpsecCryptoProfile(d, meta)
}

func deleteIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.IpsecCryptoProfile.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
