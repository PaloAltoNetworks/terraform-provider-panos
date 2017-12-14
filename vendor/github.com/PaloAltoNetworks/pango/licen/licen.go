// Package licen is the client.Licensing namespace.
package licen

import (
    "fmt"
    "encoding/xml"

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
        Cmd string `xml:"license>info"`
    }

    c.con.LogOp("(op) request license info")
    return c.returnLicenseList(lic_req{})
}

// Fetch fetches licenses from the license server.
func (c *Licen) Fetch() ([]util.License, error) {
    type fetch struct {
        XMLName xml.Name `xml:"request"`
        Cmd string `xml:"license>fetch"`
    }

    c.con.LogOp("(op) request license fetch")
    return c.returnLicenseList(fetch{})
}

// Activate updates a license using the given auth code.
func (c *Licen) Activate(auth string) error {
    type auth_req struct {
        XMLName xml.Name `xml:"request"`
        Code string `xml:"license>fetch>auth-code"`
    }

    c.con.LogOp("(op) request license fetch auth-code \"********\"")
    _, err := c.con.Op(auth_req{Code: auth}, "", nil, nil)
    return err
}

/** Structs / functions for this namespace. **/

func (c *Licen) returnLicenseList(req interface{}) ([]util.License, error) {
    type lic_resp struct {
        XMLName xml.Name `xml:"response"`
        Data []util.License `xml:"result>licenses>entry"`
    }

    ans := lic_resp{}

    if _, err := c.con.Op(req, "", nil, &ans); err != nil {
        return nil, fmt.Errorf("Failed to get licenses: %s", err)
    }

    return ans.Data, nil
}
