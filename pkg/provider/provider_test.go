package provider

import (
	"reflect"
	"testing"

	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/vault"
)

const (
	ansibleFolder = "../../examples/ansible/"
)

func TestProvider(t *testing.T) {
	var cases = []struct {
		intention string
		want      error
	}{
		{
			"configuration is valid",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			if result := Provider().InternalValidate(); result != testCase.want {
				t.Errorf("Provider() = %s, want %s", result, testCase.want)
			}
		})
	}
}

func TestConfigure(t *testing.T) {
	validVault, _ := vault.New("secret", "../../examples/ansible", "")

	var cases = []struct {
		intention    string
		path       	 string
		path_pattern string
		pass         string
		rootFolder   string
		want         interface{}
		wantErr      error
	}{
		{
			"erroneous password",
			"",
			"",
			"",
			"",
			nil,
			vault.ErrNoVaultPass,
		},
		{
			"erroneous password",
			"",
			"",
			"secret",
			"../../examples/ansible",
			validVault,
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			result, err := configure(testCase.path, testCase.path_pattern, testCase.pass, testCase.rootFolder)

			failed := false

			if testCase.wantErr == nil && err != nil {
				failed = true
			} else if testCase.wantErr != nil && err == nil {
				failed = true
			} else if testCase.wantErr != nil && testCase.wantErr.Error() != err.Error() {
				failed = true
			} else if !reflect.DeepEqual(result, testCase.want) {
				failed = true
			}

			if failed {
				t.Errorf("Configure() = (%#v, %v), want (%#v, %v)", result, err, testCase.want, testCase.wantErr)
			}
		})
	}
}
