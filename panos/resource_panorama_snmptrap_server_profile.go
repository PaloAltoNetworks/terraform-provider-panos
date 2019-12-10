package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v2c"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaSnmptrapServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaSnmptrapServerProfile,
		Read:   readPanoramaSnmptrapServerProfile,
		Update: createUpdatePanoramaSnmptrapServerProfile,
		Delete: deletePanoramaSnmptrapServerProfile,

		Schema: snmptrapServerProfileSchema(true),
	}
}

func parsePanoramaSnmptrapServerProfile(d *schema.ResourceData) (string, string, string, string, snmp.Entry, []v2c.Entry, []v3.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	o, v2cList, v3List := loadSnmptrapServerProfile(d)

	return tmpl, ts, vsys, dg, o, v2cList, v3List
}

func parsePanoramaSnmptrapServerProfileId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaSnmptrapServerProfileId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createUpdatePanoramaSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, o, v2cList, v3List := parsePanoramaSnmptrapServerProfile(d)

	if v2cList == nil && v3List == nil {
		return fmt.Errorf("Must specify at least one of v2c_server or v3_server")
	}

	if err := pano.Device.SnmpServerProfile.SetWithoutSubconfig(tmpl, ts, vsys, dg, o); err != nil {
		return err
	}

	if len(v2cList) > 0 {
		if err := pano.Device.SnmpV2cServer.Set(tmpl, ts, vsys, dg, o.Name, v2cList...); err != nil {
			return err
		}
		d.Set("auth_password_enc", nil)
		d.Set("auth_password_raw", nil)
		d.Set("priv_password_enc", nil)
		d.Set("priv_password_raw", nil)
	} else {
		if err := pano.Device.SnmpV3Server.Set(tmpl, ts, vsys, dg, o.Name, v3List...); err != nil {
			return err
		}

		auth_password_enc := make(map[string]interface{})
		auth_password_raw := make(map[string]interface{})
		priv_password_enc := make(map[string]interface{})
		priv_password_raw := make(map[string]interface{})

		for _, x := range v3List {
			live, err := pano.Device.SnmpV3Server.Get(tmpl, ts, vsys, dg, o.Name, x.Name)
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

	d.SetId(buildPanoramaSnmptrapServerProfileId(tmpl, ts, vsys, dg, o.Name))
	return readPanoramaSnmptrapServerProfile(d, meta)
}

func readPanoramaSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaSnmptrapServerProfileId(d.Id())

	o, err := pano.Device.SnmpServerProfile.Get(tmpl, ts, vsys, dg, name)
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
		list, err := pano.Device.SnmpV2cServer.GetList(tmpl, ts, vsys, dg, name)
		if err != nil {
			return err
		}
		v2cList = make([]v2c.Entry, 0, len(list))
		for i := range list {
			entry, err := pano.Device.SnmpV2cServer.Get(tmpl, ts, vsys, dg, name, list[i])
			if err != nil {
				return err
			}
			v2cList = append(v2cList, entry)
		}
	case snmp.SnmpVersionV3:
		list, err := pano.Device.SnmpV3Server.GetList(tmpl, ts, vsys, dg, name)
		if err != nil {
			return err
		}
		v3List = make([]v3.Entry, 0, len(list))
		for i := range list {
			entry, err := pano.Device.SnmpV3Server.Get(tmpl, ts, vsys, dg, name, list[i])
			if err != nil {
				return err
			}
			v3List = append(v3List, entry)
		}
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	saveSnmptrapServerProfile(d, o, v2cList, v3List)

	return nil
}

func deletePanoramaSnmptrapServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaSnmptrapServerProfileId(d.Id())

	err := pano.Device.SnmpServerProfile.Delete(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
