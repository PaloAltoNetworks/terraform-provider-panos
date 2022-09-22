---
page_title: "panos: panos_panorama_template"
subcategory: "Panorama"
---

# panos_panorama_template

This resource allows you to add/update/delete Panorama templates.

This resource has some overlap with the
[`panos_panorama_template_entry`](panorama_template_entry.html)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_template` spec does not define any
`device` blocks, and just stays as "computed".

This is the appropriate resource to use if `terraform destroy` should delete
the template.

**Note** - In PAN-OS 8.1, it looks like the `devices` field has
been removed.  Creating a template stack and specifying devices in the template
stack is still present in PAN-OS 8.1.


## PAN-OS

Panorama


## Import Name

```shell
<name>
```


## Example Usage

```hcl
# This specifies one or more device blocks, so this is applicable only for
# PAN-OS 8.0 and lower.
resource "panos_panorama_template" "example" {
    name = "template1"
    description = "description here"
    device {
        serial = "00112233"
    }
    device {
        serial = "44556677"
        vsys_list = ["vsys1", "vsys2"]
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The template's name.
* `description` - (Optional) The template's description.
* `device` - The device definition (see below).

The following arguments are valid for each `device` section:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this template.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.
