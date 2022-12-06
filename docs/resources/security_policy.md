---
page_title: "panos: panos_security_policy"
subcategory: "Policies"
---

# panos_security_policy

This resource allows you to manage the full security posture.

This resource manages the full set of security rules in a vsys, enforcing both
the contents of individual rules as well as their ordering.  Rules are defined
in a `rule` config block.

!> **Warning**: This resource will remove any security rule not defined in this resource.

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


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_security_policies`
* `panos_panorama_security_policy`
* `panos_panorama_security_policies`


## Import Name

```shell
<device_group>:<rulebase>:<vsys>
```


## Example Usage

```hcl
resource "panos_security_policy" "example" {
    rule {
        name = "the opposite of secure"
        audit_comment = "Initial config"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

Panorama specific arguments:

* `device_group` - The device group (default: `shared`).
* `rulebase` - The rulebase.  This can be either
  `pre-rulebase` (default), `rulebase`, or `post-rulebase`.

NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).

The following arguments are supported:

* `rule` - A security rule definition (see below).  The security rule
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

Each `rule` section has the following attributes:

* `uuid` - The PAN-OS UUID.
