package provider

import (
	"path"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/pkg/vault"
)

func TestInStringRead(t *testing.T) {
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
			"not found key",
			vaultRaw,
			"SECRET_KEY",
			"",
			vault.ErrKeyNotFound,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inStringResource().Data(nil)

			if err := data.Set("encrypted", testCase.input); err != nil {
				t.Errorf("unable to set encrypted: %s", err)
				return
			}

			data.Set("key", testCase.key)
			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %s", err)
				return
			}

			vaultApp, err := vault.New(path.Join(ansibleFolder, "vault_pass_test.txt"), ansibleFolder, "")
			if err != nil {
				t.Errorf("unable to create vault app: %#v", err)
				return
			}

			err = inStringRead(data, vaultApp)
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
				t.Errorf("inStringRead() = (`%s`, %#v), want (`%s`, %#v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
