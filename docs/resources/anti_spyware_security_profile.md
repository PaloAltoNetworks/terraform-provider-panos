---
page_title: "panos: panos_anti_spyware_security_profile"
subcategory: "Objects"
---

# panos_anti_spyware_security_profile

Manages anti-spyware security profiles.

## Import Name

NGFW:

```
<vsys>:<name>
```

Panorama:

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_anti_spyware_security_profile" "example" {
    name = "example"
    description = "my description"
    sinkhole_ipv4_address = "pan-sinkhole-default-ip"
    sinkhole_ipv6_address = "::1"
    botnet_list {
        name = "default-paloalto-dns"
        action = "sinkhole"
        packet_capture = "disable"
    }
    botnet_list {
        name = "default-paloalto-cloud"
        action = "allow"
        packet_capture = "disable"
    }
    rule {
        name = "foo"
        threat_name = "any"
        category = "adware"
        action = "alert"
        packet_capture = "disable"
        severities = ["any"]
    }
    exception {
        name = data.panos_predefined_threat.dot_net.threats.0.name
        action = "allow"
    }
}

data "panos_predefined_threat" "dot_net" {
    threat_type = "phone-home"
    threat_name = "Generic .NET Framework C# Webshell Upload Detection"
}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `packet_capture` - (PAN-OS 8.x only) Packet capture setting.  Valid values
  are `disable`, `single-packet`, or `extended-capture`.
* `sinkhole_ipv4_address` - IPv4 sinkhole address.
* `sinkhole_ipv6_address` - IPv6 sinkhole address.
* `threat_exceptions` - (list) List of threat exceptions.
* `bonet_list` - (repeatable) Botnet spec, as defined below.
* `dns_category` - (repeatable, PAN-OS 10.0+) DNS category spec, as defined below.
* `white_list` - (repeatable, PAN-OS 10.0+) White list spec, as defined below.
* `rule` - (repeatable) Rule list spec, as defined below.
* `exception` - (repeatable) Except list spec, as defined below.

`botnet_list` supports the following arguments:

* `name` - (Required) Name.
* `action` - Action to take.  Valid values are `alert`, `default`, `allow`,
  `block`, or `sinkhole`.
* `packet_capture` - (PAN-OS 9.0+) Packet capture setting.  Valid values
  are `disable`, `single-packet`, or `extended-capture`.


`dns_category` supports the following arguments:

* `name` - (Required) Name.
* `action` - Action to take.  Valid values are `alert`, `default`, `allow`,
  `block`, or `sinkhole`.
* `log_level` - Logging level.  Valid values are `default`, `none`, `low`,
  `informational`, `medium`, `high`, or `critical`.
* `packet_capture` - Packet capture setting.  Valid values
  are `disable`, `single-packet`, or `extended-capture`.

`white_list` supports the following arguments:

* `name` - (Required) Name.
* `description` - Description

`rule` supports the following arguments:

* `name` - (Required) Name.
* `threat_name` - Threat name.
* `category` - (Required) Category.
* `severities` - (list) Severities.
* `packet_capture` - Packet capture setting.  Valid values
  are `disable`, `single-packet`, or `extended-capture`.
* `action` - Action.  Valid values are `default`, `allow`, `alert`, `drop`,
  `reset-client`, `reset-server`, `reset-both`, or `block-ip`.
* `block_ip_track_by` - (action=`block-ip`) The track by setting.
* `block_ip_duration` - (action=`block-ip`, int) The duration.

`exception` supports the following arguments:

* `name` - (Required) Threat name.  You can use the `panos_predefined_threat` data
  source to discover the various phone home names available to use.
* `packet_capture` - Packet capture setting.  Valid values
  are `disable`, `single-packet`, or `extended-capture`.
* `action` - Action.  Valid values are `default`, `allow`, `alert`, `drop`,
  `reset-client`, `reset-server`, `reset-both`, or `block-ip`.
* `block_ip_track_by` - (action=`block-ip`) The track by setting.
* `block_ip_duration` - (action=`block-ip`, int) The duration.
* `exempt_ips` - (list) List of exempt IPs.
