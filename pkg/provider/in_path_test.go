package provider

import (
	"errors"
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/pkg/vault"
	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

func TestInPathRead(t *testing.T) {
	if err := ansible_vault.EncryptFile(path.Join(filesFolder, "InPathRead.yml"), "API_KEY:PROD_KEEP_IT_SECRET", "secret"); err != nil {
		log.Printf("unable to encrypt dev vault for testing: %v", err)
		t.Fail()
	}

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
			fmt.Errorf("open %s: no such file or directory", path.Join(filesFolder, "InPathReadNotFound.yml")),
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inPathResource().Data(nil)

			if err := data.Set("path", testCase.path); err != nil {
				t.Errorf("unable to set path: %#v", err)
				return
			}

			data.Set("key", testCase.key)
			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %#v", err)
				return
			}

			vaultApp, err := vault.New(path.Join(filesFolder, "vault_pass_test.txt"), filesFolder, "")
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
				t.Errorf("InPathRead(%#v) = (`%s`, %#v), want (`%s`, %#v)", data, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
