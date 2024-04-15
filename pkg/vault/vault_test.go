package vault

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"testing"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

const (
	ansibleFolder = "../../examples/ansible/"
)

func TestNew(t *testing.T) {
	var cases = []struct {
		intention  string
		vaultPass  string
		rootFolder string
		want       *App
		wantErr    error
	}{
		{
			"should reject empty root folder",
			"~/.vault_pass.txt",
			"",
			nil,
			ErrNoRootFolder,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			result, err := New(testCase.vaultPass, testCase.rootFolder, "")

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
				t.Errorf("New() = (%#v, %#v), want (%#v, %#v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestGetVaultValueAtPath(t *testing.T) {
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
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			result, err := getVaultValueAtPath(testCase.vaultPass)

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
				t.Errorf("getVaultPass() = (`%s`, %v), want (`%s`, %v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestGetVaultKey(t *testing.T) {
	var cases = []struct {
		intention       string
		vaultPass       string
		rootFolder      string
		input           string
		key             string
		getVaultContent func(string, string) (string, error)
		want            string
		wantErr         error
	}{
		{
			"should handle error while decrypting file",
			"secret",
			"ansible",
			"notExistingFile.txt",
			"api_key",
			ansible_vault.DecryptFile,
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should not handle simple vault file with case insensitive comparison",
			"secret",
			"./",
			path.Join(ansibleFolder, "simple_vault_test.yaml"),
			"api_key",
			ansible_vault.DecryptFile,
			"",
			ErrKeyNotFound,
		},
		{
			"should handle empty key",
			"secret",
			"./",
			path.Join(ansibleFolder, "simple_vault_test.yaml"),
			"",
			ansible_vault.DecryptFile,
			"API_KEY: NOT_IN_CLEAR_TEXT",
			nil,
		},
		{
			"should handle invalid yaml",
			"secret",
			"./",
			path.Join(ansibleFolder, "invalid_yaml_test.yaml"),
			"api_key",
			ansible_vault.DecryptFile,
			"",
			errors.New("yaml: unmarshal errors:\n  line 3: cannot unmarshal !!str `I'm not...` into map[interface {}]interface {}"),
		},
		{
			"should handle multi-line vault file",
			"secret",
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"API_secret",
			ansible_vault.DecryptFile,
			"password",
			nil,
		},
		{
			"should handle multi-line vault value",
			"secret",
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"MULTILINE_token",
			ansible_vault.DecryptFile,
			"foo\nbar",
			nil,
		},
		{
			"should handle error on not found key",
			"secret",
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"KEY_NOT_FOUND",
			ansible_vault.DecryptFile,
			"",
			ErrKeyNotFound,
		},
		{
			"double_quoted string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"double_quoted",
			ansible_vault.DecryptFile,
			"test",
			nil,
		},
		{
			"double_quoted string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"single_quoted",
			ansible_vault.DecryptFile,
			"test",
			nil,
		},
		{
			"unquoted string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"unquoted",
			ansible_vault.DecryptFile,
			"test",
			nil,
		},
		{
			"single_quote string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"single_quote",
			ansible_vault.DecryptFile,
			"'",
			nil,
		},
		{
			"integer string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"integer",
			ansible_vault.DecryptFile,
			"11",
			nil,
		},
		{
			"quote_inside string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"quote_inside",
			ansible_vault.DecryptFile,
			"abc'def",
			nil,
		},
		{
			"double_quote_inside string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"double_quote_inside",
			ansible_vault.DecryptFile,
			"abc\"def",
			nil,
		},
		{
			"boolean string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"boolean",
			ansible_vault.DecryptFile,
			"true",
			nil,
		},
		{
			"nested string",
			"secret",
			"./",
			path.Join(ansibleFolder, "sanitized_vault.yml"),
			"nested.variable",
			ansible_vault.DecryptFile,
			"value",
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

			result, err := app.getVaultKey(testCase.input, testCase.key, testCase.getVaultContent)

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
				t.Errorf("getVaultKey(`%s`, `%s`) = (`%s`, %v), want (`%s`, %v)", testCase.input, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestGetVaultPassword(t *testing.T) {
	var cases = []struct {
		intention string
		vaultPass string
		vaultPath string
		want      string
		wantErr   error
	}{
		{
			"should return vaultPass first",
			"secretValue",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"secretValue",
			nil,
		},
		{
			"should return content of vaultPath",
			"",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"secret",
			nil,
		},
		{
			"should return content of vaultPath",
			"",
			"notExistingFile.txt",
			"",
			ErrNoVaultPass,
		},
		{
			"should handle error when no vaultPass and vaultPath",
			"",
			"",
			"",
			ErrNoVaultPass,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			result, err := GetVaultPassword(testCase.vaultPath, testCase.vaultPass)

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
				t.Errorf("getVaultPass() = (`%s`, %v), want (`%s`, %v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestInEnv(t *testing.T) {
	var cases = []struct {
		intention  string
		pattern    string
		pathParams map[string]interface{}
		key        string
		want       string
		wantErr    error
	}{
		{
			"simple",
			"group_vars/tag_{{.env}}/vault.yml",
			map[string]interface{}{
				"env": "prod",
			},
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"not existing env",
			"group_vars/tag_{{.env}}/vault.yml",
			map[string]interface{}{
				"env": "dev",
			},
			"API_KEY",
			"",
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "group_vars/tag_dev/vault.yml")),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New("secret", ansibleFolder, testCase.pattern)
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.InPathPattern(testCase.pathParams, testCase.key)

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
				t.Errorf("InEnv(`%s`, `%s`) = (`%s`, %#v), want (`%s`, %#v)", testCase.pathParams, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestInPath(t *testing.T) {
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
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "group_vars", "not_found.yml")),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New("secret", path.Join(ansibleFolder, "group_vars"), "")
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

func TestInString(t *testing.T) {
	vaultRaw := `$ANSIBLE_VAULT;1.1;AES256
33623735333733316564643935636565663664376661326536303633366465343631626265303030
3464346366613935623239353334383831323036363236660a366261643665316438623431376135
32636366373330363438613439656261653932653033386132356265323937373733633834643432
6238666665373737620a653565656635373165643936303337646234663133336438343236363662
64646462623864306562623264316535653238656664383661353738623662393137`

	var cases = []struct {
		intention string
		input     string
		key       string
		want      string
		wantErr   error
	}{
		{
			"simple",
			vaultRaw,
			"API_KEY",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
		{
			"invalid format",
			"novaultformat",
			"API_KEY",
			"",
			errors.New("invalid secret format"),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New("secret", path.Join(ansibleFolder, "group_vars"), "")
			if err != nil {
				t.Errorf("unable to create App: %#v", err)
				return
			}

			result, err := app.InString(testCase.input, testCase.key)

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
				t.Errorf("InString() = (`%s`, %s), want (`%s`, %s)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
