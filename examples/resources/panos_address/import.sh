# Addresses can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name            = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
# 
#   name = "addr1"
# }
terraform import panos_address.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"name":"addr1"}' | base64)