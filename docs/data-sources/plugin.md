---
page_title: "panos: panos_plugin"
subcategory: "Operational State"
---

# panos_plugin

Retrieves information on plugins available on the PAN-OS NGFW or Panorama.

-> **Note:** Plugins for NGFW are present in PAN-OS 9.0+.

-> **Note:** `panos_panorama_plugin` is now `panos_plugin`.

## Example Usage

```hcl
data "panos_plugin" "example" {}
```

## Attribute Reference

The following attributes are present:

* `installed` - A list of installed plugins.
* `total` - (int) Total number of plugins, installed or not.
* `details` - A list of maps (see below).

The following attributes are present in each `details` entry:

* `name` - The name.
* `version` - The version number.
* `release_date` - Release date.
* `release_note_url` - Release note URL.
* `package_file` - The package file.
* `size` - The size.
* `platform` - Platform.
* `installed` - If the package is installed or not.
* `downloaded` - If the package is downloaded or not.
