---
layout: "panos"
page_title: "panos: panos_panorama_nat_rule_group"
sidebar_current: "docs-panos-panorama-resource-nat-rule-group"
description: |-
  Manages a group of Panorama NAT rules.
---

# panos_panorama_nat_rule_group

This resource allows you to add/update/delete a group of Panorama NAT rules.

This resource manages clusters of NAT rules in a single device group,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Although you cannot modify non-group NAT rules with this
resource, the `position_keyword` and `position_reference` parameters allow you
to reference some other NAT rule that already exists, using it as
a means to ensure some rough placement within the ruleset as a whole.


## Best Practices

As is to be expected, if you are separating your deployment across
multiple plan files, make sure that at most only one plan specifies any given
absolute positioning keyword such as "top" or "directly below", otherwise
they'll keep shoving each other out of the way indefinitely.

Best practices are to specify one group as `top` (if you need it), one
group as `bottom` (if needed), then
all other groups should be `above` the first rule of the bottom group.  You
do it this way because rules will natually be added at the tail end of the
rulebase, so they will always be `after` the first group, but what you want
is for them to be `before` the last group's rules.


## Example Usage

```hcl
resource "panos_panorama_nat_rule_group" "bot" {
    device_group = "${panos_panorama_device_group.dg.name}"
    rule {
        name = "second"
        original_packet {
            source_zones = ["z2"]
            destination_zone = "z3"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {}
            destination {
                static_translation {
                    address = "10.2.3.1"
                    port = 5678
                }
            }
        }
    }
    rule {
        name = "third"
        original_packet {
            source_zones = ["z3"]
            destination_zone = "z2"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {
                static_ip {
                    translated_address = "192.168.1.5"
                    bi_directional = true
                }
            }
            destination {}
        }
    }
}

resource "panos_panorama_nat_rule_group" "top" {
    device_group = "${panos_panorama_device_group.dg.name}"
    position_keyword = "directly before"
    position_reference = "${panos_panorama_nat_rule_group.bot.rule.0.name}"
    rule {
        name = "first"
        target {
            serial = "123456"
            vsys_list = ["vsys1", "vsys2"]
        }
        original_packet {
            source_zones = ["z1"]
            destination_zone = "z1"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {
                dynamic_ip_and_port {
                    interface_address {
                        interface = "ethernet1/6"
                    }
                }
            }
            destination {
                static_translation {
                    address = "10.1.1.1"
                    port = 1234
                }
            }
        }
    }
}

resource "panos_panorama_device_group" "dg" {
    name = "myDeviceGroup"
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys to put the NAT rule group into (default:
  `vsys1`).
* `device_group` - (Optional) Device group the NAT rules should be put into
  (default: `shared`).
* `rulebase` - (Optional) The rulebase the NAT rules should be put into.  Valid
  values are `pre-rulebase` (default), `rulebase`, or `post-rulebase`.
* `position_keyword` - (Optional) A positioning keyword for this group.  This
  can be `before`, `directly before`, `after`, `directly after`, `top`,
  `bottom`, or left empty (the default) to have no particular placement.  This
  param works in combination with the `position_reference` param.
* `position_reference` - (Optional) Required if `position_keyword` is one of the
  "above" or "below" variants, this is the name of a non-group rule to use
  as a reference to place this group.
* `rule` - (Repeatable) The rule definition (see below).  The rule
  ordering will match how they appear in the terraform plan file.


Each `rule` defined supports the following arguments:

* `name` - (Required) The NAT rule's name.
* `description` - (Optional) The description.
* `type` - (Optional). NAT type.  This can be `ipv4` (default), `nat64`, or
  `nptv6`.
* `tags` - (Optional) List of administrative tags.
* `disabled` - (Optional) Set to `true` to disable this rule.
* `target` - (Optional, repeatable) A target definition (see below).  If there are
  no target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (Optional, bool) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.
* `original_packet` - (Required) The original packet specification (see below).
* `translated_packet` - (Required) The translated packet spec (see below).


`target` supports the following arguments:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.


`original_packet` supports the following arguments:

* `source_zones` - (Required) The list of source zone(s).
* `destination_zone` - (Required) The destination zone.
* `destination_interface` - (Optional) Egress interface from route lookup (default:
  `any`).
* `service` - (Optional) Service (default: `any`).
* `source_addresses` - (Required) List of source address(es).
* `destination_addresses` - (Required) List of destination address(es).


`translated_packet` supports the following arguments:

* `source` - (Required) The source spec (see below).  Leave this
  empty for a destination NAT of "none".
* `destination` - (Required) The destination spec (see below).  Leave this
  empty for a destination NAT of "none".


`translated_packet.source` supports the following arguments:

* `dynamic_ip_and_port` - (Optional) Dynamic IP and port source translation
  spec (see below).
* `dynamic_ip` - (Optional) Dynamic IP source translation spec (see below).
* `static_ip` - (Optional) Static IP source translation spec (see below).


`translated_packet.source.dynamic_ip_and_port` supports the following arguments:

* `translated_address` - (Optional) Translated address source translation type
  spec (see below).
* `interface_address` - (Optional) Interface address source translation type
  spec (see below).


`translated_packet.source.dynamic_ip_and_port.translated_address` supports
the following arguments:

* `translated_addresses` - (Required) List of translated addresses.


`translated_packet.source.dynamic_ip_and_port.interface_address` supports
the following arguments:

* `interface` - (Required) The interface.
* `ip_address` - (Optional) The IP address.


`translated_packet.source.dynamic_ip` supports the following arguments:

* `translated_addresses` - (Optional) The list of translated addresses.
* `fallback` - (Optional) The fallback spec (see below).  Leaving this empty
  or omiting it means a fallback of "None".


`translated_packet.source.dynamic_ip.fallback` supports the following arguments:

* `translated_address` - (Optional) The translated address fallback spec (see below).
* `interface_address` - (Optional) The interface address fallback spec (see below).


`translated_packet.source.dynamic_ip.fallback.translated_address` supports the
following arguments:

* `translated_addresses` - (Optional) List of source address translation
  fallback translated addresses.


`translated_packet.source.dynamic_ip.fallback.interface_address` supports the
following arguments:

* `interface` - (Required) Source address translation fallback interface.
* `type` - (Optional) Type of interface fallback.  Valid values are
  `ip` (default) or `floating`.
* `ip_address` - (Optional) IP address of the fallback interface.


`translated_packet.source.static_ip` supports the following arguments:

* `translated_address` - (Required) The statically translated source
  address.
* `bi_directional` - (Optional, bool) Set to `true` to enable
  bi-directional source address translation.


`translated_packet.destination` supports the following arguments:

* `static_translation` - (Optional) Specifies a static destination NAT (see below).
* `dynamic_translation` - (Optional, PAN-OS 8.1+) Specify a dynamic destination NAT
  (see below).
* `static` - (**DEPRECATED**, Optional) Specifies a static destination NAT (see below).
  This was deprecated in provider version 1.6 in favor of `static_translation` instead.
* `dynamic` - (**DEPRECATED**, Optional, PAN-OS 8.1+) Specify a dynamic destination NAT
  (see below).  If you are using Terraform 0.12+, you cannot use this param as it
  conflicts with the new
  [dynamic](https://www.terraform.io/docs/configuration/expressions.html#dynamic-blocks) block.


`translated_packet.destination.static` and `translated_packet.destination.static_translation`
support the following arguments:

* `address` - (Required) Destination address translation address.
* `port` - (Optional, int) Destination address translation port number.


`translated_packet.destination.dynamic` and `translated_packet.destination.dynamic_translation`
support the following arguments:

* `address` - (Required) Destination address translation address.
* `port` - (Optional, int) Destination address translation port number.
* `distribution` - (Optional, PAN-OS 8.1+) Distribution algorithm
  for destination address pool.  The PAN-OS 8.1 GUI doesn't seem to set this
  anywhere, but this is added here for completeness' sake.  The GUI sets
  this to `round-robin` currently.
