#  See https://github.com/MeilleursAgents/terraform-provider-ansiblevault/blob/master/README.md for installation and usage
provider ansiblevault {
  vault_pass  = "../ansible/vault_pass_test.txt"
  root_folder = "../ansible"
}