variable "panorama_version" {
  type = "string"
}

resource "random_id" "name" {
  byte_length = 4
}

resource "tls_private_key" "default" {
  algorithm = "RSA"
}

provider "aws" {
  version = "~> 2.0"
}

resource "aws_vpc" "default" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

resource "aws_subnet" "default" {
  vpc_id                  = "${aws_vpc.default.id}"
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = "${aws_vpc.default.id}"
  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

resource "aws_route_table" "r" {
  vpc_id = "${aws_vpc.default.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.main.id}"
  }

  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

resource "aws_route_table_association" "main" {
  subnet_id      = "${aws_subnet.default.id}"
  route_table_id = "${aws_route_table.r.id}"
}

resource "aws_security_group" "main" {
  name        = "tf-default-panorama"
  description = "main"
  vpc_id      = "${aws_vpc.default.id}"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

resource "aws_key_pair" "default" {
  key_name   = "tf-acc-panorama-${random_id.name.hex}"
  public_key = "${tls_private_key.default.public_key_openssh}"
}

data "aws_ami" "panorama" {
  most_recent = true

  filter {
    name   = "name"
    values = ["Panorama-AWS-${var.panorama_version}-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["679593333241"] # Palo Alto Networks
}

resource "aws_instance" "main" {
  ami = "${data.aws_ami.panorama.id}"

  instance_type               = "m4.2xlarge"
  subnet_id                   = "${aws_subnet.default.id}"
  vpc_security_group_ids      = ["${aws_security_group.main.id}"]
  key_name                    = "${aws_key_pair.default.key_name}"
  associate_public_ip_address = true

  tags = {
    Name = "tf-acc-panorama-${random_id.name.hex}"
  }
}

output "hostname" {
  value = "${aws_instance.main.public_ip}"
}
output "ssh_private_key" {
  value = "${tls_private_key.default.private_key_pem}"
}
