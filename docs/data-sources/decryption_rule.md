---
page_title: "panos: panos_decryption_rule_group"
subcategory: "Policies"
---

# panos_decryption_rule_group

Retrieve information on the given decryption rule.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
data "panos_decryption_rule_group" "example" {
    device_group = panos_panorama_device_group.x.name
    rulebase = "pre-rulebase"
    name = "some decryption rule"
}

resource "panos_panorama_device_group" "x" {
    name = "my device group"
}
```


## Argument Reference

Panorama specific arguments:

* `device_group` - (Optional) The device group (default: `shared`).
* `rulebase` - (Optional) The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.

NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The rule name.


## Attribute Reference

The following attributes are supported:

* `rule` - The rule definition (see below).

`rule` has the following attributes:

* `name` - The rule name.
* `uuid` - The PAN-OS UUID.
* `description` - The description.
* `source_zones` - List of source zones.
* `source_addresses` - List of source addresses.
* `negate_source` - (bool) If the source should be negated.
* `source_users` - List of source users.
* `destination_zones` - List of destination zones.
* `destination_addresses` - List of destination addresses.
* `negate_destination` - (bool) Negate the destination addresses.
* `tags` - List of administrative tags.
* `disabled` - (bool) Disable this rule.
* `services` - List of services.
* `url_categories` - List of URL categories.
* `action` - Action to take.
* `decryption_type` - The decryption type.
* `ssl_certificate` - The SSL certificate.
* `decryption_profile` - The decryption profile.
* `forwarding_profile` - Forwarding profile.
* `group_tag` - The group tag.
* `source_hips` - List of source HIP devices.
* `destination_hips` - List of destination HIP devices.
* `log_successful_tls_handshakes` - Log successful TLS handshakes.
* `log_failed_tls_handshakes` - Log failed TLS handshakes.
* `log_setting` - The log setting.
* `target` - (repeatable, Panorama only) A target definition (see below).  If there
  are no target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (bool, Panorama only) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.

`rule.target` supports the following arguments:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - A listing of vsys to apply this rule to.  If `serial` is
  a VM, then this parameter should just be omitted.
