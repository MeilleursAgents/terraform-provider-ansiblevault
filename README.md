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
export ANSIBLE_ROOT_FOLDER="$(pwd)/../path/to/my/ansible/"
EOM
```

## Build and Deploy

You need following environment variables for doing a release. If not set, the `release` script will ask you.

| Name | Description |
|:--:|:--:|
| GITHUB_OAUTH_TOKEN | A Github Token with `repos` access (you can generate it [here](https://github.com/settings/tokens/new)) |
| GITHUB_REPOSITORY | The repository name for uploading assets (e.g. MeilleursAgents/terraform-provider-ansiblevault) |
| GIT_TAG | The new version to release (e.g. v1.0.0) |
| RELEASE_NAME | The version name (most of the time, the git tag) |

```bash
git tag "${GIT_TAG}"
GITHUB_REPOSITORY=MeilleursAgents/terraform-provider-ansiblevault ./script/release
```

## Contribution

You have to enable [Go modules](https://github.com/golang/go/wiki/Modules) for compiling this project.

## License

This project is licensed under the MIT license (see LICENSE file).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault?ref=badge_large)
