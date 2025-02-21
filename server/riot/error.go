package riot

import "errors"

// client errors
var NoLimiterError = errors.New("no rate limit associated with given server")

// riot api errors
var NotFoundError = errors.New("404 not found")
