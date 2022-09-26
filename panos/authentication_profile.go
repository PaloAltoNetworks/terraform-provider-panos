package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	auth "github.com/PaloAltoNetworks/pango/dev/profile/authentication"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceAuthenticationProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceAuthenticationProfilesRead,

		Schema: s,
	}
}

func dataSourceAuthenticationProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildAuthenticationProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.AuthenticationProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.AuthenticationProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Resource.
func resourceAuthenticationProfile() *schema.Resource {
	return &schema.Resource{
		Create: createAuthenticationProfile,
		Read:   readAuthenticationProfile,
		Update: updateAuthenticationProfile,
		Delete: deleteAuthenticationProfile,

		Schema: authenticationProfileSchema(),
	}
}

func createAuthenticationProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var keytabRaw, keytabEnc string
	var lo auth.Entry
	o := loadAuthenticationProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildAuthenticationProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.AuthenticationProfile.Set(vsys, o)
		if err == nil && o.SingleSignOn != nil && o.SingleSignOn.Keytab != "" {
			keytabRaw = o.SingleSignOn.Keytab
			if lo, err = con.Device.AuthenticationProfile.Get(vsys, o.Name); err == nil && lo.SingleSignOn != nil {
				keytabEnc = lo.SingleSignOn.Keytab
			}
		}
	case *pango.Panorama:
		err = con.Device.AuthenticationProfile.Set(tmpl, ts, vsys, o)
		if err == nil && o.SingleSignOn != nil && o.SingleSignOn.Keytab != "" {
			keytabRaw = o.SingleSignOn.Keytab
			if lo, err = con.Device.AuthenticationProfile.Get(tmpl, ts, vsys, o.Name); err == nil && lo.SingleSignOn != nil {
				keytabEnc = lo.SingleSignOn.Keytab
			}
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("keytab_raw", keytabRaw)
	d.Set("keytab_enc", keytabEnc)

	return readAuthenticationProfile(d, meta)
}

func readAuthenticationProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o auth.Entry

	tmpl, ts, vsys, name, err := parseAuthenticationProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.AuthenticationProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.AuthenticationProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAuthenticationProfile(d, o)
	return nil
}

func updateAuthenticationProfile(d *schema.ResourceData, meta interface{}) error {
	var keytabRaw, keytabEnc string
	var lo auth.Entry
	o := loadAuthenticationProfile(d)

	tmpl, ts, vsys, _, err := parseAuthenticationProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.AuthenticationProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.AuthenticationProfile.Edit(vsys, lo); err == nil && lo.SingleSignOn != nil && lo.SingleSignOn.Keytab != "" {
				keytabRaw = lo.SingleSignOn.Keytab
				if lo, err = con.Device.AuthenticationProfile.Get(vsys, lo.Name); err == nil && lo.SingleSignOn != nil {
					keytabEnc = lo.SingleSignOn.Keytab
				}
			}
		}
	case *pango.Panorama:
		if lo, err = con.Device.AuthenticationProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.AuthenticationProfile.Edit(tmpl, ts, vsys, lo); err == nil && lo.SingleSignOn != nil && lo.SingleSignOn.Keytab != "" {
				keytabRaw = lo.SingleSignOn.Keytab
				if lo, err = con.Device.AuthenticationProfile.Get(tmpl, ts, vsys, lo.Name); err == nil && lo.SingleSignOn != nil {
					keytabEnc = lo.SingleSignOn.Keytab
				}
			}
		}
	}

	if err != nil {
		return err
	}

	d.Set("keytab_raw", keytabRaw)
	d.Set("keytab_enc", keytabEnc)
	return readAuthenticationProfile(d, meta)
}

func deleteAuthenticationProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseAuthenticationProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.AuthenticationProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.AuthenticationProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func authenticationProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name.",
		},
		"lockout_failed_attempts": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Number of failed attempts to trigger lock-out.",
		},
		"lockout_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of minutes to lock-out.",
		},
		"allow_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of allowed users or user groups.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"type": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The type specification.",
			MinItems:    1,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"none": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "No authentication.",
					},
					"local_database": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Local database authentication.",
					},
					"radius": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Radius authentication.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Radius server profile object.",
								},
								"retrieve_user_group": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "(PAN-OS 7.0+) Retrieve user group from RADIUS.",
								},
							},
						},
					},
					"ldap": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "LDAP authentication.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "LDAP server profile object.",
								},
								"login_attribute": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Login attribute in LDAP server to authenticate against.",
								},
								"password_expiry_warning": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Number of days prior to warning a user about password expiry.",
									Default:     "7",
								},
							},
						},
					},
					"kerberos": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "LDAP authentication.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Kerberos server profile object.",
								},
								"realm": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "(PAN-OS 7.0+) Realm name to be used for authentication.",
								},
							},
						},
					},
					"tacacs_plus": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "(PAN-OS 7.0+) TACACS+ authentication.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "TACACS+ server profile object.",
								},
								"retrieve_user_group": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "(PAN-OS 8.0+) Retrieve user group from TACACS+.",
								},
							},
						},
					},
					"saml": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "LDAP authentication.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"server_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "SAML IDP server profile object.",
								},
								"enable_single_logout": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Enable single logout.",
								},
								"request_signing_certificate": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Signing certificate for SAML requests.",
								},
								"certificate_profile": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Certificate profile for IDP and SP.",
								},
								"username_attribute": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Attribute name for username to be extracted from SAML response.",
									Default:     auth.UsernameAttributeDefault,
								},
								"user_group_attribute": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "User group attribute.",
								},
								"admin_role_attribute": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Admin role attribute.",
								},
								"access_domain_attribute": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Access domain attribute.",
								},
							},
						},
					},
				},
			},
		},
		"username_modifier": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 7.0+) Username modifier to handle user domain.",
			Default:     auth.UsernameModifierInput,
			ValidateFunc: validateStringIn(
				auth.UsernameModifierInput,
				auth.UsernameModifierInputDomain,
				auth.UsernameModifierDomainInput,
			),
		},
		"user_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 7.0+) Domain name to be used for authentication.",
		},
		"single_sign_on": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "(PAN-OS 7.0+) Kerberos SSO settings.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"realm": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Kerberos realm to be used for authentication.",
					},
					"service_principal": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Kerberos service principal.",
					},
					"keytab": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Kerberos keytab.",
						Sensitive:   true,
					},
				},
			},
		},
		"multi_factor_authentication": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "(PAN-OS 8.0+) Specify MFA configuration.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable additional authentication factors.",
					},
					"factors": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of additional authentication factors.",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},

		// Computed.
		"keytab_raw": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The raw Kerberos keytab value.",
			Sensitive:   true,
		},
		"keytab_enc": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The encrypted Kerberos keytab value.",
			Sensitive:   true,
		},
	}
}

func loadAuthenticationProfile(d *schema.ResourceData) auth.Entry {
	var sso *auth.SingleSignOn
	if x := d.Get("single_sign_on"); x != nil && len(x.([]interface{})) > 1 {
		ssoData := x.([]interface{})[0].(map[string]interface{})
		sso = &auth.SingleSignOn{
			Realm:            ssoData["realm"].(string),
			ServicePrincipal: ssoData["service_principal"].(string),
			Keytab:           ssoData["keytab"].(string),
		}
	}

	var authType auth.AuthenticationType
	authData := d.Get("type").([]interface{})[0].(map[string]interface{})
	if authData["none"].(bool) {
		authType = auth.AuthenticationType{None: true}
	} else if authData["local_database"].(bool) {
		authType = auth.AuthenticationType{LocalDatabase: true}
	} else if x := asInterfaceMap(authData, "radius"); len(x) > 0 {
		authType = auth.AuthenticationType{Radius: &auth.Radius{
			ServerProfile:     x["server_profile"].(string),
			RetrieveUserGroup: x["retrieve_user_group"].(bool),
		}}
	} else if x := asInterfaceMap(authData, "ldap"); len(x) > 0 {
		authType = auth.AuthenticationType{Ldap: &auth.Ldap{
			ServerProfile:         x["server_profile"].(string),
			LoginAttribute:        x["login_attribute"].(string),
			PasswordExpiryWarning: x["password_expiry_warning"].(string),
		}}
	} else if x := asInterfaceMap(authData, "kerberos"); len(x) > 0 {
		authType = auth.AuthenticationType{Kerberos: &auth.Kerberos{
			ServerProfile: x["server_profile"].(string),
			Realm:         x["realm"].(string),
		}}
	} else if x := asInterfaceMap(authData, "tacacs_plus"); len(x) > 0 {
		authType = auth.AuthenticationType{TacacsPlus: &auth.TacacsPlus{
			ServerProfile:     x["server_profile"].(string),
			RetrieveUserGroup: x["retrieve_user_group"].(bool),
		}}
	} else if x := asInterfaceMap(authData, "saml"); len(x) > 0 {
		authType = auth.AuthenticationType{Saml: &auth.Saml{
			ServerProfile:             x["server_profile"].(string),
			EnableSingleLogout:        x["enable_single_logout"].(bool),
			RequestSigningCertificate: x["request_signing_certificate"].(string),
			CertificateProfile:        x["certificate_profile"].(string),
			UsernameAttribute:         x["username_attribute"].(string),
			UserGroupAttribute:        x["user_group_attribute"].(string),
			AdminRoleAttribute:        x["admin_role_attribute"].(string),
			AccessDomainAttribute:     x["access_domain_attribute"].(string),
		}}
	}

	var mfa *auth.MultiFactorAuthentication
	if x := configFolder(d, "multi_factor_authentication"); len(x) > 0 {
		mfa = &auth.MultiFactorAuthentication{
			Enabled: x["enabled"].(bool),
			Factors: asStringList(x["factors"].([]interface{})),
		}
	}

	return auth.Entry{
		Name:                      d.Get("name").(string),
		LockoutFailedAttempts:     d.Get("lockout_failed_attempts").(string),
		LockoutTime:               d.Get("lockout_time").(int),
		AllowList:                 asStringList(d.Get("allow_list").([]interface{})),
		Type:                      authType,
		UsernameModifier:          d.Get("username_modifier").(string),
		UserDomain:                d.Get("user_domain").(string),
		SingleSignOn:              sso,
		MultiFactorAuthentication: mfa,
	}
}

func saveAuthenticationProfile(d *schema.ResourceData, o auth.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("lockout_failed_attempts", o.LockoutFailedAttempts)
	d.Set("lockout_time", o.LockoutTime)
	d.Set("allow_list", o.AllowList)
	d.Set("username_modifier", o.UsernameModifier)
	d.Set("user_domain", o.UserDomain)

	ts := map[string]interface{}{}
	switch {
	case o.Type.None:
		ts["none"] = true
	case o.Type.LocalDatabase:
		ts["local_database"] = true
	case o.Type.Radius != nil:
		ts["radius"] = []map[string]interface{}{{
			"server_profile":      o.Type.Radius.ServerProfile,
			"retrieve_user_group": o.Type.Radius.RetrieveUserGroup,
		}}
	case o.Type.Ldap != nil:
		ts["ldap"] = []map[string]interface{}{{
			"server_profile":          o.Type.Ldap.ServerProfile,
			"login_attribute":         o.Type.Ldap.LoginAttribute,
			"password_expiry_warning": o.Type.Ldap.PasswordExpiryWarning,
		}}
	case o.Type.Kerberos != nil:
		ts["kerberos"] = []map[string]interface{}{{
			"server_profile": o.Type.Kerberos.ServerProfile,
			"realm":          o.Type.Kerberos.Realm,
		}}
	case o.Type.TacacsPlus != nil:
		ts["tacacs_plus"] = []map[string]interface{}{{
			"server_profile":      o.Type.TacacsPlus.ServerProfile,
			"retrieve_user_group": o.Type.TacacsPlus.RetrieveUserGroup,
		}}
	case o.Type.Saml != nil:
		ts["saml"] = []map[string]interface{}{{
			"server_profile":              o.Type.Saml.ServerProfile,
			"enable_single_logout":        o.Type.Saml.EnableSingleLogout,
			"request_signing_certificate": o.Type.Saml.RequestSigningCertificate,
			"certificate_profile":         o.Type.Saml.CertificateProfile,
			"username_attribute":          o.Type.Saml.UsernameAttribute,
			"user_group_attribute":        o.Type.Saml.UserGroupAttribute,
			"admin_role_attribute":        o.Type.Saml.AdminRoleAttribute,
			"access_domain_attribute":     o.Type.Saml.AccessDomainAttribute,
		}}
	}
	if err = d.Set("type", []interface{}{ts}); err != nil {
		log.Printf("[WARN] Error setting 'type' for %q: %s", d.Id(), err)
	}

	if o.SingleSignOn == nil {
		d.Set("single_sign_on", nil)
	} else {
		keytabRaw, keytabEnc := d.Get("keytab_raw").(string), d.Get("keytab_enc").(string)
		var keytabVal string

		if o.SingleSignOn.Keytab == keytabEnc {
			keytabVal = keytabRaw
		} else {
			keytabVal = "(mismatch)"
		}

		spec := map[string]interface{}{
			"realm":             o.SingleSignOn.Realm,
			"service_principal": o.SingleSignOn.ServicePrincipal,
			"keytab":            keytabVal,
		}

		if err = d.Set("single_sign_on", spec); err != nil {
			log.Printf("[WARN] Error setting 'single_sign_on' for %q: %s", d.Id(), err)
		}
	}

	var mfa []interface{}
	if o.MultiFactorAuthentication != nil {
		mfa = []interface{}{map[string]interface{}{
			"enabled": o.MultiFactorAuthentication.Enabled,
			"factors": o.MultiFactorAuthentication.Factors,
		}}
	}
	if err = d.Set("multi_factor_authentication", mfa); err != nil {
		log.Printf("[WARN] Error setting 'multi_factor_authentication' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseAuthenticationProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildAuthenticationProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
