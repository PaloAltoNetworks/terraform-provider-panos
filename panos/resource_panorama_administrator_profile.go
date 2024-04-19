package panos

import (
	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/mgtconfig/user"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

const (
	Name      = "name"
	Password  = "password"
	PublicKey = "public_key"
	RoleBased = "role_based"
	Template  = "template"
	Type      = "type"
)

func resourceAdministratorsUser() *schema.Resource {
	return &schema.Resource{
		Create: createAdministratorsUser,
		Read:   readAdministratorsUser,
		Update: updateAdministratorsUser,
		Delete: deleteAdministratorsUser,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: administratorsUserSchema(),
	}
}

func administratorsUserSchema() map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		Name: &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		Template: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		PublicKey: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		RoleBased: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Password: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		Type: &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return ans
}

func saveUser(d *schema.ResourceData, o user.Entry) {

	if err := d.Set(Name, o.Name); err != nil {
		return
	}
	if err := d.Set(PublicKey, o.PublicKey); err != nil {
		return
	}
	if err := d.Set(RoleBased, o.Role); err != nil {
		return
	}
	if err := d.Set(Type, o.Type); err != nil {
		return
	}
	if err := d.Set(Password, o.PasswordHash); err != nil {
		return
	}
}

func parseUser(d *schema.ResourceData) (user.Entry, string) {
	tmpl := d.Get(Template).(string)

	o := user.Entry{
		Name: d.Get(Name).(string),
	}

	if roleBased := d.Get(RoleBased); roleBased != nil {
		o.Role = roleBased.(string)
	}

	if publicKey := d.Get(PublicKey); publicKey != nil {
		o.PublicKey = publicKey.(string)
	}

	if roleType := d.Get(Type); roleType != nil {
		o.Type = roleType.(string)
	}

	if password := d.Get(Password); password != nil {
		o.PasswordHash = password.(string)
	}

	return o, tmpl
}

func createAdministratorsUser(d *schema.ResourceData, meta interface{}) error {
	pa := meta.(*pango.Panorama)
	o, tmpl := parseUser(d)

	if err := pa.MGTConfig.User.Set(tmpl, o); err != nil {
		return err
	}

	if tmpl != EmptyString {
		d.SetId(buildPanoramaUserId(tmpl, o.Name))
	} else {
		d.SetId(o.Name)
	}

	return readAdministratorsUser(d, meta)
}

func buildPanoramaUserId(a, c string) string {
	return strings.Join([]string{a, c}, IdSeparator)
}

func readAdministratorsUser(d *schema.ResourceData, meta interface{}) error {

	pa := meta.(*pango.Panorama)
	o, tmpl := parseUser(d)

	o, err := pa.MGTConfig.User.Get(tmpl, o.Name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId(EmptyString)
			return nil
		}
		return err
	}

	saveUser(d, o)

	return nil
}

func updateAdministratorsUser(d *schema.ResourceData, meta interface{}) error {

	pa := meta.(*pango.Panorama)
	o, tmpl := parseUser(d)

	lo, err := pa.MGTConfig.User.Get(tmpl, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pa.MGTConfig.User.Edit(tmpl, o); err != nil {
		return err
	}

	return readAdministratorsUser(d, meta)
}

func deleteAdministratorsUser(d *schema.ResourceData, meta interface{}) error {
	pa := meta.(*pango.Panorama)
	o, tmpl := parseUser(d)

	if err := pa.MGTConfig.User.Delete(tmpl, o.Name); err != nil {
		if !isObjectNotFound(err) {
			return err
		}
	}

	d.SetId(EmptyString)
	return nil
}
