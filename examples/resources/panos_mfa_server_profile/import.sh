# An MFA server profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "mfa-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "okta-mfa-profile"
# }
terraform import panos_mfa_server_profile.example $(echo '{"location":{"template":{"name":"mfa-template","panorama_device":"localhost.localdomain"}},"name":"okta-mfa-profile"}' | base64)
