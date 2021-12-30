---
page_title: "panos: panos_aws_cloud_watch"
subcategory: "Plugins"
---

# panos_aws_cloud_watch

This resource allows you to manage the AWS CloudWatch plugin config.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_aws_cloud_watch" "example" {}
```


## Argument Reference

The following arguments are supported:

* `enabled` - Enable AWS CloudWatch setup (default: `true`).
* `namespace` - Namespace (default: `VMseries`).
* `update_interval` - (int) Update time (in min) (default: `5`).
