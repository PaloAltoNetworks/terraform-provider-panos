---
layout: "panos"
page_title: "panos: panos_bfd_profile"
sidebar_current: "docs-panos-resource-bfd-profile"
description: |-
  Manages BFD profiles.
---

# panos_bfd_profile.

This resource allows you to add/update/delete BFD profiles.

~> **Note:** This resource is only applicable for PAN-OS 7.1+.


## Import Name

```
<name>
```


## Example Usage

```hcl
resource "panos_bfd_profile" "example" {
    name = "myBfdProfile"
}
```

## Argument Reference

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
