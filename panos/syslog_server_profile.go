package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceSyslogServerProfiles() *schema.Resource {
	s := listingSchema()
	for key, val := range templateWithPanoramaSharedSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceSyslogServerProfilesRead,

		Schema: s,
	}
}

func dataSourceSyslogServerProfilesRead(d *schema.ResourceData, meta interface{}) error {
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
		listing, err = con.Device.SyslogServerProfile.GetList(vsys)
	case *pango.Panorama:
		id = strings.Join([]string{tmpl, ts, vsys}, IdSeparator)
		listing, err = con.Device.SyslogServerProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceSyslogServerProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSyslogServerProfileRead,

		Schema: syslogServerProfileSchema(false, "shared", []string{"device_group"}),
	}
}

func dataSourceSyslogServerProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o syslog.Entry

	name := d.Get("name").(string)
	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		id = buildSyslogServerProfileId(vsys, name)
		o, err = con.Device.SyslogServerProfile.Get(vsys, name)
	case *pango.Panorama:
		id = buildPanoramaSyslogServerProfileId(tmpl, ts, vsys, name)
		o, err = con.Device.SyslogServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveSyslogServerProfile(d, o)

	return nil
}

// Resource.
func resourceSyslogServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSyslogServerProfile,
		Read:   readSyslogServerProfile,
		Update: updateSyslogServerProfile,
		Delete: deleteSyslogServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: syslogServerProfileSchema(true, "shared", []string{"device_group", "template", "template_stack"}),
	}
}

func resourcePanoramaSyslogServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSyslogServerProfile,
		Read:   readSyslogServerProfile,
		Update: updateSyslogServerProfile,
		Delete: deleteSyslogServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: syslogServerProfileSchema(true, "", nil),
	}
}

func createSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadSyslogServerProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildSyslogServerProfileId(vsys, o.Name)
		err = con.Device.SyslogServerProfile.Set(vsys, o)
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		id = buildPanoramaSyslogServerProfileId(tmpl, ts, vsys, o.Name)
		err = con.Device.SyslogServerProfile.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSyslogServerProfile(d, meta)
}

func readSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o syslog.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseSyslogServerProfileId(d.Id())
		d.Set("vsys", vsys)
		o, err = con.Device.SyslogServerProfile.Get(vsys, name)
	case *pango.Panorama:
		// If this is the old style ID, it will have an extra field.  So
		// we need to migrate the ID first.
		tok := strings.Split(d.Id(), IdSeparator)
		if len(tok) == 5 {
			d.SetId(buildPanoramaSyslogServerProfileId(tok[0], tok[1], tok[2], tok[4]))
		}
		// Continue on as normal.
		tmpl, ts, vsys, name := parsePanoramaSyslogServerProfileId(d.Id())
		d.Set("template", tmpl)
		d.Set("template_stack", ts)
		d.Set("vsys", vsys)
		d.Set("device_group", d.Get("device_group").(string))
		o, err = con.Device.SyslogServerProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSyslogServerProfile(d, o)

	return nil
}

func updateSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadSyslogServerProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Device.SyslogServerProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.SyslogServerProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		tmpl := d.Get("template").(string)
		ts := d.Get("template_stack").(string)
		vsys := d.Get("vsys").(string)
		lo, err := con.Device.SyslogServerProfile.Get(tmpl, ts, vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Device.SyslogServerProfile.Edit(tmpl, ts, vsys, lo); err != nil {
			return err
		}
	}

	return readSyslogServerProfile(d, meta)
}

func deleteSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseSyslogServerProfileId(d.Id())
		err = con.Device.SyslogServerProfile.Delete(vsys, name)
	case *pango.Panorama:
		tmpl, ts, vsys, name := parsePanoramaSyslogServerProfileId(d.Id())
		err = con.Device.SyslogServerProfile.Delete(tmpl, ts, vsys, name)
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
func syslogServerProfileSchema(isResource bool, vsysDefault string, rmKeys []string) map[string]*schema.Schema {
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
						Default:  syslog.TransportUdp,
						ValidateFunc: validateStringIn(
							syslog.TransportUdp,
							syslog.TransportTcp,
							syslog.TransportSsl,
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
						Default:  syslog.SyslogFormatBsd,
						ValidateFunc: validateStringIn(
							syslog.SyslogFormatBsd,
							syslog.SyslogFormatIetf,
						),
					},
					"facility": {
						Type:     schema.TypeString,
						Optional: true,
						Default:  syslog.FacilityUser,
						ValidateFunc: validateStringIn(
							syslog.FacilityUser,
							syslog.FacilityLocal0,
							syslog.FacilityLocal1,
							syslog.FacilityLocal2,
							syslog.FacilityLocal3,
							syslog.FacilityLocal4,
							syslog.FacilityLocal5,
							syslog.FacilityLocal6,
							syslog.FacilityLocal7,
						),
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

func loadSyslogServerProfile(d *schema.ResourceData) syslog.Entry {
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
	servers := make([]syslog.Server, 0, len(sl))
	for i := range sl {
		x := sl[i].(map[string]interface{})
		servers = append(servers, syslog.Server{
			Name:         x["name"].(string),
			Server:       x["server"].(string),
			Transport:    x["transport"].(string),
			Port:         x["port"].(int),
			SyslogFormat: x["syslog_format"].(string),
			Facility:     x["facility"].(string),
		})
	}
	o.Servers = servers

	return o
}

func saveSyslogServerProfile(d *schema.ResourceData, o syslog.Entry) {
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
			"server":        x.Server,
			"transport":     x.Transport,
			"port":          x.Port,
			"syslog_format": x.SyslogFormat,
			"facility":      x.Facility,
		})
	}

	if err := d.Set("syslog_server", list); err != nil {
		log.Printf("[WARN] Error setting 'syslog_server' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseSyslogServerProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func parsePanoramaSyslogServerProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildSyslogServerProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func buildPanoramaSyslogServerProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
