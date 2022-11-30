package common

import (
	"regexp"

	"github.com/pkg/errors"
)

var dmReg *regexp.Regexp

func init() {
	dmReg = regexp.MustCompile("(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]")
}

// FilterErr returns a new error that ip-address and domain-name replaced with '#.#.#.#'
func FilterErr(err error) error {
	if err == nil {
		return nil
	}
	errMsg := err.Error()
	return errors.New(dmReg.ReplaceAllLiteralString(errMsg, "#.#.#.#"))
}
