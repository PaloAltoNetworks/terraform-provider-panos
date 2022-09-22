---
page_title: "panos: panos_panorama_gke_cluster_group"
subcategory: "Plugins"
---

# panos_panorama_gke_cluster_group

!> **Note:** This is only valid for the 1.0 GCP Plugin.

This resource allows you to add/update/delete a GKE cluster group.

This resource requires that the GCP plugin be installed.


## PAN-OS

Panorama


## Import Name

```shell
<name>
```


## Example Usage

```hcl
resource "panos_panorama_gke_cluster_group" "grp" {
    name = "myCluster"
    gcp_project_credential = panos_panorama_gcp_account.gcp.name
    device_group = panos_panorama_device_group.dg.name
    template_stack = panos_panorama_template_stack.ts.name

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_device_group" "dg" {
    name = "my device group"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template_stack" "ts" {
    name = "myTemplateStack"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_gcp_account" "gcp" {
    name = "myGcpAccount"
    project_id = "gcp-project-123"
    service_account_credential_type = "gcp"
    credential_file = file("gcp-credentials.json")

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The cluster group's name.
* `description` - (Optional) The description.
* `gcp_project_credential` - (Required) The GCP account to use.
* `device_group` - (Required) The device group.
* `template_stack` - (Required) The template stack.
