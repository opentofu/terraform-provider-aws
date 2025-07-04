# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "null" {}

resource "aws_ssmcontacts_rotation" "test" {
  contact_ids = aws_ssmcontacts_contact.test[*].arn

  name = var.rName

  recurrence {
    number_of_on_calls    = 1
    recurrence_multiplier = 1
    daily_settings {
      hour_of_day    = 18
      minute_of_hour = 00
    }
  }

  tags = {
    (var.unknownTagKey) = null_resource.test.id
    (var.knownTagKey)   = var.knownTagValue
  }

  time_zone_id = "Australia/Sydney"

  depends_on = [aws_ssmincidents_replication_set.test]
}

# testAccRotationConfig_base(rName, 1)

resource "aws_ssmcontacts_contact" "test" {
  count = 1
  alias = "${var.rName}-${count.index}"
  type  = "PERSONAL"

  depends_on = [aws_ssmincidents_replication_set.test]
}

# testAccRotationConfig_replicationSetBase

resource "aws_ssmincidents_replication_set" "test" {
  region {
    name = data.aws_region.current.region
  }
}

data "aws_region" "current" {}

resource "null_resource" "test" {}

variable "rName" {
  description = "Name for resource"
  type        = string
  nullable    = false
}

variable "unknownTagKey" {
  type     = string
  nullable = false
}

variable "knownTagKey" {
  type     = string
  nullable = false
}

variable "knownTagValue" {
  type     = string
  nullable = false
}
