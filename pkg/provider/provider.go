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
			"key_separator": {
				Type:        schema.TypeString,
				Description: "Separator of key/value pair in Ansible vault",
				Optional:    true,
				Default:     ":",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ansiblevault_env": inEnvResource(),
		},
		ConfigureFunc: func(r *schema.ResourceData) (interface{}, error) {
			vaultPass := r.Get("vault_pass").(string)
			rootFolder := r.Get("root_folder").(string)
			keySeparator := r.Get("key_separator").(string)

			return vault.New(vaultPass, rootFolder, keySeparator)
		},
	}
}
