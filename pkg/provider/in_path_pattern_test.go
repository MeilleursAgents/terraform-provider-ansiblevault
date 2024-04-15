package provider

import (
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
)

func TestInEnvRead(t *testing.T) {
	var cases = []struct {
		intention  string
		pathParams map[string]interface{}
		key        string
		want       string
		wantErr    error
	}{
		{
			"simple",
			map[string]interface{}{"env": "prod"},
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"not found key",
			map[string]interface{}{"env": "prod"},
			"SECRET_KEY",
			"",
			errors.New("not found in SECRET_KEY vault"),
		},
		{
			"not found env",
			map[string]interface{}{"env": "dev"},
			"SECRET_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "group_vars/tag_dev/vault.yml")),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inPathPatternResource().Data(nil)

			if err := data.Set("path_params", testCase.pathParams); err != nil {
				t.Errorf("unable to set env: %#v", err)
				return
			}

			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %#v", err)
				return
			}

			vaultApp, err := vault.New("secret", ansibleFolder, "/group_vars/tag_{{.env}}/vault.yml")
			if err != nil {
				t.Errorf("unable to create vault app: %#v", err)
				return
			}

			err = inPathPatternRead(data, vaultApp)
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
