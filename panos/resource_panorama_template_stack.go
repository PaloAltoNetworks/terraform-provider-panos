package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template/stack"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaTemplateStack() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaTemplateStack,
		Read:   readPanoramaTemplateStack,
		Update: updatePanoramaTemplateStack,
		Delete: deletePanoramaTemplateStack,

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
			"templates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_vsys": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"devices": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parsePanoramaTemplateStack(d *schema.ResourceData) stack.Entry {
	o := stack.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		DefaultVsys: d.Get("default_vsys").(string),
		Templates:   asStringList(d.Get("templates").([]interface{})),
		Devices:     asStringList(d.Get("devices").([]interface{})),
	}

	return o
}

func createPanoramaTemplateStack(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaTemplateStack(d)

	if err = pano.Panorama.TemplateStack.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readPanoramaTemplateStack(d, meta)
}

func readPanoramaTemplateStack(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	o, err := pano.Panorama.TemplateStack.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("default_vsys", o.DefaultVsys)
	if err = d.Set("templates", o.Templates); err != nil {
		log.Printf("[WARN] Error setting 'templates' field for %q: %s", d.Id(), err)
	}
	if err = d.Set("devices", o.Devices); err != nil {
		log.Printf("[WARN] Error setting 'device' field for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaTemplateStack(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaTemplateStack(d)

	lo, err := pano.Panorama.TemplateStack.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.TemplateStack.Edit(lo); err != nil {
		return err
	}

	return readPanoramaTemplateStack(d, meta)
}

func deletePanoramaTemplateStack(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	err = pano.Panorama.TemplateStack.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
