---
layout: "panos"
page_title: "panos: panos_panorama_template_variable"
sidebar_current: "docs-panos-panorama-resource-template-variable"
description: |-
  Manages Panorama template variables.
---

# panos_panorama_template_variable

This resource allows you to add/update/delete variables for both Panorama
templates and template stacks.

Template variables are available in PAN-OS 8.1+.

## Example Usage

```hcl
resource "panos_panorama_template_variable" "example" {
    template = "${panos_panorama_template.tmpl1.name}"
    name = "$example"
    type = "ip-address"
    value = "10.1.1.1/24"
}

resource "panos_panorama_template" "tmpl1" {
    name = "MyTemplate"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The template's name.  This must start with a dollar sign ($).
* `type` - (Optional) The variable type.  Valid values are `ip-netmask`
  (default), `ip-range`, `fqdn`, `group-id`, or `interface`.
* `value` - (Required) The variable value.
