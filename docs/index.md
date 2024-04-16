# AnsibleVault Terraform Provider

The Ansible vault Terraform Provider is used to read data formatted as yaml and encrypted with [ansible-vault](https://docs.ansible.com/ansible/latest/user_guide/vault.html)

## Example Usage

See [examples](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples) directory

#### Argument Reference

| Key | Required | EnvVar | Description |
|:--:|:--:|:--:|:--:|
| vault_path |  | `ANSIBLE_VAULT_PASSWORD_FILE` | Path to ansible vault password file |
| path_pattern |  | `ANSIBLE_VAULT_PATH_PATTERN` | Vault file path pattern to be used by ansiblevault_path_pattern resources (example: /group_vars/{{.env}}/vault.yml) |
| vault_pass |  | `ANSIBLE_VAULT_PASS` | Ansible vault pass value |
| root_folder | ✅ | `ANSIBLE_ROOT_FOLDER` | Ansible root directory |

For an easy way to configure provider with environment variables, consider the following snippet:

```bash
VAULT_PASS="$(ansible-config dump | grep DEFAULT_VAULT_PASSWORD_FILE | awk '{print $3}')"

cat >> "${HOME}/.localrc" << EOM
export ANSIBLE_VAULT_PASSWORD_FILE="${VAULT_PASS}"
export ANSIBLE_ROOT_FOLDER="/path/to/my/ansible/"
EOM
```

:information_source: `vault_pass` will override `vault_path`
