package ssldecrypt

import (
	"fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Device.SslDecrypt namespace.
type Panorama struct {
	ns *namespace.Standard
}

// SetTrustedRootCa adds a certificate as a trusted root CA.
func (c *Panorama) SetTrustedRootCa(tmpl, ts, vsys, name string) error {
	path, err := c.xpath(tmpl, ts, vsys)
	if err != nil {
		return err
	}
	path = append(path, "trusted-root-CA")

	c.ns.Client.LogAction("(set) %s trusted root ca: %s", singular, name)

	_, err = c.ns.Client.Set(path, util.Member{Value: name}, nil, nil)
	return err
}

// DeleteTrustedRootCa removes a certificate as a trusted root CA.
func (c *Panorama) DeleteTrustedRootCa(tmpl, ts, vsys, name string) error {
	c.ns.Client.LogAction("(delete) %s trusted root ca: %s", singular, name)

	path, err := c.xpath(tmpl, ts, vsys)
	if err != nil {
		return err
	}
	path = append(path, "trusted-root-CA", util.AsMemberXpath([]string{name}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// Get performs GET to retrieve configuration for the given object.
func (c *Panorama) Get(tmpl, ts, vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(tmpl, ts, vsys), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Panorama) Show(tmpl, ts, vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(tmpl, ts, vsys), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Panorama) Set(tmpl, ts, vsys string, e Config) error {
	return c.ns.Set(c.pather(tmpl, ts, vsys), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Panorama) Edit(tmpl, ts, vsys string, e Config) error {
	return c.ns.Edit(c.pather(tmpl, ts, vsys), e)
}

// Delete performs DELETE to remove the config.
func (c *Panorama) Delete(tmpl, ts, vsys string) error {
	return c.ns.Delete(c.pather(tmpl, ts, vsys), nil, nil)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Panorama) FromPanosConfig(tmpl, ts, vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(tmpl, ts, vsys), "", ans)
	return first(ans, err)
}

func (c *Panorama) pather(tmpl, ts, vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(tmpl, ts, vsys)
	}
}

func (c *Panorama) xpath(tmpl, ts, vsys string) ([]string, error) {
	ans := make([]string, 0, 11)

	if tmpl == "" && ts == "" {
		if vsys != "" {
			return nil, fmt.Errorf("vsys should be empty for local %s config", singular)
		}

		ans = append(ans, "config", "shared")
	} else {
		ans = append(ans, util.TemplateXpathPrefix(tmpl, ts)...)
		ans = append(ans, util.VsysXpathPrefix(vsys)...)
	}

	ans = append(ans, "ssl-decrypt")

	return ans, nil
}

func (c *Panorama) container() normalizer {
	return container(c.ns.Client.Versioning())
}
