---
page_title: "panos: panos_virtual_router"
subcategory: "Firewall Networking"
---

# panos_virtual_router

This resource allows you to add/update/delete virtual routers.

**Note** - The `default` virtual router may be configured with this resource,
however it will not be deleted from the firewall.  It will only be unexported
from the vsys that it is currently imported in, and any interfaces imported
into the virtual router will be removed.

This resource has some overlap with the `panos_virtual_router_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_virtual_router` spec does not define the
`interfaces` field.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
# Configure a bare-bones ethernet interface.
resource "panos_virtual_router" "example" {
    name = "my virtual router"
    static_dist = 15
    interfaces = [
        panos_ethernet_interface.e1.name,
        panos_ethernet_interface.e2.name,
    ]
}

resource "panos_ethernet_interface" "e1" {
    vsys = "vsys1"
    name = "ethernet1/1"
    mode = "layer3"
}

resource "panos_ethernet_interface" "e2" {
    vsys = "vsys1"
    name = "ethernet1/2"
    mode = "layer3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The virtual router's name.
* `vsys` - (Required) The vsys that will use this virtual router.  This should
  be something like `vsys1` or `vsys3`.
* `interfaces` - (Optional) List of interfaces that should use this virtual
  router.  If you intend to use `panos_virtual_router_entry` then leave this
  field undefined.
* `static_dist` - (Optional) Admin distance - Static (default: `10`).
* `static_ipv6_dist` - (Optional) Admin distance - Static IPv6 (default: `10`).
* `ospf_int_dist` - (Optional) Admin distance - OSPF Int (default: `30`).
* `ospf_ext_dist` - (Optional) Admin distance - OSPF Ext (default: `110`).
* `ospfv3_int_dist` - (Optional) Admin distance - OSPFv3 Int (default:
  `30`).
* `ospfv3_ext_dist` - (Optional) Admin distance - OSPFv3 Ext (default:
  `110`).
* `ibgp_dist` - (Optional) Admin distance - IBGP (default: `200`).
* `ebgp_dist` - (Optional) Admin distance - EBGP (default: `20`).
* `rip_dist` - (Optional) Admin distance - RIP (default: `120`).
