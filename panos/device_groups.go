package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango/pnrm/dg"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceDeviceGroups() *schema.Resource {
	return &schema.Resource{
		Read: readDeviceGroups,

		Schema: listingSchema(),
	}
}

func readDeviceGroups(d *schema.ResourceData, meta interface{}) error {
	con, err := panorama(meta, "")
	if err != nil {
		return err
	}

	listing, err := con.Panorama.DeviceGroup.GetList()
	if err != nil {
		return err
	}

	d.SetId(con.Hostname)
	saveListing(d, listing)
	return nil
}

// Data source.
func dataSourceDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDeviceGroupRead,

		Schema: panoramaDeviceGroupSchema(false),
	}
}

func dataSourceDeviceGroupRead(d *schema.ResourceData, meta interface{}) error {
	con, err := panorama(meta, "")
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	o, err := con.Panorama.DeviceGroup.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(name)
	saveDeviceGroup(d, o)

	return nil
}

// Data source (device group parents).
func dataSourceDeviceGroupParent() *schema.Resource {
	return &schema.Resource{
		Read: readDataSourceDeviceGroupParent,

		Schema: map[string]*schema.Schema{
			"total": {
				Type:        schema.TypeInt,
				Description: "Total number of entries",
				Computed:    true,
			},
			"entries": {
				Type:        schema.TypeMap,
				Description: "Map of strings where the key is the device group and the value is the parent name; an empty string value means that the parent is the 'shared' device group",
				Computed:    true,
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
func resourceDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: createDeviceGroup,
		Read:   readDeviceGroup,
		Update: updateDeviceGroup,
		Delete: deleteDeviceGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: panoramaDeviceGroupSchema(true),
	}
}

func createDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}
	o := loadDeviceGroup(d)

	if err = pano.Panorama.DeviceGroup.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readDeviceGroup(d, meta)
}

func readDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	name := d.Id()

	o, err := pano.Panorama.DeviceGroup.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveDeviceGroup(d, o)

	return nil
}

func updateDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	o := loadDeviceGroup(d)

	lo, err := pano.Panorama.DeviceGroup.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Panorama.DeviceGroup.Edit(lo); err != nil {
		return err
	}

	return readDeviceGroup(d, meta)
}

func deleteDeviceGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	name := d.Id()

	err = pano.Panorama.DeviceGroup.Delete(name)
	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (entry).
func resourceDeviceGroupEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateDeviceGroupEntry,
		Read:   readDeviceGroupEntry,
		Update: createUpdateDeviceGroupEntry,
		Delete: deleteDeviceGroupEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_group": {
				Type:        schema.TypeString,
				Description: "The device group name.",
				Required:    true,
				ForceNew:    true,
			},
			"serial": {
				Type:        schema.TypeString,
				Description: "The NGFW serial number.",
				Required:    true,
				ForceNew:    true,
			},
			"vsys_list": {
				Type:        schema.TypeSet,
				Description: "List of vsys; leave this unspecified if the NGFW is a VM.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func createUpdateDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	dg := d.Get("device_group").(string)
	dev := d.Get("serial").(string)
	vl := setAsList(d.Get("vsys_list").(*schema.Set))

	id := buildDeviceGroupEntryId(dg, dev)

	if err = pano.Panorama.DeviceGroup.EditDeviceVsys(dg, dev, vl); err != nil {
		return err
	}

	d.SetId(id)
	return readDeviceGroupEntry(d, meta)
}

func readDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	dg, dev := parseDeviceGroupEntryId(d.Id())

	d.Set("device_group", dg)
	d.Set("serial", dev)

	o, err := pano.Panorama.DeviceGroup.Get(dg)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	x, ok := o.Devices[dev]
	if !ok {
		d.SetId("")
		return nil
	}

	if err = d.Set("vsys_list", listAsSet(x)); err != nil {
		log.Printf("[WARN] Error setting 'vsys_list' for %q: %s", d.Id(), err)
	}

	return nil
}

func deleteDeviceGroupEntry(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "")
	if err != nil {
		return err
	}

	dg, dev := parseDeviceGroupEntryId(d.Id())

	err = pano.Panorama.DeviceGroup.DeleteDeviceVsys(dg, dev, nil)
	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource (device group parent).
func resourceDeviceGroupParent() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateDeviceGroupParent,
		Read:   readDeviceGroupParent,
		Update: createUpdateDeviceGroupParent,
		Delete: deleteDeviceGroupParent,

		Schema: map[string]*schema.Schema{
			"device_group": {
				Type:        schema.TypeString,
				Description: "The device group name",
				Required:    true,
				ForceNew:    true,
			},
			"parent": {
				Type:        schema.TypeString,
				Description: "The device group's parent",
				Optional:    true,
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

// Schema handling.
func panoramaDeviceGroupSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"device": targetSchema(true),
	}

	if !isResource {
		computed(ans, "", []string{"name"})
	}

	return ans
}

func loadDeviceGroup(d *schema.ResourceData) dg.Entry {
	return dg.Entry{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Devices:     loadTarget(d.Get("device")),
	}
}

func saveDeviceGroup(d *schema.ResourceData, o dg.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	if err := d.Set("device", dumpTarget(o.Devices)); err != nil {
		log.Printf("[WARN] Error setting 'device' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildDeviceGroupEntryId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func parseDeviceGroupEntryId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}
