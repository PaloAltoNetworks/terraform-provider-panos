package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/wildfire"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceWildfireAnalysisSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema()
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceWildfireAnalysisSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceWildfireAnalysisSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.WildfireAnalysisProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.WildfireAnalysisProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceWildfireAnalysisSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWildfireAnalysisSecurityProfileRead,

		Schema: wildfireAnalysisSecurityProfileSchema(false),
	}
}

func dataSourceWildfireAnalysisSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o wildfire.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildWildfireAnalysisSecurityProfileId(vsys, name)
		o, err = con.Objects.WildfireAnalysisProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildWildfireAnalysisSecurityProfileId(dg, name)
		o, err = con.Objects.WildfireAnalysisProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveWildfireAnalysisSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceWildfireAnalysisSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createWildfireAnalysisSecurityProfile,
		Read:   readWildfireAnalysisSecurityProfile,
		Update: updateWildfireAnalysisSecurityProfile,
		Delete: deleteWildfireAnalysisSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: wildfireAnalysisSecurityProfileSchema(true),
	}
}

func createWildfireAnalysisSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadWildfireAnalysisSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildWildfireAnalysisSecurityProfileId(vsys, o.Name)
		err = con.Objects.WildfireAnalysisProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildWildfireAnalysisSecurityProfileId(dg, o.Name)
		err = con.Objects.WildfireAnalysisProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readWildfireAnalysisSecurityProfile(d, meta)
}

func readWildfireAnalysisSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o wildfire.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseWildfireAnalysisSecurityProfileId(d.Id())
		o, err = con.Objects.WildfireAnalysisProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseWildfireAnalysisSecurityProfileId(d.Id())
		o, err = con.Objects.WildfireAnalysisProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveWildfireAnalysisSecurityProfile(d, o)
	return nil
}

func updateWildfireAnalysisSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadWildfireAnalysisSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.WildfireAnalysisProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.WildfireAnalysisProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.WildfireAnalysisProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.WildfireAnalysisProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readWildfireAnalysisSecurityProfile(d, meta)
}

func deleteWildfireAnalysisSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseWildfireAnalysisSecurityProfileId(d.Id())
		err = con.Objects.WildfireAnalysisProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseWildfireAnalysisSecurityProfileId(d.Id())
		err = con.Objects.WildfireAnalysisProfile.Delete(dg, name)
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
func wildfireAnalysisSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
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
						Description: "Direction",
						Default:     wildfire.DirectionBoth,
						ValidateFunc: validateStringIn(
							wildfire.DirectionUpload,
							wildfire.DirectionDownload,
							wildfire.DirectionBoth,
						),
					},
					"analysis": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Analysis setting",
						Default:     wildfire.AnalysisPublicCloud,
						ValidateFunc: validateStringIn(
							wildfire.AnalysisPublicCloud,
							wildfire.AnalysisPrivateCloud,
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

func loadWildfireAnalysisSecurityProfile(d *schema.ResourceData) wildfire.Entry {
	var list []interface{}

	var rules []wildfire.Rule
	list = d.Get("rule").([]interface{})
	if len(list) > 0 {
		rules = make([]wildfire.Rule, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			rules = append(rules, wildfire.Rule{
				Name:         elm["name"].(string),
				Applications: asStringList(elm["applications"].([]interface{})),
				FileTypes:    asStringList(elm["file_types"].([]interface{})),
				Direction:    elm["direction"].(string),
				Analysis:     elm["analysis"].(string),
			})
		}
	}

	return wildfire.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       rules,
	}
}

func saveWildfireAnalysisSecurityProfile(d *schema.ResourceData, o wildfire.Entry) {
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
				"analysis":     x.Analysis,
			})
		}

		if err := d.Set("rule", list); err != nil {
			log.Printf("[WARN] Error setting 'rule' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildWildfireAnalysisSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseWildfireAnalysisSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
