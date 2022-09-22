---
page_title: "panos: panos_panorama_bfd_profile"
subcategory: "Network"
---

# panos_panorama_bfd_profile.

This resource allows you to add/update/delete BFD profiles on Panorama.


## Minimum PAN-OS Version

7.1


## PAN-OS

Panorama


## Import Name

```shell
<template>:<template_stack>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_bfd_profile" "example" {
    template = panos_panorama_template.t.name
    name = "myBfdProfile"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The BBFD profile's name.
* `mode` - (Optional) BFD operation mode.  Valid values are `active` (default)
  or `passive`.
* `minimum_tx_interval` - (Optional, int) Desired minimum TX interval in
  ms.  Default is `1000`.
* `minimum_rx_interval` - (Optional, int) Required minimum RX interval in
  ms.  Default is `1000`.
* `detection_multiplier` - (Optional, int) Multiplier sent to remote
  system.  Default is `3`.
* `hold_time` - (Optional, int) Delay transmission and reception of control
  packets in ms.
* `minimum_rx_ttl` - (Optional, int) Minimum accepted ttl on received BFD
  packet.
