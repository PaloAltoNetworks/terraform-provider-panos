package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceZones() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["template"] = templateSchema(true)
	s["template_stack"] = templateStackSchema()

	return &schema.Resource{
		Read: dataSourceZonesRead,

		Schema: s,
	}
}

func dataSourceZonesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildZoneId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Network.Zone.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Network.Zone.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,

		Schema: zoneSchema(false, nil),
	}
}

func dataSourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o zone.Entry

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildZoneId(tmpl, ts, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.Zone.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Network.Zone.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveZone(d, o)

	return nil
}

// Resource.
func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: createZone,
		Read:   readZone,
		Update: updateZone,
		Delete: deleteZone,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: zoneSchema(true, []string{"template", "template_stack"}),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: zoneUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: zoneSchema(true, nil),
	}
}

func zoneUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["template"]; !ok {
		raw["template"] = ""
	}
	if _, ok := raw["template_stack"]; !ok {
		raw["template_stack"] = ""
	}

	return raw, nil
}

func resourcePanoramaZone() *schema.Resource {
	return &schema.Resource{
		Create: createZone,
		Read:   readZone,
		Update: updateZone,
		Delete: deleteZone,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: zoneSchema(true, nil),
	}
}

func createZone(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadZone(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildZoneId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.Zone.Set(vsys, o)
	case *pango.Panorama:
		err = con.Network.Zone.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readZone(d, meta)
}

func readZone(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o zone.Entry

	// Migrate the firewall Id.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 2 {
		d.SetId(buildZoneId("", "", tok[0], tok[1]))
	}
	tmpl, ts, vsys, name := parseZoneId(d.Id())

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.Zone.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Network.Zone.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveZone(d, o)

	return nil
}

func updateZone(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo zone.Entry
	tmpl, ts, vsys, _ := parseZoneId(d.Id())
	o := loadZone(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err = con.Network.Zone.Get(vsys, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Network.Zone.Edit(vsys, lo)
		}
	case *pango.Panorama:
		lo, err = con.Network.Zone.Get(tmpl, ts, vsys, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Network.Zone.Edit(tmpl, ts, vsys, lo)
		}
	}

	if err != nil {
		return err
	}

	return readZone(d, meta)
}

func deleteZone(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys, name := parseZoneId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.Zone.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Network.Zone.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (entry).
func resourceZoneEntry() *schema.Resource {
	return &schema.Resource{
		Create: createZoneEntry,
		Read:   readZoneEntry,
		Delete: deleteZoneEntry,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: zoneEntrySchema(false, false),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: zoneEntryUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: zoneEntrySchema(true, true),
	}
}

func resourcePanoramaZoneEntry() *schema.Resource {
	return &schema.Resource{
		Create: createZoneEntry,
		Read:   readZoneEntry,
		Delete: deleteZoneEntry,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: zoneEntrySchema(true, false),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: zoneEntryUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: zoneEntrySchema(true, true),
	}
}

func zoneEntryUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["template"]; !ok {
		raw["template"] = ""
	}
	if _, ok := raw["template_stack"]; !ok {
		raw["template_stack"] = ""
	}

	return raw, nil
}

func createZoneEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl := d.Get("template").(string)
	ts := ""
	vsys := d.Get("vsys").(string)
	zoneName := d.Get("zone").(string)
	mode := d.Get("mode").(string)
	iface := d.Get("interface").(string)

	d.Set("template", tmpl)
	d.Set("vsys", vsys)
	d.Set("zone", zoneName)
	d.Set("mode", mode)
	d.Set("interface", iface)

	id := buildZoneEntryId(tmpl, ts, vsys, zoneName, mode, iface)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.Zone.SetInterface(vsys, zoneName, mode, iface)
	case *pango.Panorama:
		err = con.Network.Zone.SetInterface(tmpl, ts, vsys, zoneName, mode, iface)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readZoneEntry(d, meta)
}

func readZoneEntry(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o zone.Entry

	// Migrate the firewall Id to the universal Id format.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 4 {
		d.SetId(buildZoneEntryId("", "", tok[0], tok[1], tok[2], tok[3]))
	}

	tmpl, ts, vsys, zoneName, mode, iface := parseZoneEntryId(d.Id())

	d.Set("template", tmpl)
	d.Set("vsys", vsys)
	d.Set("zone", zoneName)
	d.Set("mode", mode)
	d.Set("interface", iface)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Network.Zone.Get(vsys, zoneName)
	case *pango.Panorama:
		o, err = con.Network.Zone.Get(tmpl, ts, vsys, zoneName)
	}

	/*
	   There are three possibilities to blank the Id:

	   1. The zone isn't present
	   2. The mode is incorrect
	   3. The interface isn't in the interface list
	*/
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if o.Mode != mode {
		d.SetId("")
		return nil
	}

	for _, x := range o.Interfaces {
		if x == iface {
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteZoneEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, vsys, zoneName, mode, iface := parseZoneEntryId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Network.Zone.DeleteInterface(vsys, zoneName, mode, iface)
	case *pango.Panorama:
		err = con.Network.Zone.DeleteInterface(tmpl, ts, vsys, zoneName, mode, iface)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func zoneSchema(isResource bool, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:        schema.TypeString,
			Description: "The template.",
			Optional:    true,
			ForceNew:    true,
		},
		"template_stack": {
			Type:        schema.TypeString,
			Description: "The template stack.",
			Optional:    true,
			ForceNew:    true,
		},
		"vsys": vsysSchema("vsys1"),
		"name": {
			Type:        schema.TypeString,
			Description: "The name.",
			Required:    true,
			ForceNew:    true,
		},
		"mode": {
			Type:         schema.TypeString,
			Description:  "The zone mode.",
			Optional:     true,
			ValidateFunc: validateStringIn("layer3", "layer2", "virtual-wire", "tap", "tunnel"),
		},
		"zone_profile": {
			Type:        schema.TypeString,
			Description: "The zone protection profile.",
			Optional:    true,
		},
		"log_setting": {
			Type:        schema.TypeString,
			Description: "Log setting.",
			Optional:    true,
		},
		"enable_user_id": {
			Type:        schema.TypeBool,
			Description: "Boolean to enable user identification.",
			Optional:    true,
		},
		"interfaces": {
			Type:        schema.TypeList,
			Description: "List of interfaces associated with this zone.  Leave this undefined if you want to use `panos_zone_entry` resources.",
			Optional:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"include_acls": {
			Type:        schema.TypeList,
			Description: "Users from these addresses/subnets will be identified.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"exclude_acls": {
			Type:        schema.TypeList,
			Description: "Users from these addresses/subnets will not be identified.",
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"enable_packet_buffer_protection": {
			Type:        schema.TypeBool,
			Description: "Boolean to enable packet buffer protection.",
			Optional:    true,
			Default:     true,
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "name", "template", "template_stack"})
	}

	return ans
}

func zoneEntrySchema(withTmpl, withTs bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys": vsysSchema("vsys1"),
		"zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"mode": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  zone.ModeL3,
			ValidateFunc: validateStringIn(
				zone.ModeL3,
				zone.ModeL2,
				zone.ModeVirtualWire,
				zone.ModeTap,
				zone.ModeExternal,
			),
			ForceNew: true,
		},
		"interface": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}

	if withTmpl {
		if withTs {
			ans["template"] = templateSchema(true)
			ans["template_stack"] = templateStackSchema()
		} else {
			ans["template"] = templateSchema(false)
		}
	}

	return ans
}

func loadZone(d *schema.ResourceData) zone.Entry {
	return zone.Entry{
		Name:                         d.Get("name").(string),
		Mode:                         d.Get("mode").(string),
		ZoneProfile:                  d.Get("zone_profile").(string),
		LogSetting:                   d.Get("log_setting").(string),
		EnableUserId:                 d.Get("enable_user_id").(bool),
		Interfaces:                   asStringList(d.Get("interfaces").([]interface{})),
		IncludeAcls:                  asStringList(d.Get("include_acls").([]interface{})),
		ExcludeAcls:                  asStringList(d.Get("exclude_acls").([]interface{})),
		EnablePacketBufferProtection: d.Get("enable_packet_buffer_protection").(bool),
	}
}

func saveZone(d *schema.ResourceData, o zone.Entry) {
	d.Set("name", o.Name)
	d.Set("mode", o.Mode)
	d.Set("zone_profile", o.ZoneProfile)
	d.Set("log_setting", o.LogSetting)
	d.Set("enable_user_id", o.EnableUserId)
	if err := d.Set("interfaces", o.Interfaces); err != nil {
		log.Printf("[WARN] Error setting 'interfaces' param for %q: %s", d.Id(), err)
	}
	if err := d.Set("include_acls", o.IncludeAcls); err != nil {
		log.Printf("[WARN] Error setting 'include_acls' param for %q: %s", d.Id(), err)
	}
	if err := d.Set("exclude_acls", o.ExcludeAcls); err != nil {
		log.Printf("[WARN] Error setting 'exclude_acls' param for %q: %s", d.Id(), err)
	}
	d.Set("enable_packet_buffer_protection", o.EnablePacketBufferProtection)
}

// Id functions.
func parseZoneId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildZoneId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parseZoneEntryId(v string) (string, string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4], t[5]
}

func buildZoneEntryId(a, b, c, d, e, f string) string {
	return strings.Join([]string{a, b, c, d, e, f}, IdSeparator)
}
