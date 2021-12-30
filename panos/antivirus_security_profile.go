package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	av "github.com/PaloAltoNetworks/pango/objs/profile/security/virus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceAntivirusSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceAntivirusSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceAntivirusSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.AntivirusProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.AntivirusProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceAntivirusSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntivirusSecurityProfileRead,

		Schema: antivirusSecurityProfileSchema(false),
	}
}

func dataSourceAntivirusSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o av.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntivirusSecurityProfileId(vsys, name)
		o, err = con.Objects.AntivirusProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntivirusSecurityProfileId(dg, name)
		o, err = con.Objects.AntivirusProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveAntivirusSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceAntivirusSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createAntivirusSecurityProfile,
		Read:   readAntivirusSecurityProfile,
		Update: updateAntivirusSecurityProfile,
		Delete: deleteAntivirusSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: antivirusSecurityProfileSchema(true),
	}
}

func createAntivirusSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadAntivirusSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntivirusSecurityProfileId(vsys, o.Name)
		err = con.Objects.AntivirusProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntivirusSecurityProfileId(dg, o.Name)
		err = con.Objects.AntivirusProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAntivirusSecurityProfile(d, meta)
}

func readAntivirusSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o av.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAntivirusSecurityProfileId(d.Id())
		o, err = con.Objects.AntivirusProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseAntivirusSecurityProfileId(d.Id())
		o, err = con.Objects.AntivirusProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAntivirusSecurityProfile(d, o)
	return nil
}

func updateAntivirusSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadAntivirusSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.AntivirusProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntivirusProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.AntivirusProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntivirusProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readAntivirusSecurityProfile(d, meta)
}

func deleteAntivirusSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAntivirusSecurityProfileId(d.Id())
		err = con.Objects.AntivirusProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseAntivirusSecurityProfileId(d.Id())
		err = con.Objects.AntivirusProfile.Delete(dg, name)
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
func antivirusSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
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
		"packet_capture": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Set to true to enable packet capture",
		},
		"threat_exceptions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of threat exceptions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"decoder": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Decoders",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Decoder name",
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Decoder action",
						Default:     "default",
						ValidateFunc: validateStringIn(
							"",
							av.Default,
							av.Allow,
							av.Alert,
							av.Block,
							av.Drop,
							av.ResetClient,
							av.ResetServer,
							av.ResetBoth,
						),
					},
					"wildfire_action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Wildfire action",
						Default:     "default",
						ValidateFunc: validateStringIn(
							"",
							av.Default,
							av.Allow,
							av.Alert,
							av.Block,
							av.Drop,
							av.ResetClient,
							av.ResetServer,
							av.ResetBoth,
						),
					},
					"machine_learning_action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "(PAN-OS 10.0+) ML action",
					},
				},
			},
		},
		"application_exception": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Application exception specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"application": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The application name",
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Action to take for this application",
						ValidateFunc: validateStringIn(
							"",
							av.Default,
							av.Allow,
							av.Alert,
							av.Block,
							av.Drop,
							av.ResetClient,
							av.ResetServer,
							av.ResetBoth,
						),
					},
				},
			},
		},
		"machine_learning_model": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Machine learning model spec",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"model": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "ML model",
					},
					"action": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Action to take",
					},
				},
			},
		},
		"machine_learning_exception": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Machine learning exception spec",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Machine learning exception name",
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description",
					},
					"filename": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Filename",
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

func loadAntivirusSecurityProfile(d *schema.ResourceData) av.Entry {
	var list []interface{}

	var decoders []av.Decoder
	list = d.Get("decoder").([]interface{})
	if len(list) > 0 {
		decoders = make([]av.Decoder, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			decoders = append(decoders, av.Decoder{
				Name:                  elm["name"].(string),
				Action:                elm["action"].(string),
				WildfireAction:        elm["wildfire_action"].(string),
				MachineLearningAction: elm["machine_learning_action"].(string),
			})
		}
	}

	var appExceptions []av.ApplicationException
	list = d.Get("application_exception").([]interface{})
	if len(list) > 0 {
		appExceptions = make([]av.ApplicationException, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			appExceptions = append(appExceptions, av.ApplicationException{
				Application: elm["application"].(string),
				Action:      elm["action"].(string),
			})
		}
	}

	var mlModels []av.MachineLearningModel
	list = d.Get("machine_learning_model").([]interface{})
	if len(list) > 0 {
		mlModels = make([]av.MachineLearningModel, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			mlModels = append(mlModels, av.MachineLearningModel{
				Model:  elm["model"].(string),
				Action: elm["action"].(string),
			})
		}
	}

	var mlExceptions []av.MachineLearningException
	list = d.Get("machine_learning_exception").([]interface{})
	if len(list) > 0 {
		mlExceptions = make([]av.MachineLearningException, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			mlExceptions = append(mlExceptions, av.MachineLearningException{
				Name:        elm["name"].(string),
				Description: elm["description"].(string),
				Filename:    elm["filename"].(string),
			})
		}
	}

	return av.Entry{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		PacketCapture:             d.Get("packet_capture").(bool),
		Decoders:                  decoders,
		ApplicationExceptions:     appExceptions,
		ThreatExceptions:          asStringList(d.Get("threat_exceptions").([]interface{})),
		MachineLearningModels:     mlModels,
		MachineLearningExceptions: mlExceptions,
	}
}

func saveAntivirusSecurityProfile(d *schema.ResourceData, o av.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("packet_capture", o.PacketCapture)
	if err := d.Set("threat_exceptions", o.ThreatExceptions); err != nil {
		log.Printf("[WARN] Error setting 'threat_exceptions' for %q: %s", d.Id(), err)
	}

	if len(o.Decoders) == 0 {
		d.Set("decoder", nil)
	} else {
		list := make([]interface{}, 0, len(o.Decoders))
		for _, x := range o.Decoders {
			list = append(list, map[string]interface{}{
				"name":                    x.Name,
				"action":                  x.Action,
				"wildfire_action":         x.WildfireAction,
				"machine_learning_action": x.MachineLearningAction,
			})
		}
		if err := d.Set("decoder", list); err != nil {
			log.Printf("[WARN] Error setting 'decoder' for %q: %s", d.Id(), err)
		}
	}

	if len(o.ApplicationExceptions) == 0 {
		d.Set("application_exception", nil)
	} else {
		list := make([]interface{}, 0, len(o.ApplicationExceptions))
		for _, x := range o.ApplicationExceptions {
			list = append(list, map[string]interface{}{
				"application": x.Application,
				"action":      x.Action,
			})
		}
		if err := d.Set("application_exception", list); err != nil {
			log.Printf("[WARN] Error setting 'application_exception' for %q: %s", d.Id(), err)
		}
	}

	if len(o.MachineLearningModels) == 0 {
		d.Set("machine_learning_model", nil)
	} else {
		list := make([]interface{}, 0, len(o.MachineLearningModels))
		for _, x := range o.MachineLearningModels {
			list = append(list, map[string]interface{}{
				"model":  x.Model,
				"action": x.Action,
			})
		}
		if err := d.Set("machine_learning_model", list); err != nil {
			log.Printf("[WARN] Error setting 'machine_learning_model' for %q: %s", d.Id(), err)
		}
	}

	if len(o.MachineLearningExceptions) == 0 {
		d.Set("machine_learning_exception", nil)
	} else {
		list := make([]interface{}, 0, len(o.MachineLearningExceptions))
		for _, x := range o.MachineLearningExceptions {
			list = append(list, map[string]interface{}{
				"name":        x.Name,
				"description": x.Description,
				"filename":    x.Filename,
			})
		}
		if err := d.Set("machine_learning_exception", list); err != nil {
			log.Printf("[WARN] Error setting 'machine_learning_exceptions' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildAntivirusSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseAntivirusSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
