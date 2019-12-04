variable "vault_pass" {
  type    = string
  default = "secret"
}

#  See https://github.com/MeilleursAgents/terraform-provider-ansiblevault/blob/master/README.md for installation and usage
provider ansiblevault {
  version     = "~> 2.0"
  vault_pass  = var.vault_pass
  root_folder = "../ansible"
}

provider ansiblevault {
  version     = "~> 2.0"
  alias       = "password_file"
  vault_path  = "../ansible/vault_pass_test.txt"
  root_folder = "../ansible"
}
