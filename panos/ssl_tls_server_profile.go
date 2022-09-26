package panos

import (
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/ssltls"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Data source (listing).
func dataSourceSslTlsServiceProfiles() *schema.Resource {
	s := listingSchema()
	for key, val := range templateWithPanoramaSharedSchema() {
		s[key] = val
	}

	return &schema.Resource{
		Read: dataSourceSslTlsServiceProfilesRead,

		Schema: s,
	}
}

func dataSourceSslTlsServiceProfilesRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var listing []string

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildSslTlsServiceProfileId(tmpl, ts, vsys, "")

	switch con := meta.(type) {
	case *pango.Firewall:
		listing, err = con.Device.SslTlsServiceProfile.GetList(vsys)
	case *pango.Panorama:
		listing, err = con.Device.SslTlsServiceProfile.GetList(tmpl, ts, vsys)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	saveListing(d, listing)

	return nil
}

// Data source.
func dataSourceSslTlsServiceProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSslTlsServiceProfileRead,

		Schema: sslTlsServiceProfileSchema(false),
	}
}

func dataSourceSslTlsServiceProfileRead(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ssltls.Entry

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)
	name := d.Get("name").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildSslTlsServiceProfileId(tmpl, ts, vsys, name)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslTlsServiceProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.SslTlsServiceProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(id)
	saveSslTlsServiceProfile(d, o)

	return nil
}

// Resource.
func resourceSslTlsServiceProfile() *schema.Resource {
	return &schema.Resource{
		Create: createSslTlsServiceProfile,
		Read:   readSslTlsServiceProfile,
		Update: updateSslTlsServiceProfile,
		Delete: deleteSslTlsServiceProfile,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: sslTlsServiceProfileSchema(true),
	}
}

func createSslTlsServiceProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	o := loadSslTlsServiceProfile(d)

	tmpl := d.Get("template").(string)
	ts := d.Get("template_stack").(string)
	vsys := d.Get("vsys").(string)

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	id := buildSslTlsServiceProfileId(tmpl, ts, vsys, o.Name)

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslTlsServiceProfile.Set(vsys, o)
	case *pango.Panorama:
		err = con.Device.SslTlsServiceProfile.Set(tmpl, ts, vsys, o)
	}

	if err != nil {
		return err
	}

	d.SetId(id)
	return readSslTlsServiceProfile(d, meta)
}

func readSslTlsServiceProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o ssltls.Entry

	tmpl, ts, vsys, name, err := parseSslTlsServiceProfileId(d.Id())
	if err != nil {
		return err
	}

	d.Set("template", tmpl)
	d.Set("template_stack", ts)
	d.Set("vsys", vsys)

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.SslTlsServiceProfile.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.SslTlsServiceProfile.Get(tmpl, ts, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	saveSslTlsServiceProfile(d, o)

	return nil
}

func updateSslTlsServiceProfile(d *schema.ResourceData, meta interface{}) error {
	var err error
	var lo ssltls.Entry
	o := loadSslTlsServiceProfile(d)

	tmpl, ts, vsys, _, err := parseSslTlsServiceProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		lo, err = con.Device.SslTlsServiceProfile.Get(vsys, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Device.SslTlsServiceProfile.Edit(vsys, lo)
		}
	case *pango.Panorama:
		lo, err = con.Device.SslTlsServiceProfile.Get(tmpl, ts, vsys, o.Name)
		if err == nil {
			lo.Copy(o)
			err = con.Device.SslTlsServiceProfile.Edit(tmpl, ts, vsys, lo)
		}
	}

	if err != nil {
		return err
	}

	return readSslTlsServiceProfile(d, meta)
}

func deleteSslTlsServiceProfile(d *schema.ResourceData, meta interface{}) error {
	var err error

	tmpl, ts, vsys, name, err := parseSslTlsServiceProfileId(d.Id())
	if err != nil {
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.SslTlsServiceProfile.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.SslTlsServiceProfile.Delete(tmpl, ts, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema handling.
func sslTlsServiceProfileSchema(isResource bool) map[string]*schema.Schema {
	mins := []string{ssltls.Tls1_0, ssltls.Tls1_1, ssltls.Tls1_2}
	maxes := []string{ssltls.Tls1_0, ssltls.Tls1_1, ssltls.Tls1_2, ssltls.TlsMax}

	ans := map[string]*schema.Schema{
		"template":       templateSchema(true),
		"template_stack": templateStackSchema(),
		"vsys":           vsysSchema("shared"),
		"name": {
			Type:        schema.TypeString,
			Description: "The object name.",
			Required:    true,
			ForceNew:    true,
		},
		"certificate": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SSL certificate file name.",
		},
		"min_version": {
			Type:     schema.TypeString,
			Optional: true,
			Description: addStringInSliceValidation(
				"Minimum TLS protocol version.",
				mins,
			),
			Default:      ssltls.Tls1_0,
			ValidateFunc: validateStringIn(mins...),
		},
		"max_version": {
			Type:     schema.TypeString,
			Optional: true,
			Description: addStringInSliceValidation(
				"Maximum TLS protocol version.",
				maxes,
			),
			Default:      ssltls.TlsMax,
			ValidateFunc: validateStringIn(maxes...),
		},
		"allow_algorithm_rsa": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm RSA.",
			Default:     true,
		},
		"allow_algorithm_dhe": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm DHE.",
			Default:     true,
		},
		"allow_algorithm_ecdhe": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm ECDHE.",
			Default:     true,
		},
		"allow_algorithm_3des": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm 3DES.",
			Default:     true,
		},
		"allow_algorithm_rc4": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm RC4.",
			Default:     true,
		},
		"allow_algorithm_aes_128_cbc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm AES-128-CBC.",
			Default:     true,
		},
		"allow_algorithm_aes_256_cbc": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm AES-256-CBC.",
			Default:     true,
		},
		"allow_algorithm_aes_128_gcm": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm AES-128-GCM.",
			Default:     true,
		},
		"allow_algorithm_aes_256_gcm": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow algorithm AES-256-GCM.",
			Default:     true,
		},
		"allow_authentication_sha1": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow authentication SHA1.",
			Default:     true,
		},
		"allow_authentication_sha256": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow authentication SHA256.",
			Default:     true,
		},
		"allow_authentication_sha384": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow authentication SHA384.",
			Default:     true,
		},
	}

	if !isResource {
		computed(ans, "", []string{"template", "template_stack", "vsys", "name"})
	}

	return ans
}

func loadSslTlsServiceProfile(d *schema.ResourceData) ssltls.Entry {
	return ssltls.Entry{
		Name:                      d.Get("name").(string),
		Certificate:               d.Get("certificate").(string),
		MinVersion:                d.Get("min_version").(string),
		MaxVersion:                d.Get("max_version").(string),
		AllowAlgorithmRsa:         d.Get("allow_algorithm_rsa").(bool),
		AllowAlgorithmDhe:         d.Get("allow_algorithm_dhe").(bool),
		AllowAlgorithmEcdhe:       d.Get("allow_algorithm_ecdhe").(bool),
		AllowAlgorithm3des:        d.Get("allow_algorithm_3des").(bool),
		AllowAlgorithmRc4:         d.Get("allow_algorithm_rc4").(bool),
		AllowAlgorithmAes128Cbc:   d.Get("allow_algorithm_aes_128_cbc").(bool),
		AllowAlgorithmAes256Cbc:   d.Get("allow_algorithm_aes_256_cbc").(bool),
		AllowAlgorithmAes128Gcm:   d.Get("allow_algorithm_aes_128_gcm").(bool),
		AllowAlgorithmAes256Gcm:   d.Get("allow_algorithm_aes_256_gcm").(bool),
		AllowAuthenticationSha1:   d.Get("allow_authentication_sha1").(bool),
		AllowAuthenticationSha256: d.Get("allow_authentication_sha256").(bool),
		AllowAuthenticationSha384: d.Get("allow_authentication_sha384").(bool),
	}
}

func saveSslTlsServiceProfile(d *schema.ResourceData, o ssltls.Entry) {
	d.Set("name", o.Name)
	d.Set("certificate", o.Certificate)
	d.Set("min_version", o.MinVersion)
	d.Set("max_version", o.MaxVersion)
	d.Set("allow_algorithm_rsa", o.AllowAlgorithmRsa)
	d.Set("allow_algorithm_dhe", o.AllowAlgorithmDhe)
	d.Set("allow_algorithm_ecdhe", o.AllowAlgorithmEcdhe)
	d.Set("allow_algorithm_3des", o.AllowAlgorithm3des)
	d.Set("allow_algorithm_rc4", o.AllowAlgorithmRc4)
	d.Set("allow_algorithm_aes_128_cbc", o.AllowAlgorithmAes128Cbc)
	d.Set("allow_algorithm_aes_256_cbc", o.AllowAlgorithmAes256Cbc)
	d.Set("allow_algorithm_aes_128_gcm", o.AllowAlgorithmAes128Gcm)
	d.Set("allow_algorithm_aes_256_gcm", o.AllowAlgorithmAes256Gcm)
	d.Set("allow_authentication_sha1", o.AllowAuthenticationSha1)
	d.Set("allow_authentication_sha256", o.AllowAuthenticationSha256)
	d.Set("allow_authentication_sha384", o.AllowAuthenticationSha384)
}

// Id functions.
func parseSslTlsServiceProfileId(v string) (string, string, string, string, error) {
	t := strings.Split(v, IdSeparator)
	if len(t) != 4 {
		return "", "", "", "", fmt.Errorf("Expected len(4) ID, got %d", len(t))
	}

	return t[0], t[1], t[2], t[3], nil
}

func buildSslTlsServiceProfileId(a, b, c, d string) string {
	return strings.Join([]string{a, b, c, d}, IdSeparator)
}
