package panos

import (
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	cert "github.com/PaloAltoNetworks/pango/dev/profile/certificate"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source (listing).
func dataSourceCertificateProfiles() *schema.Resource {
	s := listingSchema()
	for key, val := range templateWithPanoramaSharedSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceCertificateProfilesRead,

		Schema: s,
	}
}

func dataSourceCertificateProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildCertificateProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.CertificateProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.CertificateProfile.GetList(false, tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceCertificateProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCertificateProfileRead,

		Schema: certificateProfileSchema(false),
	}
}

func dataSourceCertificateProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o cert.Entry

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildCertificateProfileId(tmpl, ts, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.CertificateProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.CertificateProfile.Get(false, tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveCertificateProfile(d, o)

	return nil
}

// Resource.
func resourceCertificateProfile() *schema.Resource {
	return &schema.Resource{
		Create: createCertificateProfile,
		Read:   readCertificateProfile,
		Update: updateCertificateProfile,
		Delete: deleteCertificateProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: certificateProfileSchema(true),
	}
}

func createCertificateProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadCertificateProfile(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildCertificateProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.CertificateProfile.Set(vsys, o)
	case *pango.Panorama:
		err = con.Device.CertificateProfile.Set(false, tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readCertificateProfile(d, meta)
}

func readCertificateProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o cert.Entry

	tmpl, ts, vsys, name := parseCertificateProfileId(d.Id())

	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.CertificateProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.CertificateProfile.Get(false, tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveCertificateProfile(d, o)

	return nil
}

func updateCertificateProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo cert.Entry
	o := loadCertificateProfile(d)

	tmpl, ts, vsys, name := parseCertificateProfileId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err = con.Device.CertificateProfile.Get(vsys, name)
		if err == nil {
			lo.Copy(o)
			err = con.Device.CertificateProfile.Edit(vsys, lo)
		}
	case *pango.Panorama:
		lo, err = con.Device.CertificateProfile.Get(false, tmpl, ts, vsys, name)
		if err == nil {
			lo.Copy(o)
			err = con.Device.CertificateProfile.Edit(false, tmpl, ts, vsys, lo)
		}
	}

	if err != nil {
		return err
	}

	return readCertificateProfile(d, meta)
}

func deleteCertificateProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys, name := parseCertificateProfileId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.CertificateProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.CertificateProfile.Delete(false, tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func certificateProfileSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"username_field": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringIn("", cert.UsernameFieldSubject, cert.UsernameFieldSubjectAlt),
		},
		"username_field_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"use_crl": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"use_ocsp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"crl_receive_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5,
		},
		"ocsp_receive_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5,
		},
		"certificate_status_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5,
		},
		"block_unknown_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"block_certificate_timeout": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"block_unauthenticated_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"block_expired_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"ocsp_exclude_nonce": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"certificate": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"default_ocsp_url": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"ocsp_verify_certificate": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"template_name": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func loadCertificateProfile(d *schema.ResourceData) cert.Entry {
	var list []cert.Certificate
	sl := d.Get("certificate").([]interface{})
	if len(sl) > 0 {
		list = make([]cert.Certificate, 0, len(sl))
		for i := range sl {
			x := sl[i].(map[string]interface{})
			list = append(list, cert.Certificate{
				Name:                  x["name"].(string),
				DefaultOcspUrl:        x["default_ocsp_url"].(string),
				OcspVerifyCertificate: x["ocsp_verify_certificate"].(string),
				TemplateName:          x["template_name"].(string),
			})
		}
	}

	return cert.Entry{
		Name:                            d.Get("name").(string),
		UsernameField:                   d.Get("username_field").(string),
		UsernameFieldValue:              d.Get("username_field_value").(string),
		Domain:                          d.Get("domain").(string),
		Certificates:                    list,
		UseCrl:                          d.Get("use_crl").(bool),
		UseOcsp:                         d.Get("use_ocsp").(bool),
		CrlReceiveTimeout:               d.Get("crl_receive_timeout").(int),
		OcspReceiveTimeout:              d.Get("ocsp_receive_timeout").(int),
		CertificateStatusTimeout:        d.Get("certificate_status_timeout").(int),
		BlockUnknownCertificate:         d.Get("block_unknown_certificate").(bool),
		BlockCertificateTimeout:         d.Get("block_certificate_timeout").(bool),
		BlockUnauthenticatedCertificate: d.Get("block_unauthenticated_certificate").(bool),
		BlockExpiredCertificate:         d.Get("block_expired_certificate").(bool),
		OcspExcludeNonce:                d.Get("ocsp_exclude_nonce").(bool),
	}
}

func saveCertificateProfile(d *schema.ResourceData, o cert.Entry) {
	d.Set("name", o.Name)
	d.Set("username_field", o.UsernameField)
	d.Set("username_field_value", o.UsernameFieldValue)
	d.Set("domain", o.Domain)
	d.Set("use_crl", o.UseCrl)
	d.Set("use_ocsp", o.UseOcsp)
	d.Set("crl_receive_timeout", o.CrlReceiveTimeout)
	d.Set("ocsp_receive_timeout", o.OcspReceiveTimeout)
	d.Set("certificate_status_timeout", o.CertificateStatusTimeout)
	d.Set("block_unknown_certificate", o.BlockUnknownCertificate)
	d.Set("block_certificate_timeout", o.BlockCertificateTimeout)
	d.Set("block_unauthenticated_certificate", o.BlockUnauthenticatedCertificate)
	d.Set("block_expired_certificate", o.BlockExpiredCertificate)
	d.Set("ocsp_exclude_nonce", o.OcspExcludeNonce)

	var list []interface{}
	if len(o.Certificates) > 0 {
		list = make([]interface{}, 0, len(o.Certificates))
		for _, x := range o.Certificates {
			list = append(list, map[string]interface{}{
				"name":                    x.Name,
				"default_ocsp_url":        x.DefaultOcspUrl,
				"ocsp_verify_certificate": x.OcspVerifyCertificate,
				"template_name":           x.TemplateName,
			})
		}
	}

	if err := d.Set("certificate", list); err != nil {
		log.Printf("[WARN] Error setting 'certificate' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func parseCertificateProfileId(v string) (string, string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2], t[3]
}

func buildCertificateProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
