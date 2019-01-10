package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaDeviceGroupEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaDeviceGroupEntry,
		Read:   readPanoramaDeviceGroupEntry,
		Update: createUpdatePanoramaDeviceGroupEntry,
		Delete: deletePanoramaDeviceGroupEntry,

		Schema: map[string]*schema.Schema{
			"device_group": {
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

func parsePanoramaDeviceGroupEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaDeviceGroupEntryId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createUpdatePanoramaDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	grp := d.Get("device_group").(string)
	dev := d.Get("serial").(string)
	vl := asStringList(d.Get("vsys_list").(*schema.Set).List())

	if err := pano.Panorama.DeviceGroup.EditDeviceVsys(grp, dev, vl); err != nil {
		return err
	}

	d.SetId(buildPanoramaDeviceGroupEntryId(grp, dev))
	return readPanoramaDeviceGroupEntry(d, meta)
}

func readPanoramaDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	grp, dev := parsePanoramaDeviceGroupEntryId(d.Id())

	// Two possibilities:  either the group itself doesn't exist, or the
	// device is not in the group.
	o, err := pano.Panorama.DeviceGroup.Get(grp)
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
			d.Set("device_group", grp)
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

func deletePanoramaDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	grp, dev := parsePanoramaDeviceGroupEntryId(d.Id())

	err := pano.Panorama.DeviceGroup.DeleteDeviceVsys(grp, dev, nil)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
