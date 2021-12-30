---
page_title: "panos: panos_pbf_rule_group"
subcategory: "Policies"
---

# panos_pbf_rule_group

Retrieves information on the specifiec policy based forwarding rule.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_pbf_rule_group" "example" {
    name = "my-pbf"
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

The `rule` supports the following attributes:

* `name` - The rule name.
* `uuid` - (PAN-OS 9.0+) The PAN-OS UUID.
* `group_tag` - (PAN-OS 9.0+) The group tag.
* `description` - The rule description.
* `tags` - List of tags for this rule.
* `active_active_device_binding` - The active-active device binding.
* `schedule` - The schedule.
* `disabled` - (bool) Set to `true` to disable this rule.
* `source` - The source spec (defined below).
* `destination` - The destination spec (defined below).
* `forwarding` - The forwarding spec (defined below).

`rule.source` supports the following attributes:

* `zones` - List of source zones.
* `interfaces` - List of source interfaces.
* `addresses` - List of source IP addresses.
* `users` - List of source users.
* `negate` - (bool) Set to `true` to negate the source.

`rule.destination` supports the following attributes:

* `addresses` - (Required) The list of destination addresses.
* `application` - (Required) The list of applications.
* `services` - (Required) The list of services.
* `negate` - (bool) Set to `true` to negate the destination.

`rule.forwarding` supports the following attributes:

* `action` - The action to take.
* `vsys` - If `action=forward-to-vsys`, the vsys to forward to.
* `egress_interface` - If `action=forward`, the egress interface.
* `next_hop_type` - If `action=forward`, the next hop type.
* `next_hop_value` - If `action=forward` and `next_hop_type` is defined, then
  the next hop address.
* `monitor` - The monitor spec (defined below) if monitoring is enabled.
* `symmetric_return` - The symmetric return spec (defined below) if it's enforced.

`rule.forwarding.monitor` supports the following attributes:

* `profile` - The montior profile to use.
* `ip_address` - The monitor IP address.
* `disable_if_unreachable` - (bool) Set to `true` to disable this rule if
  nexthop/monitor IP is unreachable.

`rule.forwarding.symmetric_return` supports the following attributes:

* `enable` - (bool) Set to `true` to enforce symmetric return.
* `addresses` - List of next hop addresses.
