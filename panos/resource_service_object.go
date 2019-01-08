package panos

import (
	"fmt"
	"log"
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service object's name",
			},
			"vsys": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vsys1",
				ForceNew:    true,
				Description: "The vsys to put this service object in",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object's description",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The protocol (tcp or udp)",
				ValidateFunc: validateStringIn("tcp", "udp"),
			},
			"source_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source port definition",
			},
			"destination_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The destination port definition",
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
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
		Tags:            setAsList(d.Get("tags").(*schema.Set)),
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

func createServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseServiceObject(d)

	if err := fw.Objects.Services.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildServiceObjectId(vsys, o.Name))
	return readServiceObject(d, meta)
}

func readServiceObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceObjectId(d.Id())

	o, err := fw.Objects.Services.Get(vsys, name)
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
	d.Set("protocol", o.Protocol)
	d.Set("source_port", o.SourcePort)
	d.Set("destination_port", o.DestinationPort)
	if err := d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

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
	if err = fw.Objects.Services.Edit(vsys, lo); err != nil {
		return err
	}

	return readServiceObject(d, meta)
}

func deleteServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseServiceObjectId(d.Id())

	err := fw.Objects.Services.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
