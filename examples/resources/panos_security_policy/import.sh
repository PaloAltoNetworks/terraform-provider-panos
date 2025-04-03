# The entire policy can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#
#     names = [
#         "rule-1", <- the first rule in the policy
#     ]
# }
terraform import panos_security_policy.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["rule-1"]}' | base64)
