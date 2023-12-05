package panos

import (
	"log"
	"strconv"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceApplicationObjects() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceApplicationObjectsRead,

		Schema: s,
	}
}

func dataSourceApplicationObjectsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.Application.GetList(id)
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
func dataSourceApplicationObject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationObjectRead,

		Schema: applicationObjectSchema(false, "", 1),
	}
}

func dataSourceApplicationObjectRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o app.Entry

	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	name := d.Get("name").(string)

	d.Set("vsys", vsys)
	d.Set("device_group", dg)

	id := buildApplicationObjectId(dg, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.Application.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.Application.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveApplicationObject(d, o)

	return nil
}

// Resource.
func resourceApplicationObject() *schema.Resource {
	return &schema.Resource{
		Create: createApplicationObject,
		Read:   readApplicationObject,
		Update: updateApplicationObject,
		Delete: deleteApplicationObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: applicationObjectSchema(true, "device_group", 0),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: applicationObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationObjectSchema(true, "", 1),
	}
}

func resourcePanoramaApplicationObject() *schema.Resource {
	return &schema.Resource{
		Create: createApplicationObject,
		Read:   readApplicationObject,
		Update: updateApplicationObject,
		Delete: deleteApplicationObject,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: (&schema.Resource{
					Schema: applicationObjectSchema(true, "vsys", 0),
				}).CoreConfigSchema().ImpliedType(),
				Upgrade: applicationObjectUpgradeV0,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationObjectSchema(true, "", 1),
	}
}

func applicationObjectUpgradeV0(raw map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if _, ok := raw["vsys"]; !ok {
		raw["vsys"] = "vsys1"
	}
	if _, ok := raw["device_group"]; !ok {
		raw["device_group"] = "shared"
	}

	if dl := raw["defaults"].([]interface{}); len(dl) > 0 {
		def := dl[0].(map[string]interface{})
		if x := asInterfaceMap(def, "ip_protocol"); len(x) > 0 {
			x["value"] = strconv.Itoa(x["value"].(int))
		}
	}

	return raw, nil
}

func createApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadApplicationObject(d)

	dg := d.Get("device_group").(string)
	vsys := d.Get("vsys").(string)

	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	id := buildApplicationObjectId(dg, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.Application.Set(vsys, o)
	case *pango.Panorama:
		err = con.Objects.Application.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readApplicationObject(d, meta)
}

func readApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o app.Entry

	// Migrate the Id.
	tok := strings.Split(d.Id(), IdSeparator)
	if len(tok) == 2 {
		switch meta.(type) {
		case *pango.Firewall:
			d.SetId(buildApplicationObjectId("shared", tok[0], tok[1]))
		case *pango.Panorama:
			d.SetId(buildApplicationObjectId(tok[0], "vsys1", tok[1]))
		}
	}

	dg, vsys, name := parseApplicationObjectId(d.Id())
	d.Set("device_group", dg)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Objects.Application.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Objects.Application.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveApplicationObject(d, o)

	return nil
}

func updateApplicationObject(d *schema.ResourceData, meta interface{}) error {
	o := loadApplicationObject(d)

	dg, vsys, _ := parseApplicationObjectId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err := con.Objects.Application.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Application.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		lo, err := con.Objects.Application.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.Application.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readApplicationObject(d, meta)
}

func deleteApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	dg, vsys, name := parseApplicationObjectId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Objects.Application.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Objects.Application.Delete(dg, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func applicationObjectSchema(isResource bool, rmKey string, schemaVersion int) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"device_group": deviceGroupSchema(),
		"vsys":         vsysSchema("vsys1"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"defaults": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"port": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						ConflictsWith: []string{
							"defaults.ip_protocol",
							"defaults.icmp",
							"defaults.icmp6",
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ports": {
									Type:     schema.TypeList,
									MinItems: 1,
									Required: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"ip_protocol": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						ConflictsWith: []string{
							"defaults.port",
							"defaults.icmp",
							"defaults.icmp6",
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"icmp": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						ConflictsWith: []string{
							"defaults.port",
							"defaults.ip_protocol",
							"defaults.icmp6",
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:     schema.TypeInt,
									Required: true,
								},
								"code": {
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},
					"icmp6": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						ConflictsWith: []string{
							"defaults.port",
							"defaults.ip_protocol",
							"defaults.icmp",
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:     schema.TypeInt,
									Required: true,
								},
								"code": {
									Type:     schema.TypeInt,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"category": {
			Type:     schema.TypeString,
			Required: true,
		},
		"subcategory": {
			Type:     schema.TypeString,
			Required: true,
		},
		"technology": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"timeout_settings": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"timeout": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"tcp_timeout": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"udp_timeout": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"tcp_half_closed": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"tcp_time_wait": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
		"risk": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  1,
		},
		"parent_app": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"able_to_file_transfer": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"excessive_bandwidth": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"tunnels_other_applications": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"has_known_vulnerability": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"used_by_malware": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"evasive_behavior": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"pervasive_use": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"prone_to_misuse": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"continue_scanning_for_other_applications": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"scanning": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"file_types": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"viruses": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"data_patterns": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"alg_disable_capability": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"no_app_id_caching": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	if rmKey != "" {
		delete(ans, rmKey)
	}

	if schemaVersion == 0 {
		var x *schema.Resource
		x = ans["defaults"].Elem.(*schema.Resource)
		x = x.Schema["ip_protocol"].Elem.(*schema.Resource)
		x.Schema["value"].Type = schema.TypeInt
	}

	return ans
}

func loadApplicationObject(d *schema.ResourceData) app.Entry {
	ans := app.Entry{
		Name:                                 d.Get("name").(string),
		Category:                             d.Get("category").(string),
		Subcategory:                          d.Get("subcategory").(string),
		Technology:                           d.Get("technology").(string),
		Description:                          d.Get("description").(string),
		Risk:                                 d.Get("risk").(int),
		AbleToFileTransfer:                   d.Get("able_to_file_transfer").(bool),
		ExcessiveBandwidth:                   d.Get("excessive_bandwidth").(bool),
		TunnelsOtherApplications:             d.Get("tunnels_other_applications").(bool),
		HasKnownVulnerability:                d.Get("has_known_vulnerability").(bool),
		UsedByMalware:                        d.Get("used_by_malware").(bool),
		EvasiveBehavior:                      d.Get("evasive_behavior").(bool),
		PervasiveUse:                         d.Get("pervasive_use").(bool),
		ProneToMisuse:                        d.Get("prone_to_misuse").(bool),
		ContinueScanningForOtherApplications: d.Get("continue_scanning_for_other_applications").(bool),
		AlgDisableCapability:                 d.Get("alg_disable_capability").(string),
		NoAppIdCaching:                       d.Get("no_app_id_caching").(bool),
	}

	if dl := d.Get("defaults").([]interface{}); len(dl) > 0 {
		if dl[0] == nil {
			ans.DefaultType = app.DefaultTypeNone
		} else {
			def := dl[0].(map[string]interface{})
			if x := asInterfaceMap(def, "port"); len(x) > 0 {
				ans.DefaultType = app.DefaultTypePort
				ans.DefaultPorts = asStringList(x["ports"].([]interface{}))
			} else if x := asInterfaceMap(def, "ip_protocol"); len(x) > 0 {
				ans.DefaultType = app.DefaultTypeIpProtocol
				ans.DefaultIpProtocol = x["value"].(string)
			} else if x := asInterfaceMap(def, "icmp"); len(x) > 0 {
				ans.DefaultType = app.DefaultTypeIcmp
				ans.DefaultIcmpType = x["type"].(int)
				ans.DefaultIcmpCode = x["code"].(int)
			} else if x := asInterfaceMap(def, "icmp6"); len(x) > 0 {
				ans.DefaultType = app.DefaultTypeIcmp6
				ans.DefaultIcmpType = x["type"].(int)
				ans.DefaultIcmpCode = x["code"].(int)
			}
		}
	} else {
		ans.DefaultType = app.DefaultTypeNone
	}

	if tl := d.Get("timeout_settings").([]interface{}); len(tl) > 0 {
		to := tl[0].(map[string]interface{})
		ans.Timeout = to["timeout"].(int)
		ans.TcpTimeout = to["tcp_timeout"].(int)
		ans.UdpTimeout = to["udp_timeout"].(int)
		ans.TcpHalfClosedTimeout = to["tcp_half_closed"].(int)
		ans.TcpTimeWaitTimeout = to["tcp_time_wait"].(int)
	}

	if sl := d.Get("scanning").([]interface{}); len(sl) > 0 {
		scan := sl[0].(map[string]interface{})
		ans.FileTypeIdent = scan["file_types"].(bool)
		ans.VirusIdent = scan["viruses"].(bool)
		ans.DataIdent = scan["data_patterns"].(bool)
	}

	return ans
}

func saveApplicationObject(d *schema.ResourceData, o app.Entry) {
	d.Set("name", o.Name)
	d.Set("category", o.Category)
	d.Set("subcategory", o.Subcategory)
	d.Set("technology", o.Technology)
	d.Set("description", o.Description)
	d.Set("risk", o.Risk)
	d.Set("able_to_file_transfer", o.AbleToFileTransfer)
	d.Set("excessive_bandwidth", o.ExcessiveBandwidth)
	d.Set("tunnels_other_applications", o.TunnelsOtherApplications)
	d.Set("has_known_vulnerability", o.HasKnownVulnerability)
	d.Set("used_by_malware", o.UsedByMalware)
	d.Set("evasive_behavior", o.EvasiveBehavior)
	d.Set("pervasive_use", o.PervasiveUse)
	d.Set("prone_to_misuse", o.ProneToMisuse)
	d.Set("continue_scanning_for_other_applications", o.ContinueScanningForOtherApplications)
	d.Set("alg_disable_capability", o.AlgDisableCapability)
	d.Set("no_app_id_caching", o.NoAppIdCaching)

	switch o.DefaultType {
	case app.DefaultTypeNone:
		d.Set("defaults", nil)
	case app.DefaultTypePort:
		def := []interface{}{
			map[string]interface{}{
				"port": []interface{}{
					map[string]interface{}{
						"ports": o.DefaultPorts,
					},
				},
			},
		}
		if err := d.Set("defaults", def); err != nil {
			log.Printf("[WARN] Error setting 'port.defaults' for %q: %s", d.Id(), err)
		}
	case app.DefaultTypeIpProtocol:
		def := []interface{}{
			map[string]interface{}{
				"ip_protocol": []interface{}{
					map[string]interface{}{
						"value": o.DefaultIpProtocol,
					},
				},
			},
		}
		if err := d.Set("defaults", def); err != nil {
			log.Printf("[WARN] Error setting 'ip_protocol.defaults' for %q: %s", d.Id(), err)
		}
	case app.DefaultTypeIcmp:
		def := []interface{}{
			map[string]interface{}{
				"icmp": []interface{}{
					map[string]interface{}{
						"type": o.DefaultIcmpType,
						"code": o.DefaultIcmpCode,
					},
				},
			},
		}
		if err := d.Set("defaults", def); err != nil {
			log.Printf("[WARN] Error setting 'icmp.defaults' for %q: %s", d.Id(), err)
		}
	case app.DefaultTypeIcmp6:
		def := []interface{}{
			map[string]interface{}{
				"icmp6": []interface{}{
					map[string]interface{}{
						"type": o.DefaultIcmpType,
						"code": o.DefaultIcmpCode,
					},
				},
			},
		}
		if err := d.Set("defaults", def); err != nil {
			log.Printf("[WARN] Error setting 'icmp6.defaults' for %q: %s", d.Id(), err)
		}
	}

	if o.Timeout != 0 || o.TcpTimeout != 0 || o.UdpTimeout != 0 || o.TcpHalfClosedTimeout != 0 || o.TcpTimeWaitTimeout != 0 {
		to := []interface{}{
			map[string]interface{}{
				"timeout":         o.Timeout,
				"tcp_timeout":     o.TcpTimeout,
				"udp_timeout":     o.UdpTimeout,
				"tcp_half_closed": o.TcpHalfClosedTimeout,
				"tcp_time_wait":   o.TcpTimeWaitTimeout,
			},
		}
		if err := d.Set("timeout_settings", to); err != nil {
			log.Printf("[WARN] Error setting 'timeout_settings' for %q: %s", d.Id(), err)
		}
	}

	if o.FileTypeIdent || o.VirusIdent || o.DataIdent {
		scan := []interface{}{
			map[string]interface{}{
				"file_types":    o.FileTypeIdent,
				"viruses":       o.VirusIdent,
				"data_patterns": o.DataIdent,
			},
		}
		if err := d.Set("scanning", scan); err != nil {
			log.Printf("[WARN] Error setting 'scanning' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func parseApplicationObjectId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildApplicationObjectId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}
