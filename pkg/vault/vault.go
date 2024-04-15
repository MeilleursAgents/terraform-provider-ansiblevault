package vault

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
	"gopkg.in/yaml.v2"
)

var (
	// ErrNoVaultPass occurs when vault pass file or value are empty
	ErrNoVaultPass = errors.New("no vault password file or vault pass provided")

	// ErrNoRootFolder occurs when root folder is blank
	ErrNoRootFolder = errors.New("no root folder")

	// ErrKeyNotFound occurs when key is not found in vault
	ErrKeyNotFound = errors.New("key not found")
)

// App of package
type App struct {
	vaultPassword string
	rootFolder    string
	path_template template.Template
}

// New creates new App from Config
func New(vaultPassword string, rootFolder string, path_pattern string) (*App, error) {
	if rootFolder == "" {
		return nil, ErrNoRootFolder
	}

	return &App{
		vaultPassword: vaultPassword,
		rootFolder:    rootFolder,
		path_template: *template.Must(
			template.New("path_pattern").Parse(path_pattern),
		),
	}, nil
}

// GetVaultPassword is a helper for retrieve vault password value
func GetVaultPassword(vaultPath string, vaultPass string) (string, error) {
	if vaultPath == "" && vaultPass == "" {
		return "", ErrNoVaultPass
	}

	if vaultPass != "" {
		return vaultPass, nil
	}

	pass, err := getVaultValueAtPath(vaultPath)
	if err != nil {
		return "", ErrNoVaultPass
	}

	return pass, nil
}

func getVaultValueAtPath(vaultPath string) (string, error) {
	data, err := ioutil.ReadFile(vaultPath)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(data), "\n"), nil
}

func (a App) getVaultKey(filename string, key string, getVaultContent func(string, string) (string, error)) (string, error) {
	rawVault, err := getVaultContent(filename, a.vaultPassword)
	if err != nil {
		return "", err
	}

	// trim of carriage return for easier use
	if len(strings.TrimSpace(key)) == 0 {
		return strings.Trim(rawVault, "\n"), nil
	}

	var vaultContent = make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(rawVault), &vaultContent); err != nil {
		return "", err
	}

	keys := strings.Split(key, ".")
	for _, k := range keys {
		switch v := vaultContent[k].(type) {
		case map[interface{}]interface{}:
			vaultContent = v
		case string:
			return strings.Trim(v, "\n"), nil
		case int:
			return strconv.Itoa(v), nil
		case bool:
			return strconv.FormatBool(v), nil
		}
	}

	return "", ErrKeyNotFound
}

// InPathPattern retrieves given key in environment vault
func (a App) InPathPattern(pathParams map[string]interface{}, key string) (string, error) {
	var buffer bytes.Buffer
	err := a.path_template.Execute(&buffer, pathParams)
	if err != nil {
		return "", err
	}

	return a.getVaultKey(path.Join(a.rootFolder, buffer.String()), key, ansible_vault.DecryptFile)
}

// InPath retrieves given key in vault file
func (a App) InPath(vaultPath string, key string) (string, error) {
	return a.getVaultKey(path.Join(a.rootFolder, vaultPath), key, ansible_vault.DecryptFile)
}

// InString retrieves given key in vault file
func (a App) InString(rawVault string, key string) (string, error) {
	return a.getVaultKey(rawVault, key, ansible_vault.Decrypt)
}

// InString encrypts a string
func (a App) InEncString(rawValue string) (string, error) {
	return ansible_vault.Encrypt(rawValue, a.vaultPassword)
}
