package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApplicationObject() *schema.Resource {
	return &schema.Resource{
		Create: createApplicationObject,
		Read:   readApplicationObject,
		Update: updateApplicationObject,
		Delete: deleteApplicationObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationObjectSchema(false),
	}
}

func applicationObjectSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
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
									Type:     schema.TypeInt,
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

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema()
	}

	return ans
}

func parseApplicationObject(d *schema.ResourceData) (string, app.Entry) {
	vsys := d.Get("vsys").(string)
	o := loadApplicationObject(d)

	return vsys, o
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
				ans.DefaultIpProtocol = x["value"].(int)
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

func parseApplicationObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildApplicationObjectId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
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

func createApplicationObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseApplicationObject(d)

	if err := fw.Objects.Application.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildApplicationObjectId(vsys, o.Name))
	return readApplicationObject(d, meta)
}

func readApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseApplicationObjectId(d.Id())

	o, err := fw.Objects.Application.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vsys", vsys)
	saveApplicationObject(d, o)

	return nil
}

func updateApplicationObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseApplicationObject(d)

	lo, err := fw.Objects.Application.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.Application.Edit(vsys, lo); err != nil {
		return err
	}

	return readApplicationObject(d, meta)
}

func deleteApplicationObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseApplicationObjectId(d.Id())

	err := fw.Objects.Application.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
