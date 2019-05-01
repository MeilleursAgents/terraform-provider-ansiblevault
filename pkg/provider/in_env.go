package provider

import (
	"github.com/MeilleursAgents/terraform-provider-ansiblevault/pkg/vault"
	"github.com/hashicorp/terraform/helper/schema"
)

func inEnvResource() *schema.Resource {
	return &schema.Resource{
		Read: inEnvRead,
		Schema: map[string]*schema.Schema{
			"env": {
				Type:        schema.TypeString,
				Description: "Ansible environment searched",
				Required:    true,
			},
			"key": {
				Type:        schema.TypeString,
				Description: "Vault key searched",
				Required:    true,
			},
			"value": {
				Computed:    true,
				Description: "Vault value found",
				Type:        schema.TypeInt,
			},
		},
	}
}

func inEnvRead(data *schema.ResourceData, m interface{}) error {
	env := data.Get("env").(string)
	key := data.Get("key").(string)

	value, err := m.(vault.App).InEnv(env, key)
	if err != nil {
		return err
	}

	if err := data.Set("value", value); err != nil {
		return nil
	}

	return nil
}
