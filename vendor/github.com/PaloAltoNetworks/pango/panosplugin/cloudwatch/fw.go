package cloudwatch

import (
	"github.com/PaloAltoNetworks/pango/namespace"
	"github.com/PaloAltoNetworks/pango/util"
)

// Firewall is the client.Plugin.AwsCloudWatch namespace.
type Firewall struct {
	ns *namespace.Plugin
}

// Get performs GET to retrieve information for the given object.
func (c *Firewall) Get() (Config, error) {
	ans, err := c.container()
	if err != nil {
		return Config{}, err
	}
	err = c.ns.Object(util.Get, c.pather(), "", ans)
	return first(ans, err)
}

// Show performs SHOW to retrieve information for the given object.
func (c *Firewall) Show() (Config, error) {
	ans, err := c.container()
	if err != nil {
		return Config{}, err
	}
	err = c.ns.Object(util.Show, c.pather(), "", ans)
	return first(ans, err)
}

// Set performs SET to configure the specified objects.
func (c *Firewall) Set(e Config) error {
	return c.ns.Set(c.pather(), specifier(e))
}

// Edit performs EDIT to configure the specified object.
func (c *Firewall) Edit(e Config) error {
	return c.ns.Edit(c.pather(), e)
}

// Delete performs DELETE to remove the config.
func (c *Firewall) Delete() error {
	return c.ns.Delete(c.pather(), nil, nil)
}

func (c *Firewall) pather() namespace.Pather {
	return func(v []string) ([]string, error) {
		return c.xpath(v)
	}
}

func (c *Firewall) xpath(vals []string) ([]string, error) {
	return []string{
		"config",
		"devices",
		util.AsEntryXpath([]string{"localhost.localdomain"}),
		"deviceconfig",
		"plugins",
		"vm_series",
		"aws-cloudwatch",
	}, nil
}

func (c *Firewall) container() (normalizer, error) {
	return container(c.ns.Client.Plugins())
}
