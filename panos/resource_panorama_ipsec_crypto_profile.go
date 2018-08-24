package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ipsec"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaIpsecCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaIpsecCryptoProfile,
		Read:   readPanoramaIpsecCryptoProfile,
		Update: updatePanoramaIpsecCryptoProfile,
		Delete: deletePanoramaIpsecCryptoProfile,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template_stack"},
			},
			"template_stack": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template"},
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsec.ProtocolEsp,
				ValidateFunc: validateStringIn(ipsec.ProtocolEsp, ipsec.ProtocolAh),
			},
			"authentications": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"encryptions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateStringIn(ipsec.EncryptionDes, ipsec.Encryption3des, ipsec.EncryptionAes128, ipsec.EncryptionAes192, ipsec.EncryptionAes256, ipsec.EncryptionAes128Gcm, ipsec.EncryptionAes256Gcm, ipsec.EncryptionNull),
				},
			},
			"dh_group": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringHasPrefix("group"),
			},
			"lifetime_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ipsec.TimeHours,
				ValidateFunc: validateStringIn(ipsec.TimeSeconds, ipsec.TimeMinutes, ipsec.TimeHours, ipsec.TimeDays),
			},
			"lifetime_value": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lifesize_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn(ipsec.SizeKb, ipsec.SizeMb, ipsec.SizeGb, ipsec.SizeTb),
			},
			"lifesize_value": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func parsePanoramaIpsecCryptoProfile(d *schema.ResourceData) (string, string, ipsec.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

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

	return tmpl, ts, o
}

func parsePanoramaIpsecCryptoProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaIpsecCryptoProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createPanoramaIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIpsecCryptoProfile(d)

	if err := pano.Network.IpsecCryptoProfile.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaIpsecCryptoProfileId(tmpl, ts, o.Name))
	return readPanoramaIpsecCryptoProfile(d, meta)
}

func readPanoramaIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIpsecCryptoProfileId(d.Id())

	o, err := pano.Network.IpsecCryptoProfile.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
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

func updatePanoramaIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIpsecCryptoProfile(d)

	lo, err := pano.Network.IpsecCryptoProfile.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.IpsecCryptoProfile.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaIpsecCryptoProfile(d, meta)
}

func deletePanoramaIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIpsecCryptoProfileId(d.Id())

	err := pano.Network.IpsecCryptoProfile.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
