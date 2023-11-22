---
subcategory: "Lex V2 Models"
layout: "aws"
page_title: "AWS: aws_lexv2models_bot_version"
description: |-
  Terraform resource for managing an AWS Lex V2 Models Bot Version.
---

# Resource: aws_lexv2models_bot_version

Terraform resource for managing an AWS Lex V2 Models Bot Version.

## Example Usage

### Basic Usage

```terraform
resource "aws_lexv2models_bot_version" "test" {
  bot_id = aws_lexv2models.test.id
  locale_details = {
    test_key = {
      source_bot_version = "DRAFT"
    }
  }
}
```

## Argument Reference

The following arguments are required:

* `bot_id` - Idientifier of the bot to create the version for.
* `version_locale_specification` - Specifies the locales that Amazon Lex adds to this version. You can choose the draft version or any other previously published version for each locale. When you specify a source version, the locale data is copied from the source version to the new version.

The following arguments are optional:

* `description` - A description of the version. Use the description to help identify the version in lists.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `members` - List of bot members in a network to be created.
* `bot_type` - Type of a bot to create.
* `name` - Name of the bot. The bot name must be unique in the account that creates the bot. Type String. Length Constraints: Minimum length of 1. Maximum length of 100.
* `data_privacy` - Provides information on additional privacy protections Amazon Lex should use with the bot's data.
* `idle_session_ttl_in_seconds` - Time, in seconds, that Amazon Lex should keep information about a user's conversation with the bot. You can specify between 60 (1 minute) and 86,400 (24 hours) seconds.
* `role_arn` - ARN of an IAM role that has permission to access the bot.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `30m`)
* `delete` - (Default `30m`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Lex V2 Models Bot Version using the `example_id_arg`. For example:

```terraform
import {
  to = aws_lexv2models_bot_version.example
  id = "bot_version-id-12345678"
}
```

Using `terraform import`, import Lex V2 Models Bot Version using the `example_id_arg`. For example:

```console
% terraform import aws_lexv2models_bot_version.example bot_version-id-12345678
```