package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaTemplateEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaTemplateEntry,
		Read:   readPanoramaTemplateEntry,
		Update: createUpdatePanoramaTemplateEntry,
		Delete: deletePanoramaTemplateEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"serial": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parsePanoramaTemplateEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaTemplateEntryId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createUpdatePanoramaTemplateEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl := d.Get("template").(string)
	dev := d.Get("serial").(string)
	vl := asStringList(d.Get("vsys_list").(*schema.Set).List())

	if err := pano.Panorama.Template.EditDeviceVsys(tmpl, dev, vl); err != nil {
		return err
	}

	d.SetId(buildPanoramaTemplateEntryId(tmpl, dev))
	return readPanoramaTemplateEntry(d, meta)
}

func readPanoramaTemplateEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, dev := parsePanoramaTemplateEntryId(d.Id())

	// Two possibilities:  either the group itself doesn't exist, or the
	// device is not in the group.
	o, err := pano.Panorama.Template.Get(tmpl)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	for i := range o.Devices {
		if i == dev {
			d.Set("template", tmpl)
			d.Set("serial", dev)
			if err = d.Set("vsys_list", listAsSet(o.Devices[i])); err != nil {
				log.Printf("[WARN] Error setting 'vsys_list' param for %q: %s", d.Id(), err)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deletePanoramaTemplateEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, dev := parsePanoramaTemplateEntryId(d.Id())

	err := pano.Panorama.Template.DeleteDeviceVsys(tmpl, dev, nil)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
