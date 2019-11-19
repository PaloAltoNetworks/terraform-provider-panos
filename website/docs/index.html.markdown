---
layout: "panos"
page_title: "Provider: panos"
sidebar_current: "docs-panos-index"
description: |-
  Provider panos is used to interact with Palo Alto Networks NGFW and Panorama.
---

# Provider panos

PAN-OS&reg; is the operating system for Palo Alto Networks&reg; NGFWs and
Panorama&trade;. The panos provider allows you to manage various aspects
of a firewall's or a Panorama's config, such as data interfaces and security
policies.

Use the navigation to the left to read about the available Panorama and NGFW
resources.

## Versioning

The panos provider has support for PAN-OS 6.1 - 9.0.

Some resources may contain variables that are only applicable for newer
versions of PAN-OS.  If this is the case, then make sure to use
[conditionals](https://www.terraform.io/docs/configuration/expressions.html#conditional-expressions)
along with the `panos_system_info` data source to only set these variables
when the version of PAN-OS is appropriate.

One such resource is `panos_ethernet_interface` and the `ipv4_mss_adjust`
parameter.  Doing the following is one way to correctly configure this
parameter only when it's applicable:

```hcl
data "panos_system_info" "config" {}

data "panos_ethernet_interface" "eth1" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
    adjust_tcp_mss = true
    ipv4_mss_adjust = "${data.panos_system_info.config.version_major >= 8 ? 42 : 0}"
    # ...
}
```

## Commits

As of right now, Terraform does not provide native support for commits, so
commits are handled out-of-band.  Please use the following for commits:

```go
package main

import (
    "encoding/json"
    "flag"
    "log"
    "os"

    "github.com/PaloAltoNetworks/pango"
)

type Credentials struct {
    Hostname string `json:"hostname"`
    Username string `json:"username"`
    Password string `json:"password"`
    ApiKey string `json:"api_key"`
    Protocol string `json:"protocol"`
    Port uint `json:"port"`
    Timeout int `json:"timeout"`
    VerifyCertificate bool `json:"verify_certificate"`
}

func getCredentials(configFile, hostname, username, password, apiKey string) (Credentials) {
    var (
        config Credentials
        val string
        ok bool
    )

    // Auth from the config file.
    if configFile != "" {
        fd, err := os.Open(configFile)
        if err != nil {
            log.Fatalf("ERROR: %s", err)
        }
        defer fd.Close()

        dec := json.NewDecoder(fd)
        err = dec.Decode(&config)
        if err != nil {
            log.Fatalf("ERROR: %s", err)
        }
    }

    // Auth from env variables.
    if val, ok = os.LookupEnv("PANOS_HOSTNAME"); ok {
        config.Hostname = val
    }
    if val, ok = os.LookupEnv("PANOS_USERNAME"); ok {
        config.Username = val
    }
    if val, ok = os.LookupEnv("PANOS_PASSWORD"); ok {
        config.Password = val
    }
    if val, ok = os.LookupEnv("PANOS_API_KEY"); ok {
        config.ApiKey = val
    }

    // Auth from CLI args.
    if hostname != "" {
        config.Hostname = hostname
    }
    if username != "" {
        config.Username = username
    }
    if password != "" {
        config.Password = password
    }
    if apiKey != "" {
        config.ApiKey = apiKey
    }

    if config.Hostname == "" {
        log.Fatalf("ERROR: No hostname specified")
    } else if config.Username == "" && config.ApiKey == "" {
        log.Fatalf("ERROR: No username specified")
    } else if config.Password == "" && config.ApiKey == "" {
        log.Fatalf("ERROR: No password specified")
    }

    return config
}

func main() {
    var (
        err error
        configFile, hostname, username, password, apiKey string
        job uint
    )

    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

    flag.StringVar(&configFile, "config", "", "JSON config file with panos connection info")
    flag.StringVar(&hostname, "host", "", "PAN-OS hostname")
    flag.StringVar(&username, "user", "", "PAN-OS username")
    flag.StringVar(&password, "pass", "", "PAN-OS password")
    flag.StringVar(&apiKey, "key", "", "PAN-OS API key")
    flag.Parse()

    config := getCredentials(configFile, hostname, username, password, apiKey)

    fw := &pango.Firewall{Client: pango.Client{
        Hostname: config.Hostname,
        Username: config.Username,
        Password: config.Password,
        ApiKey: config.ApiKey,
        Protocol: config.Protocol,
        Port: config.Port,
        Timeout: config.Timeout,
        Logging: pango.LogOp | pango.LogAction,
    }}
    if err = fw.Initialize(); err != nil {
        log.Fatalf("Failed: %s", err)
    }

    job, err = fw.Commit(flag.Arg(0), nil, true, true, false, true)
    if err != nil {
        log.Fatalf("Error in commit: %s", err)
    } else if job == 0 {
        log.Printf("No commit needed")
    } else {
        log.Printf("Committed config successfully")
    }
}
```

Compile the above, put it somewhere in your `$PATH` (such as `$HOME/bin`),
then invoke it after `terraform apply` and `terraform destroy`:

```bash
$ go get github.com/PaloAltoNetworks/pango
$ go build commit.go
$ mv commit ~/bin
$ terraform apply && commit -config fwauth.json 'My commit comment'
```

The first trailing CLI arg is assumed to be the commit comment.  If there is
no CLI arg present then no commit comment is given to PAN-OS.

The authentication credentials can be given multiple ways, and if all are
present then this is the order, from lowest to highest priority:

1. JSON authentication credential file
2. Environment variables
3. CLI arguments (WARNING: this is insecure)

See the argument reference section below for more info on the JSON config
file and supported environment variables.

## PAN-OS API Key

API connections to PAN-OS require an
[API key](https://www.paloaltonetworks.com/documentation/71/pan-os/xml-api/get-started-with-the-pan-os-xml-api/get-your-api-key).
If you do not provide the API key to the panos provider, then the API key is
generated before every single API call.  Thus, some slight speed gains can be
realized in the panos provider by specifying the API key instead of the
username/password combo.  The following may be used to generate the API key:

```go
package main

import (
    "fmt"
    "os"

    "github.com/PaloAltoNetworks/pango"
)

func main() {
    var (
        hostname, username, password string
        ok bool
    )

    if hostname, ok = os.LookupEnv("PANOS_HOSTNAME"); !ok {
        os.Stderr.WriteString("PANOS_HOSTNAME must be set\n")
        return
    }
    if username, ok = os.LookupEnv("PANOS_USERNAME"); !ok {
        os.Stderr.WriteString("PANOS_USERNAME must be set\n")
        return
    }
    if password, ok = os.LookupEnv("PANOS_PASSWORD"); !ok {
        os.Stderr.WriteString("PANOS_PASSWORD must be set\n")
        return
    }

    fw := &pango.Firewall{Client: pango.Client{
        Hostname: hostname,
        Username: username,
        Password: password,
        Logging: pango.LogQuiet,
    }}
    if err := fw.Initialize(); err != nil {
        os.Stderr.WriteString(fmt.Sprintf("Failed initialize: %s\n", err))
        return
    }
    os.Stdout.WriteString(fmt.Sprintf("%s\n", fw.ApiKey))
}
```

Then execute it like this:

```bash
$ go get github.com/PaloAltoNetworks/pango
$ go run make_api_key.go
```

The API key is output to stdout, but you can redirect this to a file using
normal shell redirection if desired:

```bash
$ go run make_api_key.go > my_api_key.txt
```

Connection information for the above is expected to be set as environment
variables:

* `PANOS_HOSTNAME`
* `PANOS_USERNAME`
* `PANOS_PASSWORD`

## AWS / GCP Considerations

There are [a few types](https://aws.amazon.com/marketplace/seller-profile?id=0ed48363-5064-4d47-b41b-a53f7c937314)
of PAN-OS VMs available to bring up in AWS.  Both these VMs as well as the ones
that can be deployed in Google Cloud Platform are different in that
the `admin` password is unset, but it has an SSH key associated with it.  As
the panos Terraform provider package authenticates via username/password, an
initialization step of configuring a password using the given SSH key is
required.  Right now, this initialization step requires manual intervention;
the user must download this SSH key, at which point the following may be used
to automate this initialization:

```go
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

// Various prompts.
var (
    P1 *regexp.Regexp
    P2 *regexp.Regexp
    P3 *regexp.Regexp
)

func init() {
    P1 = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+> `)
    P2 = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\._\-]+@[a-zA-Z][a-zA-Z0-9\._\-]+# `)
    P3 = regexp.MustCompile(`(Enter|Confirm) password\s+:\s+?`)
}

// Globals to handle I/O.
var (
    stdin io.Writer
    stdout io.Reader
    buf [65 * 1024]byte
)

// ReadTo reads from stdout until the desired prompt is encountered.
func ReadTo(prompt *regexp.Regexp) (string, error) {
    var i int

    for {
        n, err := stdout.Read(buf[i:])
        if n > 0 {
            os.Stdout.Write(buf[i:i + n])
        }
        if err != nil {
            return "", err
        }
        i += n
        if prompt.Find(buf[:i]) != nil {
            return string(buf[:i]), nil
        }
    }
}

// Perform user initialization.
func panosInit() error {
    var err error

    // Load environment variables.
    hostname := os.Getenv("PANOS_HOSTNAME")
    username := os.Getenv("PANOS_USERNAME")
    password := os.Getenv("PANOS_PASSWORD")

    // Sanity check input.
    if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" || hostname == "" || username == "" || password == "" {
        u := []string{
            fmt.Sprintf("Usage: %s <key_file>", os.Args[0]),
            "",
            "This will connect to a PAN-OS NGFW and perform initial config:",
            "",
            " * Adds the user as a superuser (if not the admin user)",
            " * Sets the user's password",
            " * Commit",
            "",
            "The following environment variables are required:",
            "",
            " * PANOS_HOSTNAME",
            " * PANOS_USERNAME",
            " * PANOS_PASSWORD",
        }
        for i := range u {
            fmt.Printf("%s\n", u[i])
        }
        os.Exit(0)
    }

    // Read in the ssh key file.
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        return fmt.Errorf("Failed to read SSH key file %q: %s", os.Args[1], err)
    }

    signer, err := ssh.ParsePrivateKey(data)
    if err != nil {
        return fmt.Errorf("Failed to parse private key: %s", err)
    }

    useSshKey := ssh.PublicKeys(signer)

    // Configure and open the ssh connection.
    config := &ssh.ClientConfig{
        User: "admin",
        Auth: []ssh.AuthMethod{
            useSshKey,
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", hostname), config)
    if err != nil {
        return fmt.Errorf("Failed dial: %s", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return fmt.Errorf("Failed to create session: %s", err)
    }
    defer session.Close()

    modes := ssh.TerminalModes{
        ssh.ECHO: 0,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err = session.RequestPty("vt100", 80, 80, modes); err != nil {
        return fmt.Errorf("pty request failed: %s", err)
    }

    // Get input/output pipes for the ssh connection.
    stdin, err = session.StdinPipe()
    if err != nil {
        return fmt.Errorf("setup stdin err: %s", err)
    }

    stdout, err = session.StdoutPipe()
    if err != nil {
        return fmt.Errorf("setup stdout err: %s", err)
    }

    // Invoke a shell on the remote host.
    if err = session.Start("/bin/sh"); err != nil {
        return fmt.Errorf("failed session.Start: %s", err)
    }

    // Perform initial config.
    ok := true
    commands := []struct{
        Send string
        Expect *regexp.Regexp
        Validation string
        OmitIfAdmin bool
    }{
        {"", P1, "", false},
        {"set cli pager off", P1, "", false},
        {"show system info", P1, "", false},
        {"configure", P2, "", false},
        {fmt.Sprintf("set mgt-config users %s permissions role-based superuser yes", username), P2, "", true},
        {fmt.Sprintf("set mgt-config users %s password", username), P3, "", false},
        {password, P3, "", false},
        {password, P2, "", false},
        {"commit description 'initial config'", P2, "Configuration committed successfully", false},
        {"exit", P1, "", false},
        {"exit", nil, "", false},
    }

    for _, cmd := range commands {
        if cmd.OmitIfAdmin && username == "admin" {
            continue
        }
        if cmd.Send != "" {
            stdin.Write([]byte(cmd.Send + "\n"))
        }
        if cmd.Expect != nil {
            out, err := ReadTo(cmd.Expect)
            if err != nil {
                return fmt.Errorf("Error in %q: %s", cmd.Send, err)
            }
            if cmd.Validation != "" {
                ok = ok && strings.Contains(out, cmd.Validation)
            }
            // Delay slightly before sending passwords.
            if cmd.Expect == P3 {
                time.Sleep(1 * time.Second)
            }
        } else {
            fmt.Printf("exit\n")
            session.Wait()
        }
    }

    // Completed successfully.
    return nil
}

func main() {
    if err := panosInit(); err != nil {
        fmt.Printf("\nFailed initial config: %s\n", err)
        os.Exit(1)
    }
    fmt.Printf("\nConfig initialization successful")
}
```

Compile the above, put it somewhere in your `$PATH` (such as `$HOME/bin`),
then invoke it after the device is accessible in AWS:

```bash
$ go get golang.org/x/crypto/ssh
$ go build panos_init.go
$ mv panos_init ~/bin
$ panos_init my_ssh_key.pem
```

The API key is expected to be given as the first param, while the hostname is
retrieved from the following environment variable:

* `PANOS_HOSTNAME`

The username and password are expected to be in the following environment
variables:

* `PANOS_USERNAME`
* `PANOS_PASSWORD`

If `PANOS_USERNAME` is set to `admin`, then the above will skip the step that
creates the account, as the `admin` account already exists.


## Importing Resources

Many resources support being imported.  Any resource that supports `terraform
import` will have a "Import Name" section in the documentation.  The variables
given in this section directly match up with the resource params you would
specify in your plan file.  Thus if you were importing an ethernet interface
whose import name is `<vsys>:<name>`, your import name would be something like
`vsys1:ethernet1/1`.

Of special note is the Panorama resources.  The templated resources often
have both the template and the template stack in the resource name, however
only one of these can ever be present.  Thus, the one that isn't being used
should just be an empty string.  For example, if you were trying to import
a Panorama IPv4 static route whose import name is
`<template>:<template_stack>:<virtual_router>:<name>` that resides in a
template, your import name would be something like
`myTemplate::myVirtualRouter:myStaticRouteName`.


## Example Provider Usage

```hcl
# Configure the panos provider
provider "panos" {
    hostname = "127.0.0.1"
    username = "admin"
    password = "secret"
}

# Add a new zone to the firewall
resource "panos_zone" "zone1" {
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
* `logging` - (Optional) List of logging options for the provider's connection
  to the API.  If this is unspecified, then it defaults to
  `["action", "uid"]`.
* `verify_certificate` - (Optional, bool, added in v1.6.1) For HTTPS protocol
  connections, verify that the certificate is valid.
* `json_config_file` - (Optional) The path to a JSON configuration file that
  contains any number of the provider's parameters.  If specified, the params
  present act as a last resort for any other provider param that has not been
  specified yet.

The list of strings supported for `logging` are as follows:

* `quiet` - Disables logging.  This is ignored, however, if other logging
  flags are present.
* `action` - Log `set` / `edit` / `delete`.
* `query` - Log `get`.
* `op` - Log `op`.
* `uid` - Log user-id envocations.
* `xpath` - Log the XPATH associated with various actions.
* `send` - Log the raw request sent to the device.  This is probably
  only useful in development of the provider itself.
* `receive` - Log the raw response sent back from the device.  This is probably
  only useful in development of the provider itself.

## Support

This template/solution are released under an as-is, best effort, support
policy. These scripts should be seen as community supported and Palo Alto
Networks will contribute our expertise as and when possible. We do not
provide technical support or help in using or troubleshooting the components
of the project through our normal support options such as Palo Alto Networks
support teams, or ASC (Authorized Support Centers) partners and backline
support options. The underlying product used (the VM-Series firewall) by the
scripts or templates are still supported, but the support is only for the
product functionality and not for help in deploying or using the template or
script itself. Unless explicitly tagged, all projects or work posted in our
GitHub repository (at https://github.com/PaloAltoNetworks) or sites other
than our official Downloads page on https://support.paloaltonetworks.com
are provided under the best effort policy.
