package panos

import (
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/edl"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceEdls() *schema.Resource {
	s := listingSchema()
	s["device_group"] = deviceGroupSchema()
	s["vsys"] = vsysSchema("vsys1")

	return &schema.Resource{
		Read: dataSourceEdlsRead,

		Schema: s,
	}
}

func dataSourceEdlsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.Edl.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.Edl.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceEdl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEdlRead,

		Schema: edlSchema(false, 1, nil),
	}
}

func dataSourceEdlRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o edl.Entry

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildEdlId(dg, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.Edl.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.Edl.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveEdl(d, o)

	return nil
}

// Resource.
func resourceEdl() *schema.Resource {
	return &schema.Resource{
		Create: createEdl,
		Read:   readEdl,
		Update: updateEdl,
		Delete: deleteEdl,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: edlSchema(true, 0, []string{"device_group"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: edlUpgradeV0,
			},
		},

		/*
		   Importer: &schema.ResourceImporter{
		       State: schema.ImportStatePassthrough,
		   },
		*/

		Schema: edlSchema(true, 1, nil),
	}
}

func resourcePanoramaEdl() *schema.Resource {
	return &schema.Resource{
		Create: createEdl,
		Read:   readEdl,
		Update: updateEdl,
		Delete: deleteEdl,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: edlSchema(true, 0, []string{"vsys"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: edlUpgradeV0,
			},
		},

		/*
		   Importer: &schema.ResourceImporter{
		       State: schema.ImportStatePassthrough,
		   },
		*/

		Schema: edlSchema(true, 1, nil),
	}
}

func edlUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["vsys"]; !ok {
		raw["vsys"] = "vsys1"
	}
	if _, ok := raw["device_group"]; !ok {
		raw["device_group"] = "shared"
	}

	val := raw["type"].(string)
	if val == "predefined" {
		raw["type"] = edl.TypePredefinedIp
	}

	return raw, nil
}

func createEdl(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo edl.Entry
	o := loadEdl(d)

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildEdlId(dg, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		if err = con.Objects.Edl.Set(vsys, o); err == nil {
			lo, err = con.Objects.Edl.Get(vsys, o.Name)
		}
	case *pango.Panorama:
		if err = con.Objects.Edl.Set(dg, o); err == nil {
			lo, err = con.Objects.Edl.Get(dg, o.Name)
		}
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("password", o.Password)
	d.Set("password_enc", lo.Password)

	return readEdl(d, meta)
}

func readEdl(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o edl.Entry

	// Migrate the ID.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 2 {
		switch meta.(type) {
		case *pango.Firewall:
			d.SetId(buildEdlId("shared", tok[0], tok[1]))
		case *pango.Panorama:
			d.SetId(buildEdlId(tok[0], "vsys1", tok[1]))
		}
	}

	dg, vsys, name := parseEdlId(d.Id())
	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.Edl.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.Edl.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveEdl(d, o)
	return nil
}

func updateEdl(d *schema.ResourceData, meta interface{}) error {
	dg, vsys, _ := parseEdlId(d.Id())
	o := loadEdl(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Objects.Edl.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Edl.Edit(vsys, lo); err != nil {
			return err
		}
		lo2, err := con.Objects.Edl.Get(vsys, lo.Name)
		if err != nil {
			return err
		}
		d.Set("password_enc", lo2.Password)
	case *pango.Panorama:
		lo, err := con.Objects.Edl.Get(dg, o.Name)
		lo.Copy(o)
		if err = con.Objects.Edl.Edit(dg, lo); err != nil {
			return err
		}
		lo2, err := con.Objects.Edl.Get(dg, lo.Name)
		if err != nil {
			return err
		}
		d.Set("password_enc", lo2.Password)
	}

	d.Set("password", o.Password)

	return readEdl(d, meta)
}

func deleteEdl(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, vsys, name := parseEdlId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.Edl.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Objects.Edl.Delete(dg, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func edlSchema(isResource bool, schemaVersion int, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  edl.TypeIp,
			ValidateFunc: validateStringIn(
				edl.TypeIp,
				edl.TypeDomain,
				edl.TypeUrl,
				edl.TypePredefinedIp,
				edl.TypePredefinedUrl,
			),
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"source": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"certificate_profile": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"password_enc": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"repeat": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  edl.RepeatHourly,
			ValidateFunc: validateStringIn(
				edl.RepeatEveryFiveMinutes,
				edl.RepeatHourly,
				edl.RepeatDaily,
				edl.RepeatWeekly,
				edl.RepeatMonthly,
			),
		},
		"repeat_at": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"repeat_day_of_week": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validateStringIn(
				"sunday",
				"monday",
				"tuesday",
				"wednesday",
				"thursday",
				"friday",
				"saturday",
				"",
			),
		},
		"repeat_day_of_month": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateIntInRange(0, 31),
		},
		"exceptions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if schemaVersion == 0 {
		ans["type"].ValidateFunc = validateStringIn(
			edl.TypeIp,
			edl.TypeDomain,
			edl.TypeUrl,
			"predefined",
		)
	}

	return ans
}

func loadEdl(d *schema.ResourceData) edl.Entry {
	return edl.Entry{
		Name:               d.Get("name").(string),
		Type:               d.Get("type").(string),
		Description:        d.Get("description").(string),
		Source:             d.Get("source").(string),
		CertificateProfile: d.Get("certificate_profile").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		Repeat:             d.Get("repeat").(string),
		RepeatAt:           d.Get("repeat_at").(string),
		RepeatDayOfWeek:    d.Get("repeat_day_of_week").(string),
		RepeatDayOfMonth:   d.Get("repeat_day_of_month").(int),
		Exceptions:         asStringList(d.Get("exceptions").([]interface{})),
	}
}

func saveEdl(d *schema.ResourceData, o edl.Entry) {
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	d.Set("source", o.Source)
	d.Set("certificate_profile", o.CertificateProfile)
	d.Set("username", o.Username)
	d.Set("repeat", o.Repeat)
	d.Set("repeat_at", o.RepeatAt)
	d.Set("repeat_day_of_week", o.RepeatDayOfWeek)
	d.Set("repeat_day_of_month", o.RepeatDayOfMonth)
	d.Set("exceptions", o.Exceptions)

	if d.Get("password_enc").(string) != o.Password {
		d.Set("password", "")
	}
}

// Id functions.
func parseEdlId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildEdlId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
