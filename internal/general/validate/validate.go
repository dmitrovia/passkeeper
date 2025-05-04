// Package random validate functions
// working with random validation.
package validate

import (
	"fmt"
	"regexp"
)

// IsMatchesTemplate - checks
// for regular expression matches.
func IsMatchesTemplate(
	addr string,
	pattern string,
) (bool, error) {
	res, err := matchString(pattern, addr)
	if err != nil {
		return false, err
	}

	return res, err
}

func matchString(pattern string, s string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err == nil {
		return re.MatchString(s), nil
	}

	return false, fmt.Errorf("MatchString: %w", err)
}

func IsValidSecretID(sid string) bool {
	pattern := "^[0-9a-zA-Z/ ]{1,40}$"

	res, err := IsMatchesTemplate(
		sid, pattern)
	if err != nil || !res {
		return false
	}

	return true
}

func IsValidLogin(login string) bool {
	pattern := "^[0-9a-zA-Z/ ]{1,40}$"

	res, err := IsMatchesTemplate(
		login, pattern)
	if err != nil || !res {
		return false
	}

	return true
}

func IsValidPass(pass string) bool {
	pattern := "^.{1,40}$"

	res, err := IsMatchesTemplate(
		pass, pattern)
	if err != nil || !res {
		return false
	}

	return true
}

func IsValidOrigFileName(fn string) bool {
	pattern := "^.{1,207}$"

	res, err := IsMatchesTemplate(
		fn, pattern)
	if err != nil || !res {
		return false
	}

	return true
}

func IsValidHash(fn string) bool {
	pattern := "^.{1,200}$"

	res, err := IsMatchesTemplate(
		fn, pattern)
	if err != nil || !res {
		return false
	}

	return true
}

func IsValidFileName(fn string) bool {
	pattern := "^.{1,200}$"

	res, err := IsMatchesTemplate(
		fn, pattern)
	if err != nil || !res {
		return false
	}

	return true
}
