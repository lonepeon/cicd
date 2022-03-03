package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/lonepeon/cicd/internal"
	"github.com/lonepeon/cicd/internal/asdf"
	"github.com/lonepeon/cicd/internal/ghworkflow"
	"github.com/lonepeon/cicd/internal/report"
	"github.com/lonepeon/cicd/internal/system"
)

const summary = `%s [-asdf <path> -workflow-folder <folder>] <language>

The command line is in charge of checking if all the part of the system
is configured to use the same version of the language.

The source of truth is defined to be the local ASDF tool-versions file.

Arguments

	language
	  (required) one of the supported language: go

Flags

	-h
	  display this message
	-asdf
	  (default=.tool-versions) path to the local ASDF tool-version file
	-gh-workflows
	  (default=.github/workflows) path to the GitHub workflow folder.
	  The application will load all YAML files in the folder
`

type Flags struct {
	ASDFFilePath             string
	GitHubWorkflowFolderPath string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var flags Flags

	cmdline := flag.NewFlagSet("assert-versions", flag.ExitOnError)
	cmdline.Usage = func() { fmt.Fprintf(cmdline.Output(), summary, cmdline.Name()) }
	cmdline.StringVar(&flags.ASDFFilePath, "asdf", ".tool-versions", "")
	cmdline.StringVar(&flags.GitHubWorkflowFolderPath, "gh-workflows", ".github/workflows", "")

	if err := cmdline.Parse(os.Args[1:]); err != nil {
		return err
	}

	language, version, err := expectedLanguageVersion(flags.ASDFFilePath, cmdline.Arg(0))
	if err != nil {
		return err
	}

	reporter := report.NewVersionReporter(version)
	reporter.AssertVersion(report.VersionEntry{Version: version, Kind: "ASDF"})
	assertSystemVersion(reporter, language)
	if err := assertGitHubWorkflowVersions(reporter, flags.GitHubWorkflowFolderPath, language); err != nil {
		return err
	}

	if reporter.HasFailed() {
		return fmt.Errorf("one or more version doesn't match the expectation")
	}

	return nil
}

func assertGitHubWorkflowVersions(reporter *report.VersionReporter, workflowDir string, language internal.Language) error {
	return filepath.WalkDir(workflowDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if extension != ".yml" && extension != ".yaml" {
			return nil
		}

		workflow, err := ghworkflow.Parse(path)
		if err != nil {
			return err
		}

		for _, match := range workflow.GetVersions(language) {
			location := fmt.Sprintf("%s:%d", path, match.Line)
			reporter.AssertVersion(report.VersionEntry{
				Version:  match.Version,
				Kind:     "GitHub Workflow",
				Location: location,
			})
		}

		return nil
	})
}

func expectedLanguageVersion(asdfFilePath string, languageName string) (internal.Language, string, error) {
	var language internal.Language
	switch languageName {
	case "go":
		language = internal.Go
	default:
		return internal.Language(-1), "", fmt.Errorf("unsupported language '%s'", languageName)
	}

	toolVersion, err := asdf.Parse(asdfFilePath)
	if err != nil {
		return internal.Language(-1), "", err
	}

	expectedVersion, found := toolVersion.GetVersion(language)
	if !found {
		return internal.Language(-1), "", fmt.Errorf("no version defined for '%s' in '%s'", languageName, asdfFilePath)
	}

	return language, expectedVersion, nil
}

func assertSystemVersion(reporter *report.VersionReporter, language internal.Language) {
	version, err := system.GetVersion(language)
	if err != nil {
		reporter.AssertVersion(report.VersionEntry{Version: version, Kind: "System"})
		return
	}

	reporter.AssertVersion(report.VersionEntry{Version: version, Kind: "System"})
}
