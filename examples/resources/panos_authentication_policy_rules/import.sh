#!/bin/bash

# A set of authentication policy rules can be imported by providing the following base64 encoded object as the ID
# {
#     location = {
#         device_group = {
#         name = "example-device-group"
#         rulebase = "pre-rulebase"
#         panorama_device = "localhost.localdomain"
#         }
#     }
#
#     position = { where = "after", directly = true, pivot = "guest-wifi-auth" }
#
#     names = [
#         "employee-byod-auth",
#         "contractor-limited-access"
#     ]
# }
terraform import panos_authentication_policy_rules.corporate_users $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain","rulebase":"pre-rulebase"}},"names":["employee-byod-auth","contractor-limited-access"],"position":{"directly":true,"pivot":"guest-wifi-auth","where":"after"}}' | base64)
