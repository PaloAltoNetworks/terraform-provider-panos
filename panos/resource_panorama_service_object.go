package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvc"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaServiceObject() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaServiceObject,
		Read:   readPanoramaServiceObject,
		Update: updatePanoramaServiceObject,
		Delete: deletePanoramaServiceObject,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service object's name",
			},
			"device_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "shared",
				ForceNew:    true,
				Description: "The device group to put this service object in",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Object's description",
			},
			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The protocol (tcp or udp)",
				ValidateFunc: validateStringIn("tcp", "udp"),
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
			"tags": &schema.Schema{
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

func parsePanoramaServiceObject(d *schema.ResourceData) (string, srvc.Entry) {
	dg := d.Get("device_group").(string)
	o := srvc.Entry{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Protocol:        d.Get("protocol").(string),
		SourcePort:      d.Get("source_port").(string),
		DestinationPort: d.Get("destination_port").(string),
		Tags:            setAsList(d.Get("tags").(*schema.Set)),
	}

	return dg, o
}

func parsePanoramaServiceObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaServiceObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaServiceObject(d)

	if err := pano.Objects.Services.Set(dg, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaServiceObjectId(dg, o.Name))
	return readPanoramaServiceObject(d, meta)
}

func readPanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaServiceObjectId(d.Id())

	o, err := pano.Objects.Services.Get(dg, name)
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
	d.Set("description", o.Description)
	d.Set("protocol", o.Protocol)
	d.Set("source_port", o.SourcePort)
	d.Set("destination_port", o.DestinationPort)
	if err := d.Set("tags", listAsSet(o.Tags)); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}

	return nil
}

func updatePanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaServiceObject(d)

	lo, err := pano.Objects.Services.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Services.Edit(dg, lo); err != nil {
		return err
	}

	return readPanoramaServiceObject(d, meta)
}

func deletePanoramaServiceObject(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaServiceObjectId(d.Id())

	err := pano.Objects.Services.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
