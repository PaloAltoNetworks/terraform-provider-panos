package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/dug"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceDynamicUserGroups() *schema.Resource {
	s := map[string]*schema.Schema{
		"vsys":         vsysSchema(),
		"device_group": deviceGroupSchema(),
	}

	for key, val := range listingSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceDynamicUserGroupsRead,

		Schema: s,
	}
}

func dataSourceDynamicUserGroupsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.DynamicUserGroup.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.DynamicUserGroup.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceDynamicUserGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDynamicUserGroupRead,

		Schema: dynamicUserGroupSchema(false),
	}
}

func dataSourceDynamicUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o dug.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDynamicUserGroupId(vsys, name)
		o, err = con.Objects.DynamicUserGroup.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDynamicUserGroupId(dg, name)
		o, err = con.Objects.DynamicUserGroup.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveDynamicUserGroup(d, o)

	return nil
}

// Resource.
func resourceDynamicUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: createDynamicUserGroup,
		Read:   readDynamicUserGroup,
		Update: updateDynamicUserGroup,
		Delete: deleteDynamicUserGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: dynamicUserGroupSchema(true),
	}
}

func createDynamicUserGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadDynamicUserGroup(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDynamicUserGroupId(vsys, o.Name)
		err = con.Objects.DynamicUserGroup.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDynamicUserGroupId(dg, o.Name)
		err = con.Objects.DynamicUserGroup.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readDynamicUserGroup(d, meta)
}

func readDynamicUserGroup(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o dug.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDynamicUserGroupId(d.Id())
		o, err = con.Objects.DynamicUserGroup.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseDynamicUserGroupId(d.Id())
		o, err = con.Objects.DynamicUserGroup.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDynamicUserGroup(d, o)
	return nil
}

func updateDynamicUserGroup(d *schema.ResourceData, meta interface{}) error {
	o := loadDynamicUserGroup(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.DynamicUserGroup.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DynamicUserGroup.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.DynamicUserGroup.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DynamicUserGroup.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readDynamicUserGroup(d, meta)
}

func deleteDynamicUserGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDynamicUserGroupId(d.Id())
		err = con.Objects.DynamicUserGroup.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseDynamicUserGroupId(d.Id())
		err = con.Objects.DynamicUserGroup.Delete(dg, name)
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
func dynamicUserGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema(),
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The object name",
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
		"filter": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The dynamic user group filter",
		},
		"tags": tagSchema(),
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	return ans
}

func loadDynamicUserGroup(d *schema.ResourceData) dug.Entry {
	return dug.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Filter:      d.Get("filter").(string),
		Tags:        asStringList(d.Get("tags").([]interface{})),
	}
}

func saveDynamicUserGroup(d *schema.ResourceData, o dug.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("filter", o.Filter)
	if err := d.Set("tags", o.Tags); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseDynamicUserGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildDynamicUserGroupId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}
