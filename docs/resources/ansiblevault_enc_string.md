# `ansiblevault_enc_string` Resource

Use `ansiblevault_enc_string` resource to encrypt `value` using the provided ansible_vault key into `encrypted`

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

## Argument Reference

The following arguments are supported:

* `value` - (Required) the raw secret as string.

## Attributes Reference

The following attributes are exported:

* `encrypted` - the ansible vault secret.
