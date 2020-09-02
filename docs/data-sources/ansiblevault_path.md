# `ansiblevault_path` Data Source

Use `ansiblevault_path` data source to read in `path` the specified `key`.

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

## Argument Reference

The following arguments are supported:

* `path` - (Required) the relative path to the vault file.

* `key` - (Required) key to find in yaml.

## Attributes Reference

The following attributes are exported:

* `value` - the content of yaml key.
