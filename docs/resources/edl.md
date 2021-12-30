---
page_title: "panos: panos_edl"
subcategory: "Objects"
---

# panos_edl

This resource allows you to add/update/delete external dynamic lists (EDL).


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_edl`


## Setting `repeat_at`

The acceptable PAN-OS values for the `repeat_at` field is a combination of
the version of PAN-OS that you're running against and the setting of the `repeat`
parameter.

The following shorthand is used:

* `N/A` - `repeat_at` should not be set
* `minute` - A two character minute string (e.g. - `07` or `59`)
* `24hr hour` - A two character hour string in 24hr notation (e.g. - `09` or `15`)
* `24hr time` - A five character hour/minute string in 24hr notation (e.g. - `09:00` or `23:59`)

Here are the valid settings for `repeat_at` given your desired `repeat` value
and the version of PAN-OS you're running against:

* PAN-OS 6.1 - 7.0
  * `hourly` - minute
  * `daily`, `weekly`, `monthly` - 24hr time
* PAN-OS 7.1+
  * `every five minutes`, `hourly` - N/A
  * `daily`, `weekly`, `monthly` - 24hr hour


## Import Name

```
<device_group>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_edl" "example" {
    name = "example"
    type = "ip"
    description = "my edl"
    source = "https://example.com"
    repeat = "every five minutes"
    exceptions = ["10.1.1.1", "10.1.1.2"]
}
```

## Argument Reference

Panorama specific arguments:

* `device_group` - The device group (default: `shared`).


NGFW specific arguments:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The object's name
* `type` - The type of EDL.  This can be `ip` (the default; and the
  only valid value for PAN-OS 6.1 - 7.0), `domain`, `url`, `predefined-ip`
  (PAN-OS 8.0+), or `predefined-url` (PAN-OS 9.0+).
* `description` - The object's description.
* `source` - The EDL source URL
* `certificate_profile` - Profile for authenticating client certificates
* `username` - EDL username
* `password` - EDL password
* `repeat` - How often to retrieve the EDL.  This can be `hourly` (the
  default), `daily`, `weekly`, `monthly`, or `every five minutes` (valid for
  PAN-OS 7.1+)
* `repeat_at` - The time at which to retrieve the EDL.  Please refer
  to the section above for how to set this value properly.
* `repeat_day_of_week` - If `repeat` is `weekly`, then this should
  be set to the desired day of the week.  Valid values are `sunday`,
  `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, and
  `sunday`
* `repeat_day_of_month` - (int) If `repeat` is `monthly`, then this should
  be set to the desired day of the month.
* `exceptions` - (list) Provide a list of exception entries.
