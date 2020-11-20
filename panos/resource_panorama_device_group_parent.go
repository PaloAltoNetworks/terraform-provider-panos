package panos

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaDeviceGroupParent() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaDeviceGroupParent,
		Read:   readPanoramaDeviceGroupParent,
		Update: createUpdatePanoramaDeviceGroupParent,
		Delete: deletePanoramaDeviceGroupParent,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The device group whose parent to configure",
			},
			"parent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent device group",
			},
		},
	}
}

func parsePanoramaDeviceGroupParent(d *schema.ResourceData) (string, string) {
	dg := d.Get("device_group").(string)
	parent := d.Get("parent").(string)

	return dg, parent
}

func createUpdatePanoramaDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	dg, parent := parsePanoramaDeviceGroupParent(d)

	if err := pano.AssignDeviceGroupParent(dg, parent); err != nil {
		return err
	}

	d.SetId(dg)
	return readPanoramaDeviceGroupParent(d, meta)
}

func readPanoramaDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	dg := d.Id()

	info, err := pano.DeviceGroupHierarchy()
	if err != nil {
		return err
	}

	d.Set("device_group", dg)
	d.Set("parent", info[dg])
	return nil
}

func deletePanoramaDeviceGroupParent(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	dg := d.Id()

	_ = pano.AssignDeviceGroupParent(dg, "")

	d.SetId("")
	return nil
}
