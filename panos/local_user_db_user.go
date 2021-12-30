package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/localuserdb/user"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Resource.
func resourceLocalUserDbUser() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateLocalUserDbUser,
		Read:   readLocalUserDbUser,
		Update: createUpdateLocalUserDbUser,
		Delete: deleteLocalUserDbUser,

		Schema: localUserDbUserSchema(),
	}
}

func createUpdateLocalUserDbUser(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadLocalUserDbUser(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	id := buildLocalUserDbUserId(tmpl, ts, vsys, o.Name)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)
	d.Set("password", o.PasswordHash)
	d.Set("live_password", o.PasswordHash)

	switch con := meta.(type) {
	case *pango.Firewall:
		o.PasswordHash, err = con.RequestPasswordHash(o.PasswordHash)
		if err == nil {
			err = con.Device.LocalUserDbUser.Edit(vsys, o)
		}
	case *pango.Panorama:
		o.PasswordHash, err = con.RequestPasswordHash(o.PasswordHash)
		if err == nil {
			err = con.Device.LocalUserDbUser.Edit(tmpl, ts, vsys, o)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("phash", o.PasswordHash)
	return readLocalUserDbUser(d, meta)
}

func readLocalUserDbUser(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o user.Entry

	tmpl, ts, vsys, name := parseLocalUserDbUserId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.LocalUserDbUser.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.LocalUserDbUser.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveLocalUserDbUser(d, o)
	return nil
}

func deleteLocalUserDbUser(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys, name := parseLocalUserDbUserId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.LocalUserDbUser.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.LocalUserDbUser.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func localUserDbUserSchema() map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"live_password": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"phash": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"disabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	return ans
}

func loadLocalUserDbUser(d *schema.ResourceData) user.Entry {
	return user.Entry{
		Name:         d.Get("name").(string),
		PasswordHash: d.Get("password").(string),
		Disabled:     d.Get("disabled").(bool),
	}
}

func saveLocalUserDbUser(d *schema.ResourceData, o user.Entry) {
	d.Set("name", o.Name)
	if d.Get("password").(string) != d.Get("live_password").(string) || d.Get("phash") != o.PasswordHash {
		d.Set("password", "")
	}
	d.Set("disbled", o.Disabled)
}

// Id functions.
func parseLocalUserDbUserId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildLocalUserDbUserId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
