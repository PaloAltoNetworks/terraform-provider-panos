---
page_title: "panos: panos_edl"
subcategory: "Objects"
---

# panos_edl

Retrieve information on the specified EDL.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_edl`


## Example Usage

```hcl
data "panos_edl" "example" {
    name = "example"
}
```

## Argument Reference

Panorama specific arguments:

* `device_group` - (Optional) The device group (default: `shared`).


NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The object's name


## Attribute Reference

The following attributes are supported:

* `type` - The type of EDL.
* `description` - The object's description.
* `source` - The EDL source URL
* `certificate_profile` - Profile for authenticating client certificates
* `username` - EDL username
* `password` - EDL password
* `repeat` - How often to retrieve the EDL.
* `repeat_at` - The time at which to retrieve the EDL.
* `repeat_day_of_week` - Repeat day of week.
* `repeat_day_of_month` - (int) If `repeat` is `monthly`, then the
  repeat day of month.
* `exceptions` - (list) List of exceptions.
