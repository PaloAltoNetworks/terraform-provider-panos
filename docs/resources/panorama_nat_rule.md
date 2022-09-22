---
page_title: "panos: panos_panorama_nat_rule"
subcategory: "Policies"
---

# panos_panorama_nat_rule

This resource allows you to add/update/delete Panorama NAT rules.

~> **Note:** This resource has been deprecated.  Please use
`panos_panorama_nat_rule_group` instead.

~> **Note:** `panos_panorama_nat_policy` is known as `panos_panorama_nat_rule`.

The prefix `sat` stands for "Source Address Translation" while the prefix "dat"
stands for "Destination Address Translation".  The order of the params in
this resource and their naming matches how the params are presented in
the GUI.  Thus, having a GUI window open while creating your resource
definition will simplify the process.

Note that while many of the params for this resource are optional in an
absolute sense, depending on what type of NAT you wish to configure, certain
params may become necessary to correctly configure the NAT rule.


## PAN-OS

Panorama.


## Aliases

* `panos_panorama_nat_policy`


## Example Usage

```hcl
resource "panos_panorama_nat_rule" "example" {
    name = "my nat rule"
    source_zones = panos_panorama_zone.z1.name
    destination_zone = panos_panorama_zone.z2.name
    to_interface = panos_panorama_ethernet_interface.e1.name
    source_addresses = ["any"]
    destination_addresses = ["any"]
    sat_type = "none"
    dat_type = "static"
    dat_address = "my dat address object"
    target {
        serial = "123456"
        vsys_list = ["vsys1", "vsys2"]
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The NAT rule's name.
* `device_group` - (Optional) The device group to put the NAT rule into
  (default: `shared`).
* `rulebase` - (Optional) The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.
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
  This can be `interface-address` or `translated-address`.
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
* `dat_type` - (Optional) Destination address translation type.  This should
  be either `static` or `dynamic`.  The `dynamic` option is only available on
  PAN-OS 8.1+.
* `dat_address` - (Optional) Destination address translation's address.  Requires
  `dat_type` be set to "static" or "dynamic".
* `dat_port` - (Optional) Destination address translation's port number.  Requires
  `dat_type` be set to "static" or "dynamic".
* `dat_dynamic_distribution` - (Optional, PAN-OS 8.1+) Distribution algorithm
  for destination address pool.  The PAN-OS 8.1 GUI doesn't seem to set this
  anywhere, but this is added here for completeness' sake.  Requires `dat_type`
  of "dynamic".
* `disabled` - (Optional) Set to `true` to disable this rule.
* `tags` - (Optional) List of administrative tags.
* `target` - (Optional) A target definition (see below).  If there are no
  target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (Optional, bool) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.

The following arguments are valid for each `target` section:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.
