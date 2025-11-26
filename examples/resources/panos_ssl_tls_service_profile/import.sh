# An SSL/TLS Service Profile can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     template = {
#       name            = "example-template"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   name = "basic-ssl-profile"
# }
terraform import panos_ssl_tls_service_profile.basic $(echo '{"location":{"template":{"name":"example-template","panorama_device":"localhost.localdomain"}},"name":"basic-ssl-profile"}' | base64)
