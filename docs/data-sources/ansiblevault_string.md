# `ansiblevault_string` Data Source

Use `ansiblevault_string` data source to read in `encrypted` raw data give the specified `key`.

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

## Argument Reference

The following arguments are supported:

* `encrypted` - (Required) the raw vault file as string.

* `key` - (Required) key to find in yaml.

## Attributes Reference

The following attributes are exported:

* `value` - the content of yaml key.
