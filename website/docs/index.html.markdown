---
layout: "panos"
page_title: "Provider: PANOS"
sidebar_current: "docs-panos-index"
description: |-
  PANOS is used to interact with Palo Alto Networks NGFW.
---

# PANOS Provider

[PANOS](https://www.paloaltonetworks.com/) is used to interact with Palo Alto
Networks' NGFW.  The provider allows you to manage various aspects of the
firewall's config, such as data interfaces or security policies.

Use the navigation to the left to read about the available resources.

# Commits

As of right now, Terraform does not provide native support for commits, so
commits are handled out-of-band.  Please use the following script for commits:

```go
package main

import (
    "flag"
    "log"
    "os"

    "github.com/PaloAltoNetworks/pango"
)

func main() {
    var (
        hostname, username, password, apiKey, comment string
        ok bool
        err error
        job uint
    )

    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

    if hostname, ok = os.LookupEnv("PANOS_HOSTNAME"); !ok {
        log.Fatalf("PANOS_HOSTNAME must be set")
    }
    if username, ok = os.LookupEnv("PANOS_USERNAME"); !ok {
        log.Fatalf("PANOS_USERNAME must be set")
    }
    if password, ok = os.LookupEnv("PANOS_PASSWORD"); !ok {
        log.Fatalf("PANOS_PASSWORD must be set")
    }

    flag.StringVar(&comment, "c", "", "Commit comment")
    flag.Parse()

    fw := &pango.Firewall{Client: pango.Client{
        Hostname: hostname,
        Username: username,
        Password: password,
        ApiKey: apiKey,
        Logging: pango.LogOp | pango.LogAction,
    }}
    if err := fw.Initialize(); err != nil {
        log.Fatalf("Failed: %s", err)
    }

    job, err = fw.Commit(comment, true, true, false, true)
    if err != nil {
        log.Fatalf("Error in commit: %s", err)
    } else if job == 0 {
        log.Printf("No commit needed")
    } else {
        log.Printf("Committed config successfully")
    }
}
```

Compile the above, put it somewhere in your `$PATH` (such as $HOME/bin), then
invoke it after `terraform apply` / `terraform destroy`:

```bash
$ go build commit.go
$ mv commit ~/bin
$ terraform apply && commit -c 'My commit message'
```

## Example Usage

```hcl
# Configure the PANOS provider
provider "panos" {
    hostname = "127.0.0.1"
    username = "admin"
    password = "secret"
}

# Add a security policy to the firewall
resource "panos_security_policy" "myPolicy" {
    # ...
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Optional) This is the hostname / IP address of the firewall.  It
  must be provided, but can also be defined via the `PANOS_HOSTNAME`
  environment variable.
* `username` - (Optional) The username to authenticate to the firewall as.  It
  must be provided, but can also be defined via the `PANOS_USERNAME`
  environment variable.
* `password` - (Optional) The password for the given username. It must be
  provided, but can also be defined via the `PANOS_PASSWORD` environment
  variable.
* `api_key` - (Optional) The API key for the firewall.  If this is given, then
  the `username` and `password` settings are ignored.  This can also be defined
  via the `PANOS_API_KEY` environment variable.
* `protocol` - (Optional) The communication protocol.  This can be set to
  either `https` or `http`.  If left unspecified, this defaults to `https`.  
* `port` - (Optional) If the port number is non-standard for the desired
  protocol, then the port number to use.
* `timeout` - (Optional) The timeout for all communications with the
  firewall.  If left unspecified, this will be set to 10 seconds.
* `logging` - (Optional) Logging options for the API connection.
