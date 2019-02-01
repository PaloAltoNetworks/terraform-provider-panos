---
layout: "panos"
page_title: "panos: panos_bgp_import_rule_group"
sidebar_current: "docs-panos-resource-bgp-import-rule-group"
description: |-
  Manages BGP import rule groups.
---

# panos_bgp_import_rule_group

This resource allows you to add/update/delete BGP import rule groups.

This resource manages clusters of import rules in a virtual router,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Although you cannot modify non-group import rules with this
resource, the `position_keyword` and `position_reference` parameters allow you
to reference some other import rule that already exists, using it as
a means to ensure some rough placement within the ruleset as a whole.


## Best Practices

As is to be expected, if you are separating your deployment across
multiple plan files, make sure that at most only one plan specifies any given
absolute positioning keyword such as "top" or "directly below", otherwise
they'll keep shoving each other out of the way indefinitely.

Best practices are to specify one group as `top` (if you need it), one
group as `bottom`, then
all other groups should be `above` the first rule of the bottom group.  You
do it this way because rules will natually be added at the tail end of the
ruleset, so they will always be `after` the first group, but what you want
is for them to be `before` the last group's rules.


## Example Usage

```hcl
resource "panos_bgp_import_rule_group" "example" {
    virtual_router = "${panos_bgp.conf.virtual_router}"
    rule {
        name = "first"
        match_as_path_regex = "*foo*"
        match_address_prefix {
            prefix = "192.168.1.0/24"
        }
        match_address_prefix {
            prefix = "192.168.2.0/24"
            exact = true
        }
        match_route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
        local_preference = "42"
        med = "43"
        weight = 44
        origin = "incomplete"
    }
    rule {
        name = "second"
        match_as_path_regex = "*bar*"
        action = "deny"
        match_route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    }
}

data "panos_system_info" "x" {}

resource "panos_bgp" "conf" {
    virtual_router = "${panos_virtual_router.vr.name}"
    router_id = "1.2.3.4"
    as_number = 443
}

resource "panos_virtual_router" "vr" {
    name = "my vr"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to put the rule into.
* `position_keyword` - (Optional) A positioning keyword for this group.  This
  can be `before`, `directly before`, `after`, `directly after`, `top`,
  `bottom`, or left empty (the default) to have no particular placement.  This
  param works in combination with the `position_reference` param.
* `position_reference` - (Optional) Required if `position_keyword` is one of the
  "above" or "below" variants, this is the name of a non-group rule to use
  as a reference to place this group.
* `rule` - The import rule definition (see below).  The import rule
  ordering will match how they appear in the terraform plan file.

The following arguments are valid for each `rule` section:

* `name` - (Required) The security rule name.
* `enable` - (Optional, bool) Enable this import rule (default: `true`)
* `used_by` - (Optional) List of auth profiles.
* `match_as_path_regex` - (Optional) AS path to match.
* `match_community_regex` - (Optional) Community to match.
* `match_extended_community_regex` - (Optional) Extended community to match.
* `match_med` - (Optional) Match MED.
* `match_route_table` - (Optional, PAN-OS 8.0+) Route table to match rule.  Valid
  values are `unicast`, `multicast`, or `both`.  As of PAN-OS 8.1, there doesn't
  seem to be a way to configure this in the GUI, it is always set to `unicast`.
  Thus, if you're running this resource against PAN-OS 8.0+, the appropriate
  thing to do is set this value to `unicast` as well to match the GUI functionality.
* `match_address_prefix` - (Optional, repeatable) Matching address prefix definition
  (see below).
  below for the params for this section.
* `match_next_hops` - (Optional) List of next hop attributes.
* `match_from_peers` - (Optional) List of peers that advertised the route entry.
* `action` - (Optional) Rule action.  Valid values are `allow` (default) or
  `deny`.
* `dampening` - (Optional) Route flap dampening profile.
* `local_preference` - (Optional) New local preference value.
* `med` - (Optional) New MED value.
* `weight` - (Optional, int) New weight value.
* `next_hop` - (Optional) Next hop address.
* `origin` - (Optional) New route origin.  Valid values are `igp`, `egp`, or
  `incomplete`.
* `as_path_limit` - (Optional, int) Add AS path limit attribute if it does
  not exist.
* `as_path_type` - (Optional) AS path update options.  Valid values are
  `none` or `remove`.
* `community_type` - (Optional) Community update options.  Valid values are
  `none`, `remove-all`, `remove-regex`, `append`, or `overwrite`.
* `community_value` - (Optional) If `community_type` is `remove-regex`,
  `append`, or `overwrite`, the value associated with that setting.  For the
  `append` and `overwrite` types specifically, valid values for `community_value`
  are `no-export`, `no-advertise`, `local-as`, or `nopeer`.
* `extended_community_type` - (Optional) Extended community update options.  Valid
  values are `none`, `remove-all`, `remove-regex`, `append`, or `overwrite`.
* `extended_community_vaule` - (Optional) If `extended_community_type` is
  `remove-regex`, `append`, or `overwrite`, the value associated with that setting.

Each `match_address_prefix` section offers the following params:

* `prefix` - (Required) Address prefix.
* `exact` - (Optional, bool) Match exact prefix length.
