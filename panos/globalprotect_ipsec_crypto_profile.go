package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/gp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceGlobalProtectIpsecCryptoProfiles() *schema.Resource {
	s := listingSchema()
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceGlobalProtectIpsecCryptoProfilesRead,

		Schema: s,
	}
}

func dataSourceGlobalProtectIpsecCryptoProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts := d.Get("template").(string), d.Get("template_stack").(string)

	id := buildGlobalProtectIpsecCryptoProfileId(tmpl, ts, "")

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Network.GlobalProtectIpsecCryptoProfile.GetList()
	case *pango.Panorama:
		listing, err = con.Network.GlobalProtectIpsecCryptoProfile.GetList(tmpl, ts)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceGlobalProtectIpsecCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlobalProtectIpsecCryptoProfileRead,

		Schema: globalProtectIpsecCryptoProfileSchema(false),
	}
}

func dataSourceGlobalProtectIpsecCryptoProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o gp.Entry

	name := d.Get("name").(string)
	tmpl, ts := d.Get("template").(string), d.Get("template_stack").(string)
	id := buildGlobalProtectIpsecCryptoProfileId(tmpl, ts, name)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(name)
	case *pango.Panorama:
		o, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(tmpl, ts, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveGlobalProtectIpsecCryptoProfile(d, o)

	return nil
}

// Resource.
func resourceGlobalProtectIpsecCryptoProfile() *schema.Resource {
	return &schema.Resource{
		Create: createGlobalProtectIpsecCryptoProfile,
		Read:   readGlobalProtectIpsecCryptoProfile,
		Update: updateGlobalProtectIpsecCryptoProfile,
		Delete: deleteGlobalProtectIpsecCryptoProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: globalProtectIpsecCryptoProfileSchema(true),
	}
}

func createGlobalProtectIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadGlobalProtectIpsecCryptoProfile(d)
	tmpl, ts := d.Get("template").(string), d.Get("template_stack").(string)

	id := buildGlobalProtectIpsecCryptoProfileId(tmpl, ts, o.Name)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.GlobalProtectIpsecCryptoProfile.Set(o)
	case *pango.Panorama:
		err = con.Network.GlobalProtectIpsecCryptoProfile.Set(tmpl, ts, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readGlobalProtectIpsecCryptoProfile(d, meta)
}

func readGlobalProtectIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o gp.Entry

	tmpl, ts, name, err := parseGlobalProtectIpsecCryptoProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(name)
	case *pango.Panorama:
		o, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(tmpl, ts, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveGlobalProtectIpsecCryptoProfile(d, o)

	return nil
}

func updateGlobalProtectIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo gp.Entry
	o := loadGlobalProtectIpsecCryptoProfile(d)

	tmpl, ts, _, err := parseGlobalProtectIpsecCryptoProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Network.GlobalProtectIpsecCryptoProfile.Edit(lo)
		}
	case *pango.Panorama:
		lo, err = con.Network.GlobalProtectIpsecCryptoProfile.Get(tmpl, ts, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Network.GlobalProtectIpsecCryptoProfile.Edit(tmpl, ts, lo)
		}
	}

	if err != nil {
		return err
	}

	return readGlobalProtectIpsecCryptoProfile(d, meta)
}

func deleteGlobalProtectIpsecCryptoProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, name, err := parseGlobalProtectIpsecCryptoProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.GlobalProtectIpsecCryptoProfile.Delete(name)
	case *pango.Panorama:
		err = con.Network.GlobalProtectIpsecCryptoProfile.Delete(tmpl, ts, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func globalProtectIpsecCryptoProfileSchema(isResource bool) map[string]*schema.Schema {
	encs := []string{gp.EncryptionAes128Cbc, gp.EncryptionAes128Gcm, gp.EncryptionAes256Gcm}
	auths := []string{gp.AuthenticationSha1}

	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"name": {
			Type:        schema.TypeString,
			Description: "The name.",
			Required:    true,
		},
		"encryptions": {
			Type: schema.TypeList,
			Description: addStringInSliceValidation(
				"The encryptions.",
				encs,
			),
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringIn(encs...),
			},
		},
		"authentications": {
			Type: schema.TypeList,
			Description: addStringInSliceValidation(
				"The authentication algorithms.",
				auths,
			),
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringIn(auths...),
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "name"})
	}

	return ans
}

func loadGlobalProtectIpsecCryptoProfile(d *schema.ResourceData) gp.Entry {
	return gp.Entry{
		Name:            d.Get("name").(string),
		Encryptions:     asStringList(d.Get("encryptions").([]interface{})),
		Authentications: asStringList(d.Get("authentications").([]interface{})),
	}
}

func saveGlobalProtectIpsecCryptoProfile(d *schema.ResourceData, o gp.Entry) {
	d.Set("name", o.Name)
	d.Set("encryptions", o.Encryptions)
	d.Set("authentications", o.Authentications)
}

// Id functions.
func buildGlobalProtectIpsecCryptoProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseGlobalProtectIpsecCryptoProfileId(v string) (string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 3 {
		return "", "", "", fmt.Errorf("Expected 3 tokens, got %d", len(t))
	}

	return t[0], t[1], t[2], nil
}
