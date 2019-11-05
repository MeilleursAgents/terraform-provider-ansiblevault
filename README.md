# terraform-provider-ansiblevault

[![Build Status](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault.svg?branch=master)](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault)
[![codecov](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault/branch/master/graph/badge.svg)](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault)
[![Go Report Card](https://goreportcard.com/badge/github.com/MeilleursAgents/terraform-provider-ansiblevault)](https://goreportcard.com/report/github.com/MeilleursAgents/terraform-provider-ansiblevault)

This Terraform provider allows you to access secrets from an Ansible Vault from Terraform.

Made with ❤️ by [MeilleursAgents](https://www.meilleursagents.com)

## Thanks

Thanks to [ansible-vault-go](https://github.com/sosedoff/ansible-vault-go) repository for having done the hardest part.

## Installation

```bash
curl https://raw.githubusercontent.com/MeilleursAgents/terraform-provider-ansiblevault/master/install.sh | bash
```

## Usage

ansiblevault_env example:

---

```tf
provider "ansiblevault" {
  vault_pass  = "/home/username/.vault_pass.txt"
  root_folder = "/home/username/infra/ansible/"
}

data "ansiblevault_env" "api_key" {
  env = "prod"
  key = "SECRET_API_KEY"
}

${data.ansiblevault_env.api_key.value} will contain value of `SECRET_API_KEY` stored in "/home/username/infra/ansible/group_vars/tag_prod/vault.yaml"
```

ansiblevault_path example:

---

```tf
provider "ansiblevault" {
  vault_pass  = "/home/username/.vault_pass.txt"
  root_folder = "/home/username/infra/ansible/"
}

data "ansiblevault_path" "api_key" {
  path = "./passwords.yml"
  key = "USER_PASSWORD"
}

${data.ansiblevault_path.api_key.value} will contain value of `USER_PASSWORD` stored in "/home/username/infra/ansible/passwords.yml"
```

ansiblevault_string example:

---

```tf
provider "ansiblevault" {
  vault_pass  = "/home/username/.vault_pass.txt"
  root_folder = "/home/username/infra/ansible/"
}

data "ansiblevault_string" "api_key" {
  encrypted = <<EOF
$ANSIBLE_VAULT;1.1;AES256
65346463633165666232636636346631626565616132653339343961656336643930323937313231
3436383237633937636435636366386563313233366630380a316535376661653933373836633130
30336130396635363830373135643261346437366235303463643538336561356534666161353233
6133626433333965320a323966396162656332386265306539666436643033653466636335363363
35656432663266353133623834653735656534346639623233623531363332373461
EOF
  key = "API_KEY"
}

${data.ansiblevault_string.api_key.value} will contain value of `API_KEY` pass in argument vault string.
```

## Documentation

### Provider

| Key | Required | EnvVar | Description |
|:--:|:--:|:--:|:--:|
| vault_pass | ✅ | `ANSIBLE_VAULT_PASS_FILE` | Ansible vault pass file |
| root_folder | ✅ | `ANSIBLE_ROOT_FOLDER` | Ansible root directory |
| key_separator | | | Separator of key/value pair in Ansible vault (default `:`) |

For an easy way to configure provider with environment variables, consider the following snippet:

```bash
VAULT_PASS="$(ansible-config dump | grep DEFAULT_VAULT_PASSWORD_FILE | awk '{print $3}')"

cat >> "${HOME}/.localrc" << EOM
export ANSIBLE_VAULT_PASS_FILE="${VAULT_PASS}"
export ANSIBLE_ROOT_FOLDER="/path/to/my/ansible/"
EOM
```

## Contribution

You have to enable [Go modules](https://github.com/golang/go/wiki/Modules) for compiling this project.

Git hooks are availables for avoiding mistakes and ensure code quality, you can install them with:

```bash
make config
```

## Build and Deploy

You need a [Github OAuth Token](https://github.com/settings/tokens/new) for doing a GitHub release.

```bash
make github
```

## License

This project is licensed under the MIT license (see LICENSE file).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault?ref=badge_large)
