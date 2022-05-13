package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango/objs/addrgrp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaAddressGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaAddressGroup,
		Read:   readPanoramaAddressGroup,
		Update: updatePanoramaAddressGroup,
		Delete: deletePanoramaAddressGroup,

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
			"device_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "shared",
				ForceNew:    true,
				Description: "The device group to put this address object in",
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

func parsePanoramaAddressGroup(d *schema.ResourceData) (string, addrgrp.Entry) {
	dg := d.Get("device_group").(string)
	o := addrgrp.Entry{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		StaticAddresses: asStringList(d.Get("static_addresses").([]interface{})),
		DynamicMatch:    d.Get("dynamic_match").(string),
		Tags:            setAsList(d.Get("tags").(*schema.Set)),
	}

	return dg, o
}

func parsePanoramaAddressGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaAddressGroupId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaAddressGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_address_group")
	if err != nil {
		return err
	}
	dg, o := parsePanoramaAddressGroup(d)

	if err = pano.Objects.AddressGroup.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaAddressGroupId(dg, o.Name))
	return readPanoramaAddressGroup(d, meta)
}

func readPanoramaAddressGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_address_group")
	if err != nil {
		return err
	}

	dg, name := parsePanoramaAddressGroupId(d.Id())

	o, err := pano.Objects.AddressGroup.Get(dg, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", o.Name)
	d.Set("device_group", dg)
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

func updatePanoramaAddressGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_address_group")
	if err != nil {
		return err
	}

	dg, o := parsePanoramaAddressGroup(d)

	lo, err := pano.Objects.AddressGroup.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.AddressGroup.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaAddressGroup(d, meta)
}

func deletePanoramaAddressGroup(d *schema.ResourceData, meta interface{}) error {
	pano, err := panorama(meta, "panos_address_group")
	if err != nil {
		return err
	}

	dg, name := parsePanoramaAddressGroupId(d.Id())

	if err = pano.Objects.AddressGroup.Delete(dg, name); err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}
