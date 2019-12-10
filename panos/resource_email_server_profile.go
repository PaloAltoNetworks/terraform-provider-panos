package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"
	"github.com/PaloAltoNetworks/pango/dev/profile/email/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEmailServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateEmailServerProfile,
		Read:   readEmailServerProfile,
		Update: createUpdateEmailServerProfile,
		Delete: deleteEmailServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: emailServerProfileSchema(false),
	}
}

func emailServerProfileSchema(p bool) map[string]*schema.Schema {
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
		"email_server": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"display_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"from_email": {
						Type:     schema.TypeString,
						Required: true,
					},
					"to_email": {
						Type:     schema.TypeString,
						Required: true,
					},
					"also_to_email": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"email_gateway": {
						Type:     schema.TypeString,
						Required: true,
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

func parseEmailServerProfile(d *schema.ResourceData) (string, email.Entry, []server.Entry) {
	vsys := d.Get("vsys").(string)
	o, serverList := loadEmailServerProfile(d)

	return vsys, o, serverList
}

func loadEmailServerProfile(d *schema.ResourceData) (email.Entry, []server.Entry) {
	o := email.Entry{
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

	sl := d.Get("email_server").([]interface{})
	serverList := make([]server.Entry, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		serverList = append(serverList, server.Entry{
			Name:         x["name"].(string),
			DisplayName:  x["display_name"].(string),
			From:         x["from_email"].(string),
			To:           x["to_email"].(string),
			AlsoTo:       x["also_to_email"].(string),
			EmailGateway: x["email_gateway"].(string),
		})
	}

	return o, serverList
}

func saveEmailServerProfile(d *schema.ResourceData, o email.Entry, serverList []server.Entry) {
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
			"display_name":  serverList[i].DisplayName,
			"from_email":    serverList[i].From,
			"to_email":      serverList[i].To,
			"also_to_email": serverList[i].AlsoTo,
			"email_gateway": serverList[i].EmailGateway,
		})
	}

	if err := d.Set("email_server", list); err != nil {
		log.Printf("[WARN] Error setting 'email_server' for %q: %s", d.Id(), err)
	}
}

func parseEmailServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildEmailServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func createUpdateEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o, serverList := parseEmailServerProfile(d)

	if err := fw.Device.EmailServerProfile.SetWithoutSubconfig(vsys, o); err != nil {
		return err
	}

	if err := fw.Device.EmailServer.Set(vsys, o.Name, serverList...); err != nil {
		return err
	}

	d.SetId(buildEmailServerProfileId(vsys, o.Name))
	return readEmailServerProfile(d, meta)
}

func readEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseEmailServerProfileId(d.Id())

	o, err := fw.Device.EmailServerProfile.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err := fw.Device.EmailServer.GetList(vsys, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := fw.Device.EmailServer.Get(vsys, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	d.Set("vsys", vsys)
	saveEmailServerProfile(d, o, serverList)

	return nil
}

func deleteEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseEmailServerProfileId(d.Id())

	err := fw.Device.EmailServerProfile.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
