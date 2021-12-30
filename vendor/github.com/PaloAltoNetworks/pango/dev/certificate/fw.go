package certificate

import (
	"net/url"
	"time"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Device.Certificate namespace.
type Firewall struct {
	ns *namespace.Standard
}

// GetList performs GET to retrieve a list of all objects.
func (c *Firewall) GetList(vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Get, c.pather(vsys), ans)
}

// ShowList performs SHOW to retrieve a list of all objects.
func (c *Firewall) ShowList(vsys string) ([]string, error) {
	ans := c.container()
	return c.ns.Listing(util.Show, c.pather(vsys), ans)
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vsys), name, ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vsys), name, ans)
	return first(ans, err)
}

// GetAll performs GET to retrieve all objects configured.
func (c *Firewall) GetAll(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Get, c.pather(vsys), ans)
	return all(ans, err)
}

// ShowAll performs SHOW to retrieve information for all objects.
func (c *Firewall) ShowAll(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.Objects(util.Show, c.pather(vsys), ans)
	return all(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vsys string, e ...Entry) error {
	return c.ns.Set(c.pather(vsys), specifier(e...))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vsys string, e Entry) error {
	return c.ns.Edit(c.pather(vsys), e)
}

// Delete performs DELETE to remove the specified objects.
//
// Objects can be either a string or an Entry object.
func (c *Firewall) Delete(vsys string, e ...interface{}) error {
	names, nErr := toNames(e)
	return c.ns.Delete(c.pather(vsys), names, nErr)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Firewall) FromPanosConfig(vsys, name string) (Entry, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(vsys), name, ans)
	return first(ans, err)
}

// AllFromPanosConfig retrieves all objects stored in the retrieved config.
func (c *Firewall) AllFromPanosConfig(vsys string) ([]Entry, error) {
	ans := c.container()
	err := c.ns.AllFromPanosConfig(c.pather(vsys), ans)
	return all(ans, err)
}

// ImportPem imports a PEM certificate.
func (c *Firewall) ImportPem(vsys string, timeout time.Duration, cert Pem) error {
	var err error

	c.ns.Client.LogImport("(import) pem %s: %s", singular, cert.Name)

	ex := url.Values{}
	ex.Set("certificate-name", cert.Name)
	ex.Set("format", "pem")
	if vsys != "" && vsys != "shared" {
		ex.Set("vsys", vsys)
	}

	_, err = c.ns.Client.Import("certificate", cert.Certificate, cert.CertificateFilename, "file", timeout, ex, nil)

	if err != nil || cert.PrivateKey == "" {
		return err
	}

	ex.Set("passphrase", cert.Passphrase)

	_, err = c.ns.Client.Import("private-key", cert.PrivateKey, cert.PrivateKeyFilename, "file", timeout, ex, nil)

	return err
}

// ImportPkcs12 imports a PKCS12 certificate.
func (c *Firewall) ImportPkcs12(vsys string, timeout time.Duration, cert Pkcs12) error {
	var err error

	c.ns.Client.LogImport("(import) pkcs12 %s: %s", singular, cert.Name)

	ex := url.Values{}
	ex.Set("certificate-name", cert.Name)
	ex.Set("format", "pkcs12")
	ex.Set("passphrase", cert.Passphrase)
	if vsys != "" && vsys != "shared" {
		ex.Set("vsys", vsys)
	}

	_, err = c.ns.Client.Import("certificate", cert.Certificate, cert.CertificateFilename, "file", timeout, ex, nil)

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
func (c *Firewall) Export(format, vsys, name, passphrase string, includeKey bool, timeout time.Duration) (string, []byte, error) {
	c.ns.Client.LogExport("(export) %s %s: %s", format, singular, name)

	ex := url.Values{}
	ex.Set("certificate-name", name)
	ex.Set("format", format)
	ex.Set("include-key", util.YesNo(includeKey))
	if passphrase != "" {
		ex.Set("passphrase", passphrase)
	}
	if vsys != "" && vsys != "shared" {
		ex.Set("vsys", vsys)
	}

	return c.ns.Client.Export("certificate", timeout, ex, nil)
}

func (c *Firewall) pather(vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vsys, v)
	}
}

func (c *Firewall) xpath(vsys string, vals []string) ([]string, error) {
	if vsys == "" {
		vsys = "shared"
	}

	ans := make([]string, 0, 8)
	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans,
		"certificate",
		util.AsEntryXpath(vals),
	)

	return ans, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
