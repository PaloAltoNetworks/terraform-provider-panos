package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaServiceGroup,
		Read:   readPanoramaServiceGroup,
		Update: updatePanoramaServiceGroup,
		Delete: deletePanoramaServiceGroup,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service group's name",
			},
			"device_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "shared",
				ForceNew:    true,
				Description: "The device group to put this service group in",
			},
			"services": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Administrative tags for the service group",
			},
		},
	}
}

func parsePanoramaServiceGroup(d *schema.ResourceData) (string, srvcgrp.Entry) {
	dg := d.Get("device_group").(string)
	o := srvcgrp.Entry{
		Name:     d.Get("name").(string),
		Services: asStringList(d.Get("services").([]interface{})),
		Tags:     setAsList(d.Get("tags").(*schema.Set)),
	}

	return dg, o
}

func parsePanoramaServiceGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaServiceGroupId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaServiceGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaServiceGroup(d)

	if err := pano.Objects.ServiceGroup.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaServiceGroupId(dg, o.Name))
	return readPanoramaServiceGroup(d, meta)
}

func readPanoramaServiceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaServiceGroupId(d.Id())

	o, err := pano.Objects.ServiceGroup.Get(dg, name)
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
	if err = d.Set("services", o.Services); err != nil {
		log.Printf("[WARN] Error setting 'services' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaServiceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaServiceGroup(d)

	lo, err := pano.Objects.ServiceGroup.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.ServiceGroup.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaServiceGroup(d, meta)
}

func deletePanoramaServiceGroup(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaServiceGroupId(d.Id())

	err := pano.Objects.ServiceGroup.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
