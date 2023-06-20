package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/profile/security/dos"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceDosProtectionProfiles() *schema.Resource {
	s := listingSchema()
	s["vsys"] = vsysSchema("vsys1")
	s["device_group"] = deviceGroupSchema()

	return &schema.Resource{
		Read: dataSourceDosProtectionProfilesRead,

		Schema: s,
	}
}

func dataSourceDosProtectionProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string
	var id string

	switch con := meta.(type) {
	case *pango.Firewall:
		id = d.Get("vsys").(string)
		listing, err = con.Objects.DosProtectionProfile.GetList(id)
	case *pango.Panorama:
		id = d.Get("device_group").(string)
		listing, err = con.Objects.DosProtectionProfile.GetList(id)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceDosProtectionProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDosProtectionProfileRead,

		Schema: dosProtectionProfileSchema(false),
	}
}

func dataSourceDosProtectionProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	var o dos.Entry
	name := d.Get("name").(string)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDosProtectionProfileId(vsys, name)
		o, err = con.Objects.DosProtectionProfile.Get(vsys, name)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDosProtectionProfileId(dg, name)
		o, err = con.Objects.DosProtectionProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveDosProtectionProfile(d, o)

	return nil
}

// Resource.
func resourceDosProtectionProfile() *schema.Resource {
	return &schema.Resource{
		Create: createDosProtectionProfile,
		Read:   readDosProtectionProfile,
		Update: updateDosProtectionProfile,
		Delete: deleteDosProtectionProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: dosProtectionProfileSchema(true),
	}
}

func createDosProtectionProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var id string
	o := loadDosProtectionProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		id = buildDosProtectionProfileId(vsys, o.Name)
		err = con.Objects.DosProtectionProfile.Set(vsys, o)
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		id = buildDosProtectionProfileId(dg, o.Name)
		err = con.Objects.DosProtectionProfile.Set(dg, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readDosProtectionProfile(d, meta)
}

func readDosProtectionProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o dos.Entry

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDosProtectionProfileId(d.Id())
		o, err = con.Objects.DosProtectionProfile.Get(vsys, name)
	case *pango.Panorama:
		dg, name := parseDosProtectionProfileId(d.Id())
		o, err = con.Objects.DosProtectionProfile.Get(dg, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDosProtectionProfile(d, o)
	return nil
}

func updateDosProtectionProfile(d *schema.ResourceData, meta interface{}) error {
	o := loadDosProtectionProfile(d)

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys := d.Get("vsys").(string)
		lo, err := con.Objects.DosProtectionProfile.Get(vsys, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DosProtectionProfile.Edit(vsys, lo); err != nil {
			return err
		}
	case *pango.Panorama:
		dg := d.Get("device_group").(string)
		lo, err := con.Objects.DosProtectionProfile.Get(dg, o.Name)
		if err != nil {
			return err
		}
		lo.Copy(o)
		if err = con.Objects.DosProtectionProfile.Edit(dg, lo); err != nil {
			return err
		}
	}

	return readDosProtectionProfile(d, meta)
}

func deleteDosProtectionProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	switch con := meta.(type) {
	case *pango.Firewall:
		vsys, name := parseDosProtectionProfileId(d.Id())
		err = con.Objects.DosProtectionProfile.Delete(vsys, name)
	case *pango.Panorama:
		dg, name := parseDosProtectionProfileId(d.Id())
		err = con.Objects.DosProtectionProfile.Delete(dg, name)
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
func dosProtectionProfileSchema(isResource bool) map[string]*schema.Schema {
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
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The profile type",
			Default:     dos.TypeAggregate,
			ValidateFunc: validateStringIn(
				dos.TypeAggregate,
				dos.TypeClassified,
			),
		},
		"enable_sessions_protections": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable sessions protections",
		},
		"max_concurrent_sessions": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Max concurrent sessions",
		},
		"syn": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "SYN flood protection spec",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable this protection or not",
					},
					"action": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "SYN protection action",
						ValidateFunc: validateStringIn(
							dos.SynActionRed,
							dos.SynActionCookies,
						),
					},
					"alarm_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alarm rate",
					},
					"activate_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Activate rate",
					},
					"max_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Max rate",
					},
					"block_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block duration",
					},
				},
			},
		},
		"udp": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "UDP flood protection spec",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable this protection or not",
					},
					"alarm_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alarm rate",
					},
					"activate_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Activate rate",
					},
					"max_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Max rate",
					},
					"block_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block duration",
					},
				},
			},
		},
		"icmp": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "ICMP flood protection spec",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable this protection or not",
					},
					"alarm_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alarm rate",
					},
					"activate_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Activate rate",
					},
					"max_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Max rate",
					},
					"block_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block duration",
					},
				},
			},
		},
		"icmpv6": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "ICMPv6 flood protection spec",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable this protection or not",
					},
					"alarm_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alarm rate",
					},
					"activate_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Activate rate",
					},
					"max_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Max rate",
					},
					"block_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block duration",
					},
				},
			},
		},
		"other": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Other IP flood protection spec",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enable": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable this protection or not",
					},
					"alarm_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Alarm rate",
					},
					"activate_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Activate rate",
					},
					"max_rate": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Max rate",
					},
					"block_duration": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Block duration",
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

func loadDosProtectionProfile(d *schema.ResourceData) dos.Entry {
	var list []interface{}
	ans := dos.Entry{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		Type:                      d.Get("type").(string),
		EnableSessionsProtections: d.Get("enable_sessions_protections").(bool),
		MaxConcurrentSessions:     d.Get("max_concurrent_sessions").(int),
	}

	if list = d.Get("syn").([]interface{}); len(list) > 0 {
		x := list[0].(map[string]interface{})
		ans.Syn = &dos.SynProtection{
			Enable:        x["enable"].(bool),
			Action:        x["action"].(string),
			AlarmRate:     x["alarm_rate"].(int),
			ActivateRate:  x["activate_rate"].(int),
			MaxRate:       x["max_rate"].(int),
			BlockDuration: x["block_duration"].(int),
		}
	}

	if list = d.Get("udp").([]interface{}); len(list) > 0 {
		x := list[0].(map[string]interface{})
		ans.Udp = &dos.Protection{
			Enable:        x["enable"].(bool),
			AlarmRate:     x["alarm_rate"].(int),
			ActivateRate:  x["activate_rate"].(int),
			MaxRate:       x["max_rate"].(int),
			BlockDuration: x["block_duration"].(int),
		}
	}

	if list = d.Get("icmp").([]interface{}); len(list) > 0 {
		x := list[0].(map[string]interface{})
		ans.Icmp = &dos.Protection{
			Enable:        x["enable"].(bool),
			AlarmRate:     x["alarm_rate"].(int),
			ActivateRate:  x["activate_rate"].(int),
			MaxRate:       x["max_rate"].(int),
			BlockDuration: x["block_duration"].(int),
		}
	}

	if list = d.Get("icmpv6").([]interface{}); len(list) > 0 {
		x := list[0].(map[string]interface{})
		ans.Icmpv6 = &dos.Protection{
			Enable:        x["enable"].(bool),
			AlarmRate:     x["alarm_rate"].(int),
			ActivateRate:  x["activate_rate"].(int),
			MaxRate:       x["max_rate"].(int),
			BlockDuration: x["block_duration"].(int),
		}
	}

	if list = d.Get("other").([]interface{}); len(list) > 0 {
		x := list[0].(map[string]interface{})
		ans.Other = &dos.Protection{
			Enable:        x["enable"].(bool),
			AlarmRate:     x["alarm_rate"].(int),
			ActivateRate:  x["activate_rate"].(int),
			MaxRate:       x["max_rate"].(int),
			BlockDuration: x["block_duration"].(int),
		}
	}

	return ans
}

func saveDosProtectionProfile(d *schema.ResourceData, o dos.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("type", o.Type)
	d.Set("enable_sessions_protections", o.EnableSessionsProtections)
	d.Set("max_concurrent_sessions", o.MaxConcurrentSessions)

	if o.Syn == nil {
		d.Set("syn", nil)
	} else {
		prot := map[string]interface{}{
			"enable":         o.Syn.Enable,
			"action":         o.Syn.Action,
			"alarm_rate":     o.Syn.AlarmRate,
			"activate_rate":  o.Syn.ActivateRate,
			"max_rate":       o.Syn.MaxRate,
			"block_duration": o.Syn.BlockDuration,
		}

		if err := d.Set("syn", []interface{}{prot}); err != nil {
			log.Printf("[WARN] Error setting 'syn' for %q: %s", d.Id(), err)
		}
	}

	if o.Udp == nil {
		d.Set("udp", nil)
	} else {
		prot := map[string]interface{}{
			"enable":         o.Udp.Enable,
			"alarm_rate":     o.Udp.AlarmRate,
			"activate_rate":  o.Udp.ActivateRate,
			"max_rate":       o.Udp.MaxRate,
			"block_duration": o.Udp.BlockDuration,
		}

		if err := d.Set("udp", []interface{}{prot}); err != nil {
			log.Printf("[WARN] Error setting 'udp' for %q: %s", d.Id(), err)
		}
	}

	if o.Icmp == nil {
		d.Set("icmp", nil)
	} else {
		prot := map[string]interface{}{
			"enable":         o.Icmp.Enable,
			"alarm_rate":     o.Icmp.AlarmRate,
			"activate_rate":  o.Icmp.ActivateRate,
			"max_rate":       o.Icmp.MaxRate,
			"block_duration": o.Icmp.BlockDuration,
		}

		if err := d.Set("icmp", []interface{}{prot}); err != nil {
			log.Printf("[WARN] Error setting 'icmp' for %q: %s", d.Id(), err)
		}
	}

	if o.Icmpv6 == nil {
		d.Set("icmpv6", nil)
	} else {
		prot := map[string]interface{}{
			"enable":         o.Icmpv6.Enable,
			"alarm_rate":     o.Icmpv6.AlarmRate,
			"activate_rate":  o.Icmpv6.ActivateRate,
			"max_rate":       o.Icmpv6.MaxRate,
			"block_duration": o.Icmpv6.BlockDuration,
		}

		if err := d.Set("icmpv6", []interface{}{prot}); err != nil {
			log.Printf("[WARN] Error setting 'icmpv6' for %q: %s", d.Id(), err)
		}
	}

	if o.Other == nil {
		d.Set("other", nil)
	} else {
		prot := map[string]interface{}{
			"enable":         o.Other.Enable,
			"alarm_rate":     o.Other.AlarmRate,
			"activate_rate":  o.Other.ActivateRate,
			"max_rate":       o.Other.MaxRate,
			"block_duration": o.Other.BlockDuration,
		}

		if err := d.Set("other", []interface{}{prot}); err != nil {
			log.Printf("[WARN] Error setting 'other' for %q: %s", d.Id(), err)
		}
	}
}

// Id functions.
func buildDosProtectionProfileId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseDosProtectionProfileId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
