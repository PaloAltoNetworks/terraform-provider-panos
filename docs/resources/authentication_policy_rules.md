---
page_title: "panos_authentication_policy_rules Resource - panos"
subcategory: "Policies"
description: |-
  
---

# panos_authentication_policy_rules (Resource)



## Example Usage

```terraform
# Manage a group of authentication policy rules with positioning

## Place the rule group at the top of the pre-rulebase
resource "panos_authentication_policy_rules" "guest_network" {
  location = {
    device_group = {
      name     = panos_device_group.example.name
      rulebase = "pre-rulebase"
    }
  }

  position = {
    where = "first"
  }

  rules = [
    {
      name                       = "guest-wifi-auth"
      description                = "Require authentication for guest WiFi users"
      source_zones               = ["guest-zone"]
      source_addresses           = ["guest-network"]
      destination_zones          = ["untrust"]
      destination_addresses      = ["any"]
      services                   = ["any"]
      authentication_enforcement = "guest-captive-portal"
      timeout                    = 480
      log_authentication_timeout = true
      log_setting                = "authentication-log-profile"
    }
  ]
}

## Place the rule group after a specific rule
resource "panos_authentication_policy_rules" "corporate_users" {
  location = {
    device_group = {
      name     = panos_device_group.example.name
      rulebase = "pre-rulebase"
    }
  }

  position = {
    where    = "after"
    directly = true
    pivot    = "guest-wifi-auth"
  }

  rules = [
    {
      name                       = "employee-byod-auth"
      description                = "Authentication for employee BYOD devices"
      source_zones               = ["byod-zone"]
      source_addresses           = ["byod-subnet"]
      source_users               = ["any"]
      destination_zones          = ["internal", "dmz"]
      destination_addresses      = ["corporate-apps"]
      services                   = ["any"]
      category                   = ["business-and-economy", "computer-and-internet-info"]
      authentication_enforcement = "corporate-auth-profile"
      timeout                    = 1440
      log_authentication_timeout = false
      tags                       = ["byod", "corporate"]
    },
    {
      name                       = "contractor-limited-access"
      description                = "Authentication for contractors with restricted access"
      source_zones               = ["contractor-zone"]
      source_addresses           = ["contractor-subnet"]
      source_users               = ["contractor-group"]
      destination_zones          = ["dmz"]
      destination_addresses      = ["contractor-apps"]
      services                   = ["service-https"]
      authentication_enforcement = "contractor-auth-profile"
      timeout                    = 240
      log_authentication_timeout = true
      log_setting                = "authentication-log-profile"
      tags                       = ["contractor", "restricted"]
    }
  ]
}

## Advanced rule with HIP checks and target restrictions
resource "panos_authentication_policy_rules" "hip_based_auth" {
  location = {
    device_group = {
      name     = panos_device_group.example.name
      rulebase = "post-rulebase"
    }
  }

  position = {
    where = "last"
  }

  rules = [
    {
      name                       = "hip-compliant-devices"
      description                = "Allow authenticated access only for HIP-compliant devices"
      source_zones               = ["trust"]
      source_addresses           = ["corporate-subnets"]
      source_hip                 = ["compliant-hip-profile"]
      destination_zones          = ["dmz", "internal"]
      destination_addresses      = ["sensitive-servers"]
      destination_hip            = ["any"]
      services                   = ["any"]
      source_users               = ["domain\\authenticated-users"]
      authentication_enforcement = "mfa-auth-profile"
      timeout                    = 720
      log_authentication_timeout = true
      log_setting                = "security-log-profile"

      # Target specific devices in the device group
      target = {
        devices = [
          {
            name = "fw-datacenter-01"
            vsys = [
              { name = "vsys1" }
            ]
          },
          {
            name = "fw-datacenter-02"
            vsys = [
              { name = "vsys1" },
              { name = "vsys2" }
            ]
          }
        ]
        negate = false
        tags   = ["production"]
      }

      tags = ["hip-required", "production", "authenticated"]
    },
    {
      name                       = "non-compliant-redirect"
      description                = "Redirect non-compliant devices to remediation portal"
      source_zones               = ["trust"]
      source_addresses           = ["corporate-subnets"]
      negate_source              = false
      destination_zones          = ["remediation"]
      destination_addresses      = ["remediation-portal"]
      negate_destination         = false
      services                   = ["service-http", "service-https"]
      authentication_enforcement = "remediation-auth-profile"
      timeout                    = 60
      log_authentication_timeout = true
      disabled                   = false
      tags                       = ["remediation"]
    }
  ]
}

resource "panos_device_group" "example" {
  location = {
    panorama = {}
  }

  name = "example-device-group"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `location` (Attributes) The location of this object. (see [below for nested schema](#nestedatt--location))
- `position` (Attributes) (see [below for nested schema](#nestedatt--position))
- `rules` (Attributes List) (see [below for nested schema](#nestedatt--rules))

<a id="nestedatt--location"></a>
### Nested Schema for `location`

Optional:

- `device_group` (Attributes) Located in a specific device group rulebase (see [below for nested schema](#nestedatt--location--device_group))
- `shared` (Attributes) Located in a shared rulebase (see [below for nested schema](#nestedatt--location--shared))
- `vsys` (Attributes) Located in a specific vsys rulebase (see [below for nested schema](#nestedatt--location--vsys))

<a id="nestedatt--location--device_group"></a>
### Nested Schema for `location.device_group`

Optional:

- `name` (String) The device group name
- `panorama_device` (String) The panorama device
- `rulebase` (String) The rulebase


<a id="nestedatt--location--shared"></a>
### Nested Schema for `location.shared`

Optional:

- `rulebase` (String) Rulebase name


<a id="nestedatt--location--vsys"></a>
### Nested Schema for `location.vsys`

Optional:

- `name` (String) The vsys name
- `ngfw_device` (String) The NGFW device



<a id="nestedatt--position"></a>
### Nested Schema for `position`

Required:

- `where` (String)

Optional:

- `directly` (Boolean)
- `pivot` (String)


<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `name` (String)

Optional:

- `audit_comment_version` (String) Version trigger for audit comments. Change this value to send the audit_comment_wo to PAN-OS. This attribute is not sent to PAN-OS itself, but serves as a trigger to detect when the audit comment should be updated.
- `audit_comment_wo` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Write-only audit comment for this rule. This value is sent to PAN-OS but not read back. Changes are only sent when audit_comment_version is modified. Each time audit_comment_version changes, this comment is added to the audit history with a timestamp.
- `authentication_enforcement` (String) Authentication enforcement object to use for authentication.
- `category` (Set of String)
- `description` (String)
- `destination_addresses` (Set of String)
- `destination_hip` (Set of String)
- `destination_zones` (Set of String)
- `disabled` (Boolean) Disable the rule
- `group_tag` (String)
- `log_authentication_timeout` (Boolean)
- `log_setting` (String) Log setting for forwarding authentication logs
- `negate_destination` (Boolean)
- `negate_source` (Boolean)
- `services` (Set of String)
- `source_addresses` (Set of String)
- `source_hip` (Set of String)
- `source_users` (Set of String)
- `source_zones` (Set of String)
- `tags` (Set of String)
- `target` (Attributes) (see [below for nested schema](#nestedatt--rules--target))
- `timeout` (Number) expiration timer (minutes)

<a id="nestedatt--rules--target"></a>
### Nested Schema for `rules.target`

Optional:

- `devices` (Attributes List) (see [below for nested schema](#nestedatt--rules--target--devices))
- `negate` (Boolean) Target to all but these specified devices and tags
- `tags` (List of String)

<a id="nestedatt--rules--target--devices"></a>
### Nested Schema for `rules.target.devices`

Required:

- `name` (String)

Optional:

- `vsys` (Attributes List) (see [below for nested schema](#nestedatt--rules--target--devices--vsys))

<a id="nestedatt--rules--target--devices--vsys"></a>
### Nested Schema for `rules.target.devices.vsys`

Required:

- `name` (String)

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# A set of authentication policy rules can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#     position = { where = "after", directly = true, pivot = "guest-wifi-auth" }
#
#     names = [
#         "employee-byod-auth",
#         "contractor-limited-access"
#     ]
# }
terraform import panos_authentication_policy_rules.corporate_users $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["employee-byod-auth","contractor-limited-access"],"position":{"directly":true,"pivot":"guest-wifi-auth","where":"after"}}' | base64)
```
