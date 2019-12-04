package vault

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
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
}

// New creates new App from Config
func New(vaultPassword string, rootFolder string) (*App, error) {
	if rootFolder == "" {
		return nil, ErrNoRootFolder
	}

	return &App{
		vaultPassword: vaultPassword,
		rootFolder:    rootFolder,
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

	var vaultContent map[string]string
	if err := yaml.Unmarshal([]byte(rawVault), &vaultContent); err != nil {
		return "", err
	}

	for vaultKey, vaultValue := range vaultContent {
		if strings.EqualFold(key, vaultKey) {
			return strings.Trim(vaultValue, "\n"), nil
		}
	}

	return "", ErrKeyNotFound
}

// InEnv retrieves given key in environment vault
func (a App) InEnv(env string, key string) (string, error) {
	return a.getVaultKey(path.Join(a.rootFolder, fmt.Sprintf("group_vars/tag_%s/vault.yml", env)), key, ansible_vault.DecryptFile)
}

// InPath retrieves given key in vault file
func (a App) InPath(vaultPath string, key string) (string, error) {
	return a.getVaultKey(path.Join(a.rootFolder, vaultPath), key, ansible_vault.DecryptFile)
}

// InString retrieves given key in vault file
func (a App) InString(rawVault string, key string) (string, error) {
	return a.getVaultKey(rawVault, key, ansible_vault.Decrypt)
}
