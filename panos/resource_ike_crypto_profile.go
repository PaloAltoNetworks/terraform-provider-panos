package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ike"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIkeCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createIkeCryptoProfile,
		Read:   readIkeCryptoProfile,
		Update: updateIkeCryptoProfile,
		Delete: deleteIkeCryptoProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dh_groups": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateStringHasPrefix("group"),
				},
			},
			"authentications": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"encryptions": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateStringIn(ike.EncryptionDes, ike.Encryption3des, ike.EncryptionAes128, ike.EncryptionAes192, ike.EncryptionAes256),
				},
			},
			"lifetime_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ike.TimeHours,
				ValidateFunc: validateStringIn(ike.TimeSeconds, ike.TimeMinutes, ike.TimeHours, ike.TimeDays),
			},
			"lifetime_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"authentication_multiple": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parseIkeCryptoProfile(d *schema.ResourceData) ike.Entry {
	o := ike.Entry{
		Name:                   d.Get("name").(string),
		DhGroup:                asStringList(d.Get("dh_groups").([]interface{})),
		Authentication:         asStringList(d.Get("authentications").([]interface{})),
		Encryption:             asStringList(d.Get("encryptions").([]interface{})),
		LifetimeType:           d.Get("lifetime_type").(string),
		LifetimeValue:          d.Get("lifetime_value").(int),
		AuthenticationMultiple: d.Get("authentication_multiple").(int),
	}

	return o
}

func createIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseIkeCryptoProfile(d)

	if err := fw.Network.IkeCryptoProfile.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readIkeCryptoProfile(d, meta)
}

func readIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	o, err := fw.Network.IkeCryptoProfile.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	if err = d.Set("dh_groups", o.DhGroup); err != nil {
		log.Printf("[WARN] Error setting 'dh_groups' for %q: %s", d.Id(), err)
	}
	if err = d.Set("authentications", o.Authentication); err != nil {
		log.Printf("[WARN] Error setting 'authentications' for %q: %s", d.Id(), err)
	}
	if err = d.Set("encryptions", o.Encryption); err != nil {
		log.Printf("[WARN] Error setting 'encryptions' for %q: %s", d.Id(), err)
	}
	d.Set("lifetime_type", o.LifetimeType)
	d.Set("lifetime_value", o.LifetimeValue)
	d.Set("authentication_multiple", o.AuthenticationMultiple)

	return nil
}

func updateIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseIkeCryptoProfile(d)

	lo, err := fw.Network.IkeCryptoProfile.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.IkeCryptoProfile.Edit(lo); err != nil {
		return err
	}

	return readIkeCryptoProfile(d, meta)
}

func deleteIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Get("name").(string)

	err := fw.Network.IkeCryptoProfile.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
