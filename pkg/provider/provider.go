package provider

/*
	Usage:
	```
	provider "ansiblevault" {
	  vault_pass = "~/.vault_pass.txt"
	  root_folder = "~/infra/ansible/"
	}
	```
*/

import (
	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider create and returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vault_path": {
				Type:        schema.TypeString,
				Description: "Path to ansible vault password file (cf. https://docs.ansible.com/ansible/latest/user_guide/vault.html#providing-vault-passwords)",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANSIBLE_VAULT_PASSWORD_FILE", nil),
			},
			"vault_pass": {
				Type:        schema.TypeString,
				Description: "Ansible vault pass value",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANSIBLE_VAULT_PASS", nil),
			},
			"root_folder": {
				Type:        schema.TypeString,
				Description: "Ansible root directory",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANSIBLE_ROOT_FOLDER", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ansiblevault_env":    inEnvResource(),
			"ansiblevault_path":   inPathResource(),
			"ansiblevault_string": inStringResource(),
		},
		ConfigureFunc: func(r *schema.ResourceData) (interface{}, error) {
			return configure(safeString(r.Get("vault_path")), safeString(r.Get("vault_pass")), safeString(r.Get("root_folder")))
		},
	}
}

func safeString(input interface{}) string {
	switch input.(type) {
	case string:
		return input.(string)
	default:
		return ""
	}
}

func configure(path, pass, rootFolder string) (interface{}, error) {
	pass, err := vault.GetVaultPassword(path, pass)
	if err != nil {
		return nil, err
	}

	return vault.New(pass, rootFolder)
}
