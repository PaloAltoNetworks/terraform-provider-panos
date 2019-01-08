package panos

import (
	"log"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/dg"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaDeviceGroup,
		Read:   readPanoramaDeviceGroup,
		Update: updatePanoramaDeviceGroup,
		Delete: deletePanoramaDeviceGroup,

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
			"device": {
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

func parsePanoramaDeviceGroup(d *schema.ResourceData) dg.Entry {
	o := dg.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	m := make(map[string][]string)
	dl := d.Get("device").(*schema.Set).List()
	for i := range dl {
		device := dl[i].(map[string]interface{})
		key := device["serial"].(string)
		value := asStringList(device["vsys_list"].(*schema.Set).List())
		m[key] = value
	}
	o.Devices = m

	return o
}

func createPanoramaDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaDeviceGroup(d)

	if err = pano.Panorama.DeviceGroup.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readPanoramaDeviceGroup(d, meta)
}

func readPanoramaDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	o, err := pano.Panorama.DeviceGroup.Get(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	ds := d.Get("device").(*schema.Set)
	s := &schema.Set{F: ds.F}
	for key := range o.Devices {
		sg := make(map[string]interface{})
		sg["serial"] = key
		sg["vsys_list"] = listAsSet(o.Devices[key])
		s.Add(sg)
	}

	d.Set("name", o.Name)
	d.Set("description", o.Description)
	if err = d.Set("device", s); err != nil {
		log.Printf("[WARN] Error setting 'device' field for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	o := parsePanoramaDeviceGroup(d)

	lo, err := pano.Panorama.DeviceGroup.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.DeviceGroup.Edit(lo); err != nil {
		return err
	}

	return readPanoramaDeviceGroup(d, meta)
}

func deletePanoramaDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	name := d.Id()

	err = pano.Panorama.DeviceGroup.Delete(name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
