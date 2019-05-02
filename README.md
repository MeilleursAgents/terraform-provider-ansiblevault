# terraform-provider-ansiblevault

## Usage

```tf
provider "ansiblevault" {
  vault_pass  = "~/.vault_pass.txt"
  root_folder = "~/infra/ansible/"
}

data "ansiblevault_env" "api_key" {
  env = "prod"
  key = "SECRET_API_KEY"
}

${data.ansiblevault_env.api_key.value}
```
