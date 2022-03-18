---
page_title: "panos: panos_dhcp_relay"
subcategory: "Network"
---

# panos_bgp

This resource allows you to add/update/delete a dhcp relay.


## PAN-OS

NGFW


## Import Name

```
<dhcp-relay>
```


## Example Usage

```hcl
resource "panos_dhcp_relay" "relay" {
  name = "ethernet1/3"
  ipv4_enabled = true
  ipv4_servers = [
    "203.0.113.1",
    "203.0.113.254",
  ]
  ipv6_enabled = true
  ipv6_servers = [
    {
      server = "2001:db8::1",
      interface = "ethernet1/3",
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, string) The name of the interface to add the dhcp relay to.
* `ipv4_enabled` - (Optional, bool) Enable IPv4 DHCP relay.
* `ipv4_servers` - (Optional, list(string)) List of IPv4 DHCP servers.
* `ipv6_enabled` - (Optional, bool) Enable IPv6 DHCP relay.
* `ipv6_servers` - (Optional, list(object)) List of IPv6 DHCP servers.
* `vsys` - (Optional) The vsys location (default: `vsys1`).
