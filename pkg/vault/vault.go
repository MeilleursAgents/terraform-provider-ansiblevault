package vault

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

var (
	// ErrKeyNotFound occurs when key is not found in vault
	ErrKeyNotFound = errors.New("key not found")
)

// Config of package
type Config struct {
	vaultPass  *string
	rootFolder *string
}

// App of package
type App struct {
	vaultPass  string
	rootFolder string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet) Config {
	return Config{
		vaultPass:  fs.String("vaultPassFile", "~/.vault_pass.txt", "Vault pass file"),
		rootFolder: fs.String("rootFolder", "", "Ansible root directory"),
	}
}

// New creates new App from Config
func New(config Config) *App {
	return &App{
		vaultPass:  strings.TrimSpace(*config.vaultPass),
		rootFolder: strings.TrimSpace(*config.rootFolder),
	}
}

func (a App) getVaultPass() (string, error) {
	data, err := ioutil.ReadFile(a.vaultPass)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(data), "\n"), nil
}

func (a App) getVaultKey(filename string, key string) (string, error) {
	pass, err := a.getVaultPass()
	if err != nil {
		return "", err
	}

	rawVault, err := ansible_vault.DecryptFile(filename, pass)
	if err != nil {
		return "", err
	}

	for _, n := range strings.Split(rawVault, "\n") {
		parts := strings.Split(n, ":")

		if len(parts) > 1 {
			if strings.EqualFold(parts[0], key) {
				return parts[1], nil
			}
		}
	}

	return "", ErrKeyNotFound
}

// InEnv retrieves given key in environment vault
func (a App) InEnv(env string, key string) (string, error) {
	return a.getVaultKey(path.Join(a.rootFolder, fmt.Sprintf("ansible/group_vars/tag_%s/vault.yml", env)), key)
}
