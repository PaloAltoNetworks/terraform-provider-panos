package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/edl"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePanoramaEdl() *schema.Resource {
	return &schema.Resource{
		Create: createPanoramaEdl,
		Read:   readPanoramaEdl,
		Update: updatePanoramaEdl,
		Delete: deletePanoramaEdl,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "shared",
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      edl.TypeIp,
				ValidateFunc: validateStringIn(edl.TypeIp, edl.TypeDomain, edl.TypeUrl, edl.TypePredefined),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"password_enc": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"repeat": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      edl.RepeatHourly,
				ValidateFunc: validateStringIn(edl.RepeatEveryFiveMinutes, edl.RepeatHourly, edl.RepeatDaily, edl.RepeatWeekly, edl.RepeatMonthly),
			},
			"repeat_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_day_of_week": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", ""),
			},
			"repeat_day_of_month": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntInRange(0, 31),
			},
			"exceptions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parsePanoramaEdl(d *schema.ResourceData) (string, edl.Entry) {
	dg := d.Get("device_group").(string)
	o := edl.Entry{
		Name:               d.Get("name").(string),
		Type:               d.Get("type").(string),
		Description:        d.Get("description").(string),
		Source:             d.Get("source").(string),
		CertificateProfile: d.Get("certificate_profile").(string),
		Username:           d.Get("username").(string),
		Password:           d.Get("password").(string),
		Repeat:             d.Get("repeat").(string),
		RepeatAt:           d.Get("repeat_at").(string),
		RepeatDayOfWeek:    d.Get("repeat_day_of_week").(string),
		RepeatDayOfMonth:   d.Get("repeat_day_of_month").(int),
		Exceptions:         asStringList(d.Get("exceptions").([]interface{})),
	}

	return dg, o
}

func parsePanoramaEdlId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildPanoramaEdlId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createPanoramaEdl(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaEdl(d)

	if err := pano.Objects.Edl.Set(dg, o); err != nil {
		return err
	}

	eo, err := pano.Objects.Edl.Get(dg, o.Name)
	if err != nil {
		return err
	}

	d.SetId(buildPanoramaEdlId(dg, o.Name))
	d.Set("password_enc", eo.Password)
	return readPanoramaEdl(d, meta)
}

func readPanoramaEdl(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaEdlId(d.Id())

	o, err := pano.Objects.Edl.Get(dg, name)
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
	d.Set("type", o.Type)
	d.Set("description", o.Description)
	d.Set("source", o.Source)
	d.Set("certificate_profile", o.CertificateProfile)
	d.Set("username", o.Username)
	d.Set("repeat", o.Repeat)
	d.Set("repeat_at", o.RepeatAt)
	d.Set("repeat_day_of_week", o.RepeatDayOfWeek)
	d.Set("repeat_day_of_month", o.RepeatDayOfMonth)
	d.Set("exceptions", o.Exceptions)

	if d.Get("password_enc").(string) != o.Password {
		d.Set("password", "")
	}

	return nil
}

func updatePanoramaEdl(d *schema.ResourceData, meta interface{}) error {
	var err error

	pano := meta.(*pango.Panorama)
	dg, o := parsePanoramaEdl(d)

	lo, err := pano.Objects.Edl.Get(dg, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = pano.Objects.Edl.Edit(dg, lo); err != nil {
		return err
	}
	eo, err := pano.Objects.Edl.Get(dg, o.Name)
	if err != nil {
		return err
	}

	d.Set("password_enc", eo.Password)
	return readPanoramaEdl(d, meta)
}

func deletePanoramaEdl(d *schema.ResourceData, meta interface{}) error {
	pano := meta.(*pango.Panorama)
	dg, name := parsePanoramaEdlId(d.Id())

	err := pano.Objects.Edl.Delete(dg, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
