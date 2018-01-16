package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServiceGroup() *schema.Resource {
	return &schema.Resource{
		Create: createServiceGroup,
		Read:   readServiceGroup,
		Update: updateServiceGroup,
		Delete: deleteServiceGroup,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service group's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this service group in",
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": &schema.Schema{
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

func parseServiceGroup(d *schema.ResourceData) (string, srvcgrp.Entry) {
	vsys := d.Get("vsys").(string)
	o := srvcgrp.Entry{
		Name:     d.Get("name").(string),
		Services: asStringList(d.Get("services").([]interface{})),
		Tags:     setAsList(d.Get("tags").(*schema.Set)),
	}

	return vsys, o
}

func parseServiceGroupId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildServiceGroupId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createServiceGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceGroup(d)

	if err := fw.Objects.ServiceGroup.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildServiceGroupId(vsys, o.Name))
	return readServiceGroup(d, meta)
}

func readServiceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceGroupId(d.Id())

	o, err := fw.Objects.ServiceGroup.Get(vsys, name)
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
	if err = d.Set("services", o.Services); err != nil {
		log.Printf("[WARN] Error setting 'services' param for %q: %s", d.Id(), err)
	}
	if err = d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func updateServiceGroup(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceGroup(d)

	lo, err := fw.Objects.ServiceGroup.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.ServiceGroup.Edit(vsys, lo); err != nil {
		return err
	}

	return readServiceGroup(d, meta)
}

func deleteServiceGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceGroupId(d.Id())

	err := fw.Objects.ServiceGroup.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
