package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addrgrp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAddressGroup() *schema.Resource {
	return &schema.Resource{
		Create: createAddressGroup,
		Read:   readAddressGroup,
		Update: updateAddressGroup,
		Delete: deleteAddressGroup,

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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Static address group entries",
			},
			"dynamic_match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dynamic address group definition",
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

func parseAddressGroup(d *schema.ResourceData) (string, addrgrp.Entry) {
	vsys := d.Get("vsys").(string)
	o := addrgrp.Entry{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		StaticAddresses: asStringList(d.Get("static_addresses").([]interface{})),
		DynamicMatch:    d.Get("dynamic_match").(string),
		Tags:            setAsList(d.Get("tags").(*schema.Set)),
	}

	return vsys, o
}

func parseAddressGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildAddressGroupId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createAddressGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseAddressGroup(d)

	if err := fw.Objects.AddressGroup.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildAddressGroupId(vsys, o.Name))
	return readAddressGroup(d, meta)
}

func readAddressGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressGroupId(d.Id())

	o, err := fw.Objects.AddressGroup.Get(vsys, name)
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
	d.Set("description", o.Description)
	if err = d.Set("static_addresses", o.StaticAddresses); err != nil {
		log.Printf("[WARN] Error setting 'static_addresses' field for %q: %s", d.Id(), err)
	}
	d.Set("dynamic_match", o.DynamicMatch)
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' field for %q: %s", d.Id(), err)
	}

	return nil
}

func updateAddressGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseAddressGroup(d)

	lo, err := fw.Objects.AddressGroup.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.AddressGroup.Edit(vsys, lo); err != nil {
		return err
	}

	return readAddressGroup(d, meta)
}

func deleteAddressGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseAddressGroupId(d.Id())

	err := fw.Objects.AddressGroup.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
