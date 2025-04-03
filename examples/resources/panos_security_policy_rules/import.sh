# A set of rules can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#     position = { where = "after", directly = true, pivot = "rule-2" }
#
#     names = [
#         "rule-8",
#         "rule-9"
#     ]
# }
terraform import panos_security_policy_rules.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["rule-8","rule-9"],"position":{"directly":true,"pivot":"rule-2","where":"after"}}' | base64)
