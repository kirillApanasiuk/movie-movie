package gateway

import "errors"

var ErrNotFound = errors.New("not found")
var NotSuccesfull = errors.New("non-2xx response")
