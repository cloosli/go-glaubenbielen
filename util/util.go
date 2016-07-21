package util

import (
	"errors"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"os"
	"path/filepath"
	"unicode"
)

func NormalizeText(s string) string {
	isMn := func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	}
	tf := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	name, _, _ := transform.String(tf, s)
	return name
}

func CreatePathTo(s string) error {

	if len(s) == 0 {
		return errors.New("Null path")
	}

	// Ignore the end of path, which is assumed to be a file
	s = filepath.Dir(s)
	s = filepath.Clean(s)

	fmt.Printf("Creating dirs to path %s\n", s)

	// Create all directories up to path
	return os.MkdirAll(s, 0774)
}
