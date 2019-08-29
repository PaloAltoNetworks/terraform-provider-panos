---
layout: "panos"
page_title: "panos: panos_panorama_plugin"
description: |-
  Gets Panorama plugin info.
---

# panos_plugin_info

Use this data source to retrieve "show system info" from the NGFW or Panorama.

All contents of "show system info" are saved to the `info` variable.  In
addition, the version number of PAN-OS encountered is saved to multiple
fields for ease of access.

## Example Usage

```hcl
data "panos_panorama_plugin" "example" {}
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
