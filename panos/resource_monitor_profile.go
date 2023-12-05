package panos

import (
	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: createMonitorProfile,
		Read:   readMonitorProfile,
		Update: updateMonitorProfile,
		Delete: deleteMonitorProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: monitorProfileSchema(false),
	}
}

func monitorProfileSchema(p bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"interval": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"threshold": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"action": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      monitor.ActionWaitRecover,
			ValidateFunc: validateStringIn(monitor.ActionWaitRecover, monitor.ActionFailOver),
		},
	}

	if p {
		ans["template"] = templateSchema(true)
		ans["template_stack"] = templateStackSchema()
	}

	return ans
}

func parseMonitorProfile(d *schema.ResourceData) monitor.Entry {
	o := loadMonitorProfile(d)

	return o
}

func loadMonitorProfile(d *schema.ResourceData) monitor.Entry {
	return monitor.Entry{
		Name:      d.Get("name").(string),
		Interval:  d.Get("interval").(int),
		Threshold: d.Get("threshold").(int),
		Action:    d.Get("action").(string),
	}
}

func saveMonitorProfile(d *schema.ResourceData, o monitor.Entry) {
	d.Set("name", o.Name)
	d.Set("interval", o.Interval)
	d.Set("threshold", o.Threshold)
	d.Set("action", o.Action)
}

func createMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	o := parseMonitorProfile(d)

	if err := fw.Network.MonitorProfile.Set(o); err != nil {
		return err
	}

	d.SetId(o.Name)
	return readMonitorProfile(d, meta)
}

func readMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	name := d.Id()

	o, err := fw.Network.MonitorProfile.Get(name)
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveMonitorProfile(d, o)

	return nil
}

func updateMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	o := parseMonitorProfile(d)

	lo, err := fw.Network.MonitorProfile.Get(o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Network.MonitorProfile.Edit(lo); err != nil {
		return err
	}

	return readMonitorProfile(d, meta)
}

func deleteMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	name := d.Id()

	err := fw.Network.MonitorProfile.Delete(name)
	if err != nil {
		if isObjectNotFound(err) {
			return err
		}
	}

	d.SetId("")
	return nil
}
