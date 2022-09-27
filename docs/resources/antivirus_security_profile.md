---
page_title: "panos: panos_antivirus_security_profile"
subcategory: "Objects"
---

# panos_antivirus_security_profile

Manages anti-virus security profiles.

## Import Name

NGFW:

```shell
<vsys>:<name>
```

Panorama:

```shell
<device_group>:<name>
```


## Example Usage

```hcl
# NOTE: older PAN-OS does not support "http2".
resource "panos_antivirus_security_profile" "example" {
    name = "example"
    description = "my description"
    decoder { name = "smtp" }
    decoder { name = "smb" }
    decoder { name = "pop3" }
    decoder { name = "imap" }
    decoder { name = "http" }
    decoder { name = "http2" }
    decoder {
        name = "ftp"
        action = "reset-both"
    }
    application_exception {
        application = "hotmail"
        action = "alert"
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `packet_capture` - (bool) Set to `true` to enable packet capture.
* `threat_exceptions` - (list) List of threat exceptions.
* `decoder` - (Repeatable) A decoder spec, as defined below.
* `application_exception` - (Repeatable) An application exception spec, as
  defined below.
* `machine_learning_model` - (Repeatable) A machine learning model spec, as
  defined below.
* `machine_learning_exception` - (Repeatable) A machine learning exception spec, as
  defined below.

`decoder` supports the following arguments:

* `name` - (Required) Decoder name.
* `action` - Decoder action.  Valid values are `default` (default), `allow`,
  `alert`, `block` (PAN-OS 6.1 only), `drop` (PAN-OS 7.0+), `reset-client` (PAN-OS
  7.0+), `reset-server` (PAN-OS 7.0+), or `reset-both` (PAN-OS 7.0+).
* `wildfire_action` - Wildfire action.  Valid values are `default` (default), `allow`,
  `alert`, `block` (PAN-OS 6.1 only), `drop` (PAN-OS 7.0+), `reset-client` (PAN-OS
  7.0+), `reset-server` (PAN-OS 7.0+), or `reset-both` (PAN-OS 7.0+).
* `machine_learning_action` - (PAN-OS 10.0+) Machine learning action.

`application_exception` supports the following arguments:

* `application` - (Required) The application name
* `action` - The action.  Valid values are `default`, `allow`,
  `alert`, `block` (PAN-OS 6.1 only), `drop` (PAN-OS 7.0+), `reset-client` (PAN-OS
  7.0+), `reset-server` (PAN-OS 7.0+), or `reset-both` (PAN-OS 7.0+).

`machine_learning_model` supports the following arguments:

* `model` - (Required) The model.
* `action` - (Required) The action.

`machine_learning_exception` supports the following arguments:

* `name` - (Required) The name.
* `description` - The description.
* `filename` - Filename
