package panos

import (
	"fmt"
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

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address object's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this address object in",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ip-netmask",
				Description: "The type of address object (ip-netmask, ip-range, fqdn)",
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
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
		Tag:         asStringList(d, "tag"),
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

func saveDataAddressObject(d *schema.ResourceData, vsys string, o addr.Entry) {
	d.SetId(buildAddressObjectId(vsys, o.Name))
	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("value", o.Value)
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	d.Set("tag", o.Tag)
}

func createAddressObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseAddressObject(d)

	if err := fw.Objects.Address.Set(vsys, o); err != nil {
		return err
	}

	saveDataAddressObject(d, vsys, o)
	return nil
}

func readAddressObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressObjectId(d.Id())

	o, err := fw.Objects.Address.Get(vsys, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataAddressObject(d, vsys, o)
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
	err = fw.Objects.Address.Edit(vsys, lo)
	/*
	   if err == nil {
	       lo.Copy(o)
	       err = fw.Objects.Address.Edit(vsys, lo)
	   } else {
	       err = fw.Objects.Address.Set(vsys, o)
	   }
	*/

	if err == nil {
		saveDataAddressObject(d, vsys, o)
	}
	return err
}

func deleteAddressObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressObjectId(d.Id())

	_ = fw.Objects.Address.Delete(vsys, name)
	d.SetId("")
	return nil
}
