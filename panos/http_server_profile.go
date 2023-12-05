package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource.
func resourceHttpServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createHttpServerProfile,
		Read:   readHttpServerProfile,
		Update: updateHttpServerProfile,
		Delete: deleteHttpServerProfile,

		Schema: httpServerProfileSchema(true, "shared", []string{"device_group", "template", "template_stack"}),
	}
}

func resourcePanoramaHttpServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createHttpServerProfile,
		Read:   readHttpServerProfile,
		Update: updateHttpServerProfile,
		Delete: deleteHttpServerProfile,

		Schema: httpServerProfileSchema(true, "", nil),
	}
}

func createHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadHttpServerProfile(d)
	var lo http.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildHttpServerProfileId(vsys, o.Name)
		if err = con.Device.HttpServerProfile.Set(vsys, o); err != nil {
			return err
		}
		lo, err = con.Device.HttpServerProfile.Get(vsys, o.Name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		id = buildPanoramaHttpServerProfileId(tmpl, ts, vsys, o.Name)
		if err = con.Device.HttpServerProfile.Set(tmpl, ts, vsys, o); err != nil {
			return err
		}
		lo, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, o.Name)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	if err = saveHttpServerPasswords(d, o.Servers, lo.Servers); err != nil {
		return err
	}

	return readHttpServerProfile(d, meta)
}

func readHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o http.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseHttpServerProfileId(d.Id())
		d.Set("vsys", vsys)
		o, err = con.Device.HttpServerProfile.Get(vsys, name)
	case *pango.Panorama:
		// If this is the old style ID, it will have an extra field.  So
		// we need to migrate the ID first.
		tok := strings.Split(d.Id(), IdSeparator)
		if len(tok) == 5 {
			d.SetId(buildPanoramaHttpServerProfileId(tok[0], tok[1], tok[2], tok[4]))
		}
		// Continue on as normal.
		tmpl, ts, vsys, name := parsePanoramaHttpServerProfileId(d.Id())
		d.Set("template", tmpl)
		d.Set("template_stack", ts)
		d.Set("vsys", vsys)
		d.Set("device_group", d.Get("device_group").(string))
		o, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveHttpServerProfile(d, o)

	return nil
}

func updateHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadHttpServerProfile(d)
	var err error
	var lo http.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err = con.Device.HttpServerProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.HttpServerProfile.Edit(vsys, lo); err != nil {
			return err
		}
		lo, err = con.Device.HttpServerProfile.Get(vsys, lo.Name)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		lo, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.HttpServerProfile.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
		lo, err = con.Device.HttpServerProfile.Get(tmpl, ts, vsys, lo.Name)
	}

	if err != nil {
		return err
	}

	if err = saveHttpServerPasswords(d, o.Servers, lo.Servers); err != nil {
		return err
	}

	return readHttpServerProfile(d, meta)
}

func deleteHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseHttpServerProfileId(d.Id())
		err = con.Device.HttpServerProfile.Delete(vsys, name)
	case *pango.Panorama:
		tmpl, ts, vsys, name := parsePanoramaHttpServerProfileId(d.Id())
		err = con.Device.HttpServerProfile.Delete(tmpl, ts, vsys, name)
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
func httpServerProfileSchema(isResource bool, vsysDefault string, rmKeys []string) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"template_stack"},
		},
		"template_stack": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"template"},
		},
		"device_group": {
			Type:       schema.TypeString,
			Optional:   true,
			Deprecated: "This parameter is not applicable to this resource and will be removed later.",
		},
		"vsys": vsysSchema(vsysDefault),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"tag_registration": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"config_format":    httpServerProfilePayloadSchema(),
		"system_format":    httpServerProfilePayloadSchema(),
		"threat_format":    httpServerProfilePayloadSchema(),
		"traffic_format":   httpServerProfilePayloadSchema(),
		"hip_match_format": httpServerProfilePayloadSchema(),
		"url_format":       httpServerProfilePayloadSchema(),
		"data_format":      httpServerProfilePayloadSchema(),
		"wildfire_format":  httpServerProfilePayloadSchema(),
		"tunnel_format":    httpServerProfilePayloadSchema(),
		"user_id_format":   httpServerProfilePayloadSchema(),
		"gtp_format":       httpServerProfilePayloadSchema(),
		"auth_format":      httpServerProfilePayloadSchema(),
		"sctp_format":      httpServerProfilePayloadSchema(),
		"iptag_format":     httpServerProfilePayloadSchema(),
		"http_server": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"address": {
						Type:     schema.TypeString,
						Required: true,
					},
					"protocol": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  http.ProtocolHttps,
						ValidateFunc: validateStringIn(
							http.ProtocolHttps,
							http.ProtocolHttp,
						),
					},
					"port": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  443,
					},
					"http_method": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  "POST",
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
					"tls_version": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"certificate_profile": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"password_enc": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"password_raw": {
			Type:      schema.TypeMap,
			Computed:  true,
			Sensitive: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "device_group", "name"})
	}

	return ans
}

func httpServerProfilePayloadSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"uri_format": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"payload": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"headers": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"params": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func loadHttpServerProfile(d *schema.ResourceData) http.Entry {
	sl := d.Get("http_server").([]interface{})
	serverList := make([]http.Server, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		serverList = append(serverList, http.Server{
			Name:               x["name"].(string),
			Address:            x["address"].(string),
			Protocol:           x["protocol"].(string),
			Port:               x["port"].(int),
			HttpMethod:         x["http_method"].(string),
			Username:           x["username"].(string),
			Password:           x["password"].(string),
			TlsVersion:         x["tls_version"].(string),
			CertificateProfile: x["certificate_profile"].(string),
		})
	}

	return http.Entry{
		Name:            d.Get("name").(string),
		TagRegistration: d.Get("tag_registration").(bool),
		Servers:         serverList,
		Config:          loadHttpServerProfilePayload(d, "config_format"),
		System:          loadHttpServerProfilePayload(d, "system_format"),
		Threat:          loadHttpServerProfilePayload(d, "threat_format"),
		Traffic:         loadHttpServerProfilePayload(d, "traffic_format"),
		HipMatch:        loadHttpServerProfilePayload(d, "hip_match_format"),
		Url:             loadHttpServerProfilePayload(d, "url_format"),
		Data:            loadHttpServerProfilePayload(d, "data_format"),
		Wildfire:        loadHttpServerProfilePayload(d, "wildfire_format"),
		Tunnel:          loadHttpServerProfilePayload(d, "tunnel_format"),
		UserId:          loadHttpServerProfilePayload(d, "user_id_format"),
		Gtp:             loadHttpServerProfilePayload(d, "gtp_format"),
		Auth:            loadHttpServerProfilePayload(d, "auth_format"),
		Sctp:            loadHttpServerProfilePayload(d, "sctp_format"),
		Iptag:           loadHttpServerProfilePayload(d, "iptag_format"),
	}
}

func loadHttpServerProfilePayload(d *schema.ResourceData, folder string) *http.PayloadFormat {
	grp := d.Get(folder).([]interface{})
	if grp == nil || len(grp) == 0 {
		return nil
	}

	x := grp[0].(map[string]interface{})
	var headers []http.Header
	var params []http.Parameter

	if hm := x["headers"].(map[string]interface{}); len(hm) > 0 {
		headers = make([]http.Header, 0, len(hm))
		for key, value := range hm {
			headers = append(headers, http.Header{
				Name:  key,
				Value: value.(string),
			})
		}
	}

	if pm := x["params"].(map[string]interface{}); len(pm) > 0 {
		params = make([]http.Parameter, 0, len(pm))
		for key, value := range pm {
			params = append(params, http.Parameter{
				Name:  key,
				Value: value.(string),
			})
		}
	}

	return &http.PayloadFormat{
		Name:       x["name"].(string),
		UriFormat:  x["uri_format"].(string),
		Payload:    x["payload"].(string),
		Headers:    headers,
		Parameters: params,
	}
}

func saveHttpServerPasswords(d *schema.ResourceData, unencrypted, encrypted []http.Server) error {
	password_enc := make(map[string]interface{})
	password_raw := make(map[string]interface{})

	if len(unencrypted) != len(encrypted) {
		return fmt.Errorf("unencrypted len:%d vs encrypted len:%d", len(unencrypted), len(encrypted))
	}

	for i := range unencrypted {
		if unencrypted[i].Name != encrypted[i].Name {
			return fmt.Errorf("Name mismatch at index %d: config:%q vs list:%q", i, unencrypted[i].Name, encrypted[i].Name)
		}
		password_raw[unencrypted[i].Name] = unencrypted[i].Password
		password_enc[encrypted[i].Name] = encrypted[i].Password
	}

	if err := d.Set("password_enc", password_enc); err != nil {
		log.Printf("[WARN] Error setting 'password_enc' for %q: %s", d.Id(), err)
	}
	if err := d.Set("password_raw", password_raw); err != nil {
		log.Printf("[WARN] Error setting 'password_raw' for %q: %s", d.Id(), err)
	}

	return nil
}

func saveHttpServerProfile(d *schema.ResourceData, o http.Entry) {
	d.Set("name", o.Name)
	d.Set("tag_registration", o.TagRegistration)

	pMap := d.Get("password_enc").(map[string]interface{})
	pMapRaw := d.Get("password_raw").(map[string]interface{})

	servers := make([]interface{}, 0, len(o.Servers))
	for _, x := range o.Servers {
		var thePassword string
		if pMap[x.Name] != nil && (pMap[x.Name]).(string) == x.Password {
			thePassword = pMapRaw[x.Name].(string)
		}
		servers = append(servers, map[string]interface{}{
			"name":                x.Name,
			"address":             x.Address,
			"protocol":            x.Protocol,
			"port":                x.Port,
			"http_method":         x.HttpMethod,
			"username":            x.Username,
			"password":            thePassword,
			"tls_version":         x.TlsVersion,
			"certificate_profile": x.CertificateProfile,
		})
	}
	if err := d.Set("http_server", servers); err != nil {
		log.Printf("[WARN] Error setting \"http_server\" for %q: %s", d.Id(), err)
	}

	config := []struct {
		location string
		val      *http.PayloadFormat
	}{
		{"config_format", o.Config},
		{"system_format", o.System},
		{"threat_format", o.Threat},
		{"traffic_format", o.Traffic},
		{"hip_match_format", o.HipMatch},
		{"url_format", o.Url},
		{"data_format", o.Data},
		{"wildfire_format", o.Wildfire},
		{"tunnel_format", o.Tunnel},
		{"user_id_format", o.UserId},
		{"gtp_format", o.Gtp},
		{"auth_format", o.Auth},
		{"sctp_format", o.Sctp},
		{"iptag_format", o.Iptag},
	}

	for _, x := range config {
		if err := d.Set(x.location, flattenHttpServerProfilePayload(x.val)); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", x.location, d.Id(), err)
		}
	}
}

func flattenHttpServerProfilePayload(val *http.PayloadFormat) []interface{} {
	if val == nil {
		return nil
	}

	info := map[string]interface{}{
		"name":       val.Name,
		"uri_format": val.UriFormat,
		"payload":    val.Payload,
	}

	if len(val.Headers) > 0 {
		m := make(map[string]interface{})
		for _, x := range val.Headers {
			m[x.Name] = x.Value
		}
		info["headers"] = m
	} else {
		info["headers"] = nil
	}

	if len(val.Parameters) > 0 {
		m := make(map[string]interface{})
		for _, x := range val.Parameters {
			m[x.Name] = x.Value
		}
		info["params"] = m
	} else {
		info["params"] = nil
	}

	return []interface{}{info}
}

// Id functions.
func parseHttpServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaHttpServerProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildHttpServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaHttpServerProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
