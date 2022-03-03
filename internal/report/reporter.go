package report

import (
	"fmt"
	"io"
	"os"

	"github.com/lonepeon/cicd/internal/terminal/vt100"
)

type Reporter struct {
	hasFailed bool

	Stdout io.Writer
	Stderr io.Writer
}

func (r *Reporter) Success(msg string) {
	fmt.Fprintln(r.stdout(), vt100.Green("[OK] "), msg)
}

func (r *Reporter) Fail(msg string) {
	r.hasFailed = true
	fmt.Fprintln(r.stderr(), vt100.Red("[ERR]"), msg)
}

func (r *Reporter) HasFailed() bool {
	return r.hasFailed
}

func (r *Reporter) stdout() io.Writer {
	if r.Stdout == nil {
		return os.Stdout
	}

	return r.Stdout
}

func (r *Reporter) stderr() io.Writer {
	if r.Stderr == nil {
		return os.Stderr
	}

	return r.Stderr
}
