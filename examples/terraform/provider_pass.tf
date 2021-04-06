variable "vault_pass" {
  type    = string
  default = "secret"
}

#  See https://github.com/MeilleursAgents/terraform-provider-ansiblevault/blob/master/README.md for installation and usage
provider "ansiblevault" {
  vault_pass  = var.vault_pass
  root_folder = "../ansible"
}

data "ansiblevault_path" "path" {
  path = "../ansible/simple_vault_test.yaml"
  key  = "API_KEY"
}

data "ansiblevault_env" "env" {
  env = "prod"
  key = "API_KEY"
}

data "ansiblevault_string" "key_string" {
  encrypted = <<EOF
$ANSIBLE_VAULT;1.1;AES256
33623735333733316564643935636565663664376661326536303633366465343631626265303030
3464346366613935623239353334383831323036363236660a366261643665316438623431376135
32636366373330363438613439656261653932653033386132356265323937373733633834643432
6238666665373737620a653565656635373165643936303337646234663133336438343236363662
64646462623864306562623264316535653238656664383661353738623662393137
EOF
  key       = "API_KEY"
}

output "path" {
  value = data.ansiblevault_path.path.value
}

output "env" {
  value = data.ansiblevault_env.env.value
}

output "key_string" {
  value = data.ansiblevault_string.key_string.value
}
