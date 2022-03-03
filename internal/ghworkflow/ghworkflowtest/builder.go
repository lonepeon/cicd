package ghworkflowtest

import (
	"html/template"
	"strings"
	"testing"

	"github.com/lonepeon/golib/testutils"
)

type WorkflowFile struct {
	t *testing.T

	actionSetupGoV2 []string
}

func NewWorkflowFile(t *testing.T) WorkflowFile {
	return WorkflowFile{t: t}
}

func (w WorkflowFile) WithActionSetupGoV2(versions ...string) WorkflowFile {
	w.actionSetupGoV2 = versions
	return w
}

func (w WorkflowFile) Build() string {
	var out strings.Builder
	tpl := template.Must(template.New("").Parse(tpls))
	err := tpl.Execute(&out, map[string]interface{}{
		"ActionSetupGoV2": w.actionSetupGoV2,
	})
	testutils.RequireNoError(w.t, err, "can't generate github workflow file")

	return out.String()
}

const tpls = `
name: test

on:
  push:
    branches-ignore:
      - main

jobs:
  {{- range .ActionSetupGoV2}}
  integration-tests:
    runs-on: "ubuntu-latest"
    steps:
      - name: "checkout code"
        uses: actions/checkout@v2
      - name: "setup go version"
        uses: actions/setup-go@v2
        with:
          go-version: {{ . }}
      - name: "assert go version"
        run: make test-go-version
      - name: "run integration tests"
        run: make test-integration
  {{- end }}
`
