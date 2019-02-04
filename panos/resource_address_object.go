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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address object's name",
			},
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this address object in",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ip-netmask",
				Description:  "The type of address object (ip-netmask, ip-range, fqdn)",
				ValidateFunc: validateStringIn("ip-netmask", "ip-range", "fqdn"),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Administrative tags for the address object",
			},
		},
	}
}

func parseAddressObject(d *schema.ResourceData) (string, addr.Entry) {
	vsys := d.Get("vsys").(string)
	o := addr.Entry{
		Name:        d.Get("name").(string),
		Value:       d.Get("value").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        setAsList(d.Get("tags").(*schema.Set)),
	}

	return vsys, o
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

	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("value", o.Value)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

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
