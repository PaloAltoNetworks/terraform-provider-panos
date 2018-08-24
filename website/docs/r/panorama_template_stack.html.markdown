---
layout: "panos"
page_title: "panos: panos_panorama_template_stack"
sidebar_current: "docs-panos-panorama-resource-template-stack"
description: |-
  Manages Panorama template stacks.
---

# panos_panorama_template_stack

This resource allows you to add/update/delete Panorama template stacks.

This resource has some overlap with the `panos_panorama_template_stack_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_template_stack` spec does not define any
`device` blocks, and just stays as "computed".

This is the appropriate resource to use if `terraform destroy` should delete
the template stack.

## Example Usage

```hcl
resource "panos_panorama_template_stack" "example" {
    name = "myStack"
    description = "description here"
    templates = ["t1", "t2"]
    devices = ["00112233", "44556677"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The stack's name.
* `description` - (Optional) The stack's description.
* `templates` - (Optional) List of templates in this stack.
* `devices` - (Optional) List of serial numbers to include in this stack.

The following attributes are present:

* `default_vsys` - The default vsys for this template stack.
