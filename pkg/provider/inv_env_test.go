package provider

import (
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/pkg/vault"
)

func TestInEnvRead(t *testing.T) {
	var cases = []struct {
		intention string
		env       string
		key       string
		want      string
		wantErr   error
	}{
		{
			"simple",
			"prod",
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"not found key",
			"prod",
			"SECRET_KEY",
			"",
			errors.New("SECRET_KEY not found in prod vault"),
		},
		{
			"not found env",
			"dev",
			"SECRET_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "group_vars/tag_dev/vault.yml")),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inEnvResource().Data(nil)

			if err := data.Set("env", testCase.env); err != nil {
				t.Errorf("unable to set env: %#v", err)
				return
			}

			data.Set("key", testCase.key)
			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %#v", err)
				return
			}

			vaultApp, err := vault.New(path.Join(ansibleFolder, "vault_pass_test.txt"), ansibleFolder, "")
			if err != nil {
				t.Errorf("unable to create vault app: %#v", err)
				return
			}

			err = inEnvRead(data, vaultApp)
			result := data.Get("value").(string)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			} else if result != testCase.want {
				failed = true
			}

			if failed {
				t.Errorf("InEnvRead(%#v) = (`%s`, %#v), want (`%s`, %#v)", data, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
