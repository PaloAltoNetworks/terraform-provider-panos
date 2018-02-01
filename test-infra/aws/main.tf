resource "random_string" "key_name" {
  length = 16
}

variable "ssh_key" {
  default = ".ssh/id_rsa.pub"
}

provider "aws" {
  region = "us-west-2"
}

resource "aws_vpc" "default" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

resource "aws_subnet" "tf_test_subnet" {
  vpc_id                  = "${aws_vpc.default.id}"
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.default.id}"

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

resource "aws_route_table" "r" {
  vpc_id = "${aws_vpc.default.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gw.id}"
  }

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

resource "aws_route_table_association" "a" {
  subnet_id      = "${aws_subnet.tf_test_subnet.id}"
  route_table_id = "${aws_route_table.r.id}"
}

resource "aws_security_group" "tf_test_sg_ssh" {
  name        = "TFACC_PANOS_INFRA"
  description = "tf_test_sg_ssh"
  vpc_id      = "${aws_vpc.default.id}"

  ingress {
    from_port   = 0
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    self        = true
  }

  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    ipv6_cidr_blocks = ["::/0"]
    self             = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

resource "aws_key_pair" "ssh_thing" {
  key_name   = "${random_string.key_name.result}"
  public_key = "${file("${var.ssh_key}")}"
}

resource "aws_instance" "tf_test" {
  #ami = "ami-7ac6491a"
  # PanOS AMI
  ami = "ami-d258f3aa"

  instance_type               = "m3.xlarge"
  subnet_id                   = "${aws_subnet.tf_test_subnet.id}"
  vpc_security_group_ids      = ["${aws_security_group.tf_test_sg_ssh.id}"]
  key_name                    = "${aws_key_pair.ssh_thing.key_name}"
  associate_public_ip_address = true

  iam_instance_profile = "${aws_iam_instance_profile.test_profile.name}"

  depends_on = ["aws_internet_gateway.gw"]

  tags {
    Name = "TFACC_PANOS_INFRA"
  }
}

# IAM things

resource "aws_iam_instance_profile" "test_profile" {
  name = "tfacc_panos_profile"
  role = "${aws_iam_role.role.name}"
}

resource "aws_iam_role" "role" {
  name = "test_panos_role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {"AWS": "*"},
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

output "endpoint" {
  value = "${aws_instance.tf_test.public_dns}"
}
