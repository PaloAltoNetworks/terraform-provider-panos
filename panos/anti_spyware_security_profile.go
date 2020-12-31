package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceAntiSpywareSecurityProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema()
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfilesRead,

		Schema: s,
	}
}

func dataSourceAntiSpywareSecurityProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.AntiSpywareProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.AntiSpywareProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceAntiSpywareSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfileRead,

		Schema: antiSpywareSecurityProfileSchema(false),
	}
}

func dataSourceAntiSpywareSecurityProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o spyware.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileId(vsys, name)
		o, err = con.Objects.AntiSpywareProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileId(dg, name)
		o, err = con.Objects.AntiSpywareProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveAntiSpywareSecurityProfile(d, o)

	return nil
}

// Resource.
func resourceAntiSpywareSecurityProfile() *schema.Resource {
	return &schema.Resource{
		Create: createAntiSpywareSecurityProfile,
		Read:   readAntiSpywareSecurityProfile,
		Update: updateAntiSpywareSecurityProfile,
		Delete: deleteAntiSpywareSecurityProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: antiSpywareSecurityProfileSchema(true),
	}
}

func createAntiSpywareSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadAntiSpywareSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileId(vsys, o.Name)
		err = con.Objects.AntiSpywareProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileId(dg, o.Name)
		err = con.Objects.AntiSpywareProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAntiSpywareSecurityProfile(d, meta)
}

func readAntiSpywareSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o spyware.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAntiSpywareSecurityProfileId(d.Id())
		o, err = con.Objects.AntiSpywareProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseAntiSpywareSecurityProfileId(d.Id())
		o, err = con.Objects.AntiSpywareProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAntiSpywareSecurityProfile(d, o)
	return nil
}

func updateAntiSpywareSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadAntiSpywareSecurityProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.AntiSpywareProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntiSpywareProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.AntiSpywareProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntiSpywareProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readAntiSpywareSecurityProfile(d, meta)
}

func deleteAntiSpywareSecurityProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseAntiSpywareSecurityProfileId(d.Id())
		err = con.Objects.AntiSpywareProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseAntiSpywareSecurityProfileId(d.Id())
		err = con.Objects.AntiSpywareProfile.Delete(dg, name)
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
func antiSpywareSecurityProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys":         vsysSchema(),
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
		"packet_capture": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 8.x only) Packet capture config",
		},
		"sinkhole_ipv4_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv4 sinkhole address",
		},
		"sinkhole_ipv6_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv6 sinkhole address",
		},
		"threat_exceptions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of threat exceptions",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"botnet_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Botnet list structs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"action": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Action to take",
					},
					"packet_capture": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "(PAN-OS 9.0+) Packet capture config",
						ValidateFunc: validateStringIn(
							"",
							spyware.Disable,
							spyware.SinglePacket,
							spyware.ExtendedCapture,
						),
					},
				},
			},
		},
		"dns_category": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "(PAN-OS 10.0+) DNS category structs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"action": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Action to take",
					},
					"log_level": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Logging level",
					},
					"packet_capture": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Packet capture config",
						ValidateFunc: validateStringIn(
							"",
							spyware.Disable,
							spyware.SinglePacket,
							spyware.ExtendedCapture,
						),
					},
				},
			},
		},
		"white_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "White list structs",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name",
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description",
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "device_group", "name"})
	}

	return ans
}

func loadAntiSpywareSecurityProfile(d *schema.ResourceData) spyware.Entry {
	var list []interface{}

	var botnets []spyware.BotnetList
	list = d.Get("botnet_list").([]interface{})
	if len(list) > 0 {
		botnets = make([]spyware.BotnetList, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			botnets = append(botnets, spyware.BotnetList{
				Name:          elm["name"].(string),
				Action:        elm["action"].(string),
				PacketCapture: elm["packet_capture"].(string),
			})
		}
	}

	var dnsCategories []spyware.DnsCategory
	list = d.Get("dns_category").([]interface{})
	if len(list) > 0 {
		dnsCategories = make([]spyware.DnsCategory, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			dnsCategories = append(dnsCategories, spyware.DnsCategory{
				Name:          elm["name"].(string),
				Action:        elm["action"].(string),
				LogLevel:      elm["log_level"].(string),
				PacketCapture: elm["packet_capture"].(string),
			})
		}
	}

	var whiteLists []spyware.WhiteList
	list = d.Get("white_list").([]interface{})
	if len(list) > 0 {
		whiteLists = make([]spyware.WhiteList, 0, len(list))
		for i := range list {
			elm := list[i].(map[string]interface{})
			whiteLists = append(whiteLists, spyware.WhiteList{
				Name:        elm["name"].(string),
				Description: elm["description"].(string),
			})
		}
	}

	return spyware.Entry{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		PacketCapture:       d.Get("packet_capture").(string),
		BotnetLists:         botnets,
		DnsCategories:       dnsCategories,
		WhiteLists:          whiteLists,
		SinkholeIpv4Address: d.Get("sinkhole_ipv4_address").(string),
		SinkholeIpv6Address: d.Get("sinkhole_ipv6_address").(string),
		ThreatExceptions:    asStringList(d.Get("threat_exceptions").([]interface{})),
	}
}

func saveAntiSpywareSecurityProfile(d *schema.ResourceData, o spyware.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("packet_capture", o.PacketCapture)
	d.Set("sinkhole_ipv4_address", o.SinkholeIpv4Address)
	d.Set("sinkhole_ipv6_address", o.SinkholeIpv6Address)
	if err := d.Set("threat_exceptions", o.ThreatExceptions); err != nil {
		log.Printf("[WARN] Error setting 'threat_exceptions' for %q: %s", d.Id(), err)
	}

	if len(o.BotnetLists) == 0 {
		d.Set("botnet_list", nil)
	} else {
		list := make([]interface{}, 0, len(o.BotnetLists))
		for _, x := range o.BotnetLists {
			list = append(list, map[string]interface{}{
				"name":           x.Name,
				"action":         x.Action,
				"packet_capture": x.PacketCapture,
			})
		}
		if err := d.Set("botnet_list", list); err != nil {
			log.Printf("[WARN] Error setting 'botnet_list' for %q: %s", d.Id(), err)
		}
	}

	if len(o.DnsCategories) == 0 {
		d.Set("dns_category", nil)
	} else {
		list := make([]interface{}, 0, len(o.DnsCategories))
		for _, x := range o.DnsCategories {
			list = append(list, map[string]interface{}{
				"name":           x.Name,
				"action":         x.Action,
				"log_level":      x.LogLevel,
				"packet_capture": x.PacketCapture,
			})
		}
		if err := d.Set("dns_category", list); err != nil {
			log.Printf("[WARN] Error setting 'dns_category' for %q: %s", d.Id(), err)
		}
	}

	if len(o.WhiteLists) == 0 {
		d.Set("white_list", nil)
	} else {
		list := make([]interface{}, 0, len(o.WhiteLists))
		for _, x := range o.WhiteLists {
			list = append(list, map[string]interface{}{
				"name":        x.Name,
				"description": x.Description,
			})
		}
		if err := d.Set("white_list", list); err != nil {
			log.Printf("[WARN] Error setting 'white_list' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildAntiSpywareSecurityProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseAntiSpywareSecurityProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
