package provider

import (
	"fmt"
	"time"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func inPathPatternResource() *schema.Resource {
	return &schema.Resource{
		Read: inPathPatternRead,
		Schema: map[string]*schema.Schema{
			"path_params": {
				Type:        schema.TypeMap,
				Description: "Parameters for path pattern",
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

func inPathPatternRead(data *schema.ResourceData, m interface{}) error {
	pathParams := data.Get("path_params").(map[string]interface{})
	key := data.Get("key").(string)

	data.SetId(time.Now().UTC().String())

	value, err := m.(*vault.App).InPathPattern(pathParams, key)
	if err != nil {
		data.SetId("")

		if err == vault.ErrKeyNotFound {
			return fmt.Errorf("not found in %s vault", key)
		}

		return err
	}

	if err := data.Set("value", value); err != nil {
		data.SetId("")
		return err
	}

	return nil
}
