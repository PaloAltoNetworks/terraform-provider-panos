# An iptag log setting can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-iptag-settings"
# }
terraform import panos_iptag_log_settings.example $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"example-iptag-settings"}' | base64)
