package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/ike"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO(gfreeman): Add the standard array of data sources.

// Resource.
func resourceIkeCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createIkeCryptoProfile,
		Read:   readIkeCryptoProfile,
		Update: updateIkeCryptoProfile,
		Delete: deleteIkeCryptoProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: ikeCryptoProfileSchema(true),
	}
}

func createIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadIkeCryptoProfile(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	id := buildIkeCryptoProfileId(tmpl, ts, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.IkeCryptoProfile.Set(o)
	case *pango.Panorama:
		err = con.Network.IkeCryptoProfile.Set(tmpl, ts, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readIkeCryptoProfile(d, meta)
}

func readIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ike.Entry

	// Migrate the old ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 1 {
		d.SetId(buildIkeCryptoProfileId("", "", tok[0]))
	} else if len(tok) != 3 {
		return fmt.Errorf("Incorrect ID token length, expecting 1 or 3")
	}

	tmpl, ts, name := parseIkeCryptoProfileId(d.Id())

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.IkeCryptoProfile.Get(name)
	case *pango.Panorama:
		o, err = con.Network.IkeCryptoProfile.Get(tmpl, ts, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveIkeCryptoProfile(d, o)
	return nil
}

func updateIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadIkeCryptoProfile(d)

	tmpl, ts, _ := parseIkeCryptoProfileId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Network.IkeCryptoProfile.Get(o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.IkeCryptoProfile.Edit(lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Network.IkeCryptoProfile.Get(tmpl, ts, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.IkeCryptoProfile.Edit(tmpl, ts, lo); err != nil {
			return err
		}
	}

	return readIkeCryptoProfile(d, meta)
}

func deleteIkeCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, name := parseIkeCryptoProfileId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.IkeCryptoProfile.Delete(name)
	case *pango.Panorama:
		err = con.Network.IkeCryptoProfile.Delete(tmpl, ts, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func ikeCryptoProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
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
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "name"})
	}

	return ans
}

func loadIkeCryptoProfile(d *schema.ResourceData) ike.Entry {
	return ike.Entry{
		Name:                   d.Get("name").(string),
		DhGroup:                asStringList(d.Get("dh_groups").([]interface{})),
		Authentication:         asStringList(d.Get("authentications").([]interface{})),
		Encryption:             asStringList(d.Get("encryptions").([]interface{})),
		LifetimeType:           d.Get("lifetime_type").(string),
		LifetimeValue:          d.Get("lifetime_value").(int),
		AuthenticationMultiple: d.Get("authentication_multiple").(int),
	}
}

func saveIkeCryptoProfile(d *schema.ResourceData, o ike.Entry) {
	var err error

	d.Set("name", o.Name)
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
}

// Id functions.
func parseIkeCryptoProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildIkeCryptoProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
