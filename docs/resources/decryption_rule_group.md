---
page_title: "panos: panos_decryption_rule_group"
subcategory: "Policies"
---

# panos_decryption_rule_group

This resource allows you to add/update/delete decryption rule groups.

This resource manages clusters of decryption rules in a single vsys,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Because this resource only manages what it's told to, it will not manage
any rules that may already exist on the firewall.  This has
implications on the effective decryption posture of your firewall, but it
will allow you to spread your decryption rules across multiple Terraform
state files.

Although you cannot modify non-group decryption rules with this
resource, the `position_keyword` and `position_reference` parameters allow you
to reference some other decryption rule that already exists, using it as
a means to ensure some rough placement within the ruleset as a whole.


## Best Practices

As is to be expected, if you are separating your deployment across
multiple plan files, make sure that at most only one plan specifies any given
absolute positioning keyword such as "top" or "directly below", otherwise
they'll keep shoving each other out of the way indefinitely.

Best practices are to specify one group as `top` (if you need it), one
group as `bottom` (this is where you have your logging deny rule), then
all other groups should be `above` the first rule of the bottom group.  You
do it this way because rules will natually be added at the tail end of the
rulebase, so they will always be `after` the first group, but what you want
is for them to be `before` the last group's rules.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
# Panorama Example
resource "panos_decryption_rule_group" "example1" {
    device_group = panos_device_group.x.name
    rulebase = "pre-rulebase"
    position_keyword = "top"
    rule {
        name = "sampleRule"
        audit_comment = "Case id 12345"
        description = "Made by Terraform"
        source_zones = ["any"]
        source_addresses = ["192.168.10.15"]
        source_users = ["any"]
        source_hips = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["10.20.30.40"]
        destination_hips = ["any"]
        negate_destination = true
        services = ["application-default"]
        url_categories = [
            "adult",
            "dating",
        ]
        action = "decrypt"
        decryption_type = "ssl-forward-proxy"
        log_successful_tls_handshakes = true
        log_failed_tls_handshakes = true
        target {
            serial = "01234"
        }
        target {
            serial = "56789"
            vsys_list = [
                "vsys1",
                "vsys3",
            ]
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_device_group" "x" {
    name = "my device group"

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama specific arguments:

* `device_group` - The device group (default: `shared`).
* `rulebase` - The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.

NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `position_keyword` - A positioning keyword for this group.  This
  can be `before`, `directly before`, `after`, `directly after`, `top`,
  `bottom`, or left empty (the default) to have no particular placement.  This
  param works in combination with the `position_reference` param.
* `position_reference` - Required if `position_keyword` is one of the
  "above" or "below" variants, this is the name of a non-group rule to use
  as a reference to place this group.
* `rule` - The rule definition (see below).  The rule ordering will match how
  they appear in the terraform plan file.

The following arguments are valid for each `rule` section:

* `name` - (Required) The rule name.
* `audit_comment` - When this rule is created/updated, the audit comment to
  apply for this rule.
* `description` - The description.
* `source_zones` - (Required) List of source zones.
* `source_addresses` - (Required) List of source addresses.
* `negate_source` - (bool) If the source should be negated.
* `source_users` - (Required) List of source users.
* `destination_zones` - (Required) List of destination zones.
* `destination_addresses` - (Required) List of destination addresses.
* `negate_destination` - (bool) Negate the destination addresses.
* `tags` - List of administrative tags.
* `disabled` - (bool) Disable this rule.
* `services` - (Required) List of services.
* `url_categories` - (Required) List of URL categories.
* `action` - Action to take.  Valid values are `no-decrypt` (default),
  `decrypt`, or `decrypt-and-forward`.
* `decryption_type` - The decryption type.  Valid values are `ssl-forward-proxy`,
  `ssh-proxy`, or `ssl-inbound-inspection`.
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


## Attributes

Each `rule` has the following attributes:

* `uuid` - The PAN-OS UUID.
