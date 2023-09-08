package test

import (
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	reAmt := `^([[:digit:]]+)[[:space:]]*([a-z][a-z0-9]{2,15})$`

	r := regexp.MustCompile(reAmt)

	t.Log(r.FindStringSubmatch("1000   nametoken"))
}
