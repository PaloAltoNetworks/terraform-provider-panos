package certificate

import (
	"net/url"
	"time"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Device.Certificate namespace.
type Panorama struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Panorama) GetList(shared bool, tmpl, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(shared, tmpl, vsys), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Panorama) ShowList(shared bool, tmpl, vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(shared, tmpl, vsys), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Panorama) Get(shared bool, tmpl, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(shared, tmpl, vsys), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Panorama) Show(shared bool, tmpl, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(shared, tmpl, vsys), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Panorama) GetAll(shared bool, tmpl, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(shared, tmpl, vsys), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Panorama) ShowAll(shared bool, tmpl, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(shared, tmpl, vsys), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(shared bool, tmpl, vsys string, e ...Entry) error {
	return c.ns.Set(c.pather(shared, tmpl, vsys), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(shared bool, tmpl, vsys string, e Entry) error {
	return c.ns.Edit(c.pather(shared, tmpl, vsys), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Panorama) Delete(shared bool, tmpl, vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(shared, tmpl, vsys), names, nErr)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Panorama) FromPanosConfig(shared bool, tmpl, vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(shared, tmpl, vsys), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Panorama) AllFromPanosConfig(shared bool, tmpl, vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(shared, tmpl, vsys), ans)
	return all(ans, err)
}

// ImportPem imports a PEM certificate.
func (c *Panorama) ImportPem(tmpl, vsys string, timeout time.Duration, cert Pem) error {
	var err error

	c.ns.Client.LogImport("(import) pem %s: %s", singular, cert.Name)

	ex := url.Values{}
	ex.Set("certificate-name", cert.Name)
	ex.Set("format", "pem")

	if tmpl != "" {
		ex.Set("target-tpl", tmpl)
		if vsys == "" {
			vsys = "shared"
		}
		ex.Set("target-tpl-vsys", vsys)
	}

	_, err = c.ns.Client.Import("certificate", cert.Certificate, cert.CertificateFilename, "file", timeout, ex, nil)

	if err != nil || cert.PrivateKey == "" {
		return err
	}

	ex.Set("passphrase", cert.Passphrase)

	_, err = c.ns.Client.Import("certificate", cert.PrivateKey, cert.PrivateKeyFilename, "file", timeout, ex, nil)

	return err
}

// ImportPkcs12 imports a PKCS12 certificate.
func (c *Panorama) ImportPkcs12(tmpl, vsys string, timeout time.Duration, cert Pkcs12) error {
	c.ns.Client.LogImport("(import) pkcs12 %s: %s", singular, cert.Name)

	ex := url.Values{}
	ex.Set("certificate-name", cert.Name)
	ex.Set("format", "pkcs12")
	ex.Set("passphrase", cert.Passphrase)

	if tmpl != "" {
		ex.Set("target-tpl", tmpl)
		if vsys == "" {
			vsys = "shared"
		}
		ex.Set("target-tpl-vsys", vsys)
	}

	_, err := c.ns.Client.Import("certificate", cert.Certificate, cert.CertificateFilename, "file", timeout, ex, nil)

	return err
}

// Export exports a certificate.
//
// The format param should be either "pem" or "pkcs12".
//
// The public key is always exported.
//
// Attempting to export a PKCS12 cert as a PEM cert will result in an error.
//
// Return values are the filename, file contents, and an error.
func (c *Panorama) Export(format, tmpl, vsys, name, passphrase string, includeKey bool, timeout time.Duration) (string, []byte, error) {
	c.ns.Client.LogExport("(export) %s %s: %s", format, singular, name)

	ex := url.Values{}
	ex.Set("certificate-name", name)
	ex.Set("format", format)
	ex.Set("include-key", util.YesNo(includeKey))
	if passphrase != "" {
		ex.Set("passphrase", passphrase)
	}
	if tmpl != "" {
		ex.Set("target-tpl", tmpl)
		if vsys != "" && vsys != "shared" {
			// TODO: This doesn't seem to work, but it's what the docs say.
			ex.Set("target-tpl-vsys", vsys)
		}
	}

	return c.ns.Client.Export("certificate", timeout, ex, nil)
}

func (c *Panorama) pather(shared bool, tmpl, vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(shared, tmpl, vsys, v)
	}
}

func (c *Panorama) xpath(shared bool, tmpl, vsys string, vals []string) ([]string, error) {
	var ans []string

	if tmpl != "" {
		ans = make([]string, 0, 12)

		ans = append(ans, util.TemplateXpathPrefix(tmpl, "")...)
		ans = append(ans, util.VsysXpathPrefix(vsys)...)
	} else {
		ans = make([]string, 0, 4)
		if shared {
			ans = append(ans,
				"config",
				"shared",
			)
		} else {
			ans = append(ans,
				"config",
				"panorama",
			)
		}
	}

	ans = append(ans,
		"certificate",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
