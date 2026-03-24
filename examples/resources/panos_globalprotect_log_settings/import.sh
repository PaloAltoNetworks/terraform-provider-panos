# A globalprotect log setting can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-gp-settings"
# }
terraform import panos_globalprotect_log_settings.example $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"example-gp-settings"}' | base64)
