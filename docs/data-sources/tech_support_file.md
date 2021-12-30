---
page_title: "panos: panos_tech_support_file"
subcategory: "Operational State"
---

# panos_tech_support_file

Retrieve the tech support file from PAN-OS.

**WARNING**:  As this is a data source, even doing `terraform plan` will execute
retrieval of the tech support file.  It is highly recommended that this data source
not be included with the main plan file, but it's own plan file to avoid accidentally
invoking the tech support file export.


## Example Usage

```hcl
data "panos_tech_support_file" "example" {
    save_to_state = true
}

# Output the file to the screen, essentially exporting it from where Terraform
# is being executed at.
output "tech_support_file_content" {
    value = data.panos_tech_support_file.example.content
}
```


## Argument Reference

The following arguments are supported:

* `timeout` - (int) Timeout for retrieving the tech support file in seconds
  (default: `600`).
* `save_to_local_file_system` - (bool) Save the tech support file to the local
  file system where Terraform is running.
* `file_system_path` - When `save_to_local_file_system=true`, the file system path
  to place the tech support file.
* `save_to_state` - (bool) Save the tech support file to the state.


## Attribute Reference

* `filename` - The tech support file filename.
* `content` - For `save_to_state=true`, the content of the .tgz tech support file.
