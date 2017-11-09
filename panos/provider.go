package panos

import (
    "fmt"
    "github.com/PaloAltoNetworks/pango"

    "github.com/hashicorp/terraform/terraform"
    "github.com/hashicorp/terraform/helper/schema"
)


// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
    return &schema.Provider{
        Schema: map[string] *schema.Schema{
            "hostname": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("PANOS_HOSTNAME", nil),
                Description: "Hostname/IP address of the Palo Alto Networks firewall to connect to",
            },
            "username": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("PANOS_USERNAME", nil),
                Description: "The username (not used if the ApiKey is set)",
            },
            "password": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("PANOS_PASSWORD", nil),
                Description: "The password (not used if the ApiKey is set)",
            },
            "api_key": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                DefaultFunc: schema.EnvDefaultFunc("PANOS_API_KEY", nil),
                Description: "The api key of the firewall",
            },
            "protocol": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Default: "https",
                Description: "The protocol (https or http)",
            },
            "port": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Description: "If the port is non-standard for the protocol, the port number to use",
            },
            "timeout": &schema.Schema{
                Type: schema.TypeInt,
                Optional: true,
                Description: "The timeout for all communications with the firewall",
            },
            "logging": &schema.Schema{
                Type: schema.TypeList,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
                Optional: true,
                Description: "Logging options for the API connection",
            },
        },

        ResourcesMap: map[string] *schema.Resource{
            "panos_address_object": resourceAddressObject(),
            "panos_ethernet_interface": resourceEthernetInterface(),
            "panos_general_settings": resourceGeneralSettings(),
            "panos_management_profile": resourceManagementProfile(),
            "panos_security_policy": resourceSecurityPolicy(),
        },

        ConfigureFunc: providerConfigure,
    }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
    var logging uint32

    lc := d.Get("logging")
    if lc != nil {
        ll := lc.([]interface{})
        for i := range ll {
            v := ll[i].(string)
            switch v {
            case "quiet":
                logging |= pango.LogQuiet
            case "action":
                logging |= pango.LogAction
            case "query":
                logging |= pango.LogQuery
            case "op":
                logging |= pango.LogOp
            case "uid":
                logging |= pango.LogUid
            case "xpath":
                logging |= pango.LogXpath
            case "send":
                logging |= pango.LogSend
            case "receive":
                logging |= pango.LogReceive
            default:
                return nil, fmt.Errorf("Unknown logging artifact requested: %s", v)
            }
        }
    }

    fw := &pango.Firewall{Client: pango.Client{
        Hostname: d.Get("hostname").(string),
        Username: d.Get("username").(string),
        Password: d.Get("password").(string),
        ApiKey: d.Get("api_key").(string),
        Protocol: d.Get("protocol").(string),
        Port: uint(d.Get("port").(int)),
        Timeout: d.Get("timeout").(int),
        Logging: logging,
    }}
    if err := fw.Initialize(); err != nil {
        return nil, err
    }

    return fw, nil
}
