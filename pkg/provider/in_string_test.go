package provider

import (
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

func TestInStringRead(t *testing.T) {
	vaultRaw := `$ANSIBLE_VAULT;1.1;AES256
33623735333733316564643935636565663664376661326536303633366465343631626265303030
3464346366613935623239353334383831323036363236660a366261643665316438623431376135
32636366373330363438613439656261653932653033386132356265323937373733633834643432
6238666665373737620a653565656635373165643936303337646234663133336438343236363662
64646462623864306562623264316535653238656664383661353738623662393137`

	vaultRawString := `$ANSIBLE_VAULT;1.1;AES256
66306134666665663135666633346565363436333837376232613938393164353936333863653961
6563396637656665303736336463663332376463616431350a343336306234666665663038393430
66313666666366616565366536366563666135623730303462363430313532356333313734316363
6538313234313665350a313236333731656165303634616635663234636634363264383463386339
34346433386537313665666233626238613763643132346533376634356435323562`

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
			"not found key",
			vaultRaw,
			"SECRET_KEY",
			"",
			vault.ErrKeyNotFound,
		},
		{
			"not provided key",
			vaultRawString,
			"",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inStringResource().Data(nil)

			if err := data.Set("encrypted", testCase.input); err != nil {
				t.Errorf("unable to set encrypted: %s", err)
				return
			}

			if err := data.Set("key", testCase.key); err != nil {
				t.Errorf("unable to set key: %s", err)
				return
			}

			vaultApp, err := vault.New("secret", ansibleFolder, "")
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

func TestInStringEncRead(t *testing.T) {
	var cases = []struct {
		intention string
		input     string
		wantErr   error
	}{
		{
			"simple",
			"PROD_KEEP_IT_SECRET",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			data := inStringResource().Data(nil)

			if err := data.Set("value", testCase.input); err != nil {
				t.Errorf("unable to set raw value: %s", err)
				return
			}

			vaultApp, err := vault.New("secret", ansibleFolder, "")
			if err != nil {
				t.Errorf("unable to create vault app: %#v", err)
				return
			}

			err = inStringEncRead(data, vaultApp)
			result := data.Get("encrypted").(string)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			} else {
				decValue, err := ansible_vault.Decrypt(result, "secret")
				if err != nil || decValue != testCase.input {
					t.Errorf("inStringEncRead() = (`%s`, %#v), want (`%s`, %#v)", result, err, decValue, testCase.wantErr)
					failed = true
				}
			}

			if failed {
				t.Errorf("inStringEncRead() = (`%s`, %#v), want (%#v)", result, err, testCase.wantErr)
			}
		})
	}
}
