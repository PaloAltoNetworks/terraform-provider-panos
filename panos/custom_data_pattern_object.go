package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/custom/data"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceCustomDataPatternObjects() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema()
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceCustomDataPatternObjectsRead,

		Schema: s,
	}
}

func dataSourceCustomDataPatternObjectsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.DataPattern.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.DataPattern.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceCustomDataPatternObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCustomDataPatternObjectRead,

		Schema: customDataPatternObjectSchema(false),
	}
}

func dataSourceCustomDataPatternObjectRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o data.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildCustomDataPatternObjectId(vsys, name)
		o, err = con.Objects.DataPattern.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildCustomDataPatternObjectId(dg, name)
		o, err = con.Objects.DataPattern.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveCustomDataPatternObject(d, o)

	return nil
}

// Resource.
func resourceCustomDataPatternObject() *schema.Resource {
	return &schema.Resource{
		Create: createCustomDataPatternObject,
		Read:   readCustomDataPatternObject,
		Update: updateCustomDataPatternObject,
		Delete: deleteCustomDataPatternObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: customDataPatternObjectSchema(true),
	}
}

func createCustomDataPatternObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadCustomDataPatternObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildCustomDataPatternObjectId(vsys, o.Name)
		err = con.Objects.DataPattern.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildCustomDataPatternObjectId(dg, o.Name)
		err = con.Objects.DataPattern.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readCustomDataPatternObject(d, meta)
}

func readCustomDataPatternObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o data.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseCustomDataPatternObjectId(d.Id())
		o, err = con.Objects.DataPattern.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseCustomDataPatternObjectId(d.Id())
		o, err = con.Objects.DataPattern.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveCustomDataPatternObject(d, o)
	return nil
}

func updateCustomDataPatternObject(d *schema.ResourceData, meta interface{}) error {
	o := loadCustomDataPatternObject(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.DataPattern.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DataPattern.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.DataPattern.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DataPattern.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readCustomDataPatternObject(d, meta)
}

func deleteCustomDataPatternObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseCustomDataPatternObjectId(d.Id())
		err = con.Objects.DataPattern.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseCustomDataPatternObjectId(d.Id())
		err = con.Objects.DataPattern.Delete(dg, name)
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
func customDataPatternObjectSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys":         vsysSchema(),
		"device_group": deviceGroupSchema(),
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Security profile name",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type",
			Default:     data.TypeFileProperties,
			ValidateFunc: validateStringIn(
				data.TypePredefined,
				data.TypeRegex,
				data.TypeFileProperties,
			),
		},
		"predefined_pattern": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of predefined pattern specs",
			ConflictsWith: []string{
				"regex",
				"file_property",
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"file_types": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of file types",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"regex": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of predefined pattern specs",
			ConflictsWith: []string{
				"predefined_pattern",
				"file_property",
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"file_types": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of file types",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"regex": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The regex",
					},
				},
			},
		},
		"file_property": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of predefined pattern specs",
			ConflictsWith: []string{
				"predefined_pattern",
				"regex",
			},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"file_type": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The file type",
					},
					"file_property": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "File property",
					},
					"property_value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Property value",
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	return ans
}

func loadCustomDataPatternObject(d *schema.ResourceData) data.Entry {
	var list []interface{}

	var patterns []data.PredefinedPattern
	list = d.Get("predefined_pattern").([]interface{})
	if len(list) > 0 {
		patterns = make([]data.PredefinedPattern, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			patterns = append(patterns, data.PredefinedPattern{
				Name:      elm["name"].(string),
				FileTypes: asStringList(elm["file_types"].([]interface{})),
			})
		}
	}

	var regexes []data.Regex
	list = d.Get("regex").([]interface{})
	if len(list) > 0 {
		regexes = make([]data.Regex, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			regexes = append(regexes, data.Regex{
				Name:      elm["name"].(string),
				FileTypes: asStringList(elm["file_types"].([]interface{})),
				Regex:     elm["regex"].(string),
			})
		}
	}

	var props []data.FileProperty
	list = d.Get("file_property").([]interface{})
	if len(list) > 0 {
		props = make([]data.FileProperty, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			props = append(props, data.FileProperty{
				Name:          elm["name"].(string),
				FileType:      elm["file_type"].(string),
				FileProperty:  elm["file_property"].(string),
				PropertyValue: elm["property_value"].(string),
			})
		}
	}

	return data.Entry{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Type:               d.Get("type").(string),
		PredefinedPatterns: patterns,
		Regexes:            regexes,
		FileProperties:     props,
	}
}

func saveCustomDataPatternObject(d *schema.ResourceData, o data.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("type", o.Type)

	if len(o.PredefinedPatterns) == 0 {
		d.Set("predefined_pattern", nil)
	} else {
		list := make([]interface{}, 0, len(o.PredefinedPatterns))
		for _, x := range o.PredefinedPatterns {
			list = append(list, map[string]interface{}{
				"name":       x.Name,
				"file_types": x.FileTypes,
			})
		}

		if err := d.Set("predefined_pattern", list); err != nil {
			log.Printf("[WARN] Error setting 'predefined_pattern' for %q: %s", d.Id(), err)
		}
	}

	if len(o.Regexes) == 0 {
		d.Set("regex", nil)
	} else {
		list := make([]interface{}, 0, len(o.Regexes))
		for _, x := range o.Regexes {
			list = append(list, map[string]interface{}{
				"name":       x.Name,
				"file_types": x.FileTypes,
				"regex":      x.Regex,
			})
		}

		if err := d.Set("regex", list); err != nil {
			log.Printf("[WARN] Error setting 'regex' for %q: %s", d.Id(), err)
		}
	}

	if len(o.FileProperties) == 0 {
		d.Set("file_property", nil)
	} else {
		list := make([]interface{}, 0, len(o.FileProperties))
		for _, x := range o.FileProperties {
			list = append(list, map[string]interface{}{
				"name":           x.Name,
				"file_type":      x.FileType,
				"file_property":  x.FileProperty,
				"property_value": x.PropertyValue,
			})
		}

		if err := d.Set("file_property", list); err != nil {
			log.Printf("[WARN] Error setting 'file_property' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildCustomDataPatternObjectId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseCustomDataPatternObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
