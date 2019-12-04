output "path" {
  value = data.ansiblevault_path.path.value
}

output "env" {
  value = data.ansiblevault_env.env.value
}

output "key_string" {
  value = data.ansiblevault_string.key_string.value
}

output "raw_string" {
  value = data.ansiblevault_string.raw_string.value
}
