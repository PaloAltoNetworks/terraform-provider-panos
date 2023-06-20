package panos

import (
	"fmt"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/dev/profile/saml"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceSamlProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceSamlProfilesRead,

		Schema: s,
	}
}

func dataSourceSamlProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildSamlProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.SamlProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.SamlProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceSamlProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSamlProfileRead,

		Schema: samlProfileSchema(false),
	}
}

func dataSourceSamlProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o saml.Entry

	tmpl, ts, vsys, name := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string), d.Get("name").(string)

	id := buildSamlProfileId(tmpl, ts, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SamlProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.SamlProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveSamlProfile(d, o)
	return nil
}

// Resource.
func resourceSamlProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSamlProfile,
		Read:   readSamlProfile,
		Update: updateSamlProfile,
		Delete: deleteSamlProfile,

		Schema: samlProfileSchema(true),
	}
}

func createSamlProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadSamlProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildSamlProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SamlProfile.Set(vsys, o)
	case *pango.Panorama:
		err = con.Device.SamlProfile.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	return readSamlProfile(d, meta)
}

func readSamlProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o saml.Entry

	tmpl, ts, vsys, name, err := parseSamlProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SamlProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.SamlProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSamlProfile(d, o)
	return nil
}

func updateSamlProfile(d *schema.ResourceData, meta interface{}) error {
	var lo saml.Entry
	o := loadSamlProfile(d)

	tmpl, ts, vsys, _, err := parseSamlProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.SamlProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			err = con.Device.SamlProfile.Edit(vsys, lo)
		}
	case *pango.Panorama:
		if lo, err = con.Device.SamlProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			err = con.Device.SamlProfile.Edit(tmpl, ts, vsys, lo)
		}
	}

	if err != nil {
		return err
	}

	return readSamlProfile(d, meta)
}

func deleteSamlProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseSamlProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SamlProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.SamlProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func samlProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name.",
		},
		"admin_use_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Administrator use only.",
		},
		"identity_provider_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Unique identifier for SAML IdP.",
		},
		"identity_provider_certificate": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Object name of IdP signing certificate.",
		},
		"sso_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The single sign on service URL for the IdP server.",
		},
		"sso_binding": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "SAML HTTP binding for SSO requests to IdP.",
			Default:      saml.BindingPost,
			ValidateFunc: validateStringIn(saml.BindingPost, saml.BindingRedirect),
		},
		"slo_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The single logout service URL for the IdP server.",
		},
		"slo_binding": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "SAML HTTP binding for SLO requests to IdP.",
			Default:      saml.BindingPost,
			ValidateFunc: validateStringIn(saml.BindingPost, saml.BindingRedirect),
		},
		"validate_identity_provider_certificate": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Validate identity provider certificate.",
			Default:     true,
		},
		"sign_saml_message": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Sign SAML message to IdP.",
			Default:     true,
		},
		"max_clock_skew": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum allowed clock skew in seconds between SAML entities.",
			Default:     60,
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func loadSamlProfile(d *schema.ResourceData) saml.Entry {
	return saml.Entry{
		Name:                                d.Get("name").(string),
		AdminUseOnly:                        d.Get("admin_use_only").(bool),
		IdentityProviderId:                  d.Get("identity_provider_id").(string),
		IdentityProviderCertificate:         d.Get("identity_provider_certificate").(string),
		SsoUrl:                              d.Get("sso_url").(string),
		SsoBinding:                          d.Get("sso_binding").(string),
		SloUrl:                              d.Get("slo_url").(string),
		SloBinding:                          d.Get("slo_binding").(string),
		ValidateIdentityProviderCertificate: d.Get("validate_identity_provider_certificate").(bool),
		SignSamlMessage:                     d.Get("sign_saml_message").(bool),
		MaxClockSkew:                        d.Get("max_clock_skew").(int),
	}
}

func saveSamlProfile(d *schema.ResourceData, o saml.Entry) {
	d.Set("name", o.Name)
	d.Set("admin_use_only", o.AdminUseOnly)
	d.Set("identity_provider_id", o.IdentityProviderId)
	d.Set("identity_provider_certificate", o.IdentityProviderCertificate)
	d.Set("sso_url", o.SsoUrl)
	d.Set("sso_binding", o.SsoBinding)
	d.Set("slo_url", o.SloUrl)
	d.Set("slo_binding", o.SloBinding)
	d.Set("validate_identity_provider_certificate", o.ValidateIdentityProviderCertificate)
	d.Set("sign_saml_message", o.SignSamlMessage)
	d.Set("max_clock_skew", o.MaxClockSkew)
}

// Id functions.
func parseSamlProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildSamlProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
