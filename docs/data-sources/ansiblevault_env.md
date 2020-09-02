# `ansiblevault_env` Data Source

Use `ansiblevault_env` data source to read in `group_vars/tag_<env>/vault.yml` file the specified `key`.

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

## Argument Reference

The following arguments are supported:

* `env` - (Required) environment targeted for reading.

* `key` - (Required) key to find in yaml.

## Attributes Reference

The following attributes are exported:

* `value` - the content of yaml key.
