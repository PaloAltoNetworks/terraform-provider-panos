package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: createZone,
		Read:   readZone,
		Update: updateZone,
		Delete: deleteZone,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone's name",
			},
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this zone in",
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The zone's mode",
				ValidateFunc: validateStringIn("layer3", "layer2", "virtual-wire", "tap", "tunnel"),
			},
			"zone_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone's mode",
			},
			"log_setting": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The zone's mode",
			},
			"enable_user_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The zone's mode",
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User Identification include ACL list",
			},
			"include_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "User Identification include ACL list",
			},
			"exclude_acls": {
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
		Interfaces:   asStringList(d.Get("interfaces").([]interface{})),
		IncludeAcls:  asStringList(d.Get("include_acls").([]interface{})),
		ExcludeAcls:  asStringList(d.Get("exclude_acls").([]interface{})),
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

func createZone(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseZone(d)

	if err := fw.Network.Zone.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildZoneId(vsys, o.Name))
	return readZone(d, meta)
}

func readZone(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseZoneId(d.Id())

	o, err := fw.Network.Zone.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vsys", vsys)
	d.Set("name", o.Name)
	d.Set("mode", o.Mode)
	d.Set("zone_profile", o.ZoneProfile)
	d.Set("log_setting", o.LogSetting)
	d.Set("enable_user_id", o.EnableUserId)
	if err = d.Set("interfaces", o.Interfaces); err != nil {
		log.Printf("[WARN] Error setting 'interfaces' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("include_acls", o.IncludeAcls); err != nil {
		log.Printf("[WARN] Error setting 'include_acls' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("exclude_acls", o.ExcludeAcls); err != nil {
		log.Printf("[WARN] Error setting 'exclude_acls' param for %q: %s", d.Id(), err)
	}

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
	if err = fw.Network.Zone.Edit(vsys, lo); err != nil {
		return err
	}

	return readZone(d, meta)
}

func deleteZone(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseZoneId(d.Id())

	err := fw.Network.Zone.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
