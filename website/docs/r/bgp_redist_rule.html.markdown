---
layout: "panos"
page_title: "panos: panos_bgp_redist_rule"
sidebar_current: "docs-panos-resource-bgp-redist-rule"
description: |-
  Manages a BGP redistribution rule.
---

# panos_bgp_redist_rule

This resource allows you to add/update/delete a BGP redistribution rule.


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_bgp_redist_rule" "example" {
    virtual_router = "${panos_bgp.conf.virtual_router}"
    route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    name = "192.168.1.0/24"
    set_med = "42"
}

data "panos_system_info" "x" {}

resource "panos_bgp" "conf" {
    virtual_router = "${panos_virtual_router.rtr.name}"
    router_id = "5.5.5.5"
    as_number = "42"
}

resource "panos_virtual_router" "rtr" {
    name = "my virtual router"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  redist rule to.
* `name` - (Required) A subnet or a redistribution profile.
* `enable` - (Optional, bool) Enable this rule or not (default: `true`).
* `address_family` - (Optional) The address family.  Valid values are
  `ipv4` (default) or `ipv6`.
* `route_table` - (Optional, PAN-OS 8.0+) Route table to match rule.  Valid
  values are `unicast`, `multicast`, or `both`.  As of PAN-OS 8.1, there doesn't
  seem to be a way to configure this in the GUI, it is always set to `unicast`.
  Thus, if you're running this resource against PAN-OS 8.0+, the appropriate
  thing to do is set this value to `unicast` as well to match the GUI functionality.
* `metric` - (Optional, int) Metric value.
* `set_origin` - (Optional) Add the origin path attribute.  Valid values are
  `incomplete` (default), `igp`, or `egp`.
* `set_med` - (Optional) Add the MULTI_EXIT_DISC path attribute.
* `set_local_preference` - (Optional) Add the LOCAL_PREF path attribute.
* `set_as_path_limit` - (Optional, int) Add the AS_PATHLIMIT path attribute.
* `set_communities` - (Optional) List of COMMUNITY path attributes to add.
* `set_extended_communities` - (Optional) List of EXTENDED COMMUNITY path attributes to add.
