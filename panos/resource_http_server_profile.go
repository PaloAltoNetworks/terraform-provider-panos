package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/header"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/param"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/server"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceHttpServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateHttpServerProfile,
		Read:   readHttpServerProfile,
		Update: createUpdateHttpServerProfile,
		Delete: deleteHttpServerProfile,

		Schema: httpServerProfileSchema(false),
	}
}

func httpServerProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"tag_registration": {
			Type:     schema.TypeBool,
			Optional: true,
		},
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
						Type:         schema.TypeString,
						Optional:     true,
						Default:      server.ProtocolHttps,
						ValidateFunc: validateStringIn(server.ProtocolHttps, server.ProtocolHttp),
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

	folders := []string{
		"config_format",
		"system_format",
		"threat_format",
		"traffic_format",
		"hip_match_format",
		"url_format",
		"data_format",
		"wildfire_format",
		"tunnel_format",
		"user_id_format",
		"gtp_format",
		"auth_format",
		"sctp_format",
		"iptag_format",
	}

	for _, folder := range folders {
		ans[folder] = &schema.Schema{
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

	if p {
		ans["template"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template_stack",
				"device_group",
			},
		}
		ans["template_stack"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template",
				"device_group",
			},
		}
		ans["device_group"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"template",
				"template_stack",
				"vsys",
			},
		}
		ans["vsys"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ConflictsWith: []string{
				"device_group",
			},
		}
	} else {
		ans["vsys"] = &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Default:  "shared",
			ForceNew: true,
		}
	}

	return ans
}

func parseHttpServerProfile(d *schema.ResourceData) (string, http.Entry, []server.Entry, map[string][]header.Entry, map[string][]param.Entry) {
	vsys := d.Get("vsys").(string)
	o, serverList, headers, params := loadHttpServerProfile(d)

	return vsys, o, serverList, headers, params
}

func loadHttpServerProfile(d *schema.ResourceData) (http.Entry, []server.Entry, map[string][]header.Entry, map[string][]param.Entry) {
	o := http.Entry{
		Name:            d.Get("name").(string),
		TagRegistration: d.Get("tag_registration").(bool),
	}

	sl := d.Get("http_server").([]interface{})
	serverList := make([]server.Entry, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		serverList = append(serverList, server.Entry{
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

	var folder, logtype string
	headers := make(map[string][]header.Entry)
	params := make(map[string][]param.Entry)

	folder = "config_format"
	logtype = param.Config
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.ConfigName = x["name"].(string)
		o.ConfigUriFormat = x["uri_format"].(string)
		o.ConfigPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "system_format"
	logtype = param.System
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.SystemName = x["name"].(string)
		o.SystemUriFormat = x["uri_format"].(string)
		o.SystemPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "threat_format"
	logtype = param.Threat
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.ThreatName = x["name"].(string)
		o.ThreatUriFormat = x["uri_format"].(string)
		o.ThreatPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "traffic_format"
	logtype = param.Traffic
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.TrafficName = x["name"].(string)
		o.TrafficUriFormat = x["uri_format"].(string)
		o.TrafficPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "hip_match_format"
	logtype = param.HipMatch
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.HipMatchName = x["name"].(string)
		o.HipMatchUriFormat = x["uri_format"].(string)
		o.HipMatchPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "url_format"
	logtype = param.Url
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.UrlName = x["name"].(string)
		o.UrlUriFormat = x["uri_format"].(string)
		o.UrlPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "data_format"
	logtype = param.Data
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.DataName = x["name"].(string)
		o.DataUriFormat = x["uri_format"].(string)
		o.DataPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "wildfire_format"
	logtype = param.Wildfire
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.WildfireName = x["name"].(string)
		o.WildfireUriFormat = x["uri_format"].(string)
		o.WildfirePayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "tunnel_format"
	logtype = param.Tunnel
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.TunnelName = x["name"].(string)
		o.TunnelUriFormat = x["uri_format"].(string)
		o.TunnelPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "user_id_format"
	logtype = param.UserId
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.UserIdName = x["name"].(string)
		o.UserIdUriFormat = x["uri_format"].(string)
		o.UserIdPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "gtp_format"
	logtype = param.Gtp
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.GtpName = x["name"].(string)
		o.GtpUriFormat = x["uri_format"].(string)
		o.GtpPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "auth_format"
	logtype = param.Auth
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.AuthName = x["name"].(string)
		o.AuthUriFormat = x["uri_format"].(string)
		o.AuthPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "sctp_format"
	logtype = param.Sctp
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.SctpName = x["name"].(string)
		o.SctpUriFormat = x["uri_format"].(string)
		o.SctpPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	folder = "iptag_format"
	logtype = param.Iptag
	if grp := d.Get(folder).([]interface{}); grp != nil && len(grp) == 1 {
		x := grp[0].(map[string]interface{})
		o.IptagName = x["name"].(string)
		o.IptagUriFormat = x["uri_format"].(string)
		o.IptagPayload = x["payload"].(string)

		if hm := x["headers"].(map[string]interface{}); len(hm) != 0 {
			headerList := make([]header.Entry, 0, len(hm))
			for key, value := range hm {
				headerList = append(headerList, header.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			headers[logtype] = headerList
		}

		if pm := x["params"].(map[string]interface{}); len(pm) != 0 {
			paramList := make([]param.Entry, 0, len(pm))
			for key, value := range pm {
				paramList = append(paramList, param.Entry{
					Name:  key,
					Value: value.(string),
				})
			}
			params[logtype] = paramList
		}
	}

	return o, serverList, headers, params
}

func saveHttpServerProfile(d *schema.ResourceData, o http.Entry, serverList []server.Entry, headers map[string][]header.Entry, params map[string][]param.Entry) {
	d.Set("name", o.Name)
	d.Set("tag_registration", o.TagRegistration)

	pMap := d.Get("password_enc").(map[string]interface{})
	pMapRaw := d.Get("password_raw").(map[string]interface{})

	list := make([]interface{}, 0, len(serverList))
	for i := range serverList {
		var thePassword string
		if pMap[serverList[i].Name] != nil && (pMap[serverList[i].Name]).(string) == serverList[i].Password {
			thePassword = pMapRaw[serverList[i].Name].(string)
		}
		list = append(list, map[string]interface{}{
			"name":                serverList[i].Name,
			"address":             serverList[i].Address,
			"protocol":            serverList[i].Protocol,
			"port":                serverList[i].Port,
			"http_method":         serverList[i].HttpMethod,
			"username":            serverList[i].Username,
			"password":            thePassword,
			"tls_version":         serverList[i].TlsVersion,
			"certificate_profile": serverList[i].CertificateProfile,
		})
	}

	if err := d.Set("http_server", list); err != nil {
		log.Printf("[WARN] Error setting 'email_server' for %q: %s", d.Id(), err)
	}

	var folder, logtype string

	folder = "config_format"
	logtype = param.Config
	if o.ConfigName != "" || o.ConfigUriFormat != "" || o.ConfigPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.ConfigName,
			"uri_format": o.ConfigUriFormat,
			"payload":    o.ConfigPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "system_format"
	logtype = param.System
	if o.SystemName != "" || o.SystemUriFormat != "" || o.SystemPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.SystemName,
			"uri_format": o.SystemUriFormat,
			"payload":    o.SystemPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "threat_format"
	logtype = param.Threat
	if o.ThreatName != "" || o.ThreatUriFormat != "" || o.ThreatPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.ThreatName,
			"uri_format": o.ThreatUriFormat,
			"payload":    o.ThreatPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "traffic_format"
	logtype = param.Traffic
	if o.TrafficName != "" || o.TrafficUriFormat != "" || o.TrafficPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.TrafficName,
			"uri_format": o.TrafficUriFormat,
			"payload":    o.TrafficPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "hip_match_format"
	logtype = param.HipMatch
	if o.HipMatchName != "" || o.HipMatchUriFormat != "" || o.HipMatchPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.HipMatchName,
			"uri_format": o.HipMatchUriFormat,
			"payload":    o.HipMatchPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "url_format"
	logtype = param.Url
	if o.UrlName != "" || o.UrlUriFormat != "" || o.UrlPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.UrlName,
			"uri_format": o.UrlUriFormat,
			"payload":    o.UrlPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "data_format"
	logtype = param.Data
	if o.DataName != "" || o.DataUriFormat != "" || o.DataPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.DataName,
			"uri_format": o.DataUriFormat,
			"payload":    o.DataPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "wildfire_format"
	logtype = param.Wildfire
	if o.WildfireName != "" || o.WildfireUriFormat != "" || o.WildfirePayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.WildfireName,
			"uri_format": o.WildfireUriFormat,
			"payload":    o.WildfirePayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "tunnel_format"
	logtype = param.Tunnel
	if o.TunnelName != "" || o.TunnelUriFormat != "" || o.TunnelPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.TunnelName,
			"uri_format": o.TunnelUriFormat,
			"payload":    o.TunnelPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "user_id_format"
	logtype = param.UserId
	if o.UserIdName != "" || o.UserIdUriFormat != "" || o.UserIdPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.UserIdName,
			"uri_format": o.UserIdUriFormat,
			"payload":    o.UserIdPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "gtp_format"
	logtype = param.Gtp
	if o.GtpName != "" || o.GtpUriFormat != "" || o.GtpPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.GtpName,
			"uri_format": o.GtpUriFormat,
			"payload":    o.GtpPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "auth_format"
	logtype = param.Auth
	if o.AuthName != "" || o.AuthUriFormat != "" || o.AuthPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.AuthName,
			"uri_format": o.AuthUriFormat,
			"payload":    o.AuthPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "sctp_format"
	logtype = param.Sctp
	if o.SctpName != "" || o.SctpUriFormat != "" || o.SctpPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.SctpName,
			"uri_format": o.SctpUriFormat,
			"payload":    o.SctpPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}

	folder = "iptag_format"
	logtype = param.Iptag
	if o.IptagName != "" || o.IptagUriFormat != "" || o.IptagPayload != "" || len(headers[logtype]) != 0 || len(params[logtype]) != 0 {
		info := map[string]interface{}{
			"name":       o.IptagName,
			"uri_format": o.IptagUriFormat,
			"payload":    o.IptagPayload,
		}
		if headerList := headers[logtype]; len(headerList) != 0 {
			h := make(map[string]interface{})
			for _, hdr := range headerList {
				h[hdr.Name] = hdr.Value
			}
			info["headers"] = h
		}
		if paramList := params[logtype]; len(paramList) != 0 {
			p := make(map[string]interface{})
			for _, prm := range paramList {
				p[prm.Name] = prm.Value
			}
			info["params"] = p
		}

		if err := d.Set(folder, []interface{}{info}); err != nil {
			log.Printf("[WARN] Error setting %q for %q: %s", folder, d.Id(), err)
		}
	} else {
		d.Set(folder, nil)
	}
}

func parseHttpServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildHttpServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o, serverList, headers, params := parseHttpServerProfile(d)

	if err := fw.Device.HttpServerProfile.SetWithoutSubconfig(vsys, o); err != nil {
		return err
	}

	if err := fw.Device.HttpServer.Set(vsys, o.Name, serverList...); err != nil {
		return err
	}

	password_enc := make(map[string]interface{})
	password_raw := make(map[string]interface{})
	for _, x := range serverList {
		live, err := fw.Device.HttpServer.Get(vsys, o.Name, x.Name)
		if err != nil {
			return err
		}
		password_enc[x.Name] = live.Password
		password_raw[x.Name] = x.Password
	}
	if err := d.Set("password_enc", password_enc); err != nil {
		log.Printf("[WARN] Error setting 'password_enc' for %q: %s", d.Id(), err)
	}
	if err := d.Set("password_raw", password_raw); err != nil {
		log.Printf("[WARN] Error setting 'password_raw' for %q: %s", d.Id(), err)
	}

	for logtype := range headers {
		headerList := headers[logtype]
		if len(headerList) != 0 {
			if err := fw.Device.HttpHeader.Set(vsys, o.Name, logtype, headerList...); err != nil {
				return err
			}
		}
	}

	for logtype := range params {
		paramList := params[logtype]
		if len(paramList) != 0 {
			if err := fw.Device.HttpParam.Set(vsys, o.Name, logtype, paramList...); err != nil {
				return err
			}
		}
	}

	d.SetId(buildHttpServerProfileId(vsys, o.Name))
	return readHttpServerProfile(d, meta)
}

func readHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var list []string

	fw := meta.(*pango.Firewall)
	vsys, name := parseHttpServerProfileId(d.Id())

	o, err := fw.Device.HttpServerProfile.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err = fw.Device.HttpServer.GetList(vsys, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := fw.Device.HttpServer.Get(vsys, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	logtypes := []string{
		param.Config,
		param.System,
		param.Threat,
		param.Traffic,
		param.HipMatch,
		param.Url,
		param.Data,
		param.Wildfire,
		param.Tunnel,
		param.UserId,
		param.Gtp,
		param.Auth,
		param.Sctp,
		param.Iptag,
	}

	headers := make(map[string][]header.Entry)
	params := make(map[string][]param.Entry)

	for _, logtype := range logtypes {
		list, err = fw.Device.HttpHeader.GetList(vsys, name, logtype)
		if err != nil {
			return err
		}
		if len(list) != 0 {
			headerList := make([]header.Entry, 0, len(list))
			for _, hdr := range list {
				entry, err := fw.Device.HttpHeader.Get(vsys, name, logtype, hdr)
				if err != nil {
					return err
				}
				headerList = append(headerList, entry)
			}
			headers[logtype] = headerList
		}

		list, err = fw.Device.HttpParam.GetList(vsys, name, logtype)
		if err != nil {
			return err
		}
		if len(list) != 0 {
			paramList := make([]param.Entry, 0, len(list))
			for _, prm := range list {
				entry, err := fw.Device.HttpParam.Get(vsys, name, logtype, prm)
				if err != nil {
					return err
				}
				paramList = append(paramList, entry)
			}
			params[logtype] = paramList
		}
	}

	d.Set("vsys", vsys)
	saveHttpServerProfile(d, o, serverList, headers, params)

	return nil
}

func deleteHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseHttpServerProfileId(d.Id())

	err := fw.Device.HttpServerProfile.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
