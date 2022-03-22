package controllers

import "github.com/giantswarm/microerror"

var invalidObjectError = &microerror.Error{
	Kind: "invalidObjectError",
}

// IsInvalidObject asserts invalidObjectError.
func IsInvalidObject(err error) bool {
	return microerror.Cause(err) == invalidObjectError
}
