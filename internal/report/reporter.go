package report

import (
	"io"
	"os"
)

type Reporter struct {
	hasFailed bool

	Stdout io.Writer
	Stderr io.Writer
}

func (r *Reporter) Success(msg string) {
	Fsuccess(r.stdout(), msg)
}

func (r *Reporter) Fail(msg string) {
	r.hasFailed = true
	Ffail(r.stderr(), msg)
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
