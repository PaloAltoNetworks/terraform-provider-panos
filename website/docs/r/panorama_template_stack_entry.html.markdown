---
layout: "panos"
page_title: "panos: panos_panorama_template_stack_entry"
sidebar_current: "docs-panos-panorama-resource-template-stack-entry"
description: |-
  Manages a specific device in a Panorama template stack.
---

# panos_panorama_template_stack_entry

This resource allows you to add/update/delete a specific device in a Panorama
template stack.

This resource has some overlap with the `panos_panorama_template_stack`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_template_stack` spec does not define the
`devices` field.

This is the appropriate resource to use if you have a pre-existing template stack
in Panorama and don't want Terraform to delete it on `terraform destroy`.


## Import Name

```
<template_stack>:<device>
```


## Example Usage

```hcl
resource "panos_panorama_template_stack_entry" "example1" {
    template_stack = "my template stack"
    device = "00112233"
}
```

## Argument Reference

The following arguments are supported:

* `template_stack` - (Required) The template name.
* `device` - (Required) The serial number of the device to add.
