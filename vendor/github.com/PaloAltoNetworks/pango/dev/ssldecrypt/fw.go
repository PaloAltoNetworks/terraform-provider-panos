package ssldecrypt

import (
	_ "fmt"

	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Device.SslDecrypt namespace.
type Firewall struct {
	ns *namespace.Standard
}

// SetTrustedRootCa adds a certificate as a trusted root CA.
func (c *Firewall) SetTrustedRootCa(vsys, name string) error {
	path, err := c.xpath(vsys)
	if err != nil {
		return err
	}
	path = append(path, "trusted-root-CA")

	c.ns.Client.LogAction("(set) %s trusted root ca: %s", singular, name)

	_, err = c.ns.Client.Set(path, util.Member{Value: name}, nil, nil)
	return err
}

// DeleteTrustedRootCa removes a certificate as a trusted root CA.
func (c *Firewall) DeleteTrustedRootCa(vsys, name string) error {
	c.ns.Client.LogAction("(delete) %s trusted root ca: %s", singular, name)

	path, err := c.xpath(vsys)
	if err != nil {
		return err
	}
	path = append(path, "trusted-root-CA", util.AsMemberXpath([]string{name}))

	_, err = c.ns.Client.Delete(path, nil, nil)
	return err
}

// Get performs GET to retrieve configuration for the given object.
func (c *Firewall) Get(vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Get, c.pather(vsys), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve configuration for the given object.
func (c *Firewall) Show(vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.Object(util.Show, c.pather(vsys), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(vsys string, e Config) error {
	return c.ns.Set(c.pather(vsys), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(vsys string, e Config) error {
	return c.ns.Edit(c.pather(vsys), e)
}

// Delete performs DELETE to remove the config.
func (c *Firewall) Delete(vsys string) error {
	return c.ns.Delete(c.pather(vsys), nil, nil)
}

// FromPanosConfig retrieves the object stored in the retrieved config.
func (c *Firewall) FromPanosConfig(vsys string) (Config, error) {
	ans := c.container()
	err := c.ns.FromPanosConfig(c.pather(vsys), "", ans)
	return first(ans, err)
}

func (c *Firewall) pather(vsys string) namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(vsys)
	}
}

func (c *Firewall) xpath(vsys string) ([]string, error) {
	ans := make([]string, 0, 6)

	ans = append(ans, util.VsysXpathPrefix(vsys)...)
	ans = append(ans, "ssl-decrypt")

	return ans, nil
}

func (c *Firewall) container() normalizer {
	return container(c.ns.Client.Versioning())
}
