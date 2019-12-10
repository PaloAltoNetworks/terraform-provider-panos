package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaMonitorProfile() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaMonitorProfile,
		Read:   readPanoramaMonitorProfile,
		Update: updatePanoramaMonitorProfile,
		Delete: deletePanoramaMonitorProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: monitorProfileSchema(true),
	}
}

func parsePanoramaMonitorProfile(d *schema.ResourceData) (string, string, monitor.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	o := loadMonitorProfile(d)

	return tmpl, ts, o
}

func buildPanoramaMonitorProfileId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parsePanoramaMonitorProfileId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func createPanoramaMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaMonitorProfile(d)

	if err := pano.Network.MonitorProfile.Set(tmpl, ts, o); err != nil {
		return err
	}

	d.SetId(buildPanoramaMonitorProfileId(tmpl, ts, o.Name))
	return readPanoramaMonitorProfile(d, meta)
}

func readPanoramaMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaMonitorProfileId(d.Id())

	o, err := pano.Network.MonitorProfile.Get(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	saveMonitorProfile(d, o)

	return nil
}

func updatePanoramaMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, o := parsePanoramaMonitorProfile(d)

	lo, err := pano.Network.MonitorProfile.Get(tmpl, ts, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Network.MonitorProfile.Edit(tmpl, ts, lo); err != nil {
		return err
	}

	return readPanoramaMonitorProfile(d, meta)
}

func deletePanoramaMonitorProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, name := parsePanoramaMonitorProfileId(d.Id())

	err := pano.Network.MonitorProfile.Delete(tmpl, ts, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
