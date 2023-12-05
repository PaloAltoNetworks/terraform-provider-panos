package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/ldap"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceLdapProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceLdapProfilesRead,

		Schema: s,
	}
}

func dataSourceLdapProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildLdapProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.LdapProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.LdapProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Resource.
func resourceLdapProfile() *schema.Resource {
	return &schema.Resource{
		Create: createLdapProfile,
		Read:   readLdapProfile,
		Update: updateLdapProfile,
		Delete: deleteLdapProfile,

		Schema: ldapProfileSchema(),
	}
}

func createLdapProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo ldap.Entry
	o := loadLdapProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildLdapProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.LdapProfile.Set(vsys, o)
		if err == nil {
			lo, err = con.Device.LdapProfile.Get(vsys, o.Name)
		}
	case *pango.Panorama:
		err = con.Device.LdapProfile.Set(tmpl, ts, vsys, o)
		if err == nil {
			lo, err = con.Device.LdapProfile.Get(tmpl, ts, vsys, o.Name)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	d.Set("password_raw", o.Password)
	d.Set("password_enc", lo.Password)

	return readLdapProfile(d, meta)
}

func readLdapProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ldap.Entry

	tmpl, ts, vsys, name, err := parseLdapProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.LdapProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.LdapProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveLdapProfile(d, o)
	return nil
}

func updateLdapProfile(d *schema.ResourceData, meta interface{}) error {
	var lo ldap.Entry
	o := loadLdapProfile(d)

	tmpl, ts, vsys, _, err := parseLdapProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.LdapProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.LdapProfile.Edit(vsys, lo); err == nil {
				lo, err = con.Device.LdapProfile.Get(vsys, lo.Name)
			}
		}
	case *pango.Panorama:
		if lo, err = con.Device.LdapProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.LdapProfile.Edit(tmpl, ts, vsys, lo); err == nil {
				lo, err = con.Device.LdapProfile.Get(tmpl, ts, vsys, lo.Name)
			}
		}
	}

	if err != nil {
		return err
	}

	d.Set("password_raw", o.Password)
	d.Set("password_enc", lo.Password)

	return readLdapProfile(d, meta)
}

func deleteLdapProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseLdapProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.LdapProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.LdapProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func ldapProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"ldap_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "LDAP type.",
			Default:     ldap.LdapTypeOther,
			ValidateFunc: validateStringIn(
				ldap.LdapTypeActiveDirectory,
				ldap.LdapTypeEdirectory,
				ldap.LdapTypeSun,
				ldap.LdapTypeOther,
			),
		},
		"ssl": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "SSL.",
			Default:     true,
		},
		"verify_server_certificate": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Verify server certificate for SSL sessions.",
		},
		"disabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Disable this profile.",
		},
		"base_dn": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Default base distinguished name (DN) to use for searches.",
		},
		"bind_dn": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Bind distinguished name.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Bind password.",
			Sensitive:   true,
		},
		"search_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of seconds to wait for performing searches.",
			Default:     30,
		},
		"bind_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of seconds to use for connecting to servers.",
			Default:     30,
		},
		"retry_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Interval (in seconds) for reconnecting LDAP server.",
			Default:     60,
		},
		"server": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of LDAP servers.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The server name.",
					},
					"server": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Server hostname or IP address.",
					},
					"port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "LDAP server port number.",
						Default:     389,
					},
				},
			},
		},
		"password_raw": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Password, raw.",
			Sensitive:   true,
		},
		"password_enc": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Password, encrypted.",
			Sensitive:   true,
		},
	}
}

func loadLdapProfile(d *schema.ResourceData) ldap.Entry {
	var listing []ldap.Server
	sl := d.Get("server").([]interface{})
	if len(sl) > 0 {
		listing = make([]ldap.Server, 0, len(sl))
		for i := range sl {
			x := sl[i].(map[string]interface{})
			listing = append(listing, ldap.Server{
				Name:   x["name"].(string),
				Server: x["server"].(string),
				Port:   x["port"].(int),
			})
		}
	}

	return ldap.Entry{
		Name:                    d.Get("name").(string),
		AdminUseOnly:            d.Get("admin_use_only").(bool),
		LdapType:                d.Get("ldap_type").(string),
		Ssl:                     d.Get("ssl").(bool),
		VerifyServerCertificate: d.Get("verify_server_certificate").(bool),
		Disabled:                d.Get("disabled").(bool),
		BaseDn:                  d.Get("base_dn").(string),
		BindDn:                  d.Get("bind_dn").(string),
		Password:                d.Get("password").(string),
		BindTimeout:             d.Get("bind_timeout").(int),
		SearchTimeout:           d.Get("search_timeout").(int),
		RetryInterval:           d.Get("retry_interval").(int),
		Servers:                 listing,
	}
}

func saveLdapProfile(d *schema.ResourceData, o ldap.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("admin_use_only", o.AdminUseOnly)
	d.Set("ldap_type", o.LdapType)
	d.Set("ssl", o.Ssl)
	d.Set("verify_server_certificate", o.VerifyServerCertificate)
	d.Set("disabled", o.Disabled)
	d.Set("base_dn", o.BaseDn)
	d.Set("bind_dn", o.BindDn)
	d.Set("bind_timeout", o.BindTimeout)
	d.Set("search_timeout", o.SearchTimeout)
	d.Set("retry_interval", o.RetryInterval)

	var pwd string
	if d.Get("password_enc").(string) == o.Password {
		pwd = d.Get("password_raw").(string)
	} else {
		pwd = "(mismatch)"
	}
	d.Set("password", pwd)

	if len(o.Servers) == 0 {
		d.Set("server", nil)
	} else {
		listing := make([]interface{}, 0, len(o.Servers))
		for _, x := range o.Servers {
			listing = append(listing, map[string]interface{}{
				"name":   x.Name,
				"server": x.Server,
				"port":   x.Port,
			})
		}

		if err = d.Set("server", listing); err != nil {
			log.Printf("[WARN] Error setting 'server' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func parseLdapProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildLdapProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
