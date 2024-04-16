# `ansiblevault_env` Data Source

Use `ansiblevault_path_pattern` data source to read from path_pattern (see provider config) file the specified `key`.

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

## Argument Reference

The following arguments are supported:

* `path_params` - (Required) A map to render the path_pattern. Must contains all keys given in path_pattern

* `key` - (Required) key to find in yaml.

## Attributes Reference

The following attributes are exported:

* `value` - the content of yaml key.
