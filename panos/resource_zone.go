package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: createZone,
		Read:   readZone,
		Update: updateZone,
		Delete: deleteZone,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this zone in",
			},
			"mode": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone's mode",
			},
			"zone_profile": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone's mode",
			},
			"log_setting": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone's mode",
			},
			"enable_user_id": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The zone's mode",
			},
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User Identification include ACL list",
			},
			"include_acl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User Identification include ACL list",
			},
			"exclude_acl": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User Identification exclude ACL list",
			},
		},
	}
}

func parseZone(d *schema.ResourceData) (string, zone.Entry) {
	vsys := d.Get("vsys").(string)
	o := zone.Entry{
		Name:         d.Get("name").(string),
		Mode:         d.Get("mode").(string),
		ZoneProfile:  d.Get("zone_profile").(string),
		LogSetting:   d.Get("log_setting").(string),
		EnableUserId: d.Get("enable_user_id").(bool),
		Interfaces:   asStringList(d, "interfaces"),
		IncludeAcl:   asStringList(d, "include_acl"),
		ExcludeAcl:   asStringList(d, "exclude_acl"),
	}

	return vsys, o
}

func parseZoneId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildZoneId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func saveDataZone(d *schema.ResourceData, vsys string, o zone.Entry) {
	d.SetId(buildZoneId(vsys, o.Name))
	d.Set("vsys", vsys)
	d.Set("name", o.Name)
	d.Set("mode", o.Mode)
	d.Set("zone_profile", o.ZoneProfile)
	d.Set("log_setting", o.LogSetting)
	d.Set("enable_user_id", o.EnableUserId)
	d.Set("interfaces", o.Interfaces)
	d.Set("include_acl", o.IncludeAcl)
	d.Set("exclude_acl", o.ExcludeAcl)
}

func createZone(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseZone(d)

	if err := fw.Network.Zone.Set(vsys, o); err != nil {
		return err
	}

	saveDataZone(d, vsys, o)
	return nil
}

func readZone(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseZoneId(d.Id())

	o, err := fw.Network.Zone.Get(vsys, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataZone(d, vsys, o)
	return nil
}

func updateZone(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	vsys, o := parseZone(d)

	lo, err := fw.Network.Zone.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	err = fw.Network.Zone.Edit(vsys, lo)

	if err == nil {
		saveDataZone(d, vsys, o)
	}
	return err
}

func deleteZone(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseZoneId(d.Id())

	_ = fw.Network.Zone.Delete(vsys, name)
	d.SetId("")
	return nil
}
