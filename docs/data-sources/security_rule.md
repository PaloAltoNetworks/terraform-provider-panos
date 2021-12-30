---
page_title: "panos: panos_security_rule"
subcategory: "Policies"
---

# panos_security_rule

Retrieve information about the given security rule.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
data "panos_security_rule" "example" {
    name = "my rule name"
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

* `rule` - The security rule definition (see below).

`rule` has the following attributes:

* `name` - The security rule name.
* `uuid` - The PAN-OS UUID.
* `group_tag` - The group tag.
* `type` - Rule type.
* `description` - The description.
* `tags` - List of tags for this security rule.
* `source_zones` - List of source zones.
* `source_addresses` - List of source addresses.
* `negate_source` - (bool) If the source should be negated.
* `source_users` - List of source users.
* `hip_profiles` - List of HIP profiles.
* `destination_zones` - List of destination zones.
* `destination_addresses` - List of destination addresses.
* `negate_destination` - (bool) If the destination should be negated.
* `applications` - List of applications.
* `services` - List of services.
* `categories` - List of categories.
* `action` - Action for the matched traffic.
* `log_setting` - Log forwarding profile.
* `log_start` - (bool) Log the start of the traffic flow.
* `log_end` - (bool) Log the end of the traffic flow.
* `disabled` - (bool) Set to `true` to disable this rule.
* `schedule` - The security rule schedule.
* `icmp_unreachable` - (bool) ICMP unreachable.
* `disable_server_response_inspection` - (bool) If server response inspection is disabled.
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

`rule.target` supports the following attributes:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - A listing of vsys to apply this rule to.  If `serial` is
  a VM, then this parameter should just be omitted.
