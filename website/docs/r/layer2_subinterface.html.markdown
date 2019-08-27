---
layout: "panos"
page_title: "panos: panos_layer2_subinterface"
description: |-
  Manages layer2 subinterfaces.
---

# panos_layer2_subinterface

This resource allows you to add/update/delete layer2 subinterfaces.


## Import Name

```
<interface_type>:<parent_interface>:<parent_mode>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_layer2_subinterface" "example" {
    parent_interface = panos_ethernet_interface.e.name
    parent_mode = panos_ethernet_interface.e.mode
    vsys = "vsys1"
    name = "ethernet1/5.5"
    tag = 5
}

resource "panos_ethernet_interface" "e" {
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer2"
}
```

## Argument Reference

The following arguments are supported:

* `interface_type` - (Optional) The interface type.  Valid values are `ethernet` (default)
  or `aggregate-ethernet`.
* `parent_interface` - (Required) The name of the parent interface.
* `parent_mode` - (Optional) The parent's mode.  Valid values are `layer2` (default)
  or `virtual-wire`.
* `vsys` - (Required) The vsys that will use this interface.  This should be
  something like `vsys1` or `vsys3`.
* `name` - (Required) The interface's name.
* `tag` - (Optional, int) The interface's tag.
* `netflow_profile` - (Optional) The netflow profile.
* `comment` - (Optional) The interface comment.
