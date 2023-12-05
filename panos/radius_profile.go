package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/radius"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceRadiusProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceRadiusProfilesRead,

		Schema: s,
	}
}

func dataSourceRadiusProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildRadiusProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.RadiusProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.RadiusProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Resource.
func resourceRadiusProfile() *schema.Resource {
	return &schema.Resource{
		Create: createRadiusProfile,
		Read:   readRadiusProfile,
		Update: updateRadiusProfile,
		Delete: deleteRadiusProfile,

		Schema: radiusProfileSchema(),
	}
}

func createRadiusProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo radius.Entry
	o := loadRadiusProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildRadiusProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.RadiusProfile.Set(vsys, o)
		if err == nil && len(o.Servers) > 0 {
			lo, err = con.Device.RadiusProfile.Get(vsys, o.Name)
		}
	case *pango.Panorama:
		err = con.Device.RadiusProfile.Set(tmpl, ts, vsys, o)
		if err == nil && len(o.Servers) > 0 {
			lo, err = con.Device.RadiusProfile.Get(tmpl, ts, vsys, o.Name)
		}
	}

	if err != nil {
		return err
	}

	secretsEnc := map[string]interface{}{}
	secretsRaw := map[string]interface{}{}
	if len(o.Servers) != len(lo.Servers) {
		return fmt.Errorf("Servers in config is len:%d, but on live it is len:%d", len(o.Servers), len(lo.Servers))
	}
	for i := range o.Servers {
		cs, ls := o.Servers[i], lo.Servers[i]
		if cs.Name != ls.Name {
			return fmt.Errorf("index %d server mismatch: config:%q live:%q", i, cs.Name, ls.Name)
		}
		secretsEnc[ls.Name] = ls.Secret
		secretsRaw[cs.Name] = cs.Secret
	}

	d.SetId(id)
	d.Set("secrets_enc", secretsEnc)
	d.Set("secrets_raw", secretsRaw)

	return readRadiusProfile(d, meta)
}

func readRadiusProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o radius.Entry

	tmpl, ts, vsys, name, err := parseRadiusProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.RadiusProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.RadiusProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveRadiusProfile(d, o)
	return nil
}

func updateRadiusProfile(d *schema.ResourceData, meta interface{}) error {
	var lo radius.Entry
	o := loadRadiusProfile(d)

	tmpl, ts, vsys, _, err := parseRadiusProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.RadiusProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.RadiusProfile.Edit(vsys, lo); err == nil && len(lo.Servers) > 0 {
				lo, err = con.Device.RadiusProfile.Get(vsys, lo.Name)
			}
		}
	case *pango.Panorama:
		if lo, err = con.Device.RadiusProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.RadiusProfile.Edit(tmpl, ts, vsys, lo); err == nil && len(lo.Servers) > 0 {
				lo, err = con.Device.RadiusProfile.Get(tmpl, ts, vsys, lo.Name)
			}
		}
	}

	if err != nil {
		return err
	}

	secretsEnc := map[string]interface{}{}
	secretsRaw := map[string]interface{}{}
	if len(o.Servers) != len(lo.Servers) {
		return fmt.Errorf("Servers in config is len:%d, but on live it is len:%d", len(o.Servers), len(lo.Servers))
	}
	for i := range o.Servers {
		cs, ls := o.Servers[i], lo.Servers[i]
		if cs.Name != ls.Name {
			return fmt.Errorf("index %d server mismatch: config:%q live:%q", i, cs.Name, ls.Name)
		}
		secretsEnc[ls.Name] = ls.Secret
		secretsRaw[cs.Name] = cs.Secret
	}
	d.Set("secrets_enc", secretsEnc)
	d.Set("secrets_raw", secretsRaw)

	return readRadiusProfile(d, meta)
}

func deleteRadiusProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseRadiusProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.RadiusProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.RadiusProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func radiusProfileSchema() map[string]*schema.Schema {
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
		"timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Timeout (in sec).",
			Default:     3,
		},
		"retries": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of attempts before giving up authentication.",
			Default:     3,
		},
		"server": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Radius servers.",
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
					"secret": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Shared secret for Radius communication.",
						Sensitive:   true,
					},
					"port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Radius server port number.",
						Default:     1812,
					},
				},
			},
		},
		"protocol": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "(PAN-OS 8.0+) Authentication protocol settings.",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"chap": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "CHAP.",
					},
					"pap": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "PAP.",
					},
					"auto": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "(PAN-OS 8.0 only) Auto.",
					},
					"peap_mschap_v2": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "PEAP-MSCHAPv2.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"make_outer_identity_anonymous": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Make outer identity anonymous.",
								},
								"allow_expired_password_change": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Allow users to change passwords after expiry.",
								},
								"certificate_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Certificate profile for verifying the Radius server.",
								},
							},
						},
					},
					"peap_with_gtc": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "PEAP with GTC.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"make_outer_identity_anonymous": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Make outer identity anonymous.",
								},
								"certificate_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Certificate profile for verifying the Radius server.",
								},
							},
						},
					},
					"eap_ttls_with_pap": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "EAP-TTLS with PAP.",
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"make_outer_identity_anonymous": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Make outer identity anonymous.",
								},
								"certificate_profile": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Certificate profile for verifying the Radius server.",
								},
							},
						},
					},
				},
			},
		},
		"secrets_raw": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Server secrets, raw.",
			Sensitive:   true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"secrets_enc": {
			Type:        schema.TypeMap,
			Computed:    true,
			Description: "Server secrets, encrypted.",
			Sensitive:   true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func loadRadiusProfile(d *schema.ResourceData) radius.Entry {
	var proto radius.Protocol
	if p := configFolder(d, "protocol"); len(p) > 0 {
		if p["chap"].(bool) {
			proto.Chap = true
		} else if p["pap"].(bool) {
			proto.Pap = true
		} else if p["auto"].(bool) {
			proto.Auto = true
		} else if conf := asInterfaceMap(p, "peap_mschap_v2"); len(conf) > 0 {
			proto.PeapMschapv2 = &radius.PeapMschapv2{
				MakeOuterIdentityAnonymous: conf["make_outer_identity_anonymous"].(bool),
				AllowExpiredPasswordChange: conf["allow_expired_password_change"].(bool),
				CertificateProfile:         conf["certificate_profile"].(string),
			}
		} else if conf := asInterfaceMap(p, "peap_with_gtc"); len(conf) > 0 {
			proto.PeapWithGtc = &radius.PeapWithGtc{
				MakeOuterIdentityAnonymous: conf["make_outer_identity_anonymous"].(bool),
				CertificateProfile:         conf["certificate_profile"].(string),
			}
		} else if conf := asInterfaceMap(p, "eap_ttls_with_pap"); len(conf) > 0 {
			proto.EapTtlsWithPap = &radius.EapTtlsWithPap{
				MakeOuterIdentityAnonymous: conf["make_outer_identity_anonymous"].(bool),
				CertificateProfile:         conf["certificate_profile"].(string),
			}
		}
	}

	var listing []radius.Server
	sl := d.Get("server").([]interface{})
	if len(sl) > 0 {
		listing = make([]radius.Server, 0, len(sl))
		for i := range sl {
			x := sl[i].(map[string]interface{})
			listing = append(listing, radius.Server{
				Name:   x["name"].(string),
				Server: x["server"].(string),
				Secret: x["secret"].(string),
				Port:   x["port"].(int),
			})
		}
	}

	return radius.Entry{
		Name:         d.Get("name").(string),
		AdminUseOnly: d.Get("admin_use_only").(bool),
		Timeout:      d.Get("timeout").(int),
		Retries:      d.Get("retries").(int),
		Servers:      listing,
		Protocol:     proto,
	}
}

func saveRadiusProfile(d *schema.ResourceData, o radius.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("admin_use_only", o.AdminUseOnly)
	d.Set("timeout", o.Timeout)
	d.Set("retries", o.Retries)

	if len(o.Servers) == 0 {
		d.Set("server", nil)
	} else {
		secretsRaw := d.Get("secrets_raw").(map[string]interface{})
		secretsEnc := d.Get("secrets_enc").(map[string]interface{})
		listing := make([]interface{}, 0, len(o.Servers))

		for _, x := range o.Servers {
			var pwd string
			if secretsEnc[x.Name] != nil && secretsEnc[x.Name].(string) == x.Secret {
				pwd = secretsRaw[x.Name].(string)
			}

			listing = append(listing, map[string]interface{}{
				"name":   x.Name,
				"server": x.Server,
				"secret": pwd,
				"port":   x.Port,
			})
		}

		if err = d.Set("server", listing); err != nil {
			log.Printf("[WARN] Error setting 'server' for %q: %s", d.Id(), err)
		}
	}

	proto := map[string]interface{}{}
	switch {
	case o.Protocol.Chap:
		proto["chap"] = true
	case o.Protocol.Pap:
		proto["pap"] = true
	case o.Protocol.Auto:
		proto["auto"] = true
	case o.Protocol.PeapMschapv2 != nil:
		proto["peap_mschap_v2"] = []map[string]interface{}{{
			"make_outer_identity_anonymous": o.Protocol.PeapMschapv2.MakeOuterIdentityAnonymous,
			"allow_expired_password_change": o.Protocol.PeapMschapv2.AllowExpiredPasswordChange,
			"certificate_profile":           o.Protocol.PeapMschapv2.CertificateProfile,
		}}
	case o.Protocol.PeapWithGtc != nil:
		proto["peap_with_gtc"] = []map[string]interface{}{{
			"make_outer_identity_anonymous": o.Protocol.PeapWithGtc.MakeOuterIdentityAnonymous,
			"certificate_profile":           o.Protocol.PeapWithGtc.CertificateProfile,
		}}
	case o.Protocol.EapTtlsWithPap != nil:
		proto["eap_ttls_with_pap"] = []map[string]interface{}{{
			"make_outer_identity_anonymous": o.Protocol.EapTtlsWithPap.MakeOuterIdentityAnonymous,
			"certificate_profile":           o.Protocol.EapTtlsWithPap.CertificateProfile,
		}}
	}
	if err = d.Set("protocol", []interface{}{proto}); err != nil {
		log.Printf("[WARN] Error setting 'protocol' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseRadiusProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildRadiusProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
