package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"
	"github.com/PaloAltoNetworks/pango/dev/profile/email/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaEmailServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaEmailServerProfile,
		Read:   readPanoramaEmailServerProfile,
		Update: createUpdatePanoramaEmailServerProfile,
		Delete: deletePanoramaEmailServerProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: emailServerProfileSchema(true),
	}
}

func parsePanoramaEmailServerProfile(d *schema.ResourceData) (string, string, string, string, email.Entry, []server.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	o, serverList := loadEmailServerProfile(d)

	return tmpl, ts, vsys, dg, o, serverList
}

func parsePanoramaEmailServerProfileId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaEmailServerProfileId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createUpdatePanoramaEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, o, serverList := parsePanoramaEmailServerProfile(d)

	if err := pano.Device.EmailServerProfile.SetWithoutSubconfig(tmpl, ts, vsys, dg, o); err != nil {
		return err
	}

	if err := pano.Device.EmailServer.Set(tmpl, ts, vsys, dg, o.Name, serverList...); err != nil {
		return err
	}

	d.SetId(buildPanoramaEmailServerProfileId(tmpl, ts, vsys, dg, o.Name))
	return readPanoramaEmailServerProfile(d, meta)
}

func readPanoramaEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaEmailServerProfileId(d.Id())

	o, err := pano.Device.EmailServerProfile.Get(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err := pano.Device.EmailServer.GetList(tmpl, ts, vsys, dg, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := pano.Device.EmailServer.Get(tmpl, ts, vsys, dg, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	d.Set("vsys", vsys)
	saveEmailServerProfile(d, o, serverList)

	return nil
}

func deletePanoramaEmailServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaEmailServerProfileId(d.Id())

	err := pano.Device.EmailServerProfile.Delete(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
