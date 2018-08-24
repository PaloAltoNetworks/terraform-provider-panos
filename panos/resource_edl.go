package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/edl"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceEdl() *schema.Resource {
	return &schema.Resource{
		Create: createEdl,
		Read:   readEdl,
		Update: updateEdl,
		Delete: deleteEdl,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vsys": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vsys1",
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      edl.TypeIp,
				ValidateFunc: validateStringIn(edl.TypeIp, edl.TypeDomain, edl.TypeUrl, edl.TypePredefined),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"password_enc": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"repeat": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      edl.RepeatHourly,
				ValidateFunc: validateStringIn(edl.RepeatEveryFiveMinutes, edl.RepeatHourly, edl.RepeatDaily, edl.RepeatWeekly, edl.RepeatMonthly),
			},
			"repeat_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_day_of_week": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringIn("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", ""),
			},
			"repeat_day_of_month": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntInRange(0, 31),
			},
			"exceptions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func parseEdl(d *schema.ResourceData) (string, edl.Entry) {
	vsys := d.Get("vsys").(string)
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

	return vsys, o
}

func parseEdlId(v string) (string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1]
}

func buildEdlId(a, b string) string {
	return fmt.Sprintf("%s%s%s", a, IdSeparator, b)
}

func createEdl(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, o := parseEdl(d)

	if err := fw.Objects.Edl.Set(vsys, o); err != nil {
		return err
	}

	eo, err := fw.Objects.Edl.Get(vsys, o.Name)
	if err != nil {
		return err
	}

	d.SetId(buildEdlId(vsys, o.Name))
	d.Set("password_enc", eo.Password)
	return readEdl(d, meta)
}

func readEdl(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, name := parseEdlId(d.Id())

	o, err := fw.Objects.Edl.Get(vsys, name)
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

func updateEdl(d *schema.ResourceData, meta interface{}) error {
	var err error

	fw := meta.(*pango.Firewall)
	vsys, o := parseEdl(d)

	lo, err := fw.Objects.Edl.Get(vsys, o.Name)
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Objects.Edl.Edit(vsys, lo); err != nil {
		return err
	}
	eo, err := fw.Objects.Edl.Get(vsys, o.Name)
	if err != nil {
		return err
	}

	d.Set("password_enc", eo.Password)
	return readEdl(d, meta)
}

func deleteEdl(d *schema.ResourceData, meta interface{}) error {
	fw := meta.(*pango.Firewall)
	vsys, name := parseEdlId(d.Id())

	err := fw.Objects.Edl.Delete(vsys, name)
	if err != nil {
		e2, ok := err.(pango.PanosError)
		if !ok || !e2.ObjectNotFound() {
			return err
		}
	}
	d.SetId("")
	return nil
}
