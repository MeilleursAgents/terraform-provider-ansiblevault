package provider

import (
	"fmt"
	"time"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
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
				Optional:    true,
			},
			"value": {
				Computed:    true,
				Description: "Vault value found",
				Type:        schema.TypeString,
			},
		},
	}
}

func inEnvRead(data *schema.ResourceData, m interface{}) error {
	env := data.Get("env").(string)
	key := data.Get("key").(string)

	data.SetId(time.Now().UTC().String())

	value, err := m.(*vault.App).InEnv(env, key)
	if err != nil {
		data.SetId("")

		if err == vault.ErrKeyNotFound {
			return fmt.Errorf("%s not found in %s vault", key, env)
		}

		return err
	}

	if err := data.Set("value", value); err != nil {
		data.SetId("")
		return err
	}

	return nil
}
