package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/custom/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceCustomUrlCategories() *schema.Resource {
	s := listingSchema()
	s["device_group"] = deviceGroupSchema()
	s["vsys"] = vsysSchema("vsys1")

	return &schema.Resource{
		Read: dataSourceCustomUrlCategoriesRead,

		Schema: s,
	}
}

func dataSourceCustomUrlCategoriesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	id := buildCustomUrlCategoryId(dg, vsys, "")

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Objects.CustomUrlCategory.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Objects.CustomUrlCategory.GetList(dg)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceCustomUrlCategory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCustomUrlCategoryRead,

		Schema: customUrlCategorySchema(false),
	}
}

func dataSourceCustomUrlCategoryRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o url.Entry

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildCustomUrlCategoryId(dg, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.CustomUrlCategory.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.CustomUrlCategory.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveCustomUrlCategory(d, o)

	return nil
}

// Resource.
func resourceCustomUrlCategory() *schema.Resource {
	return &schema.Resource{
		Create: createCustomUrlCategory,
		Read:   readCustomUrlCategory,
		Update: updateCustomUrlCategory,
		Delete: deleteCustomUrlCategory,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: customUrlCategorySchema(true),
	}
}

func createCustomUrlCategory(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadCustomUrlCategory(d)

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildCustomUrlCategoryId(dg, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.CustomUrlCategory.Set(vsys, o)
	case *pango.Panorama:
		err = con.Objects.CustomUrlCategory.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readCustomUrlCategory(d, meta)
}

func readCustomUrlCategory(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o url.Entry

	dg, vsys, name := parseCustomUrlCategoryId(d.Id())
	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.CustomUrlCategory.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.CustomUrlCategory.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveCustomUrlCategory(d, o)

	return nil
}

func updateCustomUrlCategory(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo url.Entry
	o := loadCustomUrlCategory(d)

	dg, vsys, _ := parseCustomUrlCategoryId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err = con.Objects.CustomUrlCategory.Get(vsys, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Objects.CustomUrlCategory.Edit(vsys, lo)
		}
	case *pango.Panorama:
		lo, err = con.Objects.CustomUrlCategory.Get(dg, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Objects.CustomUrlCategory.Edit(dg, lo)
		}
	}

	if err != nil {
		return err
	}

	return readCustomUrlCategory(d, meta)
}

func deleteCustomUrlCategory(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, vsys, name := parseCustomUrlCategoryId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.CustomUrlCategory.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Objects.CustomUrlCategory.Delete(dg, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (entry).
func resourceCustomUrlCategoryEntry() *schema.Resource {
	return &schema.Resource{
		Create: createCustomUrlCategoryEntry,
		Read:   readCustomUrlCategoryEntry,
		Delete: deleteCustomUrlCategoryEntry,

		Schema: map[string]*schema.Schema{
			"device_group": deviceGroupSchema(),
			"vsys":         vsysSchema("vsys1"),
			"custom_url_category": {
				Type:        schema.TypeString,
				Description: "The custom URL category name.",
				Required:    true,
				ForceNew:    true,
			},
			"site": {
				Type:        schema.TypeString,
				Description: "The site.",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func createCustomUrlCategoryEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("custom_url_category").(string)
	site := d.Get("site").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)
	d.Set("custom_url_category", name)
	d.Set("site", site)

	id := buildCustomUrlCategoryEntryId(dg, vsys, name, site)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.CustomUrlCategory.SetSite(vsys, name, site)
	case *pango.Panorama:
		err = con.Objects.CustomUrlCategory.SetSite(dg, name, site)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readCustomUrlCategoryEntry(d, meta)
}

func readCustomUrlCategoryEntry(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o url.Entry

	dg, vsys, name, site := parseCustomUrlCategoryEntryId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.CustomUrlCategory.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.CustomUrlCategory.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for _, x := range o.Sites {
		if x == site {
			d.Set("device_group", dg)
			d.Set("vsys", vsys)
			d.Set("custom_url_category", name)
			d.Set("site", site)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteCustomUrlCategoryEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, vsys, name, site := parseCustomUrlCategoryEntryId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.CustomUrlCategory.DeleteSite(vsys, name, site)
	case *pango.Panorama:
		err = con.Objects.CustomUrlCategory.DeleteSite(dg, name, site)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func customUrlCategorySchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:        schema.TypeString,
			Description: "The name.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description.",
			Optional:    true,
		},
		"sites": {
			Type:        schema.TypeList,
			Description: "The site list.",
			Computed:    true,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"type": {
			Type:        schema.TypeString,
			Description: "(PAN-OS 9.0+) The custom URL category type.",
			Optional:    true,
		},
	}

	if !isResource {
		computed(ans, "", []string{"device_group", "vsys", "name"})
	}

	return ans
}

func loadCustomUrlCategory(d *schema.ResourceData) url.Entry {
	return url.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Sites:       asStringList(d.Get("sites").([]interface{})),
		Type:        d.Get("type").(string),
	}
}

func saveCustomUrlCategory(d *schema.ResourceData, o url.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("sites", o.Sites)
	d.Set("type", o.Type)
}

// Id functions.
func buildCustomUrlCategoryId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseCustomUrlCategoryId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildCustomUrlCategoryEntryId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parseCustomUrlCategoryEntryId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}
