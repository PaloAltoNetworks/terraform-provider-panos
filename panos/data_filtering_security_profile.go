package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/profile/security/data"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceDataFilteringSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceDataFilteringSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceDataFilteringSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.DataFilteringProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.DataFilteringProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceDataFilteringSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataFilteringSecurityProfileRead,

		Schema: dataFilteringSecurityProfileSchema(false),
	}
}

func dataSourceDataFilteringSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o data.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDataFilteringSecurityProfileId(vsys, name)
		o, err = con.Objects.DataFilteringProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDataFilteringSecurityProfileId(dg, name)
		o, err = con.Objects.DataFilteringProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveDataFilteringSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceDataFilteringSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createDataFilteringSecurityProfile,
		Read:   readDataFilteringSecurityProfile,
		Update: updateDataFilteringSecurityProfile,
		Delete: deleteDataFilteringSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: dataFilteringSecurityProfileSchema(true),
	}
}

func createDataFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadDataFilteringSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDataFilteringSecurityProfileId(vsys, o.Name)
		err = con.Objects.DataFilteringProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDataFilteringSecurityProfileId(dg, o.Name)
		err = con.Objects.DataFilteringProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readDataFilteringSecurityProfile(d, meta)
}

func readDataFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o data.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDataFilteringSecurityProfileId(d.Id())
		o, err = con.Objects.DataFilteringProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseDataFilteringSecurityProfileId(d.Id())
		o, err = con.Objects.DataFilteringProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDataFilteringSecurityProfile(d, o)
	return nil
}

func updateDataFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadDataFilteringSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.DataFilteringProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DataFilteringProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.DataFilteringProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DataFilteringProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readDataFilteringSecurityProfile(d, meta)
}

func deleteDataFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDataFilteringSecurityProfileId(d.Id())
		err = con.Objects.DataFilteringProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseDataFilteringSecurityProfileId(d.Id())
		err = con.Objects.DataFilteringProfile.Delete(dg, name)
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
func dataFilteringSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys":         vsysSchema("vsys1"),
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
		"data_capture": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Data capture",
		},
		"rule": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of rule specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Rule name",
					},
					"data_pattern": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The data pattern to use",
					},
					"applications": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of applications",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"file_types": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of file types",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"direction": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Direction",
						Default:     data.DirectionBoth,
						ValidateFunc: validateStringIn(
							data.DirectionUpload,
							data.DirectionDownload,
							data.DirectionBoth,
						),
					},
					"alert_threshold": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alert threshold",
					},
					"block_threshold": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block threshold",
					},
					"log_severity": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "(PAN-OS 8.0+) Log severity",
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

func loadDataFilteringSecurityProfile(d *schema.ResourceData) data.Entry {
	var list []interface{}

	var rules []data.Rule
	list = d.Get("rule").([]interface{})
	if len(list) > 0 {
		rules = make([]data.Rule, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			rules = append(rules, data.Rule{
				DataPattern:    elm["data_pattern"].(string),
				Applications:   asStringList(elm["applications"].([]interface{})),
				FileTypes:      asStringList(elm["file_types"].([]interface{})),
				Direction:      elm["direction"].(string),
				AlertThreshold: elm["alert_threshold"].(int),
				BlockThreshold: elm["block_threshold"].(int),
				LogSeverity:    elm["log_severity"].(string),
			})
		}
	}

	return data.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DataCapture: d.Get("data_capture").(bool),
		Rules:       rules,
	}
}

func saveDataFilteringSecurityProfile(d *schema.ResourceData, o data.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("data_capture", o.DataCapture)

	if len(o.Rules) == 0 {
		d.Set("rule", nil)
	} else {
		list := make([]interface{}, 0, len(o.Rules))
		for _, x := range o.Rules {
			list = append(list, map[string]interface{}{
				"name":            x.Name,
				"data_pattern":    x.DataPattern,
				"applications":    x.Applications,
				"file_types":      x.FileTypes,
				"direction":       x.Direction,
				"alert_threshold": x.AlertThreshold,
				"block_threshold": x.BlockThreshold,
				"log_severity":    x.LogSeverity,
			})
		}

		if err := d.Set("rule", list); err != nil {
			log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildDataFilteringSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseDataFilteringSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
