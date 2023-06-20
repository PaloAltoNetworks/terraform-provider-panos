package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/dev/profile/snmp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Resource.
func resourceSnmptrapServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSnmptrapServerProfile,
		Read:   readSnmptrapServerProfile,
		Update: updateSnmptrapServerProfile,
		Delete: deleteSnmptrapServerProfile,

		Schema: snmptrapServerProfileSchema(true, "shared", []string{"device_group", "template", "template_stack"}),
	}
}

func resourcePanoramaSnmptrapServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSnmptrapServerProfile,
		Read:   readSnmptrapServerProfile,
		Update: updateSnmptrapServerProfile,
		Delete: deleteSnmptrapServerProfile,

		Schema: snmptrapServerProfileSchema(true, "", nil),
	}
}

func createSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadSnmptrapServerProfile(d)
	var lo snmp.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildSnmptrapServerProfileId(vsys, o.Name)
		if err = con.Device.SnmpServerProfile.Set(vsys, o); err != nil {
			return err
		}
		lo, err = con.Device.SnmpServerProfile.Get(vsys, o.Name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		id = buildPanoramaSnmptrapServerProfileId(tmpl, ts, vsys, o.Name)
		if err = con.Device.SnmpServerProfile.Set(tmpl, ts, vsys, o); err != nil {
			return err
		}
		lo, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, o.Name)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	if err = saveSnmptrapServerPasswords(d, o.V3Servers, lo.V3Servers); err != nil {
		return err
	}

	return readSnmptrapServerProfile(d, meta)
}

func readSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o snmp.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseSnmptrapServerProfileId(d.Id())
		d.Set("vsys", vsys)
		o, err = con.Device.SnmpServerProfile.Get(vsys, name)
	case *pango.Panorama:
		// If this is the old style ID, it will have an extra field.  So
		// we need to migrate the ID first.
		tok := strings.Split(d.Id(), IdSeparator)
		if len(tok) == 5 {
			d.SetId(buildPanoramaSnmptrapServerProfileId(tok[0], tok[1], tok[2], tok[4]))
		}
		// Continue on as normal.
		tmpl, ts, vsys, name := parsePanoramaSnmptrapServerProfileId(d.Id())
		d.Set("template", tmpl)
		d.Set("template_stack", ts)
		d.Set("vsys", vsys)
		d.Set("device_group", d.Get("device_group").(string))
		o, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSnmptrapServerProfile(d, o)

	return nil
}

func updateSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadSnmptrapServerProfile(d)
	var err error
	var lo snmp.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err = con.Device.SnmpServerProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.SnmpServerProfile.Edit(vsys, lo); err != nil {
			return err
		}
		lo, err = con.Device.SnmpServerProfile.Get(vsys, lo.Name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		lo, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.SnmpServerProfile.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
		lo, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, lo.Name)
	}

	if err != nil {
		return err
	}

	if err = saveSnmptrapServerPasswords(d, o.V3Servers, lo.V3Servers); err != nil {
		return err
	}

	return readSnmptrapServerProfile(d, meta)
}

func deleteSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseSnmptrapServerProfileId(d.Id())
		err = con.Device.SnmpServerProfile.Delete(vsys, name)
	case *pango.Panorama:
		tmpl, ts, vsys, name := parsePanoramaSnmptrapServerProfileId(d.Id())
		err = con.Device.SnmpServerProfile.Delete(tmpl, ts, vsys, name)
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
func snmptrapServerProfileSchema(isResource bool, vsysDefault string, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
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
		"device_group": {
			Type:       schema.TypeString,
			Optional:   true,
			Deprecated: "This parameter is not applicable to this resource and will be removed later.",
		},
		"vsys": vsysSchema(vsysDefault),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"v2c_server": {
			Type:          schema.TypeList,
			Optional:      true,
			MinItems:      1,
			ConflictsWith: []string{"v3_server"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"manager": {
						Type:     schema.TypeString,
						Required: true,
					},
					"community": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"v3_server": {
			Type:          schema.TypeList,
			Optional:      true,
			MinItems:      1,
			ConflictsWith: []string{"v2c_server"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"manager": {
						Type:     schema.TypeString,
						Required: true,
					},
					"user": {
						Type:     schema.TypeString,
						Required: true,
					},
					"engine_id": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"auth_password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"priv_password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},
		"auth_password_enc": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"auth_password_raw": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"priv_password_enc": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"priv_password_raw": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "device_group", "name"})
	}

	return ans
}

func loadSnmptrapServerProfile(d *schema.ResourceData) snmp.Entry {
	o := snmp.Entry{
		Name: d.Get("name").(string),
	}

	if v2l := d.Get("v2c_server").([]interface{}); len(v2l) != 0 {
		list := make([]snmp.V2cServer, 0, len(v2l))
		for i := range v2l {
			x := v2l[i].(map[string]interface{})
			list = append(list, snmp.V2cServer{
				Name:      x["name"].(string),
				Manager:   x["manager"].(string),
				Community: x["community"].(string),
			})
		}
		o.V2cServers = list
	}

	if v3l := d.Get("v3_server").([]interface{}); len(v3l) != 0 {
		list := make([]snmp.V3Server, 0, len(v3l))
		for i := range v3l {
			x := v3l[i].(map[string]interface{})
			list = append(list, snmp.V3Server{
				Name:         x["name"].(string),
				Manager:      x["manager"].(string),
				User:         x["user"].(string),
				EngineId:     x["engine_id"].(string),
				AuthPassword: x["auth_password"].(string),
				PrivPassword: x["priv_password"].(string),
			})
		}
		o.V3Servers = list
	}

	return o
}

func saveSnmptrapServerPasswords(d *schema.ResourceData, unencrypted, encrypted []snmp.V3Server) error {
	if len(unencrypted) == 0 {
		d.Set("auth_password_enc", nil)
		d.Set("auth_password_raw", nil)
		d.Set("priv_password_enc", nil)
		d.Set("priv_password_raw", nil)
		return nil
	}

	if len(unencrypted) != len(encrypted) {
		return fmt.Errorf("unencrypted len:%d vs encrypted len:%d", len(unencrypted), len(encrypted))
	}

	auth_password_enc := make(map[string]interface{})
	auth_password_raw := make(map[string]interface{})
	priv_password_enc := make(map[string]interface{})
	priv_password_raw := make(map[string]interface{})

	for i := range unencrypted {
		if unencrypted[i].Name != encrypted[i].Name {
			return fmt.Errorf("Name mismatch at index %d: config:%q vs list:%q", i, unencrypted[i].Name, encrypted[i].Name)
		}
		auth_password_raw[unencrypted[i].Name] = unencrypted[i].AuthPassword
		auth_password_enc[encrypted[i].Name] = encrypted[i].AuthPassword
		priv_password_raw[unencrypted[i].Name] = unencrypted[i].PrivPassword
		priv_password_enc[encrypted[i].Name] = encrypted[i].PrivPassword
	}

	if err := d.Set("auth_password_enc", auth_password_enc); err != nil {
		log.Printf("[WARN] Error setting 'auth_password_enc' for %q: %s", d.Id(), err)
	}
	if err := d.Set("auth_password_raw", auth_password_raw); err != nil {
		log.Printf("[WARN] Error setting 'auth_password_raw' for %q: %s", d.Id(), err)
	}
	if err := d.Set("priv_password_enc", priv_password_enc); err != nil {
		log.Printf("[WARN] Error setting 'priv_password_enc' for %q: %s", d.Id(), err)
	}
	if err := d.Set("priv_password_raw", priv_password_raw); err != nil {
		log.Printf("[WARN] Error setting 'priv_password_raw' for %q: %s", d.Id(), err)
	}

	return nil
}

func saveSnmptrapServerProfile(d *schema.ResourceData, o snmp.Entry) {
	d.Set("name", o.Name)

	if len(o.V2cServers) == 0 {
		d.Set("v2c_server", nil)
	} else {
		list := make([]interface{}, 0, len(o.V2cServers))
		for _, x := range o.V2cServers {
			list = append(list, map[string]interface{}{
				"name":      x.Name,
				"manager":   x.Manager,
				"community": x.Community,
			})
		}

		if err := d.Set("v2c_server", list); err != nil {
			log.Printf("[WARN] Error setting 'v2c_server' for %q: %s", d.Id(), err)
		}
	}

	if len(o.V3Servers) == 0 {
		d.Set("v3_server", nil)
	} else {
		apMap := d.Get("auth_password_enc").(map[string]interface{})
		apMapRaw := d.Get("auth_password_raw").(map[string]interface{})
		ppMap := d.Get("priv_password_enc").(map[string]interface{})
		ppMapRaw := d.Get("priv_password_raw").(map[string]interface{})

		list := make([]interface{}, 0, len(o.V3Servers))
		for _, x := range o.V3Servers {
			var authPassword, privPassword string
			if apMap[x.Name] != nil && (apMap[x.Name]).(string) == x.AuthPassword {
				authPassword = apMapRaw[x.Name].(string)
			}
			if ppMap[x.Name] != nil && (ppMap[x.Name]).(string) == x.PrivPassword {
				privPassword = ppMapRaw[x.Name].(string)
			}
			entry := map[string]interface{}{
				"name":          x.Name,
				"manager":       x.Manager,
				"user":          x.User,
				"engine_id":     x.EngineId,
				"auth_password": authPassword,
				"priv_password": privPassword,
			}
			if apMap[x.Name] == nil || (apMap[x.Name]).(string) != x.AuthPassword {
				entry["auth_password"] = "(incorrect password)"
			}
			if ppMap[x.Name] == nil || (ppMap[x.Name]).(string) != x.PrivPassword {
				entry["priv_password"] = "(incorrect password)"
			}
			list = append(list, entry)
		}

		if err := d.Set("v3_server", list); err != nil {
			log.Printf("[WARN] Error setting 'v3_server' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func parseSnmptrapServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaSnmptrapServerProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildSnmptrapServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaSnmptrapServerProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
