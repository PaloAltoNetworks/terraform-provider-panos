package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAddressObject() *schema.Resource {
	return &schema.Resource{
		Create: createAddressObject,
		Read:   readAddressObject,
		Update: updateAddressObject,
		Delete: deleteAddressObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: addressObjectSchema(false),
	}
}

func addressObjectSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      addr.IpNetmask,
			ValidateFunc: validateStringIn(addr.IpNetmask, addr.IpRange, addr.Fqdn, addr.IpWildcard),
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tags": tagSchema(),
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema()
	}

	return ans
}

func parseAddressObject(d *schema.ResourceData) (string, addr.Entry) {
	vsys := d.Get("vsys").(string)
	o := loadAddressObject(d)

	return vsys, o
}

func loadAddressObject(d *schema.ResourceData) addr.Entry {
	return addr.Entry{
		Name:        d.Get("name").(string),
		Value:       d.Get("value").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        asStringList(d.Get("tags").([]interface{})),
	}
}

func saveAddressObject(d *schema.ResourceData, o addr.Entry) {
	d.Set("name", o.Name)
	d.Set("type", o.Type)
	d.Set("value", o.Value)
	d.Set("description", o.Description)
	if err := d.Set("tags", o.Tags); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
}

func parseAddressObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildAddressObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createAddressObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseAddressObject(d)

	if err := fw.Objects.Address.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildAddressObjectId(vsys, o.Name))
	return readAddressObject(d, meta)
}

func readAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressObjectId(d.Id())

	o, err := fw.Objects.Address.Get(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vsys", vsys)
	saveAddressObject(d, o)

	return nil
}

func updateAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseAddressObject(d)

	lo, err := fw.Objects.Address.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.Address.Edit(vsys, lo); err != nil {
		return err
	}

	return readAddressObject(d, meta)
}

func deleteAddressObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressObjectId(d.Id())

	err := fw.Objects.Address.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
