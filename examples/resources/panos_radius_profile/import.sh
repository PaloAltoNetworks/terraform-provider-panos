# A RADIUS profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "radius-basic"
# }
terraform import panos_radius_profile.example $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"radius-basic"}' | base64)
