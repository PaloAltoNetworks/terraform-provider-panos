---
page_title: "panos: panos_panorama_template_stack_entry"
subcategory: "Panorama"
---

# panos_panorama_template_stack_entry

This resource allows you to add/update/delete a specific device in a Panorama
template stack.

This resource has some overlap with the
[`panos_panorama_template_stack`](panorama_template_stack.html)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_template_stack` spec does not define the
`devices` field.

This is the appropriate resource to use if you have a pre-existing template stack
in Panorama and don't want Terraform to delete it on `terraform destroy`.


## PAN-OS

Panorama


## Import Name

```shell
<template_stack>:<device>
```


## Example Usage

```hcl
resource "panos_panorama_template_stack_entry" "example1" {
    template_stack = panos_panorama_template_stack.t.name
    device = "00112233"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template_stack" "t" {
    name = "my template stack"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `template_stack` - (Required) The template name.
* `device` - (Required) The serial number of the device to add.
