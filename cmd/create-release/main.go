package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/lonepeon/cicd/internal/github"
	"github.com/lonepeon/cicd/internal/report"
	"github.com/lonepeon/golib/cli"
)

const summary = `%s [-repository <repository> [-asset <path-to-binary>...]] <release-name> <commitish>

The command line is in charge of creating a release and upload all the related
binaries

Arguments

	release-name
	  (required) name of the release and tag (e.g. v0.1.0, 20210313214700,...)
	commitish
	  (required) SHA or branch to use to link to the release and tag

Flags

	-h
	  display this message
	-repository
	  (default=read from GITHUB_REPOSITORY env variable) name of the repository to
	  use. It is written with the following format: "username/project"
	-asset
	  path to the binary to upload. The flag asset can be added multiple times to
	  attach several assets tp the release. If not set, no asset will be uploaded

Environment

	- GITHUB_REPOSITORY_OWNER
	   default username used to authenticate. This variable is set automatically
	   by GitHub action but can be set manually
	- PERSONAL_TOKEN
	   convention used in my projects. The variable contains a token with write
	   access
`

type Flags struct {
	ReleaseName string
	Commitish   string
	Username    string
	Token       string
	Repository  string
	Assets      cli.ArgStrings
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var flags Flags

	cmdline := flag.NewFlagSet("create-release", flag.ExitOnError)
	cmdline.Usage = func() { fmt.Fprintf(cmdline.Output(), summary, cmdline.Name()) }
	cmdline.StringVar(&flags.Repository, "repository", os.Getenv("GITHUB_REPOSITORY"), "")
	cmdline.Var(&flags.Assets, "asset", "")

	if err := parseArguments(cmdline, &flags, os.Args[1:]); err != nil {
		cmdline.Usage()
		return err
	}

	client := github.NewClient(flags.Username, flags.Token)
	releaseID, err := github.CreateRelease(client, flags.Repository, flags.ReleaseName, flags.Commitish)
	if err != nil {
		return fmt.Errorf("can't create release '%s' based on '%s': %v", flags.ReleaseName, flags.Commitish, err)
	}
	report.Success("release created")

	for _, asset := range flags.Assets {
		if err := uploadAsset(client, flags.Repository, releaseID, asset); err != nil {
			return err
		}
	}

	return nil
}

func parseArguments(cmdline *flag.FlagSet, flags *Flags, args []string) error {
	if err := cmdline.Parse(args); err != nil {
		return err
	}

	releaseName := cmdline.Arg(0)
	if releaseName == "" {
		return fmt.Errorf("release name is required")
	}
	flags.ReleaseName = releaseName

	commitish := cmdline.Arg(1)
	if commitish == "" {
		return fmt.Errorf("commit SHA or branch is required")
	}
	flags.Commitish = commitish

	if flags.Repository == "" {
		return fmt.Errorf("repository must be set as a flag or defined using the GITHUB_REPOSITORY environment variable")
	}

	username := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if username == "" {
		return fmt.Errorf("missing environment variable GITHUB_REPOSITORY_OWNER")
	}
	flags.Username = username

	token := os.Getenv("PERSONAL_TOKEN")
	if token == "" {
		return fmt.Errorf("missing environment variable PERSONAL_TOKEN")
	}
	flags.Token = token

	return nil
}

func uploadAsset(client *github.Client, repository string, releaseID github.ReleaseID, asset string) error {
	content, err := ioutil.ReadFile(asset)
	if err != nil {
		return fmt.Errorf("can't read asset (file=%s): %v", asset, err)
	}

	err = github.UploadAsset(client, repository, releaseID, github.Asset{
		Name:        path.Base(asset),
		ContentType: "application/octet-stream",
		Content:     content,
	})
	if err != nil {
		return fmt.Errorf("can't upload asset (file=%s): %v", asset, err)
	}
	report.Success(fmt.Sprintf("asset '%s' attached to the release", asset))

	return nil
}
