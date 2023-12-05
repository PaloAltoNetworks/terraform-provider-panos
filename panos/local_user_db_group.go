package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/localuserdb/group"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceLocalUserDbGroups() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("shared")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceLocalUserDbGroupsRead,

		Schema: s,
	}
}

func dataSourceLocalUserDbGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	id := buildLocalUserDbGroupId(tmpl, ts, vsys, "")

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.LocalUserDbGroup.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.LocalUserDbGroup.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceLocalUserDbGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalUserDbGroupRead,

		Schema: localUserDbGroupSchema(false),
	}
}

func dataSourceLocalUserDbGroupRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o group.Entry

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)
	id := buildLocalUserDbGroupId(tmpl, ts, vsys, name)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.LocalUserDbGroup.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.LocalUserDbGroup.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveLocalUserDbGroup(d, o)

	return nil
}

// Resource.
func resourceLocalUserDbGroup() *schema.Resource {
	return &schema.Resource{
		Create: createLocalUserDbGroup,
		Read:   readLocalUserDbGroup,
		Update: updateLocalUserDbGroup,
		Delete: deleteLocalUserDbGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: localUserDbGroupSchema(true),
	}
}

func createLocalUserDbGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadLocalUserDbGroup(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	id := buildLocalUserDbGroupId(tmpl, ts, vsys, o.Name)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.LocalUserDbGroup.Set(vsys, o)
	case *pango.Panorama:
		err = con.Device.LocalUserDbGroup.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readLocalUserDbGroup(d, meta)
}

func readLocalUserDbGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o group.Entry
	tmpl, ts, vsys, name := parseLocalUserDbGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.LocalUserDbGroup.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.LocalUserDbGroup.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveLocalUserDbGroup(d, o)
	return nil
}

func updateLocalUserDbGroup(d *schema.ResourceData, meta interface{}) error {
	o := loadLocalUserDbGroup(d)

	tmpl, ts, vsys, _ := parseLocalUserDbGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Device.LocalUserDbGroup.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.LocalUserDbGroup.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Device.LocalUserDbGroup.Get(tmpl, ts, vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.LocalUserDbGroup.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
	}

	return readLocalUserDbGroup(d, meta)
}

func deleteLocalUserDbGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys, name := parseLocalUserDbGroupId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.LocalUserDbGroup.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.LocalUserDbGroup.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func localUserDbGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"users": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func loadLocalUserDbGroup(d *schema.ResourceData) group.Entry {
	return group.Entry{
		Name:  d.Get("name").(string),
		Users: setAsList(d.Get("users").(*schema.Set)),
	}
}

func saveLocalUserDbGroup(d *schema.ResourceData, o group.Entry) {
	d.Set("name", o.Name)
	if err := d.Set("users", listAsSet(o.Users)); err != nil {
		log.Printf("[WARN] Error setting 'users' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseLocalUserDbGroupId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildLocalUserDbGroupId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
