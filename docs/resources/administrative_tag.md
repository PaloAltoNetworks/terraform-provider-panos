---
page_title: "panos: panos_administrative_tag"
subcategory: "Objects"
---

# panos_administrative_tag

This resource allows you to add/update/delete administrative tags.

Tag colors are as follows:

* `color1`: Red
* `color2`: Green
* `color3`: Blue
* `color4`: Yellow
* `color5`: Copper
* `color6`: Orange
* `color7`: Purple
* `color8`: Gray
* `color9`: Light Green
* `color10`: Cyan
* `color11`: Light Gray
* `color12`: Blue Gray
* `color13`: Lime
* `color14`: Black
* `color15`: Gold
* `color16`: Brown
* `color17`: Olive
* `color18`: (Reserved for internal use)
* `color19`: Maroon
* `color20`: Red Orange
* `color21`: Yellow Orange
* `color22`: Forest Green
* `color23`: Turquoise Blue
* `color24`: Azure Blue
* `color25`: Cerulean Blue
* `color26`: Midnight Blue
* `color27`: Medium Blue
* `color28`: Cobalt Blue
* `color29`: Violet Blue
* `color30`: Blue Violet
* `color31`: Medium Violet
* `color32`: Medium Rose
* `color33`: Lavender
* `color34`: Orchid
* `color35`: Thistle
* `color36`: Peach
* `color37`: Salmon
* `color38`: Magenta
* `color39`: Red Violet
* `color40`: Mahogany
* `color41`: Burnt Sienna
* `color42`: Chestnut


## PAN-OS

NGFW


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_administrative_tag" "example" {
    name = "tag1"
    vsys = "vsys2"
    color = "color5"
    comment = "Internal resources"
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys to put the administrative tag into (default: `vsys1`).
* `name` - (Required) The administrative tag's name.
* `color` - (Optional) The tag's color.  This should be either an empty string
  (no color) or a string such as `color1` or `color15`.  Note that for maximum
  portability, you should limit color usage to `color16`, which was available
  in PAN-OS 6.1.  PAN-OS 8.1's colors go up to `color42`.  The value `color18`
  is reserved internally by PAN-OS and thus not available for use.
* `comment` - (Optional) The administrative tag's description.
