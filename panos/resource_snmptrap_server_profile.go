package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v2c"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v3"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSnmptrapServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSnmptrapServerProfile,
		Read:   readSnmptrapServerProfile,
		Update: createUpdateSnmptrapServerProfile,
		Delete: deleteSnmptrapServerProfile,

		Schema: snmptrapServerProfileSchema(false),
	}
}

func snmptrapServerProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
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

	if p {
		ans["template"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template_stack",
				"device_group",
			},
		}
		ans["template_stack"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template",
				"device_group",
			},
		}
		ans["device_group"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template",
				"template_stack",
				"vsys",
			},
		}
		ans["vsys"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"device_group",
			},
		}
	} else {
		ans["vsys"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "shared",
			ForceNew: true,
		}
	}

	return ans
}

func parseSnmptrapServerProfile(d *schema.ResourceData) (string, snmp.Entry, []v2c.Entry, []v3.Entry) {
	vsys := d.Get("vsys").(string)
	o, v2cList, v3List := loadSnmptrapServerProfile(d)

	return vsys, o, v2cList, v3List
}

func loadSnmptrapServerProfile(d *schema.ResourceData) (snmp.Entry, []v2c.Entry, []v3.Entry) {
	o := snmp.Entry{
		Name: d.Get("name").(string),
	}

	if v2l := d.Get("v2c_server").([]interface{}); len(v2l) != 0 {
		o.SnmpVersion = snmp.SnmpVersionV2c

		v2cList := make([]v2c.Entry, 0, len(v2l))
		for i := range v2l {
			x := v2l[i].(map[string]interface{})
			v2cList = append(v2cList, v2c.Entry{
				Name:      x["name"].(string),
				Manager:   x["manager"].(string),
				Community: x["community"].(string),
			})
		}
		return o, v2cList, nil
	} else if v3l := d.Get("v3_server").([]interface{}); len(v3l) != 0 {
		o.SnmpVersion = snmp.SnmpVersionV3

		v3List := make([]v3.Entry, 0, len(v3l))
		for i := range v3l {
			x := v3l[i].(map[string]interface{})
			v3List = append(v3List, v3.Entry{
				Name:         x["name"].(string),
				Manager:      x["manager"].(string),
				User:         x["user"].(string),
				EngineId:     x["engine_id"].(string),
				AuthPassword: x["auth_password"].(string),
				PrivPassword: x["priv_password"].(string),
			})
		}
		return o, nil, v3List
	}

	return o, nil, nil
}

func saveSnmptrapServerProfile(d *schema.ResourceData, o snmp.Entry, v2cList []v2c.Entry, v3List []v3.Entry) {
	d.Set("name", o.Name)

	if len(v2cList) > 0 {
		v2c_server := make([]interface{}, 0, len(v2cList))
		for i := range v2cList {
			v2c_server = append(v2c_server, map[string]interface{}{
				"name":      v2cList[i].Name,
				"manager":   v2cList[i].Manager,
				"community": v2cList[i].Community,
			})
		}

		if err := d.Set("v2c_server", v2c_server); err != nil {
			log.Printf("[WARN] Error setting 'v2c_server' for %q: %s", d.Id(), err)
		}
		d.Set("v3_server", nil)
	} else if len(v3List) > 0 {
		apMap := d.Get("auth_password_enc").(map[string]interface{})
		apMapRaw := d.Get("auth_password_raw").(map[string]interface{})
		ppMap := d.Get("priv_password_enc").(map[string]interface{})
		ppMapRaw := d.Get("priv_password_raw").(map[string]interface{})

		v3_server := make([]interface{}, 0, len(v3List))
		for i := range v3List {
			var authPassword, privPassword string
			if apMap[v3List[i].Name] != nil && (apMap[v3List[i].Name]).(string) == v3List[i].AuthPassword {
				authPassword = apMapRaw[v3List[i].Name].(string)
			}
			if ppMap[v3List[i].Name] != nil && (ppMap[v3List[i].Name]).(string) == v3List[i].PrivPassword {
				privPassword = ppMapRaw[v3List[i].Name].(string)
			}
			entry := map[string]interface{}{
				"name":          v3List[i].Name,
				"manager":       v3List[i].Manager,
				"user":          v3List[i].User,
				"engine_id":     v3List[i].EngineId,
				"auth_password": authPassword,
				"priv_password": privPassword,
			}
			if apMap[v3List[i].Name] == nil || (apMap[v3List[i].Name]).(string) != v3List[i].AuthPassword {
				entry["auth_password"] = "(incorrect password)"
			}
			if ppMap[v3List[i].Name] == nil || (ppMap[v3List[i].Name]).(string) != v3List[i].PrivPassword {
				entry["priv_password"] = "(incorrect password)"
			}
			v3_server = append(v3_server, entry)
		}

		d.Set("v2c_server", nil)
		if err := d.Set("v3_server", v3_server); err != nil {
			log.Printf("[WARN] Error setting 'v3_server' for %q: %s", d.Id(), err)
		}
	}
}

func parseSnmptrapServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildSnmptrapServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o, v2cList, v3List := parseSnmptrapServerProfile(d)

	if v2cList == nil && v3List == nil {
		return fmt.Errorf("Must specify at least one of v2c_server or v3_server")
	}

	if err := fw.Device.SnmpServerProfile.SetWithoutSubconfig(vsys, o); err != nil {
		return err
	}

	if len(v2cList) > 0 {
		if err := fw.Device.SnmpV2cServer.Set(vsys, o.Name, v2cList...); err != nil {
			return err
		}
		d.Set("auth_password_enc", nil)
		d.Set("auth_password_raw", nil)
		d.Set("priv_password_enc", nil)
		d.Set("priv_password_raw", nil)
	} else {
		if err := fw.Device.SnmpV3Server.Set(vsys, o.Name, v3List...); err != nil {
			return err
		}

		auth_password_enc := make(map[string]interface{})
		auth_password_raw := make(map[string]interface{})
		priv_password_enc := make(map[string]interface{})
		priv_password_raw := make(map[string]interface{})

		for _, x := range v3List {
			live, err := fw.Device.SnmpV3Server.Get(vsys, o.Name, x.Name)
			if err != nil {
				return err
			}
			auth_password_enc[x.Name] = live.AuthPassword
			auth_password_raw[x.Name] = x.AuthPassword
			priv_password_enc[x.Name] = live.PrivPassword
			priv_password_raw[x.Name] = x.PrivPassword
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
	}

	d.SetId(buildSnmptrapServerProfileId(vsys, o.Name))
	return readSnmptrapServerProfile(d, meta)
}

func readSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseSnmptrapServerProfileId(d.Id())

	o, err := fw.Device.SnmpServerProfile.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	var v2cList []v2c.Entry
	var v3List []v3.Entry
	switch o.SnmpVersion {
	case snmp.SnmpVersionV2c:
		list, err := fw.Device.SnmpV2cServer.GetList(vsys, name)
		if err != nil {
			return err
		}
		v2cList = make([]v2c.Entry, 0, len(list))
		for i := range list {
			entry, err := fw.Device.SnmpV2cServer.Get(vsys, name, list[i])
			if err != nil {
				return err
			}
			v2cList = append(v2cList, entry)
		}
	case snmp.SnmpVersionV3:
		list, err := fw.Device.SnmpV3Server.GetList(vsys, name)
		if err != nil {
			return err
		}
		v3List = make([]v3.Entry, 0, len(list))
		for i := range list {
			entry, err := fw.Device.SnmpV3Server.Get(vsys, name, list[i])
			if err != nil {
				return err
			}
			v3List = append(v3List, entry)
		}
	}

	d.Set("vsys", vsys)
	saveSnmptrapServerProfile(d, o, v2cList, v3List)

	return nil
}

func deleteSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseSnmptrapServerProfileId(d.Id())

	err := fw.Device.SnmpServerProfile.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
