---
page_title: "Audit Comment Logic for All Policy Rules"
---

# Audit Comment Logic for All Policy Rules

Audit comments are supported for the applicable PAN-OS versions and appear as a normal argument for all policy rules (aka - security rule, NAT rule, etc).  However, the way that audit comments are supported across all resources is not standard:

* the audit comment will always be an empty string in the local state file
* the audit comment is only applied if there is a change to a rule's spec (creation or update)

Let's take an example.

Let's say you start out with this (I'll omit required but otherwise not applicable arguments for clarity):

```hcl
resource "panos_security_rule_group" "g1" {
    rule {
        name = "one"
        description = "foo"
        audit_comment = "first"
    }
    rule {
        name = "two"
        description = "foo"
        audit_comment = "first"
    }
    rule {
        name = "three"
        description = "foo"
        audit_comment = "first"
    }
}
```

After `terraform apply`, all 3 rules will be created as they did not previously exist.  Because `one`, `two`, and `three` were created, all 3 will have their audit comments applied.

Now, let's say I update my plan and it now looks like this:

```hcl
resource "panos_security_rule_group" "g1" {
    rule {
        name = "one"
        description = "foo"
        audit_comment = "second"
    }
    rule {
        name = "two"
        description = "bar"
        audit_comment = "second"
    }
    rule {
        name = "three"
        description = "foo"
        audit_comment = "second"
    }
    rule {
        name = "four"
        description = "foo"
        audit_comment = "second"
    }
}
```

So we've added "four" and made an update to "two"'s description.  In this situation, `terraform apply` will apply the audit comment for "two" (due to the changed description) and "four" (because it's a new rule).  Even though the audit comments for "one" and "three" have changed, diffs are suppressed for the `audit_comment` field.  Thus, there is no change for the specs for "one" and "three", so the audit comments associated with them are not applied.

The audit comment stored in the local state file will always be an empty string. This is because of a number of reasons, but the easiest to communicate is performance.

The PAN-OS API does not allow for the retrieval of every single rule's audit comment in one API call the same way it allows for retrieving their configuration. If you have 1,000 security rules, that's 1,000 extra API calls to get the audit comment configuration. Let's say the round trip is half a second for each API call, that's 500 seconds, or just over 8 minutes just to complete the read operation.

This is why the logic for `audit_comment` is basically "apply the `audit_comment` when the rule spec changes," and why diffs are disabled for `audit_comment` params across all policy resources.
