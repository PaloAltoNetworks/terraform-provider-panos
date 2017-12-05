---
layout: "panos"
page_title: "PANOS: panos_security_policy"
sidebar_current: "docs-panos-resource-security-policy"
description: |-
  Manages security policies.
---

# panos_security_policy

This resource allows you to add/update/delete security policies.

There are three styles of profile settings:

* `None` (the default)
* `Group`
* `Profiles`

The Profile Setting is implicitly chosen based on what params are configured
for the security policy.  If you want a Profile Setting of `Group`, then the
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

## Example Usage

```hcl
resource "panos_security_policy" "example" {
    name = "allow bizdev to dmz"
    source_zone = ["bizdev"]
    destination_zone = ["dmz"]
    action = "allow"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The security policy's name.
* `vsys` - (Optional) The vsys to put the security policy into (default:
  `vsys1`).
* `rulebase` - (Optional) The rulebase.  For firewalls, there is only the
  `rulebase` value (default), but on Panorama, there is also `pre-rulebase`
  and `post-rulebase`.
* `type` - (Optional) Rule type.  This can be `universal` (default),
  `interzone`, or `intrazone`.
* `description` - (Optional) The description.
* `tags` - (Optional) List of tags for this security rule.
* `source_zone` - (Optional) List of source zones (default: `["any"]`).
* `source_address` - (Optional) List of source addresses (default: `["any"]`).
* `negate_source` - (Optional) If the source should be negated.
* `source_user` - (Optional) List of source users (default: `["any"]`).
* `hip_profile` - (Optional) List of HIP profiles.
* `destination_zone` - (Optional) List of destination zones (default: `["any"]`).
* `destination_address` - (Optional) List of destination addresses (default:
  `["any"]`).
* `negate_destination` - (Optional) If the destination should be negated.
* `application` - (Optional) List of applications (default: `["any"]`).
* `service` - (Optional) List of services (default: `["application-default"]`).
* `category` - (Optional) List of categories (default: `["any"]`).
* `action` - (Optional) Action for the matched traffic.  This can be `allow`
  (default), `deny`, `drop`, `reset-client`, `reset-server`, or `reset-both`.
* `log_setting` - (Optional) Log forwarding profile.
* `log_start` - (Optional) Log the start of the traffic flow.
* `log_end` - (Optional) Log the end of the traffic flow (default: `true`).
* `disabled` - (Optional) Set to `true` to disable this rule.
* `schedule` - (Optional) The security policy schedule.
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
