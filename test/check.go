package test

import "testing"

// Check fails a test if `err` is not nil.
func Check(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

// Assert fails a test if `ok` is not true.
func Assert(ok bool, message string, t *testing.T) {
	if !ok {
		t.Error(message)
	}
}
