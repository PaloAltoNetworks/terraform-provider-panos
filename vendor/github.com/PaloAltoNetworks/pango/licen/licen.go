// Package licen is the client.Licensing namespace.
package licen

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/util"
)

// Licen is the client.Licensing namespace.
type Licen struct {
	con util.XapiClient
}

// Initialize is invoked on client.Initialize().
func (c *Licen) Initialize(i util.XapiClient) {
	c.con = i
}

// Current returns the licenses currently installed.
func (c *Licen) Current() ([]util.License, error) {
	type lic_req struct {
		XMLName xml.Name `xml:"request"`
		Cmd     string   `xml:"license>info"`
	}

	c.con.LogOp("(op) request license info")
	return c.returnLicenseList(lic_req{})
}

// Fetch fetches licenses from the license server.
func (c *Licen) Fetch() ([]util.License, error) {
	type fetch struct {
		XMLName xml.Name `xml:"request"`
		Cmd     string   `xml:"license>fetch"`
	}

	c.con.LogOp("(op) request license fetch")
	return c.returnLicenseList(fetch{})
}

// Activate updates a license using the given auth code.
func (c *Licen) Activate(auth string) error {
	type auth_req struct {
		XMLName xml.Name `xml:"request"`
		Code    string   `xml:"license>fetch>auth-code"`
	}

	c.con.LogOp("(op) request license fetch auth-code \"********\"")
	body, err := c.con.Op(auth_req{Code: auth}, "", nil, nil)
	if err != nil {
		if string(body) == "VM Device License installed. Restarting pan services." {
			return nil
		}
	}

	return err
}

// Deactivate removes all licenses from a firewall.
//
// In order for this function to work, the following must be true:
//
//   * PAN-OS 7.1 or later
//   * PAN-OS has connectivity to the Palo Alto Networks support server
//   * Check server identity is enabled
//   * The licensing API key has been installed
func (c *Licen) Deactivate() error {
	type del_req struct {
		XMLName xml.Name `xml:"request"`
		Mode    string   `xml:"license>deactivate>VM-Capacity>mode"`
	}

	c.con.LogOp("(op) request license deactivate VM-Capacity mode auto")
	_, err := c.con.Op(del_req{Mode: "auto"}, "", nil, nil)
	return err
}

// GetApiKey returns the licensing API key.
func (c *Licen) GetApiKey() (string, error) {
	type get_req struct {
		XMLName xml.Name `xml:"request"`
		Cmd     string   `xml:"license>api-key>show"`
	}

	type get_resp struct {
		XMLName xml.Name `xml:"response"`
		Data    string   `xml:"result"`
	}

	ans := get_resp{}
	c.con.LogOp("(op) request license api-key show")
	if _, err := c.con.Op(get_req{}, "", nil, &ans); err != nil {
		if err.Error() == "API Key is not set" {
			return "", nil
		}
		return "", err
	}

	prefix := "API key: "
	if strings.HasPrefix(ans.Data, prefix) {
		return ans.Data[len(prefix):], nil
	} else {
		return ans.Data, nil
	}
}

// SetApiKey sets the licensing API key.
func (c *Licen) SetApiKey(k string) error {
	type set_req struct {
		XMLName xml.Name `xml:"request"`
		Key     string   `xml:"license>api-key>set>key"`
	}

	c.con.LogOp("(op) request license api-key set key \"********\"")
	_, err := c.con.Op(set_req{Key: k}, "", nil, nil)
	if err != nil && err.Error() == "API key is same as old" {
		return nil
	}

	return err
}

// DeleteApiKey deletes the licensing API key.
func (c *Licen) DeleteApiKey() error {
	type del_req struct {
		XMLName xml.Name `xml:"request"`
		Cmd     string   `xml:"license>api-key>delete"`
	}

	c.con.LogOp("(op) request license api-key delete")
	_, err := c.con.Op(del_req{}, "", nil, nil)

	if err != nil && err.Error() == "No API Key to be deleted" {
		return nil
	}

	return err
}

/** Structs / functions for this namespace. **/

func (c *Licen) returnLicenseList(req interface{}) ([]util.License, error) {
	type lic_resp struct {
		XMLName xml.Name       `xml:"response"`
		Data    []util.License `xml:"result>licenses>entry"`
	}

	ans := lic_resp{}

	if _, err := c.con.Op(req, "", nil, &ans); err != nil {
		return nil, fmt.Errorf("Failed to get licenses: %s", err)
	}

	return ans.Data, nil
}
