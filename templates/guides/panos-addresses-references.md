---
page_title: Referencing `panos_addresses` resources during updates
---

PAN-OS provider gives two resources that can be used to manage Address Objects on the devices: `panos_address` and `panos_addresses`. The latter one can be used to manage a batch of address objects on the device, making use of the Multi-config requests to optimize CRUD operations on objects by batching multiple Read, Update and Create operations into a smaller number of requests, allowing users to manage hundreds and thousands of objects in a more optimized manner.

This approach, however, suffers from a drawback when the list of address objects has to be modified over time. Due to how dependencies between resources work within Terraform, referencing the `panos_addresses` resource from other resources (e.g. `panos_security_policy` for managing security policies) can lead to errors in the update operations if both the `panos_addresses` resource and the resource that references it have to be changed at the same time.

During apply, terraform builds a dependency graph between resources, and referencing `panos_addresses` from another resource creates an implicit dependency between them. If a `panos_addresses` resource is modified within a plan to remove a subset of address objects, PAN-OS will reject such a modification as long as those address objects are still referenced from another resource (e.g. security policy).

Such a change has to be applied in two steps: first, the security policy has to be modified to remove all address object references that are no longer needed, and only then can another change be applied to delete those address objects from the device.

### Creating initial address objects and using them from an address group as reference

In the following example, we start with a single `panos_addresses` resource and an address group that uses it as a reference to create a group of addresses:

```hcl
resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = "example-dg"
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  name = "example-address-group"

  static = [for k, v in panos_addresses.example.addresses : k]
}
```

This configuration can be applied resulting in an expected final state:

```
$ terraform apply
[...]
Terraform will perform the following actions:

  # panos_address_group.example will be created
  + resource "panos_address_group" "example" {
      + location = {
          + device_group = {
              + name            = "example-dg"
              + panorama_device = "localhost.localdomain"
            }
        }
      + name     = "example-address-group"
      + static   = [
          + "example-addr1",
          + "example-addr2",
        ]
    }

  # panos_addresses.example will be created
  + resource "panos_addresses" "example" {
      + addresses = {
          + "example-addr1" = {
              + ip_netmask = "10.0.0.1/32"
            },
          + "example-addr2" = {
              + ip_netmask = "10.0.0.2/32"
            },
        }
      + location  = {
          + device_group = {
              + name            = "example-dg"
              + panorama_device = "localhost.localdomain"
            }
        }
    }

  # panos_device_group.example will be created
  + resource "panos_device_group" "example" {
      + location = {
          + panorama = {
              + panorama_device = "localhost.localdomain"
            }
        }
      + name     = "example-dg"
    }

Plan: 3 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

panos_device_group.example: Creating...
panos_device_group.example: Creation complete after 3s [name=example-dg]
panos_addresses.example: Creating...
panos_addresses.example: Creation complete after 5s
panos_address_group.example: Creating...
panos_address_group.example: Creation complete after 3s [name=example-address-group]

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
$
```

### Adding new address objects and using them as a reference from an address group

New address objects can be created with that configuration, resulting in an expected final state:

```hcl
resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = "example-dg"
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    },
    "example-addr3" = {
      ip_netmask = "10.0.0.3/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  name = "example-address-group"

  static = [for k, v in panos_addresses.example.addresses : k]
}
```

This configuration also applies correctly:

```
~/s/p/terraform-test ❯❯❯ terraform apply
[...]
panos_device_group.example: Refreshing state... [name=example-dg]
panos_addresses.example: Refreshing state...
panos_address_group.example: Refreshing state... [name=example-address-group]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # panos_address_group.example will be updated in-place
  ~ resource "panos_address_group" "example" {
        name     = "example-address-group"
      ~ static   = [
            # (1 unchanged element hidden)
            "example-addr2",
          + "example-addr3",
        ]
        # (1 unchanged attribute hidden)
    }

  # panos_addresses.example will be updated in-place
  ~ resource "panos_addresses" "example" {
      ~ addresses = {
          + "example-addr3" = {
              + ip_netmask = "10.0.0.3/32"
            },
            # (2 unchanged elements hidden)
        }
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 2 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

panos_addresses.example: Modifying...
panos_addresses.example: Modifications complete after 3s
panos_address_group.example: Modifying... [name=example-address-group]
panos_address_group.example: Modifications complete after 3s [name=example-address-group]

Apply complete! Resources: 0 added, 2 changed, 0 destroyed.
$
```

### Removing the address object from the device

A removal of address objects from the device will however result in an error during apply.

```hcl
resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = "example-dg"
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  name = "example-address-group"

  static = [for k, v in panos_addresses.example.addresses : k]
}
```

Note that _example-addr3_ has been removed from the panos_addresses.example addresses map. When `terraform apply` is executed, initially terraform will present a reasonable plan of changes, PAN-OS will however reject those changes during apply phase:

```hcl
$ terraform apply
[...]
panos_device_group.example: Refreshing state... [name=example-dg]
panos_addresses.example: Refreshing state...
panos_address_group.example: Refreshing state... [name=example-address-group]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # panos_address_group.example will be updated in-place
  ~ resource "panos_address_group" "example" {
        name     = "example-address-group"
      ~ static   = [
            # (1 unchanged element hidden)
            "example-addr2",
          - "example-addr3",
        ]
        # (1 unchanged attribute hidden)
    }

  # panos_addresses.example will be updated in-place
  ~ resource "panos_addresses" "example" {
      ~ addresses = {
          - "example-addr3" = {
              - ip_netmask = "10.0.0.3/32" -> null
            },
            # (2 unchanged elements hidden)
        }
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 2 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

panos_addresses.example: Modifying...
╷
│ Error: Error while updating entries
│
│   with panos_addresses.example,
│   on main.tf line 743, in resource "panos_addresses" "example":
│  743: resource "panos_addresses" "example" {
│
│ Failed to execute MultiConfig command: example-addr3 cannot be deleted because of references from: | device-group -> example-dg -> address-group -> example-address-group -> static
╵
$
```

## Workarounds

### Using `panos_address` resource to manage address objects independently

If the number of managed address objects is small enough, using `panos_address` resource instead will work around this issue. The above configuration can be translated as such:

```hcl
resource "panos_address_group" "example" {
  location = {
    device_group = {
      name = panos_device_group.example.name
    }
  }

  name        = "example-address-group"
  description = "example address group"
  static      = [for k in panos_address.example : k.name]

  depends_on = [panos_address.example]

  lifecycle {
    create_before_destroy = true
  }
}

resource "panos_address" "example" {
  location = {
    device_group = {
      name = panos_device_group.example.name
    }
  }

  for_each = tomap({
    "example-addr1" = {
      ip_netmask  = "10.0.0.1/32"
    }
    "example-addr2" = {
      ip_netmask        = "10.0.0.2/32"
    }
    "example-addr3" = {
      ip_netmask        = "10.0.0.3/32"
    }
  })

  name        = each.key
  description = each.value.description
  ip_netmask  = lookup(each.value, "ip_netmask", null)
  fqdn        = lookup(each.value, "fqdn", null)
}

resource "panos_device_group" "example" {
  location = {
    panorama = {}
  }

  name = "example-device-group"
```

By adding `create_before_destroy` terraform creates a sequence of steps where `panos_address_group` is modified after new address objects are created, but before other address objects are deleted. This cannot be done when both are managed within a single resource.

### Workflow modification

If the number of address objects makes `panos_address` performance unsatisfactory, a workflow has to be modified, where modifications to address objects and any references are split into multiple steps. For example, the original example can be modified to remove `panos_addresses` reference from `panos_address_group`:

```hcl
locals {
  address_objects = ["example-addr1", "example-addr2", "example-addr3"]
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    }
    "example-addr3" = {
      ip_netmask = "10.0.0.3/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  depends_on = [panos_addresses.example]

  name = "example-address-group"

  static = local.address_objects
}
```

If some address objects are scheduled for removal, the configuration is first modified to change `local.address_objects` and remove those addresses from there:

```hcl
locals {
  address_objects = ["example-addr1", "example-addr2"]
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    }
    "example-addr3" = {
      ip_netmask = "10.0.0.3/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  depends_on = [panos_addresses.example]

  name = "example-address-group"

  static = local.address_objects
}
```

This change will only result in `panos_address_object` resource being updated:

```
$ terraform apply
[...]
panos_device_group.example: Refreshing state... [name=example-dg]
panos_addresses.example: Refreshing state...
panos_address_group.example: Refreshing state... [name=example-address-group]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # panos_address_group.example will be updated in-place
  ~ resource "panos_address_group" "example" {
        name     = "example-address-group"
      ~ static   = [
            # (1 unchanged element hidden)
            "example-addr2",
          - "example-addr3",
        ]
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

panos_address_group.example: Modifying... [name=example-address-group]
panos_address_group.example: Modifications complete after 3s [name=example-address-group]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
$
```

Once the reference has been removed from the device, in the second step we can delete address objects:

```hcl
locals {
  address_objects = ["example-addr1", "example-addr2"]
}

resource "panos_addresses" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  addresses = {
    "example-addr1" = {
      ip_netmask = "10.0.0.1/32"
    },
    "example-addr2" = {
      ip_netmask = "10.0.0.2/32"
    }
  }
}

resource "panos_address_group" "example" {
  location = { device_group = { name = panos_device_group.example.name } }

  depends_on = [panos_addresses.example]

  name = "example-address-group"

  static = local.address_objects
}

```

The new plan will also apply on the device as expected:

```
$ terraform apply
[...]
panos_device_group.example: Refreshing state... [name=example-dg]
panos_addresses.example: Refreshing state...
panos_address_group.example: Refreshing state... [name=example-address-group]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # panos_addresses.example will be updated in-place
  ~ resource "panos_addresses" "example" {
      ~ addresses = {
          - "example-addr3" = {
              - ip_netmask = "10.0.0.3/32" -> null
            },
            # (2 unchanged elements hidden)
        }
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

panos_addresses.example: Modifying...
panos_addresses.example: Modifications complete after 2s

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
$
```
