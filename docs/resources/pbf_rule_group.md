---
page_title: "panos: panos_pbf_rule_group"
subcategory: "Policies"
---

# panos_pbf_rule_group

This resource allows you to add/update/delete policy based forwarding rule groups.

This resource manages clusters of policy based forwarding rules in a single vsys,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Although you cannot modify non-group PBF rules with this
resource, the `position_keyword` and `position_reference` parameters allow you
to reference some other PBF rule that already exists, using it as
a means to ensure some rough placement within the ruleset as a whole.


## Best Practices

As is to be expected, if you are separating your deployment across
multiple plan files, make sure that at most only one plan specifies any given
absolute positioning keyword such as "top" or "directly below", otherwise
they'll keep shoving each other out of the way indefinitely.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_pbf_rule_group`


## Example Usage

```hcl
resource "panos_pbf_rule_group" "example" {
    position_keyword = "above"
    position_reference = "deny everything else"
    rule {
        name = "my-pbf"
        audit_comment = "Initial deploy"
        description = "deployed by terraform"
        source {
            zones = [panos_zone.foo.name]
            addresses = ["10.50.50.50"]
            users = ["any"]
            negate = true
        }
        destination {
            addresses = ["10.80.80.80"]
            applications = ["any"]
            services = ["application-default"]
        }
        forwarding {
            action = "discard"
        }
    }
}

resource "panos_zone" "foo" {
    name = "myZone"
    mode = "layer2"
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

* `position_keyword` - A positioning keyword for this group.  This
  can be `before`, `directly before`, `after`, `directly after`, `top`,
  `bottom`, or left empty (the default) to have no particular placement.  This
  param works in combination with the `position_reference` param.
* `position_reference` - Required if `position_keyword` is one of the
  "above" or "below" variants, this is the name of a non-group rule to use
  as a reference to place this group.
* `rule` - The rule definition (see below).  The rule
  ordering will match how they appear in the terraform plan file.

The following arguments are valid for each `rule` section:

* `name` - (Required) The rule name.
* `audit_comment` - When this rule is created/updated, the audit comment to
  apply for this rule.
* `group_tag` - (PAN-OS 9.0+) The group tag.
* `description` - The rule description.
* `tags` - List of tags for this rule.
* `active_active_device_binding` - The active-active device binding.
* `schedule` - The schedule.
* `disabled` - (bool) Set to `true` to disable this rule.
* `source` - (Required) The source spec (defined below).
* `destination` - (Required) The destination spec (defined below).
* `forwarding` - (Required) The forwarding spec (defined below).
* `target` - A target definition (see below).  If there are no
  target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (bool) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.

`rule.source` supports the following arguments:

* `zones` - If you want a source type of "zone", then define this
  list with the desired source zones.  Mutually exclusive with `rule.interfaces`.
* `interfaces` - If you want a source type of "interface", then define this
  list with the desired source interfaces.  Mutually exclusive with `rule.zones`.
* `addresses` - (Required) List of source IP addresses.
* `users` - (Required) List of source users.
* `negate` - (bool) Set to `true` to negate the source.

`rule.destination` supports the following arguments:

* `addresses` - (Required) The list of destination addresses.
* `application` - (Required) The list of applications.
* `services` - (Required) The list of services.
* `negate` - (bool) Set to `true` to negate the destination.

`rule.forwarding` supports the following arguments:

* `action` - The action to take.  Valid values are `forward` (default),
  `forward-to-vsys`, `discard`, or `no-pbf`.
* `vsys` - If `action=forward-to-vsys`, the vsys to forward to.
* `egress_interface` - If `action=forward`, the egress interface.
* `next_hop_type` - If `action=forward`, the next hop type.  Valid values
  are `ip-address`, `fqdn`, or leaving this empty for a next hop type of None.
* `next_hop_value` - If `action=forward` and `next_hop_type` is defined, then
  the next hop address.
* `monitor` - The monitor spec (defined below).  If you do not want to enable
  monitoring, then do not specify a `monitor` config block.
* `symmetric_return` - The symmetric return spec (defined below).  If you do
  not want to enforce symmetric

`rule.forwarding.monitor` supports the following arguments:

* `profile` - The montior profile to use.
* `ip_address` - The monitor IP address.
* `disable_if_unreachable` - (bool) Set to `true` to disable this rule if
  nexthop/monitor IP is unreachable.

`rule.forwarding.symmetric_return` supports the following arguments:

* `enable` - (bool) Set to `true` to enforce symmetric return.
* `addresses` - List of next hop addresses.

`rule.target` supports the following arguments:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.


## Attribute Reference

Each `rule` supports the following attributes:

* `uuid` - (PAN-OS 9.0+) The PAN-OS UUID.
