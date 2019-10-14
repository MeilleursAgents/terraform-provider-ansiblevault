package vault

import (
	"errors"
	"fmt"
	"log"
	"path"
	"reflect"
	"testing"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

const (
	filesFolder = "../../files/"
)

func TestNew(t *testing.T) {
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

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			result, err := New(testCase.vaultPass, testCase.rootFolder, testCase.separator)

			failed := false

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
				t.Errorf("New(`%s`, `%s`, `%s`) = (%#v, %#v), want (%#v, %#v)", testCase.vaultPass, testCase.rootFolder, testCase.separator, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestGetVaultPass(t *testing.T) {
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
			path.Join(filesFolder, "vault_pass_test.txt"),
			"ansible",
			"secret",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			app, err := New(testCase.vaultPass, testCase.rootFolder, "")
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.getVaultPass()

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
				t.Errorf("getVaultPass() = (`%s`, %#v), want (`%s`, %#v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestGetVaultKey(t *testing.T) {
	if err := ansible_vault.EncryptFile(path.Join(filesFolder, "simple_vault_test.yaml"), "API_KEY:NOT_IN_CLEAR_TEXT", "secret"); err != nil {
		log.Printf("unable to encrypt simple vault for testing: %v", err)
		t.Fail()
	}

	if err := ansible_vault.EncryptFile(path.Join(filesFolder, "complex_vault_test.yaml"), "API_KEY:NOT_IN_CLEAR_TEXT\nTOKEN\nAPI_secret:password\nAPI_complex_secret:test:[!\"\"\n", "secret"); err != nil {
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
			path.Join(filesFolder, "vault_pass_test.txt"),
			"ansible",
			"notExistingFile.txt",
			"api_key",
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle simple vault file with case insensitive comparison",
			path.Join(filesFolder, "vault_pass_test.txt"),
			"./",
			path.Join(filesFolder, "simple_vault_test.yaml"),
			"api_key",
			"NOT_IN_CLEAR_TEXT",
			nil,
		},
		{
			"should handle multi-line vault file",
			path.Join(filesFolder, "vault_pass_test.txt"),
			"./",
			path.Join(filesFolder, "complex_vault_test.yaml"),
			"API_SECRET",
			"password",
			nil,
		},
		{
			"should handle error on not found key",
			path.Join(filesFolder, "vault_pass_test.txt"),
			"./",
			path.Join(filesFolder, "complex_vault_test.yaml"),
			"KEY_NOT_FOUND",
			"",
			ErrKeyNotFound,
		},
		{
			"should handle multi-line vault file with separator in password",
			path.Join(filesFolder, "vault_pass_test.txt"),
			"./",
			path.Join(filesFolder, "complex_vault_test.yaml"),
			"API_complex_secret",
			"test:[!\"\"",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			app, err := New(testCase.vaultPass, testCase.rootFolder, "")
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.getVaultKey(testCase.filename, testCase.key)

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
				t.Errorf("getVaultKey(`%s`, `%s`) = (`%s`, %#v), want (`%s`, %#v)", testCase.filename, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestInEnv(t *testing.T) {
	if err := ansible_vault.EncryptFile(path.Join(filesFolder, "group_vars/tag_prod/vault.yml"), "API_KEY:PROD_KEEP_IT_SECRET", "secret"); err != nil {
		log.Printf("unable to encrypt dev vault for testing: %v", err)
		t.Fail()
	}

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
			"not existing env",
			"dev",
			"API_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(filesFolder, "group_vars/tag_dev/vault.yml")),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New(path.Join(filesFolder, "vault_pass_test.txt"), filesFolder, "")
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.InEnv(testCase.env, testCase.key)

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
				t.Errorf("InEnv(`%s`, `%s`) = (`%s`, %#v), want (`%s`, %#v)", testCase.env, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestInPath(t *testing.T) {
	if err := ansible_vault.EncryptFile(path.Join(filesFolder, "group_vars/tag_prod/vault.yml"), "API_KEY:PROD_KEEP_IT_SECRET", "secret"); err != nil {
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
			"tag_prod/vault.yml",
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"not existing env",
			"not_found.yml",
			"API_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(filesFolder, "group_vars", "not_found.yml")),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New(path.Join(filesFolder, "vault_pass_test.txt"), path.Join(filesFolder, "group_vars"), "")
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.InPath(testCase.path, testCase.key)

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
				t.Errorf("InPath(`%s`, `%s`) = (`%s`, %#v), want (`%s`, %#v)", testCase.path, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	var cases = []struct {
		intention string
		word      string
		want      string
	}{
		{
			"should sanitize quoted string",
			" \"test\" ",
			"test",
		},
		{
			"should sanitize single quoted string",
			" 'test' ",
			"test",
		},
		{
			"should sanitize empty single quoted",
			"''",
			"",
		},
		{
			"should not sanitize single quoted",
			"'",
			"'",
		},
		{
			"should sanitize unquoted string",
			" test",
			"test",
		},
		{
			"should sanitize regular string",
			"test",
			"test",
		},
		{
			"should sanitize integer",
			"11",
			"11",
		},
		{
			"should sanitize unquoted string with inside quote",
			"1\"1",
			"1\"1",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			if result := sanitize(testCase.word); testCase.want != result {
				t.Errorf("sanitize(`%s`) = (`%s`), want (`%s`)", testCase.word, result, testCase.want)
			}
		})
	}
}
