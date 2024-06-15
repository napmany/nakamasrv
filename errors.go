package main

import "github.com/heroiclabs/nakama-common/runtime"

const (
	OK                  = 0
	CANCELED            = 1
	UNKNOWN             = 2
	INVALID_ARGUMENT    = 3
	DEADLINE_EXCEEDED   = 4
	NOT_FOUND           = 5
	ALREADY_EXISTS      = 6
	PERMISSION_DENIED   = 7
	RESOURCE_EXHAUSTED  = 8
	FAILED_PRECONDITION = 9
	ABORTED             = 10
	OUT_OF_RANGE        = 11
	UNIMPLEMENTED       = 12
	INTERNAL            = 13
	UNAVAILABLE         = 14
	DATA_LOSS           = 15
	UNAUTHENTICATED     = 16
)

var (
	errInternalError    = runtime.NewError("internal server error", INTERNAL)
	errMarshal          = runtime.NewError("cannot marshal type", INTERNAL)
	errValidationFailed = runtime.NewError("validation failed", INVALID_ARGUMENT)
	errCtxUserIdFound   = runtime.NewError("user ID in context not allowed", INVALID_ARGUMENT)
	errUnmarshal        = runtime.NewError("cannot unmarshal type", INTERNAL)
	errFileNotFound     = runtime.NewError("file not found", INTERNAL)
)
