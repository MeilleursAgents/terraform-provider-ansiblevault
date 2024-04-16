package provider

import (
	"errors"
	"fmt"
	"path"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
)

func TestInPathRead(t *testing.T) {
	var cases = []struct {
		intention string
		path      string
		key       string
		want      string
		wantErr   error
	}{
		{
			"simple",
			"InPathRead.yml",
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"not found key",
			"InPathRead.yml",
			"SECRET_KEY",
			"",
			errors.New("SECRET_KEY not found in InPathRead.yml vault"),
		},
		{
			"not found path",
			"InPathReadNotFound.yml",
			"SECRET_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "InPathReadNotFound.yml")),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inPathResource().Data(nil)

			if err := data.Set("path", testCase.path); err != nil {
				t.Errorf("unable to set path: %#v", err)
				return
			}

			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %#v", err)
				return
			}

			vaultApp, err := vault.New("secret", ansibleFolder, "")
			if err != nil {
				t.Errorf("unable to create vault app: %#v", err)
				return
			}

			err = inPathRead(data, vaultApp)
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
				t.Errorf("InPathRead() = (`%s`, %#v), want (`%s`, %#v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
