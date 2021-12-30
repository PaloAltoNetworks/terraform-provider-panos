---
page_title: "panos: panos_bgp_aggregate"
subcategory: "Network"
---

# panos_bgp_aggregate

This resource allows you to add/update/delete BGP address aggregation
rules.


## PAN-OS

NGFW


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_bgp_aggregate" "example" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "myAggRule"
    prefix = "192.168.1.0/24"
    weight = 17
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.vr.name
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
* `name` - (Required) The rule name.
* `prefix` - (Required) Aggregating address prefix.
* `enable` - (Optional, bool) Enable this rule (default: `true`)
* `as_set` - (Optional, bool) Generate AS-set attribute.
* `summary` - (Optional, bool) Summarize route.
* `local_preference` - (Optional) New local preference value.
* `med` - (Optional) New MED value.
* `weight` - (Optional, int) New weight value.
* `next_hop` - (Optional) Next hop address.
* `origin` - (Optional) New route origin.  Valid values are `incomplete`
  (default), `igp`, or `egp`.
* `as_path_limit` - (Optional, int) Add AS path limit attribute if it does
  not exist.
* `as_path_type` - (Optional) AS path update options.  Valid values are
  `none` (default) or `prepend`.
* `as_path_value` - (Optional) For `as_path_type` of `prepend`, the value to
  prepend.
* `community_type` - (Optional) Community update options.  Valid values are
  `none` (default), `remove-all`, `remove-regex`, `append`, or `overwrite`.
* `community_value` - (Optional) If `community_type` is `remove-regex`,
  `append`, or `overwrite`, the value associated with that setting.  For the
  `append` and `overwrite` types specifically, valid values are
  `no-export`, `no-advertise`, `local-as`, or `nopeer`.
* `extended_community_type` - (Optional) Extended community update options.  Valid
  values are `none` (default), `remove-all`, `remove-regex`, `append`, or `overwrite`.
* `extended_community_vaule` - (Optional) If `extended_community_type` is
  `remove-regex`, `append`, or `overwrite`, the value associated with that setting.
