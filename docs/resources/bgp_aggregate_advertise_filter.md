---
page_title: "panos: panos_bgp_aggregate_advertise_filter"
subcategory: "Network"
---

# panos_bgp_aggregate_advertise_filter

This resource allows you to add/update/delete a route advertise filter for a
BGP address aggregation rule.


## PAN-OS

NGFW


## Import Name

```shell
<virtual_router>:<bgp_aggregate>:<name>
```


## Example Usage

```hcl
resource "panos_bgp_aggregate_advertise_filter" "example" {
    virtual_router = panos_bgp_aggregate.ag.virtual_router
    bgp_aggregate = panos_bgp_aggregate.ag.name
    name = "my advertise filter"
    as_path_regex = "*42*"
    med = "443"
    address_prefix {
        prefix = "10.1.1.0/24"
        exact = true
    }
    address_prefix {
        prefix = "10.1.2.0/24"
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_bgp_aggregate" "ag" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "addyAgg1"
    prefix = "192.168.1.0/24"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.rtr.name
    router_id = "5.5.5.5"
    as_number = "42"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_virtual_router" "rtr" {
    name = "my virtual router"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this filter to.
* `bgp_aggregate` - (Required) The BGP address aggregation rule.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `as_path_regex` - (Optional) AS path to match.
* `community_regex` - (Optional) Community to match.
* `extended_community_regex` - (Optional) Extended community to match.
* `med` - (Optional) Match MED.
* `route_table` - (Optional, PAN-OS 8.0+) Route table to match rule.  Valid
  values are `unicast`, `multicast`, or `both`.
* `address_prefix` - (Optional, repeatable) Matching address prefix definition
  (see below).
* `next_hops` - (Optional) List of next hop attributes.
* `from_peers` - (Optional) List of peers that advertised the route entry.

Each `address_prefix` section offers the following params:

* `prefix` - (Required) Address prefix.
* `exact` - (Optional, bool) Match exact prefix length.
