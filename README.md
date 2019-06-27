# terraform-provider-ansiblevault

[![Build Status](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault.svg?branch=master)](https://travis-ci.org/MeilleursAgents/terraform-provider-ansiblevault)
[![codecov](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault/branch/master/graph/badge.svg)](https://codecov.io/gh/MeilleursAgents/terraform-provider-ansiblevault)
[![Go Report Card](https://goreportcard.com/badge/github.com/MeilleursAgents/terraform-provider-ansiblevault)](https://goreportcard.com/report/github.com/MeilleursAgents/terraform-provider-ansiblevault)

This Terraform provider allows you to access secrets from an Ansible Vault from Terraform.

Made with ❤️ by [MeilleursAgents](https://www.meilleursagents.com)

## Thanks

Thanks to [ansible-vault-go](https://github.com/sosedoff/ansible-vault-go) repository for having done the hardest part.

## Installation

### Golang Way

If you have [Golang installed](https://golang.org/dl/)

```bash
go install github.com/MeilleursAgents/terraform-provider-ansiblevault
mkdir -p "~/.terraform.d/plugins/$(go env GOHOSTOS)_$(go env GOHOSTARCH)/"
cp ${GOPATH}/bin/terraform-provider-ansiblevault "~/.terraform.d/plugins/$(go env GOHOSTOS)_$(go env GOHOSTARCH)/"
```

or

In the repository
```bash
make build
make install
```

### Without golang

```bash
PLUGIN_VERSION="1.0.1"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')

if [[ "${ARCH}" = "x86_64" ]]; then
  ARCH="amd64"
fi

pushd $HOME
mkdir -p ".terraform.d/plugins/${OS}_${ARCH}/"
cd ".terraform.d/plugins/${OS}_${ARCH}/"
curl -o "terraform-provider-ansiblevault_v${PLUGIN_VERSION}" "https://github.com/MeilleursAgents/terraform-provider-ansiblevault/releases/download/v${PLUGIN_VERSION}/terraform-provider-ansiblevault_${OS}-${ARCH}_v${PLUGIN_VERSION}"
popd
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

${data.ansiblevault_env.api_key.value}
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

${data.ansiblevault_path.api_key.value}
```

## Contribution

You have to enable [Go modules](https://github.com/golang/go/wiki/Modules) for compiling this project.

## License

This project is licensed under the MIT license (see LICENSE file).

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FMeilleursAgents%2Fterraform-provider-ansiblevault?ref=badge_large)
