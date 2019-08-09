package panos

import (
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/signature"
	"github.com/PaloAltoNetworks/pango/objs/app/signature/andcond"
	"github.com/PaloAltoNetworks/pango/objs/app/signature/orcond"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaApplicationSignature() *schema.Resource {
	return &schema.Resource{
		Create: createUpdatePanoramaApplicationSignature,
		Read:   readPanoramaApplicationSignature,
		Update: createUpdatePanoramaApplicationSignature,
		Delete: deletePanoramaApplicationSignature,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: applicationSignatureSchema(true),
	}
}

func parsePanoramaApplicationSignature(d *schema.ResourceData) (string, string, signature.Entry, []andcond.Entry, map[string][]orcond.Entry) {
	dg := d.Get("device_group").(string)
	ao := d.Get("application_object").(string)
	o, andList, orMap := loadApplicationSignature(d)

	return dg, ao, o, andList, orMap
}

func parsePanoramaApplicationSignatureId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildPanoramaApplicationSignatureId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func createUpdatePanoramaApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, app, o, andList, orMap := parsePanoramaApplicationSignature(d)

	if err := pano.Objects.AppSignature.Edit(dg, app, o); err != nil {
		return err
	}

	if err := pano.Objects.AppSigAndCond.Set(dg, app, o.Name, andList...); err != nil {
		return err
	}

	for i := range andList {
		orList := orMap[andList[i].Name]
		if err := pano.Objects.AppSigOrCond.Set(dg, app, o.Name, andList[i].Name, orList...); err != nil {
			return err
		}
	}

	d.SetId(buildPanoramaApplicationSignatureId(dg, app, o.Name))
	return readPanoramaApplicationSignature(d, meta)
}

func readPanoramaApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, app, name := parsePanoramaApplicationSignatureId(d.Id())

	o, err := pano.Objects.AppSignature.Get(dg, app, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if ok && e2.ObjectNotFound() {
			d.SetId("")
			return nil
		}
		return err
	}

	andNames, err := pano.Objects.AppSigAndCond.GetList(dg, app, name)
	andList := make([]andcond.Entry, 0, len(andNames))
	orMap := make(map[string][]orcond.Entry)
	for _, andName := range andNames {
		andEntry, err := pano.Objects.AppSigAndCond.Get(dg, app, name, andName)
		if err != nil {
			return err
		}
		andList = append(andList, andEntry)
		orNames, err := pano.Objects.AppSigOrCond.GetList(dg, app, name, andName)
		orList := make([]orcond.Entry, 0, len(orNames))
		if err != nil {
			return err
		}
		for _, orName := range orNames {
			orEntry, err := pano.Objects.AppSigOrCond.Get(dg, app, name, andName, orName)
			if err != nil {
				return err
			}
			orList = append(orList, orEntry)
		}
		orMap[andEntry.Name] = orList
	}

	d.Set("device_group", dg)
	d.Set("application_object", app)
	saveApplicationSignature(d, o, andList, orMap)

	return nil
}

func deletePanoramaApplicationSignature(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, app, name := parsePanoramaApplicationSignatureId(d.Id())

	err := pano.Objects.AppSignature.Delete(dg, app, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}

	d.SetId("")
	return nil
}
