---
page_title: "panos: panos_antivirus_security_profile"
subcategory: "Objects"
---

# panos_antivirus_security_profile

Gets anti-virus security profile info.


## Example Usage

```hcl
data "panos_antivirus_security_profile" "example" {
    name = panos_antivirus_security_profile.x.name
}

resource "panos_antivirus_security_profile" "x" {
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
}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys to put the address object into (default:
  `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

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

`decoder` supports the following attributes:

* `name` - Decoder name.
* `action` - Decoder action.
* `wildfire_action` - Wildfire action.
* `machine_learning_action` - (PAN-OS 10.0+) Machine learning action.

`application_exception` supports the following attributes:

* `application` - The application name
* `action` - The action.

`machine_learning_model` supports the following attributes:

* `model` - The model.
* `action` - The action.

`machine_learning_exception` supports the following attributes:

* `name` - The name.
* `description` - The description.
* `filename` - Filename
