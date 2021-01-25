package panos

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source.
func dataSourceDeviceGroupParent() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceDeviceGroupParent,

		Schema: map[string]*schema.Schema{
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of entries",
			},
			"entries": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of strings where the key is the device group and the value is the parent name; an empty string value means that the parent is the 'shared' device group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func readDataSourceDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	info, err := pano.Panorama.DeviceGroup.GetParents()
	if err != nil {
		return err
	}

	d.SetId(pano.Hostname)
	d.Set("total", len(info))
	if err := d.Set("entries", info); err != nil {
		log.Printf("[WARN] Error setting 'entries' for %q: %s", d.Id(), err)
	}

	return nil
}

// Resource.
func resourceDeviceGroupParent() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateDeviceGroupParent,
		Read:   readDeviceGroupParent,
		Update: createUpdateDeviceGroupParent,
		Delete: deleteDeviceGroupParent,

		Schema: map[string]*schema.Schema{
			"device_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The device group name",
				ForceNew:    true,
			},
			"parent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device group's parent",
			},
		},
	}
}

func createUpdateDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	dg := d.Get("device_group").(string)
	parent := d.Get("parent").(string)

	if err = pano.Panorama.DeviceGroup.AssignParent(dg, parent); err != nil {
		return err
	}

	d.SetId(dg)
	return readDeviceGroupParent(d, meta)
}

func readDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	dg := d.Id()

	info, err := pano.Panorama.DeviceGroup.GetParents()
	if err != nil {
		return err
	}

	parent, ok := info[dg]
	if !ok {
		return fmt.Errorf("Device group %q is not present", dg)
	}

	d.Set("device_group", dg)
	d.Set("parent", parent)

	return nil
}

func deleteDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	dg := d.Id()

	info, err := pano.Panorama.DeviceGroup.GetParents()
	if err != nil {
		return err
	}

	if info[dg] != "" {
		if err = pano.Panorama.DeviceGroup.AssignParent(dg, ""); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
