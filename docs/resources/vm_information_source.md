---
page_title: "panos: panos_vm_information_source"
subcategory: "Device"
---

# panos_vm_information_source

This resource allows you to add/update/delete VM information sources.


## PAN-OS

NGFW and Panorama


## Import Names

```
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
# AWS VPC example.
resource "panos_vm_information_source" "ex1" {
    name = "example1"
    aws_vpc {
        description = "Made by Terraform"
        source = "src"
        access_key_id = "abcd1234abcd1234"
        secret_access_key = "abcd1234abcd12345"
        vpc_id = "myVpcId"
    }
}
```

```hcl
# VCenter example.
resource "panos_vm_information_source" "ex2" {
    name = "example2"
    vcenter {
        description = "Made by Terraform"
        port = 8443
        enable_timeout = true
        timeout = 5
        source = "mySrc"
        username = "user"
        password = "pass"
        update_interval = 6
    }
}
```

```hcl
# ESXi example.
resource "panos_vm_information_source" "ex3" {
    name = "example3"
    esxi {
        description = "Made by Terraform"
        port = 8443
        enable_timeout = true
        timeout = 5
        source = "mySrc"
        username = "user"
        password = "pass"
        update_interval = 6
    }
}
```

```hcl
# Google Compute Engine example.
resource "panos_vm_information_source" "ex4" {
    name = "example4"
    google_compute {
        description = "Made by Terraform"
        auth_type = "service-account"
        service_account_credential = file("creds.json")
        project_id = "myProj"
        zone_name = "someZone"
        update_interval = 120
        enable_timeout = true
        timeout = 3
    }
}
```


## Argument Reference

Panorama specific arguments (one of these must be specified):

* `template` - The template name.
* `template_stack` - The template stack name.

NGFW / Panorama:

* `vsys` - The vsys (default: `vsys1`).

The following arguments are supported:

* `name` - (Required) The zone's name.
* `aws_vpc` - AWS VPC information source spec (see below).
* `esxi` - VMware ESXi information source spec (see below).
* `vcenter` - VMware vCenter information source spec (see below).
* `google_compute` - Google compute engine information source spec (see below).

`aws_vpc` supports the following arguments:

* `description` - The description.
* `disabled` - (bool) Disabled or not.
* `source` - (Required) IP address or name.
* `access_key_id` - (Required) AWS access key ID.
* `secret_access_key` - (Required) AWS secret access key.
* `update_interval` - (int) Time interval (in sec) for updates (default: `60`).
* `enable_timeout` - (bool) Enable vm-info timeout when source is disconnected.
* `timeout` - (int) The vm-info timeout value (in hours) when source is disconnected (default: `2`).
* `vpc_id` - (Required) AWS VPC name or ID.

`esxi` supports the following arguments:

* `description` - The description.
* `port` - (int) The port number (default: `443`).
* `disabled` - (bool) Disabled or not.
* `enable_timeout` - (bool) Enable vm-info timeout when source is disconnected.
* `timeout` - (int) The vm-info timeout value (in hours) when source is disconnected (default: `2`).
* `source` - (Required) IP address or source name for vm-info-source.
* `username` - (Required) The vm-info-source login username.
* `password` - (Required) The vm-info-source login password.
* `update_interval` - (int) Time interval (in sec) for updates (default: `5`).

`vcenter` supports the following arguments:

* `description` - The description.
* `port` - (int) The port number (default: `443`).
* `disabled` - (bool) Disabled or not.
* `enable_timeout` - (bool) Enable vm-info timeout when source is disconnected.
* `timeout` - (int) The vm-info timeout value (in hours) when source is disconnected (default: `2`).
* `source` - (Required) IP address or source name for vm-info-source.
* `username` - (Required) The vm-info-source login username.
* `password` - (Required) The vm-info-source login password.
* `update_interval` - (int) Time interval (in sec) for updates (default: `5`).

`google_compute` supports the following arguments:

* `description` - The description.
* `disabled` - (bool) Disabled or not.
* `auth_type` - The auth type.  Valid values are `service-in-gce` (default) or
  `service-account`.
* `service_account_credential` - GCE service account JSON file.
* `project_id` - (Required) Google Compute Engine Project-ID.
* `zone_name` - (Required) Google Compute Engine project zone name.
* `update_interval` - (int) Time interval (in sec) for updates (default: `5`).
* `enable_timeout` - (bool) Enable vm-info timeout when source is disconnected.
* `timeout` - (int) The vm-info timeout value (in hours) when source is disconnected (default: `2`).
