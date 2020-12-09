---
page_title: "panos: panos_security_rule_group"
subcategory: "Firewall Policy"
---

# panos_security_rule_group

This resource allows you to add/update/delete security rule groups.

~> **Note:** `panos_security_policy_group` is known as `panos_security_rule_group`.

This resource manages clusters of security rules in a single vsys,
enforcing both the contents of individual rules as well as their
ordering.  Rules are defined in a `rule` config block.

Because this resource only manages what it's told to, it will not manage
any rules that may already exist on the firewall.  This has
implications on the effective security posture of your firewall, but it
will allow you to spread your security rules across multiple Terraform
state files.  If you want to verify that the security rules are only
what appears in the plan file, then you should probably be using the
[panos_security_policy](security_policy.html) resource.

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

## Example Usage

```hcl
resource "panos_security_rule_group" "example" {
    position_keyword = "above"
    position_reference = "deny everything else"
    rule {
        name = "allow bizdev to dmz"
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
        name = "deny sales to eng"
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

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys to put the security rule into (default:
  `vsys1`).
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
* `type` - (Optional) Rule type.  This can be `universal` (default),
  `interzone`, or `intrazone`.
* `description` - (Optional) The description.
* `tags` - (Optional) List of tags for this security rule.
* `source_zones` - (Required) List of source zones.
* `source_addresses` - (Required) List of source addresses.
* `negate_source` - (Optional, bool) If the source should be negated.
* `source_users` - (Required) List of source users.
* `hip_profiles` - (Required) List of HIP profiles.
* `destination_zones` - (Required) List of destination zones.
* `destination_addresses` - (Required) List of destination addresses.
* `negate_destination` - (Optional, bool) If the destination should be negated.
* `applications` - (Required) List of applications.
* `services` - (Required) List of services.
* `categories` - (Required) List of categories.
* `action` - (Optional) Action for the matched traffic.  This can be `allow`
  (default), `deny`, `drop`, `reset-client`, `reset-server`, or `reset-both`.
* `log_setting` - (Optional) Log forwarding profile.
* `log_start` - (Optional, bool) Log the start of the traffic flow.
* `log_end` - (Optional, bool) Log the end of the traffic flow (default: `true`).
* `disabled` - (Optional, bool) Set to `true` to disable this rule.
* `schedule` - (Optional) The security rule schedule.
* `icmp_unreachable` - (Optional) Set to `true` to enable ICMP unreachable.
* `disable_server_response_inspection` - (Optional) Set to `true` to disable
  server response inspection.
* `group` - (Optional) Profile Setting: `Group` - The group profile name.
* `virus` - (Optional) Profile Setting: `Profiles` - The antivirus setting.
* `spyware` - (Optional) Profile Setting: `Profiles` - The anti-spyware
  setting.
* `vulnerability` - (Optional) Profile Setting: `Profiles` - The Vulnerability
  Protection setting.
* `url_filtering` - (Optional) Profile Setting: `Profiles` - The URL filtering
  setting.
* `file_blocking` - (Optional) Profile Setting: `Profiles` - The file blocking
  setting.
* `wildfire_analysis` - (Optional) Profile Setting: `Profiles` - The WildFire
  Analysis setting.
* `data_filtering` - (Optional) Profile Setting: `Profiles` - The Data
  Filtering setting.
