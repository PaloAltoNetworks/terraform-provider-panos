---
page_title: "panos: panos_bgp_conditional_adv"
subcategory: "Network"
---

# panos_bgp_conditional_adv

This resource allows you to add/update/delete a BGP conditional advertisement.

~> **Note:** In the PAN-OS GUI, this resource cannot be created without also
creating at least one non-exist filter and one advertise filter.  The API behaves
a little differently:  you can create the conditional advertisement itself, but
the API will start throwing errors if you try to update it and there is not at
least one non-exist filter and one advertise filter.  In order for a conditional
advertisement to be valid, you must specify at least one non-exist and one
advertise filter.


## PAN-OS

NGFW


## Import Name

```shell
<virtual_router>:<name>
```


## Example Usage

```hcl
data "panos_system_info" "x" {}

resource "panos_bgp_conditional_adv" "example" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "example"
    enable = false

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_bgp_conditional_adv_non_exist_filter" "nef" {
    virtual_router = panos_bgp.conf.virtual_router
    bgp_conditional_adv = panos_bgp_conditional_adv.example.name
    route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    name = "nef"
    address_prefixes = ["192.168.1.0/24"]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_bgp_conditional_adv_advertise_filter" "af" {
    virtual_router = panos_bgp.conf.virtual_router
    bgp_conditional_adv = panos_bgp_conditional_adv.example.name
    route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    name = "af"
    address_prefixes = ["192.168.2.0/24"]

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

* `virtual_router` - (Required) The virtual router to add this BGP
  conditional advertisement to.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `used_by` - (Optional) List of BGP peer groups that use this rule.
