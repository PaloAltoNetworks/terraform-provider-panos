# Returns all address objects in vsys1.
data "panos_nested_address_object_list" "example" {
    location = {
        vsys = {}
    }
}

# Returns all running config IP Netmask address objects in vsys2
# that end in "_DMZ".
data "panos_nested_address_object_list" "example" {
    action = "show"
    location = {
        vsys = {
            name = "vsys2"
        }
    }

    filter = "ip_netmask is-not-nil && name ends-with '_DMZ'"
    quote = "'"
}
