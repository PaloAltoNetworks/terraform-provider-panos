# A TACACS+ profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-tacacs-profile"
# }
terraform import panos_tacacs_plus_profile.example $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"example-tacacs-profile"}' | base64)
