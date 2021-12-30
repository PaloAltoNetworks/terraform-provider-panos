---
page_title: "panos: panos_security_rule_group"
subcategory: "Policies"
---

# panos_security_rule_group

This resource allows you to add/update/delete security rule groups.

This resource manages clusters of security rules in a single vsys,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Because this resource only manages what it's told to, it will not manage
any rules that may already exist on the firewall.  This has
implications on the effective security posture of your firewall, but it
will allow you to spread your security rules across multiple Terraform
state files.  If you want to verify that the security rules are only
what appears in the plan file, then you should probably be using the
[`panos_security_policy`](security_policy.html) resource.

Although you cannot modify non-group security rules with this
resource, the `position_keyword` and `position_reference` parameters allow you
to reference some other security rule that already exists, using it as
a means to ensure some rough placement within the ruleset as a whole.

For each security rule, there are three styles of profile settings:

* `None` (the default)
* `Group`
* `Profiles`

The Profile Setting is implicitly chosen based on what params are configured
for the security rule.  If you want a Profile Setting of `Group`, then the
`group` param should be set to the desired Group Profile.  If you want a
Profile Setting of `Profiles`, then you will need to specify one or more of
the following params:

* `virus`
* `spyware`
* `vulnerability`
* `url_filtering`
* `file_blocking`
* `wildfire_analysis`
* `data_filtering`

If the `group` param and none of the `Profiles` params are specified, then
the Profile Setting is set to `None`.


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


## Aliases

* `panos_security_policy_group`
* `panos_panorama_security_rule_group`
* `panos_panorama_security_policy_group`


## Example Usage

### NGFW Example

```hcl
resource "panos_security_rule_group" "example1" {
    position_keyword = "above"
    position_reference = panos_security_rule_group.example2.rule.0.name
    rule {
        name = "Allow bizdev to dmz"
        source_zones = [panos_zone.bizdev.name]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = [panos_zone.dmz.name]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
    rule {
        name = "Deny sales to eng"
        source_zones = [panos_zone.sales.name]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = [panos_zone.eng.name]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "deny"
    }
}

resource "panos_security_rule_group" "example2" {
    position_keyword = "bottom"
    rule {
        name = "Deny everything else"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "deny"
    }
}

resource "panos_zone" "bizdev" {
    name = "bizdev"
    mode = "layer3"
}

resource "panos_zone" "dmz" {
    name = "dmz"
    mode = "layer3"
}

resource "panos_zone" "sales" {
    name = "sales"
    mode = "layer3"
}

resource "panos_zone" "eng" {
    name = "eng"
    mode = "layer3"
}
```


### Panorama Example

```
resource "panos_security_rule_group" "example1" {
    position_keyword = "above"
    position_reference = panos_security_rule_group.example2.rule.0.name
    rule {
        name = "Allow bizdev to dmz"
        audit_comment = "Case id 12345"
        source_zones = ["bizdev"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["dmz"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
    rule {
        name = "Deny sales to eng"
        audit_comment = "Initial config"
        source_zones = ["sales"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["eng"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "deny"
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
}

resource "panos_security_rule_group" "example2" {
    position_keyword = "bottom"
    rule {
        name = "Deny everything else"
        audit_comment = "Initial config"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "deny"
    }
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

* `position_keyword` - (Optional) A positioning keyword for this group.  This
  can be `before`, `directly before`, `after`, `directly after`, `top`,
  `bottom`, or left empty (the default) to have no particular placement.  This
  param works in combination with the `position_reference` param.
* `position_reference` - (Optional) Required if `position_keyword` is one of the
  "above" or "below" variants, this is the name of a non-group rule to use
  as a reference to place this group.
* `rule` - The security rule definition (see below).  The security rule
  ordering will match how they appear in the terraform plan file.

The following arguments are valid for each `rule` section:

* `name` - (Required) The security rule name.
* `audit_comment` - When this rule is created/updated, the audit comment to
  apply for this rule.
* `group_tag` - (PAN-OS 9.0+) The group tag.
* `type` - Rule type.  This can be `universal` (default), `interzone`, or `intrazone`.
* `description` - The description.
* `tags` - List of tags for this security rule.
* `source_zones` - (Required) List of source zones.
* `source_addresses` - (Required) List of source addresses.
* `negate_source` - (bool) If the source should be negated.
* `source_users` - (Required) List of source users.
* `hip_profiles` - List of HIP profiles.
* `destination_zones` - (Required) List of destination zones.
* `destination_addresses` - (Required) List of destination addresses.
* `negate_destination` - (bool) If the destination should be negated.
* `applications` - (Required) List of applications.
* `services` - (Required) List of services.
* `categories` - (Required) List of categories.
* `action` - Action for the matched traffic.  This can be `allow`
  (default), `deny`, `drop`, `reset-client`, `reset-server`, or `reset-both`.
* `log_setting` - Log forwarding profile.
* `log_start` - (bool) Log the start of the traffic flow.
* `log_end` - (bool) Log the end of the traffic flow (default: `true`).
* `disabled` - (bool) Set to `true` to disable this rule.
* `schedule` - The security rule schedule.
* `icmp_unreachable` - (bool) Set to `true` to enable ICMP unreachable.
* `disable_server_response_inspection` - (bool) Set to `true` to disable
  server response inspection.
* `group` - Profile Setting: `Group` - The group profile name.
* `virus` - Profile Setting: `Profiles` - The antivirus setting.
* `spyware` - Profile Setting: `Profiles` - The anti-spyware setting.
* `vulnerability` - Profile Setting: `Profiles` - The Vulnerability Protection setting.
* `url_filtering` - Profile Setting: `Profiles` - The URL filtering setting.
* `file_blocking` - Profile Setting: `Profiles` - The file blocking setting.
* `wildfire_analysis` - Profile Setting: `Profiles` - The WildFire Analysis setting.
* `data_filtering` - Profile Setting: `Profiles` - The Data Filtering setting.
* `target` - (repeatable, Panorama only) A target definition (see below).  If there
  are no target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (bool, Panorama only) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.

`rule.target` supports the following arguments:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - A listing of vsys to apply this rule to.  If `serial` is
  a VM, then this parameter should just be omitted.


## Attribute Reference

Each `rule` has the following attribute:

* `uuid` - The PAN-OS UUID.
