---
layout: "panos"
page_title: "Provider: PANOS"
sidebar_current: "docs-panos-index"
description: |-
  PANOS is used to interact with Palo Alto Networks NGFW.
---

# PANOS Provider

[PANOS](https://www.paloaltonetworks.com/) is used to interact with Palo Alto
Networks' NGFW.  The provider allows you to manage various aspects of the
firewall's config, such as data interfaces or security policies.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the PANOS provider
provider "panos" {
    hostname = "127.0.0.1"
    username = "admin"
    password = "secret"
}

# Add a security policy to the firewall
resource "panos_security_policy" "myPolicy" {
    # ...
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Optional) This is the hostname / IP address of the firewall.  It
  must be provided, but can also be defined via the `PANOS_HOSTNAME`
  environment variable.
* `username` - (Optional) The username to authenticate to the firewall as.  It
  must be provided, but can also be defined via the `PANOS_USERNAME`
  environment variable.
* `password` - (Optional) The password for the given username. It must be
  provided, but can also be defined via the `PANOS_PASSWORD` environment
  variable.
* `api_key` - (Optional) The API key for the firewall.  If this is given, then
  the `username` and `password` settings are ignored.  This can also be defined
  via the `PANOS_API_KEY` environment variable.
* `protocol` - (Optional) The communication protocol.  This can be set to
  either `https` or `http`.  If left unspecified, this defaults to `https`.  
* `port` - (Optional) If the port number is non-standard for the desired
  protocol, then the port number to use.
* `timeout` - (Optional) The timeout for all communications with the
  firewall.  If left unspecified, this will be set to 10 seconds.
* `logging` - (Optional) Logging options for the API connection.
