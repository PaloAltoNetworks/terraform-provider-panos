package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ike"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaIkeCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaIkeCryptoProfile,
		Read:   readPanoramaIkeCryptoProfile,
		Update: updatePanoramaIkeCryptoProfile,
		Delete: deletePanoramaIkeCryptoProfile,

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
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template_stack"},
			},
			"template_stack": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"template"},
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
					Type: schema.TypeString,
					ValidateFunc: validateStringIn(
						ike.EncryptionDes,
						ike.Encryption3des,
						ike.EncryptionAes128,
						ike.EncryptionAes192,
						ike.EncryptionAes256,
						ike.EncryptionAes128Gcm,
						ike.EncryptionAes256Gcm,
					),
				},
			},
			"lifetime_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ike.TimeHours,
				ValidateFunc: validateStringIn(
					ike.TimeSeconds,
					ike.TimeMinutes,
					ike.TimeHours,
					ike.TimeDays,
				),
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

func parsePanoramaIkeCryptoProfile(d *schema.ResourceData) (string, string, ike.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	o := ike.Entry{
		Name:                   d.Get("name").(string),
		DhGroup:                asStringList(d.Get("dh_groups").([]interface{})),
		Authentication:         asStringList(d.Get("authentications").([]interface{})),
		Encryption:             asStringList(d.Get("encryptions").([]interface{})),
		LifetimeType:           d.Get("lifetime_type").(string),
		LifetimeValue:          d.Get("lifetime_value").(int),
		AuthenticationMultiple: d.Get("authentication_multiple").(int),
	}

	return tmpl, ts, o
}

func parsePanoramaIkeCryptoProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaIkeCryptoProfileId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func createPanoramaIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIkeCryptoProfile(d)

	if err := pano.Network.IkeCryptoProfile.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaIkeCryptoProfileId(tmpl, ts, o.Name))
	return readPanoramaIkeCryptoProfile(d, meta)
}

func readPanoramaIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIkeCryptoProfileId(d.Id())

	o, err := pano.Network.IkeCryptoProfile.Get(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", name)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)
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

func updatePanoramaIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaIkeCryptoProfile(d)

	lo, err := pano.Network.IkeCryptoProfile.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.IkeCryptoProfile.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaIkeCryptoProfile(d, meta)
}

func deletePanoramaIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaIkeCryptoProfileId(d.Id())

	err := pano.Network.IkeCryptoProfile.Delete(tmpl, ts, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
