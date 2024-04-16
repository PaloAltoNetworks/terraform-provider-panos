# Get info on the security rule named "foo" from vsys1.
data "panos_security_rule" "example1" {
  location = {
    vsys = {}
  }
  name = "foo"
}

# Get the running config spec for the post-rulebase rule from
# the shared location whose uuid is "9d35c525-631b-4dda-9543-93a819496245".
data "panos_security_rule" "example2" {
  filter = {
    config = "running"
  }
  location = {
    shared = {
      rulebase = "post-rulebase"
    }
  }
  uuid = "9d35c525-631b-4dda-9543-93a819496245"
}
