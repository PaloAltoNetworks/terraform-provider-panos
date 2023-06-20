package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/fpluchorg/pango/objs/srvc"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServiceObject() *schema.Resource {
	return &schema.Resource{
		Create: createServiceObject,
		Read:   readServiceObject,
		Update: updateServiceObject,
		Delete: deleteServiceObject,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: serviceObjectSchema(false),
	}
}

func serviceObjectSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The service object's name",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Object's description",
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateStringIn(srvc.ProtocolTcp, srvc.ProtocolUdp, srvc.ProtocolSctp),
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
		"tags": tagSchema(),
		"override_session_timeout": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"override_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"override_half_closed_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"override_time_wait_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	if p {
		ans["device_group"] = deviceGroupSchema()
	} else {
		ans["vsys"] = vsysSchema("vsys1")
	}

	return ans
}

func parseServiceObject(d *schema.ResourceData) (string, srvc.Entry) {
	vsys := d.Get("vsys").(string)
	o := loadServiceObject(d)

	return vsys, o
}

func loadServiceObject(d *schema.ResourceData) srvc.Entry {
	return srvc.Entry{
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		Protocol:                  d.Get("protocol").(string),
		SourcePort:                d.Get("source_port").(string),
		DestinationPort:           d.Get("destination_port").(string),
		Tags:                      asStringList(d.Get("tags").([]interface{})),
		OverrideSessionTimeout:    d.Get("override_session_timeout").(bool),
		OverrideTimeout:           d.Get("override_timeout").(int),
		OverrideHalfClosedTimeout: d.Get("override_half_closed_timeout").(int),
		OverrideTimeWaitTimeout:   d.Get("override_time_wait_timeout").(int),
	}
}

func parseServiceObjectId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildServiceObjectId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func saveServiceObject(d *schema.ResourceData, o srvc.Entry) {
	d.Set("name", o.Name)
	d.Set("description", o.Description)
	d.Set("protocol", o.Protocol)
	d.Set("source_port", o.SourcePort)
	d.Set("destination_port", o.DestinationPort)
	if err := d.Set("tags", o.Tags); err != nil {
		log.Printf("[WARN] Error setting 'tags' param for %q: %s", d.Id(), err)
	}
	d.Set("override_session_timeout", o.OverrideSessionTimeout)
	d.Set("override_timeout", o.OverrideTimeout)
	d.Set("override_half_closed_timeout", o.OverrideHalfClosedTimeout)
	d.Set("override_time_wait_timeout", o.OverrideTimeWaitTimeout)
}

func createServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_service_object")
	if err != nil {
		return err
	}
	vsys, o := parseServiceObject(d)

	if err := fw.Objects.Services.Set(vsys, o); err != nil {
		return err
	}

	d.SetId(buildServiceObjectId(vsys, o.Name))
	return readServiceObject(d, meta)
}

func readServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_service_object")
	if err != nil {
		return err
	}
	vsys, name := parseServiceObjectId(d.Id())

	o, err := fw.Objects.Services.Get(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vsys", vsys)
	saveServiceObject(d, o)

	return nil
}

func updateServiceObject(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "panos_panorama_service_object")
	if err != nil {
		return err
	}
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
	fw, err := firewall(meta, "panos_panorama_service_object")
	if err != nil {
		return err
	}
	vsys, name := parseServiceObjectId(d.Id())

	err = fw.Objects.Services.Delete(vsys, name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}
