# A dynamic user group can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name            = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "developers"
# }
terraform import panos_dynamic_user_group.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"name":"developers"}' | base64)
