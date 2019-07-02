package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func build(date, tag string) ([]string, error) {
	outputs := []string{}
	for goos, goarchs := range targets {
		for _, goarch := range goarchs {
			log.Println(fmt.Sprintf("==> Build on %s for %s ...", goos, goarch))
			target := fmt.Sprintf(path.Join(buildDirectory, "terraform-provider-ansiblevault_%s-%s_%s"), goos, goarch, tag)
			_, err := executeCmd(
				[]string{
					fmt.Sprintf("GOARCH=%s", goarch),
					fmt.Sprintf("GOOS=%s", goos),
					"CGO_ENABLED=0",
				},
				"go",
				"build",
				"-ldflags",
				fmt.Sprintf("-s -w -X %s.Version=%s -X %s.BuildDate=%s -installsuffix nocgo", gitRepoURL, tag, gitRepoURL, date),
				"-o",
				target,
				".",
			)

			if err != nil {
				return nil, err
			}

			outputs = append(outputs, target)
		}
	}

	return outputs, nil
}

func executeCmd(environs []string, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	cmd.Env = append(os.Environ(), environs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", errors.Wrapf(err, "%s", output)
	}

	return strings.TrimSpace(string(output)), nil
}
