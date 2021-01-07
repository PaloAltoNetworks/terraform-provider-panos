package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	tex "github.com/PaloAltoNetworks/pango/objs/profile/security/spyware/texception"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceAntiSpywareSecurityProfileThreatExceptions() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema()
	s["device_group"] = deviceGroupSchema()
	s["anti_spyware_security_profile"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The anti-spyware security profile name",
	}

	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfileThreatExceptionsRead,

		Schema: s,
	}
}

func dataSourceAntiSpywareSecurityProfileThreatExceptionsRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string
	prof := d.Get("anti_spyware_security_profile").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(vsys, prof, "")
		listing, err = con.Objects.AntiSpywareThreatException.GetList(vsys, prof)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(dg, prof, "")
		listing, err = con.Objects.AntiSpywareThreatException.GetList(dg, prof)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceAntiSpywareSecurityProfileThreatException() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntiSpywareSecurityProfileThreatExceptionRead,

		Schema: antiSpywareSecurityProfileThreatExceptionSchema(false),
	}
}

func dataSourceAntiSpywareSecurityProfileThreatExceptionRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o tex.Entry
	name := d.Get("name").(string)
	prof := d.Get("anti_spyware_security_profile").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(vsys, prof, name)
		o, err = con.Objects.AntiSpywareThreatException.Get(vsys, prof, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(dg, prof, name)
		o, err = con.Objects.AntiSpywareThreatException.Get(dg, prof, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveAntiSpywareSecurityProfileThreatException(d, o)

	return nil
}

// Resource.
func resourceAntiSpywareSecurityProfileThreatException() *schema.Resource {
	return &schema.Resource{
		Create: createAntiSpywareSecurityProfileThreatException,
		Read:   readAntiSpywareSecurityProfileThreatException,
		Update: updateAntiSpywareSecurityProfileThreatException,
		Delete: deleteAntiSpywareSecurityProfileThreatException,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: antiSpywareSecurityProfileThreatExceptionSchema(true),
	}
}

func createAntiSpywareSecurityProfileThreatException(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	prof := d.Get("anti_spyware_security_profile").(string)
	o := loadAntiSpywareSecurityProfileThreatException(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(vsys, prof, o.Name)
		err = con.Objects.AntiSpywareThreatException.Set(vsys, prof, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildAntiSpywareSecurityProfileThreatExceptionId(dg, prof, o.Name)
		err = con.Objects.AntiSpywareThreatException.Set(dg, prof, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readAntiSpywareSecurityProfileThreatException(d, meta)
}

func readAntiSpywareSecurityProfileThreatException(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o tex.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(d.Id())
		o, err = con.Objects.AntiSpywareThreatException.Get(vsys, prof, name)
	case *pango.Panorama:
		dg, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(d.Id())
		o, err = con.Objects.AntiSpywareThreatException.Get(dg, prof, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveAntiSpywareSecurityProfileThreatException(d, o)
	return nil
}

func updateAntiSpywareSecurityProfileThreatException(d *schema.ResourceData, meta interface{}) error {
	o := loadAntiSpywareSecurityProfileThreatException(d)
	prof := d.Get("anti_spyware_security_profile").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.AntiSpywareThreatException.Get(vsys, prof, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntiSpywareThreatException.Edit(vsys, prof, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.AntiSpywareThreatException.Get(dg, prof, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.AntiSpywareThreatException.Edit(dg, prof, lo); err != nil {
			return err
		}
	}

	return readAntiSpywareSecurityProfileThreatException(d, meta)
}

func deleteAntiSpywareSecurityProfileThreatException(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(d.Id())
		err = con.Objects.AntiSpywareThreatException.Delete(vsys, prof, name)
	case *pango.Panorama:
		dg, prof, name := parseAntiSpywareSecurityProfileThreatExceptionId(d.Id())
		err = con.Objects.AntiSpywareThreatException.Delete(dg, prof, name)
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
func antiSpywareSecurityProfileThreatExceptionSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"vsys":         vsysSchema(),
		"device_group": deviceGroupSchema(),
		"anti_spyware_security_profile": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The anti-spyware security profile name",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Security profile name",
		},
		"packet_capture": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(PAN-OS 8.x only) Packet capture config",
			Default:     tex.Disable,
			ValidateFunc: validateStringIn(
				"",
				tex.Disable,
				tex.SinglePacket,
				tex.ExtendedCapture,
			),
		},
		"action": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IPv4 sinkhole address",
			Default:     tex.ActionDefault,
			ValidateFunc: validateStringIn(
				"",
				tex.ActionDefault,
				tex.ActionAllow,
				tex.ActionAlert,
				tex.ActionDrop,
				tex.ActionResetClient,
				tex.ActionResetServer,
				tex.ActionResetBoth,
				tex.ActionBlockIp,
			),
		},
		"block_ip_track_by": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(action = block-ip) The track by config",
			ValidateFunc: validateStringIn(
				"",
				tex.TrackBySource,
				tex.TrackBySourceAndDestination,
			),
		},
		"block_ip_duration": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "(action = block-ip) The duration to block for",
		},
		"exempt_ips": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of exempt IP addresses",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{
			"vsys",
			"device_group",
			"name",
			"anti_spyware_security_profile",
		})
	}

	return ans
}

func loadAntiSpywareSecurityProfileThreatException(d *schema.ResourceData) tex.Entry {
	return tex.Entry{
		Name:            d.Get("name").(string),
		PacketCapture:   d.Get("packet_capture").(string),
		Action:          d.Get("action").(string),
		BlockIpTrackBy:  d.Get("block_ip_track_by").(string),
		BlockIpDuration: d.Get("block_ip_duration").(int),
		ExemptIps:       asStringList(d.Get("exempt_ips").([]interface{})),
	}
}

func saveAntiSpywareSecurityProfileThreatException(d *schema.ResourceData, o tex.Entry) {
	d.Set("name", o.Name)
	d.Set("packet_capture", o.PacketCapture)
	d.Set("action", o.Action)
	d.Set("block_ip_track_by", o.BlockIpTrackBy)
	d.Set("block_ip_duration", o.BlockIpDuration)
	if err := d.Set("exempt_ips", o.ExemptIps); err != nil {
		log.Printf("[WARN] Error setting 'exempt_ips' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildAntiSpywareSecurityProfileThreatExceptionId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseAntiSpywareSecurityProfileThreatExceptionId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}
