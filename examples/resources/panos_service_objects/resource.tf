# Create a group of one service object in vsys1.
resource "panos_service_objects" "x" {
  location = {
    vsys = {}
  }
  objects = [
    {
      name        = "bandit"
      description = "Made by Terraform"
      protocol = {
        udp = {
          destination_port = 12345
          override = {
            no = true
          }
        }
      }
    },
  ]
}
