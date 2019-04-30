package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MeilleursAgents/terraform-provider-ansible-vault/pkg/vault"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("terraform-provider-ansible-vault", flag.ExitOnError)

	ansibleFlags := vault.Flags(fs)

	env := fs.String("environment", "", "Vault environment")
	key := fs.String("key", "", "Key searched")

	err := ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
		ff.WithEnvVarPrefix("TERRAFORM_PROVIDER_ANSIBLE_VAULT"),
	)

	if err != nil {
		log.Fatalf("unable to parse flags: %+v", err)
	}

	vaultApp := vault.New(ansibleFlags)
	output, err := vaultApp.InEnv(*env, *key)
	if err != nil {
		log.Fatalf("unable to find in env: %+v", err)
	}

	fmt.Printf("%s", output)
}
