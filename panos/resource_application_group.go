package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApplicationGroup() *schema.Resource {
	return &schema.Resource{
		Create: createApplicationGroup,
		Read:   readApplicationGroup,
		Update: updateApplicationGroup,
		Delete: deleteApplicationGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationGroupSchema(false),
	}
}

func applicationGroupSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"applications": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema("vsys1")
	}

	return ans
}

func parseApplicationGroup(d *schema.ResourceData) (string, group.Entry) {
	vsys := d.Get("vsys").(string)
	o := loadApplicationGroup(d)

	return vsys, o
}

func loadApplicationGroup(d *schema.ResourceData) group.Entry {
	return group.Entry{
		Name:         d.Get("name").(string),
		Applications: asStringList(d.Get("applications").([]interface{})),
	}
}

func parseApplicationGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildApplicationGroupId(a, b string) string {
	return strings.Join([]string{a, b}, IdSeparator)
}

func saveApplicationGroup(d *schema.ResourceData, o group.Entry) {
	d.Set("name", o.Name)
	if err := d.Set("applications", o.Applications); err != nil {
		log.Printf("[WARN] Error setting 'applications' for %q: %s", d.Id(), err)
	}
}

func createApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseApplicationGroup(d)

	if err := fw.Objects.AppGroup.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildApplicationGroupId(vsys, o.Name))
	return readApplicationGroup(d, meta)
}

func readApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseApplicationGroupId(d.Id())

	o, err := fw.Objects.AppGroup.Get(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vsys", vsys)
	saveApplicationGroup(d, o)

	return nil
}

func updateApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseApplicationGroup(d)

	lo, err := fw.Objects.AppGroup.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.AppGroup.Edit(vsys, lo); err != nil {
		return err
	}

	return readApplicationGroup(d, meta)
}

func deleteApplicationGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseApplicationGroupId(d.Id())

	err := fw.Objects.AppGroup.Delete(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
