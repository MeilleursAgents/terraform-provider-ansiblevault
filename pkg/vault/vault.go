package vault

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"

	ansible_vault "github.com/sosedoff/ansible-vault-go"
)

const (
	defaultKeySeparator = ":"
)

var (
	// ErrNoVaultPass occurs when vault pass file is blank
	ErrNoVaultPass = errors.New("no vault pass file")

	// ErrNoRootFolder occurs when root folder is blank
	ErrNoRootFolder = errors.New("no root folder")

	// ErrKeyNotFound occurs when key is not found in vault
	ErrKeyNotFound = errors.New("key not found")
)

// App of package
type App struct {
	vaultPass    string
	rootFolder   string
	keySeparator string
}

// New creates new App from Config
func New(vaultPass, rootFolder, keySeparator string) (*App, error) {
	if vaultPass == "" {
		return nil, ErrNoVaultPass
	}

	if rootFolder == "" {
		return nil, ErrNoRootFolder
	}

	if keySeparator == "" {
		keySeparator = defaultKeySeparator
	}

	return &App{
		vaultPass:    vaultPass,
		rootFolder:   rootFolder,
		keySeparator: keySeparator,
	}, nil
}

func (a App) getVaultPass() (string, error) {
	data, err := ioutil.ReadFile(a.vaultPass)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(data), "\n"), nil
}

func (a App) getVaultKey(filename string, key string, getVaultContent func(string, string) (string, error)) (string, error) {
	pass, err := a.getVaultPass()
	if err != nil {
		return "", err
	}

	rawVault, err := getVaultContent(filename, pass)
	if err != nil {
		return "", err
	}

	for _, n := range strings.Split(rawVault, "\n") {
		parts := strings.SplitN(n, a.keySeparator, 2)

		if len(parts) > 1 {
			if strings.EqualFold(parts[0], key) {
				return sanitize(parts[1]), nil
			}
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

func sanitize(word string) string {
	wordTrim := strings.TrimSpace(word)

	s, err := strconv.Unquote(wordTrim)
	if err != nil {
		wordTrimLen := len(wordTrim)

		// we arrive here with integer of non protected string or single quoted string
		if wordTrimLen > 2 && wordTrim[0] == wordTrim[wordTrimLen-1] && wordTrim[0] == byte('\'') {
			return wordTrim[1 : wordTrimLen-1]
		}

		// we arrive here with integer of non protected string
		log.Printf("[WARNING] unable to unquote value: %s", err)
		return wordTrim
	}

	return s
}
