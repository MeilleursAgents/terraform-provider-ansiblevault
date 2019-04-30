package vault

import (
	"errors"
	"log"
	"testing"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

func Test_getVaultPass(t *testing.T) {
	var cases = []struct {
		intention string
		app       App
		want      string
		wantErr   error
	}{
		{
			"should handle error while reading",
			App{
				vaultPass: "notExistingFile.txt",
			},
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should sanitize vault pass",
			App{
				vaultPass: "vault_pass_test.txt",
			},
			"secret",
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := testCase.app.getVaultPass()

		failed = false

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
			t.Errorf("%s\ngetVaultPass() = (`%s`, %+v), want (`%s`, %+v)", testCase.intention, result, err, testCase.want, testCase.wantErr)
		}
	}
}

func Test_getVaultKey(t *testing.T) {
	if err := ansible_vault.EncryptFile("simple_vault_test.yaml", "API_KEY:YOU_SHOULD_NOT_SEE_ME", "secret"); err != nil {
		log.Printf("unable to encrypt simple vault for testing: %v", err)
		t.Fail()
	}

	if err := ansible_vault.EncryptFile("complex_vault_test.yaml", "API_KEY:YOU_SHOULD_NOT_SEE_ME\nTOKEN\nAPI_secret:password\n", "secret"); err != nil {
		log.Printf("unable to encrypt complex vault for testing: %v", err)
		t.Fail()
	}

	var cases = []struct {
		intention string
		app       App
		filename  string
		key       string
		want      string
		wantErr   error
	}{
		{
			"should handle error while reading vault",
			App{
				vaultPass: "notExistingFile.txt",
			},
			"",
			"api_key",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle error while decrypting file",
			App{
				vaultPass: "vault_pass_test.txt",
			},
			"notExistingFile.txt",
			"api_key",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle simple vault file with case insensitive comparison",
			App{
				vaultPass: "vault_pass_test.txt",
			},
			"simple_vault_test.yaml",
			"api_key",
			"YOU_SHOULD_NOT_SEE_ME",
			nil,
		},
		{
			"should handle multi-line vault file",
			App{
				vaultPass: "vault_pass_test.txt",
			},
			"complex_vault_test.yaml",
			"API_SECRET",
			"password",
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := testCase.app.getVaultKey(testCase.filename, testCase.key)

		failed = false

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
			t.Errorf("%s\ngetVaultKey(`%s`, `%s`) = (`%s`, %+v), want (`%s`, %+v)", testCase.intention, testCase.filename, testCase.key, result, err, testCase.want, testCase.wantErr)
		}
	}
}
