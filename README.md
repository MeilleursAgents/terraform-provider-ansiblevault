# terraform-provider-ansiblevault

[![Build Status](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault.svg?branch=master)](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault)
[![codecov](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault/branch/master/graph/badge.svg)](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault)
[![Go Report Card](https://goreportcard.com/badge/github.com/MeilleursAgents/terraform-provider-ansiblevault)](https://goreportcard.com/report/github.com/MeilleursAgents/terraform-provider-ansiblevault)

This Terraform provider allows you to access secrets from an Ansible Vault from Terraform.

Made with ❤️ by [MeilleursAgents](https://www.meilleursagents.com)

## Thanks

Thanks to [ansible-vault-go](https://github.com/sosedoff/ansible-vault-go) repository for having done the hardest part.

## Installation
### Terraform 0.13
This provider is available in [terraform registry](https://registry.terraform.io/providers/MeilleursAgents/ansiblevault/latest).

For local development follow https://www.terraform.io/docs/commands/cli-config.html#implied-local-mirror-directories

### Terrform 0.12
```bash
curl https://raw.githubusercontent.com/MeilleursAgents/terraform-provider-ansiblevault/master/install.sh | bash
```

## Usage

ansiblevault_path example:

---

```tf
provider "ansiblevault" {
  vault_path  = "/home/username/.vault_pass.txt"
  root_folder = "/home/username/infra/ansible/"
}

data "ansiblevault_path" "api_key" {
  path = "./passwords.yml"
  key = "USER_PASSWORD"
}

${data.ansiblevault_path.api_key.value} will contain value of `USER_PASSWORD` stored in "/home/username/infra/ansible/passwords.yml"
```

More examples in : [examples/terraform/](https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/examples/terraform)

## Using the provider
### Master
- https://github.com/MeilleursAgents/terraform-provider-ansiblevault/tree/master/docs

### Latest version
- https://registry.terraform.io/providers/hashicorp/http/latest/docs


## Contribution

You have to enable [Go modules](https://github.com/golang/go/wiki/Modules) for compiling this project.

Git hooks are availables for avoiding mistakes and ensure code quality, you can install them with:

```bash
make config
```

## Build and Deploy

You need a [Github OAuth Token](https://github.com/settings/tokens/new) for doing a GitHub release and [goreleaser](https://goreleaser.com/)

```bash
make github
```

## License

This project is licensed under the MIT license (see LICENSE file).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault?ref=badge_large)
