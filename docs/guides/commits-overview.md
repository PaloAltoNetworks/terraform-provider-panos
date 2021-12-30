---
page_title: "Commit Overview"
subcategory: "Commits"
---

# Commit Overview

As of right now, Terraform does not provide native support for commits, so
commits are handled out-of-band.  Refer to
[this issue](https://github.com/PaloAltoNetworks/terraform-provider-panos/issues/6)
for more information.

Please refer to the specific commit guides depending on what type
of commit you need to perform.  You'll want to save the source code something
obvious that matches it's functionality, such as `firewall-commit.go`.

Compile the source code and put it somewhere in your `$PATH` (such as
`$HOME/bin`):

```bash
$ go get github.com/PaloAltoNetworks/pango
$ go build firewall-commit.go
$ mv firewall-commit ~/bin
$ firewall-commit -h
```

Finally, you can invoke this binary after `terraform apply` or `terraform
destroy`:

```bash
$ terraform apply && firewall-commit -config fwauth.json 'My commit comment'
```

The first trailing CLI arg is the commit comment.  If there is
no CLI arg present then no commit comment is given to PAN-OS.

The authentication credentials can be given multiple ways, and if all are
present then this is the order, from highest to lowest priority:

!> Providing authentication credentials via CLI argument is insecure and
is not recommended.

1. CLI arguments
2. Environment variables
3. JSON authentication credential file

Refer to the panos provider argument reference documentation for more
information on the JSON config file and the environment variables that are used.
