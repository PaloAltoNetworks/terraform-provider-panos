---
page_title: 'Location Argument'
---

The v2 provider adds a new 'location' argument to all resources and data sources, allowing users to explicitly specify the configuration's location. This guide lists available locations that you can use based on your requirements.

#### Unmanaged Firewall

```hcl
location = {
    ngfw_device = "localhost.localdomain"
    name = "vsys1"
}
```

#### Panorama

```hcl
location = {
    panorama = {
        panorama_device = "localhost.localdomain"
    }
}
```

#### Panorama Managed Firewall

```hcl
location = {
    from_panorama_shared = {}
}

location = {
    from_panorama_vsys = {
        vsys = "vsys1"
    }
}
```

#### Specific Device Group

```hcl
location = {
  device_group = {
    panorama_device = "localhost.localdomain"
    name = ""
  }
}
```

#### Specific Template

```hcl
location = {
  template = {
    vsys = "vsys1"
    panorama_device = "localhost.localdomain"
    name = ""
    ngfw_device = "localhost.localdomain"
  }
}
```

#### Specific Template Stack

```hcl
location = {
  template_stack = {
    panorama_device = "localhost.localdomain"
    name = ""
    ngfw_device = "localhost.localdomain"
  }
}
```

#### Common (Panorama or NGFW)

```hcl
location = {
  shared = {}
}
```
