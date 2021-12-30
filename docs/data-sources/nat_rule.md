---
page_title: "panos: panos_nat_rule_group"
subcategory: "Policies"
---

# panos_nat_rule_group

Retrieve information on the specified NAT rule.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_nat_rule" "example" {
    name = "my nat rule"
}
```


## Argument Reference

Panorama specific arguments:

* `device_group` - The device group (default: `shared`).
* `rulebase` - The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.


NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The rule name.


## Attribute Reference

The following attributes are supported:

* `rule` - The rule definition (see below).


Each `rule` defined supports the following attributes:

* `name` - The NAT rule's name.
* `uuid` - (PAN-OS 9.0+) The PAN-OS UUID.
* `group_tag` - (PAN-OS 9.0+) The group tag.
* `description` - The description.
* `type` - NAT type.
* `tags` - List of administrative tags.
* `disabled` - (bool) Set to `true` to disable this rule.
* `original_packet` - The original packet specification (see below).
* `translated_packet` - The translated packet spec (see below).


`original_packet` supports the following attributes:

* `source_zones` - The list of source zone(s).
* `destination_zone` - The destination zone.
* `destination_interface` - Egress interface from route lookup.
* `service` - Service.
* `source_addresses` - List of source address(es).
* `destination_addresses` - List of destination address(es).


`translated_packet` supports the following attributes:

* `source` - The source spec (see below).
* `destination` - The destination spec (see below).


`translated_packet.source` supports the following attributes:

* `dynamic_ip_and_port` - Dynamic IP and port source translation spec (see below).
* `dynamic_ip` - Dynamic IP source translation spec (see below).
* `static_ip` - Static IP source translation spec (see below).


`translated_packet.source.dynamic_ip_and_port` supports the following attributes:

* `translated_address` - Translated address source translation type spec (see below).
* `interface_address` - Interface address source translation type spec (see below).


`translated_packet.source.dynamic_ip_and_port.translated_address` supports
the following attributes:

* `translated_addresses` - List of translated addresses.


`translated_packet.source.dynamic_ip_and_port.interface_address` supports
the following attributes:

* `interface` - The interface.
* `ip_address` - The IP address.


`translated_packet.source.dynamic_ip` supports the following attributes:

* `translated_addresses` - The list of translated addresses.
* `fallback` - The fallback spec (see below).


`translated_packet.source.dynamic_ip.fallback` supports the following attributes:

* `translated_address` - The translated address fallback spec (see below).
* `interface_address` - The interface address fallback spec (see below).


`translated_packet.source.dynamic_ip.fallback.translated_address` supports the
following attributes:

* `translated_addresses` - List of source address translation
  fallback translated addresses.


`translated_packet.source.dynamic_ip.fallback.interface_address` supports the
following attributes:

* `interface` - Source address translation fallback interface.
* `type` - Type of interface fallback.
* `ip_address` - IP address of the fallback interface.


`translated_packet.source.static_ip` supports the following attributes:

* `translated_address` - The statically translated source address.
* `bi_directional` - (bool) Set to `true` to enable
  bi-directional source address translation.


`translated_packet.destination` supports the following attributes:

* `static_translation` - Specifies a static destination NAT (see below).
* `dynamic_translation` - (PAN-OS 8.1+) Specify a dynamic destination NAT
  (see below).
* `static` - (**DEPRECATED**) Specifies a static destination NAT (see below).
  This was deprecated in provider version 1.6 in favor of `static_translation` instead.
* `dynamic` - (**DEPRECATED**, PAN-OS 8.1+) Specify a dynamic destination NAT
  (see below).  If you are using Terraform 0.12+, you cannot use this param as it
  conflicts with the new
  [dynamic](https://www.terraform.io/docs/configuration/expressions.html#dynamic-blocks) block.


`translated_packet.destination.static` and `translated_packet.destination.static_translation`
support the following attributes:

* `address` - Destination address translation address.
* `port` - (int) Destination address translation port number.


`translated_packet.destination.dynamic` and `translated_packet.destination.dynamic_translation`
support the following attributes:

* `address` - Destination address translation address.
* `port` - (int) Destination address translation port number.
* `distribution` - (PAN-OS 8.1+) Distribution algorithm
  for destination address pool.  The PAN-OS 8.1 GUI doesn't seem to set this
  anywhere, but this is added here for completeness' sake.  The GUI sets
  this to `round-robin` currently.
