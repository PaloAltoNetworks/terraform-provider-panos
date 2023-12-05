package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/file"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceFileBlockingSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceFileBlockingSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceFileBlockingSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.FileBlockingProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.FileBlockingProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceFileBlockingSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileBlockingSecurityProfileRead,

		Schema: fileBlockingSecurityProfileSchema(false),
	}
}

func dataSourceFileBlockingSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o file.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildFileBlockingSecurityProfileId(vsys, name)
		o, err = con.Objects.FileBlockingProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildFileBlockingSecurityProfileId(dg, name)
		o, err = con.Objects.FileBlockingProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveFileBlockingSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceFileBlockingSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createFileBlockingSecurityProfile,
		Read:   readFileBlockingSecurityProfile,
		Update: updateFileBlockingSecurityProfile,
		Delete: deleteFileBlockingSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: fileBlockingSecurityProfileSchema(true),
	}
}

func createFileBlockingSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadFileBlockingSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildFileBlockingSecurityProfileId(vsys, o.Name)
		err = con.Objects.FileBlockingProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildFileBlockingSecurityProfileId(dg, o.Name)
		err = con.Objects.FileBlockingProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readFileBlockingSecurityProfile(d, meta)
}

func readFileBlockingSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o file.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseFileBlockingSecurityProfileId(d.Id())
		o, err = con.Objects.FileBlockingProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseFileBlockingSecurityProfileId(d.Id())
		o, err = con.Objects.FileBlockingProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveFileBlockingSecurityProfile(d, o)
	return nil
}

func updateFileBlockingSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadFileBlockingSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.FileBlockingProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.FileBlockingProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.FileBlockingProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.FileBlockingProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readFileBlockingSecurityProfile(d, meta)
}

func deleteFileBlockingSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseFileBlockingSecurityProfileId(d.Id())
		err = con.Objects.FileBlockingProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseFileBlockingSecurityProfileId(d.Id())
		err = con.Objects.FileBlockingProfile.Delete(dg, name)
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
func fileBlockingSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
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
		"rule": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of rule specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Rule name",
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
						Description: "The direction",
						Default:     file.DirectionBoth,
						ValidateFunc: validateStringIn(
							file.DirectionUpload,
							file.DirectionDownload,
							file.DirectionBoth,
						),
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The action to take (note that forward and forward-and-continue are PAN-OS 6.1 only)",
						Default:     file.ActionAlert,
						ValidateFunc: validateStringIn(
							file.ActionAlert,
							file.ActionBlock,
							file.ActionContinue,
							file.ActionForward,
							file.ActionContinueAndForward,
						),
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

func loadFileBlockingSecurityProfile(d *schema.ResourceData) file.Entry {
	var list []interface{}

	var rules []file.Rule
	list = d.Get("rule").([]interface{})
	if len(list) > 0 {
		rules = make([]file.Rule, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			rules = append(rules, file.Rule{
				Name:         elm["name"].(string),
				Applications: asStringList(elm["applications"].([]interface{})),
				FileTypes:    asStringList(elm["file_types"].([]interface{})),
				Direction:    elm["direction"].(string),
				Action:       elm["action"].(string),
			})
		}
	}

	return file.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       rules,
	}
}

func saveFileBlockingSecurityProfile(d *schema.ResourceData, o file.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)

	if len(o.Rules) == 0 {
		d.Set("rule", nil)
	} else {
		list := make([]interface{}, 0, len(o.Rules))
		for _, x := range o.Rules {
			list = append(list, map[string]interface{}{
				"name":         x.Name,
				"applications": x.Applications,
				"file_types":   x.FileTypes,
				"direction":    x.Direction,
				"action":       x.Action,
			})
		}
		if err := d.Set("rule", list); err != nil {
			log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildFileBlockingSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseFileBlockingSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
