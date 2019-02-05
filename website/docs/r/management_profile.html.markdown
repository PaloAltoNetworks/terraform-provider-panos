---
layout: "panos"
page_title: "panos: panos_management_profile"
sidebar_current: "docs-panos-resource-management-profile"
description: |-
  Manages interface management profiles.
---

# panos_management_profile

This resource allows you to add/update/delete interface management profiles.


## Import Name

```
<name>
```


## Example Usage

```hcl
resource "panos_management_profile" "example" {
    name = "allow ping"
    ping = true
    permitted_ips = ["10.1.1.0/24", "192.168.80.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The management profile's name.
* `ping` - (Optional) Allow ping.
* `telnet` - (Optional) Allow telnet.
* `ssh` - (Optional) Allow SSH.
* `http` - (Optional) Allow HTTP.
* `http_ocsp` - (Optional) Allow HTTP OCSP.
* `https` - (Optional) Allow HTTPS.
* `snmp` - (Optional) Allow SNMP.
* `response_pages` - (Optional) Allow response pages.
* `userid_service` - (Optional) Allow User ID service.
* `userid_syslog_listener_ssl` - (Optional) Allow User ID syslog listener
  for SSL.
* `userid_syslog_listener_udp` - (Optional) Allow User ID syslog listener
  for UDP.
* `permitted_ips` - (Optional) The list of permitted IP addresses or address
  ranges for this management profile.
