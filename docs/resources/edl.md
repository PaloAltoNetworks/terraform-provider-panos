---
page_title: "panos: panos_edl"
subcategory: "Firewall Objects"
---

# panos_edl

This resource allows you to add/update/delete external dynamic lists (EDL).


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
<vsys>:<name>
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

The following arguments are supported:

* `name` - (Required) The object's name
* `vsys` - (Optional) The vsys to put the object into (default: `vsys1`)
* `type` - (Optional) The type of EDL.  This can be `ip` (the default; and the
  only valid value for PAN-OS 6.1 - 7.0), `domain`, `url`, or `predefined`
  (PAN-OS 8.0+)
* `description` - (Optional) The object's description.
* `source` - (Optional) The EDL source URL
* `certificate_profile` - (Optional) Profile for authenticating client certificates
* `username` - (Optional) EDL username
* `password` - (Optional) EDL password
* `repeat` - (Optional) How often to retrieve the EDL.  This can be `hourly` (the
  default), `daily`, `weekly`, `monthly`, or `every five minutes` (valid for
  PAN-OS 7.1+)
* `repeat_at` - (Optional) The time at which to retrieve the EDL.  Please refer
  to the section above for how to set this value properly.
* `repeat_day_of_week` - (Optional) If `repeat` is `weekly`, then this should
  be set to the desired day of the week.  Valid values are `sunday`,
  `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, and
  `sunday`
* `repeat_day_of_month` - (Optional, int) If `repeat` is `monthly`, then this should
  be set to the desired day of the month.
* `exceptions` - (Optional, list) Provide a list of exception entries.
