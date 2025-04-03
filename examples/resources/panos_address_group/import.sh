# An address group can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name            = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-address-group"
# }
terraform import panos_address_group.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"name":"example-address-group"}' | base64)
