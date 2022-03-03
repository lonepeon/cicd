package report

import (
	"fmt"
	"strings"
)

type VersionEntry struct {
	Version  string
	Kind     string
	Location string
}

type VersionReporter struct {
	expectedVersion string

	Reporter *Reporter
}

func NewVersionReporter(expectedVersion string) *VersionReporter {
	return &VersionReporter{
		expectedVersion: expectedVersion,
		Reporter:        &Reporter{},
	}
}

func (r *VersionReporter) AssertVersion(entry VersionEntry) {
	logger := r.Reporter.Success
	if entry.Version != r.expectedVersion {
		logger = r.Reporter.Fail
	}

	logger(r.buildEntryDescription(entry))
}

func (r *VersionReporter) HasFailed() bool {
	return r.Reporter.HasFailed()
}

func (r *VersionReporter) buildEntryDescription(entry VersionEntry) string {
	attrs := []string{fmt.Sprintf("version=%s", entry.Version)}

	if entry.Location != "" {
		attrs = append(attrs, fmt.Sprintf("location=%s", entry.Location))
	}

	return fmt.Sprintf("check %s version (%s)", entry.Kind, strings.Join(attrs, ", "))
}
