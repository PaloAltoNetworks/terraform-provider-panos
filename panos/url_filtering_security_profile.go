package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceUrlFilteringSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceUrlFilteringSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceUrlFilteringSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.UrlFilteringProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.UrlFilteringProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceUrlFilteringSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUrlFilteringSecurityProfileRead,

		Schema: urlFilteringSecurityProfileSchema(false),
	}
}

func dataSourceUrlFilteringSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o url.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildUrlFilteringSecurityProfileId(vsys, name)
		o, err = con.Objects.UrlFilteringProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildUrlFilteringSecurityProfileId(dg, name)
		o, err = con.Objects.UrlFilteringProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveUrlFilteringSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceUrlFilteringSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUrlFilteringSecurityProfile,
		Read:   readUrlFilteringSecurityProfile,
		Update: updateUrlFilteringSecurityProfile,
		Delete: deleteUrlFilteringSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: urlFilteringSecurityProfileSchema(true),
	}
}

func createUrlFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadUrlFilteringSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildUrlFilteringSecurityProfileId(vsys, o.Name)
		err = con.Objects.UrlFilteringProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildUrlFilteringSecurityProfileId(dg, o.Name)
		err = con.Objects.UrlFilteringProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readUrlFilteringSecurityProfile(d, meta)
}

func readUrlFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o url.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseUrlFilteringSecurityProfileId(d.Id())
		o, err = con.Objects.UrlFilteringProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseUrlFilteringSecurityProfileId(d.Id())
		o, err = con.Objects.UrlFilteringProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveUrlFilteringSecurityProfile(d, o)
	return nil
}

func updateUrlFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadUrlFilteringSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.UrlFilteringProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.UrlFilteringProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.UrlFilteringProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.UrlFilteringProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readUrlFilteringSecurityProfile(d, meta)
}

func deleteUrlFilteringSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseUrlFilteringSecurityProfileId(d.Id())
		err = con.Objects.UrlFilteringProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseUrlFilteringSecurityProfileId(d.Id())
		err = con.Objects.UrlFilteringProfile.Delete(dg, name)
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
func urlFilteringSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
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
		"dynamic_url": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "(Removed in PAN-OS 9.0) Dynamic URL",
		},
		"expired_license_action": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "(Removed in PAN-OS 9.0) Expired license action",
		},
		"block_list_action": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(Removed in PAN-OS 9.0) Block list action",
		},
		"block_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(Removed in PAN-OS 9.0) Block list",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"allow_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(Removed in PAN-OS 9.0) Allow list",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"allow_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of categories to allow",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"alert_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of categories to alert",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"block_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of categories to block",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"continue_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of categories to continue",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"override_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of categories to allow",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"track_container_page": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Set to true to track the container page",
		},
		"log_container_page_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Set to true to log container page only",
		},
		"safe_search_enforcement": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Set to true for safe search enforcement",
		},
		"log_http_header_xff": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "HTTP Header Logging: X-Forwarded-For",
		},
		"log_http_header_user_agent": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "HTTP Header Logging: User-Agent",
		},
		"log_http_header_referer": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "HTTP Header Logging: Referer",
		},
		"ucd_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 8.0+) User credential detection mode",
			Default:     url.UcdModeDisabled,
			ValidateFunc: validateStringIn(
				url.UcdModeDisabled,
				url.UcdModeIpUser,
				url.UcdModeDomainCredentials,
				url.UcdModeGroupMapping,
			),
		},
		"ucd_mode_group_mapping": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: fmt.Sprintf("(PAN-OS 8.0+, ucd_mode = %s) User credential detection: the group mapping settings", url.UcdModeGroupMapping),
		},
		"ucd_log_severity": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 8.0+) User credential detection: valid username detected log severity",
		},
		"ucd_allow_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 8.0+) Categories allowed with user credential submission",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ucd_alert_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 8.0+) Categories alerted on with user credential submission",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ucd_block_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 8.0+) Categories blocked with user credential submission",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ucd_continue_categories": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 8.0+) Categories continued with user credential submission",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"http_header_insertion": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of http header specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Header name",
					},
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Header type",
					},
					"domains": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Header domains",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"http_header": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "List of HTTP header specs",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:        schema.TypeString,
									Computed:    true,
									Description: "HTTP header name (auto-generated)",
								},
								"header": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "The header",
								},
								"value": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "The value of the header",
								},
								"log": {
									Type:        schema.TypeBool,
									Optional:    true,
									Description: "Set to true to enable logging of this header insertion",
								},
							},
						},
					},
				},
			},
		},
		"machine_learning_model": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 10.0+) List of machine learning model specs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"model": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Machine learning model",
					},
					"action": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Machine learning model action",
					},
				},
			},
		},
		"machine_learning_exceptions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 10.0+) List of machine learning exceptions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	return ans
}

func loadUrlFilteringSecurityProfile(d *schema.ResourceData) url.Entry {
	var list []interface{}

	var hhis []url.HttpHeaderInsertion
	list = d.Get("http_header_insertion").([]interface{})
	if len(list) > 0 {
		hhis = make([]url.HttpHeaderInsertion, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})

			var hhs []url.HttpHeader
			l2 := elm["http_header"].([]interface{})
			if len(l2) > 0 {
				hhs = make([]url.HttpHeader, 0, len(l2))
				for j := range l2 {
					helm := l2[j].(map[string]interface{})
					hhs = append(hhs, url.HttpHeader{
						Header: helm["header"].(string),
						Value:  helm["value"].(string),
						Log:    helm["log"].(bool),
					})
				}
			}

			hhis = append(hhis, url.HttpHeaderInsertion{
				Name:        elm["name"].(string),
				Type:        elm["type"].(string),
				Domains:     asStringList(elm["domains"].([]interface{})),
				HttpHeaders: hhs,
			})
		}
	}

	var mlms []url.MachineLearningModel
	list = d.Get("machine_learning_model").([]interface{})
	if len(list) > 0 {
		mlms = make([]url.MachineLearningModel, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			mlms = append(mlms, url.MachineLearningModel{
				Model:  elm["model"].(string),
				Action: elm["action"].(string),
			})
		}
	}

	return url.Entry{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		DynamicUrl:                d.Get("dynamic_url").(bool),
		ExpiredLicenseAction:      d.Get("expired_license_action").(bool),
		BlockListAction:           d.Get("block_list_action").(string),
		BlockList:                 asStringList(d.Get("block_list").([]interface{})),
		AllowList:                 asStringList(d.Get("allow_list").([]interface{})),
		AllowCategories:           asStringList(d.Get("allow_categories").([]interface{})),
		AlertCategories:           asStringList(d.Get("alert_categories").([]interface{})),
		BlockCategories:           asStringList(d.Get("block_categories").([]interface{})),
		ContinueCategories:        asStringList(d.Get("continue_categories").([]interface{})),
		OverrideCategories:        asStringList(d.Get("override_categories").([]interface{})),
		TrackContainerPage:        d.Get("track_container_page").(bool),
		LogContainerPageOnly:      d.Get("log_container_page_only").(bool),
		SafeSearchEnforcement:     d.Get("safe_search_enforcement").(bool),
		LogHttpHeaderXff:          d.Get("log_http_header_xff").(bool),
		LogHttpHeaderUserAgent:    d.Get("log_http_header_user_agent").(bool),
		LogHttpHeaderReferer:      d.Get("log_http_header_referer").(bool),
		UcdMode:                   d.Get("ucd_mode").(string),
		UcdModeGroupMapping:       d.Get("ucd_mode_group_mapping").(string),
		UcdLogSeverity:            d.Get("ucd_log_severity").(string),
		UcdAllowCategories:        asStringList(d.Get("ucd_allow_categories").([]interface{})),
		UcdAlertCategories:        asStringList(d.Get("ucd_alert_categories").([]interface{})),
		UcdBlockCategories:        asStringList(d.Get("ucd_block_categories").([]interface{})),
		UcdContinueCategories:     asStringList(d.Get("ucd_continue_categories").([]interface{})),
		HttpHeaderInsertions:      hhis,
		MachineLearningModels:     mlms,
		MachineLearningExceptions: asStringList(d.Get("machine_learning_exceptions").([]interface{})),
	}
}

func saveUrlFilteringSecurityProfile(d *schema.ResourceData, o url.Entry) {
	var err error

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("dynamic_url", o.DynamicUrl)
	d.Set("expired_license_action", o.ExpiredLicenseAction)
	d.Set("block_list_action", o.BlockListAction)
	if err = d.Set("block_list", o.BlockList); err != nil {
		log.Printf("[WARN] Error setting 'block_list' for %q: %s", d.Id(), err)
	}
	if err = d.Set("allow_list", o.AllowList); err != nil {
		log.Printf("[WARN] Error setting 'allow_list' for %q: %s", d.Id(), err)
	}
	if err = d.Set("allow_categories", o.AllowCategories); err != nil {
		log.Printf("[WARN] Error setting 'allow_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("alert_categories", o.AlertCategories); err != nil {
		log.Printf("[WARN] Error setting 'alert_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("block_categories", o.BlockCategories); err != nil {
		log.Printf("[WARN] Error setting 'block_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("continue_categories", o.ContinueCategories); err != nil {
		log.Printf("[WARN] Error setting 'continue_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("override_categories", o.OverrideCategories); err != nil {
		log.Printf("[WARN] Error setting 'override_categories' for %q: %s", d.Id(), err)
	}
	d.Set("track_container_page", o.TrackContainerPage)
	d.Set("log_container_page_only", o.LogContainerPageOnly)
	d.Set("safe_search_enforcement", o.SafeSearchEnforcement)
	d.Set("log_http_header_xff", o.LogHttpHeaderXff)
	d.Set("log_http_header_user_agent", o.LogHttpHeaderUserAgent)
	d.Set("log_http_header_referer", o.LogHttpHeaderReferer)
	d.Set("ucd_mode", o.UcdMode)
	d.Set("ucd_mode_group_mapping", o.UcdModeGroupMapping)
	d.Set("ucd_log_severity", o.UcdLogSeverity)
	if err = d.Set("ucd_allow_categories", o.UcdAllowCategories); err != nil {
		log.Printf("[WARN] Error setting 'ucd_allow_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ucd_alert_categories", o.UcdAlertCategories); err != nil {
		log.Printf("[WARN] Error setting 'ucd_alert_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ucd_block_categories", o.UcdBlockCategories); err != nil {
		log.Printf("[WARN] Error setting 'ucd_block_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("ucd_continue_categories", o.UcdContinueCategories); err != nil {
		log.Printf("[WARN] Error setting 'ucd_continue_categories' for %q: %s", d.Id(), err)
	}
	if err = d.Set("machine_learning_exceptions", o.MachineLearningExceptions); err != nil {
		log.Printf("[WARN] Error setting 'machine_learning_exceptions' for %q: %s", d.Id(), err)
	}

	if len(o.HttpHeaderInsertions) == 0 {
		d.Set("http_header_insertion", nil)
	} else {
		list := make([]interface{}, 0, len(o.HttpHeaderInsertions))
		for _, x := range o.HttpHeaderInsertions {
			item := map[string]interface{}{
				"name":    x.Name,
				"type":    x.Type,
				"domains": x.Domains,
			}
			if len(x.HttpHeaders) == 0 {
				item["http_header"] = nil
			} else {
				hlist := make([]interface{}, 0, len(x.HttpHeaders))
				for _, hdr := range x.HttpHeaders {
					hlist = append(hlist, map[string]interface{}{
						"name":   hdr.Name,
						"header": hdr.Header,
						"value":  hdr.Value,
						"log":    hdr.Log,
					})
				}
				item["http_header"] = hlist
			}
			list = append(list, item)
		}
		if err = d.Set("http_header_insertion", list); err != nil {
			log.Printf("[WARN] Error setting 'http_header_insertion' for %q: %s", d.Id(), err)
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
		if err = d.Set("machine_learning_model", list); err != nil {
			log.Printf("[WARN] Error setting 'machine_learning_model' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildUrlFilteringSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseUrlFilteringSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
