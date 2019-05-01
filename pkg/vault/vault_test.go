package vault

import (
	"errors"
	"log"
	"reflect"
	"testing"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

func Test_New(t *testing.T) {
	var cases = []struct {
		intention  string
		vaultPass  string
		rootFolder string
		separator  string
		want       *App
		wantErr    error
	}{
		{
			"should reject empty vault pass",
			"",
			"",
			"",
			nil,
			ErrNoVaultPass,
		},
		{
			"should reject empty root folder",
			"~/.vault_pass.txt",
			"",
			"",
			nil,
			ErrNoRootFolder,
		},
		{
			"should have default value for key separator",
			"~/.vault_pass.txt",
			"ansible",
			"",
			&App{
				vaultPass:    "~/.vault_pass.txt",
				rootFolder:   "ansible",
				keySeparator: defaultKeySeparator,
			},
			nil,
		},
		{
			"should store given key separator",
			"~/.vault_pass.txt",
			"ansible",
			"=",
			&App{
				vaultPass:    "~/.vault_pass.txt",
				rootFolder:   "ansible",
				keySeparator: "=",
			},
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		result, err := New(testCase.vaultPass, testCase.rootFolder, testCase.separator)

		failed = false

		if err == nil && testCase.wantErr != nil {
			failed = true
		} else if err != nil && testCase.wantErr == nil {
			failed = true
		} else if err != nil && err.Error() != testCase.wantErr.Error() {
			failed = true
		} else if !reflect.DeepEqual(result, testCase.want) {
			failed = true
		}

		if failed {
			t.Errorf("%s\nNew(`%s`, `%s`, `%s`) = (%+v, %+v), want (%+v, %+v)", testCase.intention, testCase.vaultPass, testCase.rootFolder, testCase.separator, result, err, testCase.want, testCase.wantErr)
		}
	}
}

func Test_getVaultPass(t *testing.T) {
	var cases = []struct {
		intention  string
		vaultPass  string
		rootFolder string
		want       string
		wantErr    error
	}{
		{
			"should handle error while reading",
			"notExistingFile.txt",
			"ansible",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should sanitize vault pass",
			"vault_pass_test.txt",
			"ansible",
			"secret",
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		app, err := New(testCase.vaultPass, testCase.rootFolder, "")
		if err != nil {
			t.Errorf("%s\nunable to create App: %+v", testCase.intention, err)
		}

		result, err := app.getVaultPass()

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
	if err := ansible_vault.EncryptFile("simple_vault_test.yaml", "API_KEY:NOT_IN_CLEAR_TEXT", "secret"); err != nil {
		log.Printf("unable to encrypt simple vault for testing: %v", err)
		t.Fail()
	}

	if err := ansible_vault.EncryptFile("complex_vault_test.yaml", "API_KEY:NOT_IN_CLEAR_TEXT\nTOKEN\nAPI_secret:password\n", "secret"); err != nil {
		log.Printf("unable to encrypt complex vault for testing: %v", err)
		t.Fail()
	}

	var cases = []struct {
		intention  string
		vaultPass  string
		rootFolder string
		filename   string
		key        string
		want       string
		wantErr    error
	}{
		{
			"should handle error while reading vault",
			"notExistingFile.txt",
			"ansible",
			"",
			"api_key",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle error while decrypting file",
			"vault_pass_test.txt",
			"ansible",
			"notExistingFile.txt",
			"api_key",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle simple vault file with case insensitive comparison",
			"vault_pass_test.txt",
			"./",
			"simple_vault_test.yaml",
			"api_key",
			"NOT_IN_CLEAR_TEXT",
			nil,
		},
		{
			"should handle multi-line vault file",
			"vault_pass_test.txt",
			"./",
			"complex_vault_test.yaml",
			"API_SECRET",
			"password",
			nil,
		},
	}

	var failed bool

	for _, testCase := range cases {
		app, err := New(testCase.vaultPass, testCase.rootFolder, "")
		if err != nil {
			t.Errorf("%s\nunable to create App: %+v", testCase.intention, err)
		}

		result, err := app.getVaultKey(testCase.filename, testCase.key)

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
