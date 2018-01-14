---
layout: "panos"
page_title: "PANOS: panos_nat_policy"
sidebar_current: "docs-panos-resource-nat-policy"
description: |-
  Manages NAT policies.
---

# panos_nat_policy

This resource allows you to add/update/delete NAT policies.

The prefix `sat` stands for "Source Address Translation" while the prefix "dat"
stands for "Destination Address Translation".  The order of the params in
this resource and their naming matches how the params are presented in
the GUI.  Thus, having a GUI window open while creating your resource
definition will simplify the process.

Note that while many of the params for this resource are optional in an
absolute sense, depending on what type of NAT you wish to configure, certain
params may become necessary for PANOS to correctly configure the NAT policy.

## Example Usage

```hcl
resource "panos_nat_policy" "example" {
    name = "my nat policy"
    source_zones = ["zone1"]
    destination_zone = "zone2"
    to_interface = "ethernet1/3"
    source_addresses = ["any"]
    destination_addresses = ["any"]
    sat_type = "none"
    dat_address = "my dat address object"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The NAT policy's name.
* `vsys` - (Optional) The vsys to put the NAT policy into (default:
  `vsys1`).
* `rulebase` - (Optional) The rulebase.  For firewalls, there is only the
  `rulebase` value (default), but on Panorama, there is also `pre-rulebase`
  and `post-rulebase`.
* `description` - (Optional) The description.
* `type` - (Optional). NAT type.  This can be `ipv4` (default), `nat64`, or
  `nptv6`.
* `source_zones` - (Required) The list of source zone(s).
* `destination_zone` - (Required) The destination zone.
* `to_interface` - (Optional) Egress interface from route lookup (default:
  `any`).
* `service` - (Optional) Service (default: `any`).
* `source_addresses` - (Required) List of source address(es).
* `destination_addresses` - (Required) List of destination address(es).
* `sat_type` - (Optional) Type of source address translation.  This can be
  `none` (default), `dynamic-ip-and-port`, `dynamic-ip`, or `static-ip`.
* `sat_address_type` - (Optional) Source address translation address type.
* `sat_translated_addresses` - (Optional) Source address translation list of
  translated addresses.
* `sat_interface` - (Optional) Source address translation interface.
* `sat_ip_address` - (Optional) Source address translation IP address.
* `sat_fallback_type` - (Optional) Source address translation fallback type.
  This can be `none`, `interface-address`, or `translated-address`.
* `sat_fallback_translated_addresses` - (Optional) Source address translation
  list of fallback translated addresses.
* `sat_fallback_interface` - (Optional) Source address translation fallback
  interface.
* `sat_fallback_ip_type` - (Optional) Source address translation fallback
  IP type.  This can be `ip` or `floating`.
* `sat_fallback_ip_address` - (Optional) The source address translation
  fallback IP address.
* `sat_static_translated_address` - (Optional) The statically translated source
  address.
* `sat_static_bi_directional` - (Optional) Set to `true` to enable
  bi-directional source address translation.
* `dat_address` - (Optional) Destination address translation's address.
* `dat_port` - (Optional) Destination address translation's port number.
* `disabled` - (Optional) Set to `true` to disable this rule.
* `tags` - (Optional) List of administrative tags.
