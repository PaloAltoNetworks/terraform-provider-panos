package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvc"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServiceObject() *schema.Resource {
	return &schema.Resource{
		Create: createServiceObject,
		Read:   readServiceObject,
		Update: updateServiceObject,
		Delete: deleteServiceObject,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service object's name",
			},
			"vsys": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this service object in",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object's description",
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol (tcp or udp)",
			},
			"source_port": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source port definition",
			},
			"destination_port": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The destination port definition",
			},
			"tag": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Administrative tags for the service object",
			},
		},
	}
}

func parseServiceObject(d *schema.ResourceData) (string, srvc.Entry) {
	vsys := d.Get("vsys").(string)
	o := srvc.Entry{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Protocol:        d.Get("protocol").(string),
		SourcePort:      d.Get("source_port").(string),
		DestinationPort: d.Get("destination_port").(string),
		Tag:             asStringList(d, "tag"),
	}

	return vsys, o
}

func parseServiceObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildServiceObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func saveDataServiceObject(d *schema.ResourceData, vsys string, o srvc.Entry) {
	d.SetId(buildServiceObjectId(vsys, o.Name))
	d.Set("name", o.Name)
	d.Set("vsys", vsys)
	d.Set("description", o.Description)
	d.Set("protocol", o.Protocol)
	d.Set("source_port", o.SourcePort)
	d.Set("destination_port", o.DestinationPort)
	d.Set("tag", o.Tag)
}

func createServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceObject(d)

	if err := fw.Objects.Services.Set(vsys, o); err != nil {
		return err
	}

	saveDataServiceObject(d, vsys, o)
	return nil
}

func readServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceObjectId(d.Id())

	o, err := fw.Objects.Services.Get(vsys, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	saveDataServiceObject(d, vsys, o)
	return nil
}

func updateServiceObject(d *schema.ResourceData, meta interface{}) error {
	var err error
	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceObject(d)

	lo, err := fw.Objects.Services.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	err = fw.Objects.Services.Edit(vsys, lo)

	if err == nil {
		saveDataServiceObject(d, vsys, o)
	}
	return err
}

func deleteServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceObjectId(d.Id())

	_ = fw.Objects.Services.Delete(vsys, name)
	d.SetId("")
	return nil
}
