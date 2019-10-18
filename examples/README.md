# Example

# Presentation
- `ansible` contains our ansible structure (vault at root and in group_vars) and the vault pass.
- `terraform` contains files with provider, data and outputs definition.

# Run

```bash
$ terraform init

Initializing the backend...

Initializing provider plugins...

The following providers do not have any version constraints in configuration,
so the latest version was installed.

To prevent automatic upgrades to new major versions that may contain breaking
changes, it is recommended to add version = "..." constraints to the
corresponding provider blocks in configuration, with the constraint strings
suggested below.

* provider.ansiblevault: version = "~> 1.2"

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.

$ terraform apply
data.ansiblevault_path.path: Refreshing state...
data.ansiblevault_env.env: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

env = {
  "env" = "prod"
  "id" = "2019-10-18 11:08:22.906240562 +0000 UTC"
  "key" = "API_KEY"
  "value" = "PROD_KEEP_IT_SECRET"
}
path = {
  "id" = "2019-10-18 11:08:22.906266316 +0000 UTC"
  "key" = "API_KEY"
  "path" = "../ansible/simple_vault_test.yaml"
  "value" = "NOT_IN_CLEAR_TEXT"
}
```