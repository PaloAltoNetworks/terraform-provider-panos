---
page_title: "panos: panos_dag_tags"
subcategory: "User-ID"
---

# panos_dag_tags

This resource allows you to add and remove dynamic address group tags.

The `ip` field should be unique in the `panos_dag_tags` block, and there
should only be one `panos_dag_tags` block defined in a given plan.

**Note** - Tags are only removed during `terraform destroy`.  Updating an
applied terraform plan to have alternative tags will leave behind the
old tags from the previously published plan(s).


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_dag_tags" "example" {
    vsys = "vsys1"
    register {
        ip = "10.1.1.1"
        tags = ["tag1", "tag2"]
    }
    register {
        ip = "10.1.1.2"
        tags = ["tag3"]
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys to put the DAG tags in (default: `vsys1`).
* `register` - (Required) A set that includes `ip`, the IP address to be tagged
  and `tags`, a list of tags to associate with the given IP.
