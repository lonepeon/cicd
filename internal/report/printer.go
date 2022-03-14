package report

import (
	"fmt"
	"io"
	"os"

	"github.com/lonepeon/cicd/internal/terminal/vt100"
)

func Success(msg string) {
	Fsuccess(os.Stdout, msg)
}

func Fsuccess(w io.Writer, msg string) {
	fmt.Fprintln(w, vt100.Green("[OK] "), msg)
}

func Fail(msg string) {
	Ffail(os.Stderr, msg)
}

func Ffail(w io.Writer, msg string) {
	fmt.Fprintln(w, vt100.Red("[ERR]"), msg)
}
