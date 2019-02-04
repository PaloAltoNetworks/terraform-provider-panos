package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaZone() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaZone,
		Read:   readPanoramaZone,
		Update: updatePanoramaZone,
		Delete: deletePanoramaZone,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
			"vsys": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("layer3", "layer2", "virtual-wire", "tap", "tunnel"),
			},
			"zone_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_setting": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_user_id": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parsePanoramaZone(d *schema.ResourceData) (string, string, string, zone.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
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

	return tmpl, ts, vsys, o
}

func parsePanoramaZoneId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildPanoramaZoneId(a, b, c, d string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c, IdSeparator, d)
}

func createPanoramaZone(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaZone(d)

	if err := pano.Network.Zone.Set(tmpl, ts, vsys, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaZoneId(tmpl, ts, vsys, o.Name))
	return readPanoramaZone(d, meta)
}

func readPanoramaZone(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaZoneId(d.Id())

	o, err := pano.Network.Zone.Get(tmpl, ts, vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
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

func updatePanoramaZone(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, o := parsePanoramaZone(d)

	lo, err := pano.Network.Zone.Get(tmpl, ts, vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.Zone.Edit(tmpl, ts, vsys, lo); err != nil {
		return err
	}

	return readPanoramaZone(d, meta)
}

func deletePanoramaZone(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, name := parsePanoramaZoneId(d.Id())

	err := pano.Network.Zone.Delete(tmpl, ts, vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
