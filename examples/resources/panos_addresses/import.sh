# Addresses can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   names = [
#     "foo",
#     "bar"
#   ]
# }
terraform import panos_addresses.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"names":["foo","bar"]}' | base64)
