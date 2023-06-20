package panos

import (
	"log"
	"strings"

	"github.com/fpluchorg/pango"
	cert "github.com/fpluchorg/pango/dev/certificate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Resource.
func resourceCertificateImport() *schema.Resource {
	return &schema.Resource{
		Create: createUpdateCertificateImport,
		Read:   readCertificateImport,
		Update: createUpdateCertificateImport,
		Delete: deleteCertificateImport,

		Schema: certificateImportSchema(),
	}
}

func createUpdateCertificateImport(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o cert.Entry

	name := d.Get("name").(string)
	tmpl := d.Get("template").(string)
	vsys := d.Get("vsys").(string)

	id := buildCertificateImportId(tmpl, vsys, name)
	d.Set("name", name)
	d.Set("template", tmpl)
	d.Set("vsys", vsys)

	pemConfig := configFolder(d, "pem")
	pkcsConfig := configFolder(d, "pkcs12")

	switch {
	case pemConfig != nil:
		data := loadPemCert(name, pemConfig)
		d.Set("cert_format", "pem")
		d.Set("cert_passphrase", data.Passphrase)
		switch con := meta.(type) {
		case *pango.Firewall:
			err = con.Device.Certificate.ImportPem(vsys, 0, data)
		case *pango.Panorama:
			err = con.Device.Certificate.ImportPem(tmpl, vsys, 0, data)
		}
		if err == nil {
			d.SetId(id)
			info := map[string]interface{}{
				"certificate":          data.Certificate,
				"certificate_filename": data.CertificateFilename,
				"private_key":          data.PrivateKey,
				"private_key_filename": data.PrivateKeyFilename,
				"passphrase":           data.Passphrase,
			}
			if e2 := d.Set("pem", []interface{}{info}); e2 != nil {
				log.Printf("[WARN] Error setting 'pem' for %q: %s", d.Id(), e2)
			}
			d.Set("pkcs12", nil)
		}
	case pkcsConfig != nil:
		data := loadPkcs12Cert(name, pkcsConfig)
		d.Set("cert_format", "pkcs12")
		d.Set("cert_passphrase", data.Passphrase)
		switch con := meta.(type) {
		case *pango.Firewall:
			err = con.Device.Certificate.ImportPkcs12(vsys, 0, data)
		case *pango.Panorama:
			err = con.Device.Certificate.ImportPkcs12(tmpl, vsys, 0, data)
		}
		if err == nil {
			d.SetId(id)
			info := map[string]interface{}{
				"certificate":          data.Certificate,
				"certificate_filename": data.CertificateFilename,
				"passphrase":           data.Passphrase,
			}
			if e2 := d.Set("pkcs12", []interface{}{info}); e2 != nil {
				log.Printf("[WARN] Error setting 'pem' for %q: %s", d.Id(), e2)
			}
			d.Set("pem", nil)
		}
	}

	if err != nil {
		d.SetId("")
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.Certificate.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.Certificate.Get(false, tmpl, vsys, name)
	}

	if err != nil {
		return err
	}

	d.Set("cert_public_key", o.PublicKey)
	return readCertificateImport(d, meta)
}

func readCertificateImport(d *schema.ResourceData, meta interface{}) error {
	var err error
	var o cert.Entry

	tmpl, vsys, name := parseCertificateImportId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		o, err = con.Device.Certificate.Get(vsys, name)
	case *pango.Panorama:
		o, err = con.Device.Certificate.Get(false, tmpl, vsys, name)
	}

	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	switch con := meta.(type) {
	case *pango.Firewall:
		_, _, err = con.Device.Certificate.Export(d.Get("cert_format").(string), vsys, name, d.Get("cert_passphrase").(string), true, 0)
	case *pango.Panorama:
		_, _, err = con.Device.Certificate.Export(d.Get("cert_format").(string), tmpl, vsys, name, d.Get("cert_passphrase").(string), true, 0)
	}

	if err != nil {
		// Failing to export the cert means that the passphrase is wrong; blank the
		// config to make Terraform redeploy the cert.
		d.Set("pem", nil)
		d.Set("pkcs12", nil)
	} else if o.PublicKey != d.Get("cert_public_key").(string) {
		// It would be better to verify that the private key is unchanged, but
		// just check that the public key is the same for right now.
		d.Set("pem", nil)
		d.Set("pkcs12", nil)
	}

	saveCertificateImport(d, o)
	return nil
}

func deleteCertificateImport(d *schema.ResourceData, meta interface{}) error {
	var err error
	tmpl, vsys, name := parseCertificateImportId(d.Id())

	switch con := meta.(type) {
	case *pango.Firewall:
		err = con.Device.Certificate.Delete(vsys, name)
	case *pango.Panorama:
		err = con.Device.Certificate.Delete(false, tmpl, vsys, name)
	}

	if err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}

// Schema functions.
func certificateImportSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vsys": vsysSchema("shared"),
		"template": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Template to import into.",
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The certificate name.",
			ForceNew:    true,
		},
		"pem": {
			Type:          schema.TypeList,
			Optional:      true,
			Description:   "PEM certificate specification.",
			MaxItems:      1,
			ConflictsWith: []string{"pkcs12"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"certificate": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The contents of the certificate file.",
					},
					"certificate_filename": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The certificate filename.",
						Default:     "cert.pem",
					},
					"private_key": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The contents of the private key file.",
						Sensitive:   true,
					},
					"private_key_filename": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The private key filename.",
						Default:     "key.pem",
					},
					"passphrase": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The private key file passphrase.",
						Sensitive:   true,
					},
				},
			},
		},
		"pkcs12": {
			Type:          schema.TypeList,
			Optional:      true,
			Description:   "PKCS12 certificate specification.",
			MaxItems:      1,
			ConflictsWith: []string{"pem"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"certificate": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The contents of the certificate file.",
					},
					"certificate_filename": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The certificate filename.",
						Default:     "cert.pfx",
					},
					"passphrase": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The passphrase.",
						Sensitive:   true,
					},
				},
			},
		},

		// Attributes.
		"cert_format": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cert_public_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cert_passphrase": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"common_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"algorithm": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ca": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"not_valid_after": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"not_valid_before": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"expiry_epoch": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subject": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subject_hash": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"issuer": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"issuer_hash": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"csr": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_key": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"private_key_on_hsm": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"revoke_date_epoch": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func loadPemCert(name string, d map[string]interface{}) cert.Pem {
	return cert.Pem{
		Name:                name,
		Certificate:         d["certificate"].(string),
		CertificateFilename: d["certificate_filename"].(string),
		PrivateKey:          d["private_key"].(string),
		PrivateKeyFilename:  d["private_key_filename"].(string),
		Passphrase:          d["passphrase"].(string),
	}
}

func loadPkcs12Cert(name string, d map[string]interface{}) cert.Pkcs12 {
	return cert.Pkcs12{
		Name:                name,
		Certificate:         d["certificate"].(string),
		CertificateFilename: d["certificate_filename"].(string),
		Passphrase:          d["passphrase"].(string),
	}
}

func saveCertificateImport(d *schema.ResourceData, o cert.Entry) {
	d.Set("name", o.Name)
	d.Set("common_name", o.CommonName)
	d.Set("algorithm", o.Algorithm)
	d.Set("ca", o.Ca)
	d.Set("not_valid_after", o.NotValidAfter)
	d.Set("not_valid_before", o.NotValidBefore)
	d.Set("expiry_epoch", o.ExpiryEpoch)
	d.Set("subject", o.Subject)
	d.Set("subject_hash", o.SubjectHash)
	d.Set("issuer", o.Issuer)
	d.Set("issuer_hash", o.IssuerHash)
	d.Set("csr", o.Csr)
	d.Set("public_key", o.PublicKey)
	d.Set("private_key", o.PrivateKey)
	d.Set("private_key_on_hsm", o.PrivateKeyOnHsm)
	d.Set("status", o.Status)
	d.Set("revoke_date_epoch", o.RevokeDateEpoch)
}

// Id functions.
func buildCertificateImportId(a, b, c string) string {
	return strings.Join([]string{a, b, c}, IdSeparator)
}

func parseCertificateImportId(v string) (string, string, string) {
	t := strings.Split(v, IdSeparator)
	return t[0], t[1], t[2]
}
