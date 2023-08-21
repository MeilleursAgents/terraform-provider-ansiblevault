package provider

import (
	"time"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func inStringResource() *schema.Resource {
	return &schema.Resource{
		Read: inStringRead,
		Schema: map[string]*schema.Schema{
			"encrypted": {
				Type:        schema.TypeString,
				Description: "Ansible-vault string representation",
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

func inStringEncResource() *schema.Resource {
	return &schema.Resource{
		Read:   inStringEncRead,
		Create: inStringEncRead,
		Delete: inStringEncDelete,
		Schema: map[string]*schema.Schema{
			"value": {
				Required:    true,
				ForceNew:    true,
				Description: "Vault value found",
				Type:        schema.TypeString,
			},
			"encrypted": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Ansible-vault string representation",
			},
		},
	}
}

func inStringRead(data *schema.ResourceData, m interface{}) error {
	raw := data.Get("encrypted").(string)
	key := data.Get("key").(string)

	data.SetId(time.Now().UTC().String())

	value, err := m.(*vault.App).InString(raw, key)

	if err != nil {
		data.SetId("")
		return err
	}

	if err := data.Set("value", value); err != nil {
		data.SetId("")
		return err
	}

	return nil
}

func inStringEncRead(data *schema.ResourceData, m interface{}) error {
	value := data.Get("value").(string)
	enc := data.Get("encrypted").(string)

	if len(enc) != 0 {
		dec, err := m.(*vault.App).InString(enc, "")
		// If there is an error, we need to update it
		if err == nil {
			if dec == value {
				return nil
			}
		}
	}

	data.SetId(time.Now().UTC().String())

	encrypted, err := m.(*vault.App).InEncString(value)

	if err != nil {
		data.SetId("")
		return err
	}

	if err := data.Set("encrypted", encrypted); err != nil {
		data.SetId("")
		return err
	}

	return nil
}

func inStringEncDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
