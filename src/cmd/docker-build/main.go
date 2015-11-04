package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	registry string
)

func main() {

	var verbose bool

	flag.StringVar(&registry, "registry-base", "docker.io/", "the registry name to prepend to each Docker image")
	flag.BoolVar(&verbose, "verbose", false, "enables verbose output")
	flag.Parse()

	repo, path, tag := detectRepoPathAndTag()
	dockerfiles := findDockerfiles()

	var err error
	for dockerfile, image := range mapDockerfileToRepo(repo, path, tag, dockerfiles...) {

		if verbose {
			err = runCmdAndPipeOutput("docker", "build", "-t", image, "-f", dockerfile, ".")
		} else {
			_, err = runCmdAndGetOutput("docker", "build", "-t", image, "-f", dockerfile, ".")
		}

		if err != nil {
			log.Println("error building", image, "from", dockerfile)
			log.Fatal(err)
		}

		if verbose {
			err = runCmdAndPipeOutput("gcloud", "docker", "push", image)
		} else {
			_, err = runCmdAndGetOutput("gcloud", "docker", "push", image)
		}

		if err != nil {
			log.Println("error pushing image", image)
			log.Fatal(err)
		}
	}
}

func detectRepoPathAndTag() (repo, path, tag string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working dir", err)
	}

	repo, err = runCmdAndGetOutput("git", "config", "--get", "remote.origin.url")
	if err != nil {
		log.Fatal("error detecting git repo", err)
	}

	if index := strings.LastIndex(repo, ":"); index != -1 {
		repo = repo[index+1:]
	}

	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-4]
	}

	path, err = runCmdAndGetOutput("git", "rev-parse", "--show-toplevel")

	switch {
	case wd == path:
		path = ""
	case strings.HasPrefix(wd, path):
		path = wd[len(path)+1:]
		path = strings.Replace(path, "/", "-", -1)
	default:
		log.Fatal("Current directory is not child of top level", wd, path)
	}

	tag, err = runCmdAndGetOutput("git", "rev-parse", "HEAD")
	if err != nil {
		log.Fatal("error detecting git HEAD revision", err)
	}

	if len(tag) > 7 {
		tag = tag[:7]
	}

	return repo, path, tag
}

func runCmdAndPipeOutput(name string, arg ...string) error {
	fmt.Println(">", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func runCmdAndGetOutput(name string, arg ...string) (string, error) {
	fmt.Println(">", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func findDockerfiles() []string {
	matches, err := filepath.Glob("Dockerfile*")
	if err != nil {
		log.Fatal("error finding Dockerfiles", err)
	}
	return matches
}

func mapDockerfileToRepo(base, path, tag string, dockerfile ...string) map[string]string {
	m := make(map[string]string)
	for _, f := range dockerfile {
		m[f] = generateRepoName(base, path, tag, f)
	}
	return m
}

func generateRepoName(base, path, tag, dockerfile string) string {
	if strings.HasSuffix(registry, "/") {
		base = registry + base
	} else {
		base = registry + "/" + base
	}

	if path != "" {
		base = base + "-" + path
	}
	// grab the suffix from the Dockerfile, if any (e.g. "Dockerfile.foo" => "foo")
	suffix := dockerfile
	if index := strings.LastIndex(suffix, "."); index != -1 {
		suffix = suffix[index+1:]
	} else {
		return base + ":" + tag
	}

	name := base + "-" + suffix + ":" + tag

	// Docker image names can't have more than 2 '/' chars in them ¯\_(ツ)_/¯
	// replace any offending '/' chars w/ '-'
	if strings.Count(name, "/") > 2 {
		nameTokens := strings.SplitN(name, "/", 3)
		nameTokens[2] = strings.Replace(nameTokens[2], "/", "-", -1)
		name = strings.Join(nameTokens, "/")
	}
	return name
}
