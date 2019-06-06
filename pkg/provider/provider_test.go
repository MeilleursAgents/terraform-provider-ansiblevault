package provider

import "testing"

const (
	filesFolder = "../../files/"
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
				t.Errorf("Provider() = %#v, want %#v", result, testCase.want)
			}
		})
	}
}
