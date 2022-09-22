---
page_title: "panos: panos_panorama_template_stack"
subcategory: "Panorama"
---

# panos_panorama_template_stack

This resource allows you to add/update/delete Panorama template stacks.

This resource has some overlap with the
[`panos_panorama_template_stack_entry`](panorama_template_stack_entry.html)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_template_stack` spec does not define any
`device` blocks, and just stays as "computed".

This is the appropriate resource to use if `terraform destroy` should delete
the template stack.


## PAN-OS

Panorama


## Import Name

```shell
<name>
```


## Example Usage

```hcl
resource "panos_panorama_template_stack" "example" {
    name = "myStack"
    description = "description here"
    templates = ["t1", "t2"]
    devices = ["00112233", "44556677"]

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The stack's name.
* `description` - (Optional) The stack's description.
* `default_vsys` - (Optional) The default virtual system template configuration
  pushed to firewalls with a single virtual system.  **Note** - you can only
  set this if there is at least one template in this stack.
* `templates` - (Optional) List of templates in this stack.
* `devices` - (Optional) List of serial numbers to include in this stack.
