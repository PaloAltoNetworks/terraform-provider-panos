---
page_title: "panos: panos_panorama_application_signature"
subcategory: "Objects"
---

# panos_panorama_application_signature

This resource allows you to add/update/delete Panorama application signatures.


## PAN-OS

Panorama


## Import Name

```
<device_group>:<application_object>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_application_signature" "example" {
    application_object = panos_panorama_application_object.myapp.name
    comment = "made by terraform"
    ordered_match = true
    and_condition {
        or_condition {
            pattern_match {
                context = "http-req-headers"
                pattern = "somepattern"
                qualifiers = {
                    "http-method": "COPY",
                    "req-hdr-type": "HOST",
                }
            }
        }
        or_condition {
            greater_than {
                // X.400-message size
                context = "cotp-req-x420-message-size"
                value = "123456"
            }
        }
        or_condition {
            less_than {
                // X.400-message size
                context = "cotp-req-x420-message-size"
                value = "42"
            }
        }
    }
    and_condition {
        or_condition {
            equal_to {
                context = "unknown-req-tcp"
                position = "first-4bytes"
                mask = "0Xff112345"
                value = "0X11bb33dd"
            }
        }
    }
}

resource "panos_panorama_application_object" "myapp" {
    name = "myApp"
    description = "made by terraform"
    category = "media"
    subcategory = "gaming"
    technology = "browser-based"
    defaults {
        port {
            ports = [
                "udp/dynamic",
            ]
        }
    }
    risk = 4
    scanning {
        viruses = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The signature's name.
* `device_group` - (Optional) The device group (default: `shared`)
* `application_object` - (Required) The applciation object for this signature.
* `comment` - (Optional) The description.
* `scope` - (Optional) The signature's scope.  Valid values are
  `transaction` (default) or `session`.
* `ordered_match` - (Optional, bool) Set to `false` to disable ordered matching
  (default: `true`).
* `and_condition` - (Optional) The and condition spec (defined below).

`and_condition` supports the following arguments:

* `name` - (Computed) And condition name, this is computed and cannot be configured.
* `or_condition` - (Required) The or condition spec (defined below).

`and_condition.or_condition` supports the following arguments:

* `name` - (Computed) Or condition name, this is computed and cannot be configured.
* `pattern_match` - (Optional) The pattern match spec (defined below).
* `greater_than` - (Optional) The greater than spec (defined below).
* `less_than` - (Optional) the less than spec (defined below).
* `equal_to` - (Optional) The equal to spec (defined below).

`and_condition.or_condition.pattern_match` supports the following arguments:

* `context` - (Required) The context.
* `pattern` - (Required) The pattern.
* `qualifiers` - (Optional, map) The qualifiers.

`and_condition.or_condition.greater_than` supports the following arguments:

* `context` - (Required) The context.
* `value` - (Required) The value.
* `qualifiers` - (Optional, map) The qualifiers.

`and_condition.or_condition.less_than` supports the following arguments:

* `context` - (Required) The context.
* `value` - (Required) The value.
* `qualifiers` - (Optional, map) The qualifiers.

`and_condition.or_condition.equal_to` supports the following arguments:

* `context` - (Required) The context.
* `value` - (Required) The value.
* `position` - (Optional) The position.
* `mask` - (Optional) The mask.
