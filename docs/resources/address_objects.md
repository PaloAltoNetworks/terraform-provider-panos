---
page_title: "panos: panos_address_objects"
subcategory: "Objects"
---

# panos_address_objects

This resource allows you to add/update/delete address objects in bulk.

If you only need a few address objects, then using
[`panos_address_object`](address_object.html) may be more appropriate.

Because the operations happen in bulk, this resource can optimize the API calls
to PAN-OS better than numerous individual `panos_address_object` resource
definitions can.  However, there are a number of additional considerations that
come with this resource, which are detailed below.


## Provider Versioning

v1.10.0+


## PAN-OS

NGFW and Panorama.


## Import Name

This resource cannot be imported.


## Provider Timeout

By default the provider's timeout is 10 seconds.  If you intend to use this
resource with thousands of address objects, the `timeout` param inside your
`panos` provider block should be increased to at least 30 seconds, if not more.


## Terraform Performance

It can be quite easy to overwhelm the host executing the `terraform` commands
when using this resource.  Besides handling the resource spec itself, the
definition that you define for any given instance of a `panos_address_objects`
resource cannot exceed 4MB, otherwise Terraform will start outputting
[grpc errors](https://github.com/hashicorp/terraform-provider-local/issues/28).

This resource does not impose limits on the number of objects contained in a single
definition, but it is highly recommended that you limit the number objects in
a single `panos_address_objects` definition to under 10k to avoid any issues.

In practice this number may need to be smaller based on the content of each individual
object (length of the name, length of the value, etc).  Experiment and find the
right number of objects, balancing the number of objects in a single resource with
Terraform execution times.


## Panorama GUI Cacheing

When using this resource against Panorama, you will notice that some of the
objects created by this resource don't show up in dropdowns where address objects
can be specified, such as source/destination addresses in security rules or NAT
rules.  This is just a display issue in the PAN-OS GUI due to a combination of
bulk object creation and notifiers on updating the GUI cache.

This is just a GUI display issue and has no effect on functionality.

If you need objects added via this resource to show up in the GUI drowdowns, then
restarting the management plane (or rebooting PAN-OS) will refresh the cache.


## Managing Dependencies When Using This Resource

The examples below show how to dynamically find objects within the
`panos_address_objects` resource and use them based on different naming conventions.

In practice, everytime your plan uses another builtin that loops over this resource
(especially if you define thousands of object definitions), you are increasing the
workload on Terraform / the host running Terraform, which will increase plan
execution time.

It may likely end up being best to use
[`depends_on`](https://www.terraform.io/language/meta-arguments/depends_on)
to link dependencies between this resource and others, which would allow you to
just put the name of the various address objects instead of using a
[resource attribute variable](https://www.terraform.io/language/expressions/references)
in your plan file:

```hcl
resource "panos_address_objects" "ao1" {
    object {
        name = "foo"
        type = "ip-netmask"
        value = "10.1.1.1"
    }
    ...
}

resource "panos_address_objects" "ao2" {
    object {
        name = "bar"
        type = "ip-netmask"
        value = "10.1.1.2"
    }
    ...
}

resource "panos_security_rule_group" "grp" {
    rule {
        name = "allow foo"
        source_addresses = "foo"
        destination_addresses = ["any"]
        action = "allow"
        ...
    }
    rule {
        name = "allow bar"
        source_addresses = "bar"
        destination_addresses = ["any"]
        action = "allow"
        ...
    }

    depends_on = [
        panos_address_objects.ao1,
        panos_address_objects.ao2,
    ]

    lifecycle {
        create_before_destroy = true
    }
}
```


## Example Usage

### Standard usage

```hcl
resource "panos_address_objects" "example" {
    object {
        name = "web dmz"
        type = "ip-netmask"
        value = "10.1.1.50"
    }
    object {
        name = "web proxy dmz"
        type = "ip-netmask"
        value = "10.1.1.200"
    }
    object {
        name = "example wildcard"
        type = "ip-wildcard"
        value = "*.example.com"
    }
    object {
        name = "sales network"
        type = "ip-range"
        value = "192.168.50.100-192.168.50.200"
    }
    object {
        name = "disneyland"
        type = "fqdn"
        value = "disneyland.com"
    }
    object {
        name = "disneyworld"
        type = "fqdn"
        value = "disneyworld.com"
    }

    lifecycle {
        create_before_destroy = true
    }
}

# Dynamically find everything that ends with " dmz"
output "dmz_stuff" {
    value = [
        for x in panos_address_objects.example.object :
        s.name if substr(s.name, -4, -1) == " dmz"
    ]
}

# Dynamically find everything that starts with "disney"
output "disney_stuff" {
    value = [
        for x in panos_address_objects.example.object :
        s.name if substr(s.name, 0, 6) == "disney"
    ]
}
```

### Make 50 Address Objects

-> Using terraform's built-in
[`setproduct`](https://www.terraform.io/language/functions/setproduct)
to create address objects in a loop is fine for a demo, but this adds overhead
to Terraform core to do this extra processing before the plan can be built and is
thus not recommended for a production deployment.

```hcl
# Make address objects like "test1_1", "test1_2", ...
resource "panos_address_objects" "example" {
    dynamic "object" {
        for_each = setproduct(range(1, 6), range(1, 11))
        content {
            name = "test${object.value[0]}_${object.value[1]}"
            type = "ip-netmask"
            value = "10.${object.value[0]}.${object.value[1]}.0/24"
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys to put the address object into (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `object` - (Required, repeatable) An `object` spec, as defined below.

The `object` spec support the following arguments:

* `name` - (Required) The address object's name.
* `type` - The type of address object.  This can be `ip-netmask`
  (default), `ip-range`, `fqdn`, or `ip-wildcard` (PAN-OS 9.0+).
* `value` - (Required) The address object's value.  This can take various
  forms depending on what type of address object this is, but can be something
  like `192.168.80.150` or `192.168.80.0/24`.
* `description` - The address object's description.
* `tags` - (list) List of administrative tags.
