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
	"github.com/MeilleursAgents/terraform-provider-ansiblevault/pkg/vault"
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider create and returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vault_pass": {
				Type:        schema.TypeString,
				Description: "Ansible vault pass file",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANSIBLE_VAULT_PASS_FILE", nil),
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
			vaultPass := r.Get("vault_pass").(string)
			rootFolder := r.Get("root_folder").(string)

			return vault.New(vaultPass, rootFolder)
		},
	}
}
