package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/header"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/param"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePanoramaHttpServerProfile() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaHttpServerProfile,
		Read:   readPanoramaHttpServerProfile,
		Update: createUpdatePanoramaHttpServerProfile,
		Delete: deletePanoramaHttpServerProfile,

		Schema: httpServerProfileSchema(true),
	}
}

func parsePanoramaHttpServerProfile(d *schema.ResourceData) (string, string, string, string, http.Entry, []server.Entry, map[string][]header.Entry, map[string][]param.Entry) {
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	dg := d.Get("device_group").(string)
	o, serverList, headers, params := loadHttpServerProfile(d)

	return tmpl, ts, vsys, dg, o, serverList, headers, params
}

func parsePanoramaHttpServerProfileId(v string) (string, string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3], t[4]
}

func buildPanoramaHttpServerProfileId(a, b, c, d, e string) string {
	return strings.Join([]string{a, b, c, d, e}, IdSeparator)
}

func createUpdatePanoramaHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, o, serverList, headers, params := parsePanoramaHttpServerProfile(d)

	if err := pano.Device.HttpServerProfile.SetWithoutSubconfig(tmpl, ts, vsys, dg, o); err != nil {
		return err
	}

	if err := pano.Device.HttpServer.Set(tmpl, ts, vsys, dg, o.Name, serverList...); err != nil {
		return err
	}

	password_enc := make(map[string]interface{})
	password_raw := make(map[string]interface{})
	for _, x := range serverList {
		live, err := pano.Device.HttpServer.Get(tmpl, ts, vsys, dg, o.Name, x.Name)
		if err != nil {
			return err
		}
		password_enc[x.Name] = live.Password
		password_raw[x.Name] = x.Password
	}
	if err := d.Set("password_enc", password_enc); err != nil {
		log.Printf("[WARN] Error setting 'password_enc' for %q: %s", d.Id(), err)
	}
	if err := d.Set("password_raw", password_raw); err != nil {
		log.Printf("[WARN] Error setting 'password_raw' for %q: %s", d.Id(), err)
	}

	for logtype := range headers {
		headerList := headers[logtype]
		if len(headerList) != 0 {
			if err := pano.Device.HttpHeader.Set(tmpl, ts, vsys, dg, o.Name, logtype, headerList...); err != nil {
				return err
			}
		}
	}

	for logtype := range params {
		paramList := params[logtype]
		if len(paramList) != 0 {
			if err := pano.Device.HttpParam.Set(tmpl, ts, vsys, dg, o.Name, logtype, paramList...); err != nil {
				return err
			}
		}
	}

	d.SetId(buildPanoramaHttpServerProfileId(tmpl, ts, vsys, dg, o.Name))
	return readPanoramaHttpServerProfile(d, meta)
}

func readPanoramaHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var list []string

	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaHttpServerProfileId(d.Id())

	o, err := pano.Device.HttpServerProfile.Get(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	list, err = pano.Device.HttpServer.GetList(tmpl, ts, vsys, dg, name)
	if err != nil {
		return err
	}
	serverList := make([]server.Entry, 0, len(list))
	for i := range list {
		entry, err := pano.Device.HttpServer.Get(tmpl, ts, vsys, dg, name, list[i])
		if err != nil {
			return err
		}
		serverList = append(serverList, entry)
	}

	logtypes := []string{
		param.Config,
		param.System,
		param.Threat,
		param.Traffic,
		param.HipMatch,
		param.Url,
		param.Data,
		param.Wildfire,
		param.Tunnel,
		param.UserId,
		param.Gtp,
		param.Auth,
		param.Sctp,
		param.Iptag,
	}

	headers := make(map[string][]header.Entry)
	params := make(map[string][]param.Entry)

	for _, logtype := range logtypes {
		list, err = pano.Device.HttpHeader.GetList(tmpl, ts, vsys, dg, name, logtype)
		if err != nil {
			return err
		}
		if len(list) != 0 {
			headerList := make([]header.Entry, 0, len(list))
			for _, hdr := range list {
				entry, err := pano.Device.HttpHeader.Get(tmpl, ts, vsys, dg, name, logtype, hdr)
				if err != nil {
					return err
				}
				headerList = append(headerList, entry)
			}
			headers[logtype] = headerList
		}

		list, err = pano.Device.HttpParam.GetList(tmpl, ts, vsys, dg, name, logtype)
		if err != nil {
			return err
		}
		if len(list) != 0 {
			paramList := make([]param.Entry, 0, len(list))
			for _, prm := range list {
				entry, err := pano.Device.HttpParam.Get(tmpl, ts, vsys, dg, name, logtype, prm)
				if err != nil {
					return err
				}
				paramList = append(paramList, entry)
			}
			params[logtype] = paramList
		}
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)
	d.Set("device_group", dg)
	saveHttpServerProfile(d, o, serverList, headers, params)

	return nil
}

func deletePanoramaHttpServerProfile(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	tmpl, ts, vsys, dg, name := parsePanoramaHttpServerProfileId(d.Id())

	err := pano.Device.HttpServerProfile.Delete(tmpl, ts, vsys, dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
