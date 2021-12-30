package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceEmailServerProfiles() *schema.Resource {
	s := listingSchema()
	for key, val := range templateWithPanoramaSharedSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceEmailServerProfilesRead,

		Schema: s,
	}
}

func dataSourceEmailServerProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = vsys
		listing, err = con.Device.EmailServerProfile.GetList(vsys)
	case *pango.Panorama:
		id = strings.Join([]string{tmpl, ts, vsys}, IdSeparator)
		listing, err = con.Device.EmailServerProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceEmailServerProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEmailServerProfileRead,

		Schema: emailServerProfileSchema(false, "shared", []string{"device_group"}),
	}
}

func dataSourceEmailServerProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o email.Entry

	name := d.Get("name").(string)
	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildEmailServerProfileId(vsys, name)
		o, err = con.Device.EmailServerProfile.Get(vsys, name)
	case *pango.Panorama:
		id = buildPanoramaEmailServerProfileId(tmpl, ts, vsys, name)
		o, err = con.Device.EmailServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveEmailServerProfile(d, o)

	return nil
}

// Resource.
func resourceEmailServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createEmailServerProfile,
		Read:   readEmailServerProfile,
		Update: updateEmailServerProfile,
		Delete: deleteEmailServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: emailServerProfileSchema(true, "shared", []string{"device_group", "template", "template_stack"}),
	}
}

func resourcePanoramaEmailServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createEmailServerProfile,
		Read:   readEmailServerProfile,
		Update: updateEmailServerProfile,
		Delete: deleteEmailServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: emailServerProfileSchema(true, "", nil),
	}
}

func createEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadEmailServerProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildEmailServerProfileId(vsys, o.Name)
		err = con.Device.EmailServerProfile.Set(vsys, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		id = buildPanoramaEmailServerProfileId(tmpl, ts, vsys, o.Name)
		err = con.Device.EmailServerProfile.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
	return readEmailServerProfile(d, meta)
}

func readEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o email.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseEmailServerProfileId(d.Id())
		d.Set("vsys", vsys)
		o, err = con.Device.EmailServerProfile.Get(vsys, name)
	case *pango.Panorama:
		// If this is the old style ID, it will have an extra field.  So
		// we need to migrate the ID first.
		tok := strings.Split(d.Id(), IdSeparator)
		if len(tok) == 5 {
			d.SetId(buildPanoramaEmailServerProfileId(tok[0], tok[1], tok[2], tok[4]))
		}
		// Continue on as normal.
		tmpl, ts, vsys, name := parsePanoramaEmailServerProfileId(d.Id())
		d.Set("template", tmpl)
		d.Set("template_stack", ts)
		d.Set("vsys", vsys)
		d.Set("device_group", d.Get("device_group").(string))
		o, err = con.Device.EmailServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveEmailServerProfile(d, o)

	return nil
}

func updateEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadEmailServerProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Device.EmailServerProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.EmailServerProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		lo, err := con.Device.EmailServerProfile.Get(tmpl, ts, vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.EmailServerProfile.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
	}

	return readEmailServerProfile(d, meta)
}

func deleteEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseEmailServerProfileId(d.Id())
		err = con.Device.EmailServerProfile.Delete(vsys, name)
	case *pango.Panorama:
		tmpl, ts, vsys, name := parsePanoramaEmailServerProfileId(d.Id())
		err = con.Device.EmailServerProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}

// Schema handling.
func emailServerProfileSchema(isResource bool, vsysDefault string, rmKeys []string) map[string]*schema.Schema {
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

	for _, rmKey := range rmKeys {
		delete(ans, rmKey)
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "device_group", "name"})
	}

	return ans
}

func loadEmailServerProfile(d *schema.ResourceData) email.Entry {
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
	serverList := make([]email.Server, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		serverList = append(serverList, email.Server{
			Name:         x["name"].(string),
			DisplayName:  x["display_name"].(string),
			From:         x["from_email"].(string),
			To:           x["to_email"].(string),
			AlsoTo:       x["also_to_email"].(string),
			EmailGateway: x["email_gateway"].(string),
		})
	}
	o.Servers = serverList

	return o
}

func saveEmailServerProfile(d *schema.ResourceData, o email.Entry) {
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

	list := make([]interface{}, 0, len(o.Servers))
	for _, x := range o.Servers {
		list = append(list, map[string]interface{}{
			"name":          x.Name,
			"display_name":  x.DisplayName,
			"from_email":    x.From,
			"to_email":      x.To,
			"also_to_email": x.AlsoTo,
			"email_gateway": x.EmailGateway,
		})
	}

	if err := d.Set("email_server", list); err != nil {
		log.Printf("[WARN] Error setting 'email_server' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseEmailServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaEmailServerProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildEmailServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaEmailServerProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
