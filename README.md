# terraform-provider-ansiblevault

## Usage

ansiblevault_env example:
-------------------------
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
-------------------------
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