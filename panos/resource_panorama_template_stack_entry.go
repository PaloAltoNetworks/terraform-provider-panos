package panos

import (
	"fmt"
	"strings"

	"github.com/fpluchorg/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaTemplateStackEntry() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaTemplateStackEntry,
		Read:   readPanoramaTemplateStackEntry,
		Delete: deletePanoramaTemplateStackEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"template_stack": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func parsePanoramaTemplateStackEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaTemplateStackEntryId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaTemplateStackEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	ts := d.Get("template_stack").(string)
	dev := d.Get("device").(string)

	if err := pano.Panorama.TemplateStack.EditDevice(ts, dev); err != nil {
		return err
	}

	d.SetId(buildPanoramaTemplateStackEntryId(ts, dev))
	return readPanoramaTemplateStackEntry(d, meta)
}

func readPanoramaTemplateStackEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	ts, dev := parsePanoramaTemplateStackEntryId(d.Id())

	// Two possibilities:  either the group itself doesn't exist, or the
	// device is not in the group.
	o, err := pano.Panorama.TemplateStack.Get(ts)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for i := range o.Devices {
		if o.Devices[i] == dev {
			d.Set("template_stack", ts)
			d.Set("device", dev)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deletePanoramaTemplateStackEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	ts, dev := parsePanoramaTemplateStackEntryId(d.Id())

	err := pano.Panorama.TemplateStack.DeleteDevice(ts, dev)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
