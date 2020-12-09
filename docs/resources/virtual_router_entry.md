---
page_title: "panos: panos_virtual_router_entry"
subcategory: "Firewall Networking"
---

# panos_virtual_router_entry

This resource allows you to add/update/delete an interface in a
virtual router.

This resource has some overlap with the `panos_virtual_router`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_virtual_router` spec does not define the
`interfaces` field.


## Import Name

```
<virtual_router>:<interface>
```


## Example Usage

```hcl
resource "panos_virtual_router_entry" "example" {
    virtual_router = panos_virtual_router.vr.name
    interface = panos_ethernet_interface.e.name
}

resource "panos_virtual_router" "vr" {
    name = "my vr"
}

resource "panos_ethernet_interface" "e" {
    name = "ethernet1/1"
    mode = "layer3"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router's name.
* `interface` - (Required) The interface to import into the virtual router.
