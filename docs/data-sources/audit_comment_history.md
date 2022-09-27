---
page_title: "panos: panos_audit_comment_history"
subcategory: "Policies"
---

# panos_audit_comment_history

Returns audit comment history about the given rule.


## PAN-OS

NGFW and Panorama.


## Example Usage

### NGFW Example

```hcl
data "panos_audit_comment_history" "example" {
    rule_type = "security"
    name = panos_security_rule_group.x.rule.0.name
}

resource "panos_security_rule_group" "x" {
    rule {
        name = "my rule"
        description = "Made by Terraform"
        ...
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

### Panorama Example

```hcl
data "panos_audit_comment_history" "example" {
    device_group = panos_panorama_device_group.x.name
    rule_type = "security"
    rulebase = panos_security_rule_group.x.rulebase
    name = panos_security_rule_group.x.rule.0.name
}

resource "panos_panorama_device_group" "x" {
    name = "my device group"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_security_rule_group" "x" {
    device_group = panos_panorama_device_group.x.name
    rulebase = "pre-rulebase"
    rule {
        name = "Allow eng to DMZ"
        description = "Made by Terraform"
        ...
    }

    lifecycle {
        create_before_destroy = true
    }
}
```
    

## Argument Reference

NGFW:

* `vsys` - The vsys. (default: `vsys1`).

Panorama:

* `device_group` - The device group location (default: `shared`)
* `rulebase` - The rulebase.  Valid values are `pre-rulebase` (default),
  `rulebase`, or `post-rulebase`.

The following arguments are supported:

* `rule_type` - The rule type.  Valid values are `security` (default),
  `pbf`, `nat`, and `decryption`.
* `name` - (Required) The rule's name.
* `direction` - Valid values are `backward` (default) to see newest logs first, or
  `forward` to see oldest first.
* `nlogs` - (int) Number of audit comments to return, max 5000 (default: `100`).
* `skip` - (int) Specify the number of logs to skip when doing log retrieval.  This
  is useful when retrieving logs in batches to skip previously retrieved logs.


## Attribute Reference

The following attributes are supported:

* `log` - (repeated) A log entry spec, defined below.

Each `log` section has the following attributes:

* `admin` - The admin who made the change.
* `comment` - The audit comment.
* `config_version` - The config version.
* `time_generated` - The time generated, as reported by PAN-OS.
* `time_generated_rfc3339` - An opportunistic representation of `time_generated`
  in RFC3339.  This is created by combining the `time_generated` with the timezone
  information of PAN-OS.
