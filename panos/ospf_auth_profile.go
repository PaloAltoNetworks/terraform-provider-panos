package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/profile/auth"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceOspfAuthProfiles() *schema.Resource {
	s := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router name",
		},
	}

	for key, value := range listingSchema() {
		s[key] = value
	}

	return &schema.Resource{
		Read: readDataSourceOspfAuthProfiles,

		Schema: s,
	}
}

func readDataSourceOspfAuthProfiles(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	vr := d.Get("virtual_router").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vr
		listing, err = con.Network.OspfAuthProfile.GetList(vr)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = base64Encode([]string{
			tmpl, ts, vr,
		})
		listing, err = con.Network.OspfAuthProfile.GetList(tmpl, ts, vr)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Resource.
func resourceOspfAuthProfile() *schema.Resource {
	return &schema.Resource{
		Create: createOspfAuthProfile,
		Read:   readOspfAuthProfile,
		Update: updateOspfAuthProfile,
		Delete: deleteOspfAuthProfile,

		Schema: ospfAuthProfileSchema(true),
	}
}

func createOspfAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	vr := d.Get("virtual_router").(string)
	o := loadOspfAuthProfile(d)
	var eo auth.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildFirewallOspfAuthProfileId(vr, o.Name)
		if err = con.Network.OspfAuthProfile.Set(vr, o); err == nil {
			eo, err = con.Network.OspfAuthProfile.Get(vr, o.Name)
		}
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		id = buildPanoramaOspfAuthProfileId(tmpl, ts, vr, o.Name)
		if err = con.Network.OspfAuthProfile.Set(tmpl, ts, vr, o); err == nil {
			eo, err = con.Network.OspfAuthProfile.Get(tmpl, ts, vr, o.Name)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	// Set encrypted values.
	if err = saveOspfAuthProfileEncryptedFields(d, o, eo); err != nil {
		return err
	}

	return readOspfAuthProfile(d, meta)
}

func readOspfAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o auth.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAuthProfileId(d.Id())
		o, err = con.Network.OspfAuthProfile.Get(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAuthProfileId(d.Id())
		o, err = con.Network.OspfAuthProfile.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveOspfAuthProfile(d, o)
	return nil
}

func updateOspfAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var eo auth.Entry
	o := loadOspfAuthProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAuthProfileId(d.Id())
		lo, err := con.Network.OspfAuthProfile.Get(vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAuthProfile.Edit(vr, o); err != nil {
			return err
		}
		eo, err = con.Network.OspfAuthProfile.Get(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAuthProfileId(d.Id())
		lo, err := con.Network.OspfAuthProfile.Get(tmpl, ts, vr, name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Network.OspfAuthProfile.Edit(tmpl, ts, vr, o); err != nil {
			return err
		}
		eo, err = con.Network.OspfAuthProfile.Get(tmpl, ts, vr, name)
	}

	if err != nil {
		return err
	}

	// Save encrypted fields.
	if err = saveOspfAuthProfileEncryptedFields(d, o, eo); err != nil {
		return err
	}

	return readOspfAuthProfile(d, meta)
}

func deleteOspfAuthProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vr, name := parseFirewallOspfAuthProfileId(d.Id())
		err = con.Network.OspfAuthProfile.Delete(vr, name)
	case *pango.Panorama:
		tmpl, ts, vr, name := parsePanoramaOspfAuthProfileId(d.Id())
		err = con.Network.OspfAuthProfile.Delete(tmpl, ts, vr, name)
	}

	if err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func ospfAuthProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"virtual_router": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The virtual router",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name",
			ForceNew:    true,
		},
		"auth_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The auth type",
			Default:     auth.AuthTypePassword,
			ValidateFunc: validateStringIn(
				auth.AuthTypePassword,
				auth.AuthTypeMd5,
			),
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The simple password",
			Sensitive:   true,
		},
		"password_enc": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Encrypted form of the simple password",
			Sensitive:   true,
		},
		"md5_key": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of MD5 key specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key_id": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "MD5 key ID",
					},
					"key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "MD5 key",
						Sensitive:   true,
					},
					"preferred": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Preferred key",
					},
				},
			},
		},
		"md5_keys_enc": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The encrypted form of the MD5 keys",
			Sensitive:   true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"raw": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"enc": {
						Type:      schema.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "virtual_router", "name"})
	}

	return ans
}

func loadOspfAuthProfile(d *schema.ResourceData) auth.Entry {
	var keys []auth.Md5Key
	if list := d.Get("md5_key").([]interface{}); len(list) > 0 {
		keys = make([]auth.Md5Key, 0, len(list))
		for i := range list {
			x := list[i].(map[string]interface{})
			keys = append(keys, auth.Md5Key{
				KeyId:     x["key_id"].(int),
				Key:       x["key"].(string),
				Preferred: x["preferred"].(bool),
			})
		}
	}

	return auth.Entry{
		Name:     d.Get("name").(string),
		AuthType: d.Get("auth_type").(string),
		Password: d.Get("password").(string),
		Md5Keys:  keys,
	}
}

func saveOspfAuthProfile(d *schema.ResourceData, o auth.Entry) {
	d.Set("name", o.Name)
	d.Set("auth_type", o.AuthType)

	pe := d.Get("password_enc").(string)
	if o.Password != pe {
		d.Set("password", "")
	}

	if len(o.Md5Keys) == 0 {
		d.Set("md5_key", nil)
	} else {
		kl := d.Get("md5_keys_enc").([]interface{})
		list := make([]interface{}, 0, len(o.Md5Keys))
		for i := range o.Md5Keys {
			var enc, raw string
			x := o.Md5Keys[i]
			if i < len(kl) {
				kli := kl[i].(map[string]interface{})
				enc = kli["enc"].(string)
				raw = kli["raw"].(string)
			}
			item := map[string]interface{}{
				"key_id":    x.KeyId,
				"preferred": x.Preferred,
			}
			if x.Key == enc {
				item["key"] = raw
			} else {
				item["key"] = ""
			}
			list = append(list, item)
		}
		if err := d.Set("md5_key", list); err != nil {
			log.Printf("[WARN] Error setting 'md5_key' for %q: %s", d.Id(), err)
		}
	}
}

func saveOspfAuthProfileEncryptedFields(d *schema.ResourceData, spec, live auth.Entry) error {
	d.Set("password_enc", live.Password)
	keys := make([]interface{}, 0, len(live.Md5Keys))
	if len(spec.Md5Keys) != len(live.Md5Keys) {
		return fmt.Errorf("Expected %d md5 keys, got %d", len(spec.Md5Keys), len(live.Md5Keys))
	}
	for i := range spec.Md5Keys {
		keys = append(keys, map[string]interface{}{
			"raw": spec.Md5Keys[i].Key,
			"enc": live.Md5Keys[i].Key,
		})
	}
	if err := d.Set("md5_keys_enc", keys); err != nil {
		log.Printf("[WARN] Error setting 'md5_keys_enc' for %q: %s", d.Id(), err)
	}

	return nil
}

// Id functions.
func parseFirewallOspfAuthProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaOspfAuthProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildFirewallOspfAuthProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaOspfAuthProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
