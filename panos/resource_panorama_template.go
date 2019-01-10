package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourcePanoramaTemplate() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaTemplate,
		Read:   readPanoramaTemplate,
		Update: updatePanoramaTemplate,
		Delete: deletePanoramaTemplate,

		SchemaVersion: 1,
		MigrateState:  migrateResourcePanoramaTemplate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_vsys": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"devices": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				// TODO(gfreeman): Uncomment once ValidateFunc is supported for TypeSet.
				//ValidateFunc: validateSetKeyIsUnique("serial"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vsys_list": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func migrateResourcePanoramaTemplate(ov int, s *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	pano := meta.(*pango.Panorama)

	if ov == 0 {
		// Previously missing ./config, so editing the template will fix this.
		name := s.ID
		info, err := pano.Panorama.Template.Get(name)
		if err != nil {
			return s, err
		}
		// Also, set the default vsys to vsys1 if it isn't already set.  This
		// field was previously Optional/Computed, and is now Optional/Default.
		if info.DefaultVsys == "" {
			info.DefaultVsys = "vsys1"
		}
		if err = pano.Panorama.Template.Edit(info); err != nil {
			return s, err
		}

		ov = 1
	}

	return s, nil
}

func parsePanoramaTemplate(d *schema.ResourceData) template.Entry {
	o := template.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DefaultVsys: "vsys1",
	}

	m := make(map[string][]string)
	dl := d.Get("devices").(*schema.Set).List()
	for i := range dl {
		device := dl[i].(map[string]interface{})
		key := device["serial"].(string)
		value := asStringList(device["vsys_list"].(*schema.Set).List())
		m[key] = value
	}
	o.Devices = m

	return o
}

func createPanoramaTemplate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaTemplate(d)

	if err = pano.Panorama.Template.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readPanoramaTemplate(d, meta)
}

func readPanoramaTemplate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	o, err := pano.Panorama.Template.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	ds := d.Get("devices").(*schema.Set)
	s := &schema.Set{F: ds.F}
	for key := range o.Devices {
		sg := make(map[string]interface{})
		sg["serial"] = key
		sg["vsys_list"] = listAsSet(o.Devices[key])
		s.Add(sg)
	}

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("default_vsys", o.DefaultVsys)
	if err = d.Set("devices", s); err != nil {
		log.Printf("[WARN] Error setting 'device' field for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaTemplate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaTemplate(d)

	lo, err := pano.Panorama.Template.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.Template.Edit(lo); err != nil {
		return err
	}

	return readPanoramaTemplate(d, meta)
}

func deletePanoramaTemplate(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	err = pano.Panorama.Template.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
