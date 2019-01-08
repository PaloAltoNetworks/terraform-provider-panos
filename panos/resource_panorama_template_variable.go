package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template/variable"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaTemplateVariable() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaTemplateVariable,
		Read:   readPanoramaTemplateVariable,
		Update: updatePanoramaTemplateVariable,
		Delete: deletePanoramaTemplateVariable,

		Schema: map[string]*schema.Schema{
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringHasPrefix("$"),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      variable.TypeIpNetmask,
				ValidateFunc: validateStringIn(variable.TypeIpNetmask, variable.TypeIpRange, variable.TypeFqdn, variable.TypeGroupId, variable.TypeInterface),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func parsePanoramaTemplateVariableId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaTemplateVariableId(a, b, c string) string {
	return fmt.Sprintf("%s%s%s%s%s", a, IdSeparator, b, IdSeparator, c)
}

func parsePanoramaTemplateVariable(d *schema.ResourceData) (string, string, variable.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	o := variable.Entry{
		Name:  d.Get("name").(string),
		Type:  d.Get("type").(string),
		Value: d.Get("value").(string),
	}

	return tmpl, ts, o
}

func createPanoramaTemplateVariable(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaTemplateVariable(d)

	if err = pano.Panorama.TemplateVariable.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaTemplateVariableId(tmpl, ts, o.Name))
	return readPanoramaTemplateVariable(d, meta)
}

func readPanoramaTemplateVariable(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaTemplateVariableId(d.Id())

	o, err := pano.Panorama.TemplateVariable.Get(tmpl, ts, name)
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
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	d.Set("value", o.Value)

	return nil
}

func updatePanoramaTemplateVariable(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaTemplateVariable(d)

	lo, err := pano.Panorama.TemplateVariable.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.TemplateVariable.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaTemplateVariable(d, meta)
}

func deletePanoramaTemplateVariable(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaTemplateVariableId(d.Id())

	err = pano.Panorama.TemplateVariable.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
