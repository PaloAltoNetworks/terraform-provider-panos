# A password profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "my-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "example-password-profile"
# }
terraform import panos_password_profile.example $(echo '{"location":{"template":{"name":"my-template","panorama_device":"localhost.localdomain"}},"name":"example-password-profile"}' | base64)
