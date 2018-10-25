---
layout: "panos"
page_title: "panos: panos_zone"
sidebar_current: "docs-panos-resource-zone"
description: |-
  Manages Zone objects.
---

# panos_zone

This resource allows you to add/update/delete zones.

This resource has some overlap with the `panos_zone_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_zone` spec does not define the
`interfaces` field.

## Example Usage

```hcl
resource "panos_zone" "example" {
    name = "my_service"
    mode = "layer3"
    interfaces = ["ethernet1/1", "ethernet1/2"]
    enable_user_id = true
    exclude_acls = ["192.168.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The zone's name.
* `vsys` - (Optional) The vsys to put the zone into (default: `vsys1`).
* `mode` - (Required) The zone's mode.  This can be `layer3`, `layer2`,
  `virtual-wire`, `tap`, or `tunnel`.
* `zone_profile` - (Optional) The zone protection profile.
* `log_setting` - (Optional) Log setting.
* `enable_user_id` - (Optional) Boolean to enable user identification.
* `interfaces` - (Optional) List of interfaces to associated with this zone.
* `include_acls` - (Optional) Users from these addresses/subnets will
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
* `exclude_acls` - (Optional) Users from these addresses/subnets will not
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
