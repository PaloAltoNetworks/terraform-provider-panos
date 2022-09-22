---
page_title: "panos: panos_panorama_gcp_account"
subcategory: "Plugins"
---

# panos_panorama_gcp_account

!> **NOTE:**  This is only valid for the GCP 1.0 Plugin.

This resource allows you to add/update/delete GCP accounts on Panorama.

This resource requires that the GCP plugin be installed.


## PAN-OS

Panorama


## Import Name

```shell
<name>
```


## Example Usage

```hcl
# A GCP account type (for cluster groups)
resource "panos_panorama_gcp_account" "gcp" {
    name = "myGcpAccount"
    project_id = "gcp-project-123"
    service_account_credential_type = "gcp"
    credential_file = file("gcp-credentials.json")

    lifecycle {
        create_before_destroy = true
    }
}

# A GKE account type (for clusters in a group).
resource "panos_panorama_gcp_account" "gke" {
    name = "myGcpAccount"
    project_id = "gcp-project-123"
    service_account_credential_type = "gke"
    credential_file = file("gcp-credentials.json")

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The account's name.
* `description` - (Optional) Account description.
* `project_id` - (Required) The GCP project ID.
* `service_account_credential_type` - (Optional) The service account credential
  type.  Valid values are `gcp` (default) or `gke`.
* `credential_file` - (Required) The contents of a GCP credentials file; use the
  `file()` function to pass in the credentials file.
