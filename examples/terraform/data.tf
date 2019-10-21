data "ansiblevault_path" "path" {
  path = "../ansible/simple_vault_test.yaml"
  key  = "API_KEY"
}

data "ansiblevault_env" "env" {
  env = "prod"
  key = "API_KEY"
}
