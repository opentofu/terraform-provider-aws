# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "aws" {
  default_tags {
    tags = var.provider_tags
  }
  ignore_tags {
    keys = var.ignore_tag_keys
  }
}

resource "aws_vpclattice_service_network_resource_association" "test" {
  resource_configuration_identifier = aws_vpclattice_resource_configuration.test.id
  service_network_identifier        = aws_vpclattice_service_network.test.id

  tags = var.resource_tags
}

resource "aws_vpclattice_service_network" "test" {
  name = var.rName

  tags = {
    Name = var.rName
  }
}


resource "aws_vpclattice_resource_configuration" "test" {
  name = var.rName

  resource_gateway_identifier = aws_vpclattice_resource_gateway.test.id

  port_ranges = ["80"]
  protocol    = "TCP"

  resource_configuration_definition {
    dns_resource {
      domain_name     = "example.com"
      ip_address_type = "IPV4"
    }
  }

  tags = var.resource_tags
}

resource "aws_vpclattice_resource_gateway" "test" {
  name       = var.rName
  vpc_id     = aws_vpc.test.id
  subnet_ids = [aws_subnet.test.id]
}

resource "aws_vpc" "test" {
  assign_generated_ipv6_cidr_block = true
  cidr_block                       = "10.0.0.0/16"

  tags = {
    Name = var.rName
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  vpc_id            = aws_vpc.test.id
  cidr_block        = "10.0.1.0/24"
  ipv6_cidr_block   = cidrsubnet(aws_vpc.test.ipv6_cidr_block, 8, 0)

  tags = {
    Name = var.rName
  }
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

variable "rName" {
  description = "Name for resource"
  type        = string
  nullable    = false
}

variable "resource_tags" {
  description = "Tags to set on resource. To specify no tags, set to `null`"
  # Not setting a default, so that this must explicitly be set to `null` to specify no tags
  type     = map(string)
  nullable = true
}

variable "provider_tags" {
  type     = map(string)
  nullable = true
  default  = null
}

variable "ignore_tag_keys" {
  type     = set(string)
  nullable = false
}
