package asdf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/lonepeon/cicd/internal"
)

var (
	ErrMultipleDeclaration = errors.New("language declared several times")

	spaceBasedRegex = regexp.MustCompile(`\s+`)
)

type ASDF struct {
	versions map[internal.Language]string
}

func Parse(path string) (ASDF, error) {
	f, err := os.Open(path)
	if err != nil {
		return ASDF{}, fmt.Errorf("can't open asdf file: %v", err)
	}
	defer f.Close()

	return ParseFromReader(f)
}

func ParseFromReader(r io.Reader) (ASDF, error) {
	versions := make(map[internal.Language]string)

	var lineNumber int
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		splittedLine := spaceBasedRegex.Split(line, -1)
		if len(splittedLine) < 2 {
			continue
		}

		languageName := splittedLine[0]
		language, err := languageFromName(languageName)
		if err != nil {
			continue
		}

		primaryVersion := splittedLine[1]

		if _, found := versions[language]; found {
			return ASDF{}, fmt.Errorf("%w (line=%d, language=%s)", ErrMultipleDeclaration, lineNumber, languageName)
		}

		versions[language] = primaryVersion
	}

	if err := scanner.Err(); err != nil {
		return ASDF{}, fmt.Errorf("can't scan asdf file: %v", err)
	}

	return ASDF{versions: versions}, nil
}

func (a ASDF) GetVersion(language internal.Language) (string, bool) {
	version, found := a.versions[language]
	return version, found
}

func languageFromName(name string) (internal.Language, error) {
	switch name {
	case "golang":
		return internal.Go, nil
	default:
		return internal.Language(-1), fmt.Errorf("unsupported language '%s'", name)
	}
}
