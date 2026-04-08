# A schedule can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name            = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "daily-schedule"
# }
terraform import panos_schedule.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"name":"daily-schedule"}' | base64)
