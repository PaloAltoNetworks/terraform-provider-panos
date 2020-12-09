---
page_title: "panos: panos_panorama_security_policy"
subcategory: "Panorama Policy"
---

# panos_panorama_security_policy

This resource allows you to manage the full security posture.

-> **Note:** `panos_panorama_security_policies` is known as `panos_panorama_security_policy`.

This resource manages the full set of security rules, enforcing both the
contents of individual rules as well as their ordering.  Rules are defined in
a `rule` config block.  As this manages the full set of security rules for
a given rulebase, any extraneous rules are removed on `terraform apply`.

!> **Note**: This resource will remove any security rule not defined in your plan file.

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


## Import Name

```
<device_group>:<rulebase>
```


## Example Usage

```hcl
resource "panos_panorama_security_policy" "example" {
    rule {
        name = "allow bizdev to dmz"
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
        name = "deny sales to eng"
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
            vsys_list = ["vsys1", "vsys3"]
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `device_group` - (Optional) The device group to put the security policy into
  (default: `shared`).
* `rulebase` - (Optional) The rulebase.  This can be `pre-rulebase` (default),
  `post-rulebase`, or `rulebase`.
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
* `target` - (Optional) A target definition (see below).  If there are no
  target sections, then the rule will apply to every vsys of every device
  in the device group.
* `negate_target` - (Optional, bool) Instead of applying the rule for the
  given serial numbers, apply it to everything except them.

The following arguments are valid for each `target` section:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.
