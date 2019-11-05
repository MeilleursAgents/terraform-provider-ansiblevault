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
			path.Join(ansibleFolder, "vault_pass_test.txt"),
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
			"should handle error while reading vault",
			"notExistingFile.txt",
			"ansible",
			"",
			"api_key",
			ansible_vault.DecryptFile,
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle error while decrypting file",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"ansible",
			"notExistingFile.txt",
			"api_key",
			ansible_vault.DecryptFile,
			"",
			errors.New("open notExistingFile.txt: no such file or directory"),
		},
		{
			"should handle simple vault file with case insensitive comparison",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"./",
			path.Join(ansibleFolder, "simple_vault_test.yaml"),
			"api_key",
			ansible_vault.DecryptFile,
			"NOT_IN_CLEAR_TEXT",
			nil,
		},
		{
			"should handle multi-line vault file",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"API_SECRET",
			ansible_vault.DecryptFile,
			"password",
			nil,
		},
		{
			"should handle error on not found key",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"KEY_NOT_FOUND",
			ansible_vault.DecryptFile,
			"",
			ErrKeyNotFound,
		},
		{
			"should handle multi-line vault file with separator in password",
			path.Join(ansibleFolder, "vault_pass_test.txt"),
			"./",
			path.Join(ansibleFolder, "complex_vault_test.yaml"),
			"API_complex_secret",
			ansible_vault.DecryptFile,
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
				t.Errorf("getVaultKey(`%s`, `%s`) = (`%s`, %#v), want (`%s`, %#v)", testCase.input, testCase.key, result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}

func TestInEnv(t *testing.T) {
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
			fmt.Errorf("open %s: no such file or directory", path.Join(ansibleFolder, "group_vars/tag_dev/vault.yml")),
		},
	}

	var failed bool

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {

			app, err := New(path.Join(ansibleFolder, "vault_pass_test.txt"), ansibleFolder, "")
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

			app, err := New(path.Join(ansibleFolder, "vault_pass_test.txt"), path.Join(ansibleFolder, "group_vars"), "")
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
61336365316161396566653134393964613564646439313333666233356463336131336537303633
6239626439383636346130653132326138313437306365310a663961653131373535633431393836
34353035376531643266383736306338333764373837656131323663396435666332343039666465
3635613231313833650a346365623861663638313830616564623663386137303735356639313163
34343639636161656230363030353763623830653838333166623234326334663338`

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
			"NOT_IN_CLEAR_TEXT",
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

			app, err := New(path.Join(ansibleFolder, "vault_pass_test.txt"), path.Join(ansibleFolder, "group_vars"), "")
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
