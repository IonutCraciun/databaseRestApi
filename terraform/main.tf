terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
    }
  }
}

# COMANDS
# terraform plan
# terraform apply -auto-approve
# terraform state list
# terraform state show "resource"
# terraform output
# terraform refresh    || refresh your states without deploying anything
# terraform destroy -target aws_instance.firstVM
# terraform apply -target aws_instance.firstVM
# terraform apply -var "cidr-block-variable=10.0.1.0/24"
# -+-terraform looks automatically after terraform.tfvars for variable substitution
# terraform apply -var-file file.tfvars

variable "access_key" {
  description = "AWS access key"
  type = string
}

variable "secret_key" {
  description = "AWS secret key"
  type = string
}

# Configure the AWS Provider
provider "aws" {
  region = "us-east-1"
  access_key = var.access_key
  secret_key = var.secret_key
}

# Create vpc
resource "aws_vpc" "production-vpc" {
  cidr_block = "10.0.0.0/16"

    tags = {
    Name = "production"
  }
}

# Create internet gateway
resource "aws_internet_gateway" "prod-internet-gateway" {
  vpc_id = aws_vpc.production-vpc.id

  tags = {
    Name = "prod"
  }
}

# Create route table
resource "aws_route_table" "prod-route-table" {
  vpc_id = aws_vpc.production-vpc.id
  
  # All the trafic from the vpc can access the internet through the internet gateway
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.prod-internet-gateway.id
  }

  route {
    ipv6_cidr_block        = "::/0"
    gateway_id = aws_internet_gateway.prod-internet-gateway.id
  }
  tags = {
    Name = "prod-route-table"
  }
}

# Variable for subnet cidr block
variable "cidr-block-variable" {
  description = "cidr block for subnet"
  default = "10.0.20.0/24" # if the user is not sending any values: through -var-file, -var or input it will use default
  type = string
}

# Create a subnet
resource "aws_subnet" "subnet-1" {
  vpc_id = aws_vpc.production-vpc.id
  # cidr_block = "10.0.1.0/24"
  cidr_block = var.cidr-block-variable
  availability_zone = "us-east-1a" // "manually select a availability zone". in a region there are multiple AZs

  tags = {
    Name = "prod-subnet-1"
  }
}

# Associate subnet with route table
resource "aws_route_table_association" "route-table-association" {
  subnet_id      = aws_subnet.subnet-1.id
  route_table_id = aws_route_table.prod-route-table.id
}

# Create security group to allow trafic on port 22,80,443
resource "aws_security_group" "allow_ssh_http" {
  name        = "allow_ssh_http"
  description = "Allow SSH/HTTP/ inbound traffic"
  vpc_id      = aws_vpc.production-vpc.id

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]   // allow access from any IP
  }

    ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]   // allow access from any IP
  }

    ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]   // allow access from any IP
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"  // any protocol
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow_ssh_http_https"
  }
}

# Create a network interface in the previously created subnet
resource "aws_network_interface" "web-server-interface" {
  subnet_id       = aws_subnet.subnet-1.id
  private_ips     = ["10.0.1.50"]
  security_groups = [aws_security_group.allow_ssh_http.id]


  # Don't attach it now because we don't have an EC2 instance deployed/declared yet
  # attachment {
  #   instance     = aws_instance.test.id
  #   device_index = 1
  # }
}

# Associate elastic ip with network interface declared previously
resource "aws_eip" "one" {
  vpc                       = true
  network_interface         = aws_network_interface.web-server-interface.id
  associate_with_private_ip = "10.0.1.50"

  # EIP may require IGW to exist prior to association. Use depends_on to set an explicit dependency on the IGW.
  depends_on = [aws_internet_gateway.prod-internet-gateway] 
}

resource "aws_instance" "firstVM" {
  ami           = "ami-0dba2cb6798deb6d8"   //ubuntu
  instance_type = "t2.micro"
  availability_zone = "us-east-1a" // the same as subnet
  key_name = "mykey"

  network_interface {
    device_index = 0
    network_interface_id = aws_network_interface.web-server-interface.id
  }

  # Here we run commands aka magic
  user_data = file("install.sh")

  tags = {
    Name = "server"
  }
}

output "elastic-ip-for-web-server" {
  value = aws_eip.one.public_ip
}

# resource "<provider>_<resource_type>" "name" {
#     config options.....
#     key = "value"
#     key2 = "another value"
# }