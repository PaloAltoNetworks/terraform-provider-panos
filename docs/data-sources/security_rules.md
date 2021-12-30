---
page_title: "panos: panos_security_rules"
subcategory: "Policies"
---

# panos_security_rules

Retrieves the list of security rules present.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
data "panos_security_rules" "example" {}
```


## Argument Reference

Panorama specific arguments:

* `device_group` - (Optional) The device group (default: `shared`).
* `rulebase` - (Optional) The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.

NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
