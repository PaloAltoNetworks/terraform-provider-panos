package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaAddressObject() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaAddressObject,
		Read:   readPanoramaAddressObject,
		Update: updatePanoramaAddressObject,
		Delete: deletePanoramaAddressObject,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address object's name",
			},
			"device_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ip-netmask",
				Description:  "The type of address object (ip-netmask, ip-range, fqdn)",
				ValidateFunc: validateStringIn("ip-netmask", "ip-range", "fqdn"),
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": &schema.Schema{
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

func parsePanoramaAddressObject(d *schema.ResourceData) (string, addr.Entry) {
	dg := d.Get("device_group").(string)
	o := addr.Entry{
		Name:        d.Get("name").(string),
		Value:       d.Get("value").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        setAsList(d.Get("tags").(*schema.Set)),
	}

	return dg, o
}

func parsePanoramaAddressObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaAddressObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaAddressObject(d)

	if err := pano.Objects.Address.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaAddressObjectId(dg, o.Name))
	return readPanoramaAddressObject(d, meta)
}

func readPanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaAddressObjectId(d.Id())

	o, err := pano.Objects.Address.Get(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("device_group", dg)
	d.Set("value", o.Value)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaAddressObject(d)

	lo, err := pano.Objects.Address.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Address.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaAddressObject(d, meta)
}

func deletePanoramaAddressObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaAddressObjectId(d.Id())

	err := pano.Objects.Address.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
