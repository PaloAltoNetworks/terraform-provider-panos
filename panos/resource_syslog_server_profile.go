package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSyslogServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSyslogServerProfile,
		Read:   readSyslogServerProfile,
		Update: createUpdateSyslogServerProfile,
		Delete: deleteSyslogServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: syslogServerProfileSchema(false),
	}
}

func syslogServerProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"config_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"system_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"threat_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"traffic_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"hip_match_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"url_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"data_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"wildfire_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tunnel_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_id_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"gtp_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"auth_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"sctp_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"iptag_format": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"escaped_characters": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"escape_character": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"syslog_server": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"server": {
						Type:     schema.TypeString,
						Required: true,
					},
					"transport": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  server.TransportUdp,
						ValidateFunc: validateStringIn(
							server.TransportUdp,
							server.TransportTcp,
							server.TransportSsl,
						),
					},
					"port": {
						Type:     schema.TypeInt,
						Optional: true,
						Default:  514,
					},
					"syslog_format": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  server.SyslogFormatBsd,
						ValidateFunc: validateStringIn(
							server.SyslogFormatBsd,
							server.SyslogFormatIetf,
						),
					},
					"facility": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  server.FacilityUser,
						ValidateFunc: validateStringIn(
							server.FacilityUser,
							server.FacilityLocal0,
							server.FacilityLocal1,
							server.FacilityLocal2,
							server.FacilityLocal3,
							server.FacilityLocal4,
							server.FacilityLocal5,
							server.FacilityLocal6,
							server.FacilityLocal7,
						),
					},
				},
			},
		},
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

func parseSyslogServerProfile(d *schema.ResourceData) (string, syslog.Entry, []server.Entry) {
	vsys := d.Get("vsys").(string)
	o, serverList := loadSyslogServerProfile(d)

	return vsys, o, serverList
}

func loadSyslogServerProfile(d *schema.ResourceData) (syslog.Entry, []server.Entry) {
	o := syslog.Entry{
		Name:              d.Get("name").(string),
		Config:            d.Get("config_format").(string),
		System:            d.Get("system_format").(string),
		Threat:            d.Get("threat_format").(string),
		Traffic:           d.Get("traffic_format").(string),
		HipMatch:          d.Get("hip_match_format").(string),
		Url:               d.Get("url_format").(string),
		Data:              d.Get("data_format").(string),
		Wildfire:          d.Get("wildfire_format").(string),
		Tunnel:            d.Get("tunnel_format").(string),
		Gtp:               d.Get("gtp_format").(string),
		Auth:              d.Get("auth_format").(string),
		Sctp:              d.Get("sctp_format").(string),
		Iptag:             d.Get("iptag_format").(string),
		EscapedCharacters: d.Get("escaped_characters").(string),
		EscapeCharacter:   d.Get("escape_character").(string),
	}

	sl := d.Get("syslog_server").([]interface{})
	serverList := make([]server.Entry, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		serverList = append(serverList, server.Entry{
			Name:         x["name"].(string),
			Server:       x["server"].(string),
			Transport:    x["transport"].(string),
			Port:         x["port"].(int),
			SyslogFormat: x["syslog_format"].(string),
			Facility:     x["facility"].(string),
		})
	}

	return o, serverList
}

func saveSyslogServerProfile(d *schema.ResourceData, o syslog.Entry, serverList []server.Entry) {
	d.Set("name", o.Name)
	d.Set("config_format", o.Config)
	d.Set("system_format", o.System)
	d.Set("threat_format", o.Threat)
	d.Set("traffic_format", o.Traffic)
	d.Set("hip_match_format", o.HipMatch)
	d.Set("url_format", o.Url)
	d.Set("data_format", o.Data)
	d.Set("wildfire_format", o.Wildfire)
	d.Set("tunnel_format", o.Tunnel)
	d.Set("user_id_format", o.UserId)
	d.Set("gtp_format", o.Gtp)
	d.Set("auth_format", o.Auth)
	d.Set("sctp_format", o.Sctp)
	d.Set("iptag_format", o.Iptag)
	d.Set("escaped_characters", o.EscapedCharacters)
	d.Set("escape_character", o.EscapeCharacter)

	list := make([]interface{}, 0, len(serverList))
	for i := range serverList {
		list = append(list, map[string]interface{}{
			"name":          serverList[i].Name,
			"server":        serverList[i].Server,
			"transport":     serverList[i].Transport,
			"port":          serverList[i].Port,
			"syslog_format": serverList[i].SyslogFormat,
			"facility":      serverList[i].Facility,
		})
	}

	if err := d.Set("syslog_server", list); err != nil {
		log.Printf("[WARN] Error setting 'syslog_server' for %q: %s", d.Id(), err)
	}
}

func parseSyslogServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildSyslogServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o, serverList := parseSyslogServerProfile(d)

	if err := fw.Device.SyslogServerProfile.SetWithoutSubconfig(vsys, o); err != nil {
		return err
	}

	if err := fw.Device.SyslogServer.Set(vsys, o.Name, serverList...); err != nil {
		return err
	}

	d.SetId(buildSyslogServerProfileId(vsys, o.Name))
	return readSyslogServerProfile(d, meta)
}

func readSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseSyslogServerProfileId(d.Id())

	o, err := fw.Device.SyslogServerProfile.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err := fw.Device.SyslogServer.GetList(vsys, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := fw.Device.SyslogServer.Get(vsys, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	d.Set("vsys", vsys)
	saveSyslogServerProfile(d, o, serverList)

	return nil
}

func deleteSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseSyslogServerProfileId(d.Id())

	err := fw.Device.SyslogServerProfile.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
