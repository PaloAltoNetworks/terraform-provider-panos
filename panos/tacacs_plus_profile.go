package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/tacplus"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceTacacsPlusProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceTacacsPlusProfilesRead,

		Schema: s,
	}
}

func dataSourceTacacsPlusProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildTacacsPlusProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.TacacsPlusProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.TacacsPlusProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Resource.
func resourceTacacsPlusProfile() *schema.Resource {
	return &schema.Resource{
		Create: createTacacsPlusProfile,
		Read:   readTacacsPlusProfile,
		Update: updateTacacsPlusProfile,
		Delete: deleteTacacsPlusProfile,

		Schema: tacacsPlusProfileSchema(),
	}
}

func createTacacsPlusProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo tacplus.Entry
	o := loadTacacsPlusProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildTacacsPlusProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.TacacsPlusProfile.Set(vsys, o)
		if err == nil && len(o.Servers) > 0 {
			lo, err = con.Device.TacacsPlusProfile.Get(vsys, o.Name)
		}
	case *pango.Panorama:
		err = con.Device.TacacsPlusProfile.Set(tmpl, ts, vsys, o)
		if err == nil && len(o.Servers) > 0 {
			lo, err = con.Device.TacacsPlusProfile.Get(tmpl, ts, vsys, o.Name)
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

	return readTacacsPlusProfile(d, meta)
}

func readTacacsPlusProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o tacplus.Entry

	tmpl, ts, vsys, name, err := parseTacacsPlusProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.TacacsPlusProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.TacacsPlusProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveTacacsPlusProfile(d, o)
	return nil
}

func updateTacacsPlusProfile(d *schema.ResourceData, meta interface{}) error {
	var lo tacplus.Entry
	o := loadTacacsPlusProfile(d)

	tmpl, ts, vsys, _, err := parseTacacsPlusProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.TacacsPlusProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.TacacsPlusProfile.Edit(vsys, lo); err == nil && len(lo.Servers) > 0 {
				lo, err = con.Device.TacacsPlusProfile.Get(vsys, lo.Name)
			}
		}
	case *pango.Panorama:
		if lo, err = con.Device.TacacsPlusProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			if err = con.Device.TacacsPlusProfile.Edit(tmpl, ts, vsys, lo); err == nil && len(lo.Servers) > 0 {
				lo, err = con.Device.TacacsPlusProfile.Get(tmpl, ts, vsys, lo.Name)
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

	return readTacacsPlusProfile(d, meta)
}

func deleteTacacsPlusProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseTacacsPlusProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.TacacsPlusProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.TacacsPlusProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func tacacsPlusProfileSchema() map[string]*schema.Schema {
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
		"use_single_connection": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Use single connection for all authentication.",
		},
		"server": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of TACACS+ servers.",
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
						Description: "Shared secret for TACACS+ communication.",
						Sensitive:   true,
					},
					"port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "TACACS+ server port number.",
						Default:     49,
					},
				},
			},
		},
		"protocol": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Authentication protocol settings.",
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

func loadTacacsPlusProfile(d *schema.ResourceData) tacplus.Entry {
	var proto tacplus.Protocol
	if p := configFolder(d, "protocol"); len(p) > 0 {
		if p["chap"].(bool) {
			proto.Chap = true
		} else if p["pap"].(bool) {
			proto.Pap = true
		} else if p["auto"].(bool) {
			proto.Auto = true
		}
	}

	var listing []tacplus.Server
	sl := d.Get("server").([]interface{})
	if len(sl) > 0 {
		listing = make([]tacplus.Server, 0, len(sl))
		for i := range sl {
			x := sl[i].(map[string]interface{})
			listing = append(listing, tacplus.Server{
				Name:   x["name"].(string),
				Server: x["server"].(string),
				Secret: x["secret"].(string),
				Port:   x["port"].(int),
			})
		}
	}

	return tacplus.Entry{
		Name:                d.Get("name").(string),
		AdminUseOnly:        d.Get("admin_use_only").(bool),
		Timeout:             d.Get("timeout").(int),
		UseSingleConnection: d.Get("use_single_connection").(bool),
		Servers:             listing,
		Protocol:            proto,
	}
}

func saveTacacsPlusProfile(d *schema.ResourceData, o tacplus.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("admin_use_only", o.AdminUseOnly)
	d.Set("timeout", o.Timeout)
	d.Set("use_single_connection", o.UseSingleConnection)

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
	}
	if err = d.Set("protocol", []interface{}{proto}); err != nil {
		log.Printf("[WARN] Error setting 'protocol' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseTacacsPlusProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildTacacsPlusProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
