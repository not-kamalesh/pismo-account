package errors

type PismoErrorCode string

const (
	PismoErrorCodeInvalidArgument    PismoErrorCode = "INVALID_ARGUMENT" // bad input
	PismoErrorCodeNotFound           PismoErrorCode = "NOT_FOUND"        // resource missing
	PismoErrorCodeAlreadyExists      PismoErrorCode = "ALREADY_EXISTS"
	PismoErrorCodePermissionDenied   PismoErrorCode = "PERMISSION_DENIED"
	PismoErrorCodeUnauthenticated    PismoErrorCode = "UNAUTHENTICATED"
	PismoErrorCodeFailedPrecondition PismoErrorCode = "FAILED_PRECONDITION"
	PismoErrorCodeInternal           PismoErrorCode = "INTERNAL"
	PismoErrorCodeUnavailable        PismoErrorCode = "UNAVAILABLE"
	PismoErrorCodeDeadlineExceeded   PismoErrorCode = "DEADLINE_EXCEEDED"
)
