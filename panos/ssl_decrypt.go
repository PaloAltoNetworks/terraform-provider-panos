package panos

import (
	"fmt"
	"log"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/ssldecrypt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source.
func dataSourceSslDecrypt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSslDecryptRead,

		Schema: sslDecryptSchema(false),
	}
}

func dataSourceSslDecryptRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ssldecrypt.Config

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildSslDecryptId(tmpl, ts, vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslDecrypt.Get(vsys)
	case *pango.Panorama:
		o, err = con.Device.SslDecrypt.Get(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveSslDecrypt(d, o)

	return nil
}

// Entry resource (exclude certificate).
func resourceSslDecryptExcludeCertificateEntry() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSslDecryptExcludeCertificate,
		Read:   readSslDecryptExcludeCertificate,
		Update: createUpdateSslDecryptExcludeCertificate,
		Delete: deleteSslDecryptExcludeCertificate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vsys":           vsysSchema("shared"),
			"template":       templateSchema(true),
			"template_stack": templateStackSchema(),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name.",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description.",
			},
			"exclude": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Exclude or not.",
				Default:     true,
			},
		},
	}
}

func createUpdateSslDecryptExcludeCertificate(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys := d.Get("template").(string), d.Get("template_stack").(string), d.Get("vsys").(string)
	e := ssldecrypt.SslDecryptExcludeCertificate{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Exclude:     d.Get("exclude").(bool),
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.SetSslDecryptExcludeCertificate(vsys, e)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.SetSslDecryptExcludeCertificate(tmpl, ts, vsys, e)
	}

	if err != nil {
		return err
	}

	id := buildSslDecryptEntryId(tmpl, ts, vsys, e.Name)
	d.SetId(id)
	return readSslDecryptExcludeCertificate(d, meta)
}

func readSslDecryptExcludeCertificate(d *schema.ResourceData, meta interface{}) error {
	var o ssldecrypt.Config
	tmpl, ts, vsys, name, err := parseSslDecryptEntryId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslDecrypt.Get(vsys)
	case *pango.Panorama:
		o, err = con.Device.SslDecrypt.Get(tmpl, ts, vsys)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for _, x := range o.SslDecryptExcludeCertificates {
		if x.Name == name {
			d.Set("template", tmpl)
			d.Set("template_stack", ts)
			d.Set("vsys", vsys)
			d.Set("name", x.Name)
			d.Set("description", x.Description)
			d.Set("exclude", x.Exclude)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteSslDecryptExcludeCertificate(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseSslDecryptEntryId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.DeleteSslDecryptExcludeCertificate(vsys, name)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.DeleteSslDecryptExcludeCertificate(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Entry resource.
func resourceSslDecryptTrustedRootCaEntry() *schema.Resource {
	return &schema.Resource{
		Create: createSslDecryptTrustedRootCaEntry,
		Read:   readSslDecryptTrustedRootCaEntry,
		Delete: deleteSslDecryptTrustedRootCaEntry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"vsys":           vsysSchema("shared"),
			"template":       templateSchema(true),
			"template_stack": templateStackSchema(),
			"certificate_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createSslDecryptTrustedRootCaEntry(d *schema.ResourceData, meta interface{}) error {
	var err error

	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	name := d.Get("certificate_name").(string)

	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("certificate_name", name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.SetTrustedRootCa(vsys, name)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.SetTrustedRootCa(tmpl, ts, vsys, name)
	}

	if err != nil {
		return err
	}

	d.SetId(buildSslDecryptEntryId(tmpl, ts, vsys, name))
	return readSslDecryptTrustedRootCaEntry(d, meta)
}

func readSslDecryptTrustedRootCaEntry(d *schema.ResourceData, meta interface{}) error {
	var o ssldecrypt.Config

	tmpl, ts, vsys, name, err := parseSslDecryptEntryId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslDecrypt.Get(vsys)
	case *pango.Panorama:
		o, err = con.Device.SslDecrypt.Get(tmpl, ts, vsys)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for _, x := range o.TrustedRootCas {
		if x == name {
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteSslDecryptTrustedRootCaEntry(d *schema.ResourceData, meta interface{}) error {
	tmpl, ts, vsys, name, err := parseSslDecryptEntryId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.DeleteTrustedRootCa(vsys, name)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.DeleteTrustedRootCa(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Resource.
func resourceSslDecrypt() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateSslDecrypt,
		Read:   readSslDecrypt,
		Update: createUpdateSslDecrypt,
		Delete: deleteSslDecrypt,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: sslDecryptSchema(true),
	}
}

func createUpdateSslDecrypt(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadSslDecrypt(d)

	vsys := d.Get("vsys").(string)
	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)

	d.Set("vsys", vsys)
	d.Set("template", tmpl)
	d.Set("template_stack", ts)

	id := buildSslDecryptId(tmpl, ts, vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.Edit(vsys, o)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.Edit(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSslDecrypt(d, meta)
}

func readSslDecrypt(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ssldecrypt.Config

	tmpl, ts, vsys := parseSslDecryptId(d.Id())

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslDecrypt.Get(vsys)
	case *pango.Panorama:
		o, err = con.Device.SslDecrypt.Get(tmpl, ts, vsys)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSslDecrypt(d, o)
	return nil
}

func deleteSslDecrypt(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, ts, vsys := parseSslDecryptId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslDecrypt.Delete(vsys)
	case *pango.Panorama:
		err = con.Device.SslDecrypt.Delete(tmpl, ts, vsys)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func sslDecryptSchema(isResource bool) map[string]*schema.Schema {
	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"forward_trust_certificate_rsa": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"forward_trust_certificate_ecdsa": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"forward_untrust_certificate_rsa": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"forward_untrust_certificate_ecdsa": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"root_ca_excludes": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"trusted_root_cas": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"disabled_predefined_exclude_certificates": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"ssl_decrypt_exclude_certificate": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"description": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"exclude": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}

	if !isResource {
		computed(ans, "", []string{"vsys", "template", "template_stack"})
	}

	return ans
}

func loadSslDecrypt(d *schema.ResourceData) ssldecrypt.Config {
	var list []ssldecrypt.SslDecryptExcludeCertificate
	slist := d.Get("ssl_decrypt_exclude_certificate").([]interface{})
	if len(slist) > 0 {
		list = make([]ssldecrypt.SslDecryptExcludeCertificate, 0, len(slist))
		for i := range slist {
			x := slist[i].(map[string]interface{})
			list = append(list, ssldecrypt.SslDecryptExcludeCertificate{
				Name:        x["name"].(string),
				Description: x["description"].(string),
				Exclude:     x["exclude"].(bool),
			})
		}
	}

	return ssldecrypt.Config{
		ForwardTrustCertificateRsa:            d.Get("forward_trust_certificate_rsa").(string),
		ForwardTrustCertificateEcdsa:          d.Get("forward_trust_certificate_ecdsa").(string),
		ForwardUntrustCertificateRsa:          d.Get("forward_untrust_certificate_rsa").(string),
		ForwardUntrustCertificateEcdsa:        d.Get("forward_untrust_certificate_ecdsa").(string),
		RootCaExcludes:                        setAsList(d.Get("root_ca_excludes").(*schema.Set)),
		TrustedRootCas:                        setAsList(d.Get("trusted_root_cas").(*schema.Set)),
		DisabledPredefinedExcludeCertificates: setAsList(d.Get("disabled_predefined_exclude_certificates").(*schema.Set)),
		SslDecryptExcludeCertificates:         list,
	}
}

func saveSslDecrypt(d *schema.ResourceData, o ssldecrypt.Config) {
	var err error

	d.Set("forward_trust_certificate_rsa", o.ForwardTrustCertificateRsa)
	d.Set("forward_trust_certificate_ecdsa", o.ForwardTrustCertificateEcdsa)
	d.Set("forward_untrust_certificate_rsa", o.ForwardUntrustCertificateRsa)
	d.Set("forward_untrust_certificate_ecdsa", o.ForwardUntrustCertificateEcdsa)
	if err = d.Set("root_ca_excludes", listAsSet(o.RootCaExcludes)); err != nil {
		log.Printf("[WARN] Error setting 'root_ca_excludes' for %q: %s", d.Id(), err)
	}
	if err = d.Set("trusted_root_cas", listAsSet(o.TrustedRootCas)); err != nil {
		log.Printf("[WARN] Error setting 'trusted_root_cas' for %q: %s", d.Id(), err)
	}
	if err = d.Set("disabled_predefined_exclude_certificates", listAsSet(o.DisabledPredefinedExcludeCertificates)); err != nil {
		log.Printf("[WARN] Error setting 'disabled_predefined_exclude_certificates' for %q: %s", d.Id(), err)
	}

	var list []interface{}
	if len(o.SslDecryptExcludeCertificates) > 0 {
		list = make([]interface{}, 0, len(o.SslDecryptExcludeCertificates))
		for _, x := range o.SslDecryptExcludeCertificates {
			list = append(list, map[string]interface{}{
				"name":        x.Name,
				"description": x.Description,
				"exclude":     x.Exclude,
			})
		}
	}
	if err = d.Set("ssl_decrypt_exclude_certificate", list); err != nil {
		log.Printf("[WARN] Error setting 'ssl_decrypt_exclude_certificate' for %q: %s", d.Id(), err)
	}
}

// Id functions.
func buildSslDecryptId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseSslDecryptId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}

func buildSslDecryptEntryId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}

func parseSslDecryptEntryId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected id of len() = 4")
	}

	return t[0], t[1], t[2], t[3], nil
}
