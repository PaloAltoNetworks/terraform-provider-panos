package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaSyslogServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaSyslogServerProfile,
		Read:   readPanoramaSyslogServerProfile,
		Update: createUpdatePanoramaSyslogServerProfile,
		Delete: deletePanoramaSyslogServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: syslogServerProfileSchema(true),
	}
}

func parsePanoramaSyslogServerProfile(d *schema.ResourceData) (string, string, string, string, syslog.Entry, []server.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	o, serverList := loadSyslogServerProfile(d)

	return tmpl, ts, vsys, dg, o, serverList
}

func parsePanoramaSyslogServerProfileId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaSyslogServerProfileId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createUpdatePanoramaSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, o, serverList := parsePanoramaSyslogServerProfile(d)

	if err := pano.Device.SyslogServerProfile.SetWithoutSubconfig(tmpl, ts, vsys, dg, o); err != nil {
		return err
	}

	if err := pano.Device.SyslogServer.Set(tmpl, ts, vsys, dg, o.Name, serverList...); err != nil {
		return err
	}

	d.SetId(buildPanoramaSyslogServerProfileId(tmpl, ts, vsys, dg, o.Name))
	return readPanoramaSyslogServerProfile(d, meta)
}

func readPanoramaSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaSyslogServerProfileId(d.Id())

	o, err := pano.Device.SyslogServerProfile.Get(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err := pano.Device.SyslogServer.GetList(tmpl, ts, vsys, dg, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := pano.Device.SyslogServer.Get(tmpl, ts, vsys, dg, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	saveSyslogServerProfile(d, o, serverList)

	return nil
}

func deletePanoramaSyslogServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaSyslogServerProfileId(d.Id())

	err := pano.Device.SyslogServerProfile.Delete(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
