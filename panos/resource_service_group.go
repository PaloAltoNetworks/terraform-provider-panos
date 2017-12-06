package panos

import (
	"fmt"
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
                Type: schema.TypeList,
                Required: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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
		Name:        d.Get("name").(string),
        Services: asStringList(d, "services"),
		Tag:         asStringList(d, "tag"),
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

func saveDataServiceGroup(d *schema.ResourceData, vsys string, o srvcgrp.Entry) {
	d.SetId(buildServiceGroupId(vsys, o.Name))
	d.Set("name", o.Name)
    d.Set("services", o.Services)
	d.Set("tag", o.Tag)
}

func createServiceGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceGroup(d)

	if err := fw.Objects.ServiceGroup.Set(vsys, o); err != nil {
		return err
	}

	saveDataServiceGroup(d, vsys, o)
	return nil
}

func readServiceGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceGroupId(d.Id())

	o, err := fw.Objects.ServiceGroup.Get(vsys, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataServiceGroup(d, vsys, o)
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
	err = fw.Objects.ServiceGroup.Edit(vsys, lo)

	if err == nil {
		saveDataServiceGroup(d, vsys, o)
	}
	return err
}

func deleteServiceGroup(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceGroupId(d.Id())

	_ = fw.Objects.ServiceGroup.Delete(vsys, name)
	d.SetId("")
	return nil
}
