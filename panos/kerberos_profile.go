package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/kerberos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceKerberosProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceKerberosProfilesRead,

		Schema: s,
	}
}

func dataSourceKerberosProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildKerberosProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.KerberosProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.KerberosProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceKerberosProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKerberosProfileRead,

		Schema: kerberosProfileSchema(false),
	}
}

func dataSourceKerberosProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o kerberos.Entry

	tmpl, ts, vsys, name := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string), d.Get("name").(string)

	id := buildKerberosProfileId(tmpl, ts, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.KerberosProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.KerberosProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveKerberosProfile(d, o)
	return nil
}

// Resource.
func resourceKerberosProfile() *schema.Resource {
	return &schema.Resource{
		Create: createKerberosProfile,
		Read:   readKerberosProfile,
		Update: updateKerberosProfile,
		Delete: deleteKerberosProfile,

		Schema: kerberosProfileSchema(true),
	}
}

func createKerberosProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadKerberosProfile(d)
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)

	id := buildKerberosProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.KerberosProfile.Set(vsys, o)
	case *pango.Panorama:
		err = con.Device.KerberosProfile.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)

	return readKerberosProfile(d, meta)
}

func readKerberosProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o kerberos.Entry

	tmpl, ts, vsys, name, err := parseKerberosProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.KerberosProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.KerberosProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveKerberosProfile(d, o)
	return nil
}

func updateKerberosProfile(d *schema.ResourceData, meta interface{}) error {
	var lo kerberos.Entry
	o := loadKerberosProfile(d)

	tmpl, ts, vsys, _, err := parseKerberosProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		if lo, err = con.Device.KerberosProfile.Get(vsys, o.Name); err == nil {
			lo.Copy(o)
			err = con.Device.KerberosProfile.Edit(vsys, lo)
		}
	case *pango.Panorama:
		if lo, err = con.Device.KerberosProfile.Get(tmpl, ts, vsys, o.Name); err == nil {
			lo.Copy(o)
			err = con.Device.KerberosProfile.Edit(tmpl, ts, vsys, lo)
		}
	}

	if err != nil {
		return err
	}

	return readKerberosProfile(d, meta)
}

func deleteKerberosProfile(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseKerberosProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.KerberosProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.KerberosProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func kerberosProfileSchema(isResource bool) map[string]*schema.Schema {
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
		"server": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Kerberos servers.",
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
						Description: "Kerberos server port number.",
						Default:     88,
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func loadKerberosProfile(d *schema.ResourceData) kerberos.Entry {
	var listing []kerberos.Server
	sl := d.Get("server").([]interface{})
	if len(sl) > 0 {
		listing = make([]kerberos.Server, 0, len(sl))
		for i := range sl {
			x := sl[i].(map[string]interface{})
			listing = append(listing, kerberos.Server{
				Name:   x["name"].(string),
				Server: x["server"].(string),
				Port:   x["port"].(int),
			})
		}
	}

	return kerberos.Entry{
		Name:         d.Get("name").(string),
		AdminUseOnly: d.Get("admin_use_only").(bool),
		Servers:      listing,
	}
}

func saveKerberosProfile(d *schema.ResourceData, o kerberos.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("admin_use_only", o.AdminUseOnly)

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
func parseKerberosProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len-4 ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildKerberosProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
